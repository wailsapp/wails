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
    var message;
    try {
      message = JSON.parse(incomingMessage);
    } catch (e) {
      const error = `Invalid JSON passed to callback: ${e.message}. Message: ${incomingMessage}`;
      wails.LogDebug(error);
      throw new Error(error);
    }
    var callbackID = message.callbackid;
    var callbackData = callbacks[callbackID];
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
    WindowSetMaxSize: () => WindowSetMaxSize,
    WindowSetMinSize: () => WindowSetMinSize,
    WindowSetPosition: () => WindowSetPosition,
    WindowSetRGBA: () => WindowSetRGBA,
    WindowSetSize: () => WindowSetSize,
    WindowSetTitle: () => WindowSetTitle,
    WindowShow: () => WindowShow,
    WindowUnFullscreen: () => WindowUnFullscreen,
    WindowUnmaximise: () => WindowUnmaximise,
    WindowUnminimise: () => WindowUnminimise
  });
  function WindowReload() {
    window.location.reload();
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
  function WindowUnFullscreen() {
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
  function WindowUnmaximise() {
    window.WailsInvoke("WU");
  }
  function WindowMinimise() {
    window.WailsInvoke("Wm");
  }
  function WindowUnminimise() {
    window.WailsInvoke("Wu");
  }
  function WindowSetRGBA(RGBA) {
    let rgba = JSON.stringify(RGBA);
    window.WailsInvoke("Wr:" + rgba);
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
  window.runtime = {
    ...log_exports,
    ...window_exports,
    ...browser_exports,
    EventsOn,
    EventsOnce,
    EventsOnMultiple,
    EventsEmit,
    EventsOff,
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
      borderThickness: 6
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
        e.preventDefault();
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfX1xyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhIGxvZyBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHdpdGggdGhlIGdpdmVuIGxldmVsICsgbWVzc2FnZVxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbGV2ZWxcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmZ1bmN0aW9uIHNlbmRMb2dNZXNzYWdlKGxldmVsLCBtZXNzYWdlKSB7XHJcblxyXG5cdC8vIExvZyBNZXNzYWdlIGZvcm1hdDpcclxuXHQvLyBsW3R5cGVdW21lc3NhZ2VdXHJcblx0d2luZG93LldhaWxzSW52b2tlKCdMJyArIGxldmVsICsgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIHRyYWNlIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nVHJhY2UobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdUJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nUHJpbnQobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdQJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGRlYnVnIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRGVidWcobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdEJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGluZm8gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dJbmZvKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnSScsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiB3YXJuaW5nIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nV2FybmluZyhtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1cnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZXJyb3IgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dFcnJvcihtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0UnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZmF0YWwgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dGYXRhbChtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0YnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIExvZyBsZXZlbCB0byB0aGUgZ2l2ZW4gbG9nIGxldmVsXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IGxvZ2xldmVsXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0TG9nTGV2ZWwobG9nbGV2ZWwpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnUycsIGxvZ2xldmVsKTtcclxufVxyXG5cclxuLy8gTG9nIGxldmVsc1xyXG5leHBvcnQgY29uc3QgTG9nTGV2ZWwgPSB7XHJcblx0VFJBQ0U6IDEsXHJcblx0REVCVUc6IDIsXHJcblx0SU5GTzogMyxcclxuXHRXQVJOSU5HOiA0LFxyXG5cdEVSUk9SOiA1LFxyXG59O1xyXG4iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xyXG5cclxuLyoqXHJcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxyXG4gKlxyXG4gKiBAY2xhc3MgTGlzdGVuZXJcclxuICovXHJcbmNsYXNzIExpc3RlbmVyIHtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cclxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcclxuICAgICAgICBtYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XHJcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxyXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcclxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcclxuICAgICAgICAgICAgY2FsbGJhY2suYXBwbHkobnVsbCwgZGF0YSk7XHJcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcclxuICAgICAgICAgICAgaWYgKG1heENhbGxiYWNrcyA9PT0gLTEpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXHJcbiAgICAgICAgICAgIG1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgICAgICByZXR1cm4gbWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IHt9O1xyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xyXG4gICAgRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcclxufVxyXG5cclxuZnVuY3Rpb24gbm90aWZ5TGlzdGVuZXJzKGV2ZW50RGF0YSkge1xyXG5cclxuICAgIC8vIEdldCB0aGUgZXZlbnQgbmFtZVxyXG4gICAgbGV0IGV2ZW50TmFtZSA9IGV2ZW50RGF0YS5uYW1lO1xyXG5cclxuICAgIC8vIENoZWNrIGlmIHdlIGhhdmUgYW55IGxpc3RlbmVycyBmb3IgdGhpcyBldmVudFxyXG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0pIHtcclxuXHJcbiAgICAgICAgLy8gS2VlcCBhIGxpc3Qgb2YgbGlzdGVuZXIgaW5kZXhlcyB0byBkZXN0cm95XHJcbiAgICAgICAgY29uc3QgbmV3RXZlbnRMaXN0ZW5lckxpc3QgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnNsaWNlKCk7XHJcblxyXG4gICAgICAgIC8vIEl0ZXJhdGUgbGlzdGVuZXJzXHJcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSAwOyBjb3VudCA8IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoOyBjb3VudCArPSAxKSB7XHJcblxyXG4gICAgICAgICAgICAvLyBHZXQgbmV4dCBsaXN0ZW5lclxyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV1bY291bnRdO1xyXG5cclxuICAgICAgICAgICAgbGV0IGRhdGEgPSBldmVudERhdGEuZGF0YTtcclxuXHJcbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xyXG4gICAgICAgICAgICBjb25zdCBkZXN0cm95ID0gbGlzdGVuZXIuQ2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgICAgIGlmIChkZXN0cm95KSB7XHJcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxyXG4gICAgICAgICAgICAgICAgbmV3RXZlbnRMaXN0ZW5lckxpc3Quc3BsaWNlKGNvdW50LCAxKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gVXBkYXRlIGNhbGxiYWNrcyB3aXRoIG5ldyBsaXN0IG9mIGxpc3RlbmVyc1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBuZXdFdmVudExpc3RlbmVyTGlzdDtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcclxuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXHJcbiAgICBsZXQgbWVzc2FnZTtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcbiAgICB9XHJcbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcclxuXHJcbiAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcclxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxyXG4gICAgfTtcclxuXHJcbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXHJcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmYoZXZlbnROYW1lKSB7XHJcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJzXHJcbiAgICBkZWxldGUgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXTtcclxuXHJcbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VYJyArIGV2ZW50TmFtZSk7XHJcbn0iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cclxuICpcclxuICogQHJldHVybnMgbnVtYmVyXHJcbiAqL1xyXG5mdW5jdGlvbiBjcnlwdG9SYW5kb20oKSB7XHJcblx0dmFyIGFycmF5ID0gbmV3IFVpbnQzMkFycmF5KDEpO1xyXG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgbnVtYmVyIHVzaW5nIGRhIG9sZC1za29vbCBNYXRoLlJhbmRvbVxyXG4gKiBJIGxpa2VzIHRvIGNhbGwgaXQgTE9MUmFuZG9tXHJcbiAqXHJcbiAqIEByZXR1cm5zIG51bWJlclxyXG4gKi9cclxuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XHJcblx0cmV0dXJuIE1hdGgucmFuZG9tKCkgKiA5MDA3MTk5MjU0NzQwOTkxO1xyXG59XHJcblxyXG4vLyBQaWNrIGEgcmFuZG9tIG51bWJlciBmdW5jdGlvbiBiYXNlZCBvbiBicm93c2VyIGNhcGFiaWxpdHlcclxudmFyIHJhbmRvbUZ1bmM7XHJcbmlmICh3aW5kb3cuY3J5cHRvKSB7XHJcblx0cmFuZG9tRnVuYyA9IGNyeXB0b1JhbmRvbTtcclxufSBlbHNlIHtcclxuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogQ2FsbCBzZW5kcyBhIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgdG8gY2FsbCB0aGUgYmluZGluZyB3aXRoIHRoZVxyXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXHJcbiAqIGJhY2tlbmQgcmVzcG9uZHMuIFRoaXMgd2lsbCBiZSByZXNvbHZlZCB3aGVuIHRoZSBjYWxsIHdhcyBzdWNjZXNzZnVsXHJcbiAqIG9yIHJlamVjdGVkIGlmIGFuIGVycm9yIGlzIHBhc3NlZCBiYWNrLlxyXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXHJcbiAqIHRpbWUgKGluIG1pbGxpc2Vjb25kcykgdGhlbiB0aGUgcHJvbWlzZSBpcyByZWplY3RlZC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZVxyXG4gKiBAcGFyYW0ge2FueT19IGFyZ3NcclxuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XHJcbiAqIEByZXR1cm5zXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XHJcblxyXG5cdC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxyXG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcclxuXHRcdHRpbWVvdXQgPSAwO1xyXG5cdH1cclxuXHJcblx0Ly8gQ3JlYXRlIGEgcHJvbWlzZVxyXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XHJcblxyXG5cdFx0Ly8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcclxuXHRcdHZhciBjYWxsYmFja0lEO1xyXG5cdFx0ZG8ge1xyXG5cdFx0XHRjYWxsYmFja0lEID0gbmFtZSArICctJyArIHJhbmRvbUZ1bmMoKTtcclxuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XHJcblxyXG5cdFx0dmFyIHRpbWVvdXRIYW5kbGU7XHJcblx0XHQvLyBTZXQgdGltZW91dFxyXG5cdFx0aWYgKHRpbWVvdXQgPiAwKSB7XHJcblx0XHRcdHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcclxuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcclxuXHRcdFx0fSwgdGltZW91dCk7XHJcblx0XHR9XHJcblxyXG5cdFx0Ly8gU3RvcmUgY2FsbGJhY2tcclxuXHRcdGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcclxuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcclxuXHRcdFx0cmVqZWN0OiByZWplY3QsXHJcblx0XHRcdHJlc29sdmU6IHJlc29sdmVcclxuXHRcdH07XHJcblxyXG5cdFx0dHJ5IHtcclxuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcclxuXHRcdFx0XHRuYW1lLFxyXG5cdFx0XHRcdGFyZ3MsXHJcblx0XHRcdFx0Y2FsbGJhY2tJRCxcclxuXHRcdFx0fTtcclxuXHJcblx0XHRcdC8vIE1ha2UgdGhlIGNhbGxcclxuXHRcdFx0d2luZG93LldhaWxzSW52b2tlKCdDJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcclxuXHRcdH0gY2F0Y2ggKGUpIHtcclxuXHRcdFx0Ly8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXHJcblx0XHRcdGNvbnNvbGUuZXJyb3IoZSk7XHJcblx0XHR9XHJcblx0fSk7XHJcbn1cclxuXHJcblxyXG5cclxuLyoqXHJcbiAqIENhbGxlZCBieSB0aGUgYmFja2VuZCB0byByZXR1cm4gZGF0YSB0byBhIHByZXZpb3VzbHkgY2FsbGVkXHJcbiAqIGJpbmRpbmcgaW52b2NhdGlvblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBpbmNvbWluZ01lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsYmFjayhpbmNvbWluZ01lc3NhZ2UpIHtcclxuXHQvLyBEZWNvZGUgdGhlIG1lc3NhZ2UgLSBDcmVkaXQ6IGh0dHBzOi8vc3RhY2tvdmVyZmxvdy5jb20vYS8xMzg2NTY4MFxyXG5cdC8vaW5jb21pbmdNZXNzYWdlID0gZGVjb2RlVVJJQ29tcG9uZW50KGluY29taW5nTWVzc2FnZS5yZXBsYWNlKC9cXHMrL2csICcnKS5yZXBsYWNlKC9bMC05YS1mXXsyfS9nLCAnJSQmJykpO1xyXG5cclxuXHQvLyBQYXJzZSB0aGUgbWVzc2FnZVxyXG5cdHZhciBtZXNzYWdlO1xyXG5cdHRyeSB7XHJcblx0XHRtZXNzYWdlID0gSlNPTi5wYXJzZShpbmNvbWluZ01lc3NhZ2UpO1xyXG5cdH0gY2F0Y2ggKGUpIHtcclxuXHRcdGNvbnN0IGVycm9yID0gYEludmFsaWQgSlNPTiBwYXNzZWQgdG8gY2FsbGJhY2s6ICR7ZS5tZXNzYWdlfS4gTWVzc2FnZTogJHtpbmNvbWluZ01lc3NhZ2V9YDtcclxuXHRcdHdhaWxzLkxvZ0RlYnVnKGVycm9yKTtcclxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcblx0fVxyXG5cdHZhciBjYWxsYmFja0lEID0gbWVzc2FnZS5jYWxsYmFja2lkO1xyXG5cdHZhciBjYWxsYmFja0RhdGEgPSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XHJcblx0aWYgKCFjYWxsYmFja0RhdGEpIHtcclxuXHRcdGNvbnN0IGVycm9yID0gYENhbGxiYWNrICcke2NhbGxiYWNrSUR9JyBub3QgcmVnaXN0ZXJlZCEhIWA7XHJcblx0XHRjb25zb2xlLmVycm9yKGVycm9yKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZVxyXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcclxuXHR9XHJcblx0Y2xlYXJUaW1lb3V0KGNhbGxiYWNrRGF0YS50aW1lb3V0SGFuZGxlKTtcclxuXHJcblx0ZGVsZXRlIGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcclxuXHJcblx0aWYgKG1lc3NhZ2UuZXJyb3IpIHtcclxuXHRcdGNhbGxiYWNrRGF0YS5yZWplY3QobWVzc2FnZS5lcnJvcik7XHJcblx0fSBlbHNlIHtcclxuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcclxuXHR9XHJcbn1cclxuIiwgIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfXyAgICBcclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKSBcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vICBcclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSAnLi9jYWxscyc7XHJcblxyXG4vLyBUaGlzIGlzIHdoZXJlIHdlIGJpbmQgZ28gbWV0aG9kIHdyYXBwZXJzXHJcbndpbmRvdy5nbyA9IHt9O1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFNldEJpbmRpbmdzKGJpbmRpbmdzTWFwKSB7XHJcblx0dHJ5IHtcclxuXHRcdGJpbmRpbmdzTWFwID0gSlNPTi5wYXJzZShiaW5kaW5nc01hcCk7XHJcblx0fSBjYXRjaCAoZSkge1xyXG5cdFx0Y29uc29sZS5lcnJvcihlKTtcclxuXHR9XHJcblxyXG5cdC8vIEluaXRpYWxpc2UgdGhlIGJpbmRpbmdzIG1hcFxyXG5cdHdpbmRvdy5nbyA9IHdpbmRvdy5nbyB8fCB7fTtcclxuXHJcblx0Ly8gSXRlcmF0ZSBwYWNrYWdlIG5hbWVzXHJcblx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXApLmZvckVhY2goKHBhY2thZ2VOYW1lKSA9PiB7XHJcblxyXG5cdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XHJcblx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXSB8fCB7fTtcclxuXHJcblx0XHQvLyBJdGVyYXRlIHN0cnVjdCBuYW1lc1xyXG5cdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdKS5mb3JFYWNoKChzdHJ1Y3ROYW1lKSA9PiB7XHJcblxyXG5cdFx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3RcclxuXHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gfHwge307XHJcblxyXG5cdFx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0pLmZvckVhY2goKG1ldGhvZE5hbWUpID0+IHtcclxuXHJcblx0XHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXVttZXRob2ROYW1lXSA9IGZ1bmN0aW9uICgpIHtcclxuXHJcblx0XHRcdFx0XHQvLyBObyB0aW1lb3V0IGJ5IGRlZmF1bHRcclxuXHRcdFx0XHRcdGxldCB0aW1lb3V0ID0gMDtcclxuXHJcblx0XHRcdFx0XHQvLyBBY3R1YWwgZnVuY3Rpb25cclxuXHRcdFx0XHRcdGZ1bmN0aW9uIGR5bmFtaWMoKSB7XHJcblx0XHRcdFx0XHRcdGNvbnN0IGFyZ3MgPSBbXS5zbGljZS5jYWxsKGFyZ3VtZW50cyk7XHJcblx0XHRcdFx0XHRcdHJldHVybiBDYWxsKFtwYWNrYWdlTmFtZSwgc3RydWN0TmFtZSwgbWV0aG9kTmFtZV0uam9pbignLicpLCBhcmdzLCB0aW1lb3V0KTtcclxuXHRcdFx0XHRcdH1cclxuXHJcblx0XHRcdFx0XHQvLyBBbGxvdyBzZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cclxuXHRcdFx0XHRcdGR5bmFtaWMuc2V0VGltZW91dCA9IGZ1bmN0aW9uIChuZXdUaW1lb3V0KSB7XHJcblx0XHRcdFx0XHRcdHRpbWVvdXQgPSBuZXdUaW1lb3V0O1xyXG5cdFx0XHRcdFx0fTtcclxuXHJcblx0XHRcdFx0XHQvLyBBbGxvdyBnZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cclxuXHRcdFx0XHRcdGR5bmFtaWMuZ2V0VGltZW91dCA9IGZ1bmN0aW9uICgpIHtcclxuXHRcdFx0XHRcdFx0cmV0dXJuIHRpbWVvdXQ7XHJcblx0XHRcdFx0XHR9O1xyXG5cclxuXHRcdFx0XHRcdHJldHVybiBkeW5hbWljO1xyXG5cdFx0XHRcdH0oKTtcclxuXHRcdFx0fSk7XHJcblx0XHR9KTtcclxuXHR9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWQoKSB7XHJcbiAgICB3aW5kb3cubG9jYXRpb24ucmVsb2FkKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBQbGFjZSB0aGUgd2luZG93IGluIHRoZSBjZW50ZXIgb2YgdGhlIHNjcmVlblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93Q2VudGVyKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXYycpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgd2luZG93IHRpdGxlXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZVxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0VGl0bGUodGl0bGUpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1QnICsgdGl0bGUpO1xyXG59XHJcblxyXG4vKipcclxuICogTWFrZXMgdGhlIHdpbmRvdyBnbyBmdWxsc2NyZWVuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dGdWxsc2NyZWVuKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXRicpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV2ZXJ0cyB0aGUgd2luZG93IGZyb20gZnVsbHNjcmVlblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5GdWxsc2NyZWVuKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXZicpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFNpemUod2lkdGgsIGhlaWdodCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXczonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHt3OiBudW1iZXIsIGg6IG51bWJlcn0+fSBUaGUgc2l6ZSBvZiB0aGUgd2luZG93XHJcblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFNpemUoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dHZXRTaXplXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1daOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB5XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0UG9zaXRpb24oeCwgeSkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXcDonICsgeCArICc6JyArIHkpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTx7eDogbnVtYmVyLCB5OiBudW1iZXJ9Pn0gVGhlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRQb3NpdGlvbigpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFBvc1wiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEhpZGUgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SGlkZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0gnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3cgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2hvdygpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1MnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1heGltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd01heGltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXTScpO1xyXG59XHJcblxyXG4vKipcclxuICogVW5tYXhpbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1heGltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVScpO1xyXG59XHJcblxyXG4vKipcclxuICogTWluaW1pc2UgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWluaW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dtJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBVbm1pbmltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWluaW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d1Jyk7XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7UkdCQX0gUkdCQSBiYWNrZ3JvdW5kIGNvbG91clxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFJHQkEoUkdCQSkge1xyXG4gICAgbGV0IHJnYmEgPSBKU09OLnN0cmluZ2lmeShSR0JBKTtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xyXG59XHJcblxyXG4iLCAiLyoqXHJcbiAqIEBkZXNjcmlwdGlvbjogVXNlIHRoZSBzeXN0ZW0gZGVmYXVsdCBicm93c2VyIHRvIG9wZW4gdGhlIHVybFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIFxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJyb3dzZXJPcGVuVVJMKHVybCkge1xyXG4gIHdpbmRvdy5XYWlsc0ludm9rZSgnQk86JyArIHVybCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XHJcbmltcG9ydCB7ZXZlbnRMaXN0ZW5lcnMsIEV2ZW50c0VtaXQsIEV2ZW50c05vdGlmeSwgRXZlbnRzT2ZmLCBFdmVudHNPbiwgRXZlbnRzT25jZSwgRXZlbnRzT25NdWx0aXBsZX0gZnJvbSAnLi9ldmVudHMnO1xyXG5pbXBvcnQge0NhbGxiYWNrLCBjYWxsYmFja3N9IGZyb20gJy4vY2FsbHMnO1xyXG5pbXBvcnQge1NldEJpbmRpbmdzfSBmcm9tIFwiLi9iaW5kaW5nc1wiO1xyXG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XHJcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3NlclwiO1xyXG5cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdRJyk7XHJcbn1cclxuXHJcbi8vIFRoZSBKUyBydW50aW1lXHJcbndpbmRvdy5ydW50aW1lID0ge1xyXG4gICAgLi4uTG9nLFxyXG4gICAgLi4uV2luZG93LFxyXG4gICAgLi4uQnJvd3NlcixcclxuICAgIEV2ZW50c09uLFxyXG4gICAgRXZlbnRzT25jZSxcclxuICAgIEV2ZW50c09uTXVsdGlwbGUsXHJcbiAgICBFdmVudHNFbWl0LFxyXG4gICAgRXZlbnRzT2ZmLFxyXG4gICAgUXVpdFxyXG59O1xyXG5cclxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXHJcbndpbmRvdy53YWlscyA9IHtcclxuICAgIENhbGxiYWNrLFxyXG4gICAgRXZlbnRzTm90aWZ5LFxyXG4gICAgU2V0QmluZGluZ3MsXHJcbiAgICBldmVudExpc3RlbmVycyxcclxuICAgIGNhbGxiYWNrcyxcclxuICAgIGZsYWdzOiB7XHJcbiAgICAgICAgZGlzYWJsZVNjcm9sbGJhckRyYWc6IGZhbHNlLFxyXG4gICAgICAgIGRpc2FibGVXYWlsc0RlZmF1bHRDb250ZXh0TWVudTogZmFsc2UsXHJcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcclxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxyXG4gICAgICAgIGJvcmRlclRoaWNrbmVzczogNlxyXG4gICAgfVxyXG59O1xyXG5cclxuLy8gU2V0IHRoZSBiaW5kaW5nc1xyXG53aW5kb3cud2FpbHMuU2V0QmluZGluZ3Mod2luZG93LndhaWxzYmluZGluZ3MpO1xyXG5kZWxldGUgd2luZG93LndhaWxzLlNldEJpbmRpbmdzO1xyXG5cclxuLy8gVGhpcyBpcyBldmFsdWF0ZWQgYXQgYnVpbGQgdGltZSBpbiBwYWNrYWdlLmpzb25cclxuLy8gY29uc3QgZGV2ID0gMDtcclxuLy8gY29uc3QgcHJvZHVjdGlvbiA9IDE7XHJcbmlmIChFTlYgPT09IDApIHtcclxuICAgIGRlbGV0ZSB3aW5kb3cud2FpbHNiaW5kaW5ncztcclxufVxyXG5cclxuLy8gU2V0dXAgZHJhZyBoYW5kbGVyXHJcbi8vIEJhc2VkIG9uIGNvZGUgZnJvbTogaHR0cHM6Ly9naXRodWIuY29tL3BhdHIwbnVzL0Rlc2tHYXBcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIChlKSA9PiB7XHJcblxyXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nXHJcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJyZXNpemU6XCIgKyB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSk7XHJcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBmb3IgZHJhZ2dpbmdcclxuICAgIGxldCBjdXJyZW50RWxlbWVudCA9IGUudGFyZ2V0O1xyXG4gICAgd2hpbGUgKGN1cnJlbnRFbGVtZW50ICE9IG51bGwpIHtcclxuICAgICAgICBpZiAoY3VycmVudEVsZW1lbnQuaGFzQXR0cmlidXRlKCdkYXRhLXdhaWxzLW5vLWRyYWcnKSkge1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9IGVsc2UgaWYgKGN1cnJlbnRFbGVtZW50Lmhhc0F0dHJpYnV0ZSgnZGF0YS13YWlscy1kcmFnJykpIHtcclxuICAgICAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlU2Nyb2xsYmFyRHJhZykge1xyXG4gICAgICAgICAgICAgICAgLy8gVGhpcyBjaGVja3MgZm9yIGNsaWNrcyBvbiB0aGUgc2Nyb2xsIGJhclxyXG4gICAgICAgICAgICAgICAgaWYgKGUub2Zmc2V0WCA+IGUudGFyZ2V0LmNsaWVudFdpZHRoIHx8IGUub2Zmc2V0WSA+IGUudGFyZ2V0LmNsaWVudEhlaWdodCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XHJcbiAgICAgICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGN1cnJlbnRFbGVtZW50ID0gY3VycmVudEVsZW1lbnQucGFyZW50RWxlbWVudDtcclxuICAgIH1cclxufSk7XHJcblxyXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XHJcbiAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvcjtcclxuICAgIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlID0gY3Vyc29yO1xyXG59XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgZnVuY3Rpb24gKGUpIHtcclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVJlc2l6ZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XHJcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvcjtcclxuICAgIH1cclxuICAgIGlmICh3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MgJiYgd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcykge1xyXG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gXCJzZS1yZXNpemVcIjtcclxuICAgIH1cclxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xyXG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcblxyXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcclxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH0gZWxzZSBpZiAocmlnaHRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiB0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcEJvcmRlciAmJiByaWdodEJvcmRlcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcclxuXHJcbn0pO1xyXG5cclxuLy8gU2V0dXAgY29udGV4dCBtZW51IGhvb2tcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgZnVuY3Rpb24gKGUpIHtcclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZVdhaWxzRGVmYXVsdENvbnRleHRNZW51KSB7XHJcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG59KTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBa0JBLDBCQUF3QixPQUFPLFNBQVM7QUFJdkMsV0FBTyxZQUFZLE1BQU0sUUFBUTtBQUFBO0FBUzNCLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG1CQUFpQixTQUFTO0FBQ2hDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLHNCQUFvQixTQUFTO0FBQ25DLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLG9CQUFrQixTQUFTO0FBQ2pDLG1CQUFlLEtBQUs7QUFBQTtBQVNkLHVCQUFxQixVQUFVO0FBQ3JDLG1CQUFlLEtBQUs7QUFBQTtBQUlkLE1BQU0sV0FBVztBQUFBLElBQ3ZCLE9BQU87QUFBQSxJQUNQLE9BQU87QUFBQSxJQUNQLE1BQU07QUFBQSxJQUNOLFNBQVM7QUFBQSxJQUNULE9BQU87QUFBQTs7O0FDN0ZSLHVCQUFlO0FBQUEsSUFPWCxZQUFZLFVBQVUsY0FBYztBQUVoQyxxQkFBZSxnQkFBZ0I7QUFHL0IsV0FBSyxXQUFXLENBQUMsU0FBUztBQUN0QixpQkFBUyxNQUFNLE1BQU07QUFFckIsWUFBSSxpQkFBaUIsSUFBSTtBQUNyQixpQkFBTztBQUFBO0FBR1gsd0JBQWdCO0FBQ2hCLGVBQU8saUJBQWlCO0FBQUE7QUFBQTtBQUFBO0FBSzdCLE1BQU0saUJBQWlCO0FBVXZCLDRCQUEwQixXQUFXLFVBQVUsY0FBYztBQUNoRSxtQkFBZSxhQUFhLGVBQWUsY0FBYztBQUN6RCxVQUFNLGVBQWUsSUFBSSxTQUFTLFVBQVU7QUFDNUMsbUJBQWUsV0FBVyxLQUFLO0FBQUE7QUFVNUIsb0JBQWtCLFdBQVcsVUFBVTtBQUMxQyxxQkFBaUIsV0FBVyxVQUFVO0FBQUE7QUFVbkMsc0JBQW9CLFdBQVcsVUFBVTtBQUM1QyxxQkFBaUIsV0FBVyxVQUFVO0FBQUE7QUFHMUMsMkJBQXlCLFdBQVc7QUFHaEMsUUFBSSxZQUFZLFVBQVU7QUFHMUIsUUFBSSxlQUFlLFlBQVk7QUFHM0IsWUFBTSx1QkFBdUIsZUFBZSxXQUFXO0FBR3ZELGVBQVMsUUFBUSxHQUFHLFFBQVEsZUFBZSxXQUFXLFFBQVEsU0FBUyxHQUFHO0FBR3RFLGNBQU0sV0FBVyxlQUFlLFdBQVc7QUFFM0MsWUFBSSxPQUFPLFVBQVU7QUFHckIsY0FBTSxVQUFVLFNBQVMsU0FBUztBQUNsQyxZQUFJLFNBQVM7QUFFVCwrQkFBcUIsT0FBTyxPQUFPO0FBQUE7QUFBQTtBQUszQyxxQkFBZSxhQUFhO0FBQUE7QUFBQTtBQVc3Qix3QkFBc0IsZUFBZTtBQUV4QyxRQUFJO0FBQ0osUUFBSTtBQUNBLGdCQUFVLEtBQUssTUFBTTtBQUFBLGFBQ2hCLEdBQVA7QUFDRSxZQUFNLFFBQVEsb0NBQW9DO0FBQ2xELFlBQU0sSUFBSSxNQUFNO0FBQUE7QUFFcEIsb0JBQWdCO0FBQUE7QUFTYixzQkFBb0IsV0FBVztBQUVsQyxVQUFNLFVBQVU7QUFBQSxNQUNaLE1BQU07QUFBQSxNQUNOLE1BQU0sR0FBRyxNQUFNLE1BQU0sV0FBVyxNQUFNO0FBQUE7QUFJMUMsb0JBQWdCO0FBR2hCLFdBQU8sWUFBWSxPQUFPLEtBQUssVUFBVTtBQUFBO0FBR3RDLHFCQUFtQixXQUFXO0FBRWpDLFdBQU8sZUFBZTtBQUd0QixXQUFPLFlBQVksT0FBTztBQUFBOzs7QUNsSnZCLE1BQU0sWUFBWTtBQU96QiwwQkFBd0I7QUFDdkIsUUFBSSxRQUFRLElBQUksWUFBWTtBQUM1QixXQUFPLE9BQU8sT0FBTyxnQkFBZ0IsT0FBTztBQUFBO0FBUzdDLHlCQUF1QjtBQUN0QixXQUFPLEtBQUssV0FBVztBQUFBO0FBSXhCLE1BQUk7QUFDSixNQUFJLE9BQU8sUUFBUTtBQUNsQixpQkFBYTtBQUFBLFNBQ1A7QUFDTixpQkFBYTtBQUFBO0FBa0JQLGdCQUFjLE1BQU0sTUFBTSxTQUFTO0FBR3pDLFFBQUksV0FBVyxNQUFNO0FBQ3BCLGdCQUFVO0FBQUE7QUFJWCxXQUFPLElBQUksUUFBUSxTQUFVLFNBQVMsUUFBUTtBQUc3QyxVQUFJO0FBQ0osU0FBRztBQUNGLHFCQUFhLE9BQU8sTUFBTTtBQUFBLGVBQ2xCLFVBQVU7QUFFbkIsVUFBSTtBQUVKLFVBQUksVUFBVSxHQUFHO0FBQ2hCLHdCQUFnQixXQUFXLFdBQVk7QUFDdEMsaUJBQU8sTUFBTSxhQUFhLE9BQU8sNkJBQTZCO0FBQUEsV0FDNUQ7QUFBQTtBQUlKLGdCQUFVLGNBQWM7QUFBQSxRQUN2QjtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUE7QUFHRCxVQUFJO0FBQ0gsY0FBTSxVQUFVO0FBQUEsVUFDZjtBQUFBLFVBQ0E7QUFBQSxVQUNBO0FBQUE7QUFJRCxlQUFPLFlBQVksTUFBTSxLQUFLLFVBQVU7QUFBQSxlQUNoQyxHQUFQO0FBRUQsZ0JBQVEsTUFBTTtBQUFBO0FBQUE7QUFBQTtBQWNWLG9CQUFrQixpQkFBaUI7QUFLekMsUUFBSTtBQUNKLFFBQUk7QUFDSCxnQkFBVSxLQUFLLE1BQU07QUFBQSxhQUNiLEdBQVA7QUFDRCxZQUFNLFFBQVEsb0NBQW9DLEVBQUUscUJBQXFCO0FBQ3pFLFlBQU0sU0FBUztBQUNmLFlBQU0sSUFBSSxNQUFNO0FBQUE7QUFFakIsUUFBSSxhQUFhLFFBQVE7QUFDekIsUUFBSSxlQUFlLFVBQVU7QUFDN0IsUUFBSSxDQUFDLGNBQWM7QUFDbEIsWUFBTSxRQUFRLGFBQWE7QUFDM0IsY0FBUSxNQUFNO0FBQ2QsWUFBTSxJQUFJLE1BQU07QUFBQTtBQUVqQixpQkFBYSxhQUFhO0FBRTFCLFdBQU8sVUFBVTtBQUVqQixRQUFJLFFBQVEsT0FBTztBQUNsQixtQkFBYSxPQUFPLFFBQVE7QUFBQSxXQUN0QjtBQUNOLG1CQUFhLFFBQVEsUUFBUTtBQUFBO0FBQUE7OztBQzdIL0IsU0FBTyxLQUFLO0FBRUwsdUJBQXFCLGFBQWE7QUFDeEMsUUFBSTtBQUNILG9CQUFjLEtBQUssTUFBTTtBQUFBLGFBQ2pCLEdBQVA7QUFDRCxjQUFRLE1BQU07QUFBQTtBQUlmLFdBQU8sS0FBSyxPQUFPLE1BQU07QUFHekIsV0FBTyxLQUFLLGFBQWEsUUFBUSxDQUFDLGdCQUFnQjtBQUdqRCxhQUFPLEdBQUcsZUFBZSxPQUFPLEdBQUcsZ0JBQWdCO0FBR25ELGFBQU8sS0FBSyxZQUFZLGNBQWMsUUFBUSxDQUFDLGVBQWU7QUFHN0QsZUFBTyxHQUFHLGFBQWEsY0FBYyxPQUFPLEdBQUcsYUFBYSxlQUFlO0FBRTNFLGVBQU8sS0FBSyxZQUFZLGFBQWEsYUFBYSxRQUFRLENBQUMsZUFBZTtBQUV6RSxpQkFBTyxHQUFHLGFBQWEsWUFBWSxjQUFjLFdBQVk7QUFHNUQsZ0JBQUksVUFBVTtBQUdkLCtCQUFtQjtBQUNsQixvQkFBTSxPQUFPLEdBQUcsTUFBTSxLQUFLO0FBQzNCLHFCQUFPLEtBQUssQ0FBQyxhQUFhLFlBQVksWUFBWSxLQUFLLE1BQU0sTUFBTTtBQUFBO0FBSXBFLG9CQUFRLGFBQWEsU0FBVSxZQUFZO0FBQzFDLHdCQUFVO0FBQUE7QUFJWCxvQkFBUSxhQUFhLFdBQVk7QUFDaEMscUJBQU87QUFBQTtBQUdSLG1CQUFPO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDN0RaO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWVPLDBCQUF3QjtBQUMzQixXQUFPLFNBQVM7QUFBQTtBQVFiLDBCQUF3QjtBQUMzQixXQUFPLFlBQVk7QUFBQTtBQVNoQiwwQkFBd0IsT0FBTztBQUNsQyxXQUFPLFlBQVksT0FBTztBQUFBO0FBUXZCLDhCQUE0QjtBQUMvQixXQUFPLFlBQVk7QUFBQTtBQVFoQixnQ0FBOEI7QUFDakMsV0FBTyxZQUFZO0FBQUE7QUFVaEIseUJBQXVCLE9BQU8sUUFBUTtBQUN6QyxXQUFPLFlBQVksUUFBUSxRQUFRLE1BQU07QUFBQTtBQVV0QywyQkFBeUI7QUFDNUIsV0FBTyxLQUFLO0FBQUE7QUFVVCw0QkFBMEIsT0FBTyxRQUFRO0FBQzVDLFdBQU8sWUFBWSxRQUFRLFFBQVEsTUFBTTtBQUFBO0FBVXRDLDRCQUEwQixPQUFPLFFBQVE7QUFDNUMsV0FBTyxZQUFZLFFBQVEsUUFBUSxNQUFNO0FBQUE7QUFVdEMsNkJBQTJCLEdBQUcsR0FBRztBQUNwQyxXQUFPLFlBQVksUUFBUSxJQUFJLE1BQU07QUFBQTtBQVNsQywrQkFBNkI7QUFDaEMsV0FBTyxLQUFLO0FBQUE7QUFRVCx3QkFBc0I7QUFDekIsV0FBTyxZQUFZO0FBQUE7QUFRaEIsd0JBQXNCO0FBQ3pCLFdBQU8sWUFBWTtBQUFBO0FBUWhCLDRCQUEwQjtBQUM3QixXQUFPLFlBQVk7QUFBQTtBQVFoQiw4QkFBNEI7QUFDL0IsV0FBTyxZQUFZO0FBQUE7QUFRaEIsNEJBQTBCO0FBQzdCLFdBQU8sWUFBWTtBQUFBO0FBUWhCLDhCQUE0QjtBQUMvQixXQUFPLFlBQVk7QUFBQTtBQVVoQix5QkFBdUIsTUFBTTtBQUNoQyxRQUFJLE9BQU8sS0FBSyxVQUFVO0FBQzFCLFdBQU8sWUFBWSxRQUFRO0FBQUE7OztBQ3hML0I7QUFBQTtBQUFBO0FBQUE7QUFLTywwQkFBd0IsS0FBSztBQUNsQyxXQUFPLFlBQVksUUFBUTtBQUFBOzs7QUNZdEIsa0JBQWdCO0FBQ25CLFdBQU8sWUFBWTtBQUFBO0FBSXZCLFNBQU8sVUFBVTtBQUFBLE9BQ1Y7QUFBQSxPQUNBO0FBQUEsT0FDQTtBQUFBLElBQ0g7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBO0FBSUosU0FBTyxRQUFRO0FBQUEsSUFDWDtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBLE9BQU87QUFBQSxNQUNILHNCQUFzQjtBQUFBLE1BQ3RCLGdDQUFnQztBQUFBLE1BQ2hDLGNBQWM7QUFBQSxNQUNkLGVBQWU7QUFBQSxNQUNmLGlCQUFpQjtBQUFBO0FBQUE7QUFLekIsU0FBTyxNQUFNLFlBQVksT0FBTztBQUNoQyxTQUFPLE9BQU8sTUFBTTtBQUtwQixNQUFJLE1BQVc7QUFDWCxXQUFPLE9BQU87QUFBQTtBQUtsQixTQUFPLGlCQUFpQixhQUFhLENBQUMsTUFBTTtBQUd4QyxRQUFJLE9BQU8sTUFBTSxNQUFNLFlBQVk7QUFDL0IsYUFBTyxZQUFZLFlBQVksT0FBTyxNQUFNLE1BQU07QUFDbEQsUUFBRTtBQUNGO0FBQUE7QUFJSixRQUFJLGlCQUFpQixFQUFFO0FBQ3ZCLFdBQU8sa0JBQWtCLE1BQU07QUFDM0IsVUFBSSxlQUFlLGFBQWEsdUJBQXVCO0FBQ25EO0FBQUEsaUJBQ08sZUFBZSxhQUFhLG9CQUFvQjtBQUN2RCxZQUFJLE9BQU8sTUFBTSxNQUFNLHNCQUFzQjtBQUV6QyxjQUFJLEVBQUUsVUFBVSxFQUFFLE9BQU8sZUFBZSxFQUFFLFVBQVUsRUFBRSxPQUFPLGNBQWM7QUFDdkU7QUFBQTtBQUFBO0FBR1IsZUFBTyxZQUFZO0FBQ25CLFVBQUU7QUFDRjtBQUFBO0FBRUosdUJBQWlCLGVBQWU7QUFBQTtBQUFBO0FBSXhDLHFCQUFtQixRQUFRO0FBQ3ZCLGFBQVMsS0FBSyxNQUFNLFNBQVMsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUMxRCxXQUFPLE1BQU0sTUFBTSxhQUFhO0FBQUE7QUFHcEMsU0FBTyxpQkFBaUIsYUFBYSxTQUFVLEdBQUc7QUFDOUMsUUFBSSxDQUFDLE9BQU8sTUFBTSxNQUFNLGNBQWM7QUFDbEM7QUFBQTtBQUVKLFFBQUksT0FBTyxNQUFNLE1BQU0saUJBQWlCLE1BQU07QUFDMUMsYUFBTyxNQUFNLE1BQU0sZ0JBQWdCLFNBQVMsS0FBSyxNQUFNO0FBQUE7QUFFM0QsUUFBSSxPQUFPLGFBQWEsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNLG1CQUFtQixPQUFPLGNBQWMsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNLGlCQUFpQjtBQUMzSSxlQUFTLEtBQUssTUFBTSxTQUFTO0FBQUE7QUFFakMsUUFBSSxjQUFjLE9BQU8sYUFBYSxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU07QUFDckUsUUFBSSxhQUFhLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUNoRCxRQUFJLFlBQVksRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBQy9DLFFBQUksZUFBZSxPQUFPLGNBQWMsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBR3ZFLFFBQUksQ0FBQyxjQUFjLENBQUMsZUFBZSxDQUFDLGFBQWEsQ0FBQyxnQkFBZ0IsT0FBTyxNQUFNLE1BQU0sZUFBZSxRQUFXO0FBQzNHO0FBQUEsZUFDTyxlQUFlO0FBQWMsZ0JBQVU7QUFBQSxhQUN6QyxjQUFjO0FBQWMsZ0JBQVU7QUFBQSxhQUN0QyxjQUFjO0FBQVcsZ0JBQVU7QUFBQSxhQUNuQyxhQUFhO0FBQWEsZ0JBQVU7QUFBQSxhQUNwQztBQUFZLGdCQUFVO0FBQUEsYUFDdEI7QUFBVyxnQkFBVTtBQUFBLGFBQ3JCO0FBQWMsZ0JBQVU7QUFBQSxhQUN4QjtBQUFhLGdCQUFVO0FBQUE7QUFLcEMsU0FBTyxpQkFBaUIsZUFBZSxTQUFVLEdBQUc7QUFDaEQsUUFBSSxPQUFPLE1BQU0sTUFBTSxnQ0FBZ0M7QUFDbkQsUUFBRTtBQUFBO0FBQUE7IiwKICAibmFtZXMiOiBbXQp9Cg==
