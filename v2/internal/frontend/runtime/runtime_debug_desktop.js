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
    return new Promise(function(resolve, reject2) {
      var callbackID;
      do {
        callbackID = name + "-" + randomFunc();
      } while (callbacks[callbackID]);
      var timeoutHandle;
      if (timeout > 0) {
        timeoutHandle = setTimeout(function() {
          reject2(Error("Call to " + name + " timed out. Request ID: " + callbackID));
        }, timeout);
      }
      callbacks[callbackID] = {
        timeoutHandle,
        reject: reject2,
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
    return new Promise(function(resolve, reject2) {
      var callbackID;
      do {
        callbackID = id + "-" + randomFunc();
      } while (callbacks[callbackID]);
      var timeoutHandle;
      if (timeout > 0) {
        timeoutHandle = setTimeout(function() {
          reject2(Error("Call to method " + id + " timed out. Request ID: " + callbackID));
        }, timeout);
      }
      callbacks[callbackID] = {
        timeoutHandle,
        reject: reject2,
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

  // desktop/draganddrop.js
  var draganddrop_exports = {};
  __export(draganddrop_exports, {
    CanResolveFilePaths: () => CanResolveFilePaths,
    ResolveFilePaths: () => ResolveFilePaths
  });
  function CanResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
  }
  function ResolveFilePaths(files) {
    if (!window.chrome?.webview?.postMessageWithAdditionalObjects) {
      reject(new Error("Unsupported Platform"));
      return;
    }
    chrome.webview.postMessageWithAdditionalObjects(`file:drop:`, files);
  }

  // desktop/contextmenu.js
  function processDefaultContextMenu(event) {
    const element = event.target;
    const computedStyle = window.getComputedStyle(element);
    const defaultContextMenuAction = computedStyle.getPropertyValue("--default-contextmenu").trim();
    switch (defaultContextMenuAction) {
      case "show":
        return;
      case "hide":
        event.preventDefault();
        return;
      default:
        if (element.isContentEditable) {
          return;
        }
        const selection = window.getSelection();
        const hasSelection = selection.toString().length > 0;
        if (hasSelection) {
          for (let i = 0; i < selection.rangeCount; i++) {
            const range = selection.getRangeAt(i);
            const rects = range.getClientRects();
            for (let j = 0; j < rects.length; j++) {
              const rect = rects[j];
              if (document.elementFromPoint(rect.left, rect.top) === element) {
                return;
              }
            }
          }
        }
        if (element.tagName === "INPUT" || element.tagName === "TEXTAREA") {
          if (hasSelection || !element.readOnly && !element.disabled) {
            return;
          }
        }
        event.preventDefault();
    }
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
    ...draganddrop_exports,
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
      disableDefaultContextMenu: false,
      enableResize: false,
      defaultCursor: null,
      borderThickness: 6,
      shouldDrag: false,
      deferDragToMouseMove: true,
      cssDragProperty: "--wails-draggable",
      cssDragValue: "drag",
      cssDropProperty: "--wails-drop-target",
      cssDropValue: "drop",
      enableWailsDragAndDrop: false,
      wailsDropPreviousElement: null
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
  window.wails.setCSSDropProperties = function(property, value) {
    window.wails.flags.cssDropProperty = property;
    window.wails.flags.cssDropValue = value;
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
    if (true)
      return;
    if (window.wails.flags.disableDefaultContextMenu) {
      e.preventDefault();
    } else {
      processDefaultContextMenu(e);
    }
  });
  window.addEventListener("dragover", function(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
      return;
    }
    e.preventDefault();
    let targetElement = document.elementFromPoint(e.x, e.y);
    if (targetElement === window.wails.flags.wailsDropPreviousElement) {
      return;
    }
    const style = targetElement.style;
    let cssDropValue = null;
    if (Object.keys(style).findIndex((key) => style[key] === window.wails.flags.cssDropProperty) < 0) {
      targetElement = targetElement.closest(`[style*='${window.wails.flags.cssDropProperty}']`);
    }
    if (targetElement == null) {
      return;
    }
    cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
      cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue === window.wails.flags.cssDropValue) {
      targetElement.classList.add("wails-drop-target-active");
    } else if (window.wails.flags.wailsDropPreviousElement) {
      window.wails.flags.wailsDropPreviousElement.classList.remove("wails-drop-target-active");
    }
    window.wails.flags.wailsDropPreviousElement = targetElement;
  });
  window.addEventListener("dragleave", function(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
      return;
    }
    e.preventDefault();
    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
      cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue && window.wails.flags.wailsDropPreviousElement) {
      window.wails.flags.wailsDropPreviousElement.classList.remove("wails-drop-target-active");
    }
  });
  window.addEventListener("drop", function(e) {
    if (!window.wails.flags.enableWailsDragAndDrop) {
      return;
    }
    e.preventDefault();
    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
      cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue) {
      return;
    }
    let files = [];
    if (e.dataTransfer.items) {
      files = [...e.dataTransfer.items].map((item, i) => {
        if (item.kind === "file") {
          const file = item.getAsFile();
          return file;
        }
      });
    } else {
      files = [...e.dataTransfer.files];
    }
    window.runtime.ResolveFilePaths(files);
    if (window.wails.flags.wailsDropPreviousElement) {
      window.wails.flags.wailsDropPreviousElement.classList.remove("wails-drop-target-active");
    }
  });
  window.WailsInvoke("runtime:ready");
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL2RyYWdhbmRkcm9wLmpzIiwgImRlc2t0b3AvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9tYWluLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX19cclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG4vKipcclxuICogU2VuZHMgYSBsb2cgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB3aXRoIHRoZSBnaXZlbiBsZXZlbCArIG1lc3NhZ2VcclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGxldmVsXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5mdW5jdGlvbiBzZW5kTG9nTWVzc2FnZShsZXZlbCwgbWVzc2FnZSkge1xyXG5cclxuXHQvLyBMb2cgTWVzc2FnZSBmb3JtYXQ6XHJcblx0Ly8gbFt0eXBlXVttZXNzYWdlXVxyXG5cdHdpbmRvdy5XYWlsc0ludm9rZSgnTCcgKyBsZXZlbCArIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiB0cmFjZSBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ1RyYWNlKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnVCcsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ1ByaW50KG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnUCcsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBkZWJ1ZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ0RlYnVnKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnRCcsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBpbmZvIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nSW5mbyhtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0knLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gd2FybmluZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ1dhcm5pbmcobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdXJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGVycm9yIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRXJyb3IobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdFJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGZhdGFsIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRmF0YWwobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdGJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSBMb2cgbGV2ZWwgdG8gdGhlIGdpdmVuIGxvZyBsZXZlbFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBsb2dsZXZlbFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldExvZ0xldmVsKGxvZ2xldmVsKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1MnLCBsb2dsZXZlbCk7XHJcbn1cclxuXHJcbi8vIExvZyBsZXZlbHNcclxuZXhwb3J0IGNvbnN0IExvZ0xldmVsID0ge1xyXG5cdFRSQUNFOiAxLFxyXG5cdERFQlVHOiAyLFxyXG5cdElORk86IDMsXHJcblx0V0FSTklORzogNCxcclxuXHRFUlJPUjogNSxcclxufTtcclxuIiwgIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfX1xyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xyXG5cclxuLy8gRGVmaW5lcyBhIHNpbmdsZSBsaXN0ZW5lciB3aXRoIGEgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdG8gY2FsbGJhY2tcclxuXHJcbi8qKlxyXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcclxuICpcclxuICogQGNsYXNzIExpc3RlbmVyXHJcbiAqL1xyXG5jbGFzcyBMaXN0ZW5lciB7XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXHJcbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xyXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xyXG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXHJcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxyXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xyXG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcclxuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxyXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxyXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XHJcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XHJcblxyXG4gICAgLy8gR2V0IHRoZSBldmVudCBuYW1lXHJcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgd2UgaGF2ZSBhbnkgbGlzdGVuZXJzIGZvciB0aGlzIGV2ZW50XHJcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSkge1xyXG5cclxuICAgICAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lciBpbmRleGVzIHRvIGRlc3Ryb3lcclxuICAgICAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uc2xpY2UoKTtcclxuXHJcbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcclxuICAgICAgICBmb3IgKGxldCBjb3VudCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoIC0gMTsgY291bnQgPj0gMDsgY291bnQgLT0gMSkge1xyXG5cclxuICAgICAgICAgICAgLy8gR2V0IG5leHQgbGlzdGVuZXJcclxuICAgICAgICAgICAgY29uc3QgbGlzdGVuZXIgPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdW2NvdW50XTtcclxuXHJcbiAgICAgICAgICAgIGxldCBkYXRhID0gZXZlbnREYXRhLmRhdGE7XHJcblxyXG4gICAgICAgICAgICAvLyBEbyB0aGUgY2FsbGJhY2tcclxuICAgICAgICAgICAgY29uc3QgZGVzdHJveSA9IGxpc3RlbmVyLkNhbGxiYWNrKGRhdGEpO1xyXG4gICAgICAgICAgICBpZiAoZGVzdHJveSkge1xyXG4gICAgICAgICAgICAgICAgLy8gaWYgdGhlIGxpc3RlbmVyIGluZGljYXRlZCB0byBkZXN0cm95IGl0c2VsZiwgYWRkIGl0IHRvIHRoZSBkZXN0cm95IGxpc3RcclxuICAgICAgICAgICAgICAgIG5ld0V2ZW50TGlzdGVuZXJMaXN0LnNwbGljZShjb3VudCwgMSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIFVwZGF0ZSBjYWxsYmFja3Mgd2l0aCBuZXcgbGlzdCBvZiBsaXN0ZW5lcnNcclxuICAgICAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0O1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcclxuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXHJcbiAgICBsZXQgbWVzc2FnZTtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcbiAgICB9XHJcbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcclxuXHJcbiAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcclxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxyXG4gICAgfTtcclxuXHJcbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXHJcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSkge1xyXG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyc1xyXG4gICAgZGVsZXRlIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV07XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFWCcgKyBldmVudE5hbWUpO1xyXG59XHJcblxyXG4vKipcclxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXHJcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJlcyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XHJcbiAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXHJcblxyXG4gICAgaWYgKGFkZGl0aW9uYWxFdmVudE5hbWVzLmxlbmd0aCA+IDApIHtcclxuICAgICAgICBhZGRpdGlvbmFsRXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XHJcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcclxuICAgICAgICB9KVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogT2ZmIHVucmVnaXN0ZXJzIGFsbCBldmVudCBsaXN0ZW5lcnMgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cclxuICovXHJcbiBleHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmQWxsKCkge1xyXG4gICAgY29uc3QgZXZlbnROYW1lcyA9IE9iamVjdC5rZXlzKGV2ZW50TGlzdGVuZXJzKTtcclxuICAgIGZvciAobGV0IGkgPSAwOyBpICE9PSBldmVudE5hbWVzLmxlbmd0aDsgaSsrKSB7XHJcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lc1tpXSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIEV2ZW50c09uXHJcbiAqXHJcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXHJcbiAqL1xyXG4gZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcclxuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcclxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xyXG5cclxuICAgIC8vIENsZWFuIHVwIGlmIHRoZXJlIGFyZSBubyBldmVudCBsaXN0ZW5lcnMgbGVmdFxyXG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cclxuICpcclxuICogQHJldHVybnMgbnVtYmVyXHJcbiAqL1xyXG5mdW5jdGlvbiBjcnlwdG9SYW5kb20oKSB7XHJcblx0dmFyIGFycmF5ID0gbmV3IFVpbnQzMkFycmF5KDEpO1xyXG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgbnVtYmVyIHVzaW5nIGRhIG9sZC1za29vbCBNYXRoLlJhbmRvbVxyXG4gKiBJIGxpa2VzIHRvIGNhbGwgaXQgTE9MUmFuZG9tXHJcbiAqXHJcbiAqIEByZXR1cm5zIG51bWJlclxyXG4gKi9cclxuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XHJcblx0cmV0dXJuIE1hdGgucmFuZG9tKCkgKiA5MDA3MTk5MjU0NzQwOTkxO1xyXG59XHJcblxyXG4vLyBQaWNrIGEgcmFuZG9tIG51bWJlciBmdW5jdGlvbiBiYXNlZCBvbiBicm93c2VyIGNhcGFiaWxpdHlcclxudmFyIHJhbmRvbUZ1bmM7XHJcbmlmICh3aW5kb3cuY3J5cHRvKSB7XHJcblx0cmFuZG9tRnVuYyA9IGNyeXB0b1JhbmRvbTtcclxufSBlbHNlIHtcclxuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogQ2FsbCBzZW5kcyBhIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgdG8gY2FsbCB0aGUgYmluZGluZyB3aXRoIHRoZVxyXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXHJcbiAqIGJhY2tlbmQgcmVzcG9uZHMuIFRoaXMgd2lsbCBiZSByZXNvbHZlZCB3aGVuIHRoZSBjYWxsIHdhcyBzdWNjZXNzZnVsXHJcbiAqIG9yIHJlamVjdGVkIGlmIGFuIGVycm9yIGlzIHBhc3NlZCBiYWNrLlxyXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXHJcbiAqIHRpbWUgKGluIG1pbGxpc2Vjb25kcykgdGhlbiB0aGUgcHJvbWlzZSBpcyByZWplY3RlZC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZVxyXG4gKiBAcGFyYW0ge2FueT19IGFyZ3NcclxuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XHJcbiAqIEByZXR1cm5zXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XHJcblxyXG5cdC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxyXG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcclxuXHRcdHRpbWVvdXQgPSAwO1xyXG5cdH1cclxuXHJcblx0Ly8gQ3JlYXRlIGEgcHJvbWlzZVxyXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XHJcblxyXG5cdFx0Ly8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcclxuXHRcdHZhciBjYWxsYmFja0lEO1xyXG5cdFx0ZG8ge1xyXG5cdFx0XHRjYWxsYmFja0lEID0gbmFtZSArICctJyArIHJhbmRvbUZ1bmMoKTtcclxuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XHJcblxyXG5cdFx0dmFyIHRpbWVvdXRIYW5kbGU7XHJcblx0XHQvLyBTZXQgdGltZW91dFxyXG5cdFx0aWYgKHRpbWVvdXQgPiAwKSB7XHJcblx0XHRcdHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcclxuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcclxuXHRcdFx0fSwgdGltZW91dCk7XHJcblx0XHR9XHJcblxyXG5cdFx0Ly8gU3RvcmUgY2FsbGJhY2tcclxuXHRcdGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcclxuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcclxuXHRcdFx0cmVqZWN0OiByZWplY3QsXHJcblx0XHRcdHJlc29sdmU6IHJlc29sdmVcclxuXHRcdH07XHJcblxyXG5cdFx0dHJ5IHtcclxuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcclxuXHRcdFx0XHRuYW1lLFxyXG5cdFx0XHRcdGFyZ3MsXHJcblx0XHRcdFx0Y2FsbGJhY2tJRCxcclxuXHRcdFx0fTtcclxuXHJcbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcclxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKCdDJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcclxuICAgICAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGUpO1xyXG4gICAgICAgIH1cclxuICAgIH0pO1xyXG59XHJcblxyXG53aW5kb3cuT2JmdXNjYXRlZENhbGwgPSAoaWQsIGFyZ3MsIHRpbWVvdXQpID0+IHtcclxuXHJcbiAgICAvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcclxuICAgIGlmICh0aW1lb3V0ID09IG51bGwpIHtcclxuICAgICAgICB0aW1lb3V0ID0gMDtcclxuICAgIH1cclxuXHJcbiAgICAvLyBDcmVhdGUgYSBwcm9taXNlXHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xyXG5cclxuICAgICAgICAvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxyXG4gICAgICAgIHZhciBjYWxsYmFja0lEO1xyXG4gICAgICAgIGRvIHtcclxuICAgICAgICAgICAgY2FsbGJhY2tJRCA9IGlkICsgJy0nICsgcmFuZG9tRnVuYygpO1xyXG4gICAgICAgIH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XHJcblxyXG4gICAgICAgIHZhciB0aW1lb3V0SGFuZGxlO1xyXG4gICAgICAgIC8vIFNldCB0aW1lb3V0XHJcbiAgICAgICAgaWYgKHRpbWVvdXQgPiAwKSB7XHJcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcignQ2FsbCB0byBtZXRob2QgJyArIGlkICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XHJcbiAgICAgICAgICAgIH0sIHRpbWVvdXQpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gU3RvcmUgY2FsbGJhY2tcclxuICAgICAgICBjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XHJcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXHJcbiAgICAgICAgICAgIHJlamVjdDogcmVqZWN0LFxyXG4gICAgICAgICAgICByZXNvbHZlOiByZXNvbHZlXHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcclxuXHRcdFx0XHRpZCxcclxuXHRcdFx0XHRhcmdzLFxyXG5cdFx0XHRcdGNhbGxiYWNrSUQsXHJcblx0XHRcdH07XHJcblxyXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXHJcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnYycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcclxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcclxuICAgICAgICB9XHJcbiAgICB9KTtcclxufTtcclxuXHJcblxyXG4vKipcclxuICogQ2FsbGVkIGJ5IHRoZSBiYWNrZW5kIHRvIHJldHVybiBkYXRhIHRvIGEgcHJldmlvdXNseSBjYWxsZWRcclxuICogYmluZGluZyBpbnZvY2F0aW9uXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGluY29taW5nTWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGxiYWNrKGluY29taW5nTWVzc2FnZSkge1xyXG5cdC8vIFBhcnNlIHRoZSBtZXNzYWdlXHJcblx0bGV0IG1lc3NhZ2U7XHJcblx0dHJ5IHtcclxuXHRcdG1lc3NhZ2UgPSBKU09OLnBhcnNlKGluY29taW5nTWVzc2FnZSk7XHJcblx0fSBjYXRjaCAoZSkge1xyXG5cdFx0Y29uc3QgZXJyb3IgPSBgSW52YWxpZCBKU09OIHBhc3NlZCB0byBjYWxsYmFjazogJHtlLm1lc3NhZ2V9LiBNZXNzYWdlOiAke2luY29taW5nTWVzc2FnZX1gO1xyXG5cdFx0cnVudGltZS5Mb2dEZWJ1ZyhlcnJvcik7XHJcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xyXG5cdH1cclxuXHRsZXQgY2FsbGJhY2tJRCA9IG1lc3NhZ2UuY2FsbGJhY2tpZDtcclxuXHRsZXQgY2FsbGJhY2tEYXRhID0gY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xyXG5cdGlmICghY2FsbGJhY2tEYXRhKSB7XHJcblx0XHRjb25zdCBlcnJvciA9IGBDYWxsYmFjayAnJHtjYWxsYmFja0lEfScgbm90IHJlZ2lzdGVyZWQhISFgO1xyXG5cdFx0Y29uc29sZS5lcnJvcihlcnJvcik7IC8vIGVzbGludC1kaXNhYmxlLWxpbmVcclxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcblx0fVxyXG5cdGNsZWFyVGltZW91dChjYWxsYmFja0RhdGEudGltZW91dEhhbmRsZSk7XHJcblxyXG5cdGRlbGV0ZSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XHJcblxyXG5cdGlmIChtZXNzYWdlLmVycm9yKSB7XHJcblx0XHRjYWxsYmFja0RhdGEucmVqZWN0KG1lc3NhZ2UuZXJyb3IpO1xyXG5cdH0gZWxzZSB7XHJcblx0XHRjYWxsYmFja0RhdGEucmVzb2x2ZShtZXNzYWdlLnJlc3VsdCk7XHJcblx0fVxyXG59XHJcbiIsICIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX18gICAgXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gICkgXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fLyAgXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG5pbXBvcnQge0NhbGx9IGZyb20gJy4vY2FsbHMnO1xyXG5cclxuLy8gVGhpcyBpcyB3aGVyZSB3ZSBiaW5kIGdvIG1ldGhvZCB3cmFwcGVyc1xyXG53aW5kb3cuZ28gPSB7fTtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRCaW5kaW5ncyhiaW5kaW5nc01hcCkge1xyXG5cdHRyeSB7XHJcblx0XHRiaW5kaW5nc01hcCA9IEpTT04ucGFyc2UoYmluZGluZ3NNYXApO1xyXG5cdH0gY2F0Y2ggKGUpIHtcclxuXHRcdGNvbnNvbGUuZXJyb3IoZSk7XHJcblx0fVxyXG5cclxuXHQvLyBJbml0aWFsaXNlIHRoZSBiaW5kaW5ncyBtYXBcclxuXHR3aW5kb3cuZ28gPSB3aW5kb3cuZ28gfHwge307XHJcblxyXG5cdC8vIEl0ZXJhdGUgcGFja2FnZSBuYW1lc1xyXG5cdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwKS5mb3JFYWNoKChwYWNrYWdlTmFtZSkgPT4ge1xyXG5cclxuXHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxyXG5cdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gfHwge307XHJcblxyXG5cdFx0Ly8gSXRlcmF0ZSBzdHJ1Y3QgbmFtZXNcclxuXHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXSkuZm9yRWFjaCgoc3RydWN0TmFtZSkgPT4ge1xyXG5cclxuXHRcdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XHJcblx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdIHx8IHt9O1xyXG5cclxuXHRcdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdKS5mb3JFYWNoKChtZXRob2ROYW1lKSA9PiB7XHJcblxyXG5cdFx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV1bbWV0aG9kTmFtZV0gPSBmdW5jdGlvbiAoKSB7XHJcblxyXG5cdFx0XHRcdFx0Ly8gTm8gdGltZW91dCBieSBkZWZhdWx0XHJcblx0XHRcdFx0XHRsZXQgdGltZW91dCA9IDA7XHJcblxyXG5cdFx0XHRcdFx0Ly8gQWN0dWFsIGZ1bmN0aW9uXHJcblx0XHRcdFx0XHRmdW5jdGlvbiBkeW5hbWljKCkge1xyXG5cdFx0XHRcdFx0XHRjb25zdCBhcmdzID0gW10uc2xpY2UuY2FsbChhcmd1bWVudHMpO1xyXG5cdFx0XHRcdFx0XHRyZXR1cm4gQ2FsbChbcGFja2FnZU5hbWUsIHN0cnVjdE5hbWUsIG1ldGhvZE5hbWVdLmpvaW4oJy4nKSwgYXJncywgdGltZW91dCk7XHJcblx0XHRcdFx0XHR9XHJcblxyXG5cdFx0XHRcdFx0Ly8gQWxsb3cgc2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXHJcblx0XHRcdFx0XHRkeW5hbWljLnNldFRpbWVvdXQgPSBmdW5jdGlvbiAobmV3VGltZW91dCkge1xyXG5cdFx0XHRcdFx0XHR0aW1lb3V0ID0gbmV3VGltZW91dDtcclxuXHRcdFx0XHRcdH07XHJcblxyXG5cdFx0XHRcdFx0Ly8gQWxsb3cgZ2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXHJcblx0XHRcdFx0XHRkeW5hbWljLmdldFRpbWVvdXQgPSBmdW5jdGlvbiAoKSB7XHJcblx0XHRcdFx0XHRcdHJldHVybiB0aW1lb3V0O1xyXG5cdFx0XHRcdFx0fTtcclxuXHJcblx0XHRcdFx0XHRyZXR1cm4gZHluYW1pYztcclxuXHRcdFx0XHR9KCk7XHJcblx0XHRcdH0pO1xyXG5cdFx0fSk7XHJcblx0fSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcblxyXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93UmVsb2FkKCkge1xyXG4gICAgd2luZG93LmxvY2F0aW9uLnJlbG9hZCgpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93UmVsb2FkQXBwKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUicpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U3lzdGVtRGVmYXVsdFRoZW1lKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVNEVCcpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TGlnaHRUaGVtZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FMVCcpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0RGFya1RoZW1lKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQURUJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBQbGFjZSB0aGUgd2luZG93IGluIHRoZSBjZW50ZXIgb2YgdGhlIHNjcmVlblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93Q2VudGVyKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXYycpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgd2luZG93IHRpdGxlXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZVxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0VGl0bGUodGl0bGUpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1QnICsgdGl0bGUpO1xyXG59XHJcblxyXG4vKipcclxuICogTWFrZXMgdGhlIHdpbmRvdyBnbyBmdWxsc2NyZWVuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dGdWxsc2NyZWVuKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXRicpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV2ZXJ0cyB0aGUgd2luZG93IGZyb20gZnVsbHNjcmVlblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5mdWxsc2NyZWVuKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXZicpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgaW4gZnVsbCBzY3JlZW4gbW9kZSBvciBub3QuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc0Z1bGxzY3JlZW4oKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc0Z1bGxzY3JlZW5cIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dzOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8e3c6IG51bWJlciwgaDogbnVtYmVyfT59IFRoZSBzaXplIG9mIHRoZSB3aW5kb3dcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0U2l6ZSgpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFNpemVcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNYXhTaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1pblNpemUod2lkdGgsIGhlaWdodCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXejonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xyXG59XHJcblxyXG5cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIHdpbmRvdyBBbHdheXNPblRvcCBvciBub3Qgb24gdG9wXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRBbHdheXNPblRvcChiKSB7XHJcblxyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVRQOicgKyAoYiA/ICcxJyA6ICcwJykpO1xyXG59XHJcblxyXG5cclxuXHJcblxyXG4vKipcclxuICogU2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IHhcclxuICogQHBhcmFtIHtudW1iZXJ9IHlcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRQb3NpdGlvbih4LCB5KSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dwOicgKyB4ICsgJzonICsgeSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHt4OiBudW1iZXIsIHk6IG51bWJlcn0+fSBUaGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFBvc2l0aW9uKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0UG9zXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogSGlkZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dIaWRlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXSCcpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvdyB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTaG93KCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUycpO1xyXG59XHJcblxyXG4vKipcclxuICogTWF4aW1pc2UgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWF4aW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dNJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUb2dnbGUgdGhlIE1heGltaXNlIG9mIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1RvZ2dsZU1heGltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdCcpO1xyXG59XHJcblxyXG4vKipcclxuICogVW5tYXhpbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1heGltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVScpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkIG9yIG5vdC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWF4aW1pc2VkKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNYXhpbWlzZWRcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNaW5pbWlzZSB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNaW5pbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV20nKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFVubWluaW1pc2UgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5taW5pbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3UnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1pbmltaXNlZCBvciBub3QuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc01pbmltaXNlZCgpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWluaW1pc2VkXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbm9ybWFsIG9yIG5vdC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTm9ybWFsKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNOb3JtYWxcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IFIgUmVkXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBHIEdyZWVuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBCIEJsdWVcclxuICogQHBhcmFtIHtudW1iZXJ9IEEgQWxwaGFcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRCYWNrZ3JvdW5kQ29sb3VyKFIsIEcsIEIsIEEpIHtcclxuICAgIGxldCByZ2JhID0gSlNPTi5zdHJpbmdpZnkoe3I6IFIgfHwgMCwgZzogRyB8fCAwLCBiOiBCIHx8IDAsIGE6IEEgfHwgMjU1fSk7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dyOicgKyByZ2JhKTtcclxufVxyXG5cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcblxyXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XHJcblxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIGFsbCBzY3JlZW5zLiBDYWxsIHRoaXMgYW5ldyBlYWNoIHRpbWUgeW91IHdhbnQgdG8gcmVmcmVzaCBkYXRhIGZyb20gdGhlIHVuZGVybHlpbmcgd2luZG93aW5nIHN5c3RlbS5cclxuICogQGV4cG9ydFxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KCcuLi93cmFwcGVyL3J1bnRpbWUnKS5TY3JlZW59IFNjcmVlblxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHtTY3JlZW5bXX0+fSBUaGUgc2NyZWVuc1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNjcmVlbkdldEFsbCgpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlNjcmVlbkdldEFsbFwiKTtcclxufVxyXG4iLCAiLyoqXHJcbiAqIEBkZXNjcmlwdGlvbjogVXNlIHRoZSBzeXN0ZW0gZGVmYXVsdCBicm93c2VyIHRvIG9wZW4gdGhlIHVybFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIFxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJyb3dzZXJPcGVuVVJMKHVybCkge1xyXG4gIHdpbmRvdy5XYWlsc0ludm9rZSgnQk86JyArIHVybCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IHRleHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRTZXRUZXh0KHRleHQpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsaXBib2FyZFNldFRleHRcIiwgW3RleHRdKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgdGV4dCBjb250ZW50IG9mIHRoZSBjbGlwYm9hcmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHtzdHJpbmd9Pn0gVGV4dCBjb250ZW50IG9mIHRoZSBjbGlwYm9hcmRcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2xpcGJvYXJkR2V0VGV4dCgpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsaXBib2FyZEdldFRleHRcIik7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtFdmVudHNPbn0gZnJvbSBcIi4vZXZlbnRzXCI7XHJcblxyXG4vKipcclxuICogcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgY2hlY2tzIHRoZSBicm93c2VyJ3MgY2FwYWJpbGl0eSBvZiBzZW5kaW5nIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufVxyXG4gKiBAY29uc3RydWN0b3JcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYW5SZXNvbHZlRmlsZVBhdGhzKCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzICE9IG51bGw7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBSZXNvbHZlRmlsZVBhdGhzKGZpbGVzKSB7XHJcbiAgICAvLyBPbmx5IGZvciB3aW5kb3dzIHdlYnZpZXcyID49IDEuMC4xNzc0LjMwXHJcbiAgICAvLyBodHRwczovL2xlYXJuLm1pY3Jvc29mdC5jb20vZW4tdXMvbWljcm9zb2Z0LWVkZ2Uvd2VidmlldzIvcmVmZXJlbmNlL3dpbjMyL2ljb3Jld2VidmlldzJ3ZWJtZXNzYWdlcmVjZWl2ZWRldmVudGFyZ3MyP3ZpZXc9d2VidmlldzItMS4wLjE4MjMuMzIjYXBwbGllcy10b1xyXG4gICAgaWYgKCF3aW5kb3cuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cykge1xyXG4gICAgICAgIHJlamVjdChuZXcgRXJyb3IoXCJVbnN1cHBvcnRlZCBQbGF0Zm9ybVwiKSk7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKGBmaWxlOmRyb3A6YCwgZmlsZXMpO1xyXG59XHJcbiIsICIvKlxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuXHJcblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXHJcbiovXHJcbmV4cG9ydCBmdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KSB7XHJcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICBjb25zdCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xyXG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xyXG4gICAgY29uc3QgZGVmYXVsdENvbnRleHRNZW51QWN0aW9uID0gY29tcHV0ZWRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKTtcclxuICAgIHN3aXRjaCAoZGVmYXVsdENvbnRleHRNZW51QWN0aW9uKSB7XHJcbiAgICAgICAgY2FzZSBcInNob3dcIjpcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGNhc2UgXCJoaWRlXCI6XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICBkZWZhdWx0OlxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxyXG4gICAgICAgICAgICBpZiAoZWxlbWVudC5pc0NvbnRlbnRFZGl0YWJsZSkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkIGFuZCBhY3Rpb24gaXMgb24gdGhlIHNlbGVjdGVkIGVsZW1lbnRzXHJcbiAgICAgICAgICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcclxuICAgICAgICAgICAgY29uc3QgaGFzU2VsZWN0aW9uID0gKHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDApXHJcbiAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcclxuICAgICAgICAgICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xyXG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xyXG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gZWxlbWVudCkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRhZ25hbWUgaXMgaW5wdXQgb3IgdGV4dGFyZWFcclxuICAgICAgICAgICAgaWYgKGVsZW1lbnQudGFnTmFtZSA9PT0gXCJJTlBVVFwiIHx8IGVsZW1lbnQudGFnTmFtZSA9PT0gXCJURVhUQVJFQVwiKSB7XHJcbiAgICAgICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghZWxlbWVudC5yZWFkT25seSAmJiAhZWxlbWVudC5kaXNhYmxlZCkpIHtcclxuICAgICAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XHJcbmltcG9ydCB7ZXZlbnRMaXN0ZW5lcnMsIEV2ZW50c0VtaXQsIEV2ZW50c05vdGlmeSwgRXZlbnRzT2ZmLCBFdmVudHNPbiwgRXZlbnRzT25jZSwgRXZlbnRzT25NdWx0aXBsZX0gZnJvbSAnLi9ldmVudHMnO1xyXG5pbXBvcnQge0NhbGwsIENhbGxiYWNrLCBjYWxsYmFja3N9IGZyb20gJy4vY2FsbHMnO1xyXG5pbXBvcnQge1NldEJpbmRpbmdzfSBmcm9tIFwiLi9iaW5kaW5nc1wiO1xyXG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XHJcbmltcG9ydCAqIGFzIFNjcmVlbiBmcm9tIFwiLi9zY3JlZW5cIjtcclxuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyXCI7XHJcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmRcIjtcclxuaW1wb3J0ICogYXMgRHJhZ0FuZERyb3AgZnJvbSBcIi4vZHJhZ2FuZGRyb3BcIjtcclxuaW1wb3J0ICogYXMgQ29udGV4dE1lbnUgZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcclxuXHJcblxyXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUScpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUycpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnSCcpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpFbnZpcm9ubWVudFwiKTtcclxufVxyXG5cclxuLy8gVGhlIEpTIHJ1bnRpbWVcclxud2luZG93LnJ1bnRpbWUgPSB7XHJcbiAgICAuLi5Mb2csXHJcbiAgICAuLi5XaW5kb3csXHJcbiAgICAuLi5Ccm93c2VyLFxyXG4gICAgLi4uU2NyZWVuLFxyXG4gICAgLi4uQ2xpcGJvYXJkLFxyXG4gICAgLi4uRHJhZ0FuZERyb3AsXHJcbiAgICBFdmVudHNPbixcclxuICAgIEV2ZW50c09uY2UsXHJcbiAgICBFdmVudHNPbk11bHRpcGxlLFxyXG4gICAgRXZlbnRzRW1pdCxcclxuICAgIEV2ZW50c09mZixcclxuICAgIEVudmlyb25tZW50LFxyXG4gICAgU2hvdyxcclxuICAgIEhpZGUsXHJcbiAgICBRdWl0XHJcbn07XHJcblxyXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcclxud2luZG93LndhaWxzID0ge1xyXG4gICAgQ2FsbGJhY2ssXHJcbiAgICBFdmVudHNOb3RpZnksXHJcbiAgICBTZXRCaW5kaW5ncyxcclxuICAgIGV2ZW50TGlzdGVuZXJzLFxyXG4gICAgY2FsbGJhY2tzLFxyXG4gICAgZmxhZ3M6IHtcclxuICAgICAgICBkaXNhYmxlU2Nyb2xsYmFyRHJhZzogZmFsc2UsXHJcbiAgICAgICAgZGlzYWJsZURlZmF1bHRDb250ZXh0TWVudTogZmFsc2UsXHJcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcclxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxyXG4gICAgICAgIGJvcmRlclRoaWNrbmVzczogNixcclxuICAgICAgICBzaG91bGREcmFnOiBmYWxzZSxcclxuICAgICAgICBkZWZlckRyYWdUb01vdXNlTW92ZTogdHJ1ZSxcclxuICAgICAgICBjc3NEcmFnUHJvcGVydHk6IFwiLS13YWlscy1kcmFnZ2FibGVcIixcclxuICAgICAgICBjc3NEcmFnVmFsdWU6IFwiZHJhZ1wiLFxyXG4gICAgICAgIGNzc0Ryb3BQcm9wZXJ0eTogXCItLXdhaWxzLWRyb3AtdGFyZ2V0XCIsXHJcbiAgICAgICAgY3NzRHJvcFZhbHVlOiBcImRyb3BcIixcclxuICAgICAgICBlbmFibGVXYWlsc0RyYWdBbmREcm9wOiBmYWxzZSxcclxuICAgICAgICB3YWlsc0Ryb3BQcmV2aW91c0VsZW1lbnQ6IG51bGwsXHJcbiAgICB9XHJcbn07XHJcblxyXG4vLyBTZXQgdGhlIGJpbmRpbmdzXHJcbmlmICh3aW5kb3cud2FpbHNiaW5kaW5ncykge1xyXG4gICAgd2luZG93LndhaWxzLlNldEJpbmRpbmdzKHdpbmRvdy53YWlsc2JpbmRpbmdzKTtcclxuICAgIGRlbGV0ZSB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3M7XHJcbn1cclxuXHJcbi8vIChib29sKSBUaGlzIGlzIGV2YWx1YXRlZCBhdCBidWlsZCB0aW1lIGluIHBhY2thZ2UuanNvblxyXG5pZiAoIURFQlVHKSB7XHJcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XHJcbn1cclxuXHJcbmxldCBkcmFnVGVzdCA9IGZ1bmN0aW9uIChlKSB7XHJcbiAgICB2YXIgdmFsID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdQcm9wZXJ0eSk7XHJcbiAgICBpZiAodmFsKSB7XHJcbiAgICAgIHZhbCA9IHZhbC50cmltKCk7XHJcbiAgICB9XHJcbiAgICBcclxuICAgIGlmICh2YWwgIT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGUuYnV0dG9ucyAhPT0gMSkge1xyXG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBub3QgdGhlIHByaW1hcnkgYnV0dG9uIGhhcyBiZWVuIGNsaWNrZWQuXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChlLmRldGFpbCAhPT0gMSkge1xyXG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBtb3JlIHRoYW4gb25jZSBoYXMgYmVlbiBjbGlja2VkLCBlLmcuIHdoZW4gZG91YmxlIGNsaWNraW5nXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiB0cnVlO1xyXG59O1xyXG5cclxud2luZG93LndhaWxzLnNldENTU0RyYWdQcm9wZXJ0aWVzID0gZnVuY3Rpb24gKHByb3BlcnR5LCB2YWx1ZSkge1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdQcm9wZXJ0eSA9IHByb3BlcnR5O1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdWYWx1ZSA9IHZhbHVlO1xyXG59XHJcblxyXG53aW5kb3cud2FpbHMuc2V0Q1NTRHJvcFByb3BlcnRpZXMgPSBmdW5jdGlvbiAocHJvcGVydHksIHZhbHVlKSB7XHJcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFByb3BlcnR5ID0gcHJvcGVydHk7XHJcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlID0gdmFsdWU7XHJcbn1cclxuXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xyXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nXHJcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJyZXNpemU6XCIgKyB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSk7XHJcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoZHJhZ1Rlc3QoZSkpIHtcclxuICAgICAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRpc2FibGVTY3JvbGxiYXJEcmFnKSB7XHJcbiAgICAgICAgICAgIC8vIFRoaXMgY2hlY2tzIGZvciBjbGlja3Mgb24gdGhlIHNjcm9sbCBiYXJcclxuICAgICAgICAgICAgaWYgKGUub2Zmc2V0WCA+IGUudGFyZ2V0LmNsaWVudFdpZHRoIHx8IGUub2Zmc2V0WSA+IGUudGFyZ2V0LmNsaWVudEhlaWdodCkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmZXJEcmFnVG9Nb3VzZU1vdmUpIHtcclxuICAgICAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSB0cnVlO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIGUucHJldmVudERlZmF1bHQoKVxyXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJkcmFnXCIpO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XHJcbiAgICB9XHJcbn0pO1xyXG5cclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCAoKSA9PiB7XHJcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG59KTtcclxuXHJcbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IpIHtcclxuICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBjdXJzb3IgfHwgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3I7XHJcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSA9IGN1cnNvcjtcclxufVxyXG5cclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIGZ1bmN0aW9uIChlKSB7XHJcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcpIHtcclxuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XHJcbiAgICAgICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcclxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVJlc2l6ZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XHJcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yO1xyXG4gICAgfVxyXG4gICAgaWYgKHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcyAmJiB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzKSB7XHJcbiAgICAgICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IFwic2UtcmVzaXplXCI7XHJcbiAgICB9XHJcbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xyXG5cclxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XHJcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9IGVsc2UgaWYgKHJpZ2h0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgdG9wQm9yZGVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIgJiYgcmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XHJcblxyXG59KTtcclxuXHJcbi8vIFNldHVwIGNvbnRleHQgbWVudSBob29rXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGZ1bmN0aW9uIChlKSB7XHJcbiAgICAvLyBhbHdheXMgc2hvdyB0aGUgY29udGV4dG1lbnUgaW4gZGVidWcgJiBkZXZcclxuICAgIGlmIChERUJVRykgcmV0dXJuO1xyXG5cclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZURlZmF1bHRDb250ZXh0TWVudSkge1xyXG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgQ29udGV4dE1lbnUucHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShlKTtcclxuICAgIH1cclxufSk7XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBmdW5jdGlvbiAoZSkge1xyXG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlV2FpbHNEcmFnQW5kRHJvcCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgIGxldCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChlLngsIGUueSk7XHJcbiAgICBpZiAodGFyZ2V0RWxlbWVudCA9PT0gd2luZG93LndhaWxzLmZsYWdzLndhaWxzRHJvcFByZXZpb3VzRWxlbWVudCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGNvbnN0IHN0eWxlID0gdGFyZ2V0RWxlbWVudC5zdHlsZTtcclxuICAgIGxldCBjc3NEcm9wVmFsdWUgPSBudWxsO1xyXG4gICAgaWYgKE9iamVjdC5rZXlzKHN0eWxlKS5maW5kSW5kZXgoa2V5ID0+IHN0eWxlW2tleV0gPT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkpIDwgMCkge1xyXG4gICAgICAgIHRhcmdldEVsZW1lbnQgPSB0YXJnZXRFbGVtZW50LmNsb3Nlc3QoYFtzdHlsZSo9JyR7d2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eX0nXWApO1xyXG4gICAgfVxyXG4gICAgaWYgKHRhcmdldEVsZW1lbnQgPT0gbnVsbCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGNzc0Ryb3BWYWx1ZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldEVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSk7XHJcbiAgICBpZiAoY3NzRHJvcFZhbHVlKSB7XHJcbiAgICAgICAgY3NzRHJvcFZhbHVlID0gY3NzRHJvcFZhbHVlLnRyaW0oKTtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoY3NzRHJvcFZhbHVlID09PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlKSB7XHJcbiAgICAgICAgdGFyZ2V0RWxlbWVudC5jbGFzc0xpc3QuYWRkKFwid2FpbHMtZHJvcC10YXJnZXQtYWN0aXZlXCIpO1xyXG4gICAgfSBlbHNlIGlmICh3aW5kb3cud2FpbHMuZmxhZ3Mud2FpbHNEcm9wUHJldmlvdXNFbGVtZW50KSB7XHJcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLndhaWxzRHJvcFByZXZpb3VzRWxlbWVudC5jbGFzc0xpc3QucmVtb3ZlKFwid2FpbHMtZHJvcC10YXJnZXQtYWN0aXZlXCIpO1xyXG4gICAgfVxyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLndhaWxzRHJvcFByZXZpb3VzRWxlbWVudCA9IHRhcmdldEVsZW1lbnQ7XHJcbn0pXHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgZnVuY3Rpb24gKGUpIHtcclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3ApIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcbiAgICBlLnByZXZlbnREZWZhdWx0KCk7XHJcblxyXG4gICAgbGV0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGUueCwgZS55KTtcclxuICAgIGxldCBjc3NEcm9wVmFsdWUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXRFbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkpO1xyXG4gICAgaWYgKGNzc0Ryb3BWYWx1ZSkge1xyXG4gICAgICAgIGNzc0Ryb3BWYWx1ZSA9IGNzc0Ryb3BWYWx1ZS50cmltKCk7XHJcbiAgICB9XHJcbiAgICBpZiAoY3NzRHJvcFZhbHVlICE9PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlICYmIHdpbmRvdy53YWlscy5mbGFncy53YWlsc0Ryb3BQcmV2aW91c0VsZW1lbnQpIHtcclxuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Mud2FpbHNEcm9wUHJldmlvdXNFbGVtZW50LmNsYXNzTGlzdC5yZW1vdmUoXCJ3YWlscy1kcm9wLXRhcmdldC1hY3RpdmVcIik7XHJcbiAgICB9XHJcbn0pO1xyXG5cclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCBmdW5jdGlvbiAoZSkge1xyXG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlV2FpbHNEcmFnQW5kRHJvcCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgIGxldCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChlLngsIGUueSk7XHJcbiAgICBsZXQgY3NzRHJvcFZhbHVlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0RWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZSh3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFByb3BlcnR5KTtcclxuICAgIGlmIChjc3NEcm9wVmFsdWUpIHtcclxuICAgICAgICBjc3NEcm9wVmFsdWUgPSBjc3NEcm9wVmFsdWUudHJpbSgpO1xyXG4gICAgfVxyXG4gICAgaWYgKGNzc0Ryb3BWYWx1ZSAhPT0gd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BWYWx1ZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIC8vIHByb2Nlc3MgZmlsZXNcclxuICAgIGxldCBmaWxlcyA9IFtdO1xyXG4gICAgaWYgKGUuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XHJcbiAgICAgICAgZmlsZXMgPSBbLi4uZS5kYXRhVHJhbnNmZXIuaXRlbXNdLm1hcCgoaXRlbSwgaSkgPT4ge1xyXG4gICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcclxuICAgICAgICAgICAgICAgIGNvbnN0IGZpbGUgPSBpdGVtLmdldEFzRmlsZSgpO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuIGZpbGU7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgZmlsZXMgPSBbLi4uZS5kYXRhVHJhbnNmZXIuZmlsZXNdO1xyXG4gICAgfVxyXG5cclxuICAgIHdpbmRvdy5ydW50aW1lLlJlc29sdmVGaWxlUGF0aHMoZmlsZXMpO1xyXG4gICAgaWYod2luZG93LndhaWxzLmZsYWdzLndhaWxzRHJvcFByZXZpb3VzRWxlbWVudCkge1xyXG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy53YWlsc0Ryb3BQcmV2aW91c0VsZW1lbnQuY2xhc3NMaXN0LnJlbW92ZShcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcclxuICAgIH1cclxufSk7XHJcblxyXG53aW5kb3cuV2FpbHNJbnZva2UoXCJydW50aW1lOnJlYWR5XCIpOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBa0JBLFdBQVMsZUFBZSxPQUFPLFNBQVM7QUFJdkMsV0FBTyxZQUFZLE1BQU0sUUFBUSxPQUFPO0FBQUEsRUFDekM7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFFBQVEsU0FBUztBQUNoQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsV0FBVyxTQUFTO0FBQ25DLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsWUFBWSxVQUFVO0FBQ3JDLG1CQUFlLEtBQUssUUFBUTtBQUFBLEVBQzdCO0FBR08sTUFBTSxXQUFXO0FBQUEsSUFDdkIsT0FBTztBQUFBLElBQ1AsT0FBTztBQUFBLElBQ1AsTUFBTTtBQUFBLElBQ04sU0FBUztBQUFBLElBQ1QsT0FBTztBQUFBLEVBQ1I7OztBQzlGQSxNQUFNLFdBQU4sTUFBZTtBQUFBLElBUVgsWUFBWSxXQUFXLFVBQVUsY0FBYztBQUMzQyxX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
