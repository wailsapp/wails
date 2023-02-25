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

  // desktop/window.js
  function newWindow(id) {
    let call8 = newRuntimeCaller("window", id);
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
      Center: () => call8("Center"),
      SetTitle: (title) => call8("SetTitle", { title }),
      Fullscreen: () => call8("Fullscreen"),
      UnFullscreen: () => call8("UnFullscreen"),
      SetSize: (width, height) => call8("SetSize", { width, height }),
      Size: () => {
        return call8("Size");
      },
      SetMaxSize: (width, height) => call8("SetMaxSize", { width, height }),
      SetMinSize: (width, height) => call8("SetMinSize", { width, height }),
      SetAlwaysOnTop: (b) => call8("SetAlwaysOnTop", { alwaysOnTop: b }),
      SetPosition: (x, y) => call8("SetPosition", { x, y }),
      Position: () => {
        return call8("Position");
      },
      Screen: () => {
        return call8("Screen");
      },
      Hide: () => call8("Hide"),
      Maximise: () => call8("Maximise"),
      Show: () => call8("Show"),
      Close: () => call8("Close"),
      ToggleMaximise: () => call8("ToggleMaximise"),
      UnMaximise: () => call8("UnMaximise"),
      Minimise: () => call8("Minimise"),
      UnMinimise: () => call8("UnMinimise"),
      SetBackgroundColour: (r, g, b, a) => call8("SetBackgroundColour", { r, g, b, a })
    };
  }

  // desktop/events.js
  var call5 = newRuntimeCaller("events");
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
    return call5("Emit", event);
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

  // desktop/dialogs.js
  var call6 = newRuntimeCaller("dialog");
  var dialogResponses = /* @__PURE__ */ new Map();
  function generateID() {
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
      let id = generateID();
      options = options || {};
      options["dialog-id"] = id;
      dialogResponses.set(id, { resolve, reject });
      call6(type, options).catch((error) => {
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
  var call7 = newRuntimeCaller("contextmenu");
  function openContextMenu(id, x, y, data) {
    return call7("OpenContextMenu", { id, x, y, data });
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
    dispatchCustomEvent
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2RpYWxvZ3MuanMiLCAiZGVza3RvcC9jb250ZXh0bWVudS5qcyIsICJkZXNrdG9wL3dtbC5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjbGlwYm9hcmRcIik7XG5cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gY2FsbChcIlNldFRleHRcIiwge3RleHR9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBpZihhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybClcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3QsIGlkKSB7XG4gICAgaWYgKCFpZCB8fCBpZCA9PT0gLTEpIHtcbiAgICAgICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgYXJncyk7XG4gICAgICAgIH07XG4gICAgfVxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzKSB7XG4gICAgICAgIGFyZ3MgPSBhcmdzIHx8IHt9O1xuICAgICAgICBhcmdzW1wid2luZG93SURcIl0gPSBpZDtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICB9XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJhcHBsaWNhdGlvblwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJIaWRlXCIpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChcIlNob3dcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHJldHVybiBjYWxsKFwiUXVpdFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImxvZ1wiKTtcblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZS5cbiAqIEBwYXJhbSB7bWVzc2FnZX0gTWVzc2FnZSB0byBsb2dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZyhtZXNzYWdlKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJMb2dcIiwgbWVzc2FnZSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcInNjcmVlbnNcIik7XG5cbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRBbGxcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0UHJpbWFyeVwiKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRDdXJyZW50XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1dpbmRvdyhpZCkge1xuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCBpZCk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBSZWxvYWRBcHA6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFNldFN5c3RlbURlZmF1bHRUaGVtZTogKCkgPT4gY2FsbCgnV0FTRFQnKSxcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxuICAgICAgICAvLyBTZXREYXJrVGhlbWU6ICgpID0+IGNhbGwoJ1dBRFQnKSxcbiAgICAgICAgLy8gSXNGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxuICAgICAgICAvLyBJc01pbmltaXplZDogKCkgPT4gY2FsbCgnV0lNTicpLFxuICAgICAgICAvLyBJc1dpbmRvd2VkOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKCdDZW50ZXInKSxcbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbCgnU2V0VGl0bGUnLCB7dGl0bGV9KSxcbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnRnVsbHNjcmVlbicpLFxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2l6ZTogKCkgPT4geyByZXR1cm4gY2FsbCgnU2l6ZScpIH0sXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRNYXhTaXplJywge3dpZHRoLGhlaWdodH0pLFxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWluU2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0QWx3YXlzT25Ub3A6IChiKSA9PiBjYWxsKCdTZXRBbHdheXNPblRvcCcsIHthbHdheXNPblRvcDpifSksXG4gICAgICAgIFNldFBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbCgnU2V0UG9zaXRpb24nLCB7eCx5fSksXG4gICAgICAgIFBvc2l0aW9uOiAoKSA9PiB7IHJldHVybiBjYWxsKCdQb3NpdGlvbicpIH0sXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJykgfSxcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbCgnSGlkZScpLFxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gY2FsbCgnTWF4aW1pc2UnKSxcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbCgnU2hvdycpLFxuICAgICAgICBDbG9zZTogKCkgPT4gY2FsbCgnQ2xvc2UnKSxcbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1RvZ2dsZU1heGltaXNlJyksXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1VuTWF4aW1pc2UnKSxcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IGNhbGwoJ01pbmltaXNlJyksXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwoJ1VuTWluaW1pc2UnKSxcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IGNhbGwoJ1NldEJhY2tncm91bmRDb2xvdXInLCB7ciwgZywgYiwgYX0pLFxuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZXZlbnRzXCIpO1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuXG4vKipcbiAqIEN1c3RvbUV2ZW50IGRlZmluZXMgYSBjdXN0b20gZXZlbnQuIEl0IGlzIHBhc3NlZCB0byBldmVudCBsaXN0ZW5lcnMuXG4gKlxuICogQGNsYXNzIEN1c3RvbUV2ZW50XG4gKi9cbmV4cG9ydCBjbGFzcyBDdXN0b21FdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBDdXN0b21FdmVudC5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gICAgICogQHBhcmFtIHthbnl9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcbiAgICAgKiBAbWVtYmVyb2YgQ3VzdG9tRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihDdXN0b21FdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbmZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG5cbi8qKlxuICogZGlzcGF0Y2hlcyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtDdXN0b21FdmVudH0gZXZlbnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGRpc3BhdGNoQ3VzdG9tRXZlbnQoZXZlbnQpIHtcbiAgICBjb25zb2xlLmxvZyhcImRpc3BhdGNoaW5nIGV2ZW50OiBcIiwge2V2ZW50fSk7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBbXTtcbiAgICAgICAgbGlzdGVuZXJzLmZvckVhY2gobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KVxuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xuICAgICAgICBpZiAodG9SZW1vdmUubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuIFt2MyBDSEFOR0VdIE9mZiBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9KVxufVxuXG4vKipcbiAqIE9mZkFsbCB1bnJlZ2lzdGVycyBhbGwgbGlzdGVuZXJzXG4gKiBbdjMgQ0hBTkdFXSBPZmZBbGwgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpIHtcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xufVxuXG4vKlxuICAgRW1pdCBlbWl0cyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XG4gICAgcmV0dXJuIGNhbGwoXCJFbWl0XCIsIGV2ZW50KTtcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSdcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZGlhbG9nXCIpO1xuXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zKSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pXG4gICAgfSk7XG59XG5cblxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJJbmZvXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIldhcm5pbmdcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkVycm9yXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJRdWVzdGlvblwiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiT3BlbkZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlNhdmVGaWxlXCIsIG9wdGlvbnMpO1xufVxuXG4iLCAiaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNvbnRleHRtZW51XCIpO1xuXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQsIHgsIHksIGRhdGEpIHtcbiAgICByZXR1cm4gY2FsbChcIk9wZW5Db250ZXh0TWVudVwiLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGVuYWJsZUNvbnRleHRNZW51cyhlbmFibGVkKSB7XG4gICAgaWYgKGVuYWJsZWQpIHtcbiAgICAgICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50KSB7XG4gICAgcHJvY2Vzc0NvbnRleHRNZW51KGV2ZW50LnRhcmdldCwgZXZlbnQpO1xufVxuXG5mdW5jdGlvbiBwcm9jZXNzQ29udGV4dE1lbnUoZWxlbWVudCwgZXZlbnQpIHtcbiAgICBsZXQgaWQgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS1jb250ZXh0bWVudScpO1xuICAgIGlmIChpZCkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoaWQsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51LWRhdGEnKSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgbGV0IHBhcmVudCA9IGVsZW1lbnQucGFyZW50RWxlbWVudDtcbiAgICAgICAgaWYgKHBhcmVudCkge1xuICAgICAgICAgICAgcHJvY2Vzc0NvbnRleHRNZW51KHBhcmVudCwgZXZlbnQpO1xuICAgICAgICB9XG4gICAgfVxufVxuIiwgIlxuaW1wb3J0IHtFbWl0fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcblxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50KSB7XG4gICBsZXQgXyA9IEVtaXQoe25hbWU6IGV2ZW50fSApO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC1ldmVudF0nKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGVsZW1lbnRzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGNvbnN0IGVsZW1lbnQgPSBlbGVtZW50c1tpXTtcbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWV2ZW50Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgfVxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG5cbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKG1ldGhvZCkge1xuICAgIGlmICh3YWlscy5XaW5kb3dbbWV0aG9kXSA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIGNvbnNvbGUubG9nKFwiV2luZG93IG1ldGhvZCBcIiArIG1ldGhvZCArIFwiIG5vdCBmb3VuZFwiKTtcbiAgICB9XG4gICAgd2FpbHMuV2luZG93W21ldGhvZF0oKTtcbn1cblxuZnVuY3Rpb24gYWRkV01MV2luZG93TGlzdGVuZXJzKCkge1xuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW2RhdGEtd21sLXdpbmRvd10nKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGVsZW1lbnRzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGNvbnN0IGVsZW1lbnQgPSBlbGVtZW50c1tpXTtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVsb2FkV01MKCkge1xuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gJy4vY2xpcGJvYXJkJztcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gJy4vYXBwbGljYXRpb24nO1xuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSAnLi9zY3JlZW5zJztcblxuaW1wb3J0IHtuZXdXaW5kb3d9IGZyb20gXCIuL3dpbmRvd1wiO1xuaW1wb3J0IHtkaXNwYXRjaEN1c3RvbUV2ZW50LCBFbWl0LCBPZmYsIE9mZkFsbCwgT24sIE9uY2UsIE9uTXVsdGlwbGV9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtkaWFsb2dDYWxsYmFjaywgZGlhbG9nRXJyb3JDYWxsYmFjaywgRXJyb3IsIEluZm8sIE9wZW5GaWxlLCBRdWVzdGlvbiwgU2F2ZUZpbGUsIFdhcm5pbmcsfSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5pbXBvcnQge2VuYWJsZUNvbnRleHRNZW51c30gZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCB7cmVsb2FkV01MfSBmcm9tIFwiLi93bWxcIjtcblxud2luZG93LndhaWxzID0ge1xuICAgIC4uLm5ld1J1bnRpbWUoLTEpLFxufTtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cuX3dhaWxzID0ge1xuICAgIGRpYWxvZ0NhbGxiYWNrLFxuICAgIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssXG4gICAgZGlzcGF0Y2hDdXN0b21FdmVudCxcbn1cblxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZShpZCkge1xuICAgIHJldHVybiB7XG4gICAgICAgIENsaXBib2FyZDoge1xuICAgICAgICAgICAgLi4uQ2xpcGJvYXJkXG4gICAgICAgIH0sXG4gICAgICAgIEFwcGxpY2F0aW9uOiB7XG4gICAgICAgICAgICAuLi5BcHBsaWNhdGlvblxuICAgICAgICB9LFxuICAgICAgICBMb2csXG4gICAgICAgIFNjcmVlbnMsXG4gICAgICAgIFdNTDoge1xuICAgICAgICAgICAgUmVsb2FkOiByZWxvYWRXTUwsXG4gICAgICAgIH0sXG4gICAgICAgIERpYWxvZzoge1xuICAgICAgICAgICAgSW5mbyxcbiAgICAgICAgICAgIFdhcm5pbmcsXG4gICAgICAgICAgICBFcnJvcixcbiAgICAgICAgICAgIFF1ZXN0aW9uLFxuICAgICAgICAgICAgT3BlbkZpbGUsXG4gICAgICAgICAgICBTYXZlRmlsZSxcbiAgICAgICAgfSxcbiAgICAgICAgRXZlbnRzOiB7XG4gICAgICAgICAgICBFbWl0LFxuICAgICAgICAgICAgT24sXG4gICAgICAgICAgICBPbmNlLFxuICAgICAgICAgICAgT25NdWx0aXBsZSxcbiAgICAgICAgICAgIE9mZixcbiAgICAgICAgICAgIE9mZkFsbCxcbiAgICAgICAgfSxcbiAgICAgICAgV2luZG93OiBuZXdXaW5kb3coaWQpLFxuICAgIH1cbn1cblxuaWYgKERFQlVHKSB7XG4gICAgY29uc29sZS5sb2coXCJXYWlscyB2My4wLjAgRGVidWcgTW9kZSBFbmFibGVkXCIpO1xufVxuXG5lbmFibGVDb250ZXh0TWVudXModHJ1ZSk7XG5cbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoXCJET01Db250ZW50TG9hZGVkXCIsIGZ1bmN0aW9uKGV2ZW50KSB7XG4gICAgcmVsb2FkV01MKCk7XG59KTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ1lBLE1BQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUU1QyxXQUFTLFlBQVksUUFBUSxNQUFNO0FBQy9CLFFBQUksTUFBTSxJQUFJLElBQUksVUFBVTtBQUM1QixRQUFJLGFBQWEsT0FBTyxVQUFVLE1BQU07QUFDeEMsUUFBRyxNQUFNO0FBQ0wsVUFBSSxhQUFhLE9BQU8sUUFBUSxLQUFLLFVBQVUsSUFBSSxDQUFDO0FBQUEsSUFDeEQ7QUFDQSxXQUFP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QVMsUUFBUTtBQUc3QyxjQUFRLGlCQUFpQixTQUFTLFFBQVE7QUFBQSxJQUM5QztBQUFBLEVBQ0o7QUFFTyxXQUFTLFlBQVk7QUFDeEIseUJBQXFCO0FBQ3JCLDBCQUFzQjtBQUFBLEVBQzFCOzs7QUNuREEsU0FBTyxRQUFRO0FBQUEsSUFDWCxHQUFHLFdBQVcsRUFBRTtBQUFBLEVBQ3BCO0FBR0EsU0FBTyxTQUFTO0FBQUEsSUFDWjtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsRUFDSjtBQUdPLFdBQVMsV0FBVyxJQUFJO0FBQzNCLFdBQU87QUFBQSxNQUNILFdBQVc7QUFBQSxRQUNQLEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQSxhQUFhO0FBQUEsUUFDVCxHQUFHO0FBQUEsTUFDUDtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQSxLQUFLO0FBQUEsUUFDRCxRQUFRO0FBQUEsTUFDWjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQSxPQUFBQztBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRLFVBQVUsRUFBRTtBQUFBLElBQ3hCO0FBQUEsRUFDSjtBQUVBLE1BQUksTUFBTztBQUNQLFlBQVEsSUFBSSxpQ0FBaUM7QUFBQSxFQUNqRDtBQUVBLHFCQUFtQixJQUFJO0FBRXZCLFdBQVMsaUJBQWlCLG9CQUFvQixTQUFTLE9BQU87QUFDMUQsY0FBVTtBQUFBLEVBQ2QsQ0FBQzsiLAogICJuYW1lcyI6IFsiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImNhbGwiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciJdCn0K
