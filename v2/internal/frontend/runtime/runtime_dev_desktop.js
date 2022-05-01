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
    WindowSetDarkTheme: () => WindowSetDarkTheme,
    WindowSetLightTheme: () => WindowSetLightTheme,
    WindowSetMaxSize: () => WindowSetMaxSize,
    WindowSetMinSize: () => WindowSetMinSize,
    WindowSetPosition: () => WindowSetPosition,
    WindowSetRGBA: () => WindowSetRGBA,
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
  function WindowSetRGBA(R, G, B, A) {
    let rgba = JSON.stringify({ r: R || 0, g: G || 0, b: B || 0, a: A || 255 });
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
  function Environment() {
    return Call(":wails:Environment");
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
    Environment,
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
      dbClickInterval: 100
    }
  };
  window.wails.SetBindings(window.wailsbindings);
  delete window.wails.SetBindings;
  if (true) {
    delete window.wailsbindings;
  }
  var dragTimeOut;
  var dragLastTime = 0;
  function drag() {
    window.WailsInvoke("drag");
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
        if (new Date().getTime() - dragLastTime < window.wails.flags.dbClickInterval) {
          clearTimeout(dragTimeOut);
          break;
        }
        dragTimeOut = setTimeout(drag, window.wails.flags.dbClickInterval);
        dragLastTime = new Date().getTime();
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbi8qKlxuICogU2VuZHMgYSBsb2cgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB3aXRoIHRoZSBnaXZlbiBsZXZlbCArIG1lc3NhZ2VcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gbGV2ZWxcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmZ1bmN0aW9uIHNlbmRMb2dNZXNzYWdlKGxldmVsLCBtZXNzYWdlKSB7XG5cblx0Ly8gTG9nIE1lc3NhZ2UgZm9ybWF0OlxuXHQvLyBsW3R5cGVdW21lc3NhZ2VdXG5cdHdpbmRvdy5XYWlsc0ludm9rZSgnTCcgKyBsZXZlbCArIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gdHJhY2UgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1RyYWNlKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1QnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dQcmludChtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdQJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBkZWJ1ZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRGVidWcobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gaW5mbyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nSW5mbyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdJJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiB3YXJuaW5nIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dXYXJuaW5nKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1cnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGVycm9yIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dFcnJvcihtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdFJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBmYXRhbCBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRmF0YWwobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRicsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIExvZyBsZXZlbCB0byB0aGUgZ2l2ZW4gbG9nIGxldmVsXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IGxvZ2xldmVsXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRMb2dMZXZlbChsb2dsZXZlbCkge1xuXHRzZW5kTG9nTWVzc2FnZSgnUycsIGxvZ2xldmVsKTtcbn1cblxuLy8gTG9nIGxldmVsc1xuZXhwb3J0IGNvbnN0IExvZ0xldmVsID0ge1xuXHRUUkFDRTogMSxcblx0REVCVUc6IDIsXG5cdElORk86IDMsXG5cdFdBUk5JTkc6IDQsXG5cdEVSUk9SOiA1LFxufTtcbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuLy8gRGVmaW5lcyBhIHNpbmdsZSBsaXN0ZW5lciB3aXRoIGEgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdG8gY2FsbGJhY2tcblxuLyoqXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcbiAqXG4gKiBAY2xhc3MgTGlzdGVuZXJcbiAqL1xuY2xhc3MgTGlzdGVuZXIge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIG1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgLy8gQ2FsbGJhY2sgaW52b2tlcyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xuICAgICAgICAgICAgY2FsbGJhY2suYXBwbHkobnVsbCwgZGF0YSk7XG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXG4gICAgICAgICAgICBpZiAobWF4Q2FsbGJhY2tzID09PSAtMSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcbiAgICAgICAgICAgIG1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIG1heENhbGxiYWNrcyA9PT0gMDtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IHt9O1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ucHVzaCh0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XG5cbiAgICAvLyBHZXQgdGhlIGV2ZW50IG5hbWVcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XG5cbiAgICAvLyBDaGVjayBpZiB3ZSBoYXZlIGFueSBsaXN0ZW5lcnMgZm9yIHRoaXMgZXZlbnRcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSkge1xuXG4gICAgICAgIC8vIEtlZXAgYSBsaXN0IG9mIGxpc3RlbmVyIGluZGV4ZXMgdG8gZGVzdHJveVxuICAgICAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uc2xpY2UoKTtcblxuICAgICAgICAvLyBJdGVyYXRlIGxpc3RlbmVyc1xuICAgICAgICBmb3IgKGxldCBjb3VudCA9IDA7IGNvdW50IDwgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGg7IGNvdW50ICs9IDEpIHtcblxuICAgICAgICAgICAgLy8gR2V0IG5leHQgbGlzdGVuZXJcbiAgICAgICAgICAgIGNvbnN0IGxpc3RlbmVyID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXVtjb3VudF07XG5cbiAgICAgICAgICAgIGxldCBkYXRhID0gZXZlbnREYXRhLmRhdGE7XG5cbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xuICAgICAgICAgICAgY29uc3QgZGVzdHJveSA9IGxpc3RlbmVyLkNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgaWYgKGRlc3Ryb3kpIHtcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxuICAgICAgICAgICAgICAgIG5ld0V2ZW50TGlzdGVuZXJMaXN0LnNwbGljZShjb3VudCwgMSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cblxuICAgICAgICAvLyBVcGRhdGUgY2FsbGJhY2tzIHdpdGggbmV3IGxpc3Qgb2YgbGlzdGVuZXJzXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBuZXdFdmVudExpc3RlbmVyTGlzdDtcbiAgICB9XG59XG5cbi8qKlxuICogTm90aWZ5IGluZm9ybXMgZnJvbnRlbmQgbGlzdGVuZXJzIHRoYXQgYW4gZXZlbnQgd2FzIGVtaXR0ZWQgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBub3RpZnlNZXNzYWdlIC0gZW5jb2RlZCBub3RpZmljYXRpb24gbWVzc2FnZVxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNOb3RpZnkobm90aWZ5TWVzc2FnZSkge1xuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXG4gICAgbGV0IG1lc3NhZ2U7XG4gICAgdHJ5IHtcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zdCBlcnJvciA9ICdJbnZhbGlkIEpTT04gcGFzc2VkIHRvIE5vdGlmeTogJyArIG5vdGlmeU1lc3NhZ2U7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XG4gICAgfVxuICAgIG5vdGlmeUxpc3RlbmVycyhtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNFbWl0KGV2ZW50TmFtZSkge1xuXG4gICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgbmFtZTogZXZlbnROYW1lLFxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxuICAgIH07XG5cbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXG4gICAgbm90aWZ5TGlzdGVuZXJzKHBheWxvYWQpO1xuXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRUUnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmKGV2ZW50TmFtZSkge1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lcnNcbiAgICBkZWxldGUgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXTtcblxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VYJyArIGV2ZW50TmFtZSk7XG59IiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gY3J5cHRvUmFuZG9tKCkge1xuXHR2YXIgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciB1c2luZyBkYSBvbGQtc2tvb2wgTWF0aC5SYW5kb21cbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XG5cdHJldHVybiBNYXRoLnJhbmRvbSgpICogOTAwNzE5OTI1NDc0MDk5MTtcbn1cblxuLy8gUGljayBhIHJhbmRvbSBudW1iZXIgZnVuY3Rpb24gYmFzZWQgb24gYnJvd3NlciBjYXBhYmlsaXR5XG52YXIgcmFuZG9tRnVuYztcbmlmICh3aW5kb3cuY3J5cHRvKSB7XG5cdHJhbmRvbUZ1bmMgPSBjcnlwdG9SYW5kb207XG59IGVsc2Uge1xuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XG59XG5cblxuLyoqXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXG4gKiBiYWNrZW5kIHJlc3BvbmRzLiBUaGlzIHdpbGwgYmUgcmVzb2x2ZWQgd2hlbiB0aGUgY2FsbCB3YXMgc3VjY2Vzc2Z1bFxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXG4gKiB0aW1lIChpbiBtaWxsaXNlY29uZHMpIHRoZW4gdGhlIHByb21pc2UgaXMgcmVqZWN0ZWQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWVcbiAqIEBwYXJhbSB7YW55PX0gYXJnc1xuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XG4gKiBAcmV0dXJuc1xuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XG5cblx0Ly8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcblx0XHR0aW1lb3V0ID0gMDtcblx0fVxuXG5cdC8vIENyZWF0ZSBhIHByb21pc2Vcblx0cmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcblxuXHRcdC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXG5cdFx0dmFyIGNhbGxiYWNrSUQ7XG5cdFx0ZG8ge1xuXHRcdFx0Y2FsbGJhY2tJRCA9IG5hbWUgKyAnLScgKyByYW5kb21GdW5jKCk7XG5cdFx0fSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcblxuXHRcdHZhciB0aW1lb3V0SGFuZGxlO1xuXHRcdC8vIFNldCB0aW1lb3V0XG5cdFx0aWYgKHRpbWVvdXQgPiAwKSB7XG5cdFx0XHR0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XG5cdFx0XHRcdHJlamVjdChFcnJvcignQ2FsbCB0byAnICsgbmFtZSArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xuXHRcdFx0fSwgdGltZW91dCk7XG5cdFx0fVxuXG5cdFx0Ly8gU3RvcmUgY2FsbGJhY2tcblx0XHRjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XG5cdFx0XHR0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxuXHRcdFx0cmVqZWN0OiByZWplY3QsXG5cdFx0XHRyZXNvbHZlOiByZXNvbHZlXG5cdFx0fTtcblxuXHRcdHRyeSB7XG5cdFx0XHRjb25zdCBwYXlsb2FkID0ge1xuXHRcdFx0XHRuYW1lLFxuXHRcdFx0XHRhcmdzLFxuXHRcdFx0XHRjYWxsYmFja0lELFxuXHRcdFx0fTtcblxuXHRcdFx0Ly8gTWFrZSB0aGUgY2FsbFxuXHRcdFx0d2luZG93LldhaWxzSW52b2tlKCdDJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcblx0XHR9IGNhdGNoIChlKSB7XG5cdFx0XHQvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcblx0XHRcdGNvbnNvbGUuZXJyb3IoZSk7XG5cdFx0fVxuXHR9KTtcbn1cblxuXG5cbi8qKlxuICogQ2FsbGVkIGJ5IHRoZSBiYWNrZW5kIHRvIHJldHVybiBkYXRhIHRvIGEgcHJldmlvdXNseSBjYWxsZWRcbiAqIGJpbmRpbmcgaW52b2NhdGlvblxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBpbmNvbWluZ01lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGxiYWNrKGluY29taW5nTWVzc2FnZSkge1xuXHQvLyBQYXJzZSB0aGUgbWVzc2FnZVxuXHRsZXQgbWVzc2FnZTtcblx0dHJ5IHtcblx0XHRtZXNzYWdlID0gSlNPTi5wYXJzZShpbmNvbWluZ01lc3NhZ2UpO1xuXHR9IGNhdGNoIChlKSB7XG5cdFx0Y29uc3QgZXJyb3IgPSBgSW52YWxpZCBKU09OIHBhc3NlZCB0byBjYWxsYmFjazogJHtlLm1lc3NhZ2V9LiBNZXNzYWdlOiAke2luY29taW5nTWVzc2FnZX1gO1xuXHRcdHJ1bnRpbWUuTG9nRGVidWcoZXJyb3IpO1xuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XG5cdH1cblx0bGV0IGNhbGxiYWNrSUQgPSBtZXNzYWdlLmNhbGxiYWNraWQ7XG5cdGxldCBjYWxsYmFja0RhdGEgPSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XG5cdGlmICghY2FsbGJhY2tEYXRhKSB7XG5cdFx0Y29uc3QgZXJyb3IgPSBgQ2FsbGJhY2sgJyR7Y2FsbGJhY2tJRH0nIG5vdCByZWdpc3RlcmVkISEhYDtcblx0XHRjb25zb2xlLmVycm9yKGVycm9yKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZVxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XG5cdH1cblx0Y2xlYXJUaW1lb3V0KGNhbGxiYWNrRGF0YS50aW1lb3V0SGFuZGxlKTtcblxuXHRkZWxldGUgY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xuXG5cdGlmIChtZXNzYWdlLmVycm9yKSB7XG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChtZXNzYWdlLmVycm9yKTtcblx0fSBlbHNlIHtcblx0XHRjYWxsYmFja0RhdGEucmVzb2x2ZShtZXNzYWdlLnJlc3VsdCk7XG5cdH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fICAgIFxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApIFxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vICBcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSAnLi9jYWxscyc7XG5cbi8vIFRoaXMgaXMgd2hlcmUgd2UgYmluZCBnbyBtZXRob2Qgd3JhcHBlcnNcbndpbmRvdy5nbyA9IHt9O1xuXG5leHBvcnQgZnVuY3Rpb24gU2V0QmluZGluZ3MoYmluZGluZ3NNYXApIHtcblx0dHJ5IHtcblx0XHRiaW5kaW5nc01hcCA9IEpTT04ucGFyc2UoYmluZGluZ3NNYXApO1xuXHR9IGNhdGNoIChlKSB7XG5cdFx0Y29uc29sZS5lcnJvcihlKTtcblx0fVxuXG5cdC8vIEluaXRpYWxpc2UgdGhlIGJpbmRpbmdzIG1hcFxuXHR3aW5kb3cuZ28gPSB3aW5kb3cuZ28gfHwge307XG5cblx0Ly8gSXRlcmF0ZSBwYWNrYWdlIG5hbWVzXG5cdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwKS5mb3JFYWNoKChwYWNrYWdlTmFtZSkgPT4ge1xuXG5cdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XG5cdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gfHwge307XG5cblx0XHQvLyBJdGVyYXRlIHN0cnVjdCBuYW1lc1xuXHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXSkuZm9yRWFjaCgoc3RydWN0TmFtZSkgPT4ge1xuXG5cdFx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3Rcblx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdIHx8IHt9O1xuXG5cdFx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0pLmZvckVhY2goKG1ldGhvZE5hbWUpID0+IHtcblxuXHRcdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdW21ldGhvZE5hbWVdID0gZnVuY3Rpb24gKCkge1xuXG5cdFx0XHRcdFx0Ly8gTm8gdGltZW91dCBieSBkZWZhdWx0XG5cdFx0XHRcdFx0bGV0IHRpbWVvdXQgPSAwO1xuXG5cdFx0XHRcdFx0Ly8gQWN0dWFsIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZnVuY3Rpb24gZHluYW1pYygpIHtcblx0XHRcdFx0XHRcdGNvbnN0IGFyZ3MgPSBbXS5zbGljZS5jYWxsKGFyZ3VtZW50cyk7XG5cdFx0XHRcdFx0XHRyZXR1cm4gQ2FsbChbcGFja2FnZU5hbWUsIHN0cnVjdE5hbWUsIG1ldGhvZE5hbWVdLmpvaW4oJy4nKSwgYXJncywgdGltZW91dCk7XG5cdFx0XHRcdFx0fVxuXG5cdFx0XHRcdFx0Ly8gQWxsb3cgc2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZHluYW1pYy5zZXRUaW1lb3V0ID0gZnVuY3Rpb24gKG5ld1RpbWVvdXQpIHtcblx0XHRcdFx0XHRcdHRpbWVvdXQgPSBuZXdUaW1lb3V0O1xuXHRcdFx0XHRcdH07XG5cblx0XHRcdFx0XHQvLyBBbGxvdyBnZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cblx0XHRcdFx0XHRkeW5hbWljLmdldFRpbWVvdXQgPSBmdW5jdGlvbiAoKSB7XG5cdFx0XHRcdFx0XHRyZXR1cm4gdGltZW91dDtcblx0XHRcdFx0XHR9O1xuXG5cdFx0XHRcdFx0cmV0dXJuIGR5bmFtaWM7XG5cdFx0XHRcdH0oKTtcblx0XHRcdH0pO1xuXHRcdH0pO1xuXHR9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZCgpIHtcbiAgICB3aW5kb3cubG9jYXRpb24ucmVsb2FkKCk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWRBcHAoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUicpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U3lzdGVtRGVmYXVsdFRoZW1lKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FTRFQnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldExpZ2h0VGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQUxUJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXREYXJrVGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQURUJyk7XG59XG5cbi8qKlxuICogUGxhY2UgdGhlIHdpbmRvdyBpbiB0aGUgY2VudGVyIG9mIHRoZSBzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dDZW50ZXIoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXYycpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIHdpbmRvdyB0aXRsZVxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZVxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0VGl0bGUodGl0bGUpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dUJyArIHRpdGxlKTtcbn1cblxuLyoqXG4gKiBNYWtlcyB0aGUgd2luZG93IGdvIGZ1bGxzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dGdWxsc2NyZWVuKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0YnKTtcbn1cblxuLyoqXG4gKiBSZXZlcnRzIHRoZSB3aW5kb3cgZnJvbSBmdWxsc2NyZWVuXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5mdWxsc2NyZWVuKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2YnKTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXczonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG4vKipcbiAqIEdldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7dzogbnVtYmVyLCBoOiBudW1iZXJ9Pn0gVGhlIHNpemUgb2YgdGhlIHdpbmRvd1xuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRTaXplKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFNpemVcIik7XG59XG5cbi8qKlxuICogU2V0IHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXWjonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1pblNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0geFxuICogQHBhcmFtIHtudW1iZXJ9IHlcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFBvc2l0aW9uKHgsIHkpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dwOicgKyB4ICsgJzonICsgeSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7eDogbnVtYmVyLCB5OiBudW1iZXJ9Pn0gVGhlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFBvc2l0aW9uKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFBvc1wiKTtcbn1cblxuLyoqXG4gKiBIaWRlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dIaWRlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0gnKTtcbn1cblxuLyoqXG4gKiBTaG93IHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTaG93KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1MnKTtcbn1cblxuLyoqXG4gKiBNYXhpbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXTScpO1xufVxuXG4vKipcbiAqIFRvZ2dsZSB0aGUgTWF4aW1pc2Ugb2YgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1RvZ2dsZU1heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3QnKTtcbn1cblxuLyoqXG4gKiBVbm1heGltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1UnKTtcbn1cblxuLyoqXG4gKiBNaW5pbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWluaW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXbScpO1xufVxuXG4vKipcbiAqIFVubWluaW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWluaW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdScpO1xufVxuXG5cbi8qKlxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSBSIFJlZFxuICogQHBhcmFtIHtudW1iZXJ9IEcgR3JlZW5cbiAqIEBwYXJhbSB7bnVtYmVyfSBCIEJsdWVcbiAqIEBwYXJhbSB7bnVtYmVyfSBBIEFscGhhXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRSR0JBKFIsIEcsIEIsIEEpIHtcbiAgICBsZXQgcmdiYSA9IEpTT04uc3RyaW5naWZ5KHtyOiBSIHx8IDAsIGc6IEcgfHwgMCwgYjogQiB8fCAwLCBhOiBBIHx8IDI1NX0pO1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xufVxuXG4iLCAiLyoqXG4gKiBAZGVzY3JpcHRpb246IFVzZSB0aGUgc3lzdGVtIGRlZmF1bHQgYnJvd3NlciB0byBvcGVuIHRoZSB1cmxcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gQnJvd3Nlck9wZW5VUkwodXJsKSB7XG4gIHdpbmRvdy5XYWlsc0ludm9rZSgnQk86JyArIHVybCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcbmltcG9ydCB7ZXZlbnRMaXN0ZW5lcnMsIEV2ZW50c0VtaXQsIEV2ZW50c05vdGlmeSwgRXZlbnRzT2ZmLCBFdmVudHNPbiwgRXZlbnRzT25jZSwgRXZlbnRzT25NdWx0aXBsZX0gZnJvbSAnLi9ldmVudHMnO1xuaW1wb3J0IHtDYWxsLCBDYWxsYmFjaywgY2FsbGJhY2tzfSBmcm9tICcuL2NhbGxzJztcbmltcG9ydCB7U2V0QmluZGluZ3N9IGZyb20gXCIuL2JpbmRpbmdzXCI7XG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXJcIjtcblxuXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1EnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkVudmlyb25tZW50XCIpO1xufVxuXG4vLyBUaGUgSlMgcnVudGltZVxud2luZG93LnJ1bnRpbWUgPSB7XG4gICAgLi4uTG9nLFxuICAgIC4uLldpbmRvdyxcbiAgICAuLi5Ccm93c2VyLFxuICAgIEV2ZW50c09uLFxuICAgIEV2ZW50c09uY2UsXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcbiAgICBFdmVudHNFbWl0LFxuICAgIEV2ZW50c09mZixcbiAgICBFbnZpcm9ubWVudCxcbiAgICBRdWl0XG59O1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy53YWlscyA9IHtcbiAgICBDYWxsYmFjayxcbiAgICBFdmVudHNOb3RpZnksXG4gICAgU2V0QmluZGluZ3MsXG4gICAgZXZlbnRMaXN0ZW5lcnMsXG4gICAgY2FsbGJhY2tzLFxuICAgIGZsYWdzOiB7XG4gICAgICAgIGRpc2FibGVTY3JvbGxiYXJEcmFnOiBmYWxzZSxcbiAgICAgICAgZGlzYWJsZVdhaWxzRGVmYXVsdENvbnRleHRNZW51OiBmYWxzZSxcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcbiAgICAgICAgZGVmYXVsdEN1cnNvcjogbnVsbCxcbiAgICAgICAgYm9yZGVyVGhpY2tuZXNzOiA2LFxuICAgICAgICBkYkNsaWNrSW50ZXJ2YWw6IDEwMCxcbiAgICB9XG59O1xuXG4vLyBTZXQgdGhlIGJpbmRpbmdzXG53aW5kb3cud2FpbHMuU2V0QmluZGluZ3Mod2luZG93LndhaWxzYmluZGluZ3MpO1xuZGVsZXRlIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncztcblxuLy8gVGhpcyBpcyBldmFsdWF0ZWQgYXQgYnVpbGQgdGltZSBpbiBwYWNrYWdlLmpzb25cbi8vIGNvbnN0IGRldiA9IDA7XG4vLyBjb25zdCBwcm9kdWN0aW9uID0gMTtcbmlmIChFTlYgPT09IDApIHtcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XG59XG5cbnZhciBkcmFnVGltZU91dDtcbnZhciBkcmFnTGFzdFRpbWUgPSAwO1xuXG5mdW5jdGlvbiBkcmFnKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XG59XG5cbi8vIFNldHVwIGRyYWcgaGFuZGxlclxuLy8gQmFzZWQgb24gY29kZSBmcm9tOiBodHRwczovL2dpdGh1Yi5jb20vcGF0cjBudXMvRGVza0dhcFxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIChlKSA9PiB7XG5cbiAgICAvLyBDaGVjayBmb3IgcmVzaXppbmdcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwicmVzaXplOlwiICsgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpO1xuICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBDaGVjayBmb3IgZHJhZ2dpbmdcbiAgICBsZXQgY3VycmVudEVsZW1lbnQgPSBlLnRhcmdldDtcbiAgICB3aGlsZSAoY3VycmVudEVsZW1lbnQgIT0gbnVsbCkge1xuICAgICAgICBpZiAoY3VycmVudEVsZW1lbnQuaGFzQXR0cmlidXRlKCdkYXRhLXdhaWxzLW5vLWRyYWcnKSkge1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH0gZWxzZSBpZiAoY3VycmVudEVsZW1lbnQuaGFzQXR0cmlidXRlKCdkYXRhLXdhaWxzLWRyYWcnKSkge1xuICAgICAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlU2Nyb2xsYmFyRHJhZykge1xuICAgICAgICAgICAgICAgIC8vIFRoaXMgY2hlY2tzIGZvciBjbGlja3Mgb24gdGhlIHNjcm9sbCBiYXJcbiAgICAgICAgICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XG4gICAgICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmIChuZXcgRGF0ZSgpLmdldFRpbWUoKSAtIGRyYWdMYXN0VGltZSA8IHdpbmRvdy53YWlscy5mbGFncy5kYkNsaWNrSW50ZXJ2YWwpIHtcbiAgICAgICAgICAgICAgICBjbGVhclRpbWVvdXQoZHJhZ1RpbWVPdXQpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZHJhZ1RpbWVPdXQgPSBzZXRUaW1lb3V0KGRyYWcsIHdpbmRvdy53YWlscy5mbGFncy5kYkNsaWNrSW50ZXJ2YWwpO1xuICAgICAgICAgICAgZHJhZ0xhc3RUaW1lID0gbmV3IERhdGUoKS5nZXRUaW1lKCk7XG4gICAgICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgICAgICBjdXJyZW50RWxlbWVudCA9IGN1cnJlbnRFbGVtZW50LnBhcmVudEVsZW1lbnQ7XG4gICAgfVxufSk7XG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IpIHtcbiAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvcjtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSA9IGN1cnNvcjtcbn1cblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIGZ1bmN0aW9uIChlKSB7XG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlUmVzaXplKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yID09IG51bGwpIHtcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvcjtcbiAgICB9XG4gICAgaWYgKHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcyAmJiB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzKSB7XG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gXCJzZS1yZXNpemVcIjtcbiAgICB9XG4gICAgbGV0IHJpZ2h0Qm9yZGVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcblxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9IGVsc2UgaWYgKHJpZ2h0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiB0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIgJiYgcmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xuXG59KTtcblxuLy8gU2V0dXAgY29udGV4dCBtZW51IGhvb2tcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGZ1bmN0aW9uIChlKSB7XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlV2FpbHNEZWZhdWx0Q29udGV4dE1lbnUpIHtcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn0pOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFrQkEsMEJBQXdCLE9BQU8sU0FBUztBQUl2QyxXQUFPLFlBQVksTUFBTSxRQUFRO0FBQUE7QUFTM0Isb0JBQWtCLFNBQVM7QUFDakMsbUJBQWUsS0FBSztBQUFBO0FBU2Qsb0JBQWtCLFNBQVM7QUFDakMsbUJBQWUsS0FBSztBQUFBO0FBU2Qsb0JBQWtCLFNBQVM7QUFDakMsbUJBQWUsS0FBSztBQUFBO0FB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
