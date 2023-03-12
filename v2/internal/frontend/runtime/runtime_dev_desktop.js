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
      for (let count = eventListeners[eventName].length - 1; count >= 0; count -= 1) {
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
      window.wails.flags.shouldDrag = false;
      let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
      if (mousePressed > 0) {
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbi8qKlxuICogU2VuZHMgYSBsb2cgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB3aXRoIHRoZSBnaXZlbiBsZXZlbCArIG1lc3NhZ2VcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gbGV2ZWxcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmZ1bmN0aW9uIHNlbmRMb2dNZXNzYWdlKGxldmVsLCBtZXNzYWdlKSB7XG5cblx0Ly8gTG9nIE1lc3NhZ2UgZm9ybWF0OlxuXHQvLyBsW3R5cGVdW21lc3NhZ2VdXG5cdHdpbmRvdy5XYWlsc0ludm9rZSgnTCcgKyBsZXZlbCArIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gdHJhY2UgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1RyYWNlKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1QnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dQcmludChtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdQJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBkZWJ1ZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRGVidWcobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gaW5mbyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nSW5mbyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdJJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiB3YXJuaW5nIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dXYXJuaW5nKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1cnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGVycm9yIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dFcnJvcihtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdFJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBmYXRhbCBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRmF0YWwobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRicsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIExvZyBsZXZlbCB0byB0aGUgZ2l2ZW4gbG9nIGxldmVsXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IGxvZ2xldmVsXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRMb2dMZXZlbChsb2dsZXZlbCkge1xuXHRzZW5kTG9nTWVzc2FnZSgnUycsIGxvZ2xldmVsKTtcbn1cblxuLy8gTG9nIGxldmVsc1xuZXhwb3J0IGNvbnN0IExvZ0xldmVsID0ge1xuXHRUUkFDRTogMSxcblx0REVCVUc6IDIsXG5cdElORk86IDMsXG5cdFdBUk5JTkc6IDQsXG5cdEVSUk9SOiA1LFxufTtcbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuLy8gRGVmaW5lcyBhIHNpbmdsZSBsaXN0ZW5lciB3aXRoIGEgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdG8gY2FsbGJhY2tcblxuLyoqXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcbiAqXG4gKiBAY2xhc3MgTGlzdGVuZXJcbiAqL1xuY2xhc3MgTGlzdGVuZXIge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICAgICAqIEBtZW1iZXJvZiBMaXN0ZW5lclxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xuICAgICAgICAgICAgY2FsbGJhY2suYXBwbHkobnVsbCwgZGF0YSk7XG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IHt9O1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuZnVuY3Rpb24gbm90aWZ5TGlzdGVuZXJzKGV2ZW50RGF0YSkge1xuXG4gICAgLy8gR2V0IHRoZSBldmVudCBuYW1lXG4gICAgbGV0IGV2ZW50TmFtZSA9IGV2ZW50RGF0YS5uYW1lO1xuXG4gICAgLy8gQ2hlY2sgaWYgd2UgaGF2ZSBhbnkgbGlzdGVuZXJzIGZvciB0aGlzIGV2ZW50XG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0pIHtcblxuICAgICAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lciBpbmRleGVzIHRvIGRlc3Ryb3lcbiAgICAgICAgY29uc3QgbmV3RXZlbnRMaXN0ZW5lckxpc3QgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnNsaWNlKCk7XG5cbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmxlbmd0aCAtIDE7IGNvdW50ID49IDA7IGNvdW50IC09IDEpIHtcblxuICAgICAgICAgICAgLy8gR2V0IG5leHQgbGlzdGVuZXJcbiAgICAgICAgICAgIGNvbnN0IGxpc3RlbmVyID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXVtjb3VudF07XG5cbiAgICAgICAgICAgIGxldCBkYXRhID0gZXZlbnREYXRhLmRhdGE7XG5cbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xuICAgICAgICAgICAgY29uc3QgZGVzdHJveSA9IGxpc3RlbmVyLkNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgaWYgKGRlc3Ryb3kpIHtcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxuICAgICAgICAgICAgICAgIG5ld0V2ZW50TGlzdGVuZXJMaXN0LnNwbGljZShjb3VudCwgMSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cblxuICAgICAgICAvLyBVcGRhdGUgY2FsbGJhY2tzIHdpdGggbmV3IGxpc3Qgb2YgbGlzdGVuZXJzXG4gICAgICAgIGlmIChuZXdFdmVudExpc3RlbmVyTGlzdC5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gbmV3RXZlbnRMaXN0ZW5lckxpc3Q7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogTm90aWZ5IGluZm9ybXMgZnJvbnRlbmQgbGlzdGVuZXJzIHRoYXQgYW4gZXZlbnQgd2FzIGVtaXR0ZWQgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBub3RpZnlNZXNzYWdlIC0gZW5jb2RlZCBub3RpZmljYXRpb24gbWVzc2FnZVxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNOb3RpZnkobm90aWZ5TWVzc2FnZSkge1xuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXG4gICAgbGV0IG1lc3NhZ2U7XG4gICAgdHJ5IHtcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zdCBlcnJvciA9ICdJbnZhbGlkIEpTT04gcGFzc2VkIHRvIE5vdGlmeTogJyArIG5vdGlmeU1lc3NhZ2U7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XG4gICAgfVxuICAgIG5vdGlmeUxpc3RlbmVycyhtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNFbWl0KGV2ZW50TmFtZSkge1xuXG4gICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgbmFtZTogZXZlbnROYW1lLFxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxuICAgIH07XG5cbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXG4gICAgbm90aWZ5TGlzdGVuZXJzKHBheWxvYWQpO1xuXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRUUnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xufVxuXG5mdW5jdGlvbiByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpIHtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJzXG4gICAgZGVsZXRlIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV07XG5cbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFWCcgKyBldmVudE5hbWUpO1xufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcmVzIGNhbiBiZSB1bnJlZ2lzdGVyZWQgdmlhIGBhZGRpdGlvbmFsRXZlbnROYW1lc2BcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcblxuICAgIGlmIChhZGRpdGlvbmFsRXZlbnROYW1lcy5sZW5ndGggPiAwKSB7XG4gICAgICAgIGFkZGl0aW9uYWxFdmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IHtcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcbiAgICAgICAgfSlcbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGFsbCBldmVudCBsaXN0ZW5lcnMgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cbiAqL1xuIGV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmZBbGwoKSB7XG4gICAgY29uc3QgZXZlbnROYW1lcyA9IE9iamVjdC5rZXlzKGV2ZW50TGlzdGVuZXJzKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSAhPT0gZXZlbnROYW1lcy5sZW5ndGg7IGkrKykge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWVzW2ldKTtcbiAgICB9XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBFdmVudHNPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbiBmdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcbiAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG5cbiAgICAvLyBDbGVhbiB1cCBpZiB0aGVyZSBhcmUgbm8gZXZlbnQgbGlzdGVuZXJzIGxlZnRcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGggPT09IDApIHtcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcbiAgICB9XG59XG4iLCAiLypcbiBfICAgICAgIF9fICAgICAgXyBfX1xufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbmV4cG9ydCBjb25zdCBjYWxsYmFja3MgPSB7fTtcblxuLyoqXG4gKiBSZXR1cm5zIGEgbnVtYmVyIGZyb20gdGhlIG5hdGl2ZSBicm93c2VyIHJhbmRvbSBmdW5jdGlvblxuICpcbiAqIEByZXR1cm5zIG51bWJlclxuICovXG5mdW5jdGlvbiBjcnlwdG9SYW5kb20oKSB7XG5cdHZhciBhcnJheSA9IG5ldyBVaW50MzJBcnJheSgxKTtcblx0cmV0dXJuIHdpbmRvdy5jcnlwdG8uZ2V0UmFuZG9tVmFsdWVzKGFycmF5KVswXTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgbnVtYmVyIHVzaW5nIGRhIG9sZC1za29vbCBNYXRoLlJhbmRvbVxuICogSSBsaWtlcyB0byBjYWxsIGl0IExPTFJhbmRvbVxuICpcbiAqIEByZXR1cm5zIG51bWJlclxuICovXG5mdW5jdGlvbiBiYXNpY1JhbmRvbSgpIHtcblx0cmV0dXJuIE1hdGgucmFuZG9tKCkgKiA5MDA3MTk5MjU0NzQwOTkxO1xufVxuXG4vLyBQaWNrIGEgcmFuZG9tIG51bWJlciBmdW5jdGlvbiBiYXNlZCBvbiBicm93c2VyIGNhcGFiaWxpdHlcbnZhciByYW5kb21GdW5jO1xuaWYgKHdpbmRvdy5jcnlwdG8pIHtcblx0cmFuZG9tRnVuYyA9IGNyeXB0b1JhbmRvbTtcbn0gZWxzZSB7XG5cdHJhbmRvbUZ1bmMgPSBiYXNpY1JhbmRvbTtcbn1cblxuXG4vKipcbiAqIENhbGwgc2VuZHMgYSBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHRvIGNhbGwgdGhlIGJpbmRpbmcgd2l0aCB0aGVcbiAqIGdpdmVuIGRhdGEuIEEgcHJvbWlzZSBpcyByZXR1cm5lZCBhbmQgd2lsbCBiZSBjb21wbGV0ZWQgd2hlbiB0aGVcbiAqIGJhY2tlbmQgcmVzcG9uZHMuIFRoaXMgd2lsbCBiZSByZXNvbHZlZCB3aGVuIHRoZSBjYWxsIHdhcyBzdWNjZXNzZnVsXG4gKiBvciByZWplY3RlZCBpZiBhbiBlcnJvciBpcyBwYXNzZWQgYmFjay5cbiAqIFRoZXJlIGlzIGEgdGltZW91dCBtZWNoYW5pc20uIElmIHRoZSBjYWxsIGRvZXNuJ3QgcmVzcG9uZCBpbiB0aGUgZ2l2ZW5cbiAqIHRpbWUgKGluIG1pbGxpc2Vjb25kcykgdGhlbiB0aGUgcHJvbWlzZSBpcyByZWplY3RlZC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZVxuICogQHBhcmFtIHthbnk9fSBhcmdzXG4gKiBAcGFyYW0ge251bWJlcj19IHRpbWVvdXRcbiAqIEByZXR1cm5zXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG5hbWUsIGFyZ3MsIHRpbWVvdXQpIHtcblxuXHQvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcblx0aWYgKHRpbWVvdXQgPT0gbnVsbCkge1xuXHRcdHRpbWVvdXQgPSAwO1xuXHR9XG5cblx0Ly8gQ3JlYXRlIGEgcHJvbWlzZVxuXHRyZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xuXG5cdFx0Ly8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcblx0XHR2YXIgY2FsbGJhY2tJRDtcblx0XHRkbyB7XG5cdFx0XHRjYWxsYmFja0lEID0gbmFtZSArICctJyArIHJhbmRvbUZ1bmMoKTtcblx0XHR9IHdoaWxlIChjYWxsYmFja3NbY2FsbGJhY2tJRF0pO1xuXG5cdFx0dmFyIHRpbWVvdXRIYW5kbGU7XG5cdFx0Ly8gU2V0IHRpbWVvdXRcblx0XHRpZiAodGltZW91dCA+IDApIHtcblx0XHRcdHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcblx0XHRcdFx0cmVqZWN0KEVycm9yKCdDYWxsIHRvICcgKyBuYW1lICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XG5cdFx0XHR9LCB0aW1lb3V0KTtcblx0XHR9XG5cblx0XHQvLyBTdG9yZSBjYWxsYmFja1xuXHRcdGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcblx0XHRcdHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXG5cdFx0XHRyZWplY3Q6IHJlamVjdCxcblx0XHRcdHJlc29sdmU6IHJlc29sdmVcblx0XHR9O1xuXG5cdFx0dHJ5IHtcblx0XHRcdGNvbnN0IHBheWxvYWQgPSB7XG5cdFx0XHRcdG5hbWUsXG5cdFx0XHRcdGFyZ3MsXG5cdFx0XHRcdGNhbGxiYWNrSUQsXG5cdFx0XHR9O1xuXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XG4gICAgICAgIH1cbiAgICB9KTtcbn1cblxud2luZG93Lk9iZnVzY2F0ZWRDYWxsID0gKGlkLCBhcmdzLCB0aW1lb3V0KSA9PiB7XG5cbiAgICAvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcbiAgICBpZiAodGltZW91dCA9PSBudWxsKSB7XG4gICAgICAgIHRpbWVvdXQgPSAwO1xuICAgIH1cblxuICAgIC8vIENyZWF0ZSBhIHByb21pc2VcbiAgICByZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xuXG4gICAgICAgIC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXG4gICAgICAgIHZhciBjYWxsYmFja0lEO1xuICAgICAgICBkbyB7XG4gICAgICAgICAgICBjYWxsYmFja0lEID0gaWQgKyAnLScgKyByYW5kb21GdW5jKCk7XG4gICAgICAgIH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cbiAgICAgICAgdmFyIHRpbWVvdXRIYW5kbGU7XG4gICAgICAgIC8vIFNldCB0aW1lb3V0XG4gICAgICAgIGlmICh0aW1lb3V0ID4gMCkge1xuICAgICAgICAgICAgdGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcignQ2FsbCB0byBtZXRob2QgJyArIGlkICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XG4gICAgICAgICAgICB9LCB0aW1lb3V0KTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFN0b3JlIGNhbGxiYWNrXG4gICAgICAgIGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXG4gICAgICAgICAgICByZWplY3Q6IHJlamVjdCxcbiAgICAgICAgICAgIHJlc29sdmU6IHJlc29sdmVcbiAgICAgICAgfTtcblxuICAgICAgICB0cnkge1xuICAgICAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0aWQsXG5cdFx0XHRcdGFyZ3MsXG5cdFx0XHRcdGNhbGxiYWNrSUQsXG5cdFx0XHR9O1xuXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ2MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XG4gICAgICAgIH1cbiAgICB9KTtcbn07XG5cblxuLyoqXG4gKiBDYWxsZWQgYnkgdGhlIGJhY2tlbmQgdG8gcmV0dXJuIGRhdGEgdG8gYSBwcmV2aW91c2x5IGNhbGxlZFxuICogYmluZGluZyBpbnZvY2F0aW9uXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGluY29taW5nTWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbGJhY2soaW5jb21pbmdNZXNzYWdlKSB7XG5cdC8vIFBhcnNlIHRoZSBtZXNzYWdlXG5cdGxldCBtZXNzYWdlO1xuXHR0cnkge1xuXHRcdG1lc3NhZ2UgPSBKU09OLnBhcnNlKGluY29taW5nTWVzc2FnZSk7XG5cdH0gY2F0Y2ggKGUpIHtcblx0XHRjb25zdCBlcnJvciA9IGBJbnZhbGlkIEpTT04gcGFzc2VkIHRvIGNhbGxiYWNrOiAke2UubWVzc2FnZX0uIE1lc3NhZ2U6ICR7aW5jb21pbmdNZXNzYWdlfWA7XG5cdFx0cnVudGltZS5Mb2dEZWJ1ZyhlcnJvcik7XG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcblx0fVxuXHRsZXQgY2FsbGJhY2tJRCA9IG1lc3NhZ2UuY2FsbGJhY2tpZDtcblx0bGV0IGNhbGxiYWNrRGF0YSA9IGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcblx0aWYgKCFjYWxsYmFja0RhdGEpIHtcblx0XHRjb25zdCBlcnJvciA9IGBDYWxsYmFjayAnJHtjYWxsYmFja0lEfScgbm90IHJlZ2lzdGVyZWQhISFgO1xuXHRcdGNvbnNvbGUuZXJyb3IoZXJyb3IpOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcblx0fVxuXHRjbGVhclRpbWVvdXQoY2FsbGJhY2tEYXRhLnRpbWVvdXRIYW5kbGUpO1xuXG5cdGRlbGV0ZSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XG5cblx0aWYgKG1lc3NhZ2UuZXJyb3IpIHtcblx0XHRjYWxsYmFja0RhdGEucmVqZWN0KG1lc3NhZ2UuZXJyb3IpO1xuXHR9IGVsc2Uge1xuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX18gICAgXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gICkgXG58X18vfF9fL1xcX18sXy9fL18vX19fXy8gIFxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuaW1wb3J0IHtDYWxsfSBmcm9tICcuL2NhbGxzJztcblxuLy8gVGhpcyBpcyB3aGVyZSB3ZSBiaW5kIGdvIG1ldGhvZCB3cmFwcGVyc1xud2luZG93LmdvID0ge307XG5cbmV4cG9ydCBmdW5jdGlvbiBTZXRCaW5kaW5ncyhiaW5kaW5nc01hcCkge1xuXHR0cnkge1xuXHRcdGJpbmRpbmdzTWFwID0gSlNPTi5wYXJzZShiaW5kaW5nc01hcCk7XG5cdH0gY2F0Y2ggKGUpIHtcblx0XHRjb25zb2xlLmVycm9yKGUpO1xuXHR9XG5cblx0Ly8gSW5pdGlhbGlzZSB0aGUgYmluZGluZ3MgbWFwXG5cdHdpbmRvdy5nbyA9IHdpbmRvdy5nbyB8fCB7fTtcblxuXHQvLyBJdGVyYXRlIHBhY2thZ2UgbmFtZXNcblx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXApLmZvckVhY2goKHBhY2thZ2VOYW1lKSA9PiB7XG5cblx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3Rcblx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXSB8fCB7fTtcblxuXHRcdC8vIEl0ZXJhdGUgc3RydWN0IG5hbWVzXG5cdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdKS5mb3JFYWNoKChzdHJ1Y3ROYW1lKSA9PiB7XG5cblx0XHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxuXHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gfHwge307XG5cblx0XHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSkuZm9yRWFjaCgobWV0aG9kTmFtZSkgPT4ge1xuXG5cdFx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV1bbWV0aG9kTmFtZV0gPSBmdW5jdGlvbiAoKSB7XG5cblx0XHRcdFx0XHQvLyBObyB0aW1lb3V0IGJ5IGRlZmF1bHRcblx0XHRcdFx0XHRsZXQgdGltZW91dCA9IDA7XG5cblx0XHRcdFx0XHQvLyBBY3R1YWwgZnVuY3Rpb25cblx0XHRcdFx0XHRmdW5jdGlvbiBkeW5hbWljKCkge1xuXHRcdFx0XHRcdFx0Y29uc3QgYXJncyA9IFtdLnNsaWNlLmNhbGwoYXJndW1lbnRzKTtcblx0XHRcdFx0XHRcdHJldHVybiBDYWxsKFtwYWNrYWdlTmFtZSwgc3RydWN0TmFtZSwgbWV0aG9kTmFtZV0uam9pbignLicpLCBhcmdzLCB0aW1lb3V0KTtcblx0XHRcdFx0XHR9XG5cblx0XHRcdFx0XHQvLyBBbGxvdyBzZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cblx0XHRcdFx0XHRkeW5hbWljLnNldFRpbWVvdXQgPSBmdW5jdGlvbiAobmV3VGltZW91dCkge1xuXHRcdFx0XHRcdFx0dGltZW91dCA9IG5ld1RpbWVvdXQ7XG5cdFx0XHRcdFx0fTtcblxuXHRcdFx0XHRcdC8vIEFsbG93IGdldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxuXHRcdFx0XHRcdGR5bmFtaWMuZ2V0VGltZW91dCA9IGZ1bmN0aW9uICgpIHtcblx0XHRcdFx0XHRcdHJldHVybiB0aW1lb3V0O1xuXHRcdFx0XHRcdH07XG5cblx0XHRcdFx0XHRyZXR1cm4gZHluYW1pYztcblx0XHRcdFx0fSgpO1xuXHRcdFx0fSk7XG5cdFx0fSk7XG5cdH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93UmVsb2FkKCkge1xuICAgIHdpbmRvdy5sb2NhdGlvbi5yZWxvYWQoKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZEFwcCgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dSJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTeXN0ZW1EZWZhdWx0VGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVNEVCcpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TGlnaHRUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBTFQnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldERhcmtUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBRFQnKTtcbn1cblxuLyoqXG4gKiBQbGFjZSB0aGUgd2luZG93IGluIHRoZSBjZW50ZXIgb2YgdGhlIHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0NlbnRlcigpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1djJyk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgd2luZG93IHRpdGxlXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHRpdGxlXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRUaXRsZSh0aXRsZSkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1QnICsgdGl0bGUpO1xufVxuXG4vKipcbiAqIE1ha2VzIHRoZSB3aW5kb3cgZ28gZnVsbHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0Z1bGxzY3JlZW4oKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXRicpO1xufVxuXG4vKipcbiAqIFJldmVydHMgdGhlIHdpbmRvdyBmcm9tIGZ1bGxzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbmZ1bGxzY3JlZW4oKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXZicpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIGluIGZ1bGwgc2NyZWVuIG1vZGUgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc0Z1bGxzY3JlZW4oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNGdWxsc2NyZWVuXCIpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dzOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHt3OiBudW1iZXIsIGg6IG51bWJlcn0+fSBUaGUgc2l6ZSBvZiB0aGUgd2luZG93XG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFNpemUoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0U2l6ZVwiKTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNYXhTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1daOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cbi8qKlxuICogU2V0IHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXejonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG5cblxuLyoqXG4gKiBTZXQgdGhlIHdpbmRvdyBBbHdheXNPblRvcCBvciBub3Qgb24gdG9wXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QWx3YXlzT25Ub3AoYikge1xuXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVRQOicgKyAoYiA/ICcxJyA6ICcwJykpO1xufVxuXG5cblxuXG4vKipcbiAqIFNldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gKiBAcGFyYW0ge251bWJlcn0geVxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0UG9zaXRpb24oeCwgeSkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3A6JyArIHggKyAnOicgKyB5KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHt4OiBudW1iZXIsIHk6IG51bWJlcn0+fSBUaGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0UG9zaXRpb24oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0UG9zXCIpO1xufVxuXG4vKipcbiAqIEhpZGUgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0hpZGUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXSCcpO1xufVxuXG4vKipcbiAqIFNob3cgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1Nob3coKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUycpO1xufVxuXG4vKipcbiAqIE1heGltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNYXhpbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dNJyk7XG59XG5cbi8qKlxuICogVG9nZ2xlIHRoZSBNYXhpbWlzZSBvZiB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VG9nZ2xlTWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdCcpO1xufVxuXG4vKipcbiAqIFVubWF4aW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVScpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1heGltaXNlZCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWF4aW1pc2VkKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWF4aW1pc2VkXCIpO1xufVxuXG4vKipcbiAqIE1pbmltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNaW5pbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dtJyk7XG59XG5cbi8qKlxuICogVW5taW5pbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5taW5pbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d1Jyk7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNNaW5pbWlzZWQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNaW5pbWlzZWRcIik7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbm9ybWFsIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNOb3JtYWwoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNOb3JtYWxcIik7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSBSIFJlZFxuICogQHBhcmFtIHtudW1iZXJ9IEcgR3JlZW5cbiAqIEBwYXJhbSB7bnVtYmVyfSBCIEJsdWVcbiAqIEBwYXJhbSB7bnVtYmVyfSBBIEFscGhhXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRCYWNrZ3JvdW5kQ29sb3VyKFIsIEcsIEIsIEEpIHtcbiAgICBsZXQgcmdiYSA9IEpTT04uc3RyaW5naWZ5KHtyOiBSIHx8IDAsIGc6IEcgfHwgMCwgYjogQiB8fCAwLCBhOiBBIHx8IDI1NX0pO1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cblxuLyoqXG4gKiBHZXRzIHRoZSBhbGwgc2NyZWVucy4gQ2FsbCB0aGlzIGFuZXcgZWFjaCB0aW1lIHlvdSB3YW50IHRvIHJlZnJlc2ggZGF0YSBmcm9tIHRoZSB1bmRlcmx5aW5nIHdpbmRvd2luZyBzeXN0ZW0uXG4gKiBAZXhwb3J0XG4gKiBAdHlwZWRlZiB7aW1wb3J0KCcuLi93cmFwcGVyL3J1bnRpbWUnKS5TY3JlZW59IFNjcmVlblxuICogQHJldHVybiB7UHJvbWlzZTx7U2NyZWVuW119Pn0gVGhlIHNjcmVlbnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNjcmVlbkdldEFsbCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpTY3JlZW5HZXRBbGxcIik7XG59XG4iLCAiLyoqXG4gKiBAZGVzY3JpcHRpb246IFVzZSB0aGUgc3lzdGVtIGRlZmF1bHQgYnJvd3NlciB0byBvcGVuIHRoZSB1cmxcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gQnJvd3Nlck9wZW5VUkwodXJsKSB7XG4gIHdpbmRvdy5XYWlsc0ludm9rZSgnQk86JyArIHVybCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuLyoqXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRTZXRUZXh0XCIsIFt0ZXh0XSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSB0ZXh0IGNvbnRlbnQgb2YgdGhlIGNsaXBib2FyZFxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8e3N0cmluZ30+fSBUZXh0IGNvbnRlbnQgb2YgdGhlIGNsaXBib2FyZFxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRHZXRUZXh0KCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsaXBib2FyZEdldFRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcbmltcG9ydCB7ZXZlbnRMaXN0ZW5lcnMsIEV2ZW50c0VtaXQsIEV2ZW50c05vdGlmeSwgRXZlbnRzT2ZmLCBFdmVudHNPbiwgRXZlbnRzT25jZSwgRXZlbnRzT25NdWx0aXBsZX0gZnJvbSAnLi9ldmVudHMnO1xuaW1wb3J0IHtDYWxsLCBDYWxsYmFjaywgY2FsbGJhY2tzfSBmcm9tICcuL2NhbGxzJztcbmltcG9ydCB7U2V0QmluZGluZ3N9IGZyb20gXCIuL2JpbmRpbmdzXCI7XG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQgKiBhcyBTY3JlZW4gZnJvbSBcIi4vc2NyZWVuXCI7XG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXJcIjtcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmRcIjtcblxuXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1EnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdTJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnSCcpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6RW52aXJvbm1lbnRcIik7XG59XG5cbi8vIFRoZSBKUyBydW50aW1lXG53aW5kb3cucnVudGltZSA9IHtcbiAgICAuLi5Mb2csXG4gICAgLi4uV2luZG93LFxuICAgIC4uLkJyb3dzZXIsXG4gICAgLi4uU2NyZWVuLFxuICAgIC4uLkNsaXBib2FyZCxcbiAgICBFdmVudHNPbixcbiAgICBFdmVudHNPbmNlLFxuICAgIEV2ZW50c09uTXVsdGlwbGUsXG4gICAgRXZlbnRzRW1pdCxcbiAgICBFdmVudHNPZmYsXG4gICAgRW52aXJvbm1lbnQsXG4gICAgU2hvdyxcbiAgICBIaWRlLFxuICAgIFF1aXRcbn07XG5cbi8vIEludGVybmFsIHdhaWxzIGVuZHBvaW50c1xud2luZG93LndhaWxzID0ge1xuICAgIENhbGxiYWNrLFxuICAgIEV2ZW50c05vdGlmeSxcbiAgICBTZXRCaW5kaW5ncyxcbiAgICBldmVudExpc3RlbmVycyxcbiAgICBjYWxsYmFja3MsXG4gICAgZmxhZ3M6IHtcbiAgICAgICAgZGlzYWJsZVNjcm9sbGJhckRyYWc6IGZhbHNlLFxuICAgICAgICBkaXNhYmxlV2FpbHNEZWZhdWx0Q29udGV4dE1lbnU6IGZhbHNlLFxuICAgICAgICBlbmFibGVSZXNpemU6IGZhbHNlLFxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxuICAgICAgICBib3JkZXJUaGlja25lc3M6IDYsXG4gICAgICAgIHNob3VsZERyYWc6IGZhbHNlLFxuICAgICAgICBkZWZlckRyYWdUb01vdXNlTW92ZTogZmFsc2UsXG4gICAgICAgIGNzc0RyYWdQcm9wZXJ0eTogXCItLXdhaWxzLWRyYWdnYWJsZVwiLFxuICAgICAgICBjc3NEcmFnVmFsdWU6IFwiZHJhZ1wiLFxuICAgIH1cbn07XG5cbi8vIFNldCB0aGUgYmluZGluZ3NcbmlmICh3aW5kb3cud2FpbHNiaW5kaW5ncykge1xuICAgIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncyh3aW5kb3cud2FpbHNiaW5kaW5ncyk7XG4gICAgZGVsZXRlIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncztcbn1cblxuLy8gVGhpcyBpcyBldmFsdWF0ZWQgYXQgYnVpbGQgdGltZSBpbiBwYWNrYWdlLmpzb25cbi8vIGNvbnN0IGRldiA9IDA7XG4vLyBjb25zdCBwcm9kdWN0aW9uID0gMTtcbmlmIChFTlYgPT09IDEpIHtcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XG59XG5cbmxldCBkcmFnVGVzdCA9IGZ1bmN0aW9uIChlKSB7XG4gICAgdmFyIHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnUHJvcGVydHkpO1xuICAgIGlmICh2YWwpIHtcbiAgICAgIHZhbCA9IHZhbC50cmltKCk7XG4gICAgfVxuICAgIFxuICAgIGlmICh2YWwgIT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIGlmIChlLmJ1dHRvbnMgIT09IDEpIHtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWdnaW5nIGlmIG5vdCB0aGUgcHJpbWFyeSBidXR0b24gaGFzIGJlZW4gY2xpY2tlZC5cbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIGlmIChlLmRldGFpbCAhPT0gMSkge1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZ2dpbmcgaWYgbW9yZSB0aGFuIG9uY2UgaGFzIGJlZW4gY2xpY2tlZCwgZS5nLiB3aGVuIGRvdWJsZSBjbGlja2luZ1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgcmV0dXJuIHRydWU7XG59O1xuXG53aW5kb3cud2FpbHMuc2V0Q1NTRHJhZ1Byb3BlcnRpZXMgPSBmdW5jdGlvbiAocHJvcGVydHksIHZhbHVlKSB7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdQcm9wZXJ0eSA9IHByb3BlcnR5O1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUgPSB2YWx1ZTtcbn1cblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIChlKSA9PiB7XG5cbiAgICAvLyBDaGVjayBmb3IgcmVzaXppbmdcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwicmVzaXplOlwiICsgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpO1xuICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoZHJhZ1Rlc3QoZSkpIHtcbiAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlU2Nyb2xsYmFyRHJhZykge1xuICAgICAgICAgICAgLy8gVGhpcyBjaGVja3MgZm9yIGNsaWNrcyBvbiB0aGUgc2Nyb2xsIGJhclxuICAgICAgICAgICAgaWYgKGUub2Zmc2V0WCA+IGUudGFyZ2V0LmNsaWVudFdpZHRoIHx8IGUub2Zmc2V0WSA+IGUudGFyZ2V0LmNsaWVudEhlaWdodCkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRlZmVyRHJhZ1RvTW91c2VNb3ZlKSB7XG4gICAgICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IHRydWU7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBlLnByZXZlbnREZWZhdWx0KClcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufSk7XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgKCkgPT4ge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG59KTtcblxuZnVuY3Rpb24gc2V0UmVzaXplKGN1cnNvcikge1xuICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBjdXJzb3IgfHwgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3I7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBmdW5jdGlvbiAoZSkge1xuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZykge1xuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgICAgICBpZiAobW91c2VQcmVzc2VkID4gMCkge1xuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVSZXNpemUpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPT0gbnVsbCkge1xuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3I7XG4gICAgfVxuICAgIGlmICh3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MgJiYgd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcykge1xuICAgICAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gXCJzZS1yZXNpemVcIjtcbiAgICB9XG4gICAgbGV0IHJpZ2h0Qm9yZGVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcblxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9IGVsc2UgaWYgKHJpZ2h0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiB0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIgJiYgcmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xuXG59KTtcblxuLy8gU2V0dXAgY29udGV4dCBtZW51IGhvb2tcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGZ1bmN0aW9uIChlKSB7XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlV2FpbHNEZWZhdWx0Q29udGV4dE1lbnUpIHtcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn0pO1xuXG53aW5kb3cuV2FpbHNJbnZva2UoXCJydW50aW1lOnJlYWR5XCIpOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFB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QUEsRUFDcEM7QUFFQSxTQUFPLGlCQUFpQixhQUFhLFNBQVUsR0FBRztBQUM5QyxRQUFJLE9BQU8sTUFBTSxNQUFNLFlBQVk7QUFDL0IsYUFBTyxNQUFNLE1BQU0sYUFBYTtBQUNoQyxVQUFJLGVBQWUsRUFBRSxZQUFZLFNBQVksRUFBRSxVQUFVLEVBQUU7QUFDM0QsVUFBSSxlQUFlLEdBQUc7QUFDbEIsZUFBTyxZQUFZLE1BQU07QUFDekI7QUFBQSxNQUNKO0FBQUEsSUFDSjtBQUNBLFFBQUksQ0FBQyxPQUFPLE1BQU0sTUFBTSxjQUFjO0FBQ2xDO0FBQUEsSUFDSjtBQUNBLFFBQUksT0FBTyxNQUFNLE1BQU0saUJBQWlCLE1BQU07QUFDMUMsYUFBTyxNQUFNLE1BQU0sZ0JBQWdCLFNBQVMsZ0JBQWdCLE1BQU07QUFBQSxJQUN0RTtBQUNBLFFBQUksT0FBTyxhQUFhLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTSxtQkFBbUIsT0FBTyxjQUFjLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTSxpQkFBaUI7QUFDM0ksZUFBUyxnQkFBZ0IsTUFBTSxTQUFTO0FBQUEsSUFDNUM7QUFDQSxRQUFJLGNBQWMsT0FBTyxhQUFhLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUNyRSxRQUFJLGFBQWEsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBQ2hELFFBQUksWUFBWSxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU07QUFDL0MsUUFBSSxlQUFlLE9BQU8sY0FBYyxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU07QUFHdkUsUUFBSSxDQUFDLGNBQWMsQ0FBQyxlQUFlLENBQUMsYUFBYSxDQUFDLGdCQUFnQixPQUFPLE1BQU0sTUFBTSxlQUFlLFFBQVc7QUFDM0csZ0JBQVU7QUFBQSxJQUNkLFdBQVcsZUFBZTtBQUFjLGdCQUFVLFdBQVc7QUFBQSxhQUNwRCxjQUFjO0FBQWMsZ0JBQVUsV0FBVztBQUFBLGFBQ2pELGNBQWM7QUFBVyxnQkFBVSxXQUFXO0FBQUEsYUFDOUMsYUFBYTtBQUFhLGdCQUFVLFdBQVc7QUFBQSxhQUMvQztBQUFZLGdCQUFVLFVBQVU7QUFBQSxhQUNoQztBQUFXLGdCQUFVLFVBQVU7QUFBQSxhQUMvQjtBQUFjLGdCQUFVLFVBQVU7QUFBQSxhQUNsQztBQUFhLGdCQUFVLFVBQVU7QUFBQSxFQUU5QyxDQUFDO0FBR0QsU0FBTyxpQkFBaUIsZUFBZSxTQUFVLEdBQUc7QUFDaEQsUUFBSSxPQUFPLE1BQU0sTUFBTSxnQ0FBZ0M7QUFDbkQsUUFBRSxlQUFlO0FBQUEsSUFDckI7QUFBQSxFQUNKLENBQUM7QUFFRCxTQUFPLFlBQVksZUFBZTsiLAogICJuYW1lcyI6IFsiZXZlbnROYW1lIl0KfQo=
