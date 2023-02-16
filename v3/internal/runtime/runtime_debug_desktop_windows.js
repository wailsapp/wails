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
  function refresh() {
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
        Refresh: refresh
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
    refresh();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2RpYWxvZ3MuanMiLCAiZGVza3RvcC9jb250ZXh0bWVudS5qcyIsICJkZXNrdG9wL3dtbC5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjbGlwYm9hcmRcIik7XG5cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gY2FsbChcIlNldFRleHRcIiwge3RleHR9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBpZihhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybClcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3QsIGlkKSB7XG4gICAgaWYgKCFpZCB8fCBpZCA9PT0gLTEpIHtcbiAgICAgICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgYXJncyk7XG4gICAgICAgIH07XG4gICAgfVxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzKSB7XG4gICAgICAgIGFyZ3MgPSBhcmdzIHx8IHt9O1xuICAgICAgICBhcmdzW1wid2luZG93SURcIl0gPSBpZDtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICB9XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJhcHBsaWNhdGlvblwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJIaWRlXCIpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChcIlNob3dcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHJldHVybiBjYWxsKFwiUXVpdFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImxvZ1wiKTtcblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZS5cbiAqIEBwYXJhbSB7bWVzc2FnZX0gTWVzc2FnZSB0byBsb2dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZyhtZXNzYWdlKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJMb2dcIiwgbWVzc2FnZSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcInNjcmVlbnNcIik7XG5cbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRBbGxcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0UHJpbWFyeVwiKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRDdXJyZW50XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1dpbmRvdyhpZCkge1xuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCBpZCk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBSZWxvYWRBcHA6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFNldFN5c3RlbURlZmF1bHRUaGVtZTogKCkgPT4gY2FsbCgnV0FTRFQnKSxcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxuICAgICAgICAvLyBTZXREYXJrVGhlbWU6ICgpID0+IGNhbGwoJ1dBRFQnKSxcbiAgICAgICAgLy8gSXNGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxuICAgICAgICAvLyBJc01pbmltaXplZDogKCkgPT4gY2FsbCgnV0lNTicpLFxuICAgICAgICAvLyBJc1dpbmRvd2VkOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKCdDZW50ZXInKSxcbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbCgnU2V0VGl0bGUnLCB7dGl0bGV9KSxcbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnRnVsbHNjcmVlbicpLFxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2l6ZTogKCkgPT4geyByZXR1cm4gY2FsbCgnU2l6ZScpIH0sXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRNYXhTaXplJywge3dpZHRoLGhlaWdodH0pLFxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWluU2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0QWx3YXlzT25Ub3A6IChiKSA9PiBjYWxsKCdTZXRBbHdheXNPblRvcCcsIHthbHdheXNPblRvcDpifSksXG4gICAgICAgIFNldFBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbCgnU2V0UG9zaXRpb24nLCB7eCx5fSksXG4gICAgICAgIFBvc2l0aW9uOiAoKSA9PiB7IHJldHVybiBjYWxsKCdQb3NpdGlvbicpIH0sXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJykgfSxcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbCgnSGlkZScpLFxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gY2FsbCgnTWF4aW1pc2UnKSxcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbCgnU2hvdycpLFxuICAgICAgICBDbG9zZTogKCkgPT4gY2FsbCgnQ2xvc2UnKSxcbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1RvZ2dsZU1heGltaXNlJyksXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1VuTWF4aW1pc2UnKSxcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IGNhbGwoJ01pbmltaXNlJyksXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwoJ1VuTWluaW1pc2UnKSxcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IGNhbGwoJ1NldEJhY2tncm91bmRDb2xvdXInLCB7ciwgZywgYiwgYX0pLFxuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZXZlbnRzXCIpO1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuXG4vKipcbiAqIEN1c3RvbUV2ZW50IGRlZmluZXMgYSBjdXN0b20gZXZlbnQuIEl0IGlzIHBhc3NlZCB0byBldmVudCBsaXN0ZW5lcnMuXG4gKlxuICogQGNsYXNzIEN1c3RvbUV2ZW50XG4gKi9cbmV4cG9ydCBjbGFzcyBDdXN0b21FdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBDdXN0b21FdmVudC5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gICAgICogQHBhcmFtIHthbnl9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcbiAgICAgKiBAbWVtYmVyb2YgQ3VzdG9tRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihDdXN0b21FdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbmZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG5cbi8qKlxuICogZGlzcGF0Y2hlcyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtDdXN0b21FdmVudH0gZXZlbnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGRpc3BhdGNoQ3VzdG9tRXZlbnQoZXZlbnQpIHtcbiAgICBjb25zb2xlLmxvZyhcImRpc3BhdGNoaW5nIGV2ZW50OiBcIiwge2V2ZW50fSk7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBbXTtcbiAgICAgICAgbGlzdGVuZXJzLmZvckVhY2gobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KVxuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xuICAgICAgICBpZiAodG9SZW1vdmUubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuIFt2MyBDSEFOR0VdIE9mZiBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9KVxufVxuXG4vKipcbiAqIE9mZkFsbCB1bnJlZ2lzdGVycyBhbGwgbGlzdGVuZXJzXG4gKiBbdjMgQ0hBTkdFXSBPZmZBbGwgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpIHtcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xufVxuXG4vKlxuICAgRW1pdCBlbWl0cyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XG4gICAgcmV0dXJuIGNhbGwoXCJFbWl0XCIsIGV2ZW50KTtcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSdcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZGlhbG9nXCIpO1xuXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zKSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pXG4gICAgfSk7XG59XG5cblxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJJbmZvXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIldhcm5pbmdcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkVycm9yXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJRdWVzdGlvblwiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiT3BlbkZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlNhdmVGaWxlXCIsIG9wdGlvbnMpO1xufVxuXG4iLCAiaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNvbnRleHRtZW51XCIpO1xuXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQsIHgsIHksIGRhdGEpIHtcbiAgICByZXR1cm4gY2FsbChcIk9wZW5Db250ZXh0TWVudVwiLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGVuYWJsZUNvbnRleHRNZW51cyhlbmFibGVkKSB7XG4gICAgaWYgKGVuYWJsZWQpIHtcbiAgICAgICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50KSB7XG4gICAgcHJvY2Vzc0NvbnRleHRNZW51KGV2ZW50LnRhcmdldCwgZXZlbnQpO1xufVxuXG5mdW5jdGlvbiBwcm9jZXNzQ29udGV4dE1lbnUoZWxlbWVudCwgZXZlbnQpIHtcbiAgICBsZXQgaWQgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS1jb250ZXh0bWVudScpO1xuICAgIGlmIChpZCkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoaWQsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51LWRhdGEnKSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgbGV0IHBhcmVudCA9IGVsZW1lbnQucGFyZW50RWxlbWVudDtcbiAgICAgICAgaWYgKHBhcmVudCkge1xuICAgICAgICAgICAgcHJvY2Vzc0NvbnRleHRNZW51KHBhcmVudCwgZXZlbnQpO1xuICAgICAgICB9XG4gICAgfVxufVxuIiwgIlxuaW1wb3J0IHtFbWl0fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcblxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50KSB7XG4gICBsZXQgXyA9IEVtaXQoe25hbWU6IGV2ZW50fSApO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC1ldmVudF0nKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGVsZW1lbnRzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGNvbnN0IGVsZW1lbnQgPSBlbGVtZW50c1tpXTtcbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWV2ZW50Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgfVxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG5cbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKG1ldGhvZCkge1xuICAgIGlmICh3YWlscy5XaW5kb3dbbWV0aG9kXSA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIGNvbnNvbGUubG9nKFwiV2luZG93IG1ldGhvZCBcIiArIG1ldGhvZCArIFwiIG5vdCBmb3VuZFwiKTtcbiAgICB9XG4gICAgd2FpbHMuV2luZG93W21ldGhvZF0oKTtcbn1cblxuZnVuY3Rpb24gYWRkV01MV2luZG93TGlzdGVuZXJzKCkge1xuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW2RhdGEtd21sLXdpbmRvd10nKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGVsZW1lbnRzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgIGNvbnN0IGVsZW1lbnQgPSBlbGVtZW50c1tpXTtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVmcmVzaCgpIHtcbiAgICBhZGRXTUxFdmVudExpc3RlbmVycygpO1xuICAgIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tICcuL2NsaXBib2FyZCc7XG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tICcuL2FwcGxpY2F0aW9uJztcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gJy4vc2NyZWVucyc7XG5cbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCB7ZGlzcGF0Y2hDdXN0b21FdmVudCwgRW1pdCwgT2ZmLCBPZmZBbGwsIE9uLCBPbmNlLCBPbk11bHRpcGxlfSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7ZGlhbG9nQ2FsbGJhY2ssIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssIEVycm9yLCBJbmZvLCBPcGVuRmlsZSwgUXVlc3Rpb24sIFNhdmVGaWxlLCBXYXJuaW5nLH0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuaW1wb3J0IHtlbmFibGVDb250ZXh0TWVudXN9IGZyb20gXCIuL2NvbnRleHRtZW51XCI7XG5pbXBvcnQge3JlZnJlc2h9IGZyb20gXCIuL3dtbFwiO1xuXG53aW5kb3cud2FpbHMgPSB7XG4gICAgLi4ubmV3UnVudGltZSgtMSksXG59O1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy5fd2FpbHMgPSB7XG4gICAgZGlhbG9nQ2FsbGJhY2ssXG4gICAgZGlhbG9nRXJyb3JDYWxsYmFjayxcbiAgICBkaXNwYXRjaEN1c3RvbUV2ZW50LFxufVxuXG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lKGlkKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgQ2xpcGJvYXJkOiB7XG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcbiAgICAgICAgfSxcbiAgICAgICAgQXBwbGljYXRpb246IHtcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uXG4gICAgICAgIH0sXG4gICAgICAgIExvZyxcbiAgICAgICAgU2NyZWVucyxcbiAgICAgICAgV01MOiB7XG4gICAgICAgICAgICBSZWZyZXNoOiByZWZyZXNoLFxuICAgICAgICB9LFxuICAgICAgICBEaWFsb2c6IHtcbiAgICAgICAgICAgIEluZm8sXG4gICAgICAgICAgICBXYXJuaW5nLFxuICAgICAgICAgICAgRXJyb3IsXG4gICAgICAgICAgICBRdWVzdGlvbixcbiAgICAgICAgICAgIE9wZW5GaWxlLFxuICAgICAgICAgICAgU2F2ZUZpbGUsXG4gICAgICAgIH0sXG4gICAgICAgIEV2ZW50czoge1xuICAgICAgICAgICAgRW1pdCxcbiAgICAgICAgICAgIE9uLFxuICAgICAgICAgICAgT25jZSxcbiAgICAgICAgICAgIE9uTXVsdGlwbGUsXG4gICAgICAgICAgICBPZmYsXG4gICAgICAgICAgICBPZmZBbGwsXG4gICAgICAgIH0sXG4gICAgICAgIFdpbmRvdzogbmV3V2luZG93KGlkKSxcbiAgICB9XG59XG5cbmlmIChERUJVRykge1xuICAgIGNvbnNvbGUubG9nKFwiV2FpbHMgdjMuMC4wIERlYnVnIE1vZGUgRW5hYmxlZFwiKTtcbn1cblxuZW5hYmxlQ29udGV4dE1lbnVzKHRydWUpO1xuXG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKFwiRE9NQ29udGVudExvYWRlZFwiLCBmdW5jdGlvbihldmVudCkge1xuICAgIHJlZnJlc2goKTtcbn0pOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDWUEsTUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBRTVDLFdBQVMsWUFBWSxRQUFRLE1BQU07QUFDL0IsUUFBSSxNQUFNLElBQUksSUFBSSxVQUFVO0FBQzVCLFFBQUksYUFBYSxPQUFPLFVBQVUsTUFBTTtBQUN4QyxRQUFHLE1BQU07QUFDTCxVQUFJLGFBQWEsT0FBTyxRQUFRLEtBQUssVUFBVSxJQUFJLENBQUM7QUFBQSxJQUN4RDtBQUNBLFdBQU8sSUFBSSxR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O0FBRzdDLGNBQVEsaUJBQWlCLFNBQVMsUUFBUTtBQUFBLElBQzlDO0FBQUEsRUFDSjtBQUVPLFdBQVMsVUFBVTtBQUN0Qix5QkFBcUI7QUFDckIsMEJBQXNCO0FBQUEsRUFDMUI7OztBQ25EQSxTQUFPLFFBQVE7QUFBQSxJQUNYLEdBQUcsV0FBVyxFQUFFO0FBQUEsRUFDcEI7QUFHQSxTQUFPLFNBQVM7QUFBQSxJQUNaO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxFQUNKO0FBR08sV0FBUyxXQUFXLElBQUk7QUFDM0IsV0FBTztBQUFBLE1BQ0gsV0FBVztBQUFBLFFBQ1AsR0FBRztBQUFBLE1BQ1A7QUFBQSxNQUNBLGFBQWE7QUFBQSxRQUNULEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFNBQVM7QUFBQSxNQUNiO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxFQUFFO0FBQUEsSUFDeEI7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEscUJBQW1CLElBQUk7QUFFdkIsV0FBUyxpQkFBaUIsb0JBQW9CLFNBQVMsT0FBTztBQUMxRCxZQUFRO0FBQUEsRUFDWixDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiZXZlbnROYW1lIiwgImNhbGwiLCAiY2FsbCIsICJFcnJvciIsICJjYWxsIiwgIkVycm9yIl0KfQo=
