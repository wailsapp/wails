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
    constructor(eventName, callback, maxCallbacks) {
      this.eventName = eventName;
      this.maxCallbacks = maxCallbacks || -1;
      this.Callback = (data) => {
        callback.apply(null, data);
        if (this.maxCallbacks === -1) {
          return false;
        }
        this.maxCallbacks -= 1;
        return this.maxCallbacks === 0;
      };
    }
  };
  var eventListeners = {};
  function EventsOnMultiple(eventName, callback, maxCallbacks) {
    eventListeners[eventName] = eventListeners[eventName] || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    eventListeners[eventName].push(thisListener);
    return () => listenerOff(thisListener);
  }
  function EventsOn(eventName, callback) {
    return EventsOnMultiple(eventName, callback, -1);
  }
  function EventsOnce(eventName, callback) {
    return EventsOnMultiple(eventName, callback, 1);
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
      if (newEventListenerList.length === 0) {
        removeListener(eventName);
      } else {
        eventListeners[eventName] = newEventListenerList;
      }
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
  function listenerOff(listener) {
    const eventName = listener.eventName;
    eventListeners[eventName] = eventListeners[eventName].filter((l) => l !== listener);
    if (eventListeners[eventName].length === 0) {
      removeListener(eventName);
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
    return new Promise(function(resolve, reject) {
      var callbackID;
      do {
        callbackID = id + "-" + randomFunc();
      } while (callbacks[callbackID]);
      var timeoutHandle;
      if (timeout > 0) {
        timeoutHandle = setTimeout(function() {
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

  // desktop/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    ClipboardGetText: () => ClipboardGetText,
    ClipboardSetText: () => ClipboardSetText
  });
  function ClipboardSetText(text) {
    return Call(":wails:ClipboardSetText", [text]);
  }
  function ClipboardGetText() {
    return Call(":wails:ClipboardGetText");
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
    ...clipboard_exports,
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
      shouldDrag: false,
      deferDragToMouseMove: false,
      cssDragProperty: "--wails-draggable",
      cssDragValue: "drag"
    }
  };
  if (window.wailsbindings) {
    window.wails.SetBindings(window.wailsbindings);
    delete window.wails.SetBindings;
  }
  if (false) {
    delete window.wailsbindings;
  }
  var dragTest = function(e) {
    var val = window.getComputedStyle(e.target).getPropertyValue(window.wails.flags.cssDragProperty);
    if (val) {
      val = val.trim();
    }
    if (val !== window.wails.flags.cssDragValue) {
      return false;
    }
    if (e.buttons !== 1) {
      return false;
    }
    if (e.detail !== 1) {
      return false;
    }
    return true;
  };
  window.wails.setCSSDragProperties = function(property, value) {
    window.wails.flags.cssDragProperty = property;
    window.wails.flags.cssDragValue = value;
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
      if (window.wails.flags.deferDragToMouseMove) {
        window.wails.flags.shouldDrag = true;
      } else {
        e.preventDefault();
        window.WailsInvoke("drag");
      }
      return;
    } else {
      window.wails.flags.shouldDrag = false;
    }
  });
  window.addEventListener("mouseup", () => {
    window.wails.flags.shouldDrag = false;
  });
  function setResize(cursor) {
    document.documentElement.style.cursor = cursor || window.wails.flags.defaultCursor;
    window.wails.flags.resizeEdge = cursor;
  }
  window.addEventListener("mousemove", function(e) {
    if (window.wails.flags.shouldDrag) {
      let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
      if (mousePressed <= 0) {
        window.wails.flags.shouldDrag = false;
      } else {
        window.WailsInvoke("drag");
        return;
      }
    }
    if (!window.wails.flags.enableResize) {
      return;
    }
    if (window.wails.flags.defaultCursor == null) {
      window.wails.flags.defaultCursor = document.documentElement.style.cursor;
    }
    if (window.outerWidth - e.clientX < window.wails.flags.borderThickness && window.outerHeight - e.clientY < window.wails.flags.borderThickness) {
      document.documentElement.style.cursor = "se-resize";
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbi8qKlxuICogU2VuZHMgYSBsb2cgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB3aXRoIHRoZSBnaXZlbiBsZXZlbCArIG1lc3NhZ2VcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gbGV2ZWxcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmZ1bmN0aW9uIHNlbmRMb2dNZXNzYWdlKGxldmVsLCBtZXNzYWdlKSB7XG5cblx0Ly8gTG9nIE1lc3NhZ2UgZm9ybWF0OlxuXHQvLyBsW3R5cGVdW21lc3NhZ2VdXG5cdHdpbmRvdy5XYWlsc0ludm9rZSgnTCcgKyBsZXZlbCArIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gdHJhY2UgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1RyYWNlKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1QnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dQcmludChtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdQJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBkZWJ1ZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRGVidWcobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gaW5mbyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nSW5mbyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdJJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiB3YXJuaW5nIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dXYXJuaW5nKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1cnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGVycm9yIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dFcnJvcihtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdFJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBmYXRhbCBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRmF0YWwobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRicsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIExvZyBsZXZlbCB0byB0aGUgZ2l2ZW4gbG9nIGxldmVsXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IGxvZ2xldmVsXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRMb2dMZXZlbChsb2dsZXZlbCkge1xuXHRzZW5kTG9nTWVzc2FnZSgnUycsIGxvZ2xldmVsKTtcbn1cblxuLy8gTG9nIGxldmVsc1xuZXhwb3J0IGNvbnN0IExvZ0xldmVsID0ge1xuXHRUUkFDRTogMSxcblx0REVCVUc6IDIsXG5cdElORk86IDMsXG5cdFdBUk5JTkc6IDQsXG5cdEVSUk9SOiA1LFxufTtcbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuLy8gRGVmaW5lcyBhIHNpbmdsZSBsaXN0ZW5lciB3aXRoIGEgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdG8gY2FsbGJhY2tcblxuLyoqXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcbiAqXG4gKiBAY2xhc3MgTGlzdGVuZXJcbiAqL1xuY2xhc3MgTGlzdGVuZXIge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICAgICAqIEBtZW1iZXJvZiBMaXN0ZW5lclxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xuICAgICAgICAgICAgY2FsbGJhY2suYXBwbHkobnVsbCwgZGF0YSk7XG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IHt9O1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuZnVuY3Rpb24gbm90aWZ5TGlzdGVuZXJzKGV2ZW50RGF0YSkge1xuXG4gICAgLy8gR2V0IHRoZSBldmVudCBuYW1lXG4gICAgbGV0IGV2ZW50TmFtZSA9IGV2ZW50RGF0YS5uYW1lO1xuXG4gICAgLy8gQ2hlY2sgaWYgd2UgaGF2ZSBhbnkgbGlzdGVuZXJzIGZvciB0aGlzIGV2ZW50XG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0pIHtcblxuICAgICAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lciBpbmRleGVzIHRvIGRlc3Ryb3lcbiAgICAgICAgY29uc3QgbmV3RXZlbnRMaXN0ZW5lckxpc3QgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnNsaWNlKCk7XG5cbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSAwOyBjb3VudCA8IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoOyBjb3VudCArPSAxKSB7XG5cbiAgICAgICAgICAgIC8vIEdldCBuZXh0IGxpc3RlbmVyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV1bY291bnRdO1xuXG4gICAgICAgICAgICBsZXQgZGF0YSA9IGV2ZW50RGF0YS5kYXRhO1xuXG4gICAgICAgICAgICAvLyBEbyB0aGUgY2FsbGJhY2tcbiAgICAgICAgICAgIGNvbnN0IGRlc3Ryb3kgPSBsaXN0ZW5lci5DYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIGlmIChkZXN0cm95KSB7XG4gICAgICAgICAgICAgICAgLy8gaWYgdGhlIGxpc3RlbmVyIGluZGljYXRlZCB0byBkZXN0cm95IGl0c2VsZiwgYWRkIGl0IHRvIHRoZSBkZXN0cm95IGxpc3RcbiAgICAgICAgICAgICAgICBuZXdFdmVudExpc3RlbmVyTGlzdC5zcGxpY2UoY291bnQsIDEpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gVXBkYXRlIGNhbGxiYWNrcyB3aXRoIG5ldyBsaXN0IG9mIGxpc3RlbmVyc1xuICAgICAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoID09PSAwKSB7XG4gICAgICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0O1xuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcblxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcbiAgICAvLyBQYXJzZSB0aGUgbWVzc2FnZVxuICAgIGxldCBtZXNzYWdlO1xuICAgIHRyeSB7XG4gICAgICAgIG1lc3NhZ2UgPSBKU09OLnBhcnNlKG5vdGlmeU1lc3NhZ2UpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuICAgIH1cbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XG59XG5cbi8qKlxuICogRW1pdCBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBkYXRhXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcblxuICAgIGNvbnN0IHBheWxvYWQgPSB7XG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcbiAgICAgICAgZGF0YTogW10uc2xpY2UuYXBwbHkoYXJndW1lbnRzKS5zbGljZSgxKSxcbiAgICB9O1xuXG4gICAgLy8gTm90aWZ5IEpTIGxpc3RlbmVyc1xuICAgIG5vdGlmeUxpc3RlbmVycyhwYXlsb2FkKTtcblxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VFJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcbn1cblxuZnVuY3Rpb24gcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKSB7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyc1xuICAgIGRlbGV0ZSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdO1xuXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRVgnICsgZXZlbnROYW1lKTtcbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPbixcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJlcyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG5cbiAgICBpZiAoYWRkaXRpb25hbEV2ZW50TmFtZXMubGVuZ3RoID4gMCkge1xuICAgICAgICBhZGRpdGlvbmFsRXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XG4gICAgICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG4gICAgICAgIH0pXG4gICAgfVxufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhbGwgZXZlbnQgbGlzdGVuZXJzIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uXG4gKi9cbiBleHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmQWxsKCkge1xuICAgIGNvbnN0IGV2ZW50TmFtZXMgPSBPYmplY3Qua2V5cyhldmVudExpc3RlbmVycyk7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgIT09IGV2ZW50TmFtZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lc1tpXSk7XG4gICAgfVxufVxuXG4vKipcbiAqIGxpc3RlbmVyT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggRXZlbnRzT25cbiAqXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxuICovXG4gZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuXG4gICAgLy8gQ2xlYW4gdXAgaWYgdGhlcmUgYXJlIG5vIGV2ZW50IGxpc3RlbmVycyBsZWZ0XG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XG4gICAgfVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gY3J5cHRvUmFuZG9tKCkge1xuXHR2YXIgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciB1c2luZyBkYSBvbGQtc2tvb2wgTWF0aC5SYW5kb21cbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XG5cdHJldHVybiBNYXRoLnJhbmRvbSgpICogOTAwNzE5OTI1NDc0MDk5MTtcbn1cblxuLy8gUGljayBhIHJhbmRvbSBudW1iZXIgZnVuY3Rpb24gYmFzZWQgb24gYnJvd3NlciBjYXBhYmlsaXR5XG52YXIgcmFuZG9tRnVuYztcbmlmICh3aW5kb3cuY3J5cHRvKSB7XG5cdHJhbmRvbUZ1bmMgPSBjcnlwdG9SYW5kb207XG59IGVsc2Uge1xuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XG59XG5cblxuLyoqXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXG4gKiBiYWNrZW5kIHJlc3BvbmRzLiBUaGlzIHdpbGwgYmUgcmVzb2x2ZWQgd2hlbiB0aGUgY2FsbCB3YXMgc3VjY2Vzc2Z1bFxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXG4gKiB0aW1lIChpbiBtaWxsaXNlY29uZHMpIHRoZW4gdGhlIHByb21pc2UgaXMgcmVqZWN0ZWQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWVcbiAqIEBwYXJhbSB7YW55PX0gYXJnc1xuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XG4gKiBAcmV0dXJuc1xuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XG5cblx0Ly8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcblx0XHR0aW1lb3V0ID0gMDtcblx0fVxuXG5cdC8vIENyZWF0ZSBhIHByb21pc2Vcblx0cmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcblxuXHRcdC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXG5cdFx0dmFyIGNhbGxiYWNrSUQ7XG5cdFx0ZG8ge1xuXHRcdFx0Y2FsbGJhY2tJRCA9IG5hbWUgKyAnLScgKyByYW5kb21GdW5jKCk7XG5cdFx0fSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcblxuXHRcdHZhciB0aW1lb3V0SGFuZGxlO1xuXHRcdC8vIFNldCB0aW1lb3V0XG5cdFx0aWYgKHRpbWVvdXQgPiAwKSB7XG5cdFx0XHR0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XG5cdFx0XHRcdHJlamVjdChFcnJvcignQ2FsbCB0byAnICsgbmFtZSArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xuXHRcdFx0fSwgdGltZW91dCk7XG5cdFx0fVxuXG5cdFx0Ly8gU3RvcmUgY2FsbGJhY2tcblx0XHRjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XG5cdFx0XHR0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxuXHRcdFx0cmVqZWN0OiByZWplY3QsXG5cdFx0XHRyZXNvbHZlOiByZXNvbHZlXG5cdFx0fTtcblxuXHRcdHRyeSB7XG5cdFx0XHRjb25zdCBwYXlsb2FkID0ge1xuXHRcdFx0XHRuYW1lLFxuXHRcdFx0XHRhcmdzLFxuXHRcdFx0XHRjYWxsYmFja0lELFxuXHRcdFx0fTtcblxuICAgICAgICAgICAgLy8gTWFrZSB0aGUgY2FsbFxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKCdDJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcbiAgICAgICAgfSBjYXRjaCAoZSkge1xuICAgICAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGUpO1xuICAgICAgICB9XG4gICAgfSk7XG59XG5cbndpbmRvdy5PYmZ1c2NhdGVkQ2FsbCA9IChpZCwgYXJncywgdGltZW91dCkgPT4ge1xuXG4gICAgLy8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG4gICAgaWYgKHRpbWVvdXQgPT0gbnVsbCkge1xuICAgICAgICB0aW1lb3V0ID0gMDtcbiAgICB9XG5cbiAgICAvLyBDcmVhdGUgYSBwcm9taXNlXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcblxuICAgICAgICAvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuICAgICAgICB2YXIgY2FsbGJhY2tJRDtcbiAgICAgICAgZG8ge1xuICAgICAgICAgICAgY2FsbGJhY2tJRCA9IGlkICsgJy0nICsgcmFuZG9tRnVuYygpO1xuICAgICAgICB9IHdoaWxlIChjYWxsYmFja3NbY2FsbGJhY2tJRF0pO1xuXG4gICAgICAgIHZhciB0aW1lb3V0SGFuZGxlO1xuICAgICAgICAvLyBTZXQgdGltZW91dFxuICAgICAgICBpZiAodGltZW91dCA+IDApIHtcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgICAgICByZWplY3QoRXJyb3IoJ0NhbGwgdG8gbWV0aG9kICcgKyBpZCArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xuICAgICAgICAgICAgfSwgdGltZW91dCk7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBTdG9yZSBjYWxsYmFja1xuICAgICAgICBjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxuICAgICAgICAgICAgcmVqZWN0OiByZWplY3QsXG4gICAgICAgICAgICByZXNvbHZlOiByZXNvbHZlXG4gICAgICAgIH07XG5cbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XG5cdFx0XHRcdGlkLFxuXHRcdFx0XHRhcmdzLFxuXHRcdFx0XHRjYWxsYmFja0lELFxuXHRcdFx0fTtcblxuICAgICAgICAgICAgLy8gTWFrZSB0aGUgY2FsbFxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKCdjJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcbiAgICAgICAgfSBjYXRjaCAoZSkge1xuICAgICAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGUpO1xuICAgICAgICB9XG4gICAgfSk7XG59O1xuXG5cbi8qKlxuICogQ2FsbGVkIGJ5IHRoZSBiYWNrZW5kIHRvIHJldHVybiBkYXRhIHRvIGEgcHJldmlvdXNseSBjYWxsZWRcbiAqIGJpbmRpbmcgaW52b2NhdGlvblxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBpbmNvbWluZ01lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGxiYWNrKGluY29taW5nTWVzc2FnZSkge1xuXHQvLyBQYXJzZSB0aGUgbWVzc2FnZVxuXHRsZXQgbWVzc2FnZTtcblx0dHJ5IHtcblx0XHRtZXNzYWdlID0gSlNPTi5wYXJzZShpbmNvbWluZ01lc3NhZ2UpO1xuXHR9IGNhdGNoIChlKSB7XG5cdFx0Y29uc3QgZXJyb3IgPSBgSW52YWxpZCBKU09OIHBhc3NlZCB0byBjYWxsYmFjazogJHtlLm1lc3NhZ2V9LiBNZXNzYWdlOiAke2luY29taW5nTWVzc2FnZX1gO1xuXHRcdHJ1bnRpbWUuTG9nRGVidWcoZXJyb3IpO1xuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XG5cdH1cblx0bGV0IGNhbGxiYWNrSUQgPSBtZXNzYWdlLmNhbGxiYWNraWQ7XG5cdGxldCBjYWxsYmFja0RhdGEgPSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XG5cdGlmICghY2FsbGJhY2tEYXRhKSB7XG5cdFx0Y29uc3QgZXJyb3IgPSBgQ2FsbGJhY2sgJyR7Y2FsbGJhY2tJRH0nIG5vdCByZWdpc3RlcmVkISEhYDtcblx0XHRjb25zb2xlLmVycm9yKGVycm9yKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZVxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XG5cdH1cblx0Y2xlYXJUaW1lb3V0KGNhbGxiYWNrRGF0YS50aW1lb3V0SGFuZGxlKTtcblxuXHRkZWxldGUgY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xuXG5cdGlmIChtZXNzYWdlLmVycm9yKSB7XG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChtZXNzYWdlLmVycm9yKTtcblx0fSBlbHNlIHtcblx0XHRjYWxsYmFja0RhdGEucmVzb2x2ZShtZXNzYWdlLnJlc3VsdCk7XG5cdH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fICAgIFxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApIFxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vICBcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSAnLi9jYWxscyc7XG5cbi8vIFRoaXMgaXMgd2hlcmUgd2UgYmluZCBnbyBtZXRob2Qgd3JhcHBlcnNcbndpbmRvdy5nbyA9IHt9O1xuXG5leHBvcnQgZnVuY3Rpb24gU2V0QmluZGluZ3MoYmluZGluZ3NNYXApIHtcblx0dHJ5IHtcblx0XHRiaW5kaW5nc01hcCA9IEpTT04ucGFyc2UoYmluZGluZ3NNYXApO1xuXHR9IGNhdGNoIChlKSB7XG5cdFx0Y29uc29sZS5lcnJvcihlKTtcblx0fVxuXG5cdC8vIEluaXRpYWxpc2UgdGhlIGJpbmRpbmdzIG1hcFxuXHR3aW5kb3cuZ28gPSB3aW5kb3cuZ28gfHwge307XG5cblx0Ly8gSXRlcmF0ZSBwYWNrYWdlIG5hbWVzXG5cdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwKS5mb3JFYWNoKChwYWNrYWdlTmFtZSkgPT4ge1xuXG5cdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XG5cdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gfHwge307XG5cblx0XHQvLyBJdGVyYXRlIHN0cnVjdCBuYW1lc1xuXHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXSkuZm9yRWFjaCgoc3RydWN0TmFtZSkgPT4ge1xuXG5cdFx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3Rcblx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdIHx8IHt9O1xuXG5cdFx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0pLmZvckVhY2goKG1ldGhvZE5hbWUpID0+IHtcblxuXHRcdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdW21ldGhvZE5hbWVdID0gZnVuY3Rpb24gKCkge1xuXG5cdFx0XHRcdFx0Ly8gTm8gdGltZW91dCBieSBkZWZhdWx0XG5cdFx0XHRcdFx0bGV0IHRpbWVvdXQgPSAwO1xuXG5cdFx0XHRcdFx0Ly8gQWN0dWFsIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZnVuY3Rpb24gZHluYW1pYygpIHtcblx0XHRcdFx0XHRcdGNvbnN0IGFyZ3MgPSBbXS5zbGljZS5jYWxsKGFyZ3VtZW50cyk7XG5cdFx0XHRcdFx0XHRyZXR1cm4gQ2FsbChbcGFja2FnZU5hbWUsIHN0cnVjdE5hbWUsIG1ldGhvZE5hbWVdLmpvaW4oJy4nKSwgYXJncywgdGltZW91dCk7XG5cdFx0XHRcdFx0fVxuXG5cdFx0XHRcdFx0Ly8gQWxsb3cgc2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZHluYW1pYy5zZXRUaW1lb3V0ID0gZnVuY3Rpb24gKG5ld1RpbWVvdXQpIHtcblx0XHRcdFx0XHRcdHRpbWVvdXQgPSBuZXdUaW1lb3V0O1xuXHRcdFx0XHRcdH07XG5cblx0XHRcdFx0XHQvLyBBbGxvdyBnZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cblx0XHRcdFx0XHRkeW5hbWljLmdldFRpbWVvdXQgPSBmdW5jdGlvbiAoKSB7XG5cdFx0XHRcdFx0XHRyZXR1cm4gdGltZW91dDtcblx0XHRcdFx0XHR9O1xuXG5cdFx0XHRcdFx0cmV0dXJuIGR5bmFtaWM7XG5cdFx0XHRcdH0oKTtcblx0XHRcdH0pO1xuXHRcdH0pO1xuXHR9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZCgpIHtcbiAgICB3aW5kb3cubG9jYXRpb24ucmVsb2FkKCk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWRBcHAoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUicpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U3lzdGVtRGVmYXVsdFRoZW1lKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FTRFQnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldExpZ2h0VGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQUxUJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXREYXJrVGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQURUJyk7XG59XG5cbi8qKlxuICogUGxhY2UgdGhlIHdpbmRvdyBpbiB0aGUgY2VudGVyIG9mIHRoZSBzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dDZW50ZXIoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXYycpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIHdpbmRvdyB0aXRsZVxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZVxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0VGl0bGUodGl0bGUpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dUJyArIHRpdGxlKTtcbn1cblxuLyoqXG4gKiBNYWtlcyB0aGUgd2luZG93IGdvIGZ1bGxzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dGdWxsc2NyZWVuKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0YnKTtcbn1cblxuLyoqXG4gKiBSZXZlcnRzIHRoZSB3aW5kb3cgZnJvbSBmdWxsc2NyZWVuXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5mdWxsc2NyZWVuKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2YnKTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBpbiBmdWxsIHNjcmVlbiBtb2RlIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNGdWxsc2NyZWVuKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzRnVsbHNjcmVlblwiKTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXczonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG4vKipcbiAqIEdldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7dzogbnVtYmVyLCBoOiBudW1iZXJ9Pn0gVGhlIHNpemUgb2YgdGhlIHdpbmRvd1xuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRTaXplKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFNpemVcIik7XG59XG5cbi8qKlxuICogU2V0IHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXWjonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1pblNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcbn1cblxuXG5cbi8qKlxuICogU2V0IHRoZSB3aW5kb3cgQWx3YXlzT25Ub3Agb3Igbm90IG9uIHRvcFxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEFsd2F5c09uVG9wKGIpIHtcblxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FUUDonICsgKGIgPyAnMScgOiAnMCcpKTtcbn1cblxuXG5cblxuLyoqXG4gKiBTZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0geFxuICogQHBhcmFtIHtudW1iZXJ9IHlcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFBvc2l0aW9uKHgsIHkpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dwOicgKyB4ICsgJzonICsgeSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7eDogbnVtYmVyLCB5OiBudW1iZXJ9Pn0gVGhlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFBvc2l0aW9uKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFBvc1wiKTtcbn1cblxuLyoqXG4gKiBIaWRlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dIaWRlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0gnKTtcbn1cblxuLyoqXG4gKiBTaG93IHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTaG93KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1MnKTtcbn1cblxuLyoqXG4gKiBNYXhpbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXTScpO1xufVxuXG4vKipcbiAqIFRvZ2dsZSB0aGUgTWF4aW1pc2Ugb2YgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1RvZ2dsZU1heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3QnKTtcbn1cblxuLyoqXG4gKiBVbm1heGltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1UnKTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc01heGltaXNlZCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc01heGltaXNlZFwiKTtcbn1cblxuLyoqXG4gKiBNaW5pbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWluaW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXbScpO1xufVxuXG4vKipcbiAqIFVubWluaW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWluaW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdScpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1pbmltaXNlZCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWluaW1pc2VkKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWluaW1pc2VkXCIpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG5vcm1hbCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTm9ybWFsKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTm9ybWFsXCIpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gUiBSZWRcbiAqIEBwYXJhbSB7bnVtYmVyfSBHIEdyZWVuXG4gKiBAcGFyYW0ge251bWJlcn0gQiBCbHVlXG4gKiBAcGFyYW0ge251bWJlcn0gQSBBbHBoYVxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QmFja2dyb3VuZENvbG91cihSLCBHLCBCLCBBKSB7XG4gICAgbGV0IHJnYmEgPSBKU09OLnN0cmluZ2lmeSh7cjogUiB8fCAwLCBnOiBHIHx8IDAsIGI6IEIgfHwgMCwgYTogQSB8fCAyNTV9KTtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dyOicgKyByZ2JhKTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG5cbi8qKlxuICogR2V0cyB0aGUgYWxsIHNjcmVlbnMuIENhbGwgdGhpcyBhbmV3IGVhY2ggdGltZSB5b3Ugd2FudCB0byByZWZyZXNoIGRhdGEgZnJvbSB0aGUgdW5kZXJseWluZyB3aW5kb3dpbmcgc3lzdGVtLlxuICogQGV4cG9ydFxuICogQHR5cGVkZWYge2ltcG9ydCgnLi4vd3JhcHBlci9ydW50aW1lJykuU2NyZWVufSBTY3JlZW5cbiAqIEByZXR1cm4ge1Byb21pc2U8e1NjcmVlbltdfT59IFRoZSBzY3JlZW5zXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTY3JlZW5HZXRBbGwoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6U2NyZWVuR2V0QWxsXCIpO1xufVxuIiwgIi8qKlxuICogQGRlc2NyaXB0aW9uOiBVc2UgdGhlIHN5c3RlbSBkZWZhdWx0IGJyb3dzZXIgdG8gb3BlbiB0aGUgdXJsXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIFxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJyb3dzZXJPcGVuVVJMKHVybCkge1xuICB3aW5kb3cuV2FpbHNJbnZva2UoJ0JPOicgKyB1cmwpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cbi8qKlxuICogU2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dFxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2xpcGJvYXJkU2V0VGV4dCh0ZXh0KSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6Q2xpcGJvYXJkU2V0VGV4dFwiLCBbdGV4dF0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgdGV4dCBjb250ZW50IG9mIHRoZSBjbGlwYm9hcmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHtzdHJpbmd9Pn0gVGV4dCBjb250ZW50IG9mIHRoZSBjbGlwYm9hcmRcblxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2xpcGJvYXJkR2V0VGV4dCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRHZXRUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XG5pbXBvcnQge2V2ZW50TGlzdGVuZXJzLCBFdmVudHNFbWl0LCBFdmVudHNOb3RpZnksIEV2ZW50c09mZiwgRXZlbnRzT24sIEV2ZW50c09uY2UsIEV2ZW50c09uTXVsdGlwbGV9IGZyb20gJy4vZXZlbnRzJztcbmltcG9ydCB7Q2FsbCwgQ2FsbGJhY2ssIGNhbGxiYWNrc30gZnJvbSAnLi9jYWxscyc7XG5pbXBvcnQge1NldEJpbmRpbmdzfSBmcm9tIFwiLi9iaW5kaW5nc1wiO1xuaW1wb3J0ICogYXMgV2luZG93IGZyb20gXCIuL3dpbmRvd1wiO1xuaW1wb3J0ICogYXMgU2NyZWVuIGZyb20gXCIuL3NjcmVlblwiO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkXCI7XG5cblxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdRJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUycpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0gnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkVudmlyb25tZW50XCIpO1xufVxuXG4vLyBUaGUgSlMgcnVudGltZVxud2luZG93LnJ1bnRpbWUgPSB7XG4gICAgLi4uTG9nLFxuICAgIC4uLldpbmRvdyxcbiAgICAuLi5Ccm93c2VyLFxuICAgIC4uLlNjcmVlbixcbiAgICAuLi5DbGlwYm9hcmQsXG4gICAgRXZlbnRzT24sXG4gICAgRXZlbnRzT25jZSxcbiAgICBFdmVudHNPbk11bHRpcGxlLFxuICAgIEV2ZW50c0VtaXQsXG4gICAgRXZlbnRzT2ZmLFxuICAgIEVudmlyb25tZW50LFxuICAgIFNob3csXG4gICAgSGlkZSxcbiAgICBRdWl0XG59O1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy53YWlscyA9IHtcbiAgICBDYWxsYmFjayxcbiAgICBFdmVudHNOb3RpZnksXG4gICAgU2V0QmluZGluZ3MsXG4gICAgZXZlbnRMaXN0ZW5lcnMsXG4gICAgY2FsbGJhY2tzLFxuICAgIGZsYWdzOiB7XG4gICAgICAgIGRpc2FibGVTY3JvbGxiYXJEcmFnOiBmYWxzZSxcbiAgICAgICAgZGlzYWJsZVdhaWxzRGVmYXVsdENvbnRleHRNZW51OiBmYWxzZSxcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcbiAgICAgICAgZGVmYXVsdEN1cnNvcjogbnVsbCxcbiAgICAgICAgYm9yZGVyVGhpY2tuZXNzOiA2LFxuICAgICAgICBzaG91bGREcmFnOiBmYWxzZSxcbiAgICAgICAgZGVmZXJEcmFnVG9Nb3VzZU1vdmU6IGZhbHNlLFxuICAgICAgICBjc3NEcmFnUHJvcGVydHk6IFwiLS13YWlscy1kcmFnZ2FibGVcIixcbiAgICAgICAgY3NzRHJhZ1ZhbHVlOiBcImRyYWdcIixcbiAgICB9XG59O1xuXG4vLyBTZXQgdGhlIGJpbmRpbmdzXG5pZiAod2luZG93LndhaWxzYmluZGluZ3MpIHtcbiAgICB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3Mod2luZG93LndhaWxzYmluZGluZ3MpO1xuICAgIGRlbGV0ZSB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3M7XG59XG5cbi8vIFRoaXMgaXMgZXZhbHVhdGVkIGF0IGJ1aWxkIHRpbWUgaW4gcGFja2FnZS5qc29uXG4vLyBjb25zdCBkZXYgPSAwO1xuLy8gY29uc3QgcHJvZHVjdGlvbiA9IDE7XG5pZiAoRU5WID09PSAxKSB7XG4gICAgZGVsZXRlIHdpbmRvdy53YWlsc2JpbmRpbmdzO1xufVxuXG5sZXQgZHJhZ1Rlc3QgPSBmdW5jdGlvbiAoZSkge1xuICAgIHZhciB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZSh3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1Byb3BlcnR5KTtcbiAgICBpZiAodmFsKSB7XG4gICAgICB2YWwgPSB2YWwudHJpbSgpO1xuICAgIH1cbiAgICBcbiAgICBpZiAodmFsICE9PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1ZhbHVlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICBpZiAoZS5idXR0b25zICE9PSAxKSB7XG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBub3QgdGhlIHByaW1hcnkgYnV0dG9uIGhhcyBiZWVuIGNsaWNrZWQuXG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICBpZiAoZS5kZXRhaWwgIT09IDEpIHtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWdnaW5nIGlmIG1vcmUgdGhhbiBvbmNlIGhhcyBiZWVuIGNsaWNrZWQsIGUuZy4gd2hlbiBkb3VibGUgY2xpY2tpbmdcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIHJldHVybiB0cnVlO1xufTtcblxud2luZG93LndhaWxzLnNldENTU0RyYWdQcm9wZXJ0aWVzID0gZnVuY3Rpb24gKHByb3BlcnR5LCB2YWx1ZSkge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnUHJvcGVydHkgPSBwcm9wZXJ0eTtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1ZhbHVlID0gdmFsdWU7XG59XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xuXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKSB7XG4gICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcInJlc2l6ZTpcIiArIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKTtcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKGRyYWdUZXN0KGUpKSB7XG4gICAgICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZVNjcm9sbGJhckRyYWcpIHtcbiAgICAgICAgICAgIC8vIFRoaXMgY2hlY2tzIGZvciBjbGlja3Mgb24gdGhlIHNjcm9sbCBiYXJcbiAgICAgICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kZWZlckRyYWdUb01vdXNlTW92ZSkge1xuICAgICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IHRydWU7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBlLnByZXZlbnREZWZhdWx0KClcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufSk7XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgKCkgPT4ge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG59KTtcblxuZnVuY3Rpb24gc2V0UmVzaXplKGN1cnNvcikge1xuICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBjdXJzb3IgfHwgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3I7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBmdW5jdGlvbiAoZSkge1xuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZykge1xuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgICAgICBpZihtb3VzZVByZXNzZWQgPD0gMCkge1xuICAgICAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlUmVzaXplKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yID09IG51bGwpIHtcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yO1xuICAgIH1cbiAgICBpZiAod2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzICYmIHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MpIHtcbiAgICAgICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IFwic2UtcmVzaXplXCI7XG4gICAgfVxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG4gICAgbGV0IHRvcEJvcmRlciA9IGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG5cbiAgICAvLyBJZiB3ZSBhcmVuJ3Qgb24gYW4gZWRnZSwgYnV0IHdlcmUsIHJlc2V0IHRoZSBjdXJzb3IgdG8gZGVmYXVsdFxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfSBlbHNlIGlmIChyaWdodEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgdG9wQm9yZGVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyICYmIHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcblxufSk7XG5cbi8vIFNldHVwIGNvbnRleHQgbWVudSBob29rXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBmdW5jdGlvbiAoZSkge1xuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZVdhaWxzRGVmYXVsdENvbnRleHRNZW51KSB7XG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG59KTtcblxud2luZG93LldhaWxzSW52b2tlKFwicnVudGltZTpyZWFkeVwiKTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFB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
