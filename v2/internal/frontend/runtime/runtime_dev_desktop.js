(() => {
  var __defProp = Object.defineProperty;
  var __markAsModule = (target) => __defProp(target, "__esModule", { value: true });
  var __export = (target, all) => {
    __markAsModule(target);
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
  function EventsOff(eventName) {
    delete eventListeners[eventName];
    window.WailsInvoke("EX" + eventName);
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
  function WindowMinimise() {
    window.WailsInvoke("Wm");
  }
  function WindowUnminimise() {
    window.WailsInvoke("Wu");
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
      dbClickInterval: 150
    }
  };
  window.wails.SetBindings(window.wailsbindings);
  delete window.wails.SetBindings;
  if (true) {
    delete window.wailsbindings;
  }
  window.addEventListener("mousedown", (e) => {
    if (window.wails.flags.resizeEdge) {
      window.WailsInvoke("resize:" + window.wails.flags.resizeEdge);
      e.preventDefault();
      return;
    }
    let currentElement = e.target;
    while (currentElement != null) {
      if (currentElement.hasAttribute("data-wails-no-drag")) {
        break;
      } else if (currentElement.hasAttribute("data-wails-drag")) {
        if (window.wails.flags.disableScrollbarDrag) {
          if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
            break;
          }
        }
        window.WailsInvoke("drag");
        break;
      }
      currentElement = currentElement.parentElement;
    }
  });
  function setResize(cursor) {
    document.body.style.cursor = cursor || window.wails.flags.defaultCursor;
    window.wails.flags.resizeEdge = cursor;
  }
  window.addEventListener("mousemove", function(e) {
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
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfX1xyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhIGxvZyBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHdpdGggdGhlIGdpdmVuIGxldmVsICsgbWVzc2FnZVxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbGV2ZWxcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmZ1bmN0aW9uIHNlbmRMb2dNZXNzYWdlKGxldmVsLCBtZXNzYWdlKSB7XHJcblxyXG5cdC8vIExvZyBNZXNzYWdlIGZvcm1hdDpcclxuXHQvLyBsW3R5cGVdW21lc3NhZ2VdXHJcblx0d2luZG93LldhaWxzSW52b2tlKCdMJyArIGxldmVsICsgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIHRyYWNlIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nVHJhY2UobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdUJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nUHJpbnQobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdQJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGRlYnVnIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRGVidWcobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdEJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGluZm8gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dJbmZvKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnSScsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiB3YXJuaW5nIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nV2FybmluZyhtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1cnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZXJyb3IgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dFcnJvcihtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0UnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZmF0YWwgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dGYXRhbChtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0YnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIExvZyBsZXZlbCB0byB0aGUgZ2l2ZW4gbG9nIGxldmVsXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IGxvZ2xldmVsXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0TG9nTGV2ZWwobG9nbGV2ZWwpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnUycsIGxvZ2xldmVsKTtcclxufVxyXG5cclxuLy8gTG9nIGxldmVsc1xyXG5leHBvcnQgY29uc3QgTG9nTGV2ZWwgPSB7XHJcblx0VFJBQ0U6IDEsXHJcblx0REVCVUc6IDIsXHJcblx0SU5GTzogMyxcclxuXHRXQVJOSU5HOiA0LFxyXG5cdEVSUk9SOiA1LFxyXG59O1xyXG4iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xyXG5cclxuLyoqXHJcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxyXG4gKlxyXG4gKiBAY2xhc3MgTGlzdGVuZXJcclxuICovXHJcbmNsYXNzIExpc3RlbmVyIHtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cclxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcclxuICAgICAgICBtYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XHJcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxyXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcclxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcclxuICAgICAgICAgICAgY2FsbGJhY2suYXBwbHkobnVsbCwgZGF0YSk7XHJcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcclxuICAgICAgICAgICAgaWYgKG1heENhbGxiYWNrcyA9PT0gLTEpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXHJcbiAgICAgICAgICAgIG1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgICAgICByZXR1cm4gbWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IHt9O1xyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xyXG4gICAgRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcclxufVxyXG5cclxuZnVuY3Rpb24gbm90aWZ5TGlzdGVuZXJzKGV2ZW50RGF0YSkge1xyXG5cclxuICAgIC8vIEdldCB0aGUgZXZlbnQgbmFtZVxyXG4gICAgbGV0IGV2ZW50TmFtZSA9IGV2ZW50RGF0YS5uYW1lO1xyXG5cclxuICAgIC8vIENoZWNrIGlmIHdlIGhhdmUgYW55IGxpc3RlbmVycyBmb3IgdGhpcyBldmVudFxyXG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0pIHtcclxuXHJcbiAgICAgICAgLy8gS2VlcCBhIGxpc3Qgb2YgbGlzdGVuZXIgaW5kZXhlcyB0byBkZXN0cm95XHJcbiAgICAgICAgY29uc3QgbmV3RXZlbnRMaXN0ZW5lckxpc3QgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnNsaWNlKCk7XHJcblxyXG4gICAgICAgIC8vIEl0ZXJhdGUgbGlzdGVuZXJzXHJcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSAwOyBjb3VudCA8IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoOyBjb3VudCArPSAxKSB7XHJcblxyXG4gICAgICAgICAgICAvLyBHZXQgbmV4dCBsaXN0ZW5lclxyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV1bY291bnRdO1xyXG5cclxuICAgICAgICAgICAgbGV0IGRhdGEgPSBldmVudERhdGEuZGF0YTtcclxuXHJcbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xyXG4gICAgICAgICAgICBjb25zdCBkZXN0cm95ID0gbGlzdGVuZXIuQ2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgICAgIGlmIChkZXN0cm95KSB7XHJcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxyXG4gICAgICAgICAgICAgICAgbmV3RXZlbnRMaXN0ZW5lckxpc3Quc3BsaWNlKGNvdW50LCAxKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gVXBkYXRlIGNhbGxiYWNrcyB3aXRoIG5ldyBsaXN0IG9mIGxpc3RlbmVyc1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBuZXdFdmVudExpc3RlbmVyTGlzdDtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcclxuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXHJcbiAgICBsZXQgbWVzc2FnZTtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcbiAgICB9XHJcbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcclxuXHJcbiAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcclxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxyXG4gICAgfTtcclxuXHJcbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXHJcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmYoZXZlbnROYW1lKSB7XHJcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJzXHJcbiAgICBkZWxldGUgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXTtcclxuXHJcbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VYJyArIGV2ZW50TmFtZSk7XHJcbn0iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cclxuICpcclxuICogQHJldHVybnMgbnVtYmVyXHJcbiAqL1xyXG5mdW5jdGlvbiBjcnlwdG9SYW5kb20oKSB7XHJcblx0dmFyIGFycmF5ID0gbmV3IFVpbnQzMkFycmF5KDEpO1xyXG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgbnVtYmVyIHVzaW5nIGRhIG9sZC1za29vbCBNYXRoLlJhbmRvbVxyXG4gKiBJIGxpa2VzIHRvIGNhbGwgaXQgTE9MUmFuZG9tXHJcbiAqXHJcbiAqIEByZXR1cm5zIG51bWJlclxyXG4gKi9cclxuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XHJcblx0cmV0dXJuIE1hdGgucmFuZG9tKCkgKiA5MDA3MTk5MjU0NzQwOTkxO1xyXG59XHJcblxyXG4vLyBQaWNrIGEgcmFuZG9tIG51bWJlciBmdW5jdGlvbiBiYXNlZCBvbiBicm93c2VyIGNhcGFiaWxpdHlcclxudmFyIHJhbmRvbUZ1bmM7XHJcbmlmICh3aW5kb3cuY3J5cHRvKSB7XHJcblx0cmFuZG9tRnVuYyA9IGNyeXB0b1JhbmRvbTtcclxufSBlbHNlIHtcclxuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogQ2FsbCBzZW5kcyBhIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgdG8gY2FsbCB0aGUgYmluZGluZyB3aXRoIHRoZVxyXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXHJcbiAqIGJhY2tlbmQgcmVzcG9uZHMuIFRoaXMgd2lsbCBiZSByZXNvbHZlZCB3aGVuIHRoZSBjYWxsIHdhcyBzdWNjZXNzZnVsXHJcbiAqIG9yIHJlamVjdGVkIGlmIGFuIGVycm9yIGlzIHBhc3NlZCBiYWNrLlxyXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXHJcbiAqIHRpbWUgKGluIG1pbGxpc2Vjb25kcykgdGhlbiB0aGUgcHJvbWlzZSBpcyByZWplY3RlZC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZVxyXG4gKiBAcGFyYW0ge2FueT19IGFyZ3NcclxuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XHJcbiAqIEByZXR1cm5zXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XHJcblxyXG5cdC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxyXG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcclxuXHRcdHRpbWVvdXQgPSAwO1xyXG5cdH1cclxuXHJcblx0Ly8gQ3JlYXRlIGEgcHJvbWlzZVxyXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XHJcblxyXG5cdFx0Ly8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcclxuXHRcdHZhciBjYWxsYmFja0lEO1xyXG5cdFx0ZG8ge1xyXG5cdFx0XHRjYWxsYmFja0lEID0gbmFtZSArICctJyArIHJhbmRvbUZ1bmMoKTtcclxuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XHJcblxyXG5cdFx0dmFyIHRpbWVvdXRIYW5kbGU7XHJcblx0XHQvLyBTZXQgdGltZW91dFxyXG5cdFx0aWYgKHRpbWVvdXQgPiAwKSB7XHJcblx0XHRcdHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcclxuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcclxuXHRcdFx0fSwgdGltZW91dCk7XHJcblx0XHR9XHJcblxyXG5cdFx0Ly8gU3RvcmUgY2FsbGJhY2tcclxuXHRcdGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcclxuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcclxuXHRcdFx0cmVqZWN0OiByZWplY3QsXHJcblx0XHRcdHJlc29sdmU6IHJlc29sdmVcclxuXHRcdH07XHJcblxyXG5cdFx0dHJ5IHtcclxuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcclxuXHRcdFx0XHRuYW1lLFxyXG5cdFx0XHRcdGFyZ3MsXHJcblx0XHRcdFx0Y2FsbGJhY2tJRCxcclxuXHRcdFx0fTtcclxuXHJcblx0XHRcdC8vIE1ha2UgdGhlIGNhbGxcclxuXHRcdFx0d2luZG93LldhaWxzSW52b2tlKCdDJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcclxuXHRcdH0gY2F0Y2ggKGUpIHtcclxuXHRcdFx0Ly8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXHJcblx0XHRcdGNvbnNvbGUuZXJyb3IoZSk7XHJcblx0XHR9XHJcblx0fSk7XHJcbn1cclxuXHJcblxyXG5cclxuLyoqXHJcbiAqIENhbGxlZCBieSB0aGUgYmFja2VuZCB0byByZXR1cm4gZGF0YSB0byBhIHByZXZpb3VzbHkgY2FsbGVkXHJcbiAqIGJpbmRpbmcgaW52b2NhdGlvblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBpbmNvbWluZ01lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsYmFjayhpbmNvbWluZ01lc3NhZ2UpIHtcclxuXHQvLyBQYXJzZSB0aGUgbWVzc2FnZVxyXG5cdGxldCBtZXNzYWdlO1xyXG5cdHRyeSB7XHJcblx0XHRtZXNzYWdlID0gSlNPTi5wYXJzZShpbmNvbWluZ01lc3NhZ2UpO1xyXG5cdH0gY2F0Y2ggKGUpIHtcclxuXHRcdGNvbnN0IGVycm9yID0gYEludmFsaWQgSlNPTiBwYXNzZWQgdG8gY2FsbGJhY2s6ICR7ZS5tZXNzYWdlfS4gTWVzc2FnZTogJHtpbmNvbWluZ01lc3NhZ2V9YDtcclxuXHRcdHJ1bnRpbWUuTG9nRGVidWcoZXJyb3IpO1xyXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcclxuXHR9XHJcblx0bGV0IGNhbGxiYWNrSUQgPSBtZXNzYWdlLmNhbGxiYWNraWQ7XHJcblx0bGV0IGNhbGxiYWNrRGF0YSA9IGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcclxuXHRpZiAoIWNhbGxiYWNrRGF0YSkge1xyXG5cdFx0Y29uc3QgZXJyb3IgPSBgQ2FsbGJhY2sgJyR7Y2FsbGJhY2tJRH0nIG5vdCByZWdpc3RlcmVkISEhYDtcclxuXHRcdGNvbnNvbGUuZXJyb3IoZXJyb3IpOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lXHJcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xyXG5cdH1cclxuXHRjbGVhclRpbWVvdXQoY2FsbGJhY2tEYXRhLnRpbWVvdXRIYW5kbGUpO1xyXG5cclxuXHRkZWxldGUgY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xyXG5cclxuXHRpZiAobWVzc2FnZS5lcnJvcikge1xyXG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChtZXNzYWdlLmVycm9yKTtcclxuXHR9IGVsc2Uge1xyXG5cdFx0Y2FsbGJhY2tEYXRhLnJlc29sdmUobWVzc2FnZS5yZXN1bHQpO1xyXG5cdH1cclxufVxyXG4iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fICAgIFxyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApIFxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy8gIFxyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xyXG5cclxuaW1wb3J0IHtDYWxsfSBmcm9tICcuL2NhbGxzJztcclxuXHJcbi8vIFRoaXMgaXMgd2hlcmUgd2UgYmluZCBnbyBtZXRob2Qgd3JhcHBlcnNcclxud2luZG93LmdvID0ge307XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gU2V0QmluZGluZ3MoYmluZGluZ3NNYXApIHtcclxuXHR0cnkge1xyXG5cdFx0YmluZGluZ3NNYXAgPSBKU09OLnBhcnNlKGJpbmRpbmdzTWFwKTtcclxuXHR9IGNhdGNoIChlKSB7XHJcblx0XHRjb25zb2xlLmVycm9yKGUpO1xyXG5cdH1cclxuXHJcblx0Ly8gSW5pdGlhbGlzZSB0aGUgYmluZGluZ3MgbWFwXHJcblx0d2luZG93LmdvID0gd2luZG93LmdvIHx8IHt9O1xyXG5cclxuXHQvLyBJdGVyYXRlIHBhY2thZ2UgbmFtZXNcclxuXHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcCkuZm9yRWFjaCgocGFja2FnZU5hbWUpID0+IHtcclxuXHJcblx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3RcclxuXHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdIHx8IHt9O1xyXG5cclxuXHRcdC8vIEl0ZXJhdGUgc3RydWN0IG5hbWVzXHJcblx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV0pLmZvckVhY2goKHN0cnVjdE5hbWUpID0+IHtcclxuXHJcblx0XHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxyXG5cdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSB8fCB7fTtcclxuXHJcblx0XHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSkuZm9yRWFjaCgobWV0aG9kTmFtZSkgPT4ge1xyXG5cclxuXHRcdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdW21ldGhvZE5hbWVdID0gZnVuY3Rpb24gKCkge1xyXG5cclxuXHRcdFx0XHRcdC8vIE5vIHRpbWVvdXQgYnkgZGVmYXVsdFxyXG5cdFx0XHRcdFx0bGV0IHRpbWVvdXQgPSAwO1xyXG5cclxuXHRcdFx0XHRcdC8vIEFjdHVhbCBmdW5jdGlvblxyXG5cdFx0XHRcdFx0ZnVuY3Rpb24gZHluYW1pYygpIHtcclxuXHRcdFx0XHRcdFx0Y29uc3QgYXJncyA9IFtdLnNsaWNlLmNhbGwoYXJndW1lbnRzKTtcclxuXHRcdFx0XHRcdFx0cmV0dXJuIENhbGwoW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXS5qb2luKCcuJyksIGFyZ3MsIHRpbWVvdXQpO1xyXG5cdFx0XHRcdFx0fVxyXG5cclxuXHRcdFx0XHRcdC8vIEFsbG93IHNldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxyXG5cdFx0XHRcdFx0ZHluYW1pYy5zZXRUaW1lb3V0ID0gZnVuY3Rpb24gKG5ld1RpbWVvdXQpIHtcclxuXHRcdFx0XHRcdFx0dGltZW91dCA9IG5ld1RpbWVvdXQ7XHJcblx0XHRcdFx0XHR9O1xyXG5cclxuXHRcdFx0XHRcdC8vIEFsbG93IGdldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxyXG5cdFx0XHRcdFx0ZHluYW1pYy5nZXRUaW1lb3V0ID0gZnVuY3Rpb24gKCkge1xyXG5cdFx0XHRcdFx0XHRyZXR1cm4gdGltZW91dDtcclxuXHRcdFx0XHRcdH07XHJcblxyXG5cdFx0XHRcdFx0cmV0dXJuIGR5bmFtaWM7XHJcblx0XHRcdFx0fSgpO1xyXG5cdFx0XHR9KTtcclxuXHRcdH0pO1xyXG5cdH0pO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5cclxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZCgpIHtcclxuICAgIHdpbmRvdy5sb2NhdGlvbi5yZWxvYWQoKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZEFwcCgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1InKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFN5c3RlbURlZmF1bHRUaGVtZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FTRFQnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldExpZ2h0VGhlbWUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBTFQnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldERhcmtUaGVtZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FEVCcpO1xyXG59XHJcblxyXG4vKipcclxuICogUGxhY2UgdGhlIHdpbmRvdyBpbiB0aGUgY2VudGVyIG9mIHRoZSBzY3JlZW5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0NlbnRlcigpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2MnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHdpbmRvdyB0aXRsZVxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdGl0bGVcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFRpdGxlKHRpdGxlKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dUJyArIHRpdGxlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1ha2VzIHRoZSB3aW5kb3cgZ28gZnVsbHNjcmVlblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93RnVsbHNjcmVlbigpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0YnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldmVydHMgdGhlIHdpbmRvdyBmcm9tIGZ1bGxzY3JlZW5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VuZnVsbHNjcmVlbigpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2YnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3M6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTx7dzogbnVtYmVyLCBoOiBudW1iZXJ9Pn0gVGhlIHNpemUgb2YgdGhlIHdpbmRvd1xyXG5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRTaXplKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0U2l6ZVwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1heFNpemUod2lkdGgsIGhlaWdodCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXWjonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d6OicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XHJcbn1cclxuXHJcblxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgd2luZG93IEFsd2F5c09uVG9wIG9yIG5vdCBvbiB0b3BcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEFsd2F5c09uVG9wKGIpIHtcclxuXHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBVFA6JyArIChiID8gJzEnIDogJzAnKSk7XHJcbn1cclxuXHJcblxyXG5cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0geFxyXG4gKiBAcGFyYW0ge251bWJlcn0geVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFBvc2l0aW9uKHgsIHkpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3A6JyArIHggKyAnOicgKyB5KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8e3g6IG51bWJlciwgeTogbnVtYmVyfT59IFRoZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0UG9zaXRpb24oKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dHZXRQb3NcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIaWRlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0hpZGUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dIJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93IHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1Nob3coKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dTJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXhpbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNYXhpbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV00nKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFRvZ2dsZSB0aGUgTWF4aW1pc2Ugb2YgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VG9nZ2xlTWF4aW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d0Jyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBVbm1heGltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWF4aW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dVJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNaW5pbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNaW5pbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV20nKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFVubWluaW1pc2UgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5taW5pbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3UnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gUiBSZWRcclxuICogQHBhcmFtIHtudW1iZXJ9IEcgR3JlZW5cclxuICogQHBhcmFtIHtudW1iZXJ9IEIgQmx1ZVxyXG4gKiBAcGFyYW0ge251bWJlcn0gQSBBbHBoYVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEJhY2tncm91bmRDb2xvdXIoUiwgRywgQiwgQSkge1xyXG4gICAgbGV0IHJnYmEgPSBKU09OLnN0cmluZ2lmeSh7cjogUiB8fCAwLCBnOiBHIHx8IDAsIGI6IEIgfHwgMCwgYTogQSB8fCAyNTV9KTtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcclxuXHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgYWxsIHNjcmVlbnMuIENhbGwgdGhpcyBhbmV3IGVhY2ggdGltZSB5b3Ugd2FudCB0byByZWZyZXNoIGRhdGEgZnJvbSB0aGUgdW5kZXJseWluZyB3aW5kb3dpbmcgc3lzdGVtLlxyXG4gKiBAZXhwb3J0XHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoJy4uL3dyYXBwZXIvcnVudGltZScpLlNjcmVlbn0gU2NyZWVuXHJcbiAqIEByZXR1cm4ge1Byb21pc2U8e1NjcmVlbltdfT59IFRoZSBzY3JlZW5zXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2NyZWVuR2V0QWxsKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6U2NyZWVuR2V0QWxsXCIpO1xyXG59XHJcbiIsICIvKipcclxuICogQGRlc2NyaXB0aW9uOiBVc2UgdGhlIHN5c3RlbSBkZWZhdWx0IGJyb3dzZXIgdG8gb3BlbiB0aGUgdXJsXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnJvd3Nlck9wZW5VUkwodXJsKSB7XHJcbiAgd2luZG93LldhaWxzSW52b2tlKCdCTzonICsgdXJsKTtcclxufSIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcclxuaW1wb3J0IHtldmVudExpc3RlbmVycywgRXZlbnRzRW1pdCwgRXZlbnRzTm90aWZ5LCBFdmVudHNPZmYsIEV2ZW50c09uLCBFdmVudHNPbmNlLCBFdmVudHNPbk11bHRpcGxlfSBmcm9tICcuL2V2ZW50cyc7XHJcbmltcG9ydCB7Q2FsbCwgQ2FsbGJhY2ssIGNhbGxiYWNrc30gZnJvbSAnLi9jYWxscyc7XHJcbmltcG9ydCB7U2V0QmluZGluZ3N9IGZyb20gXCIuL2JpbmRpbmdzXCI7XHJcbmltcG9ydCAqIGFzIFdpbmRvdyBmcm9tIFwiLi93aW5kb3dcIjtcclxuaW1wb3J0ICogYXMgU2NyZWVuIGZyb20gXCIuL3NjcmVlblwiO1xyXG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXJcIjtcclxuXHJcblxyXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUScpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUycpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnSCcpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpFbnZpcm9ubWVudFwiKTtcclxufVxyXG5cclxuLy8gVGhlIEpTIHJ1bnRpbWVcclxud2luZG93LnJ1bnRpbWUgPSB7XHJcbiAgICAuLi5Mb2csXHJcbiAgICAuLi5XaW5kb3csXHJcbiAgICAuLi5Ccm93c2VyLFxyXG4gICAgLi4uU2NyZWVuLFxyXG4gICAgRXZlbnRzT24sXHJcbiAgICBFdmVudHNPbmNlLFxyXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcclxuICAgIEV2ZW50c0VtaXQsXHJcbiAgICBFdmVudHNPZmYsXHJcbiAgICBFbnZpcm9ubWVudCxcclxuICAgIFNob3csXHJcbiAgICBIaWRlLFxyXG4gICAgUXVpdFxyXG59O1xyXG5cclxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXHJcbndpbmRvdy53YWlscyA9IHtcclxuICAgIENhbGxiYWNrLFxyXG4gICAgRXZlbnRzTm90aWZ5LFxyXG4gICAgU2V0QmluZGluZ3MsXHJcbiAgICBldmVudExpc3RlbmVycyxcclxuICAgIGNhbGxiYWNrcyxcclxuICAgIGZsYWdzOiB7XHJcbiAgICAgICAgZGlzYWJsZVNjcm9sbGJhckRyYWc6IGZhbHNlLFxyXG4gICAgICAgIGRpc2FibGVXYWlsc0RlZmF1bHRDb250ZXh0TWVudTogZmFsc2UsXHJcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcclxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxyXG4gICAgICAgIGJvcmRlclRoaWNrbmVzczogNixcclxuICAgICAgICBkYkNsaWNrSW50ZXJ2YWw6IDE1MCxcclxuICAgIH1cclxufTtcclxuXHJcbi8vIFNldCB0aGUgYmluZGluZ3Ncclxud2luZG93LndhaWxzLlNldEJpbmRpbmdzKHdpbmRvdy53YWlsc2JpbmRpbmdzKTtcclxuZGVsZXRlIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncztcclxuXHJcbi8vIFRoaXMgaXMgZXZhbHVhdGVkIGF0IGJ1aWxkIHRpbWUgaW4gcGFja2FnZS5qc29uXHJcbi8vIGNvbnN0IGRldiA9IDA7XHJcbi8vIGNvbnN0IHByb2R1Y3Rpb24gPSAxO1xyXG5pZiAoRU5WID09PSAwKSB7XHJcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XHJcbn1cclxuXHJcbi8vIFNldHVwIGRyYWcgaGFuZGxlclxyXG4vLyBCYXNlZCBvbiBjb2RlIGZyb206IGh0dHBzOi8vZ2l0aHViLmNvbS9wYXRyMG51cy9EZXNrR2FwXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xyXG5cclxuICAgIC8vIENoZWNrIGZvciByZXNpemluZ1xyXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKSB7XHJcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwicmVzaXplOlwiICsgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpO1xyXG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgZm9yIGRyYWdnaW5nXHJcbiAgICBsZXQgY3VycmVudEVsZW1lbnQgPSBlLnRhcmdldDtcclxuICAgIHdoaWxlIChjdXJyZW50RWxlbWVudCAhPSBudWxsKSB7XHJcbiAgICAgICAgaWYgKGN1cnJlbnRFbGVtZW50Lmhhc0F0dHJpYnV0ZSgnZGF0YS13YWlscy1uby1kcmFnJykpIHtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfSBlbHNlIGlmIChjdXJyZW50RWxlbWVudC5oYXNBdHRyaWJ1dGUoJ2RhdGEtd2FpbHMtZHJhZycpKSB7XHJcbiAgICAgICAgICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZVNjcm9sbGJhckRyYWcpIHtcclxuICAgICAgICAgICAgICAgIC8vIFRoaXMgY2hlY2tzIGZvciBjbGlja3Mgb24gdGhlIHNjcm9sbCBiYXJcclxuICAgICAgICAgICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcclxuICAgICAgICAgICAgICAgICAgICBicmVhaztcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGN1cnJlbnRFbGVtZW50ID0gY3VycmVudEVsZW1lbnQucGFyZW50RWxlbWVudDtcclxuICAgIH1cclxufSk7XHJcblxyXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XHJcbiAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvcjtcclxuICAgIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlID0gY3Vyc29yO1xyXG59XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgZnVuY3Rpb24gKGUpIHtcclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVJlc2l6ZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XHJcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvcjtcclxuICAgIH1cclxuICAgIGlmICh3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MgJiYgd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcykge1xyXG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gXCJzZS1yZXNpemVcIjtcclxuICAgIH1cclxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xyXG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcblxyXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcclxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH0gZWxzZSBpZiAocmlnaHRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiB0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcEJvcmRlciAmJiByaWdodEJvcmRlcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcclxuXHJcbn0pO1xyXG5cclxuLy8gU2V0dXAgY29udGV4dCBtZW51IGhvb2tcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgZnVuY3Rpb24gKGUpIHtcclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZVdhaWxzRGVmYXVsdENvbnRleHRNZW51KSB7XHJcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG59KTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBa0JBLDBCQUF3QixPQUFPLFNBQVM7QUFJdkMsV0FBTyxZQUFZLE1BQU0sUUFBUTtBQUFBO0FBUzNCLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG1CQUFpQixTQUFTO0FBQ2hDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLHNCQUFvQixTQUFTO0FBQ25DLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLHVCQUFxQixVQUFVO0FBQ3JDLG1CQUFlLEtBQUs7QUFBQTtBQUlkLE1BQU0sV0FBVztBQUFBLElBQ3ZCLE9BQU87QUFBQSxJQUNQLE9BQU87QUFBQSxJQUNQLE1BQU07QUFBQSxJQUNOLFNBQVM7QUFBQSxJQUNULE9BQU87QUFBQTs7O0FDN0ZSLHVCQUFlO0FBQUEsSUFPWCxZQUFZLFVB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
