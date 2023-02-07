(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

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
          if (response.headers.get("content-type") && response.headers.get("content-type").indexOf("application/json") !== -1) {
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
  var call = newRuntimeCaller("dialog");
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
      call(type, options).catch((error) => {
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

  // desktop/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    SetText: () => SetText,
    Text: () => Text
  });
  var call2 = newRuntimeCaller("clipboard");
  function SetText(text) {
    return call2("SetText", { text });
  }
  function Text() {
    return call2("Text");
  }

  // desktop/application.js
  var application_exports = {};
  __export(application_exports, {
    Hide: () => Hide,
    Quit: () => Quit,
    Show: () => Show
  });
  var call3 = newRuntimeCaller("application");
  function Hide() {
    return call3("Hide");
  }
  function Show() {
    return call3("Show");
  }
  function Quit() {
    return call3("Quit");
  }

  // desktop/log.js
  var log_exports = {};
  __export(log_exports, {
    Log: () => Log
  });
  var call4 = newRuntimeCaller("log");
  function Log(message) {
    return call4("Log", message);
  }

  // desktop/window.js
  function newWindow(id) {
    let call6 = newRuntimeCaller("window", id);
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
      Center: () => call6("Center"),
      SetTitle: (title) => call6("SetTitle", { title }),
      Fullscreen: () => call6("Fullscreen"),
      UnFullscreen: () => call6("UnFullscreen"),
      SetSize: (width, height) => call6("SetSize", { width, height }),
      Size: () => {
        return call6("Size");
      },
      SetMaxSize: (width, height) => call6("SetMaxSize", { width, height }),
      SetMinSize: (width, height) => call6("SetMinSize", { width, height }),
      SetAlwaysOnTop: (b) => call6("SetAlwaysOnTop", { alwaysOnTop: b }),
      SetPosition: (x, y) => call6("SetPosition", { x, y }),
      Position: () => {
        return call6("Position");
      },
      Screen: () => {
        return call6("Screen");
      },
      Hide: () => call6("Hide"),
      Maximise: () => call6("Maximise"),
      Show: () => call6("Show"),
      ToggleMaximise: () => call6("ToggleMaximise"),
      UnMaximise: () => call6("UnMaximise"),
      Minimise: () => call6("Minimise"),
      UnMinimise: () => call6("UnMinimise"),
      SetBackgroundColour: (r, g, b, a) => call6("SetBackgroundColour", { r, g, b, a })
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
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9ydW50aW1lLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2RpYWxvZ3MuanMiLCAiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9hcHBsaWNhdGlvbi5qcyIsICJkZXNrdG9wL2xvZy5qcyIsICJkZXNrdG9wL3dpbmRvdy5qcyIsICJkZXNrdG9wL2V2ZW50cy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIGlmKGFyZ3MpIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcbiAgICB9XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZmV0Y2godXJsKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiY29udGVudC10eXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiY29udGVudC10eXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgaWQpIHtcbiAgICBpZiAoIWlkIHx8IGlkID09PSAtMSkge1xuICAgICAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncykge1xuICAgICAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICAgICAgfTtcbiAgICB9XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgYXJncyA9IGFyZ3MgfHwge307XG4gICAgICAgIGFyZ3NbXCJ3aW5kb3dJRFwiXSA9IGlkO1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGwob2JqZWN0ICsgXCIuXCIgKyBtZXRob2QsIGFyZ3MpO1xuICAgIH1cbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSdcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZGlhbG9nXCIpO1xuXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zKSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pXG4gICAgfSk7XG59XG5cblxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJJbmZvXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIldhcm5pbmdcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkVycm9yXCIsIG9wdGlvbnMpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJRdWVzdGlvblwiLCBvcHRpb25zKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiT3BlbkZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlNhdmVGaWxlXCIsIG9wdGlvbnMpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChcIlRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJhcHBsaWNhdGlvblwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJIaWRlXCIpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChcIlNob3dcIik7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHJldHVybiBjYWxsKFwiUXVpdFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImxvZ1wiKTtcblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZS5cbiAqIEBwYXJhbSB7bGV2ZWx9IExvZyBsZXZlbFxuICogQHBhcmFtIHttZXNzYWdlfSBNZXNzYWdlIHRvIGxvZ1xuICogQHBhcmFtIHtkYXRhfSBPcHRpb25hbCBkYXRhIHRvIGxvZ1xuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nKG1lc3NhZ2UpIHtcbiAgICByZXR1cm4gY2FsbChcIkxvZ1wiLCBtZXNzYWdlKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1dpbmRvdyhpZCkge1xuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCBpZCk7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBSZWxvYWRBcHA6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFNldFN5c3RlbURlZmF1bHRUaGVtZTogKCkgPT4gY2FsbCgnV0FTRFQnKSxcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxuICAgICAgICAvLyBTZXREYXJrVGhlbWU6ICgpID0+IGNhbGwoJ1dBRFQnKSxcbiAgICAgICAgLy8gSXNGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxuICAgICAgICAvLyBJc01pbmltaXplZDogKCkgPT4gY2FsbCgnV0lNTicpLFxuICAgICAgICAvLyBJc1dpbmRvd2VkOiAoKSA9PiBjYWxsKCdXSUYnKSxcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKCdDZW50ZXInKSxcbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbCgnU2V0VGl0bGUnLCB7dGl0bGV9KSxcbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnRnVsbHNjcmVlbicpLFxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2l6ZTogKCkgPT4geyByZXR1cm4gY2FsbCgnU2l6ZScpIH0sXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRNYXhTaXplJywge3dpZHRoLGhlaWdodH0pLFxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWluU2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0QWx3YXlzT25Ub3A6IChiKSA9PiBjYWxsKCdTZXRBbHdheXNPblRvcCcsIHthbHdheXNPblRvcDpifSksXG4gICAgICAgIFNldFBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbCgnU2V0UG9zaXRpb24nLCB7eCx5fSksXG4gICAgICAgIFBvc2l0aW9uOiAoKSA9PiB7IHJldHVybiBjYWxsKCdQb3NpdGlvbicpIH0sXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJykgfSxcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbCgnSGlkZScpLFxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gY2FsbCgnTWF4aW1pc2UnKSxcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbCgnU2hvdycpLFxuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gY2FsbCgnVG9nZ2xlTWF4aW1pc2UnKSxcbiAgICAgICAgVW5NYXhpbWlzZTogKCkgPT4gY2FsbCgnVW5NYXhpbWlzZScpLFxuICAgICAgICBNaW5pbWlzZTogKCkgPT4gY2FsbCgnTWluaW1pc2UnKSxcbiAgICAgICAgVW5NaW5pbWlzZTogKCkgPT4gY2FsbCgnVW5NaW5pbWlzZScpLFxuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gY2FsbCgnU2V0QmFja2dyb3VuZENvbG91cicsIHtyLCBnLCBiLCBhfSksXG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJldmVudHNcIik7XG5cbi8qKlxuICogVGhlIExpc3RlbmVyIGNsYXNzIGRlZmluZXMgYSBsaXN0ZW5lciEgOi0pXG4gKlxuICogQGNsYXNzIExpc3RlbmVyXG4gKi9cbmNsYXNzIExpc3RlbmVyIHtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIExpc3RlbmVyLlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxuICAgICAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5cbi8qKlxuICogQ3VzdG9tRXZlbnQgZGVmaW5lcyBhIGN1c3RvbSBldmVudC4gSXQgaXMgcGFzc2VkIHRvIGV2ZW50IGxpc3RlbmVycy5cbiAqXG4gKiBAY2xhc3MgQ3VzdG9tRXZlbnRcbiAqL1xuZXhwb3J0IGNsYXNzIEN1c3RvbUV2ZW50IHtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIEN1c3RvbUV2ZW50LlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gTmFtZSBvZiB0aGUgZXZlbnRcbiAgICAgKiBAcGFyYW0ge2FueX0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxuICAgICAqIEBtZW1iZXJvZiBDdXN0b21FdmVudFxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG5hbWUsIGRhdGEpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcbiAgICB9XG59XG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXAoKTtcblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgYG1heENhbGxiYWNrc2AgdGltZXMgYmVmb3JlIGJlaW5nIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb24oQ3VzdG9tRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb24oQ3VzdG9tRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKEN1c3RvbUV2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuLyoqXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uXG4gKlxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcbiAqL1xuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBkaXNwYXRjaGVzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge0N1c3RvbUV2ZW50fSBldmVudFxuICovXG5leHBvcnQgZnVuY3Rpb24gZGlzcGF0Y2hDdXN0b21FdmVudChldmVudCkge1xuICAgIGNvbnNvbGUubG9nKFwiZGlzcGF0Y2hpbmcgZXZlbnQ6IFwiLCB7ZXZlbnR9KTtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmIChsaXN0ZW5lcnMpIHtcbiAgICAgICAgLy8gaXRlcmF0ZSBsaXN0ZW5lcnMgYW5kIGNhbGwgY2FsbGJhY2suIElmIGNhbGxiYWNrIHJldHVybnMgdHJ1ZSwgcmVtb3ZlIGxpc3RlbmVyXG4gICAgICAgIGxldCB0b1JlbW92ZSA9IFtdO1xuICAgICAgICBsaXN0ZW5lcnMuZm9yRWFjaChsaXN0ZW5lciA9PiB7XG4gICAgICAgICAgICBsZXQgcmVtb3ZlID0gbGlzdGVuZXIuQ2FsbGJhY2soZXZlbnQpXG4gICAgICAgICAgICBpZiAocmVtb3ZlKSB7XG4gICAgICAgICAgICAgICAgdG9SZW1vdmUucHVzaChsaXN0ZW5lcik7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgICAgICAvLyByZW1vdmUgbGlzdGVuZXJzXG4gICAgICAgIGlmICh0b1JlbW92ZS5sZW5ndGggPiAwKSB7XG4gICAgICAgICAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gIXRvUmVtb3ZlLmluY2x1ZGVzKGwpKTtcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPbixcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJzIGNhbiBiZSB1bnJlZ2lzdGVyZWQgdmlhIGBhZGRpdGlvbmFsRXZlbnROYW1lc2BcbiAqXG4gW3YzIENIQU5HRV0gT2ZmIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIGxldCBldmVudHNUb1JlbW92ZSA9IFtldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzXTtcbiAgICBldmVudHNUb1JlbW92ZS5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xuICAgIH0pXG59XG5cbi8qKlxuICogT2ZmQWxsIHVucmVnaXN0ZXJzIGFsbCBsaXN0ZW5lcnNcbiAqIFt2MyBDSEFOR0VdIE9mZkFsbCBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkge1xuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XG59XG5cbi8qXG4gICBFbWl0IGVtaXRzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVtaXQoZXZlbnQpIHtcbiAgICByZXR1cm4gY2FsbChcIkVtaXRcIiwgZXZlbnQpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtkaWFsb2dDYWxsYmFjaywgZGlhbG9nRXJyb3JDYWxsYmFjaywgRXJyb3IsIEluZm8sIE9wZW5GaWxlLCBRdWVzdGlvbiwgU2F2ZUZpbGUsIFdhcm5pbmcsfSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5cbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tICcuL2NsaXBib2FyZCc7XG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tICcuL2FwcGxpY2F0aW9uJztcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XG5cbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCB7ZGlzcGF0Y2hDdXN0b21FdmVudCwgRW1pdCwgT2ZmLCBPZmZBbGwsIE9uLCBPbmNlLCBPbk11bHRpcGxlfSBmcm9tIFwiLi9ldmVudHNcIjtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cud2FpbHMgPSB7XG4gICAgLi4ubmV3UnVudGltZSgtMSksXG59O1xuXG53aW5kb3cuX3dhaWxzID0ge1xuICAgIGRpYWxvZ0NhbGxiYWNrLFxuICAgIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssXG4gICAgZGlzcGF0Y2hDdXN0b21FdmVudCxcbn1cblxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZShpZCkge1xuICAgIHJldHVybiB7XG4gICAgICAgIENsaXBib2FyZDoge1xuICAgICAgICAgICAgLi4uQ2xpcGJvYXJkXG4gICAgICAgIH0sXG4gICAgICAgIEFwcGxpY2F0aW9uOiB7XG4gICAgICAgICAgICAuLi5BcHBsaWNhdGlvblxuICAgICAgICB9LFxuICAgICAgICBMb2csXG4gICAgICAgIERpYWxvZzoge1xuICAgICAgICAgICAgSW5mbyxcbiAgICAgICAgICAgIFdhcm5pbmcsXG4gICAgICAgICAgICBFcnJvcixcbiAgICAgICAgICAgIFF1ZXN0aW9uLFxuICAgICAgICAgICAgT3BlbkZpbGUsXG4gICAgICAgICAgICBTYXZlRmlsZSxcbiAgICAgICAgfSxcbiAgICAgICAgRXZlbnRzOiB7XG4gICAgICAgICAgICBFbWl0LFxuICAgICAgICAgICAgT24sXG4gICAgICAgICAgICBPbmNlLFxuICAgICAgICAgICAgT25NdWx0aXBsZSxcbiAgICAgICAgICAgIE9mZixcbiAgICAgICAgICAgIE9mZkFsbCxcbiAgICAgICAgfSxcbiAgICAgICAgV2luZG93OiBuZXdXaW5kb3coaWQpLFxuICAgIH1cbn1cblxuaWYgKERFQlVHKSB7XG4gICAgY29uc29sZS5sb2coXCJXYWlscyB2My4wLjAgRGVidWcgTW9kZSBFbmFibGVkXCIpO1xufVxuXG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQVlBLE1BQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUU1QyxXQUFTLFlBQVksUUFBUSxNQUFNO0FBQy9CLFFBQUksTUFBTSxJQUFJLElBQUksVUFBVTtBQUM1QixRQUFJLGFBQWEsT0FBTyxVQUFVLE1BQU07QUFDeEMsUUFBRyxNQUFNO0FBQ0wsVUFBSSxhQUFhLE9BQU8sUUFBUSxLQUFLLFVBQVUsSUFBSSxDQUFDO0FBQUEsSUFDeEQ7QUFDQSxXQUFPLElBQUksUUFBUSxDQUFDLFNBQVMsV0FBVztBQUNwQyxZQUFNLEdBQUcsRUFDSixLQUFLLGNBQVk7QUFDZCxZQUFJLFNBQVMsSUFBSTtBQUViLGNBQUksU0FBUyxRQUFRLElBQUksY0FBYyxLQUFLLFNBQVMsUUFBUSxJQUFJLGNBQWMsRUFBRSxRQUFRLGtCQUFrQixNQUFNLElBQUk7QUFDakgsbUJBQU8sU0FBUyxLQUFLO0FBQUEsVUFDekIsT0FBTztBQUNILG1CQUFPLFNBQVMsS0FBSztBQUFBLFVBQ3pCO0FBQUEsUUFDSjtBQUNBLGVBQU8sTUFBTSxTQUFTLFVBQVUsQ0FBQztBQUFBLE1BQ3JDLENBQUMsRUFDQSxLQUFLLFVBQVEsUUFBUSxJQUFJLENBQUMsRUFDMUIsTUFBTSxXQUFTLE9BQU8sS0FBSyxDQUFDO0FBQUEsSUFDckMsQ0FBQztBQUFBLEVBQ0w7QUFFTyxXQUFTLGlCQUFpQixRQUFRLElBQUk7QUFDekMsUUFBSSxDQUFDLE1BQU0sT0FBTyxJQUFJO0FBQ2xCLGFBQU8sU0FBVSxRQUFRLE1BQU07QUFDM0IsZUFBTyxZQUFZLFNBQVMsTUFBTSxRQUFRLElBQUk7QUFBQSxNQUNsRDtBQUFBLElBQ0o7QUFDQSxXQUFPLFNBQVUsUUFBUSxNQUFNO0FBQzNCLGFBQU8sUUFBUSxDQUFDO0FBQ2hCLFdBQUssVUFBVSxJQUFJO0FBQ25CLGFBQU8sWUFBWSxTQUFTLE1BQU0sUUFBUSxJQUFJO0FBQUEsSUFDbEQ7QUFBQSxFQUNKOzs7QUNqREEsTUFBSSxjQUNGO0FBV0ssTUFBSSxTQUFTLENBQUMsT0FBTyxPQUFPO0FBQ2pDLFFBQUksS0FBSztBQUNULFFBQUksSUFBSTtBQUNSLFdBQU8sS0FBSztBQUNWLFlBQU0sWUFBYSxLQUFLLE9BQU8sSUFBSSxLQUFNLENBQUM7QUFBQSxJQUM1QztBQUNBLFdBQU87QUFBQSxFQUNUOzs7QUNIQSxNQUFJLE9BQU8saUJBQWlCLFFBQVE7QUFFcEMsTUFBSSxrQkFBa0Isb0JBQUksSUFBSTtBQUU5QixXQUFTLGFBQWE7QUFDbEIsUUFBSTtBQUNKLE9BQUc7QUFDQyxlQUFTLE9BQU87QUFBQSxJQUNwQixTQUFTLGdCQUFnQixJQUFJLE1BQU07QUFDbkMsV0FBTztBQUFBLEVBQ1g7QUFFTyxXQUFTLGVBQWUsSUFBSSxNQUFNLFFBQVE7QUFDN0MsUUFBSSxJQUFJLGdCQUFnQixJQUFJLEVBQUU7QUFDOUIsUUFBSSxHQUFHO0FBQ0gsVUFBSSxRQUFRO0FBQ1IsVUFBRSxRQUFRLEtBQUssTUFBTSxJQUFJLENBQUM7QUFBQSxNQUM5QixPQUFPO0FBQ0gsVUFBRSxRQUFRLElBQUk7QUFBQSxNQUNsQjtBQUNBLHNCQUFnQixPQUFPLEVBQUU7QUFBQSxJQUM3QjtBQUFBLEVBQ0o7QUFDTyxXQUFTLG9CQUFvQixJQUFJLFNBQVM7QUFDN0MsUUFBSSxJQUFJLGdCQUFnQixJQUFJLEVBQUU7QUFDOUIsUUFBSSxHQUFHO0FBQ0gsUUFBRSxPQUFPLE9BQU87QUFDaEIsc0JBQWdCLE9BQU8sRUFBRTtBQUFBLElBQzdCO0FBQUEsRUFDSjtBQUVBLFdBQVMsT0FBTyxNQUFNLFNBQVM7QUFDM0IsV0FBTyxJQUFJLFFBQVEsQ0FBQyxTQUFTLFdBQVc7QUFDcEMsVUFBSSxLQUFLLFdBQVc7QUFDcEIsZ0JBQVUsV0FBVyxDQUFDO0FBQ3RCLGNBQVEsV0FBVyxJQUFJO0FBQ3ZCLHNCQUFnQixJQUFJLElBQUksRUFBQyxTQUFTLE9BQU0sQ0FBQztBQUN6QyxXQUFLLE1BQU0sT0FBTyxFQUFFLE1BQU0sQ0FBQyxVQUFVO0FBQ2pDLGVBQU8sS0FBSztBQUNaLHdCQUFnQixPQUFPLEVBQUU7QUFBQSxNQUM3QixDQUFDO0FBQUEsSUFDTCxDQUFDO0FBQUEsRUFDTDtBQUdPLFdBQVMsS0FBSyxTQUFTO0FBQzFCLFdBQU8sT0FBTyxRQUFRLE9BQU87QUFBQSxFQUNqQztBQUVPLFdBQVMsUUFBUSxTQUFTO0FBQzdCLFdBQU8sT0FBTyxXQUFXLE9BQU87QUFBQSxFQUNwQztBQUVPLFdBQVNBLE9BQU0sU0FBUztBQUMzQixXQUFPLE9BQU8sU0FBUyxPQUFPO0FBQUEsRUFDbEM7QUFFTyxXQUFTLFNBQVMsU0FBUztBQUM5QixXQUFPLE9BQU8sWUFBWSxPQUFPO0FBQUEsRUFDckM7QUFFTyxXQUFTLFNBQVMsU0FBUztBQUM5QixXQUFPLE9BQU8sWUFBWSxPQUFPO0FBQUEsRUFDckM7QUFFTyxXQUFTLFNBQVMsU0FBUztBQUM5QixXQUFPLE9BQU8sWUFBWSxPQUFPO0FBQUEsRUFDckM7OztBQ25GQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBY0EsTUFBSUMsUUFBTyxpQkFBaUIsV0FBVztBQUVoQyxXQUFTLFFBQVEsTUFBTTtBQUMxQixXQUFPQSxNQUFLLFdBQVcsRUFBQyxLQUFJLENBQUM7QUFBQSxFQUNqQztBQUVPLFdBQVMsT0FBTztBQUNuQixXQUFPQSxNQUFLLE1BQU07QUFBQSxFQUN0Qjs7O0FDdEJBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWNBLE1BQUlDLFFBQU8saUJBQWlCLGFBQWE7QUFFbEMsV0FBUyxPQUFPO0FBQ25CLFdBQU9BLE1BQUssTUFBTTtBQUFBLEVBQ3RCO0FBRU8sV0FBUyxPQUFPO0FBQ25CLFdBQU9BLE1BQUssTUFBTTtBQUFBLEVBQ3RCO0FBRU8sV0FBUyxPQUFPO0FBQ25CLFdBQU9BLE1BQUssTUFBTTtBQUFBLEVBQ3RCOzs7QUMxQkE7QUFBQTtBQUFBO0FBQUE7QUFjQSxNQUFJQyxRQUFPLGlCQUFpQixLQUFLO0FBUTFCLFdBQVMsSUFBSSxTQUFTO0FBQ3pCLFdBQU9BLE1BQUssT0FBTyxPQUFPO0FBQUEsRUFDOUI7OztBQ1ZPLFdBQVMsVUFBVSxJQUFJO0FBQzFCLFFBQUlDLFFBQU8saUJBQWlCLFVBQVUsRUFBRTtBQUN4QyxXQUFPO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsTUFVSCxRQUFRLE1BQU1BLE1BQUssUUFBUTtBQUFBLE1BQzNCLFVBQVUsQ0FBQyxVQUFVQSxNQUFLLFlBQVksRUFBQyxNQUFLLENBQUM7QUFBQSxNQUM3QyxZQUFZLE1BQU1BLE1BQUssWUFBWTtBQUFBLE1BQ25DLGNBQWMsTUFBTUEsTUFBSyxjQUFjO0FBQUEsTUFDdkMsU0FBUyxDQUFDLE9BQU8sV0FBV0EsTUFBSyxXQUFXLEVBQUMsT0FBTSxPQUFNLENBQUM7QUFBQSxNQUMxRCxNQUFNLE1BQU07QUFBRSxlQUFPQSxNQUFLLE1BQU07QUFBQSxNQUFFO0FBQUEsTUFDbEMsWUFBWSxDQUFDLE9BQU8sV0FBV0EsTUFBSyxjQUFjLEVBQUMsT0FBTSxPQUFNLENBQUM7QUFBQSxNQUNoRSxZQUFZLENBQUMsT0FBTyxXQUFXQSxNQUFLLGNBQWMsRUFBQyxPQUFNLE9BQU0sQ0FBQztBQUFBLE1BQ2hFLGdCQUFnQixDQUFDLE1BQU1BLE1BQUssa0JBQWtCLEVBQUMsYUFBWSxFQUFDLENBQUM7QUFBQSxNQUM3RCxhQUFhLENBQUMsR0FBRyxNQUFNQSxNQUFLLGVBQWUsRUFBQyxHQUFFLEVBQUMsQ0FBQztBQUFBLE1BQ2hELFVBQVUsTUFBTTtBQUFFLGVBQU9BLE1BQUssVUFBVTtBQUFBLE1BQUU7QUFBQSxNQUMxQyxRQUFRLE1BQU07QUFBRSxlQUFPQSxNQUFLLFFBQVE7QUFBQSxNQUFFO0FBQUEsTUFDdEMsTUFBTSxNQUFNQSxNQUFLLE1BQU07QUFBQSxNQUN2QixVQUFVLE1BQU1BLE1BQUssVUFBVTtBQUFBLE1BQy9CLE1BQU0sTUFBTUEsTUFBSyxNQUFNO0FBQUEsTUFDdkIsZ0JBQWdCLE1BQU1BLE1BQUssZ0JBQWdCO0FBQUEsTUFDM0MsWUFBWSxNQUFNQSxNQUFLLFlBQVk7QUFBQSxNQUNuQyxVQUFVLE1BQU1BLE1BQUssVUFBVTtBQUFBLE1BQy9CLFlBQVksTUFBTUEsTUFBSyxZQUFZO0FBQUEsTUFDbkMscUJBQXFCLENBQUMsR0FBRyxHQUFHLEdBQUcsTUFBTUEsTUFBSyx1QkFBdUIsRUFBQyxHQUFHLEdBQUcsR0FBRyxFQUFDLENBQUM7QUFBQSxJQUNqRjtBQUFBLEVBQ0o7OztBQ2pDQSxNQUFJQyxRQUFPLGlCQUFpQixRQUFRO0FBT3BDLE1BQU0sV0FBTixNQUFlO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxJQVFYLFlBQVksV0FBVyxVQUFVLGNBQWM7QUFDM0MsV0FBSyxZQUFZO0FBRWpCLFdBQUssZUFBZSxnQkFBZ0I7QUFHcEMsV0FBSyxXQUFXLENBQUMsU0FBUztBQUN0QixpQkFBUyxJQUFJO0FBRWIsWUFBSSxLQUFLLGlCQUFpQixJQUFJO0FBQzFCLGlCQUFPO0FBQUEsUUFDWDtBQUVBLGFBQUssZ0JBQWdCO0FBQ3JCLGVBQU8sS0FBSyxpQkFBaUI7QUFBQSxNQUNqQztBQUFBLElBQ0o7QUFBQSxFQUNKO0FBcUJPLE1BQU0saUJBQWlCLG9CQUFJLElBQUk7QUFXL0IsV0FBUyxXQUFXLFdBQVcsVUFBVSxjQUFjO0FBQzFELFFBQUksWUFBWSxlQUFlLElBQUksU0FBUyxLQUFLLENBQUM7QUFDbEQsVUFBTSxlQUFlLElBQUksU0FBUyxXQUFXLFVBQVUsWUFBWTtBQUNuRSxjQUFVLEtBQUssWUFBWTtBQUMzQixtQkFBZSxJQUFJLFdBQVcsU0FBUztBQUN2QyxXQUFPLE1BQU0sWUFBWSxZQUFZO0FBQUEsRUFDekM7QUFVTyxXQUFTLEdBQUcsV0FBVyxVQUFVO0FBQ3BDLFdBQU8sV0FBVyxXQUFXLFVBQVUsRUFBRTtBQUFBLEVBQzdDO0FBVU8sV0FBUyxLQUFLLFdBQVcsVUFBVTtBQUN0QyxXQUFPLFdBQVcsV0FBVyxVQUFVLENBQUM7QUFBQSxFQUM1QztBQU9BLFdBQVMsWUFBWSxVQUFVO0FBQzNCLFVBQU0sWUFBWSxTQUFTO0FBRTNCLFFBQUksWUFBWSxlQUFlLElBQUksU0FBUyxFQUFFLE9BQU8sT0FBSyxNQUFNLFFBQVE7QUFDeEUsUUFBSSxVQUFVLFdBQVcsR0FBRztBQUN4QixxQkFBZSxPQUFPLFNBQVM7QUFBQSxJQUNuQyxPQUFPO0FBQ0gscUJBQWUsSUFBSSxXQUFXLFNBQVM7QUFBQSxJQUMzQztBQUFBLEVBQ0o7QUFRTyxXQUFTLG9CQUFvQixPQUFPO0FBQ3ZDLFlBQVEsSUFBSSx1QkFBdUIsRUFBQyxNQUFLLENBQUM7QUFDMUMsUUFBSSxZQUFZLGVBQWUsSUFBSSxNQUFNLElBQUk7QUFDN0MsUUFBSSxXQUFXO0FBRVgsVUFBSSxXQUFXLENBQUM7QUFDaEIsZ0JBQVUsUUFBUSxjQUFZO0FBQzFCLFlBQUksU0FBUyxTQUFTLFNBQVMsS0FBSztBQUNwQyxZQUFJLFFBQVE7QUFDUixtQkFBUyxLQUFLLFFBQVE7QUFBQSxRQUMxQjtBQUFBLE1BQ0osQ0FBQztBQUVELFVBQUksU0FBUyxTQUFTLEdBQUc7QUFDckIsb0JBQVksVUFBVSxPQUFPLE9BQUssQ0FBQyxTQUFTLFNBQVMsQ0FBQyxDQUFDO0FBQ3ZELFlBQUksVUFBVSxXQUFXLEdBQUc7QUFDeEIseUJBQWUsT0FBTyxNQUFNLElBQUk7QUFBQSxRQUNwQyxPQUFPO0FBQ0gseUJBQWUsSUFBSSxNQUFNLE1BQU0sU0FBUztBQUFBLFFBQzVDO0FBQUEsTUFDSjtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBV08sV0FBUyxJQUFJLGNBQWMsc0JBQXNCO0FBQ3BELFFBQUksaUJBQWlCLENBQUMsV0FBVyxHQUFHLG9CQUFvQjtBQUN4RCxtQkFBZSxRQUFRLENBQUFDLGVBQWE7QUFDaEMscUJBQWUsT0FBT0EsVUFBUztBQUFBLElBQ25DLENBQUM7QUFBQSxFQUNMO0FBT08sV0FBUyxTQUFTO0FBQ3JCLG1CQUFlLE1BQU07QUFBQSxFQUN6QjtBQUtPLFdBQVMsS0FBSyxPQUFPO0FBQ3hCLFdBQU9DLE1BQUssUUFBUSxLQUFLO0FBQUEsRUFDN0I7OztBQ3JLQSxTQUFPLFFBQVE7QUFBQSxJQUNYLEdBQUcsV0FBVyxFQUFFO0FBQUEsRUFDcEI7QUFFQSxTQUFPLFNBQVM7QUFBQSxJQUNaO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxFQUNKO0FBR08sV0FBUyxXQUFXLElBQUk7QUFDM0IsV0FBTztBQUFBLE1BQ0gsV0FBVztBQUFBLFFBQ1AsR0FBRztBQUFBLE1BQ1A7QUFBQSxNQUNBLGFBQWE7QUFBQSxRQUNULEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQTtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQSxPQUFBQztBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRLFVBQVUsRUFBRTtBQUFBLElBQ3hCO0FBQUEsRUFDSjtBQUVBLE1BQUksTUFBTztBQUNQLFlBQVEsSUFBSSxpQ0FBaUM7QUFBQSxFQUNqRDsiLAogICJuYW1lcyI6IFsiRXJyb3IiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgIkVycm9yIl0KfQo=
