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
  document.addEventListener("DOMContentLoaded", function(event) {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuLyoqXG4gKiBTZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICB2b2lkIGNhbGwoXCJTZXRUZXh0XCIsIHt0ZXh0fSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBpZiAoYXJncykge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xuICAgIH1cbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiYXBwbGljYXRpb25cIik7XG5cbi8qKlxuICogSGlkZSB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgdm9pZCBjYWxsKFwiSGlkZVwiKTtcbn1cblxuLyoqXG4gKiBTaG93IHRoZSBhcHBsaWNhdGlvblxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICB2b2lkIGNhbGwoXCJTaG93XCIpO1xufVxuXG5cbi8qKlxuICogUXVpdCB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgdm9pZCBjYWxsKFwiUXVpdFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImxvZ1wiKTtcblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZS5cbiAqIEBwYXJhbSB7bWVzc2FnZX0gTWVzc2FnZSB0byBsb2dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZyhtZXNzYWdlKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJMb2dcIiwgbWVzc2FnZSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuU2NyZWVufSBTY3JlZW5cbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwic2NyZWVuc1wiKTtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0QWxsXCIpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRQcmltYXJ5XCIpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gKiBAY29uc3RydWN0b3JcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRDdXJyZW50XCIpO1xufSIsICJsZXQgdXJsQWxwaGFiZXQgPVxuICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcbmV4cG9ydCBsZXQgY3VzdG9tQWxwaGFiZXQgPSAoYWxwaGFiZXQsIGRlZmF1bHRTaXplID0gMjEpID0+IHtcbiAgcmV0dXJuIChzaXplID0gZGVmYXVsdFNpemUpID0+IHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIGxldCBpID0gc2l6ZVxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgIGlkICs9IGFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogYWxwaGFiZXQubGVuZ3RoKSB8IDBdXG4gICAgfVxuICAgIHJldHVybiBpZFxuICB9XG59XG5leHBvcnQgbGV0IG5hbm9pZCA9IChzaXplID0gMjEpID0+IHtcbiAgbGV0IGlkID0gJydcbiAgbGV0IGkgPSBzaXplXG4gIHdoaWxlIChpLS0pIHtcbiAgICBpZCArPSB1cmxBbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIDY0KSB8IDBdXG4gIH1cbiAgcmV0dXJuIGlkXG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY2FsbFwiKTtcblxubGV0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwKCk7XG5cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBjYWxsQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcbiAgICAgICAgfVxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gY2FsbEVycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucykge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhcIkNhbGxcIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogQ2FsbCBhIHBsdWdpbiBtZXRob2RcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gbmFtZSBvZiB0aGUgcGx1Z2luXG4gKiBAcGFyYW0ge3N0cmluZ30gbWV0aG9kTmFtZSAtIG5hbWUgb2YgdGhlIG1ldGhvZFxuICogQHBhcmFtIHsuLi5hbnl9IGFyZ3MgLSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxhbnk+fSAtIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSByZXN1bHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFBsdWdpbihwbHVnaW5OYW1lLCBtZXRob2ROYW1lLCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCB7XG4gICAgICAgIHBhY2thZ2VOYW1lOiBcIndhaWxzLXBsdWdpbnNcIixcbiAgICAgICAgc3RydWN0TmFtZTogcGx1Z2luTmFtZSxcbiAgICAgICAgbWV0aG9kTmFtZTogbWV0aG9kTmFtZSxcbiAgICAgICAgYXJnczogYXJncyxcbiAgICB9KTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNpemV9IFNpemVcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuUG9zaXRpb259IFBvc2l0aW9uXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXG4gKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdXaW5kb3cod2luZG93TmFtZSkge1xuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCB3aW5kb3dOYW1lKTtcbiAgICByZXR1cm4ge1xuICAgICAgICAvLyBSZWxvYWQ6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFJlbG9hZEFwcDogKCkgPT4gY2FsbCgnV1InKSxcbiAgICAgICAgLy8gU2V0U3lzdGVtRGVmYXVsdFRoZW1lOiAoKSA9PiBjYWxsKCdXQVNEVCcpLFxuICAgICAgICAvLyBTZXRMaWdodFRoZW1lOiAoKSA9PiBjYWxsKCdXQUxUJyksXG4gICAgICAgIC8vIFNldERhcmtUaGVtZTogKCkgPT4gY2FsbCgnV0FEVCcpLFxuICAgICAgICAvLyBJc0Z1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1dJRicpLFxuICAgICAgICAvLyBJc01heGltaXplZDogKCkgPT4gY2FsbCgnV0lNJyksXG4gICAgICAgIC8vIElzTWluaW1pemVkOiAoKSA9PiBjYWxsKCdXSU1OJyksXG4gICAgICAgIC8vIElzV2luZG93ZWQ6ICgpID0+IGNhbGwoJ1dJRicpLFxuXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdy5cbiAgICAgICAgICovXG4gICAgICAgIENlbnRlcjogKCkgPT4gdm9pZCBjYWxsKCdDZW50ZXInKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgdGl0bGUuXG4gICAgICAgICAqIEBwYXJhbSB0aXRsZVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gdm9pZCBjYWxsKCdTZXRUaXRsZScsIHt0aXRsZX0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBNYWtlcyB0aGUgd2luZG93IGZ1bGxzY3JlZW4uXG4gICAgICAgICAqL1xuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ0Z1bGxzY3JlZW4nKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVW5mdWxsc2NyZWVuIHRoZSB3aW5kb3cuXG4gICAgICAgICAqL1xuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IHZvaWQgY2FsbCgnVW5GdWxsc2NyZWVuJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHNpemUuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCBUaGUgd2luZG93IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgVGhlIHdpbmRvdyBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRTaXplJywge3dpZHRoLGhlaWdodH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBzaXplLlxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxTaXplPn0gVGhlIHdpbmRvdyBzaXplXG4gICAgICAgICAqL1xuICAgICAgICBTaXplOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTaXplJyk7IH0sXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1heGltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoJ1NldE1heFNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1pbmltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1pblNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoJ1NldE1pblNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB3aW5kb3cgdG8gYmUgYWx3YXlzIG9uIHRvcC5cbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSBvblRvcCBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGFsd2F5cyBvbiB0b3BcbiAgICAgICAgICovXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAob25Ub3ApID0+IHZvaWQgY2FsbCgnU2V0QWx3YXlzT25Ub3AnLCB7YWx3YXlzT25Ub3A6b25Ub3B9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgcG9zaXRpb24uXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB5XG4gICAgICAgICAqL1xuICAgICAgICBTZXRQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoJ1NldFBvc2l0aW9uJywge3gseX0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBwb3NpdGlvbi5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8UG9zaXRpb24+fSBUaGUgd2luZG93IHBvc2l0aW9uXG4gICAgICAgICAqL1xuICAgICAgICBQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUG9zaXRpb24nKTsgfSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBvbi5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cbiAgICAgICAgICovXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJyk7IH0sXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEhpZGUgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgSGlkZTogKCkgPT4gdm9pZCBjYWxsKCdIaWRlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIE1heGltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIE1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ01heGltaXNlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNob3cgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgU2hvdzogKCkgPT4gdm9pZCBjYWxsKCdTaG93JyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIENsb3NlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIENsb3NlOiAoKSA9PiB2b2lkIGNhbGwoJ0Nsb3NlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFRvZ2dsZSB0aGUgd2luZG93IG1heGltaXNlIHN0YXRlXG4gICAgICAgICAqL1xuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdUb2dnbGVNYXhpbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbm1heGltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVW5NYXhpbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBNaW5pbWlzZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBNaW5pbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdNaW5pbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbm1pbmltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVW5NaW5pbWlzZScpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSByIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gZyAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBhIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IHZvaWQgY2FsbCgnU2V0QmFja2dyb3VuZENvbG91cicsIHtyLCBnLCBiLCBhfSksXG4gICAgfTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5XYWlsc0V2ZW50fSBXYWlsc0V2ZW50XG4gKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImV2ZW50c1wiKTtcblxuLyoqXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcbiAqXG4gKiBAY2xhc3MgTGlzdGVuZXJcbiAqL1xuY2xhc3MgTGlzdGVuZXIge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICAgICAqIEBtZW1iZXJvZiBMaXN0ZW5lclxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xuICAgICAgICAgICAgY2FsbGJhY2soZGF0YSk7XG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cblxuLyoqXG4gKiBXYWlsc0V2ZW50IGRlZmluZXMgYSBjdXN0b20gZXZlbnQuIEl0IGlzIHBhc3NlZCB0byBldmVudCBsaXN0ZW5lcnMuXG4gKlxuICogQGNsYXNzIFdhaWxzRXZlbnRcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBuYW1lIC0gTmFtZSBvZiB0aGUgZXZlbnRcbiAqIEBwcm9wZXJ0eSB7YW55fSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XG4gKi9cbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIFdhaWxzRXZlbnQuXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IG5hbWUgLSBOYW1lIG9mIHRoZSBldmVudFxuICAgICAqIEBwYXJhbSB7YW55PW51bGx9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcbiAgICAgKiBAbWVtYmVyb2YgV2FpbHNFdmVudFxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG5hbWUsIGRhdGEgPSBudWxsKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb24oV2FpbHNFdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb24oV2FpbHNFdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG4vKipcbiAqIGxpc3RlbmVyT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cbiAqXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxuICovXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuXG4vKipcbiAqIGRpc3BhdGNoZXMgYW4gZXZlbnQgdG8gYWxsIGxpc3RlbmVyc1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7V2FpbHNFdmVudH0gZXZlbnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudCkge1xuICAgIGNvbnNvbGUubG9nKFwiZGlzcGF0Y2hpbmcgZXZlbnQ6IFwiLCB7ZXZlbnR9KTtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmIChsaXN0ZW5lcnMpIHtcbiAgICAgICAgLy8gaXRlcmF0ZSBsaXN0ZW5lcnMgYW5kIGNhbGwgY2FsbGJhY2suIElmIGNhbGxiYWNrIHJldHVybnMgdHJ1ZSwgcmVtb3ZlIGxpc3RlbmVyXG4gICAgICAgIGxldCB0b1JlbW92ZSA9IFtdO1xuICAgICAgICBsaXN0ZW5lcnMuZm9yRWFjaChsaXN0ZW5lciA9PiB7XG4gICAgICAgICAgICBsZXQgcmVtb3ZlID0gbGlzdGVuZXIuQ2FsbGJhY2soZXZlbnQpO1xuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xuICAgICAgICBpZiAodG9SZW1vdmUubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuIFt2MyBDSEFOR0VdIE9mZiBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBPZmZBbGwgdW5yZWdpc3RlcnMgYWxsIGxpc3RlbmVyc1xuICogW3YzIENIQU5HRV0gT2ZmQWxsIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcbiAqXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLyoqXG4gKiBFbWl0IGFuIGV2ZW50XG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50IFRoZSBldmVudCB0byBlbWl0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XG4gICAgdm9pZCBjYWxsKFwiRW1pdFwiLCBldmVudCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk1lc3NhZ2VEaWFsb2dPcHRpb25zfSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk9wZW5EaWFsb2dPcHRpb25zfSBPcGVuRGlhbG9nT3B0aW9uc1xuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNhdmVEaWFsb2dPcHRpb25zfSBTYXZlRGlhbG9nT3B0aW9uc1xuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImRpYWxvZ1wiKTtcblxubGV0IGRpYWxvZ1Jlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucykge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XG4gICAgICAgIG9wdGlvbnNbXCJkaWFsb2ctaWRcIl0gPSBpZDtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuXG4vKipcbiAqIFNob3dzIGFuIEluZm8gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiSW5mb1wiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhbiBXYXJuaW5nIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIldhcm5pbmdcIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYW4gRXJyb3IgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkVycm9yXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGEgUXVlc3Rpb24gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlF1ZXN0aW9uXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGFuIE9wZW4gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlIG9yIGFuIGFycmF5IG9mIHNlbGVjdGVkIGZpbGVzIGlmIEFsbG93c011bHRpcGxlU2VsZWN0aW9uIGlzIHRydWUuIEEgYmxhbmsgc3RyaW5nIGlzIHJldHVybmVkIGlmIG5vIGZpbGUgd2FzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJPcGVuRmlsZVwiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhIFNhdmUgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiU2F2ZUZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbiIsICJpbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY29udGV4dG1lbnVcIik7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHJldHVybiBjYWxsKFwiT3BlbkNvbnRleHRNZW51XCIsIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZW5hYmxlQ29udGV4dE1lbnVzKGVuYWJsZWQpIHtcbiAgICBpZiAoZW5hYmxlZCkge1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcbiAgICBwcm9jZXNzQ29udGV4dE1lbnUoZXZlbnQudGFyZ2V0LCBldmVudCk7XG59XG5cbmZ1bmN0aW9uIHByb2Nlc3NDb250ZXh0TWVudShlbGVtZW50LCBldmVudCkge1xuICAgIGxldCBpZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51Jyk7XG4gICAgaWYgKGlkKSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShpZCwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtY29udGV4dG1lbnUtZGF0YScpKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBsZXQgcGFyZW50ID0gZWxlbWVudC5wYXJlbnRFbGVtZW50O1xuICAgICAgICBpZiAocGFyZW50KSB7XG4gICAgICAgICAgICBwcm9jZXNzQ29udGV4dE1lbnUocGFyZW50LCBldmVudCk7XG4gICAgICAgIH1cbiAgICB9XG59XG4iLCAiXG5pbXBvcnQge0VtaXQsIFdhaWxzRXZlbnR9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lLCBkYXRhPW51bGwpIHtcbiAgICBsZXQgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpO1xuICAgIEVtaXQoZXZlbnQpO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC1ldmVudF0nKTtcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIH07XG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcblxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9KTtcbn1cblxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcbiAgICBpZiAod2FpbHMuV2luZG93W21ldGhvZF0gPT09IHVuZGVmaW5lZCkge1xuICAgICAgICBjb25zb2xlLmxvZyhcIldpbmRvdyBtZXRob2QgXCIgKyBtZXRob2QgKyBcIiBub3QgZm91bmRcIik7XG4gICAgfVxuICAgIHdhaWxzLldpbmRvd1ttZXRob2RdKCk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC13aW5kb3ddJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiByZWxvYWRXTUwoKSB7XG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcbiAgICBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSAnLi9jbGlwYm9hcmQnO1xuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSAnLi9hcHBsaWNhdGlvbic7XG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tICcuL3NjcmVlbnMnO1xuaW1wb3J0IHtQbHVnaW4sIENhbGwsIGNhbGxFcnJvckNhbGxiYWNrLCBjYWxsQ2FsbGJhY2t9IGZyb20gXCIuL2NhbGxzXCI7XG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQge2Rpc3BhdGNoV2FpbHNFdmVudCwgRW1pdCwgT2ZmLCBPZmZBbGwsIE9uLCBPbmNlLCBPbk11bHRpcGxlfSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7ZGlhbG9nQ2FsbGJhY2ssIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssIEVycm9yLCBJbmZvLCBPcGVuRmlsZSwgUXVlc3Rpb24sIFNhdmVGaWxlLCBXYXJuaW5nLH0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuaW1wb3J0IHtlbmFibGVDb250ZXh0TWVudXN9IGZyb20gXCIuL2NvbnRleHRtZW51XCI7XG5pbXBvcnQge3JlbG9hZFdNTH0gZnJvbSBcIi4vd21sXCI7XG5cbndpbmRvdy53YWlscyA9IHtcbiAgICAuLi5uZXdSdW50aW1lKG51bGwpLFxufTtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cuX3dhaWxzID0ge1xuICAgIGRpYWxvZ0NhbGxiYWNrLFxuICAgIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssXG4gICAgZGlzcGF0Y2hXYWlsc0V2ZW50LFxuICAgIGNhbGxDYWxsYmFjayxcbiAgICBjYWxsRXJyb3JDYWxsYmFjayxcbn07XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lKHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4ge1xuICAgICAgICBDbGlwYm9hcmQ6IHtcbiAgICAgICAgICAgIC4uLkNsaXBib2FyZFxuICAgICAgICB9LFxuICAgICAgICBBcHBsaWNhdGlvbjoge1xuICAgICAgICAgICAgLi4uQXBwbGljYXRpb24sXG4gICAgICAgICAgICBHZXRXaW5kb3dCeU5hbWUod2luZG93TmFtZSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBuZXdSdW50aW1lKHdpbmRvd05hbWUpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9LFxuICAgICAgICBMb2csXG4gICAgICAgIFNjcmVlbnMsXG4gICAgICAgIENhbGwsXG4gICAgICAgIFBsdWdpbixcbiAgICAgICAgV01MOiB7XG4gICAgICAgICAgICBSZWxvYWQ6IHJlbG9hZFdNTCxcbiAgICAgICAgfSxcbiAgICAgICAgRGlhbG9nOiB7XG4gICAgICAgICAgICBJbmZvLFxuICAgICAgICAgICAgV2FybmluZyxcbiAgICAgICAgICAgIEVycm9yLFxuICAgICAgICAgICAgUXVlc3Rpb24sXG4gICAgICAgICAgICBPcGVuRmlsZSxcbiAgICAgICAgICAgIFNhdmVGaWxlLFxuICAgICAgICB9LFxuICAgICAgICBFdmVudHM6IHtcbiAgICAgICAgICAgIEVtaXQsXG4gICAgICAgICAgICBPbixcbiAgICAgICAgICAgIE9uY2UsXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxuICAgICAgICAgICAgT2ZmLFxuICAgICAgICAgICAgT2ZmQWxsLFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyh3aW5kb3dOYW1lKSxcbiAgICB9O1xufVxuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbmVuYWJsZUNvbnRleHRNZW51cyh0cnVlKTtcblxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcihcIkRPTUNvbnRlbnRMb2FkZWRcIiwgZnVuY3Rpb24oZXZlbnQpIHtcbiAgICByZWxvYWRXTUwoKTtcbn0pOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDWUEsTUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBRTVDLFdBQVMsWUFBWSxRQUFRLFlBQVksTUFBTTtBQUMzQyxRQUFJLE1BQU0sSUFBSSxJQUFJLFVBQVU7QUFDNUIsUUFBSSxhQUFhLE9BQU8sVUFBVSxNQUFNO0FBQ3hDLFFBQUksTUFBTTtBQUNOLFVBQUksYUFBYSxPQUFPLFFBQVEsS0FBSyxVQUFVLElBQUksQ0FBQztBQUFBLElBQ3hEO0FBQ0EsUUFBSSxlQUFlO0FBQUEsTUFDZixTQUFTLENBQUM7QUFBQSxJQUNkO0FBQ0EsUUFBSSxZQUFZO0FBQ1osbUJBQWEsUUFBUSxxQkFBcUIsSUFBSTtBQUFBLElBQ2xEO0FBQ0EsV0FBTyxJQUFJLFFBQVEsQ0FBQyxTQUFTLFdBQVc7QUFDcEMsWUFBTSxLQUFLLFlBQVksRUFDbEIsS0FBSyxjQUFZO0FBQ2QsWUFBSSxTQUFTLElBQUk7QUFFYixjQUFJLFNBQVMsUUFBUSxJQUFJLGNBQWMsS0FBSyxTQUFTLFFBQVEsSUFBSSxjQUFjLEVBQUUsUUFBUSxrQkFBa0IsTUFBTSxJQUFJO0FBQ2pILG1CQUFPLFNBQVMsS0FBSztBQUFBLFVBQ3pCLE9BQU87QUFDSCxtQkFBTyxTQUFTLEtBQUs7QUFBQSxVQUN6QjtBQUFBLFFBQ0o7QUFDQSxlQUFPLE1BQU0sU0FBUyxVQUFVLENBQUM7QUFBQSxNQUNyQyxDQUFDLEVBQ0EsS0FBSyxVQUFRLFFBQVEsSUFBSSxDQUFDLEVBQzFCLE1BQU0sV0FBUyxPQUFPLEtBQUss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QUFJLENBQUMsRUFBQyxDQUFDLEVBQUUsS0FBSyxTQUFVLFFBQVE7QUFDdkgsZ0JBQUksV0FBVyxNQUFNO0FBQ2pCLCtCQUFpQixZQUFZO0FBQUEsWUFDakM7QUFBQSxVQUNKLENBQUM7QUFDRDtBQUFBLFFBQ0o7QUFDQSx5QkFBaUIsWUFBWTtBQUFBLE1BQ2pDO0FBR0EsY0FBUSxvQkFBb0IsU0FBUyxRQUFRO0FBRzdDLGNBQVEsaUJBQWlCLFNBQVMsUUFBUTtBQUFBLElBQzlDLENBQUM7QUFBQSxFQUNMO0FBRU8sV0FBUyxZQUFZO0FBQ3hCLHlCQUFxQjtBQUNyQiwwQkFBc0I7QUFBQSxFQUMxQjs7O0FDbERBLFNBQU8sUUFBUTtBQUFBLElBQ1gsR0FBRyxXQUFXLElBQUk7QUFBQSxFQUN0QjtBQUdBLFNBQU8sU0FBUztBQUFBLElBQ1o7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsRUFDSjtBQUVPLFdBQVMsV0FBVyxZQUFZO0FBQ25DLFdBQU87QUFBQSxNQUNILFdBQVc7QUFBQSxRQUNQLEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQSxhQUFhO0FBQUEsUUFDVCxHQUFHO0FBQUEsUUFDSCxnQkFBZ0JDLGFBQVk7QUFDeEIsaUJBQU8sV0FBV0EsV0FBVTtBQUFBLFFBQ2hDO0FBQUEsTUFDSjtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFFBQVE7QUFBQSxNQUNaO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxVQUFVO0FBQUEsSUFDaEM7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEscUJBQW1CLElBQUk7QUFFdkIsV0FBUyxpQkFBaUIsb0JBQW9CLFNBQVMsT0FBTztBQUMxRCxjQUFVO0FBQUEsRUFDZCxDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJldmVudE5hbWUiLCAiY2FsbCIsICJnZW5lcmF0ZUlEIiwgIkVycm9yIiwgImNhbGwiLCAid2luZG93TmFtZSIsICJFcnJvciJdCn0K
