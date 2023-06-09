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
      chrome.external.invoke(input);
    } else {
      webkit.messageHandlers.external.postMessage(input);
    }
  };

  // desktop/drag.js
  var shouldDrag = false;
  function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--wails-draggable");
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
  function onMouseDown(e) {
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
    document.body.style.cursor = window.wails.previousCursor || "auto";
    shouldDrag = false;
  }
  function onMouseMove(e) {
    if (shouldDrag) {
      shouldDrag = false;
      let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
      if (mousePressed > 0) {
        window.wails.previousCursor = document.body.style.cursor;
        document.body.style.cursor = "grab";
        invoke("drag");
        return;
      }
    }
  }

  // desktop/main.js
  window.wails = {
    ...newRuntime(null)
  };
  window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchWailsEvent,
    callCallback,
    callErrorCallback
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvaW52b2tlLmpzIiwgImRlc2t0b3AvZHJhZy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjbGlwYm9hcmRcIik7XG5cbi8qKlxuICogU2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0KSB7XG4gICAgdm9pZCBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBUZXh0KCkge1xuICAgIHJldHVybiBjYWxsKFwiVGV4dFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuY29uc3QgcnVudGltZVVSTCA9IHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XG5cbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsKG1ldGhvZCwgd2luZG93TmFtZSwgYXJncykge1xuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xuICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwibWV0aG9kXCIsIG1ldGhvZCk7XG4gICAgaWYgKGFyZ3MpIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcbiAgICB9XG4gICAgbGV0IGZldGNoT3B0aW9ucyA9IHtcbiAgICAgICAgaGVhZGVyczoge30sXG4gICAgfTtcbiAgICBpZiAod2luZG93TmFtZSkge1xuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xuICAgIH1cbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3QsIHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImFwcGxpY2F0aW9uXCIpO1xuXG4vKipcbiAqIEhpZGUgdGhlIGFwcGxpY2F0aW9uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xuICAgIHZvaWQgY2FsbChcIkhpZGVcIik7XG59XG5cbi8qKlxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XG4gICAgdm9pZCBjYWxsKFwiU2hvd1wiKTtcbn1cblxuXG4vKipcbiAqIFF1aXQgdGhlIGFwcGxpY2F0aW9uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHZvaWQgY2FsbChcIlF1aXRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJsb2dcIik7XG5cbi8qKlxuICogTG9ncyBhIG1lc3NhZ2UuXG4gKiBAcGFyYW0ge21lc3NhZ2V9IE1lc3NhZ2UgdG8gbG9nXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2cobWVzc2FnZSkge1xuICAgIHJldHVybiBjYWxsKFwiTG9nXCIsIG1lc3NhZ2UpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXG4gKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcInNjcmVlbnNcIik7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldEFsbFwiKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0UHJpbWFyeVwiKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxuICogQGNvbnN0cnVjdG9yXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0Q3VycmVudFwiKTtcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNhbGxcIik7XG5cbmxldCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gY2FsbENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoXCJDYWxsXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIENhbGwgYSBwbHVnaW4gbWV0aG9kXG4gKiBAcGFyYW0ge3N0cmluZ30gcGx1Z2luTmFtZSAtIG5hbWUgb2YgdGhlIHBsdWdpblxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBuYW1lIG9mIHRoZSBtZXRob2RcbiAqIEBwYXJhbSB7Li4uYW55fSBhcmdzIC0gYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gLSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzdWx0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhcIkNhbGxcIiwge1xuICAgICAgICBwYWNrYWdlTmFtZTogXCJ3YWlscy1wbHVnaW5zXCIsXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXG4gICAgICAgIG1ldGhvZE5hbWU6IG1ldGhvZE5hbWUsXG4gICAgICAgIGFyZ3M6IGFyZ3MsXG4gICAgfSk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5TaXplfSBTaXplXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlBvc2l0aW9ufSBQb3NpdGlvblxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5TY3JlZW59IFNjcmVlblxuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5leHBvcnQgZnVuY3Rpb24gbmV3V2luZG93KHdpbmRvd05hbWUpIHtcbiAgICBsZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJ3aW5kb3dcIiwgd2luZG93TmFtZSk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBSZWxvYWRBcHA6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFNldFN5c3RlbURlZmF1bHRUaGVtZTogKCkgPT4gY2FsbCgnV0FTRFQnKSxcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxuICAgICAgICAvLyBTZXREYXJrVGhlbWU6ICgpID0+IGNhbGwoJ1dBRFQnKSxcbiAgICAgICAgLy8gSXNGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxuICAgICAgICAvLyBJc01pbmltaXplZDogKCkgPT4gY2FsbCgnV0lNTicpLFxuICAgICAgICAvLyBJc1dpbmRvd2VkOiAoKSA9PiBjYWxsKCdXSUYnKSxcblxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cuXG4gICAgICAgICAqL1xuICAgICAgICBDZW50ZXI6ICgpID0+IHZvaWQgY2FsbCgnQ2VudGVyJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHRpdGxlLlxuICAgICAgICAgKiBAcGFyYW0gdGl0bGVcbiAgICAgICAgICovXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IHZvaWQgY2FsbCgnU2V0VGl0bGUnLCB7dGl0bGV9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWFrZXMgdGhlIHdpbmRvdyBmdWxsc2NyZWVuLlxuICAgICAgICAgKi9cbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gdm9pZCBjYWxsKCdGdWxsc2NyZWVuJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFVuZnVsbHNjcmVlbiB0aGUgd2luZG93LlxuICAgICAgICAgKi9cbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBzaXplLlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGggVGhlIHdpbmRvdyB3aWR0aFxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IFRoZSB3aW5kb3cgaGVpZ2h0XG4gICAgICAgICAqL1xuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgc2l6ZS5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8U2l6ZT59IFRoZSB3aW5kb3cgc2l6ZVxuICAgICAgICAgKi9cbiAgICAgICAgU2l6ZTogKCkgPT4geyByZXR1cm4gY2FsbCgnU2l6ZScpOyB9LFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBtYXhpbXVtIHNpemUuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gICAgICAgICAqL1xuICAgICAgICBTZXRNYXhTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gdm9pZCBjYWxsKCdTZXRNYXhTaXplJywge3dpZHRoLGhlaWdodH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBtaW5pbXVtIHNpemUuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gICAgICAgICAqL1xuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gdm9pZCBjYWxsKCdTZXRNaW5TaXplJywge3dpZHRoLGhlaWdodH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXG4gICAgICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gb25Ub3AgV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBhbHdheXMgb24gdG9wXG4gICAgICAgICAqL1xuICAgICAgICBTZXRBbHdheXNPblRvcDogKG9uVG9wKSA9PiB2b2lkIGNhbGwoJ1NldEFsd2F5c09uVG9wJywge2Fsd2F5c09uVG9wOm9uVG9wfSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHBvc2l0aW9uLlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0geFxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0geVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0UG9zaXRpb246ICh4LCB5KSA9PiBjYWxsKCdTZXRQb3NpdGlvbicsIHt4LHl9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgcG9zaXRpb24uXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFBvc2l0aW9uPn0gVGhlIHdpbmRvdyBwb3NpdGlvblxuICAgICAgICAgKi9cbiAgICAgICAgUG9zaXRpb246ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1Bvc2l0aW9uJyk7IH0sXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEdldCB0aGUgc2NyZWVuIHRoZSB3aW5kb3cgaXMgb24uXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gICAgICAgICAqL1xuICAgICAgICBTY3JlZW46ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1NjcmVlbicpOyB9LFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBIaWRlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIEhpZGU6ICgpID0+IHZvaWQgY2FsbCgnSGlkZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBNYXhpbWlzZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBNYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdNYXhpbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTaG93IHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFNob3c6ICgpID0+IHZvaWQgY2FsbCgnU2hvdycpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBDbG9zZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBDbG9zZTogKCkgPT4gdm9pZCBjYWxsKCdDbG9zZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBUb2dnbGUgdGhlIHdpbmRvdyBtYXhpbWlzZSBzdGF0ZVxuICAgICAgICAgKi9cbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVG9nZ2xlTWF4aW1pc2UnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVW5tYXhpbWlzZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBVbk1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ1VuTWF4aW1pc2UnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWluaW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnTWluaW1pc2UnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVW5taW5pbWlzZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBVbk1pbmltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ1VuTWluaW1pc2UnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gciAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGcgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBiIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gYSAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICovXG4gICAgICAgIFNldEJhY2tncm91bmRDb2xvdXI6IChyLCBnLCBiLCBhKSA9PiB2b2lkIGNhbGwoJ1NldEJhY2tncm91bmRDb2xvdXInLCB7ciwgZywgYiwgYX0pLFxuICAgIH07XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuV2FpbHNFdmVudH0gV2FpbHNFdmVudFxuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJldmVudHNcIik7XG5cbi8qKlxuICogVGhlIExpc3RlbmVyIGNsYXNzIGRlZmluZXMgYSBsaXN0ZW5lciEgOi0pXG4gKlxuICogQGNsYXNzIExpc3RlbmVyXG4gKi9cbmNsYXNzIExpc3RlbmVyIHtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIExpc3RlbmVyLlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxuICAgICAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5cbi8qKlxuICogV2FpbHNFdmVudCBkZWZpbmVzIGEgY3VzdG9tIGV2ZW50LiBJdCBpcyBwYXNzZWQgdG8gZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBjbGFzcyBXYWlsc0V2ZW50XG4gKiBAcHJvcGVydHkge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gKiBAcHJvcGVydHkge2FueX0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxuICovXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBXYWlsc0V2ZW50LlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gTmFtZSBvZiB0aGUgZXZlbnRcbiAgICAgKiBAcGFyYW0ge2FueT1udWxsfSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XG4gICAgICogQG1lbWJlcm9mIFdhaWxzRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhID0gbnVsbCkge1xuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuLyoqXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uXG4gKlxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcbiAqL1xuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBkaXNwYXRjaGVzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQpIHtcbiAgICBjb25zb2xlLmxvZyhcImRpc3BhdGNoaW5nIGV2ZW50OiBcIiwge2V2ZW50fSk7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBbXTtcbiAgICAgICAgbGlzdGVuZXJzLmZvckVhY2gobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHtcbiAgICAgICAgICAgICAgICB0b1JlbW92ZS5wdXNoKGxpc3RlbmVyKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgICAgIC8vIHJlbW92ZSBsaXN0ZW5lcnNcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xuICAgICAgICAgICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCBsaXN0ZW5lcnMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcnMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxuICpcbiBbdjMgQ0hBTkdFXSBPZmYgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSk7XG59XG5cbi8qKlxuICogT2ZmQWxsIHVucmVnaXN0ZXJzIGFsbCBsaXN0ZW5lcnNcbiAqIFt2MyBDSEFOR0VdIE9mZkFsbCBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkge1xuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XG59XG5cbi8qKlxuICogRW1pdCBhbiBldmVudFxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCBUaGUgZXZlbnQgdG8gZW1pdFxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkge1xuICAgIHZvaWQgY2FsbChcIkVtaXRcIiwgZXZlbnQpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5NZXNzYWdlRGlhbG9nT3B0aW9uc30gTWVzc2FnZURpYWxvZ09wdGlvbnNcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5PcGVuRGlhbG9nT3B0aW9uc30gT3BlbkRpYWxvZ09wdGlvbnNcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5TYXZlRGlhbG9nT3B0aW9uc30gU2F2ZURpYWxvZ09wdGlvbnNcbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJkaWFsb2dcIik7XG5cbmxldCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwKCk7XG5cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0NhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcbiAgICAgICAgfVxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gZGlhbG9nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5jYXRjaCgoZXJyb3IpID0+IHtcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cblxuLyoqXG4gKiBTaG93cyBhbiBJbmZvIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkluZm9cIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYW4gV2FybmluZyBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdhcm5pbmcob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJXYXJuaW5nXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGFuIEVycm9yIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gRXJyb3Iob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJFcnJvclwiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhIFF1ZXN0aW9uIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJRdWVzdGlvblwiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhbiBPcGVuIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtPcGVuRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nW118c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZSBvciBhbiBhcnJheSBvZiBzZWxlY3RlZCBmaWxlcyBpZiBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbiBpcyB0cnVlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiT3BlbkZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYSBTYXZlIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtPcGVuRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZS4gQSBibGFuayBzdHJpbmcgaXMgcmV0dXJuZWQgaWYgbm8gZmlsZSB3YXMgc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlNhdmVGaWxlXCIsIG9wdGlvbnMpO1xufVxuXG4iLCAiaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNvbnRleHRtZW51XCIpO1xuXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQsIHgsIHksIGRhdGEpIHtcbiAgICByZXR1cm4gY2FsbChcIk9wZW5Db250ZXh0TWVudVwiLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGVuYWJsZUNvbnRleHRNZW51cyhlbmFibGVkKSB7XG4gICAgaWYgKGVuYWJsZWQpIHtcbiAgICAgICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50KSB7XG4gICAgcHJvY2Vzc0NvbnRleHRNZW51KGV2ZW50LnRhcmdldCwgZXZlbnQpO1xufVxuXG5mdW5jdGlvbiBwcm9jZXNzQ29udGV4dE1lbnUoZWxlbWVudCwgZXZlbnQpIHtcbiAgICBsZXQgaWQgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS1jb250ZXh0bWVudScpO1xuICAgIGlmIChpZCkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoaWQsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51LWRhdGEnKSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgbGV0IHBhcmVudCA9IGVsZW1lbnQucGFyZW50RWxlbWVudDtcbiAgICAgICAgaWYgKHBhcmVudCkge1xuICAgICAgICAgICAgcHJvY2Vzc0NvbnRleHRNZW51KHBhcmVudCwgZXZlbnQpO1xuICAgICAgICB9XG4gICAgfVxufVxuIiwgIlxuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcblxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZSwgZGF0YT1udWxsKSB7XG4gICAgbGV0IGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKTtcbiAgICBFbWl0KGV2ZW50KTtcbn1cblxuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtZXZlbnRdJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICB9O1xuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG5cbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2QobWV0aG9kKSB7XG4gICAgaWYgKHdhaWxzLldpbmRvd1ttZXRob2RdID09PSB1bmRlZmluZWQpIHtcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xuICAgIH1cbiAgICB3YWlscy5XaW5kb3dbbWV0aG9kXSgpO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtd2luZG93XScpO1xuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVsb2FkV01MKCkge1xuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLy8gZGVmaW5lZCBpbiB0aGUgVGFza2ZpbGVcbmV4cG9ydCBsZXQgaW52b2tlID0gZnVuY3Rpb24oaW5wdXQpIHtcbiAgICBpZihXSU5ET1dTKSB7XG4gICAgICAgIGNocm9tZS5leHRlcm5hbC5pbnZva2UoaW5wdXQpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2UoaW5wdXQpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge2ludm9rZX0gZnJvbSBcIi4vaW52b2tlXCI7XG5cbmxldCBzaG91bGREcmFnID0gZmFsc2U7XG5cbmV4cG9ydCBmdW5jdGlvbiBkcmFnVGVzdChlKSB7XG4gICAgbGV0IHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13YWlscy1kcmFnZ2FibGVcIik7XG4gICAgaWYgKHZhbCkge1xuICAgICAgICB2YWwgPSB2YWwudHJpbSgpO1xuICAgIH1cblxuICAgIGlmICh2YWwgIT09IFwiZHJhZ1wiKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICAvLyBPbmx5IHByb2Nlc3MgdGhlIHByaW1hcnkgYnV0dG9uXG4gICAgaWYgKGUuYnV0dG9ucyAhPT0gMSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgcmV0dXJuIGUuZGV0YWlsID09PSAxO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc2V0dXBEcmFnKCkge1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCBvbk1vdXNlRG93bik7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIG9uTW91c2VNb3ZlKTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIG9uTW91c2VVcCk7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VEb3duKGUpIHtcbiAgICBpZiAoZHJhZ1Rlc3QoZSkpIHtcbiAgICAgICAgLy8gSWdub3JlIGRyYWcgb24gc2Nyb2xsYmFyc1xuICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgc2hvdWxkRHJhZyA9IHRydWU7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gb25Nb3VzZVVwKGUpIHtcbiAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IHdpbmRvdy53YWlscy5wcmV2aW91c0N1cnNvciB8fCAnYXV0byc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShlKSB7XG4gICAgaWYgKHNob3VsZERyYWcpIHtcbiAgICAgICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgICAgICBpZiAobW91c2VQcmVzc2VkID4gMCkge1xuICAgICAgICAgICAgd2luZG93LndhaWxzLnByZXZpb3VzQ3Vyc29yID0gZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3I7XG4gICAgICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9ICdncmFiJztcbiAgICAgICAgICAgIGludm9rZShcImRyYWdcIik7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tICcuL2NsaXBib2FyZCc7XG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tICcuL2FwcGxpY2F0aW9uJztcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gJy4vc2NyZWVucyc7XG5pbXBvcnQge1BsdWdpbiwgQ2FsbCwgY2FsbEVycm9yQ2FsbGJhY2ssIGNhbGxDYWxsYmFja30gZnJvbSBcIi4vY2FsbHNcIjtcbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCB7ZGlzcGF0Y2hXYWlsc0V2ZW50LCBFbWl0LCBPZmYsIE9mZkFsbCwgT24sIE9uY2UsIE9uTXVsdGlwbGV9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtkaWFsb2dDYWxsYmFjaywgZGlhbG9nRXJyb3JDYWxsYmFjaywgRXJyb3IsIEluZm8sIE9wZW5GaWxlLCBRdWVzdGlvbiwgU2F2ZUZpbGUsIFdhcm5pbmcsfSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5pbXBvcnQge2VuYWJsZUNvbnRleHRNZW51c30gZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCB7cmVsb2FkV01MfSBmcm9tIFwiLi93bWxcIjtcbmltcG9ydCB7c2V0dXBEcmFnfSBmcm9tIFwiLi9kcmFnXCI7XG5cbndpbmRvdy53YWlscyA9IHtcbiAgICAuLi5uZXdSdW50aW1lKG51bGwpLFxufTtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cuX3dhaWxzID0ge1xuICAgIGRpYWxvZ0NhbGxiYWNrLFxuICAgIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssXG4gICAgZGlzcGF0Y2hXYWlsc0V2ZW50LFxuICAgIGNhbGxDYWxsYmFjayxcbiAgICBjYWxsRXJyb3JDYWxsYmFjayxcbn07XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lKHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4ge1xuICAgICAgICBDbGlwYm9hcmQ6IHtcbiAgICAgICAgICAgIC4uLkNsaXBib2FyZFxuICAgICAgICB9LFxuICAgICAgICBBcHBsaWNhdGlvbjoge1xuICAgICAgICAgICAgLi4uQXBwbGljYXRpb24sXG4gICAgICAgICAgICBHZXRXaW5kb3dCeU5hbWUod2luZG93TmFtZSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBuZXdSdW50aW1lKHdpbmRvd05hbWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9LFxuICAgICAgICBMb2csXG4gICAgICAgIFNjcmVlbnMsXG4gICAgICAgIENhbGwsXG4gICAgICAgIFBsdWdpbixcbiAgICAgICAgV01MOiB7XG4gICAgICAgICAgICBSZWxvYWQ6IHJlbG9hZFdNTCxcbiAgICAgICAgfSxcbiAgICAgICAgRGlhbG9nOiB7XG4gICAgICAgICAgICBJbmZvLFxuICAgICAgICAgICAgV2FybmluZyxcbiAgICAgICAgICAgIEVycm9yLFxuICAgICAgICAgICAgUXVlc3Rpb24sXG4gICAgICAgICAgICBPcGVuRmlsZSxcbiAgICAgICAgICAgIFNhdmVGaWxlLFxuICAgICAgICB9LFxuICAgICAgICBFdmVudHM6IHtcbiAgICAgICAgICAgIEVtaXQsXG4gICAgICAgICAgICBPbixcbiAgICAgICAgICAgIE9uY2UsXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxuICAgICAgICAgICAgT2ZmLFxuICAgICAgICAgICAgT2ZmQWxsLFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyh3aW5kb3dOYW1lKSxcbiAgICB9O1xufVxuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbmVuYWJsZUNvbnRleHRNZW51cyh0cnVlKTtcblxuc2V0dXBEcmFnKCk7XG5cbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoXCJET01Db250ZW50TG9hZGVkXCIsIGZ1bmN0aW9uKGV2ZW50KSB7XG4gICAgcmVsb2FkV01MKCk7XG59KTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ1lBLE1BQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUU1QyxXQUFTLFlBQVksUUFBUSxZQUFZLE1BQU07QUFDM0MsUUFBSSxNQUFNLElBQUksSUFBSSxVQUFVO0FBQzVCLFFBQUksYUFBYSxPQUFPLFVBQVUsTUFBTTtBQUN4QyxRQUFJLE1BQU07QUFDTixVQUFJLGFBQWEsT0FBTyxRQUFRLEtBQUssVUFBVSxJQUFJLENBQUM7QUFBQSxJQUN4RDtBQUNBLFFBQUksZUFBZTtBQUFBLE1BQ2YsU0FBUyxDQUFDO0FBQUEsSUFDZDtBQUNBLFFBQUksWUFBWTtBQUNa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QUNKO0FBRU8sV0FBUyxXQUFXLFlBQVk7QUFDbkMsV0FBTztBQUFBLE1BQ0gsV0FBVztBQUFBLFFBQ1AsR0FBRztBQUFBLE1BQ1A7QUFBQSxNQUNBLGFBQWE7QUFBQSxRQUNULEdBQUc7QUFBQSxRQUNILGdCQUFnQkMsYUFBWTtBQUN4QixpQkFBTyxXQUFXQSxXQUFVO0FBQUEsUUFDaEM7QUFBQSxNQUNKO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0EsS0FBSztBQUFBLFFBQ0QsUUFBUTtBQUFBLE1BQ1o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0EsT0FBQUM7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUSxVQUFVLFVBQVU7QUFBQSxJQUNoQztBQUFBLEVBQ0o7QUFFQSxNQUFJLE1BQU87QUFDUCxZQUFRLElBQUksaUNBQWlDO0FBQUEsRUFDakQ7QUFFQSxxQkFBbUIsSUFBSTtBQUV2QixZQUFVO0FBRVYsV0FBUyxpQkFBaUIsb0JBQW9CLFNBQVMsT0FBTztBQUMxRCxjQUFVO0FBQUEsRUFDZCxDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJldmVudE5hbWUiLCAiY2FsbCIsICJnZW5lcmF0ZUlEIiwgIkVycm9yIiwgImNhbGwiLCAid2luZG93TmFtZSIsICJFcnJvciJdCn0K
