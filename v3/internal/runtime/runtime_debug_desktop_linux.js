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
      args["windowID"] = id;
      return runtimeCall(object + "." + method, args);
    };
  }

  // desktop/clipboard.js
  var call = newRuntimeCaller("clipboard");
  function SetText(text) {
    return call("SetText", { text });
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
    return call2("Hide");
  }
  function Show() {
    return call2("Show");
  }
  function Quit() {
    return call2("Quit");
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
  function Plugin(options) {
    return callBinding("Plugin", options);
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
      Center: () => call9("Center"),
      SetTitle: (title) => call9("SetTitle", { title }),
      Fullscreen: () => call9("Fullscreen"),
      UnFullscreen: () => call9("UnFullscreen"),
      SetSize: (width, height) => call9("SetSize", { width, height }),
      Size: () => {
        return call9("Size");
      },
      SetMaxSize: (width, height) => call9("SetMaxSize", { width, height }),
      SetMinSize: (width, height) => call9("SetMinSize", { width, height }),
      SetAlwaysOnTop: (b) => call9("SetAlwaysOnTop", { alwaysOnTop: b }),
      SetPosition: (x, y) => call9("SetPosition", { x, y }),
      Position: () => {
        return call9("Position");
      },
      Screen: () => {
        return call9("Screen");
      },
      Hide: () => call9("Hide"),
      Maximise: () => call9("Maximise"),
      Show: () => call9("Show"),
      Close: () => call9("Close"),
      ToggleMaximise: () => call9("ToggleMaximise"),
      UnMaximise: () => call9("UnMaximise"),
      Minimise: () => call9("Minimise"),
      UnMinimise: () => call9("UnMinimise"),
      SetBackgroundColour: (r, g, b, a) => call9("SetBackgroundColour", { r, g, b, a })
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
    return call6("Emit", event);
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
    let _ = Emit({ name: event });
  }
  function addWMLEventListeners() {
    const elements = document.querySelectorAll("[data-wml-event]");
    for (let i = 0; i < elements.length; i++) {
      const element = elements[i];
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
    }
  }
  function callWindowMethod(method) {
    if (wails.Window[method] === void 0) {
      console.log("Window method " + method + " not found");
    }
    wails.Window[method]();
  }
  function addWMLWindowListeners() {
    const elements = document.querySelectorAll("[data-wml-window]");
    for (let i = 0; i < elements.length; i++) {
      const element = elements[i];
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
    }
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChcIlRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIGlmKGFyZ3MpIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcbiAgICB9XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZmV0Y2godXJsKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgaWQpIHtcbiAgICBpZiAoIWlkIHx8IGlkID09PSAtMSkge1xuICAgICAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncykge1xuICAgICAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICAgICAgfTtcbiAgICB9XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgYXJncyA9IGFyZ3MgfHwge307XG4gICAgICAgIGFyZ3NbXCJ3aW5kb3dJRFwiXSA9IGlkO1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGwob2JqZWN0ICsgXCIuXCIgKyBtZXRob2QsIGFyZ3MpO1xuICAgIH1cbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImFwcGxpY2F0aW9uXCIpO1xuXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcbiAgICByZXR1cm4gY2FsbChcIkhpZGVcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xuICAgIHJldHVybiBjYWxsKFwiU2hvd1wiKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJRdWl0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwibG9nXCIpO1xuXG4vKipcbiAqIExvZ3MgYSBtZXNzYWdlLlxuICogQHBhcmFtIHttZXNzYWdlfSBNZXNzYWdlIHRvIGxvZ1xuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nKG1lc3NhZ2UpIHtcbiAgICByZXR1cm4gY2FsbChcIkxvZ1wiLCBtZXNzYWdlKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwic2NyZWVuc1wiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldEFsbFwiKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRQcmltYXJ5XCIpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldEN1cnJlbnRcIik7XG59IiwgImxldCB1cmxBbHBoYWJldCA9XG4gICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuZXhwb3J0IGxldCBjdXN0b21BbHBoYWJldCA9IChhbHBoYWJldCwgZGVmYXVsdFNpemUgPSAyMSkgPT4ge1xuICByZXR1cm4gKHNpemUgPSBkZWZhdWx0U2l6ZSkgPT4ge1xuICAgIGxldCBpZCA9ICcnXG4gICAgbGV0IGkgPSBzaXplXG4gICAgd2hpbGUgKGktLSkge1xuICAgICAgaWQgKz0gYWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiBhbHBoYWJldC5sZW5ndGgpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG4gIH1cbn1cbmV4cG9ydCBsZXQgbmFub2lkID0gKHNpemUgPSAyMSkgPT4ge1xuICBsZXQgaWQgPSAnJ1xuICBsZXQgaSA9IHNpemVcbiAgd2hpbGUgKGktLSkge1xuICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgfVxuICByZXR1cm4gaWRcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjYWxsXCIpO1xuXG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBjYWxsRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xuICAgIGxldCBwID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjYWxsQmluZGluZyh0eXBlLCBvcHRpb25zKSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcbiAgICAgICAgb3B0aW9uc1tcImNhbGwtaWRcIl0gPSBpZDtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5jYXRjaCgoZXJyb3IpID0+IHtcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XG4gICAgICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFBsdWdpbihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiUGx1Z2luXCIsIG9wdGlvbnMpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1dpbmRvdyhpZCkge1xuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCBpZCk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBSZWxvYWRBcHA6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFNldFN5c3RlbURlZmF1bHRUaGVtZTogKCkgPT4gY2FsbCgnV0FTRFQnKSxcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxuICAgICAgICAvLyBTZXREYXJrVGhlbWU6ICgpID0+IGNhbGwoJ1dBRFQnKSxcbiAgICAgICAgLy8gSXNGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxuICAgICAgICAvLyBJc01pbmltaXplZDogKCkgPT4gY2FsbCgnV0lNTicpLFxuICAgICAgICAvLyBJc1dpbmRvd2VkOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKCdDZW50ZXInKSxcbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbCgnU2V0VGl0bGUnLCB7dGl0bGV9KSxcbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnRnVsbHNjcmVlbicpLFxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2l6ZTogKCkgPT4geyByZXR1cm4gY2FsbCgnU2l6ZScpIH0sXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRNYXhTaXplJywge3dpZHRoLGhlaWdodH0pLFxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWluU2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0QWx3YXlzT25Ub3A6IChiKSA9PiBjYWxsKCdTZXRBbHdheXNPblRvcCcsIHthbHdheXNPblRvcDpifSksXG4gICAgICAgIFNldFBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbCgnU2V0UG9zaXRpb24nLCB7eCx5fSksXG4gICAgICAgIFBvc2l0aW9uOiAoKSA9PiB7IHJldHVybiBjYWxsKCdQb3NpdGlvbicpIH0sXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJykgfSxcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbCgnSGlkZScpLFxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gY2FsbCgnTWF4aW1pc2UnKSxcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbCgnU2hvdycpLFxuICAgICAgICBDbG9zZTogKCkgPT4gY2FsbCgnQ2xvc2UnKSxcbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1RvZ2dsZU1heGltaXNlJyksXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1VuTWF4aW1pc2UnKSxcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IGNhbGwoJ01pbmltaXNlJyksXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwoJ1VuTWluaW1pc2UnKSxcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IGNhbGwoJ1NldEJhY2tncm91bmRDb2xvdXInLCB7ciwgZywgYiwgYX0pLFxuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZXZlbnRzXCIpO1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuXG4vKipcbiAqIEN1c3RvbUV2ZW50IGRlZmluZXMgYSBjdXN0b20gZXZlbnQuIEl0IGlzIHBhc3NlZCB0byBldmVudCBsaXN0ZW5lcnMuXG4gKlxuICogQGNsYXNzIEN1c3RvbUV2ZW50XG4gKi9cbmV4cG9ydCBjbGFzcyBDdXN0b21FdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBDdXN0b21FdmVudC5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gICAgICogQHBhcmFtIHthbnl9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcbiAgICAgKiBAbWVtYmVyb2YgQ3VzdG9tRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihDdXN0b21FdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbmZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG5cbi8qKlxuICogZGlzcGF0Y2hlcyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtDdXN0b21FdmVudH0gZXZlbnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGRpc3BhdGNoQ3VzdG9tRXZlbnQoZXZlbnQpIHtcbiAgICBjb25zb2xlLmxvZyhcImRpc3BhdGNoaW5nIGV2ZW50OiBcIiwge2V2ZW50fSk7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBbXTtcbiAgICAgICAgbGlzdGVuZXJzLmZvckVhY2gobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KVxuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xuICAgICAgICBpZiAodG9SZW1vdmUubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuIFt2MyBDSEFOR0VdIE9mZiBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBPZmZBbGwgdW5yZWdpc3RlcnMgYWxsIGxpc3RlbmVyc1xuICogW3YzIENIQU5HRV0gT2ZmQWxsIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcbiAqXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLypcbiAgIEVtaXQgZW1pdHMgYW4gZXZlbnQgdG8gYWxsIGxpc3RlbmVyc1xuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkge1xuICAgIHJldHVybiBjYWxsKFwiRW1pdFwiLCBldmVudCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImRpYWxvZ1wiKTtcblxubGV0IGRpYWxvZ1Jlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucykge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XG4gICAgICAgIG9wdGlvbnNbXCJkaWFsb2ctaWRcIl0gPSBpZDtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkluZm9cIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXYXJuaW5nKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiV2FybmluZ1wiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiRXJyb3JcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlF1ZXN0aW9uXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJPcGVuRmlsZVwiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiU2F2ZUZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbiIsICJpbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY29udGV4dG1lbnVcIik7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHJldHVybiBjYWxsKFwiT3BlbkNvbnRleHRNZW51XCIsIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZW5hYmxlQ29udGV4dE1lbnVzKGVuYWJsZWQpIHtcbiAgICBpZiAoZW5hYmxlZCkge1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcbiAgICBwcm9jZXNzQ29udGV4dE1lbnUoZXZlbnQudGFyZ2V0LCBldmVudCk7XG59XG5cbmZ1bmN0aW9uIHByb2Nlc3NDb250ZXh0TWVudShlbGVtZW50LCBldmVudCkge1xuICAgIGxldCBpZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51Jyk7XG4gICAgaWYgKGlkKSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShpZCwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtY29udGV4dG1lbnUtZGF0YScpKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBsZXQgcGFyZW50ID0gZWxlbWVudC5wYXJlbnRFbGVtZW50O1xuICAgICAgICBpZiAocGFyZW50KSB7XG4gICAgICAgICAgICBwcm9jZXNzQ29udGV4dE1lbnUocGFyZW50LCBldmVudCk7XG4gICAgICAgIH1cbiAgICB9XG59XG4iLCAiXG5pbXBvcnQge0VtaXR9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnQpIHtcbiAgIGxldCBfID0gRW1pdCh7bmFtZTogZXZlbnR9ICk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCkge1xuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW2RhdGEtd21sLWV2ZW50XScpO1xuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZWxlbWVudHMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGVsZW1lbnRzW2ldO1xuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICB9XG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcblxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2QobWV0aG9kKSB7XG4gICAgaWYgKHdhaWxzLldpbmRvd1ttZXRob2RdID09PSB1bmRlZmluZWQpIHtcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xuICAgIH1cbiAgICB3YWlscy5XaW5kb3dbbWV0aG9kXSgpO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtd2luZG93XScpO1xuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZWxlbWVudHMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGVsZW1lbnRzW2ldO1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiByZWxvYWRXTUwoKSB7XG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcbiAgICBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSAnLi9jbGlwYm9hcmQnO1xuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSAnLi9hcHBsaWNhdGlvbic7XG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tICcuL3NjcmVlbnMnO1xuaW1wb3J0IHtQbHVnaW4sIENhbGwsIGNhbGxFcnJvckNhbGxiYWNrLCBjYWxsQ2FsbGJhY2t9IGZyb20gXCIuL2NhbGxzXCI7XG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQge2Rpc3BhdGNoQ3VzdG9tRXZlbnQsIEVtaXQsIE9mZiwgT2ZmQWxsLCBPbiwgT25jZSwgT25NdWx0aXBsZX0gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQge2RpYWxvZ0NhbGxiYWNrLCBkaWFsb2dFcnJvckNhbGxiYWNrLCBFcnJvciwgSW5mbywgT3BlbkZpbGUsIFF1ZXN0aW9uLCBTYXZlRmlsZSwgV2FybmluZyx9IGZyb20gXCIuL2RpYWxvZ3NcIjtcbmltcG9ydCB7ZW5hYmxlQ29udGV4dE1lbnVzfSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xuaW1wb3J0IHtyZWxvYWRXTUx9IGZyb20gXCIuL3dtbFwiO1xuXG53aW5kb3cud2FpbHMgPSB7XG4gICAgLi4ubmV3UnVudGltZSgtMSksXG59O1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy5fd2FpbHMgPSB7XG4gICAgZGlhbG9nQ2FsbGJhY2ssXG4gICAgZGlhbG9nRXJyb3JDYWxsYmFjayxcbiAgICBkaXNwYXRjaEN1c3RvbUV2ZW50LFxuICAgIGNhbGxDYWxsYmFjayxcbiAgICBjYWxsRXJyb3JDYWxsYmFjayxcbn07XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lKGlkKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgQ2xpcGJvYXJkOiB7XG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcbiAgICAgICAgfSxcbiAgICAgICAgQXBwbGljYXRpb246IHtcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uXG4gICAgICAgIH0sXG4gICAgICAgIExvZyxcbiAgICAgICAgU2NyZWVucyxcbiAgICAgICAgQ2FsbCxcbiAgICAgICAgUGx1Z2luLFxuICAgICAgICBXTUw6IHtcbiAgICAgICAgICAgIFJlbG9hZDogcmVsb2FkV01MLFxuICAgICAgICB9LFxuICAgICAgICBEaWFsb2c6IHtcbiAgICAgICAgICAgIEluZm8sXG4gICAgICAgICAgICBXYXJuaW5nLFxuICAgICAgICAgICAgRXJyb3IsXG4gICAgICAgICAgICBRdWVzdGlvbixcbiAgICAgICAgICAgIE9wZW5GaWxlLFxuICAgICAgICAgICAgU2F2ZUZpbGUsXG4gICAgICAgIH0sXG4gICAgICAgIEV2ZW50czoge1xuICAgICAgICAgICAgRW1pdCxcbiAgICAgICAgICAgIE9uLFxuICAgICAgICAgICAgT25jZSxcbiAgICAgICAgICAgIE9uTXVsdGlwbGUsXG4gICAgICAgICAgICBPZmYsXG4gICAgICAgICAgICBPZmZBbGwsXG4gICAgICAgIH0sXG4gICAgICAgIFdpbmRvdzogbmV3V2luZG93KGlkKSxcbiAgICB9O1xufVxuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbmVuYWJsZUNvbnRleHRNZW51cyh0cnVlKTtcblxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcihcIkRPTUNvbnRlbnRMb2FkZWRcIiwgZnVuY3Rpb24oZXZlbnQpIHtcbiAgICByZWxvYWRXTUwoKTtcbn0pOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDWUEsTUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBRTVDLFdBQVMsWUFBWSxRQUFRLE1BQU07QUFDL0IsUUFBSSxNQUFNLElBQUksSUFBSSxVQUFVO0FBQzVCLFFBQUksYUFBYSxPQUFPLFVBQVUsTUFBTTtBQUN4QyxRQUFHLE1BQU07QUFDTCxVQUFJLGFBQWEsT0FBTyxRQUFRLEtBQUssVUFBVSxJQUFJLENBQUM7QUFBQSxJQUN4RDtBQUNBLFdBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLFlBQU0sR0FBRyxFQUNKLEtBQUssY0FBWTtBQUNkLFlBQUksU0FBUyxJQUFJO0FBRWIsY0FBSSxTQUFT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UyxZQUFZO0FBQ3hCLHlCQUFxQjtBQUNyQiwwQkFBc0I7QUFBQSxFQUMxQjs7O0FDbkRBLFNBQU8sUUFBUTtBQUFBLElBQ1gsR0FBRyxXQUFXLEVBQUU7QUFBQSxFQUNwQjtBQUdBLFNBQU8sU0FBUztBQUFBLElBQ1o7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsRUFDSjtBQUVPLFdBQVMsV0FBVyxJQUFJO0FBQzNCLFdBQU87QUFBQSxNQUNILFdBQVc7QUFBQSxRQUNQLEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQSxhQUFhO0FBQUEsUUFDVCxHQUFHO0FBQUEsTUFDUDtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFFBQVE7QUFBQSxNQUNaO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxFQUFFO0FBQUEsSUFDeEI7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEscUJBQW1CLElBQUk7QUFFdkIsV0FBUyxpQkFBaUIsb0JBQW9CLFNBQVMsT0FBTztBQUMxRCxjQUFVO0FBQUEsRUFDZCxDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJldmVudE5hbWUiLCAiY2FsbCIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciJdCn0K
