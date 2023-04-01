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
  function runtimeCall(method, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("method", method);
    if (args) {
      url.searchParams.append("args", JSON.stringify(args));
    }
    return new Promise((resolve, reject) => {
      fetch(url).then((response) => {
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
  function newRuntimeCaller(object, id) {
    if (!id || id === -1) {
      return function(method, args) {
        return runtimeCall(object + "." + method, args);
      };
    }
    return function(method, args) {
      args = args || {};
      args.windowID = id;
      return runtimeCall(object + "." + method, args);
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
  function newWindow(id) {
    let call9 = newRuntimeCaller("window", id);
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
  function dispatchCustomEvent(event) {
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
  function sendEvent(event) {
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

  // desktop/main.js
  window.wails = {
    ...newRuntime(-1)
  };
  window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchCustomEvent,
    callCallback,
    callErrorCallback
  };
  function newRuntime(id) {
    return {
      Clipboard: {
        ...clipboard_exports
      },
      Application: {
        ...application_exports
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
      Window: newWindow(id)
    };
  }
  if (true) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
  }
  enableContextMenus(true);
  document.addEventListener("DOMContentLoaded", function(event) {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuLyoqXG4gKiBTZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICB2b2lkIGNhbGwoXCJTZXRUZXh0XCIsIHt0ZXh0fSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBpZihhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybClcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3QsIGlkKSB7XG4gICAgaWYgKCFpZCB8fCBpZCA9PT0gLTEpIHtcbiAgICAgICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgYXJncyk7XG4gICAgICAgIH07XG4gICAgfVxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzKSB7XG4gICAgICAgIGFyZ3MgPSBhcmdzIHx8IHt9O1xuICAgICAgICBhcmdzLndpbmRvd0lEID0gaWQ7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgYXJncyk7XG4gICAgfTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImFwcGxpY2F0aW9uXCIpO1xuXG4vKipcbiAqIEhpZGUgdGhlIGFwcGxpY2F0aW9uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xuICAgIHZvaWQgY2FsbChcIkhpZGVcIik7XG59XG5cbi8qKlxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XG4gICAgdm9pZCBjYWxsKFwiU2hvd1wiKTtcbn1cblxuXG4vKipcbiAqIFF1aXQgdGhlIGFwcGxpY2F0aW9uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHZvaWQgY2FsbChcIlF1aXRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJsb2dcIik7XG5cbi8qKlxuICogTG9ncyBhIG1lc3NhZ2UuXG4gKiBAcGFyYW0ge21lc3NhZ2V9IE1lc3NhZ2UgdG8gbG9nXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2cobWVzc2FnZSkge1xuICAgIHJldHVybiBjYWxsKFwiTG9nXCIsIG1lc3NhZ2UpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXG4gKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcInNjcmVlbnNcIik7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldEFsbFwiKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0UHJpbWFyeVwiKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxuICogQGNvbnN0cnVjdG9yXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0Q3VycmVudFwiKTtcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNhbGxcIik7XG5cbmxldCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gY2FsbENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoXCJDYWxsXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIENhbGwgYSBwbHVnaW4gbWV0aG9kXG4gKiBAcGFyYW0ge3N0cmluZ30gcGx1Z2luTmFtZSAtIG5hbWUgb2YgdGhlIHBsdWdpblxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBuYW1lIG9mIHRoZSBtZXRob2RcbiAqIEBwYXJhbSB7Li4uYW55fSBhcmdzIC0gYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gLSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzdWx0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhcIkNhbGxcIiwge1xuICAgICAgICBwYWNrYWdlTmFtZTogXCJ3YWlscy1wbHVnaW5zXCIsXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXG4gICAgICAgIG1ldGhvZE5hbWU6IG1ldGhvZE5hbWUsXG4gICAgICAgIGFyZ3M6IGFyZ3MsXG4gICAgfSk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5TaXplfSBTaXplXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlBvc2l0aW9ufSBQb3NpdGlvblxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5TY3JlZW59IFNjcmVlblxuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5leHBvcnQgZnVuY3Rpb24gbmV3V2luZG93KGlkKSB7XG4gICAgbGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwid2luZG93XCIsIGlkKTtcbiAgICByZXR1cm4ge1xuICAgICAgICAvLyBSZWxvYWQ6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFJlbG9hZEFwcDogKCkgPT4gY2FsbCgnV1InKSxcbiAgICAgICAgLy8gU2V0U3lzdGVtRGVmYXVsdFRoZW1lOiAoKSA9PiBjYWxsKCdXQVNEVCcpLFxuICAgICAgICAvLyBTZXRMaWdodFRoZW1lOiAoKSA9PiBjYWxsKCdXQUxUJyksXG4gICAgICAgIC8vIFNldERhcmtUaGVtZTogKCkgPT4gY2FsbCgnV0FEVCcpLFxuICAgICAgICAvLyBJc0Z1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1dJRicpLFxuICAgICAgICAvLyBJc01heGltaXplZDogKCkgPT4gY2FsbCgnV0lNJyksXG4gICAgICAgIC8vIElzTWluaW1pemVkOiAoKSA9PiBjYWxsKCdXSU1OJyksXG4gICAgICAgIC8vIElzV2luZG93ZWQ6ICgpID0+IGNhbGwoJ1dJRicpLFxuXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdy5cbiAgICAgICAgICovXG4gICAgICAgIENlbnRlcjogKCkgPT4gdm9pZCBjYWxsKCdDZW50ZXInKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgdGl0bGUuXG4gICAgICAgICAqIEBwYXJhbSB0aXRsZVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gdm9pZCBjYWxsKCdTZXRUaXRsZScsIHt0aXRsZX0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBNYWtlcyB0aGUgd2luZG93IGZ1bGxzY3JlZW4uXG4gICAgICAgICAqL1xuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ0Z1bGxzY3JlZW4nKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVW5mdWxsc2NyZWVuIHRoZSB3aW5kb3cuXG4gICAgICAgICAqL1xuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IHZvaWQgY2FsbCgnVW5GdWxsc2NyZWVuJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHNpemUuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCBUaGUgd2luZG93IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgVGhlIHdpbmRvdyBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRTaXplJywge3dpZHRoLGhlaWdodH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBzaXplLlxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxTaXplPn0gVGhlIHdpbmRvdyBzaXplXG4gICAgICAgICAqL1xuICAgICAgICBTaXplOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTaXplJyk7IH0sXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1heGltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoJ1NldE1heFNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1pbmltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1pblNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoJ1NldE1pblNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB3aW5kb3cgdG8gYmUgYWx3YXlzIG9uIHRvcC5cbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSBvblRvcCBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGFsd2F5cyBvbiB0b3BcbiAgICAgICAgICovXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAob25Ub3ApID0+IHZvaWQgY2FsbCgnU2V0QWx3YXlzT25Ub3AnLCB7YWx3YXlzT25Ub3A6b25Ub3B9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgcG9zaXRpb24uXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB5XG4gICAgICAgICAqL1xuICAgICAgICBTZXRQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoJ1NldFBvc2l0aW9uJywge3gseX0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBwb3NpdGlvbi5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8UG9zaXRpb24+fSBUaGUgd2luZG93IHBvc2l0aW9uXG4gICAgICAgICAqL1xuICAgICAgICBQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUG9zaXRpb24nKTsgfSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBvbi5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cbiAgICAgICAgICovXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJyk7IH0sXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEhpZGUgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgSGlkZTogKCkgPT4gdm9pZCBjYWxsKCdIaWRlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIE1heGltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIE1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ01heGltaXNlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNob3cgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgU2hvdzogKCkgPT4gdm9pZCBjYWxsKCdTaG93JyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIENsb3NlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIENsb3NlOiAoKSA9PiB2b2lkIGNhbGwoJ0Nsb3NlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFRvZ2dsZSB0aGUgd2luZG93IG1heGltaXNlIHN0YXRlXG4gICAgICAgICAqL1xuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdUb2dnbGVNYXhpbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbm1heGltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVW5NYXhpbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBNaW5pbWlzZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBNaW5pbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdNaW5pbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbm1pbmltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVW5NaW5pbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSByIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gZyAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBhIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IHZvaWQgY2FsbCgnU2V0QmFja2dyb3VuZENvbG91cicsIHtyLCBnLCBiLCBhfSksXG4gICAgfTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5DdXN0b21FdmVudH0gQ3VzdG9tRXZlbnRcbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZXZlbnRzXCIpO1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuXG4vKipcbiAqIEN1c3RvbUV2ZW50IGRlZmluZXMgYSBjdXN0b20gZXZlbnQuIEl0IGlzIHBhc3NlZCB0byBldmVudCBsaXN0ZW5lcnMuXG4gKlxuICogQGNsYXNzIEN1c3RvbUV2ZW50XG4gKi9cbmV4cG9ydCBjbGFzcyBDdXN0b21FdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBDdXN0b21FdmVudC5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gICAgICogQHBhcmFtIHthbnl9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcbiAgICAgKiBAbWVtYmVyb2YgQ3VzdG9tRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihDdXN0b21FdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG4vKipcbiAqIGxpc3RlbmVyT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cbiAqXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxuICovXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuXG4vKipcbiAqIGRpc3BhdGNoZXMgYW4gZXZlbnQgdG8gYWxsIGxpc3RlbmVyc1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7Q3VzdG9tRXZlbnR9IGV2ZW50XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkaXNwYXRjaEN1c3RvbUV2ZW50KGV2ZW50KSB7XG4gICAgY29uc29sZS5sb2coXCJkaXNwYXRjaGluZyBldmVudDogXCIsIHtldmVudH0pO1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XG4gICAgaWYgKGxpc3RlbmVycykge1xuICAgICAgICAvLyBpdGVyYXRlIGxpc3RlbmVycyBhbmQgY2FsbCBjYWxsYmFjay4gSWYgY2FsbGJhY2sgcmV0dXJucyB0cnVlLCByZW1vdmUgbGlzdGVuZXJcbiAgICAgICAgbGV0IHRvUmVtb3ZlID0gW107XG4gICAgICAgIGxpc3RlbmVycy5mb3JFYWNoKGxpc3RlbmVyID0+IHtcbiAgICAgICAgICAgIGxldCByZW1vdmUgPSBsaXN0ZW5lci5DYWxsYmFjayhldmVudCk7XG4gICAgICAgICAgICBpZiAocmVtb3ZlKSB7XG4gICAgICAgICAgICAgICAgdG9SZW1vdmUucHVzaChsaXN0ZW5lcik7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgICAgICAvLyByZW1vdmUgbGlzdGVuZXJzXG4gICAgICAgIGlmICh0b1JlbW92ZS5sZW5ndGggPiAwKSB7XG4gICAgICAgICAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gIXRvUmVtb3ZlLmluY2x1ZGVzKGwpKTtcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPbixcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJzIGNhbiBiZSB1bnJlZ2lzdGVyZWQgdmlhIGBhZGRpdGlvbmFsRXZlbnROYW1lc2BcbiAqXG4gW3YzIENIQU5HRV0gT2ZmIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIGxldCBldmVudHNUb1JlbW92ZSA9IFtldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzXTtcbiAgICBldmVudHNUb1JlbW92ZS5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xuICAgIH0pO1xufVxuXG4vKipcbiAqIE9mZkFsbCB1bnJlZ2lzdGVycyBhbGwgbGlzdGVuZXJzXG4gKiBbdjMgQ0hBTkdFXSBPZmZBbGwgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpIHtcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xufVxuXG4vKipcbiAqIEVtaXQgYW4gZXZlbnRcbiAqIEBwYXJhbSB7Q3VzdG9tRXZlbnR9IGV2ZW50IFRoZSBldmVudCB0byBlbWl0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XG4gICAgdm9pZCBjYWxsKFwiRW1pdFwiLCBldmVudCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk1lc3NhZ2VEaWFsb2dPcHRpb25zfSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk9wZW5EaWFsb2dPcHRpb25zfSBPcGVuRGlhbG9nT3B0aW9uc1xuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNhdmVEaWFsb2dPcHRpb25zfSBTYXZlRGlhbG9nT3B0aW9uc1xuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImRpYWxvZ1wiKTtcblxubGV0IGRpYWxvZ1Jlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucykge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XG4gICAgICAgIG9wdGlvbnNbXCJkaWFsb2ctaWRcIl0gPSBpZDtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuXG4vKipcbiAqIFNob3dzIGFuIEluZm8gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiSW5mb1wiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhbiBXYXJuaW5nIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIldhcm5pbmdcIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYW4gRXJyb3IgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkVycm9yXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGEgUXVlc3Rpb24gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlF1ZXN0aW9uXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGFuIE9wZW4gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlIG9yIGFuIGFycmF5IG9mIHNlbGVjdGVkIGZpbGVzIGlmIEFsbG93c011bHRpcGxlU2VsZWN0aW9uIGlzIHRydWUuIEEgYmxhbmsgc3RyaW5nIGlzIHJldHVybmVkIGlmIG5vIGZpbGUgd2FzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJPcGVuRmlsZVwiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhIFNhdmUgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiU2F2ZUZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbiIsICJpbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY29udGV4dG1lbnVcIik7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHJldHVybiBjYWxsKFwiT3BlbkNvbnRleHRNZW51XCIsIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZW5hYmxlQ29udGV4dE1lbnVzKGVuYWJsZWQpIHtcbiAgICBpZiAoZW5hYmxlZCkge1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcbiAgICBwcm9jZXNzQ29udGV4dE1lbnUoZXZlbnQudGFyZ2V0LCBldmVudCk7XG59XG5cbmZ1bmN0aW9uIHByb2Nlc3NDb250ZXh0TWVudShlbGVtZW50LCBldmVudCkge1xuICAgIGxldCBpZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51Jyk7XG4gICAgaWYgKGlkKSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShpZCwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtY29udGV4dG1lbnUtZGF0YScpKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBsZXQgcGFyZW50ID0gZWxlbWVudC5wYXJlbnRFbGVtZW50O1xuICAgICAgICBpZiAocGFyZW50KSB7XG4gICAgICAgICAgICBwcm9jZXNzQ29udGV4dE1lbnUocGFyZW50LCBldmVudCk7XG4gICAgICAgIH1cbiAgICB9XG59XG4iLCAiXG5pbXBvcnQge0VtaXR9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnQpIHtcbiAgICBFbWl0KGV2ZW50KTtcbn1cblxuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtZXZlbnRdJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICB9O1xuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG5cbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2QobWV0aG9kKSB7XG4gICAgaWYgKHdhaWxzLldpbmRvd1ttZXRob2RdID09PSB1bmRlZmluZWQpIHtcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xuICAgIH1cbiAgICB3YWlscy5XaW5kb3dbbWV0aG9kXSgpO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtd2luZG93XScpO1xuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVsb2FkV01MKCkge1xuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gJy4vY2xpcGJvYXJkJztcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gJy4vYXBwbGljYXRpb24nO1xuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSAnLi9zY3JlZW5zJztcbmltcG9ydCB7UGx1Z2luLCBDYWxsLCBjYWxsRXJyb3JDYWxsYmFjaywgY2FsbENhbGxiYWNrfSBmcm9tIFwiLi9jYWxsc1wiO1xuaW1wb3J0IHtuZXdXaW5kb3d9IGZyb20gXCIuL3dpbmRvd1wiO1xuaW1wb3J0IHtkaXNwYXRjaEN1c3RvbUV2ZW50LCBFbWl0LCBPZmYsIE9mZkFsbCwgT24sIE9uY2UsIE9uTXVsdGlwbGV9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtkaWFsb2dDYWxsYmFjaywgZGlhbG9nRXJyb3JDYWxsYmFjaywgRXJyb3IsIEluZm8sIE9wZW5GaWxlLCBRdWVzdGlvbiwgU2F2ZUZpbGUsIFdhcm5pbmcsfSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5pbXBvcnQge2VuYWJsZUNvbnRleHRNZW51c30gZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCB7cmVsb2FkV01MfSBmcm9tIFwiLi93bWxcIjtcblxud2luZG93LndhaWxzID0ge1xuICAgIC4uLm5ld1J1bnRpbWUoLTEpLFxufTtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cuX3dhaWxzID0ge1xuICAgIGRpYWxvZ0NhbGxiYWNrLFxuICAgIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssXG4gICAgZGlzcGF0Y2hDdXN0b21FdmVudCxcbiAgICBjYWxsQ2FsbGJhY2ssXG4gICAgY2FsbEVycm9yQ2FsbGJhY2ssXG59O1xuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZShpZCkge1xuICAgIHJldHVybiB7XG4gICAgICAgIENsaXBib2FyZDoge1xuICAgICAgICAgICAgLi4uQ2xpcGJvYXJkXG4gICAgICAgIH0sXG4gICAgICAgIEFwcGxpY2F0aW9uOiB7XG4gICAgICAgICAgICAuLi5BcHBsaWNhdGlvblxuICAgICAgICB9LFxuICAgICAgICBMb2csXG4gICAgICAgIFNjcmVlbnMsXG4gICAgICAgIENhbGwsXG4gICAgICAgIFBsdWdpbixcbiAgICAgICAgV01MOiB7XG4gICAgICAgICAgICBSZWxvYWQ6IHJlbG9hZFdNTCxcbiAgICAgICAgfSxcbiAgICAgICAgRGlhbG9nOiB7XG4gICAgICAgICAgICBJbmZvLFxuICAgICAgICAgICAgV2FybmluZyxcbiAgICAgICAgICAgIEVycm9yLFxuICAgICAgICAgICAgUXVlc3Rpb24sXG4gICAgICAgICAgICBPcGVuRmlsZSxcbiAgICAgICAgICAgIFNhdmVGaWxlLFxuICAgICAgICB9LFxuICAgICAgICBFdmVudHM6IHtcbiAgICAgICAgICAgIEVtaXQsXG4gICAgICAgICAgICBPbixcbiAgICAgICAgICAgIE9uY2UsXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxuICAgICAgICAgICAgT2ZmLFxuICAgICAgICAgICAgT2ZmQWxsLFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyhpZCksXG4gICAgfTtcbn1cblxuaWYgKERFQlVHKSB7XG4gICAgY29uc29sZS5sb2coXCJXYWlscyB2My4wLjAgRGVidWcgTW9kZSBFbmFibGVkXCIpO1xufVxuXG5lbmFibGVDb250ZXh0TWVudXModHJ1ZSk7XG5cbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoXCJET01Db250ZW50TG9hZGVkXCIsIGZ1bmN0aW9uKGV2ZW50KSB7XG4gICAgcmVsb2FkV01MKCk7XG59KTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ1lBLE1BQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUU1QyxXQUFTLFlBQVksUUFBUSxNQUFNO0FBQy9CLFFBQUks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R0FBRztBQUFBLE1BQ1A7QUFBQSxNQUNBLGFBQWE7QUFBQSxRQUNULEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0EsS0FBSztBQUFBLFFBQ0QsUUFBUTtBQUFBLE1BQ1o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0EsT0FBQUM7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUSxVQUFVLEVBQUU7QUFBQSxJQUN4QjtBQUFBLEVBQ0o7QUFFQSxNQUFJLE1BQU87QUFDUCxZQUFRLElBQUksaUNBQWlDO0FBQUEsRUFDakQ7QUFFQSxxQkFBbUIsSUFBSTtBQUV2QixXQUFTLGlCQUFpQixvQkFBb0IsU0FBUyxPQUFPO0FBQzFELGNBQVU7QUFBQSxFQUNkLENBQUM7IiwKICAibmFtZXMiOiBbImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImNhbGwiLCAiZ2VuZXJhdGVJRCIsICJFcnJvciIsICJjYWxsIiwgIkVycm9yIl0KfQo=
