(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

  // desktop/log.js
  var log_exports = {};
  __export(log_exports, {
    LogDebug: () => LogDebug,
    LogError: () => LogError,
    LogFatal: () => LogFatal,
    LogInfo: () => LogInfo,
    LogLevel: () => LogLevel,
    LogPrint: () => LogPrint,
    LogTrace: () => LogTrace,
    LogWarning: () => LogWarning,
    SetLogLevel: () => SetLogLevel
  });
  function sendLogMessage(level, message) {
    window.WailsInvoke("L" + level + message);
  }
  function LogTrace(message) {
    sendLogMessage("T", message);
  }
  function LogPrint(message) {
    sendLogMessage("P", message);
  }
  function LogDebug(message) {
    sendLogMessage("D", message);
  }
  function LogInfo(message) {
    sendLogMessage("I", message);
  }
  function LogWarning(message) {
    sendLogMessage("W", message);
  }
  function LogError(message) {
    sendLogMessage("E", message);
  }
  function LogFatal(message) {
    sendLogMessage("F", message);
  }
  function SetLogLevel(loglevel) {
    sendLogMessage("S", loglevel);
  }
  var LogLevel = {
    TRACE: 1,
    DEBUG: 2,
    INFO: 3,
    WARNING: 4,
    ERROR: 5
  };

  // desktop/events.js
  var Listener = class {
    constructor(callback, maxCallbacks) {
      maxCallbacks = maxCallbacks || -1;
      this.Callback = (data) => {
        callback.apply(null, data);
        if (maxCallbacks === -1) {
          return false;
        }
        maxCallbacks -= 1;
        return maxCallbacks === 0;
      };
    }
  };
  var eventListeners = {};
  function EventsOnMultiple(eventName, callback, maxCallbacks) {
    eventListeners[eventName] = eventListeners[eventName] || [];
    const thisListener = new Listener(callback, maxCallbacks);
    eventListeners[eventName].push(thisListener);
  }
  function EventsOn(eventName, callback) {
    EventsOnMultiple(eventName, callback, -1);
  }
  function EventsOnce(eventName, callback) {
    EventsOnMultiple(eventName, callback, 1);
  }
  function notifyListeners(eventData) {
    let eventName = eventData.name;
    if (eventListeners[eventName]) {
      const newEventListenerList = eventListeners[eventName].slice();
      for (let count = 0; count < eventListeners[eventName].length; count += 1) {
        const listener = eventListeners[eventName][count];
        let data = eventData.data;
        const destroy = listener.Callback(data);
        if (destroy) {
          newEventListenerList.splice(count, 1);
        }
      }
      eventListeners[eventName] = newEventListenerList;
    }
  }
  function EventsNotify(notifyMessage) {
    let message;
    try {
      message = JSON.parse(notifyMessage);
    } catch (e) {
      const error = "Invalid JSON passed to Notify: " + notifyMessage;
      throw new Error(error);
    }
    notifyListeners(message);
  }
  function EventsEmit(eventName) {
    const payload = {
      name: eventName,
      data: [].slice.apply(arguments).slice(1)
    };
    notifyListeners(payload);
    window.WailsInvoke("EE" + JSON.stringify(payload));
  }
  function removeListener(eventName) {
    delete eventListeners[eventName];
    window.WailsInvoke("EX" + eventName);
  }
  function EventsOff(eventName, ...additionalEventNames) {
    removeListener(eventName);
    if (additionalEventNames.length > 0) {
      additionalEventNames.forEach((eventName2) => {
        removeListener(eventName2);
      });
    }
  }

  // desktop/calls.js
  var callbacks = {};
  function cryptoRandom() {
    var array = new Uint32Array(1);
    return window.crypto.getRandomValues(array)[0];
  }
  function basicRandom() {
    return Math.random() * 9007199254740991;
  }
  var randomFunc;
  if (window.crypto) {
    randomFunc = cryptoRandom;
  } else {
    randomFunc = basicRandom;
  }
  function Call(name, args, timeout) {
    if (timeout == null) {
      timeout = 0;
    }
    return new Promise(function(resolve, reject) {
      var callbackID;
      do {
        callbackID = name + "-" + randomFunc();
      } while (callbacks[callbackID]);
      var timeoutHandle;
      if (timeout > 0) {
        timeoutHandle = setTimeout(function() {
          reject(Error("Call to " + name + " timed out. Request ID: " + callbackID));
        }, timeout);
      }
      callbacks[callbackID] = {
        timeoutHandle,
        reject,
        resolve
      };
      try {
        const payload = {
          name,
          args,
          callbackID
        };
        window.WailsInvoke("C" + JSON.stringify(payload));
      } catch (e) {
        console.error(e);
      }
    });
  }
  window.ObfuscatedCall = (id, args, timeout) => {
    if (timeout == null) {
      timeout = 0;
    }
    return new Promise(function (resolve, reject) {
      var callbackID;
      do {
        callbackID = id + "-" + randomFunc();
      } while (callbacks[callbackID]);
      var timeoutHandle;
      if (timeout > 0) {
        timeoutHandle = setTimeout(function () {
          reject(Error("Call to method " + id + " timed out. Request ID: " + callbackID));
        }, timeout);
      }
      callbacks[callbackID] = {
        timeoutHandle,
        reject,
        resolve
      };
      try {
        const payload = {
          id,
          args,
          callbackID
        };
        window.WailsInvoke("c" + JSON.stringify(payload));
      } catch (e) {
        console.error(e);
      }
    });
  };
  function Callback(incomingMessage) {
    let message;
    try {
      message = JSON.parse(incomingMessage);
    } catch (e) {
      const error = `Invalid JSON passed to callback: ${e.message}. Message: ${incomingMessage}`;
      runtime.LogDebug(error);
      throw new Error(error);
    }
    let callbackID = message.callbackid;
    let callbackData = callbacks[callbackID];
    if (!callbackData) {
      const error = `Callback '${callbackID}' not registered!!!`;
      console.error(error);
      throw new Error(error);
    }
    clearTimeout(callbackData.timeoutHandle);
    delete callbacks[callbackID];
    if (message.error) {
      callbackData.reject(message.error);
    } else {
      callbackData.resolve(message.result);
    }
  }

  // desktop/bindings.js
  window.go = {};
  function SetBindings(bindingsMap) {
    try {
      bindingsMap = JSON.parse(bindingsMap);
    } catch (e) {
      console.error(e);
    }
    window.go = window.go || {};
    Object.keys(bindingsMap).forEach((packageName) => {
      window.go[packageName] = window.go[packageName] || {};
      Object.keys(bindingsMap[packageName]).forEach((structName) => {
        window.go[packageName][structName] = window.go[packageName][structName] || {};
        Object.keys(bindingsMap[packageName][structName]).forEach((methodName) => {
          window.go[packageName][structName][methodName] = function() {
            let timeout = 0;
            function dynamic() {
              const args = [].slice.call(arguments);
              return Call([packageName, structName, methodName].join("."), args, timeout);
            }
            dynamic.setTimeout = function(newTimeout) {
              timeout = newTimeout;
            };
            dynamic.getTimeout = function() {
              return timeout;
            };
            return dynamic;
          }();
        });
      });
    });
  }

  // desktop/window.js
  var window_exports = {};
  __export(window_exports, {
    WindowCenter: () => WindowCenter,
    WindowFullscreen: () => WindowFullscreen,
    WindowGetPosition: () => WindowGetPosition,
    WindowGetSize: () => WindowGetSize,
    WindowHide: () => WindowHide,
    WindowIsFullscreen: () => WindowIsFullscreen,
    WindowIsMaximised: () => WindowIsMaximised,
    WindowIsMinimised: () => WindowIsMinimised,
    WindowIsNormal: () => WindowIsNormal,
    WindowMaximise: () => WindowMaximise,
    WindowMinimise: () => WindowMinimise,
    WindowReload: () => WindowReload,
    WindowReloadApp: () => WindowReloadApp,
    WindowSetAlwaysOnTop: () => WindowSetAlwaysOnTop,
    WindowSetBackgroundColour: () => WindowSetBackgroundColour,
    WindowSetDarkTheme: () => WindowSetDarkTheme,
    WindowSetLightTheme: () => WindowSetLightTheme,
    WindowSetMaxSize: () => WindowSetMaxSize,
    WindowSetMinSize: () => WindowSetMinSize,
    WindowSetPosition: () => WindowSetPosition,
    WindowSetSize: () => WindowSetSize,
    WindowSetSystemDefaultTheme: () => WindowSetSystemDefaultTheme,
    WindowSetTitle: () => WindowSetTitle,
    WindowShow: () => WindowShow,
    WindowToggleMaximise: () => WindowToggleMaximise,
    WindowUnfullscreen: () => WindowUnfullscreen,
    WindowUnmaximise: () => WindowUnmaximise,
    WindowUnminimise: () => WindowUnminimise
  });
  function WindowReload() {
    window.location.reload();
  }
  function WindowReloadApp() {
    window.WailsInvoke("WR");
  }
  function WindowSetSystemDefaultTheme() {
    window.WailsInvoke("WASDT");
  }
  function WindowSetLightTheme() {
    window.WailsInvoke("WALT");
  }
  function WindowSetDarkTheme() {
    window.WailsInvoke("WADT");
  }
  function WindowCenter() {
    window.WailsInvoke("Wc");
  }
  function WindowSetTitle(title) {
    window.WailsInvoke("WT" + title);
  }
  function WindowFullscreen() {
    window.WailsInvoke("WF");
  }
  function WindowUnfullscreen() {
    window.WailsInvoke("Wf");
  }
  function WindowIsFullscreen() {
    return Call(":wails:WindowIsFullscreen");
  }
  function WindowSetSize(width, height) {
    window.WailsInvoke("Ws:" + width + ":" + height);
  }
  function WindowGetSize() {
    return Call(":wails:WindowGetSize");
  }
  function WindowSetMaxSize(width, height) {
    window.WailsInvoke("WZ:" + width + ":" + height);
  }
  function WindowSetMinSize(width, height) {
    window.WailsInvoke("Wz:" + width + ":" + height);
  }
  function WindowSetAlwaysOnTop(b) {
    window.WailsInvoke("WATP:" + (b ? "1" : "0"));
  }
  function WindowSetPosition(x, y) {
    window.WailsInvoke("Wp:" + x + ":" + y);
  }
  function WindowGetPosition() {
    return Call(":wails:WindowGetPos");
  }
  function WindowHide() {
    window.WailsInvoke("WH");
  }
  function WindowShow() {
    window.WailsInvoke("WS");
  }
  function WindowMaximise() {
    window.WailsInvoke("WM");
  }
  function WindowToggleMaximise() {
    window.WailsInvoke("Wt");
  }
  function WindowUnmaximise() {
    window.WailsInvoke("WU");
  }
  function WindowIsMaximised() {
    return Call(":wails:WindowIsMaximised");
  }
  function WindowMinimise() {
    window.WailsInvoke("Wm");
  }
  function WindowUnminimise() {
    window.WailsInvoke("Wu");
  }
  function WindowIsMinimised() {
    return Call(":wails:WindowIsMinimised");
  }
  function WindowIsNormal() {
    return Call(":wails:WindowIsNormal");
  }
  function WindowSetBackgroundColour(R, G, B, A) {
    let rgba = JSON.stringify({ r: R || 0, g: G || 0, b: B || 0, a: A || 255 });
    window.WailsInvoke("Wr:" + rgba);
  }

  // desktop/screen.js
  var screen_exports = {};
  __export(screen_exports, {
    ScreenGetAll: () => ScreenGetAll
  });
  function ScreenGetAll() {
    return Call(":wails:ScreenGetAll");
  }

  // desktop/browser.js
  var browser_exports = {};
  __export(browser_exports, {
    BrowserOpenURL: () => BrowserOpenURL
  });
  function BrowserOpenURL(url) {
    window.WailsInvoke("BO:" + url);
  }

  // desktop/main.js
  function Quit() {
    window.WailsInvoke("Q");
  }
  function Show() {
    window.WailsInvoke("S");
  }
  function Hide() {
    window.WailsInvoke("H");
  }
  function Environment() {
    return Call(":wails:Environment");
  }
  window.runtime = {
    ...log_exports,
    ...window_exports,
    ...browser_exports,
    ...screen_exports,
    EventsOn,
    EventsOnce,
    EventsOnMultiple,
    EventsEmit,
    EventsOff,
    Environment,
    Show,
    Hide,
    Quit
  };
  window.wails = {
    Callback,
    EventsNotify,
    SetBindings,
    eventListeners,
    callbacks,
    flags: {
      disableScrollbarDrag: false,
      disableWailsDefaultContextMenu: false,
      enableResize: false,
      defaultCursor: null,
      borderThickness: 6,
      shouldDrag: false
    }
  };
  if (window.wailsbindings) {
    window.wails.SetBindings(window.wailsbindings);
    delete window.wails.SetBindings;
  }
  if (false) {
    delete window.wailsbindings;
  }
  window.addEventListener("mouseup", () => {
    window.wails.flags.shouldDrag = false;
  });
  var cssDragTest = function (e) {
    return window.getComputedStyle(e.target).getPropertyValue("--wails-draggable") === "drag";
  };
  var elementDragTest = function (e) {
    let currentElement = e.target;
    while (currentElement != null) {
      if (currentElement.hasAttribute("data-wails-no-drag")) {
        break;
      } else if (currentElement.hasAttribute("data-wails-drag")) {
        return true;
      }
      currentElement = currentElement.parentElement;
    }
    return false;
  };
  var dragTest = elementDragTest;
  window.wails.useCSSDrag = function (t) {
    if (t === false) {
      console.log("Using original drag detection");
      dragTest = elementDragTest;
    } else {
      console.log("Using CSS drag detection");
      dragTest = cssDragTest;
    }
  };
  window.addEventListener("mousedown", (e) => {
    if (window.wails.flags.resizeEdge) {
      window.WailsInvoke("resize:" + window.wails.flags.resizeEdge);
      e.preventDefault();
      return;
    }
    if (dragTest(e)) {
      if (window.wails.flags.disableScrollbarDrag) {
        if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
          return;
        }
      }
      window.wails.flags.shouldDrag = true;
    }
  });
  function setResize(cursor) {
    document.body.style.cursor = cursor || window.wails.flags.defaultCursor;
    window.wails.flags.resizeEdge = cursor;
  }
  window.addEventListener("mousemove", function(e) {
    if (window.wails.flags.shouldDrag) {
      window.WailsInvoke("drag");
      return;
    }
    if (!window.wails.flags.enableResize) {
      return;
    }
    if (window.wails.flags.defaultCursor == null) {
      window.wails.flags.defaultCursor = document.body.style.cursor;
    }
    if (window.outerWidth - e.clientX < window.wails.flags.borderThickness && window.outerHeight - e.clientY < window.wails.flags.borderThickness) {
      document.body.style.cursor = "se-resize";
    }
    let rightBorder = window.outerWidth - e.clientX < window.wails.flags.borderThickness;
    let leftBorder = e.clientX < window.wails.flags.borderThickness;
    let topBorder = e.clientY < window.wails.flags.borderThickness;
    let bottomBorder = window.outerHeight - e.clientY < window.wails.flags.borderThickness;
    if (!leftBorder && !rightBorder && !topBorder && !bottomBorder && window.wails.flags.resizeEdge !== void 0) {
      setResize();
    } else if (rightBorder && bottomBorder)
      setResize("se-resize");
    else if (leftBorder && bottomBorder)
      setResize("sw-resize");
    else if (leftBorder && topBorder)
      setResize("nw-resize");
    else if (topBorder && rightBorder)
      setResize("ne-resize");
    else if (leftBorder)
      setResize("w-resize");
    else if (topBorder)
      setResize("n-resize");
    else if (bottomBorder)
      setResize("s-resize");
    else if (rightBorder)
      setResize("e-resize");
  });
  window.addEventListener("contextmenu", function(e) {
    if (window.wails.flags.disableWailsDefaultContextMenu) {
      e.preventDefault();
    }
  });
  window.WailsInvoke("runtime:ready");
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfX1xyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhIGxvZyBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHdpdGggdGhlIGdpdmVuIGxldmVsICsgbWVzc2FnZVxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbGV2ZWxcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmZ1bmN0aW9uIHNlbmRMb2dNZXNzYWdlKGxldmVsLCBtZXNzYWdlKSB7XHJcblxyXG5cdC8vIExvZyBNZXNzYWdlIGZvcm1hdDpcclxuXHQvLyBsW3R5cGVdW21lc3NhZ2VdXHJcblx0d2luZG93LldhaWxzSW52b2tlKCdMJyArIGxldmVsICsgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIHRyYWNlIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nVHJhY2UobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdUJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nUHJpbnQobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdQJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGRlYnVnIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRGVidWcobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdEJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGluZm8gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dJbmZvKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnSScsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiB3YXJuaW5nIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nV2FybmluZyhtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1cnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZXJyb3IgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dFcnJvcihtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0UnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZmF0YWwgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dGYXRhbChtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0YnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIExvZyBsZXZlbCB0byB0aGUgZ2l2ZW4gbG9nIGxldmVsXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IGxvZ2xldmVsXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0TG9nTGV2ZWwobG9nbGV2ZWwpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnUycsIGxvZ2xldmVsKTtcclxufVxyXG5cclxuLy8gTG9nIGxldmVsc1xyXG5leHBvcnQgY29uc3QgTG9nTGV2ZWwgPSB7XHJcblx0VFJBQ0U6IDEsXHJcblx0REVCVUc6IDIsXHJcblx0SU5GTzogMyxcclxuXHRXQVJOSU5HOiA0LFxyXG5cdEVSUk9SOiA1LFxyXG59O1xyXG4iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xyXG5cclxuLyoqXHJcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxyXG4gKlxyXG4gKiBAY2xhc3MgTGlzdGVuZXJcclxuICovXHJcbmNsYXNzIExpc3RlbmVyIHtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cclxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcclxuICAgICAgICBtYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XHJcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxyXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcclxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcclxuICAgICAgICAgICAgY2FsbGJhY2suYXBwbHkobnVsbCwgZGF0YSk7XHJcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcclxuICAgICAgICAgICAgaWYgKG1heENhbGxiYWNrcyA9PT0gLTEpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXHJcbiAgICAgICAgICAgIG1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgICAgICByZXR1cm4gbWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IHt9O1xyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xyXG4gICAgRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcclxufVxyXG5cclxuZnVuY3Rpb24gbm90aWZ5TGlzdGVuZXJzKGV2ZW50RGF0YSkge1xyXG5cclxuICAgIC8vIEdldCB0aGUgZXZlbnQgbmFtZVxyXG4gICAgbGV0IGV2ZW50TmFtZSA9IGV2ZW50RGF0YS5uYW1lO1xyXG5cclxuICAgIC8vIENoZWNrIGlmIHdlIGhhdmUgYW55IGxpc3RlbmVycyBmb3IgdGhpcyBldmVudFxyXG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0pIHtcclxuXHJcbiAgICAgICAgLy8gS2VlcCBhIGxpc3Qgb2YgbGlzdGVuZXIgaW5kZXhlcyB0byBkZXN0cm95XHJcbiAgICAgICAgY29uc3QgbmV3RXZlbnRMaXN0ZW5lckxpc3QgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnNsaWNlKCk7XHJcblxyXG4gICAgICAgIC8vIEl0ZXJhdGUgbGlzdGVuZXJzXHJcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSAwOyBjb3VudCA8IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoOyBjb3VudCArPSAxKSB7XHJcblxyXG4gICAgICAgICAgICAvLyBHZXQgbmV4dCBsaXN0ZW5lclxyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV1bY291bnRdO1xyXG5cclxuICAgICAgICAgICAgbGV0IGRhdGEgPSBldmVudERhdGEuZGF0YTtcclxuXHJcbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xyXG4gICAgICAgICAgICBjb25zdCBkZXN0cm95ID0gbGlzdGVuZXIuQ2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgICAgIGlmIChkZXN0cm95KSB7XHJcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxyXG4gICAgICAgICAgICAgICAgbmV3RXZlbnRMaXN0ZW5lckxpc3Quc3BsaWNlKGNvdW50LCAxKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gVXBkYXRlIGNhbGxiYWNrcyB3aXRoIG5ldyBsaXN0IG9mIGxpc3RlbmVyc1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBuZXdFdmVudExpc3RlbmVyTGlzdDtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcclxuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXHJcbiAgICBsZXQgbWVzc2FnZTtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcbiAgICB9XHJcbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcclxuXHJcbiAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcclxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxyXG4gICAgfTtcclxuXHJcbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXHJcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSkge1xyXG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyc1xyXG4gICAgZGVsZXRlIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV07XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFWCcgKyBldmVudE5hbWUpO1xyXG59XHJcblxyXG4vKipcclxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXHJcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJlcyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXHJcbiAqIFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIFxyXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzIFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XHJcbiAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXHJcbiAgICBcclxuICAgIGlmIChhZGRpdGlvbmFsRXZlbnROYW1lcy5sZW5ndGggPiAwKSB7XHJcbiAgICAgICAgYWRkaXRpb25hbEV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xyXG4gICAgICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXHJcbiAgICAgICAgfSlcclxuICAgIH1cclxufSIsICIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX19cclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbmV4cG9ydCBjb25zdCBjYWxsYmFja3MgPSB7fTtcclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgbnVtYmVyIGZyb20gdGhlIG5hdGl2ZSBicm93c2VyIHJhbmRvbSBmdW5jdGlvblxyXG4gKlxyXG4gKiBAcmV0dXJucyBudW1iZXJcclxuICovXHJcbmZ1bmN0aW9uIGNyeXB0b1JhbmRvbSgpIHtcclxuXHR2YXIgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XHJcblx0cmV0dXJuIHdpbmRvdy5jcnlwdG8uZ2V0UmFuZG9tVmFsdWVzKGFycmF5KVswXTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBudW1iZXIgdXNpbmcgZGEgb2xkLXNrb29sIE1hdGguUmFuZG9tXHJcbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cclxuICpcclxuICogQHJldHVybnMgbnVtYmVyXHJcbiAqL1xyXG5mdW5jdGlvbiBiYXNpY1JhbmRvbSgpIHtcclxuXHRyZXR1cm4gTWF0aC5yYW5kb20oKSAqIDkwMDcxOTkyNTQ3NDA5OTE7XHJcbn1cclxuXHJcbi8vIFBpY2sgYSByYW5kb20gbnVtYmVyIGZ1bmN0aW9uIGJhc2VkIG9uIGJyb3dzZXIgY2FwYWJpbGl0eVxyXG52YXIgcmFuZG9tRnVuYztcclxuaWYgKHdpbmRvdy5jcnlwdG8pIHtcclxuXHRyYW5kb21GdW5jID0gY3J5cHRvUmFuZG9tO1xyXG59IGVsc2Uge1xyXG5cdHJhbmRvbUZ1bmMgPSBiYXNpY1JhbmRvbTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXHJcbiAqIGdpdmVuIGRhdGEuIEEgcHJvbWlzZSBpcyByZXR1cm5lZCBhbmQgd2lsbCBiZSBjb21wbGV0ZWQgd2hlbiB0aGVcclxuICogYmFja2VuZCByZXNwb25kcy4gVGhpcyB3aWxsIGJlIHJlc29sdmVkIHdoZW4gdGhlIGNhbGwgd2FzIHN1Y2Nlc3NmdWxcclxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXHJcbiAqIFRoZXJlIGlzIGEgdGltZW91dCBtZWNoYW5pc20uIElmIHRoZSBjYWxsIGRvZXNuJ3QgcmVzcG9uZCBpbiB0aGUgZ2l2ZW5cclxuICogdGltZSAoaW4gbWlsbGlzZWNvbmRzKSB0aGVuIHRoZSBwcm9taXNlIGlzIHJlamVjdGVkLlxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lXHJcbiAqIEBwYXJhbSB7YW55PX0gYXJnc1xyXG4gKiBAcGFyYW0ge251bWJlcj19IHRpbWVvdXRcclxuICogQHJldHVybnNcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG5hbWUsIGFyZ3MsIHRpbWVvdXQpIHtcclxuXHJcblx0Ly8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XHJcblx0aWYgKHRpbWVvdXQgPT0gbnVsbCkge1xyXG5cdFx0dGltZW91dCA9IDA7XHJcblx0fVxyXG5cclxuXHQvLyBDcmVhdGUgYSBwcm9taXNlXHJcblx0cmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcclxuXHJcblx0XHQvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxyXG5cdFx0dmFyIGNhbGxiYWNrSUQ7XHJcblx0XHRkbyB7XHJcblx0XHRcdGNhbGxiYWNrSUQgPSBuYW1lICsgJy0nICsgcmFuZG9tRnVuYygpO1xyXG5cdFx0fSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcclxuXHJcblx0XHR2YXIgdGltZW91dEhhbmRsZTtcclxuXHRcdC8vIFNldCB0aW1lb3V0XHJcblx0XHRpZiAodGltZW91dCA+IDApIHtcclxuXHRcdFx0dGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xyXG5cdFx0XHRcdHJlamVjdChFcnJvcignQ2FsbCB0byAnICsgbmFtZSArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xyXG5cdFx0XHR9LCB0aW1lb3V0KTtcclxuXHRcdH1cclxuXHJcblx0XHQvLyBTdG9yZSBjYWxsYmFja1xyXG5cdFx0Y2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xyXG5cdFx0XHR0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxyXG5cdFx0XHRyZWplY3Q6IHJlamVjdCxcclxuXHRcdFx0cmVzb2x2ZTogcmVzb2x2ZVxyXG5cdFx0fTtcclxuXHJcblx0XHR0cnkge1xyXG5cdFx0XHRjb25zdCBwYXlsb2FkID0ge1xyXG5cdFx0XHRcdG5hbWUsXHJcblx0XHRcdFx0YXJncyxcclxuXHRcdFx0XHRjYWxsYmFja0lELFxyXG5cdFx0XHR9O1xyXG5cclxuICAgICAgICAgICAgLy8gTWFrZSB0aGUgY2FsbFxyXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xyXG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXHJcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XHJcbiAgICAgICAgfVxyXG4gICAgfSk7XHJcbn1cclxuXHJcbndpbmRvdy5PYmZ1c2NhdGVkQ2FsbCA9IChpZCwgYXJncywgdGltZW91dCkgPT4ge1xyXG5cclxuICAgIC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxyXG4gICAgaWYgKHRpbWVvdXQgPT0gbnVsbCkge1xyXG4gICAgICAgIHRpbWVvdXQgPSAwO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIENyZWF0ZSBhIHByb21pc2VcclxuICAgIHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XHJcblxyXG4gICAgICAgIC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXHJcbiAgICAgICAgdmFyIGNhbGxiYWNrSUQ7XHJcbiAgICAgICAgZG8ge1xyXG4gICAgICAgICAgICBjYWxsYmFja0lEID0gaWQgKyAnLScgKyByYW5kb21GdW5jKCk7XHJcbiAgICAgICAgfSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcclxuXHJcbiAgICAgICAgdmFyIHRpbWVvdXRIYW5kbGU7XHJcbiAgICAgICAgLy8gU2V0IHRpbWVvdXRcclxuICAgICAgICBpZiAodGltZW91dCA+IDApIHtcclxuICAgICAgICAgICAgdGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKCdDYWxsIHRvIG1ldGhvZCAnICsgaWQgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcclxuICAgICAgICAgICAgfSwgdGltZW91dCk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBTdG9yZSBjYWxsYmFja1xyXG4gICAgICAgIGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcclxuICAgICAgICAgICAgdGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcclxuICAgICAgICAgICAgcmVqZWN0OiByZWplY3QsXHJcbiAgICAgICAgICAgIHJlc29sdmU6IHJlc29sdmVcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgICAgICAgICAgaWQsXHJcbiAgICAgICAgICAgICAgICBhcmdzLFxyXG4gICAgICAgICAgICAgICAgY2FsbGJhY2tJRCxcclxuICAgICAgICAgICAgfTtcclxuXHJcbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcclxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKCdjJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcclxuICAgICAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGUpO1xyXG4gICAgICAgIH1cclxuICAgIH0pO1xyXG59O1xyXG5cclxuXHJcbi8qKlxyXG4gKiBDYWxsZWQgYnkgdGhlIGJhY2tlbmQgdG8gcmV0dXJuIGRhdGEgdG8gYSBwcmV2aW91c2x5IGNhbGxlZFxyXG4gKiBiaW5kaW5nIGludm9jYXRpb25cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gaW5jb21pbmdNZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbGJhY2soaW5jb21pbmdNZXNzYWdlKSB7XHJcblx0Ly8gUGFyc2UgdGhlIG1lc3NhZ2VcclxuXHRsZXQgbWVzc2FnZTtcclxuXHR0cnkge1xyXG5cdFx0bWVzc2FnZSA9IEpTT04ucGFyc2UoaW5jb21pbmdNZXNzYWdlKTtcclxuXHR9IGNhdGNoIChlKSB7XHJcblx0XHRjb25zdCBlcnJvciA9IGBJbnZhbGlkIEpTT04gcGFzc2VkIHRvIGNhbGxiYWNrOiAke2UubWVzc2FnZX0uIE1lc3NhZ2U6ICR7aW5jb21pbmdNZXNzYWdlfWA7XHJcblx0XHRydW50aW1lLkxvZ0RlYnVnKGVycm9yKTtcclxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcblx0fVxyXG5cdGxldCBjYWxsYmFja0lEID0gbWVzc2FnZS5jYWxsYmFja2lkO1xyXG5cdGxldCBjYWxsYmFja0RhdGEgPSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XHJcblx0aWYgKCFjYWxsYmFja0RhdGEpIHtcclxuXHRcdGNvbnN0IGVycm9yID0gYENhbGxiYWNrICcke2NhbGxiYWNrSUR9JyBub3QgcmVnaXN0ZXJlZCEhIWA7XHJcblx0XHRjb25zb2xlLmVycm9yKGVycm9yKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZVxyXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcclxuXHR9XHJcblx0Y2xlYXJUaW1lb3V0KGNhbGxiYWNrRGF0YS50aW1lb3V0SGFuZGxlKTtcclxuXHJcblx0ZGVsZXRlIGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcclxuXHJcblx0aWYgKG1lc3NhZ2UuZXJyb3IpIHtcclxuXHRcdGNhbGxiYWNrRGF0YS5yZWplY3QobWVzc2FnZS5lcnJvcik7XHJcblx0fSBlbHNlIHtcclxuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcclxuXHR9XHJcbn1cclxuIiwgIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfXyAgICBcclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKSBcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vICBcclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSAnLi9jYWxscyc7XHJcblxyXG4vLyBUaGlzIGlzIHdoZXJlIHdlIGJpbmQgZ28gbWV0aG9kIHdyYXBwZXJzXHJcbndpbmRvdy5nbyA9IHt9O1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFNldEJpbmRpbmdzKGJpbmRpbmdzTWFwKSB7XHJcblx0dHJ5IHtcclxuXHRcdGJpbmRpbmdzTWFwID0gSlNPTi5wYXJzZShiaW5kaW5nc01hcCk7XHJcblx0fSBjYXRjaCAoZSkge1xyXG5cdFx0Y29uc29sZS5lcnJvcihlKTtcclxuXHR9XHJcblxyXG5cdC8vIEluaXRpYWxpc2UgdGhlIGJpbmRpbmdzIG1hcFxyXG5cdHdpbmRvdy5nbyA9IHdpbmRvdy5nbyB8fCB7fTtcclxuXHJcblx0Ly8gSXRlcmF0ZSBwYWNrYWdlIG5hbWVzXHJcblx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXApLmZvckVhY2goKHBhY2thZ2VOYW1lKSA9PiB7XHJcblxyXG5cdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XHJcblx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXSB8fCB7fTtcclxuXHJcblx0XHQvLyBJdGVyYXRlIHN0cnVjdCBuYW1lc1xyXG5cdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdKS5mb3JFYWNoKChzdHJ1Y3ROYW1lKSA9PiB7XHJcblxyXG5cdFx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3RcclxuXHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gfHwge307XHJcblxyXG5cdFx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0pLmZvckVhY2goKG1ldGhvZE5hbWUpID0+IHtcclxuXHJcblx0XHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXVttZXRob2ROYW1lXSA9IGZ1bmN0aW9uICgpIHtcclxuXHJcblx0XHRcdFx0XHQvLyBObyB0aW1lb3V0IGJ5IGRlZmF1bHRcclxuXHRcdFx0XHRcdGxldCB0aW1lb3V0ID0gMDtcclxuXHJcblx0XHRcdFx0XHQvLyBBY3R1YWwgZnVuY3Rpb25cclxuXHRcdFx0XHRcdGZ1bmN0aW9uIGR5bmFtaWMoKSB7XHJcblx0XHRcdFx0XHRcdGNvbnN0IGFyZ3MgPSBbXS5zbGljZS5jYWxsKGFyZ3VtZW50cyk7XHJcblx0XHRcdFx0XHRcdHJldHVybiBDYWxsKFtwYWNrYWdlTmFtZSwgc3RydWN0TmFtZSwgbWV0aG9kTmFtZV0uam9pbignLicpLCBhcmdzLCB0aW1lb3V0KTtcclxuXHRcdFx0XHRcdH1cclxuXHJcblx0XHRcdFx0XHQvLyBBbGxvdyBzZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cclxuXHRcdFx0XHRcdGR5bmFtaWMuc2V0VGltZW91dCA9IGZ1bmN0aW9uIChuZXdUaW1lb3V0KSB7XHJcblx0XHRcdFx0XHRcdHRpbWVvdXQgPSBuZXdUaW1lb3V0O1xyXG5cdFx0XHRcdFx0fTtcclxuXHJcblx0XHRcdFx0XHQvLyBBbGxvdyBnZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cclxuXHRcdFx0XHRcdGR5bmFtaWMuZ2V0VGltZW91dCA9IGZ1bmN0aW9uICgpIHtcclxuXHRcdFx0XHRcdFx0cmV0dXJuIHRpbWVvdXQ7XHJcblx0XHRcdFx0XHR9O1xyXG5cclxuXHRcdFx0XHRcdHJldHVybiBkeW5hbWljO1xyXG5cdFx0XHRcdH0oKTtcclxuXHRcdFx0fSk7XHJcblx0XHR9KTtcclxuXHR9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWQoKSB7XHJcbiAgICB3aW5kb3cubG9jYXRpb24ucmVsb2FkKCk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWRBcHAoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dSJyk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTeXN0ZW1EZWZhdWx0VGhlbWUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBU0RUJyk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRMaWdodFRoZW1lKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQUxUJyk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXREYXJrVGhlbWUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBRFQnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFBsYWNlIHRoZSB3aW5kb3cgaW4gdGhlIGNlbnRlciBvZiB0aGUgc2NyZWVuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dDZW50ZXIoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1djJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSB3aW5kb3cgdGl0bGVcclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IHRpdGxlXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRUaXRsZSh0aXRsZSkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVCcgKyB0aXRsZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYWtlcyB0aGUgd2luZG93IGdvIGZ1bGxzY3JlZW5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0Z1bGxzY3JlZW4oKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dGJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXZlcnRzIHRoZSB3aW5kb3cgZnJvbSBmdWxsc2NyZWVuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbmZ1bGxzY3JlZW4oKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dmJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBpbiBmdWxsIHNjcmVlbiBtb2RlIG9yIG5vdC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzRnVsbHNjcmVlbigpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzRnVsbHNjcmVlblwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3M6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTx7dzogbnVtYmVyLCBoOiBudW1iZXJ9Pn0gVGhlIHNpemUgb2YgdGhlIHdpbmRvd1xyXG5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRTaXplKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0U2l6ZVwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1heFNpemUod2lkdGgsIGhlaWdodCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXWjonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d6OicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XHJcbn1cclxuXHJcblxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgd2luZG93IEFsd2F5c09uVG9wIG9yIG5vdCBvbiB0b3BcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEFsd2F5c09uVG9wKGIpIHtcclxuXHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBVFA6JyArIChiID8gJzEnIDogJzAnKSk7XHJcbn1cclxuXHJcblxyXG5cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0geFxyXG4gKiBAcGFyYW0ge251bWJlcn0geVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFBvc2l0aW9uKHgsIHkpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3A6JyArIHggKyAnOicgKyB5KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8e3g6IG51bWJlciwgeTogbnVtYmVyfT59IFRoZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0UG9zaXRpb24oKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dHZXRQb3NcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIaWRlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0hpZGUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dIJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93IHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1Nob3coKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dTJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXhpbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNYXhpbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV00nKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFRvZ2dsZSB0aGUgTWF4aW1pc2Ugb2YgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VG9nZ2xlTWF4aW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d0Jyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBVbm1heGltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWF4aW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dVJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQgb3Igbm90LlxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNNYXhpbWlzZWQoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc01heGltaXNlZFwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1pbmltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd01pbmltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXbScpO1xyXG59XHJcblxyXG4vKipcclxuICogVW5taW5pbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1pbmltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdScpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkIG9yIG5vdC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWluaW1pc2VkKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNaW5pbWlzZWRcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBub3JtYWwgb3Igbm90LlxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNOb3JtYWwoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc05vcm1hbFwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gUiBSZWRcclxuICogQHBhcmFtIHtudW1iZXJ9IEcgR3JlZW5cclxuICogQHBhcmFtIHtudW1iZXJ9IEIgQmx1ZVxyXG4gKiBAcGFyYW0ge251bWJlcn0gQSBBbHBoYVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEJhY2tncm91bmRDb2xvdXIoUiwgRywgQiwgQSkge1xyXG4gICAgbGV0IHJnYmEgPSBKU09OLnN0cmluZ2lmeSh7cjogUiB8fCAwLCBnOiBHIHx8IDAsIGI6IEIgfHwgMCwgYTogQSB8fCAyNTV9KTtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcclxuXHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgYWxsIHNjcmVlbnMuIENhbGwgdGhpcyBhbmV3IGVhY2ggdGltZSB5b3Ugd2FudCB0byByZWZyZXNoIGRhdGEgZnJvbSB0aGUgdW5kZXJseWluZyB3aW5kb3dpbmcgc3lzdGVtLlxyXG4gKiBAZXhwb3J0XHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoJy4uL3dyYXBwZXIvcnVudGltZScpLlNjcmVlbn0gU2NyZWVuXHJcbiAqIEByZXR1cm4ge1Byb21pc2U8e1NjcmVlbltdfT59IFRoZSBzY3JlZW5zXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2NyZWVuR2V0QWxsKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6U2NyZWVuR2V0QWxsXCIpO1xyXG59XHJcbiIsICIvKipcclxuICogQGRlc2NyaXB0aW9uOiBVc2UgdGhlIHN5c3RlbSBkZWZhdWx0IGJyb3dzZXIgdG8gb3BlbiB0aGUgdXJsXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnJvd3Nlck9wZW5VUkwodXJsKSB7XHJcbiAgd2luZG93LldhaWxzSW52b2tlKCdCTzonICsgdXJsKTtcclxufSIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcclxuaW1wb3J0IHtldmVudExpc3RlbmVycywgRXZlbnRzRW1pdCwgRXZlbnRzTm90aWZ5LCBFdmVudHNPZmYsIEV2ZW50c09uLCBFdmVudHNPbmNlLCBFdmVudHNPbk11bHRpcGxlfSBmcm9tICcuL2V2ZW50cyc7XHJcbmltcG9ydCB7Q2FsbCwgQ2FsbGJhY2ssIGNhbGxiYWNrc30gZnJvbSAnLi9jYWxscyc7XHJcbmltcG9ydCB7U2V0QmluZGluZ3N9IGZyb20gXCIuL2JpbmRpbmdzXCI7XHJcbmltcG9ydCAqIGFzIFdpbmRvdyBmcm9tIFwiLi93aW5kb3dcIjtcclxuaW1wb3J0ICogYXMgU2NyZWVuIGZyb20gXCIuL3NjcmVlblwiO1xyXG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXJcIjtcclxuXHJcblxyXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUScpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUycpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnSCcpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpFbnZpcm9ubWVudFwiKTtcclxufVxyXG5cclxuLy8gVGhlIEpTIHJ1bnRpbWVcclxud2luZG93LnJ1bnRpbWUgPSB7XHJcbiAgICAuLi5Mb2csXHJcbiAgICAuLi5XaW5kb3csXHJcbiAgICAuLi5Ccm93c2VyLFxyXG4gICAgLi4uU2NyZWVuLFxyXG4gICAgRXZlbnRzT24sXHJcbiAgICBFdmVudHNPbmNlLFxyXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcclxuICAgIEV2ZW50c0VtaXQsXHJcbiAgICBFdmVudHNPZmYsXHJcbiAgICBFbnZpcm9ubWVudCxcclxuICAgIFNob3csXHJcbiAgICBIaWRlLFxyXG4gICAgUXVpdFxyXG59O1xyXG5cclxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXHJcbndpbmRvdy53YWlscyA9IHtcclxuICAgIENhbGxiYWNrLFxyXG4gICAgRXZlbnRzTm90aWZ5LFxyXG4gICAgU2V0QmluZGluZ3MsXHJcbiAgICBldmVudExpc3RlbmVycyxcclxuICAgIGNhbGxiYWNrcyxcclxuICAgIGZsYWdzOiB7XHJcbiAgICAgICAgZGlzYWJsZVNjcm9sbGJhckRyYWc6IGZhbHNlLFxyXG4gICAgICAgIGRpc2FibGVXYWlsc0RlZmF1bHRDb250ZXh0TWVudTogZmFsc2UsXHJcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcclxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxyXG4gICAgICAgIGJvcmRlclRoaWNrbmVzczogNixcclxuICAgICAgICBzaG91bGREcmFnOiBmYWxzZVxyXG4gICAgfVxyXG59O1xyXG5cclxuLy8gU2V0IHRoZSBiaW5kaW5nc1xyXG5pZiAod2luZG93LndhaWxzYmluZGluZ3MpIHtcclxuICAgIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncyh3aW5kb3cud2FpbHNiaW5kaW5ncyk7XHJcbiAgICBkZWxldGUgd2luZG93LndhaWxzLlNldEJpbmRpbmdzO1xyXG59XHJcblxyXG4vLyBUaGlzIGlzIGV2YWx1YXRlZCBhdCBidWlsZCB0aW1lIGluIHBhY2thZ2UuanNvblxyXG4vLyBjb25zdCBkZXYgPSAwO1xyXG4vLyBjb25zdCBwcm9kdWN0aW9uID0gMTtcclxuaWYgKEVOViA9PT0gMSkge1xyXG4gICAgZGVsZXRlIHdpbmRvdy53YWlsc2JpbmRpbmdzO1xyXG59XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsICgpID0+IHtcclxuICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XHJcbn0pO1xyXG5cclxubGV0IGNzc0RyYWdUZXN0ID0gZnVuY3Rpb24gKGUpIHtcclxuICAgIHJldHVybiB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZSgnLS13YWlscy1kcmFnZ2FibGUnKSA9PT0gXCJkcmFnXCI7XHJcbn07XHJcblxyXG4vLyBFbGVtZW50IGRyYWcgaGFuZGxlclxyXG4vLyBCYXNlZCBvbiBjb2RlIGZyb206IGh0dHBzOi8vZ2l0aHViLmNvbS9wYXRyMG51cy9EZXNrR2FwXHJcbmxldCBlbGVtZW50RHJhZ1Rlc3QgPSBmdW5jdGlvbiAoZSkge1xyXG4gICAgLy8gQ2hlY2sgZm9yIGRyYWdnaW5nXHJcbiAgICBsZXQgY3VycmVudEVsZW1lbnQgPSBlLnRhcmdldDtcclxuICAgIHdoaWxlIChjdXJyZW50RWxlbWVudCAhPSBudWxsKSB7XHJcbiAgICAgICAgaWYgKGN1cnJlbnRFbGVtZW50Lmhhc0F0dHJpYnV0ZSgnZGF0YS13YWlscy1uby1kcmFnJykpIHtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfSBlbHNlIGlmIChjdXJyZW50RWxlbWVudC5oYXNBdHRyaWJ1dGUoJ2RhdGEtd2FpbHMtZHJhZycpKSB7XHJcbiAgICAgICAgICAgIHJldHVybiB0cnVlO1xyXG4gICAgICAgIH1cclxuICAgICAgICBjdXJyZW50RWxlbWVudCA9IGN1cnJlbnRFbGVtZW50LnBhcmVudEVsZW1lbnQ7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gZmFsc2U7XHJcbn07XHJcblxyXG5sZXQgZHJhZ1Rlc3QgPSBlbGVtZW50RHJhZ1Rlc3Q7XHJcblxyXG53aW5kb3cud2FpbHMudXNlQ1NTRHJhZyA9IGZ1bmN0aW9uICh0KSB7XHJcbiAgICBpZiAodCA9PT0gZmFsc2UpIHtcclxuICAgICAgICBjb25zb2xlLmxvZyhcIlVzaW5nIG9yaWdpbmFsIGRyYWcgZGV0ZWN0aW9uXCIpO1xyXG4gICAgICAgIGRyYWdUZXN0ID0gZWxlbWVudERyYWdUZXN0O1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBjb25zb2xlLmxvZyhcIlVzaW5nIENTUyBkcmFnIGRldGVjdGlvblwiKTtcclxuICAgICAgICBkcmFnVGVzdCA9IGNzc0RyYWdUZXN0O1xyXG4gICAgfVxyXG59O1xyXG5cclxuXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xyXG5cclxuICAgIC8vIENoZWNrIGZvciByZXNpemluZ1xyXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKSB7XHJcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwicmVzaXplOlwiICsgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpO1xyXG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGRyYWdUZXN0KGUpKSB7XHJcbiAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlU2Nyb2xsYmFyRHJhZykge1xyXG4gICAgICAgICAgICAvLyBUaGlzIGNoZWNrcyBmb3IgY2xpY2tzIG9uIHRoZSBzY3JvbGwgYmFyXHJcbiAgICAgICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IHRydWU7XHJcbiAgICB9XHJcblxyXG59KTtcclxuXHJcbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IpIHtcclxuICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gY3Vyc29yIHx8IHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yO1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgPSBjdXJzb3I7XHJcbn1cclxuXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBmdW5jdGlvbiAoZSkge1xyXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnKSB7XHJcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVSZXNpemUpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPT0gbnVsbCkge1xyXG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yID0gZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3I7XHJcbiAgICB9XHJcbiAgICBpZiAod2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzICYmIHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MpIHtcclxuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IFwic2UtcmVzaXplXCI7XHJcbiAgICB9XHJcbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xyXG5cclxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XHJcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9IGVsc2UgaWYgKHJpZ2h0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgdG9wQm9yZGVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIgJiYgcmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XHJcblxyXG59KTtcclxuXHJcbi8vIFNldHVwIGNvbnRleHQgbWVudSBob29rXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGZ1bmN0aW9uIChlKSB7XHJcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRpc2FibGVXYWlsc0RlZmF1bHRDb250ZXh0TWVudSkge1xyXG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgIH1cclxufSk7XHJcblxyXG53aW5kb3cuV2FpbHNJbnZva2UoXCJydW50aW1lOnJlYWR5XCIpOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBa0JBLFdBQVMsZUFBZSxPQUFPLFNBQVM7QUFJdkMsV0FBTyxZQUFZLE1BQU0sUUFBUSxPQUFPO0FBQUEsRUFDekM7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFFBQVEsU0FBUztBQUNoQyxtQkFBZSxL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
