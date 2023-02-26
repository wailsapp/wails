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
  function dialog(type, options) {
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
    return dialog("Call", options);
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
  function dialog2(type, options) {
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
    return dialog2("Info", options);
  }
  function Warning(options) {
    return dialog2("Warning", options);
  }
  function Error2(options) {
    return dialog2("Error", options);
  }
  function Question(options) {
    return dialog2("Question", options);
  }
  function OpenFile(options) {
    return dialog2("OpenFile", options);
  }
  function SaveFile(options) {
    return dialog2("SaveFile", options);
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChcIlRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIGlmKGFyZ3MpIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcbiAgICB9XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZmV0Y2godXJsKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgaWQpIHtcbiAgICBpZiAoIWlkIHx8IGlkID09PSAtMSkge1xuICAgICAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncykge1xuICAgICAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICAgICAgfTtcbiAgICB9XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgYXJncyA9IGFyZ3MgfHwge307XG4gICAgICAgIGFyZ3NbXCJ3aW5kb3dJRFwiXSA9IGlkO1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGwob2JqZWN0ICsgXCIuXCIgKyBtZXRob2QsIGFyZ3MpO1xuICAgIH1cbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImFwcGxpY2F0aW9uXCIpO1xuXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcbiAgICByZXR1cm4gY2FsbChcIkhpZGVcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xuICAgIHJldHVybiBjYWxsKFwiU2hvd1wiKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJRdWl0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwibG9nXCIpO1xuXG4vKipcbiAqIExvZ3MgYSBtZXNzYWdlLlxuICogQHBhcmFtIHttZXNzYWdlfSBNZXNzYWdlIHRvIGxvZ1xuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nKG1lc3NhZ2UpIHtcbiAgICByZXR1cm4gY2FsbChcIkxvZ1wiLCBtZXNzYWdlKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwic2NyZWVuc1wiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldEFsbFwiKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRQcmltYXJ5XCIpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldEN1cnJlbnRcIik7XG59IiwgImxldCB1cmxBbHBoYWJldCA9XG4gICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuZXhwb3J0IGxldCBjdXN0b21BbHBoYWJldCA9IChhbHBoYWJldCwgZGVmYXVsdFNpemUgPSAyMSkgPT4ge1xuICByZXR1cm4gKHNpemUgPSBkZWZhdWx0U2l6ZSkgPT4ge1xuICAgIGxldCBpZCA9ICcnXG4gICAgbGV0IGkgPSBzaXplXG4gICAgd2hpbGUgKGktLSkge1xuICAgICAgaWQgKz0gYWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiBhbHBoYWJldC5sZW5ndGgpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG4gIH1cbn1cbmV4cG9ydCBsZXQgbmFub2lkID0gKHNpemUgPSAyMSkgPT4ge1xuICBsZXQgaWQgPSAnJ1xuICBsZXQgaSA9IHNpemVcbiAgd2hpbGUgKGktLSkge1xuICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgfVxuICByZXR1cm4gaWRcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjYWxsXCIpO1xuXG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5leHBvcnQgZnVuY3Rpb24gY2FsbEVycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gZGlhbG9nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJDYWxsXCIsIG9wdGlvbnMpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdXaW5kb3coaWQpIHtcbiAgICBsZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJ3aW5kb3dcIiwgaWQpO1xuICAgIHJldHVybiB7XG4gICAgICAgIC8vIFJlbG9hZDogKCkgPT4gY2FsbCgnV1InKSxcbiAgICAgICAgLy8gUmVsb2FkQXBwOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBTZXRTeXN0ZW1EZWZhdWx0VGhlbWU6ICgpID0+IGNhbGwoJ1dBU0RUJyksXG4gICAgICAgIC8vIFNldExpZ2h0VGhlbWU6ICgpID0+IGNhbGwoJ1dBTFQnKSxcbiAgICAgICAgLy8gU2V0RGFya1RoZW1lOiAoKSA9PiBjYWxsKCdXQURUJyksXG4gICAgICAgIC8vIElzRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnV0lGJyksXG4gICAgICAgIC8vIElzTWF4aW1pemVkOiAoKSA9PiBjYWxsKCdXSU0nKSxcbiAgICAgICAgLy8gSXNNaW5pbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTU4nKSxcbiAgICAgICAgLy8gSXNXaW5kb3dlZDogKCkgPT4gY2FsbCgnV0lGJyksXG4gICAgICAgIENlbnRlcjogKCkgPT4gY2FsbCgnQ2VudGVyJyksXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IGNhbGwoJ1NldFRpdGxlJywge3RpdGxlfSksXG4gICAgICAgIEZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ0Z1bGxzY3JlZW4nKSxcbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdVbkZ1bGxzY3JlZW4nKSxcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoJ1NldFNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG4gICAgICAgIFNpemU6ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1NpemUnKSB9LFxuICAgICAgICBTZXRNYXhTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWF4U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0TWluU2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoJ1NldE1pblNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAoYikgPT4gY2FsbCgnU2V0QWx3YXlzT25Ub3AnLCB7YWx3YXlzT25Ub3A6Yn0pLFxuICAgICAgICBTZXRQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoJ1NldFBvc2l0aW9uJywge3gseX0pLFxuICAgICAgICBQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUG9zaXRpb24nKSB9LFxuICAgICAgICBTY3JlZW46ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1NjcmVlbicpIH0sXG4gICAgICAgIEhpZGU6ICgpID0+IGNhbGwoJ0hpZGUnKSxcbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IGNhbGwoJ01heGltaXNlJyksXG4gICAgICAgIFNob3c6ICgpID0+IGNhbGwoJ1Nob3cnKSxcbiAgICAgICAgQ2xvc2U6ICgpID0+IGNhbGwoJ0Nsb3NlJyksXG4gICAgICAgIFRvZ2dsZU1heGltaXNlOiAoKSA9PiBjYWxsKCdUb2dnbGVNYXhpbWlzZScpLFxuICAgICAgICBVbk1heGltaXNlOiAoKSA9PiBjYWxsKCdVbk1heGltaXNlJyksXG4gICAgICAgIE1pbmltaXNlOiAoKSA9PiBjYWxsKCdNaW5pbWlzZScpLFxuICAgICAgICBVbk1pbmltaXNlOiAoKSA9PiBjYWxsKCdVbk1pbmltaXNlJyksXG4gICAgICAgIFNldEJhY2tncm91bmRDb2xvdXI6IChyLCBnLCBiLCBhKSA9PiBjYWxsKCdTZXRCYWNrZ3JvdW5kQ29sb3VyJywge3IsIGcsIGIsIGF9KSxcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImV2ZW50c1wiKTtcblxuLyoqXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcbiAqXG4gKiBAY2xhc3MgTGlzdGVuZXJcbiAqL1xuY2xhc3MgTGlzdGVuZXIge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICAgICAqIEBtZW1iZXJvZiBMaXN0ZW5lclxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xuICAgICAgICAgICAgY2FsbGJhY2soZGF0YSk7XG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cblxuLyoqXG4gKiBDdXN0b21FdmVudCBkZWZpbmVzIGEgY3VzdG9tIGV2ZW50LiBJdCBpcyBwYXNzZWQgdG8gZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBjbGFzcyBDdXN0b21FdmVudFxuICovXG5leHBvcnQgY2xhc3MgQ3VzdG9tRXZlbnQge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgQ3VzdG9tRXZlbnQuXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IG5hbWUgLSBOYW1lIG9mIHRoZSBldmVudFxuICAgICAqIEBwYXJhbSB7YW55fSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XG4gICAgICogQG1lbWJlcm9mIEN1c3RvbUV2ZW50XG4gICAgICovXG4gICAgY29uc3RydWN0b3IobmFtZSwgZGF0YSkge1xuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihDdXN0b21FdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihDdXN0b21FdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb24oQ3VzdG9tRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG4vKipcbiAqIGxpc3RlbmVyT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cbiAqXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxuICovXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuXG4vKipcbiAqIGRpc3BhdGNoZXMgYW4gZXZlbnQgdG8gYWxsIGxpc3RlbmVyc1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7Q3VzdG9tRXZlbnR9IGV2ZW50XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkaXNwYXRjaEN1c3RvbUV2ZW50KGV2ZW50KSB7XG4gICAgY29uc29sZS5sb2coXCJkaXNwYXRjaGluZyBldmVudDogXCIsIHtldmVudH0pO1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XG4gICAgaWYgKGxpc3RlbmVycykge1xuICAgICAgICAvLyBpdGVyYXRlIGxpc3RlbmVycyBhbmQgY2FsbCBjYWxsYmFjay4gSWYgY2FsbGJhY2sgcmV0dXJucyB0cnVlLCByZW1vdmUgbGlzdGVuZXJcbiAgICAgICAgbGV0IHRvUmVtb3ZlID0gW107XG4gICAgICAgIGxpc3RlbmVycy5mb3JFYWNoKGxpc3RlbmVyID0+IHtcbiAgICAgICAgICAgIGxldCByZW1vdmUgPSBsaXN0ZW5lci5DYWxsYmFjayhldmVudClcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHtcbiAgICAgICAgICAgICAgICB0b1JlbW92ZS5wdXNoKGxpc3RlbmVyKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgICAgIC8vIHJlbW92ZSBsaXN0ZW5lcnNcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xuICAgICAgICAgICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCBsaXN0ZW5lcnMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcnMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxuICpcbiBbdjMgQ0hBTkdFXSBPZmYgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSlcbn1cblxuLyoqXG4gKiBPZmZBbGwgdW5yZWdpc3RlcnMgYWxsIGxpc3RlbmVyc1xuICogW3YzIENIQU5HRV0gT2ZmQWxsIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcbiAqXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLypcbiAgIEVtaXQgZW1pdHMgYW4gZXZlbnQgdG8gYWxsIGxpc3RlbmVyc1xuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkge1xuICAgIHJldHVybiBjYWxsKFwiRW1pdFwiLCBldmVudCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSdcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZGlhbG9nXCIpO1xuXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zKSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pO1xuICAgIH0pO1xufVxuXG5cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiSW5mb1wiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdhcm5pbmcob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJXYXJuaW5nXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gRXJyb3Iob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJFcnJvclwiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiUXVlc3Rpb25cIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIk9wZW5GaWxlXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJTYXZlRmlsZVwiLCBvcHRpb25zKTtcbn1cblxuIiwgImltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjb250ZXh0bWVudVwiKTtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkLCB4LCB5LCBkYXRhKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJPcGVuQ29udGV4dE1lbnVcIiwge2lkLCB4LCB5LCBkYXRhfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBlbmFibGVDb250ZXh0TWVudXMoZW5hYmxlZCkge1xuICAgIGlmIChlbmFibGVkKSB7XG4gICAgICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgd2luZG93LnJlbW92ZUV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudCkge1xuICAgIHByb2Nlc3NDb250ZXh0TWVudShldmVudC50YXJnZXQsIGV2ZW50KTtcbn1cblxuZnVuY3Rpb24gcHJvY2Vzc0NvbnRleHRNZW51KGVsZW1lbnQsIGV2ZW50KSB7XG4gICAgbGV0IGlkID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtY29udGV4dG1lbnUnKTtcbiAgICBpZiAoaWQpIHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGlkLCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS1jb250ZXh0bWVudS1kYXRhJykpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGxldCBwYXJlbnQgPSBlbGVtZW50LnBhcmVudEVsZW1lbnQ7XG4gICAgICAgIGlmIChwYXJlbnQpIHtcbiAgICAgICAgICAgIHByb2Nlc3NDb250ZXh0TWVudShwYXJlbnQsIGV2ZW50KTtcbiAgICAgICAgfVxuICAgIH1cbn1cbiIsICJcbmltcG9ydCB7RW1pdH0gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQge1F1ZXN0aW9ufSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5cbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudCkge1xuICAgbGV0IF8gPSBFbWl0KHtuYW1lOiBldmVudH0gKTtcbn1cblxuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtZXZlbnRdJyk7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCBlbGVtZW50cy5sZW5ndGg7IGkrKykge1xuICAgICAgICBjb25zdCBlbGVtZW50ID0gZWxlbWVudHNbaV07XG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIH1cbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xuXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcbiAgICBpZiAod2FpbHMuV2luZG93W21ldGhvZF0gPT09IHVuZGVmaW5lZCkge1xuICAgICAgICBjb25zb2xlLmxvZyhcIldpbmRvdyBtZXRob2QgXCIgKyBtZXRob2QgKyBcIiBub3QgZm91bmRcIik7XG4gICAgfVxuICAgIHdhaWxzLldpbmRvd1ttZXRob2RdKCk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC13aW5kb3ddJyk7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCBlbGVtZW50cy5sZW5ndGg7IGkrKykge1xuICAgICAgICBjb25zdCBlbGVtZW50ID0gZWxlbWVudHNbaV07XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC13aW5kb3cnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHdpbmRvd01ldGhvZCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHdpbmRvd01ldGhvZCk7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHJlbG9hZFdNTCgpIHtcbiAgICBhZGRXTUxFdmVudExpc3RlbmVycygpO1xuICAgIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tICcuL2NsaXBib2FyZCc7XG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tICcuL2FwcGxpY2F0aW9uJztcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gJy4vc2NyZWVucyc7XG5pbXBvcnQge0NhbGwsIGNhbGxFcnJvckNhbGxiYWNrLCBjYWxsQ2FsbGJhY2t9IGZyb20gXCIuL2NhbGxzXCI7XG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQge2Rpc3BhdGNoQ3VzdG9tRXZlbnQsIEVtaXQsIE9mZiwgT2ZmQWxsLCBPbiwgT25jZSwgT25NdWx0aXBsZX0gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQge2RpYWxvZ0NhbGxiYWNrLCBkaWFsb2dFcnJvckNhbGxiYWNrLCBFcnJvciwgSW5mbywgT3BlbkZpbGUsIFF1ZXN0aW9uLCBTYXZlRmlsZSwgV2FybmluZyx9IGZyb20gXCIuL2RpYWxvZ3NcIjtcbmltcG9ydCB7ZW5hYmxlQ29udGV4dE1lbnVzfSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xuaW1wb3J0IHtyZWxvYWRXTUx9IGZyb20gXCIuL3dtbFwiO1xuXG53aW5kb3cud2FpbHMgPSB7XG4gICAgLi4ubmV3UnVudGltZSgtMSksXG59O1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy5fd2FpbHMgPSB7XG4gICAgZGlhbG9nQ2FsbGJhY2ssXG4gICAgZGlhbG9nRXJyb3JDYWxsYmFjayxcbiAgICBkaXNwYXRjaEN1c3RvbUV2ZW50LFxuICAgIGNhbGxDYWxsYmFjayxcbiAgICBjYWxsRXJyb3JDYWxsYmFjayxcbn07XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lKGlkKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgQ2xpcGJvYXJkOiB7XG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcbiAgICAgICAgfSxcbiAgICAgICAgQXBwbGljYXRpb246IHtcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uXG4gICAgICAgIH0sXG4gICAgICAgIExvZyxcbiAgICAgICAgU2NyZWVucyxcbiAgICAgICAgQ2FsbCxcbiAgICAgICAgV01MOiB7XG4gICAgICAgICAgICBSZWxvYWQ6IHJlbG9hZFdNTCxcbiAgICAgICAgfSxcbiAgICAgICAgRGlhbG9nOiB7XG4gICAgICAgICAgICBJbmZvLFxuICAgICAgICAgICAgV2FybmluZyxcbiAgICAgICAgICAgIEVycm9yLFxuICAgICAgICAgICAgUXVlc3Rpb24sXG4gICAgICAgICAgICBPcGVuRmlsZSxcbiAgICAgICAgICAgIFNhdmVGaWxlLFxuICAgICAgICB9LFxuICAgICAgICBFdmVudHM6IHtcbiAgICAgICAgICAgIEVtaXQsXG4gICAgICAgICAgICBPbixcbiAgICAgICAgICAgIE9uY2UsXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxuICAgICAgICAgICAgT2ZmLFxuICAgICAgICAgICAgT2ZmQWxsLFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyhpZCksXG4gICAgfVxufVxuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbmVuYWJsZUNvbnRleHRNZW51cyh0cnVlKTtcblxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcihcIkRPTUNvbnRlbnRMb2FkZWRcIiwgZnVuY3Rpb24oZXZlbnQpIHtcbiAgICByZWxvYWRXTUwoKTtcbn0pOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDWUEsTUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBRTVDLFdBQVMsWUFBWSxRQUFRLE1BQU07QUFDL0IsUUFBSSxNQUFNLElBQUksSUFBSSxVQUFVO0FBQzVCLFFBQUksYUFBYSxPQUFPLFVBQVUsTUFBTTtBQUN4QyxRQUFHLE1BQU07QUFDTCxVQUFJLGFBQWEsT0FBTyxRQUFRLEtBQUssVUFBVSxJQUFJLENBQUM7QUFBQSxJQUN4RDtBQUNBLFdBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLFlBQU0sR0FBRyxFQUNKLEtBQUssY0FBWTtBQUNkLFlBQUksU0FBUyxJQUFJO0FBRWIsY0FBSSxTQUFTLFFBQVEsSUFBSSxjQUFjLEtBQUssU0FBUyxRQUFRLElBQUksY0FBYyxFQUFFLFFBQVEsa0JBQWtCLE1BQU0sSUFBSTtBQUNqSCxtQkFBTyxTQUFTLEtBQUs7QUFBQSxVQUN6QixPQUFPO0FBQ0gsbUJBQU8sU0FB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Q0o7QUFFTyxXQUFTLFdBQVcsSUFBSTtBQUMzQixXQUFPO0FBQUEsTUFDSCxXQUFXO0FBQUEsUUFDUCxHQUFHO0FBQUEsTUFDUDtBQUFBLE1BQ0EsYUFBYTtBQUFBLFFBQ1QsR0FBRztBQUFBLE1BQ1A7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFFBQVE7QUFBQSxNQUNaO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxFQUFFO0FBQUEsSUFDeEI7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEscUJBQW1CLElBQUk7QUFFdkIsV0FBUyxpQkFBaUIsb0JBQW9CLFNBQVMsT0FBTztBQUMxRCxjQUFVO0FBQUEsRUFDZCxDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJldmVudE5hbWUiLCAiY2FsbCIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiZGlhbG9nIiwgIkVycm9yIiwgImNhbGwiLCAiRXJyb3IiXQp9Cg==
