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
  function EventsOff2(eventName, ...additionalEventNames) {
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

  // desktop/draganddrop.js
  var draganddrop_exports = {};
  __export(draganddrop_exports, {
    CanResolveFilePaths: () => CanResolveFilePaths,
    DragAndDropOff: () => DragAndDropOff,
    DragAndDropOn: () => DragAndDropOn,
    ResolveFilePaths: () => ResolveFilePaths
  });
  var flags = {
    registered: false,
    defaultUseDropTarget: true,
    useDropTarget: true,
    prevElement: null
  };
  function onDragOver(e) {
    e.preventDefault();
    e.stopPropagation();
    if (!flags.useDropTarget) {
      return;
    }
    let targetElement = document.elementFromPoint(e.x, e.y);
    if (targetElement === flags.prevElement) {
      return;
    }
    const style = targetElement.style;
    let cssDropValue = null;
    if (Object.keys(style).findIndex((key) => style[key] === window.wails.flags.cssDropProperty) < 0) {
      targetElement = targetElement.closest(`[style*='${window.wails.flags.cssDropProperty}']`);
    }
    if (targetElement === null) {
      return;
    }
    cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
      cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue === window.wails.flags.cssDropValue) {
      targetElement.classList.add("wails-drop-target-active");
    } else if (flags.prevElement) {
      targetElement.classList.remove("wails-drop-target-active");
      flags.prevElement.classList.remove("wails-drop-target-active");
    }
    flags.prevElement = targetElement;
  }
  function onDragLeave(e) {
    e.preventDefault();
    e.stopPropagation();
    if (!flags.useDropTarget) {
      return;
    }
    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
      cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue && flags.prevElement) {
      targetElement.classList.remove("wails-drop-target-active");
      flags.prevElement.classList.remove("wails-drop-target-active");
    }
  }
  function onDrop(e) {
    e.preventDefault();
    e.stopPropagation();
    if (!flags.useDropTarget) {
      return;
    }
    let targetElement = document.elementFromPoint(e.x, e.y);
    let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
    if (cssDropValue) {
      cssDropValue = cssDropValue.trim();
    }
    if (cssDropValue !== window.wails.flags.cssDropValue) {
      if (flags.prevElement) {
        targetElement.classList.remove("wails-drop-target-active");
        flags.prevElement.classList.remove("wails-drop-target-active");
      }
      return;
    }
    if (CanResolveFilePaths()) {
      let files = [];
      if (e.dataTransfer.items) {
        files = [...e.dataTransfer.items].map((item, i) => {
          if (item.kind === "file") {
            return item.getAsFile();
          }
        });
      } else {
        files = [...e.dataTransfer.files];
      }
      window.runtime.ResolveFilePaths(e.x, e.y, files);
    }
    if (flags.prevElement) {
      flags.prevElement.classList.remove("wails-drop-target-active");
    }
  }
  function CanResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
  }
  function ResolveFilePaths(x, y, files) {
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
      chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
      return;
    }
    console.warn("unsupported platform");
  }
  function DragAndDropOn(callback, useDropTarget) {
    if (!window.wails.flags.enableWailsDragAndDrop || typeof callback !== "function") {
      return;
    }
    if (flags.registered) {
      return;
    }
    flags.registered = true;
    const uDTPT = typeof useDropTarget;
    flags.useDropTarget = uDTPT === "undefined" || uDTPT !== "boolean" ? flags.defaultUseDropTarget : useDropTarget;
    window.addEventListener("dragover", onDragOver);
    window.addEventListener("dragleave", onDragLeave);
    window.addEventListener("drop", onDrop);
    let cb = callback;
    if (flags.useDropTarget) {
      cb = function(x, y, paths) {
        let targetElement = document.elementFromPoint(x, y);
        if (!targetElement) {
          return;
        }
        let cssDropValue = window.getComputedStyle(targetElement).getPropertyValue(window.wails.flags.cssDropProperty);
        if (cssDropValue) {
          cssDropValue = cssDropValue.trim();
        }
        if (cssDropValue !== window.wails.flags.cssDropValue) {
          return;
        }
        callback(x, y, paths);
      };
    }
    EventsOn("wails:dnd:drop", cb);
  }
  function DragAndDropOff() {
    window.removeEventListener("dragover", onDragOver);
    window.removeEventListener("dragleave", onDragLeave);
    window.removeEventListener("drop", onDrop);
    EventsOff("wails:dnd:drop");
    flags.registered = false;
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
    EventsOff: EventsOff2,
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
  window.WailsInvoke("runtime:ready");
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL2RyYWdhbmRkcm9wLmpzIiwgImRlc2t0b3AvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9tYWluLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG4vKipcbiAqIFNlbmRzIGEgbG9nIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgd2l0aCB0aGUgZ2l2ZW4gbGV2ZWwgKyBtZXNzYWdlXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGxldmVsXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5mdW5jdGlvbiBzZW5kTG9nTWVzc2FnZShsZXZlbCwgbWVzc2FnZSkge1xuXG5cdC8vIExvZyBNZXNzYWdlIGZvcm1hdDpcblx0Ly8gbFt0eXBlXVttZXNzYWdlXVxuXHR3aW5kb3cuV2FpbHNJbnZva2UoJ0wnICsgbGV2ZWwgKyBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIHRyYWNlIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dUcmFjZShtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdUJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nUHJpbnQobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnUCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gZGVidWcgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ0RlYnVnKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0QnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGluZm8gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ0luZm8obWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnSScsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gd2FybmluZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nV2FybmluZyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdXJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBlcnJvciBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nRXJyb3IobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnRScsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gZmF0YWwgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ0ZhdGFsKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0YnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSBMb2cgbGV2ZWwgdG8gdGhlIGdpdmVuIGxvZyBsZXZlbFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSBsb2dsZXZlbFxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0TG9nTGV2ZWwobG9nbGV2ZWwpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1MnLCBsb2dsZXZlbCk7XG59XG5cbi8vIExvZyBsZXZlbHNcbmV4cG9ydCBjb25zdCBMb2dMZXZlbCA9IHtcblx0VFJBQ0U6IDEsXG5cdERFQlVHOiAyLFxuXHRJTkZPOiAzLFxuXHRXQVJOSU5HOiA0LFxuXHRFUlJPUjogNSxcbn07XG4iLCAiLypcbiBfICAgICAgIF9fICAgICAgXyBfX1xufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbi8vIERlZmluZXMgYSBzaW5nbGUgbGlzdGVuZXIgd2l0aCBhIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRvIGNhbGxiYWNrXG5cbi8qKlxuICogVGhlIExpc3RlbmVyIGNsYXNzIGRlZmluZXMgYSBsaXN0ZW5lciEgOi0pXG4gKlxuICogQGNsYXNzIExpc3RlbmVyXG4gKi9cbmNsYXNzIExpc3RlbmVyIHtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIExpc3RlbmVyLlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcbiAgICAgICAgICAgIGNhbGxiYWNrLmFwcGx5KG51bGwsIGRhdGEpO1xuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxuICAgICAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSB7fTtcblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgYG1heENhbGxiYWNrc2AgdGltZXMgYmVmb3JlIGJlaW5nIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT24oZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbmZ1bmN0aW9uIG5vdGlmeUxpc3RlbmVycyhldmVudERhdGEpIHtcblxuICAgIC8vIEdldCB0aGUgZXZlbnQgbmFtZVxuICAgIGxldCBldmVudE5hbWUgPSBldmVudERhdGEubmFtZTtcblxuICAgIC8vIENoZWNrIGlmIHdlIGhhdmUgYW55IGxpc3RlbmVycyBmb3IgdGhpcyBldmVudFxuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdKSB7XG5cbiAgICAgICAgLy8gS2VlcCBhIGxpc3Qgb2YgbGlzdGVuZXIgaW5kZXhlcyB0byBkZXN0cm95XG4gICAgICAgIGNvbnN0IG5ld0V2ZW50TGlzdGVuZXJMaXN0ID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5zbGljZSgpO1xuXG4gICAgICAgIC8vIEl0ZXJhdGUgbGlzdGVuZXJzXG4gICAgICAgIGZvciAobGV0IGNvdW50ID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGggLSAxOyBjb3VudCA+PSAwOyBjb3VudCAtPSAxKSB7XG5cbiAgICAgICAgICAgIC8vIEdldCBuZXh0IGxpc3RlbmVyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV1bY291bnRdO1xuXG4gICAgICAgICAgICBsZXQgZGF0YSA9IGV2ZW50RGF0YS5kYXRhO1xuXG4gICAgICAgICAgICAvLyBEbyB0aGUgY2FsbGJhY2tcbiAgICAgICAgICAgIGNvbnN0IGRlc3Ryb3kgPSBsaXN0ZW5lci5DYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIGlmIChkZXN0cm95KSB7XG4gICAgICAgICAgICAgICAgLy8gaWYgdGhlIGxpc3RlbmVyIGluZGljYXRlZCB0byBkZXN0cm95IGl0c2VsZiwgYWRkIGl0IHRvIHRoZSBkZXN0cm95IGxpc3RcbiAgICAgICAgICAgICAgICBuZXdFdmVudExpc3RlbmVyTGlzdC5zcGxpY2UoY291bnQsIDEpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG5cbiAgICAgICAgLy8gVXBkYXRlIGNhbGxiYWNrcyB3aXRoIG5ldyBsaXN0IG9mIGxpc3RlbmVyc1xuICAgICAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoID09PSAwKSB7XG4gICAgICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0O1xuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcblxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcbiAgICAvLyBQYXJzZSB0aGUgbWVzc2FnZVxuICAgIGxldCBtZXNzYWdlO1xuICAgIHRyeSB7XG4gICAgICAgIG1lc3NhZ2UgPSBKU09OLnBhcnNlKG5vdGlmeU1lc3NhZ2UpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuICAgIH1cbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XG59XG5cbi8qKlxuICogRW1pdCBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBkYXRhXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcblxuICAgIGNvbnN0IHBheWxvYWQgPSB7XG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcbiAgICAgICAgZGF0YTogW10uc2xpY2UuYXBwbHkoYXJndW1lbnRzKS5zbGljZSgxKSxcbiAgICB9O1xuXG4gICAgLy8gTm90aWZ5IEpTIGxpc3RlbmVyc1xuICAgIG5vdGlmeUxpc3RlbmVycyhwYXlsb2FkKTtcblxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VFJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcbn1cblxuZnVuY3Rpb24gcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKSB7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyc1xuICAgIGRlbGV0ZSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdO1xuXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRVgnICsgZXZlbnROYW1lKTtcbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPbixcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJlcyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG5cbiAgICBpZiAoYWRkaXRpb25hbEV2ZW50TmFtZXMubGVuZ3RoID4gMCkge1xuICAgICAgICBhZGRpdGlvbmFsRXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XG4gICAgICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG4gICAgICAgIH0pXG4gICAgfVxufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhbGwgZXZlbnQgbGlzdGVuZXJzIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uXG4gKi9cbiBleHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmQWxsKCkge1xuICAgIGNvbnN0IGV2ZW50TmFtZXMgPSBPYmplY3Qua2V5cyhldmVudExpc3RlbmVycyk7XG4gICAgZm9yIChsZXQgaSA9IDA7IGkgIT09IGV2ZW50TmFtZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lc1tpXSk7XG4gICAgfVxufVxuXG4vKipcbiAqIGxpc3RlbmVyT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggRXZlbnRzT25cbiAqXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxuICovXG4gZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuXG4gICAgLy8gQ2xlYW4gdXAgaWYgdGhlcmUgYXJlIG5vIGV2ZW50IGxpc3RlbmVycyBsZWZ0XG4gICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XG4gICAgfVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gY3J5cHRvUmFuZG9tKCkge1xuXHR2YXIgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciB1c2luZyBkYSBvbGQtc2tvb2wgTWF0aC5SYW5kb21cbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XG5cdHJldHVybiBNYXRoLnJhbmRvbSgpICogOTAwNzE5OTI1NDc0MDk5MTtcbn1cblxuLy8gUGljayBhIHJhbmRvbSBudW1iZXIgZnVuY3Rpb24gYmFzZWQgb24gYnJvd3NlciBjYXBhYmlsaXR5XG52YXIgcmFuZG9tRnVuYztcbmlmICh3aW5kb3cuY3J5cHRvKSB7XG5cdHJhbmRvbUZ1bmMgPSBjcnlwdG9SYW5kb207XG59IGVsc2Uge1xuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XG59XG5cblxuLyoqXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXG4gKiBiYWNrZW5kIHJlc3BvbmRzLiBUaGlzIHdpbGwgYmUgcmVzb2x2ZWQgd2hlbiB0aGUgY2FsbCB3YXMgc3VjY2Vzc2Z1bFxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXG4gKiB0aW1lIChpbiBtaWxsaXNlY29uZHMpIHRoZW4gdGhlIHByb21pc2UgaXMgcmVqZWN0ZWQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWVcbiAqIEBwYXJhbSB7YW55PX0gYXJnc1xuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XG4gKiBAcmV0dXJuc1xuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XG5cblx0Ly8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcblx0XHR0aW1lb3V0ID0gMDtcblx0fVxuXG5cdC8vIENyZWF0ZSBhIHByb21pc2Vcblx0cmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcblxuXHRcdC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXG5cdFx0dmFyIGNhbGxiYWNrSUQ7XG5cdFx0ZG8ge1xuXHRcdFx0Y2FsbGJhY2tJRCA9IG5hbWUgKyAnLScgKyByYW5kb21GdW5jKCk7XG5cdFx0fSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcblxuXHRcdHZhciB0aW1lb3V0SGFuZGxlO1xuXHRcdC8vIFNldCB0aW1lb3V0XG5cdFx0aWYgKHRpbWVvdXQgPiAwKSB7XG5cdFx0XHR0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XG5cdFx0XHRcdHJlamVjdChFcnJvcignQ2FsbCB0byAnICsgbmFtZSArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xuXHRcdFx0fSwgdGltZW91dCk7XG5cdFx0fVxuXG5cdFx0Ly8gU3RvcmUgY2FsbGJhY2tcblx0XHRjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XG5cdFx0XHR0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxuXHRcdFx0cmVqZWN0OiByZWplY3QsXG5cdFx0XHRyZXNvbHZlOiByZXNvbHZlXG5cdFx0fTtcblxuXHRcdHRyeSB7XG5cdFx0XHRjb25zdCBwYXlsb2FkID0ge1xuXHRcdFx0XHRuYW1lLFxuXHRcdFx0XHRhcmdzLFxuXHRcdFx0XHRjYWxsYmFja0lELFxuXHRcdFx0fTtcblxuICAgICAgICAgICAgLy8gTWFrZSB0aGUgY2FsbFxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKCdDJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcbiAgICAgICAgfSBjYXRjaCAoZSkge1xuICAgICAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGUpO1xuICAgICAgICB9XG4gICAgfSk7XG59XG5cbndpbmRvdy5PYmZ1c2NhdGVkQ2FsbCA9IChpZCwgYXJncywgdGltZW91dCkgPT4ge1xuXG4gICAgLy8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG4gICAgaWYgKHRpbWVvdXQgPT0gbnVsbCkge1xuICAgICAgICB0aW1lb3V0ID0gMDtcbiAgICB9XG5cbiAgICAvLyBDcmVhdGUgYSBwcm9taXNlXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcblxuICAgICAgICAvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuICAgICAgICB2YXIgY2FsbGJhY2tJRDtcbiAgICAgICAgZG8ge1xuICAgICAgICAgICAgY2FsbGJhY2tJRCA9IGlkICsgJy0nICsgcmFuZG9tRnVuYygpO1xuICAgICAgICB9IHdoaWxlIChjYWxsYmFja3NbY2FsbGJhY2tJRF0pO1xuXG4gICAgICAgIHZhciB0aW1lb3V0SGFuZGxlO1xuICAgICAgICAvLyBTZXQgdGltZW91dFxuICAgICAgICBpZiAodGltZW91dCA+IDApIHtcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGUgPSBzZXRUaW1lb3V0KGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgICAgICByZWplY3QoRXJyb3IoJ0NhbGwgdG8gbWV0aG9kICcgKyBpZCArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xuICAgICAgICAgICAgfSwgdGltZW91dCk7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBTdG9yZSBjYWxsYmFja1xuICAgICAgICBjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxuICAgICAgICAgICAgcmVqZWN0OiByZWplY3QsXG4gICAgICAgICAgICByZXNvbHZlOiByZXNvbHZlXG4gICAgICAgIH07XG5cbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XG5cdFx0XHRcdGlkLFxuXHRcdFx0XHRhcmdzLFxuXHRcdFx0XHRjYWxsYmFja0lELFxuXHRcdFx0fTtcblxuICAgICAgICAgICAgLy8gTWFrZSB0aGUgY2FsbFxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKCdjJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcbiAgICAgICAgfSBjYXRjaCAoZSkge1xuICAgICAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGUpO1xuICAgICAgICB9XG4gICAgfSk7XG59O1xuXG5cbi8qKlxuICogQ2FsbGVkIGJ5IHRoZSBiYWNrZW5kIHRvIHJldHVybiBkYXRhIHRvIGEgcHJldmlvdXNseSBjYWxsZWRcbiAqIGJpbmRpbmcgaW52b2NhdGlvblxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBpbmNvbWluZ01lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGxiYWNrKGluY29taW5nTWVzc2FnZSkge1xuXHQvLyBQYXJzZSB0aGUgbWVzc2FnZVxuXHRsZXQgbWVzc2FnZTtcblx0dHJ5IHtcblx0XHRtZXNzYWdlID0gSlNPTi5wYXJzZShpbmNvbWluZ01lc3NhZ2UpO1xuXHR9IGNhdGNoIChlKSB7XG5cdFx0Y29uc3QgZXJyb3IgPSBgSW52YWxpZCBKU09OIHBhc3NlZCB0byBjYWxsYmFjazogJHtlLm1lc3NhZ2V9LiBNZXNzYWdlOiAke2luY29taW5nTWVzc2FnZX1gO1xuXHRcdHJ1bnRpbWUuTG9nRGVidWcoZXJyb3IpO1xuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XG5cdH1cblx0bGV0IGNhbGxiYWNrSUQgPSBtZXNzYWdlLmNhbGxiYWNraWQ7XG5cdGxldCBjYWxsYmFja0RhdGEgPSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XG5cdGlmICghY2FsbGJhY2tEYXRhKSB7XG5cdFx0Y29uc3QgZXJyb3IgPSBgQ2FsbGJhY2sgJyR7Y2FsbGJhY2tJRH0nIG5vdCByZWdpc3RlcmVkISEhYDtcblx0XHRjb25zb2xlLmVycm9yKGVycm9yKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZVxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XG5cdH1cblx0Y2xlYXJUaW1lb3V0KGNhbGxiYWNrRGF0YS50aW1lb3V0SGFuZGxlKTtcblxuXHRkZWxldGUgY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xuXG5cdGlmIChtZXNzYWdlLmVycm9yKSB7XG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChtZXNzYWdlLmVycm9yKTtcblx0fSBlbHNlIHtcblx0XHRjYWxsYmFja0RhdGEucmVzb2x2ZShtZXNzYWdlLnJlc3VsdCk7XG5cdH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fICAgIFxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApIFxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vICBcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSAnLi9jYWxscyc7XG5cbi8vIFRoaXMgaXMgd2hlcmUgd2UgYmluZCBnbyBtZXRob2Qgd3JhcHBlcnNcbndpbmRvdy5nbyA9IHt9O1xuXG5leHBvcnQgZnVuY3Rpb24gU2V0QmluZGluZ3MoYmluZGluZ3NNYXApIHtcblx0dHJ5IHtcblx0XHRiaW5kaW5nc01hcCA9IEpTT04ucGFyc2UoYmluZGluZ3NNYXApO1xuXHR9IGNhdGNoIChlKSB7XG5cdFx0Y29uc29sZS5lcnJvcihlKTtcblx0fVxuXG5cdC8vIEluaXRpYWxpc2UgdGhlIGJpbmRpbmdzIG1hcFxuXHR3aW5kb3cuZ28gPSB3aW5kb3cuZ28gfHwge307XG5cblx0Ly8gSXRlcmF0ZSBwYWNrYWdlIG5hbWVzXG5cdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwKS5mb3JFYWNoKChwYWNrYWdlTmFtZSkgPT4ge1xuXG5cdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XG5cdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gfHwge307XG5cblx0XHQvLyBJdGVyYXRlIHN0cnVjdCBuYW1lc1xuXHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXSkuZm9yRWFjaCgoc3RydWN0TmFtZSkgPT4ge1xuXG5cdFx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3Rcblx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdIHx8IHt9O1xuXG5cdFx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0pLmZvckVhY2goKG1ldGhvZE5hbWUpID0+IHtcblxuXHRcdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdW21ldGhvZE5hbWVdID0gZnVuY3Rpb24gKCkge1xuXG5cdFx0XHRcdFx0Ly8gTm8gdGltZW91dCBieSBkZWZhdWx0XG5cdFx0XHRcdFx0bGV0IHRpbWVvdXQgPSAwO1xuXG5cdFx0XHRcdFx0Ly8gQWN0dWFsIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZnVuY3Rpb24gZHluYW1pYygpIHtcblx0XHRcdFx0XHRcdGNvbnN0IGFyZ3MgPSBbXS5zbGljZS5jYWxsKGFyZ3VtZW50cyk7XG5cdFx0XHRcdFx0XHRyZXR1cm4gQ2FsbChbcGFja2FnZU5hbWUsIHN0cnVjdE5hbWUsIG1ldGhvZE5hbWVdLmpvaW4oJy4nKSwgYXJncywgdGltZW91dCk7XG5cdFx0XHRcdFx0fVxuXG5cdFx0XHRcdFx0Ly8gQWxsb3cgc2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZHluYW1pYy5zZXRUaW1lb3V0ID0gZnVuY3Rpb24gKG5ld1RpbWVvdXQpIHtcblx0XHRcdFx0XHRcdHRpbWVvdXQgPSBuZXdUaW1lb3V0O1xuXHRcdFx0XHRcdH07XG5cblx0XHRcdFx0XHQvLyBBbGxvdyBnZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cblx0XHRcdFx0XHRkeW5hbWljLmdldFRpbWVvdXQgPSBmdW5jdGlvbiAoKSB7XG5cdFx0XHRcdFx0XHRyZXR1cm4gdGltZW91dDtcblx0XHRcdFx0XHR9O1xuXG5cdFx0XHRcdFx0cmV0dXJuIGR5bmFtaWM7XG5cdFx0XHRcdH0oKTtcblx0XHRcdH0pO1xuXHRcdH0pO1xuXHR9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZCgpIHtcbiAgICB3aW5kb3cubG9jYXRpb24ucmVsb2FkKCk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWRBcHAoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUicpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U3lzdGVtRGVmYXVsdFRoZW1lKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FTRFQnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldExpZ2h0VGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQUxUJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXREYXJrVGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQURUJyk7XG59XG5cbi8qKlxuICogUGxhY2UgdGhlIHdpbmRvdyBpbiB0aGUgY2VudGVyIG9mIHRoZSBzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dDZW50ZXIoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXYycpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIHdpbmRvdyB0aXRsZVxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZVxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0VGl0bGUodGl0bGUpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dUJyArIHRpdGxlKTtcbn1cblxuLyoqXG4gKiBNYWtlcyB0aGUgd2luZG93IGdvIGZ1bGxzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dGdWxsc2NyZWVuKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0YnKTtcbn1cblxuLyoqXG4gKiBSZXZlcnRzIHRoZSB3aW5kb3cgZnJvbSBmdWxsc2NyZWVuXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5mdWxsc2NyZWVuKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2YnKTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBpbiBmdWxsIHNjcmVlbiBtb2RlIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNGdWxsc2NyZWVuKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzRnVsbHNjcmVlblwiKTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXczonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG4vKipcbiAqIEdldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7dzogbnVtYmVyLCBoOiBudW1iZXJ9Pn0gVGhlIHNpemUgb2YgdGhlIHdpbmRvd1xuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRTaXplKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFNpemVcIik7XG59XG5cbi8qKlxuICogU2V0IHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXWjonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1pblNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcbn1cblxuXG5cbi8qKlxuICogU2V0IHRoZSB3aW5kb3cgQWx3YXlzT25Ub3Agb3Igbm90IG9uIHRvcFxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEFsd2F5c09uVG9wKGIpIHtcblxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FUUDonICsgKGIgPyAnMScgOiAnMCcpKTtcbn1cblxuXG5cblxuLyoqXG4gKiBTZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0geFxuICogQHBhcmFtIHtudW1iZXJ9IHlcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFBvc2l0aW9uKHgsIHkpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dwOicgKyB4ICsgJzonICsgeSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7eDogbnVtYmVyLCB5OiBudW1iZXJ9Pn0gVGhlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFBvc2l0aW9uKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFBvc1wiKTtcbn1cblxuLyoqXG4gKiBIaWRlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dIaWRlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0gnKTtcbn1cblxuLyoqXG4gKiBTaG93IHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTaG93KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1MnKTtcbn1cblxuLyoqXG4gKiBNYXhpbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXTScpO1xufVxuXG4vKipcbiAqIFRvZ2dsZSB0aGUgTWF4aW1pc2Ugb2YgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1RvZ2dsZU1heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3QnKTtcbn1cblxuLyoqXG4gKiBVbm1heGltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1UnKTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc01heGltaXNlZCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc01heGltaXNlZFwiKTtcbn1cblxuLyoqXG4gKiBNaW5pbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWluaW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXbScpO1xufVxuXG4vKipcbiAqIFVubWluaW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWluaW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdScpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1pbmltaXNlZCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWluaW1pc2VkKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWluaW1pc2VkXCIpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG5vcm1hbCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTm9ybWFsKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTm9ybWFsXCIpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gUiBSZWRcbiAqIEBwYXJhbSB7bnVtYmVyfSBHIEdyZWVuXG4gKiBAcGFyYW0ge251bWJlcn0gQiBCbHVlXG4gKiBAcGFyYW0ge251bWJlcn0gQSBBbHBoYVxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QmFja2dyb3VuZENvbG91cihSLCBHLCBCLCBBKSB7XG4gICAgbGV0IHJnYmEgPSBKU09OLnN0cmluZ2lmeSh7cjogUiB8fCAwLCBnOiBHIHx8IDAsIGI6IEIgfHwgMCwgYTogQSB8fCAyNTV9KTtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dyOicgKyByZ2JhKTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG5cbi8qKlxuICogR2V0cyB0aGUgYWxsIHNjcmVlbnMuIENhbGwgdGhpcyBhbmV3IGVhY2ggdGltZSB5b3Ugd2FudCB0byByZWZyZXNoIGRhdGEgZnJvbSB0aGUgdW5kZXJseWluZyB3aW5kb3dpbmcgc3lzdGVtLlxuICogQGV4cG9ydFxuICogQHR5cGVkZWYge2ltcG9ydCgnLi4vd3JhcHBlci9ydW50aW1lJykuU2NyZWVufSBTY3JlZW5cbiAqIEByZXR1cm4ge1Byb21pc2U8e1NjcmVlbltdfT59IFRoZSBzY3JlZW5zXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTY3JlZW5HZXRBbGwoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6U2NyZWVuR2V0QWxsXCIpO1xufVxuIiwgIi8qKlxuICogQGRlc2NyaXB0aW9uOiBVc2UgdGhlIHN5c3RlbSBkZWZhdWx0IGJyb3dzZXIgdG8gb3BlbiB0aGUgdXJsXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIFxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJyb3dzZXJPcGVuVVJMKHVybCkge1xuICB3aW5kb3cuV2FpbHNJbnZva2UoJ0JPOicgKyB1cmwpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cbi8qKlxuICogU2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dFxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2xpcGJvYXJkU2V0VGV4dCh0ZXh0KSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6Q2xpcGJvYXJkU2V0VGV4dFwiLCBbdGV4dF0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgdGV4dCBjb250ZW50IG9mIHRoZSBjbGlwYm9hcmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHtzdHJpbmd9Pn0gVGV4dCBjb250ZW50IG9mIHRoZSBjbGlwYm9hcmRcblxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2xpcGJvYXJkR2V0VGV4dCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRHZXRUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge0V2ZW50c09ufSBmcm9tIFwiLi9ldmVudHNcIjtcblxuY29uc3QgZmxhZ3MgPSB7XG4gICAgcmVnaXN0ZXJlZDogZmFsc2UsXG4gICAgZGVmYXVsdFVzZURyb3BUYXJnZXQ6IHRydWUsXG4gICAgdXNlRHJvcFRhcmdldDogdHJ1ZSxcbiAgICBwcmV2RWxlbWVudDogbnVsbFxufTtcblxuZnVuY3Rpb24gb25EcmFnT3ZlcihlKSB7XG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIGUuc3RvcFByb3BhZ2F0aW9uKCk7XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChlLngsIGUueSk7XG5cbiAgICBpZiAodGFyZ2V0RWxlbWVudCA9PT0gZmxhZ3MucHJldkVsZW1lbnQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGNvbnN0IHN0eWxlID0gdGFyZ2V0RWxlbWVudC5zdHlsZTtcbiAgICBsZXQgY3NzRHJvcFZhbHVlID0gbnVsbDtcbiAgICBpZiAoT2JqZWN0LmtleXMoc3R5bGUpLmZpbmRJbmRleChrZXkgPT4gc3R5bGVba2V5XSA9PT0gd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSkgPCAwKSB7XG4gICAgICAgIHRhcmdldEVsZW1lbnQgPSB0YXJnZXRFbGVtZW50LmNsb3Nlc3QoYFtzdHlsZSo9JyR7d2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eX0nXWApO1xuICAgIH1cblxuICAgIGlmICh0YXJnZXRFbGVtZW50ID09PSBudWxsKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBjc3NEcm9wVmFsdWUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXRFbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkpO1xuICAgIGlmIChjc3NEcm9wVmFsdWUpIHtcbiAgICAgICAgY3NzRHJvcFZhbHVlID0gY3NzRHJvcFZhbHVlLnRyaW0oKTtcbiAgICB9XG5cbiAgICBpZiAoY3NzRHJvcFZhbHVlID09PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlKSB7XG4gICAgICAgIHRhcmdldEVsZW1lbnQuY2xhc3NMaXN0LmFkZChcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcbiAgICB9IGVsc2UgaWYgKGZsYWdzLnByZXZFbGVtZW50KSB7XG4gICAgICAgIHRhcmdldEVsZW1lbnQuY2xhc3NMaXN0LnJlbW92ZShcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcbiAgICAgICAgZmxhZ3MucHJldkVsZW1lbnQuY2xhc3NMaXN0LnJlbW92ZShcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcbiAgICB9XG4gICAgZmxhZ3MucHJldkVsZW1lbnQgPSB0YXJnZXRFbGVtZW50O1xufVxuXG5mdW5jdGlvbiBvbkRyYWdMZWF2ZShlKSB7XG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIGUuc3RvcFByb3BhZ2F0aW9uKCk7XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChlLngsIGUueSk7XG4gICAgbGV0IGNzc0Ryb3BWYWx1ZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldEVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSk7XG4gICAgaWYgKGNzc0Ryb3BWYWx1ZSkge1xuICAgICAgICBjc3NEcm9wVmFsdWUgPSBjc3NEcm9wVmFsdWUudHJpbSgpO1xuICAgIH1cbiAgICBpZiAoY3NzRHJvcFZhbHVlICE9PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlICYmIGZsYWdzLnByZXZFbGVtZW50KSB7XG4gICAgICAgIHRhcmdldEVsZW1lbnQuY2xhc3NMaXN0LnJlbW92ZShcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcbiAgICAgICAgZmxhZ3MucHJldkVsZW1lbnQuY2xhc3NMaXN0LnJlbW92ZShcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIG9uRHJvcChlKSB7XG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIGUuc3RvcFByb3BhZ2F0aW9uKCk7XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChlLngsIGUueSk7XG4gICAgbGV0IGNzc0Ryb3BWYWx1ZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldEVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSk7XG4gICAgaWYgKGNzc0Ryb3BWYWx1ZSkge1xuICAgICAgICBjc3NEcm9wVmFsdWUgPSBjc3NEcm9wVmFsdWUudHJpbSgpO1xuICAgIH1cbiAgICBpZiAoY3NzRHJvcFZhbHVlICE9PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlKSB7XG4gICAgICAgIGlmIChmbGFncy5wcmV2RWxlbWVudCkge1xuICAgICAgICAgICAgdGFyZ2V0RWxlbWVudC5jbGFzc0xpc3QucmVtb3ZlKFwid2FpbHMtZHJvcC10YXJnZXQtYWN0aXZlXCIpO1xuICAgICAgICAgICAgZmxhZ3MucHJldkVsZW1lbnQuY2xhc3NMaXN0LnJlbW92ZShcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKENhblJlc29sdmVGaWxlUGF0aHMoKSkge1xuICAgICAgICAvLyBwcm9jZXNzIGZpbGVzXG4gICAgICAgIGxldCBmaWxlcyA9IFtdO1xuICAgICAgICBpZiAoZS5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcbiAgICAgICAgICAgIGZpbGVzID0gWy4uLmUuZGF0YVRyYW5zZmVyLml0ZW1zXS5tYXAoKGl0ZW0sIGkpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuIGl0ZW0uZ2V0QXNGaWxlKCk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBmaWxlcyA9IFsuLi5lLmRhdGFUcmFuc2Zlci5maWxlc107XG4gICAgICAgIH1cbiAgICAgICAgd2luZG93LnJ1bnRpbWUuUmVzb2x2ZUZpbGVQYXRocyhlLngsIGUueSwgZmlsZXMpO1xuICAgIH1cblxuICAgIGlmIChmbGFncy5wcmV2RWxlbWVudCkge1xuICAgICAgICBmbGFncy5wcmV2RWxlbWVudC5jbGFzc0xpc3QucmVtb3ZlKFwid2FpbHMtZHJvcC10YXJnZXQtYWN0aXZlXCIpO1xuICAgIH1cbn1cblxuXG4vKipcbiAqIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIGNoZWNrcyB0aGUgYnJvd3NlcidzIGNhcGFiaWxpdHkgb2Ygc2VuZGluZyBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0c1xuICpcbiAqIEByZXR1cm5zIHtib29sZWFufVxuICogQGNvbnN0cnVjdG9yXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYW5SZXNvbHZlRmlsZVBhdGhzKCkge1xuICAgIHJldHVybiB3aW5kb3cuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyAhPSBudWxsO1xufVxuXG4vKipcbiAqIFJlc29sdmVGaWxlUGF0aHMgc2VuZHMgZHJvcCBldmVudHMgdG8gdGhlIEdPIHNpZGUgdG8gcmVzb2x2ZSBmaWxlIHBhdGhzIG9uIHdpbmRvd3MuXG4gKlxuICogQHBhcmFtIHhcbiAqIEBwYXJhbSB5XG4gKiBAcGFyYW0gZmlsZXNcbiAqIEBjb25zdHJ1Y3RvclxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVzb2x2ZUZpbGVQYXRocyh4LCB5LCBmaWxlcykge1xuICAgIC8vIE9ubHkgZm9yIHdpbmRvd3Mgd2VidmlldzIgPj0gMS4wLjE3NzQuMzBcbiAgICAvLyBodHRwczovL2xlYXJuLm1pY3Jvc29mdC5jb20vZW4tdXMvbWljcm9zb2Z0LWVkZ2Uvd2VidmlldzIvcmVmZXJlbmNlL3dpbjMyL2ljb3Jld2VidmlldzJ3ZWJtZXNzYWdlcmVjZWl2ZWRldmVudGFyZ3MyP3ZpZXc9d2VidmlldzItMS4wLjE4MjMuMzIjYXBwbGllcy10b1xuICAgIGlmICh3aW5kb3cuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cykge1xuICAgICAgICBjaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyhgZmlsZTpkcm9wOiR7eH06JHt5fWAsIGZpbGVzKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBjb25zb2xlLndhcm4oXCJ1bnN1cHBvcnRlZCBwbGF0Zm9ybVwiKTtcbn1cblxuLyoqXG4gKiBDYWxsYmFjayBmb3IgRHJhZ0FuZERyb3BPbiByZXR1cm5zIGEgc2xpY2Ugb2YgZmlsZSBwYXRoIHN0cmluZ3Mgd2hlbiBhIGRyb3AgaXMgZmluaXNoZWQuXG4gKlxuICogQGV4cG9ydFxuICogQGNhbGxiYWNrIERyYWdBbmREcm9wQ2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSB4IC0geCBjb29yZGluYXRlIG9mIHRoZSBkcm9wXG4gKiBAcGFyYW0ge251bWJlcn0geSAtIHkgY29vcmRpbmF0ZSBvZiB0aGUgZHJvcFxuICogQHBhcmFtIHtzdHJpbmdbXX0gcGF0aHMgLSBBIGxpc3Qgb2YgZmlsZSBwYXRocy5cbiAqL1xuXG4vKipcbiAqIERyYWdBbmREcm9wT24gbGlzdGVucyB0byBkcmFnIGFuZCBkcm9wIGV2ZW50cyBhbmQgY2FsbHMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGNvb3JkaW5hdGVzIG9mIHRoZSBkcm9wIGFuZCBhbiBhcnJheSBvZiBwYXRoIHN0cmluZ3MuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtEcmFnQW5kRHJvcENhbGxiYWNrfSBjYWxsYmFjayAtIENhbGxiYWNrIGZvciBEcmFnQW5kRHJvcE9uIHJldHVybnMgYSBzbGljZSBvZiBmaWxlIHBhdGggc3RyaW5ncyB3aGVuIGEgZHJvcCBpcyBmaW5pc2hlZC5cbiAqIEBwYXJhbSB7Ym9vbGVhbn0gW3VzZURyb3BUYXJnZXQ9dHJ1ZV0gLSBPbmx5IGNhbGwgdGhlIGNhbGxiYWNrIHdoZW4gdGhlIGRyb3AgZmluaXNoZWQgb24gYW4gZWxlbWVudCB0aGF0IGhhcyB0aGUgZHJvcCB0YXJnZXQgc3R5bGUuICgtLXdhaWxzLWRyb3AtdGFyZ2V0KVxuICovXG5leHBvcnQgZnVuY3Rpb24gRHJhZ0FuZERyb3BPbihjYWxsYmFjaywgdXNlRHJvcFRhcmdldCkge1xuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3AgfHwgdHlwZW9mIGNhbGxiYWNrICE9PSBcImZ1bmN0aW9uXCIpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmIChmbGFncy5yZWdpc3RlcmVkKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG4gICAgZmxhZ3MucmVnaXN0ZXJlZCA9IHRydWU7XG5cblxuICAgIGNvbnN0IHVEVFBUID0gdHlwZW9mIHVzZURyb3BUYXJnZXQ7XG4gICAgZmxhZ3MudXNlRHJvcFRhcmdldCA9IHVEVFBUID09PSBcInVuZGVmaW5lZFwiIHx8IHVEVFBUICE9PSBcImJvb2xlYW5cIiA/IGZsYWdzLmRlZmF1bHRVc2VEcm9wVGFyZ2V0IDogdXNlRHJvcFRhcmdldDtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBvbkRyYWdPdmVyKTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgb25EcmFnTGVhdmUpO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgb25Ecm9wKTtcblxuICAgIGxldCBjYiA9IGNhbGxiYWNrO1xuICAgIGlmIChmbGFncy51c2VEcm9wVGFyZ2V0KSB7XG4gICAgICAgIGNiID0gZnVuY3Rpb24gKHgsIHksIHBhdGhzKSB7XG4gICAgICAgICAgICBsZXQgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgICAgICAgICBpZiAoIXRhcmdldEVsZW1lbnQpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBsZXQgY3NzRHJvcFZhbHVlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0RWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZSh3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFByb3BlcnR5KTtcbiAgICAgICAgICAgIGlmIChjc3NEcm9wVmFsdWUpIHtcbiAgICAgICAgICAgICAgICBjc3NEcm9wVmFsdWUgPSBjc3NEcm9wVmFsdWUudHJpbSgpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgaWYgKGNzc0Ryb3BWYWx1ZSAhPT0gd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BWYWx1ZSkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhbGxiYWNrKHgsIHksIHBhdGhzKTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIEV2ZW50c09uKFwid2FpbHM6ZG5kOmRyb3BcIiwgY2IpO1xufVxuXG4vKipcbiAqIERyYWdBbmREcm9wT2ZmIHJlbW92ZXMgdGhlIGRyYWcgYW5kIGRyb3AgbGlzdGVuZXJzIGFuZCBoYW5kbGVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIERyYWdBbmREcm9wT2ZmKCkge1xuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcmFnb3ZlcicsIG9uRHJhZ092ZXIpO1xuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCBvbkRyYWdMZWF2ZSk7XG4gICAgd2luZG93LnJlbW92ZUV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCBvbkRyb3ApO1xuICAgIEV2ZW50c09mZihcIndhaWxzOmRuZDpkcm9wXCIpO1xuICAgIGZsYWdzLnJlZ2lzdGVyZWQgPSBmYWxzZTtcbn0iLCAiLypcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG5cblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXG4qL1xuZXhwb3J0IGZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgY29uc3QgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcbiAgICBjb25zdCBjb21wdXRlZFN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XG4gICAgY29uc3QgZGVmYXVsdENvbnRleHRNZW51QWN0aW9uID0gY29tcHV0ZWRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKTtcbiAgICBzd2l0Y2ggKGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbikge1xuICAgICAgICBjYXNlIFwic2hvd1wiOlxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBjYXNlIFwiaGlkZVwiOlxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgICAgICAgICBpZiAoZWxlbWVudC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBhbmQgYWN0aW9uIGlzIG9uIHRoZSBzZWxlY3RlZCBlbGVtZW50c1xuICAgICAgICAgICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgICAgICAgICAgY29uc3QgaGFzU2VsZWN0aW9uID0gKHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDApXG4gICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgICAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXG4gICAgICAgICAgICBpZiAoZWxlbWVudC50YWdOYW1lID09PSBcIklOUFVUXCIgfHwgZWxlbWVudC50YWdOYW1lID09PSBcIlRFWFRBUkVBXCIpIHtcbiAgICAgICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghZWxlbWVudC5yZWFkT25seSAmJiAhZWxlbWVudC5kaXNhYmxlZCkpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xuaW1wb3J0IHtldmVudExpc3RlbmVycywgRXZlbnRzRW1pdCwgRXZlbnRzTm90aWZ5LCBFdmVudHNPZmYsIEV2ZW50c09uLCBFdmVudHNPbmNlLCBFdmVudHNPbk11bHRpcGxlfSBmcm9tICcuL2V2ZW50cyc7XG5pbXBvcnQge0NhbGwsIENhbGxiYWNrLCBjYWxsYmFja3N9IGZyb20gJy4vY2FsbHMnO1xuaW1wb3J0IHtTZXRCaW5kaW5nc30gZnJvbSBcIi4vYmluZGluZ3NcIjtcbmltcG9ydCAqIGFzIFdpbmRvdyBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCAqIGFzIFNjcmVlbiBmcm9tIFwiLi9zY3JlZW5cIjtcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3NlclwiO1xuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZFwiO1xuaW1wb3J0ICogYXMgRHJhZ0FuZERyb3AgZnJvbSBcIi4vZHJhZ2FuZGRyb3BcIjtcbmltcG9ydCAqIGFzIENvbnRleHRNZW51IGZyb20gXCIuL2NvbnRleHRtZW51XCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUScpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1MnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdIJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpFbnZpcm9ubWVudFwiKTtcbn1cblxuLy8gVGhlIEpTIHJ1bnRpbWVcbndpbmRvdy5ydW50aW1lID0ge1xuICAgIC4uLkxvZyxcbiAgICAuLi5XaW5kb3csXG4gICAgLi4uQnJvd3NlcixcbiAgICAuLi5TY3JlZW4sXG4gICAgLi4uQ2xpcGJvYXJkLFxuICAgIC4uLkRyYWdBbmREcm9wLFxuICAgIEV2ZW50c09uLFxuICAgIEV2ZW50c09uY2UsXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcbiAgICBFdmVudHNFbWl0LFxuICAgIEV2ZW50c09mZixcbiAgICBFbnZpcm9ubWVudCxcbiAgICBTaG93LFxuICAgIEhpZGUsXG4gICAgUXVpdFxufTtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cud2FpbHMgPSB7XG4gICAgQ2FsbGJhY2ssXG4gICAgRXZlbnRzTm90aWZ5LFxuICAgIFNldEJpbmRpbmdzLFxuICAgIGV2ZW50TGlzdGVuZXJzLFxuICAgIGNhbGxiYWNrcyxcbiAgICBmbGFnczoge1xuICAgICAgICBkaXNhYmxlU2Nyb2xsYmFyRHJhZzogZmFsc2UsXG4gICAgICAgIGRpc2FibGVEZWZhdWx0Q29udGV4dE1lbnU6IGZhbHNlLFxuICAgICAgICBlbmFibGVSZXNpemU6IGZhbHNlLFxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxuICAgICAgICBib3JkZXJUaGlja25lc3M6IDYsXG4gICAgICAgIHNob3VsZERyYWc6IGZhbHNlLFxuICAgICAgICBkZWZlckRyYWdUb01vdXNlTW92ZTogdHJ1ZSxcbiAgICAgICAgY3NzRHJhZ1Byb3BlcnR5OiBcIi0td2FpbHMtZHJhZ2dhYmxlXCIsXG4gICAgICAgIGNzc0RyYWdWYWx1ZTogXCJkcmFnXCIsXG4gICAgICAgIGNzc0Ryb3BQcm9wZXJ0eTogXCItLXdhaWxzLWRyb3AtdGFyZ2V0XCIsXG4gICAgICAgIGNzc0Ryb3BWYWx1ZTogXCJkcm9wXCIsXG4gICAgICAgIGVuYWJsZVdhaWxzRHJhZ0FuZERyb3A6IGZhbHNlLFxuICAgICAgICB3YWlsc0Ryb3BQcmV2aW91c0VsZW1lbnQ6IG51bGwsXG4gICAgfVxufTtcblxuLy8gU2V0IHRoZSBiaW5kaW5nc1xuaWYgKHdpbmRvdy53YWlsc2JpbmRpbmdzKSB7XG4gICAgd2luZG93LndhaWxzLlNldEJpbmRpbmdzKHdpbmRvdy53YWlsc2JpbmRpbmdzKTtcbiAgICBkZWxldGUgd2luZG93LndhaWxzLlNldEJpbmRpbmdzO1xufVxuXG4vLyAoYm9vbCkgVGhpcyBpcyBldmFsdWF0ZWQgYXQgYnVpbGQgdGltZSBpbiBwYWNrYWdlLmpzb25cbmlmICghREVCVUcpIHtcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XG59XG5cbmxldCBkcmFnVGVzdCA9IGZ1bmN0aW9uIChlKSB7XG4gICAgdmFyIHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnUHJvcGVydHkpO1xuICAgIGlmICh2YWwpIHtcbiAgICAgIHZhbCA9IHZhbC50cmltKCk7XG4gICAgfVxuICAgIFxuICAgIGlmICh2YWwgIT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIGlmIChlLmJ1dHRvbnMgIT09IDEpIHtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWdnaW5nIGlmIG5vdCB0aGUgcHJpbWFyeSBidXR0b24gaGFzIGJlZW4gY2xpY2tlZC5cbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIGlmIChlLmRldGFpbCAhPT0gMSkge1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZ2dpbmcgaWYgbW9yZSB0aGFuIG9uY2UgaGFzIGJlZW4gY2xpY2tlZCwgZS5nLiB3aGVuIGRvdWJsZSBjbGlja2luZ1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgcmV0dXJuIHRydWU7XG59O1xuXG53aW5kb3cud2FpbHMuc2V0Q1NTRHJhZ1Byb3BlcnRpZXMgPSBmdW5jdGlvbiAocHJvcGVydHksIHZhbHVlKSB7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdQcm9wZXJ0eSA9IHByb3BlcnR5O1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUgPSB2YWx1ZTtcbn1cblxud2luZG93LndhaWxzLnNldENTU0Ryb3BQcm9wZXJ0aWVzID0gZnVuY3Rpb24gKHByb3BlcnR5LCB2YWx1ZSkge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkgPSBwcm9wZXJ0eTtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlID0gdmFsdWU7XG59XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xuICAgIC8vIENoZWNrIGZvciByZXNpemluZ1xuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSkge1xuICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJyZXNpemU6XCIgKyB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSk7XG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmIChkcmFnVGVzdChlKSkge1xuICAgICAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRpc2FibGVTY3JvbGxiYXJEcmFnKSB7XG4gICAgICAgICAgICAvLyBUaGlzIGNoZWNrcyBmb3IgY2xpY2tzIG9uIHRoZSBzY3JvbGwgYmFyXG4gICAgICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmZXJEcmFnVG9Nb3VzZU1vdmUpIHtcbiAgICAgICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gdHJ1ZTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGUucHJldmVudERlZmF1bHQoKVxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm47XG4gICAgfSBlbHNlIHtcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcbiAgICB9XG59KTtcblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCAoKSA9PiB7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcbn0pO1xuXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvcjtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSA9IGN1cnNvcjtcbn1cblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIGZ1bmN0aW9uIChlKSB7XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnKSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJkcmFnXCIpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVJlc2l6ZSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yID0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvcjtcbiAgICB9XG4gICAgaWYgKHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcyAmJiB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzKSB7XG4gICAgICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBcInNlLXJlc2l6ZVwiO1xuICAgIH1cbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH0gZWxzZSBpZiAocmlnaHRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyICYmIHRvcEJvcmRlcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlciAmJiByaWdodEJvcmRlcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XG5cbn0pO1xuXG4vLyBTZXR1cCBjb250ZXh0IG1lbnUgaG9va1xud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgZnVuY3Rpb24gKGUpIHtcbiAgICAvLyBhbHdheXMgc2hvdyB0aGUgY29udGV4dG1lbnUgaW4gZGVidWcgJiBkZXZcbiAgICBpZiAoREVCVUcpIHJldHVybjtcblxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZURlZmF1bHRDb250ZXh0TWVudSkge1xuICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgQ29udGV4dE1lbnUucHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShlKTtcbiAgICB9XG59KTtcblxud2luZG93LldhaWxzSW52b2tlKFwicnVudGltZTpyZWFkeVwiKTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWtCQSxXQUFTLGVBQWUsT0FBTyxTQUFTO0FBSXZDLFdBQU8sWUFBWSxNQUFNLFFBQVEsT0FBTztBQUFBLEVBQ3pDO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxRQUFRLFNBQVM7QUFDaEMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFdBQVcsU0FBUztBQUNuQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMs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QUEsTUFDSjtBQUFBLElBQ0o7QUFDQSxRQUFJLENBQUMsT0FBTyxNQUFNLE1BQU0sY0FBYztBQUNsQztBQUFBLElBQ0o7QUFDQSxRQUFJLE9BQU8sTUFBTSxNQUFNLGlCQUFpQixNQUFNO0FBQzFDLGFBQU8sTUFBTSxNQUFNLGdCQUFnQixTQUFTLGdCQUFnQixNQUFNO0FBQUEsSUFDdEU7QUFDQSxRQUFJLE9BQU8sYUFBYSxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU0sbUJBQW1CLE9BQU8sY0FBYyxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU0saUJBQWlCO0FBQzNJLGVBQVMsZ0JBQWdCLE1BQU0sU0FBUztBQUFBLElBQzVDO0FBQ0EsUUFBSSxjQUFjLE9BQU8sYUFBYSxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU07QUFDckUsUUFBSSxhQUFhLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUNoRCxRQUFJLFlBQVksRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBQy9DLFFBQUksZUFBZSxPQUFPLGNBQWMsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBR3ZFLFFBQUksQ0FBQyxjQUFjLENBQUMsZUFBZSxDQUFDLGFBQWEsQ0FBQyxnQkFBZ0IsT0FBTyxNQUFNLE1BQU0sZUFBZSxRQUFXO0FBQzNHLGdCQUFVO0FBQUEsSUFDZCxXQUFXLGVBQWU7QUFBYyxnQkFBVSxXQUFXO0FBQUEsYUFDcEQsY0FBYztBQUFjLGdCQUFVLFdBQVc7QUFBQSxhQUNqRCxjQUFjO0FBQVcsZ0JBQVUsV0FBVztBQUFBLGFBQzlDLGFBQWE7QUFBYSxnQkFBVSxXQUFXO0FBQUEsYUFDL0M7QUFBWSxnQkFBVSxVQUFVO0FBQUEsYUFDaEM7QUFBVyxnQkFBVSxVQUFVO0FBQUEsYUFDL0I7QUFBYyxnQkFBVSxVQUFVO0FBQUEsYUFDbEM7QUFBYSxnQkFBVSxVQUFVO0FBQUEsRUFFOUMsQ0FBQztBQUdELFNBQU8saUJBQWlCLGVBQWUsU0FBVSxHQUFHO0FBRWhELFFBQUk7QUFBTztBQUVYLFFBQUksT0FBTyxNQUFNLE1BQU0sMkJBQTJCO0FBQzlDLFFBQUUsZUFBZTtBQUFBLElBQ3JCLE9BQU87QUFDSCxNQUFZLDBCQUEwQixDQUFDO0FBQUEsSUFDM0M7QUFBQSxFQUNKLENBQUM7QUFFRCxTQUFPLFlBQVksZUFBZTsiLAogICJuYW1lcyI6IFsiRXZlbnRzT2ZmIiwgImV2ZW50TmFtZSIsICJFdmVudHNPZmYiXQp9Cg==
