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
    const newEventListenerList = eventListeners[eventName]?.slice() || [];
    if (newEventListenerList.length) {
      for (let count = newEventListenerList.length - 1; count >= 0; count -= 1) {
        const listener = newEventListenerList[count];
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
  function EventsOffAll() {
    const eventNames = Object.keys(eventListeners);
    eventNames.forEach((eventName) => {
      removeListener(eventName);
    });
  }
  function listenerOff(listener) {
    const eventName = listener.eventName;
    if (eventListeners[eventName] === void 0)
      return;
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
    OnFileDrop: () => OnFileDrop,
    OnFileDropOff: () => OnFileDropOff,
    ResolveFilePaths: () => ResolveFilePaths
  });
  var flags = {
    registered: false,
    defaultUseDropTarget: true,
    useDropTarget: true,
    nextDeactivate: null,
    nextDeactivateTimeout: null
  };
  var DROP_TARGET_ACTIVE = "wails-drop-target-active";
  function checkStyleDropTarget(style) {
    const cssDropValue = style.getPropertyValue(window.wails.flags.cssDropProperty).trim();
    if (cssDropValue) {
      if (cssDropValue === window.wails.flags.cssDropValue) {
        return true;
      }
      return false;
    }
    return false;
  }
  function onDragOver(e) {
    const isFileDrop = e.dataTransfer.types.includes("Files");
    if (!isFileDrop) {
      return;
    }
    e.preventDefault();
    e.dataTransfer.dropEffect = "copy";
    if (!window.wails.flags.enableWailsDragAndDrop) {
      return;
    }
    if (!flags.useDropTarget) {
      return;
    }
    const element = e.target;
    if (flags.nextDeactivate)
      flags.nextDeactivate();
    if (!element || !checkStyleDropTarget(getComputedStyle(element))) {
      return;
    }
    let currentElement = element;
    while (currentElement) {
      if (checkStyleDropTarget(getComputedStyle(currentElement))) {
        currentElement.classList.add(DROP_TARGET_ACTIVE);
      }
      currentElement = currentElement.parentElement;
    }
  }
  function onDragLeave(e) {
    const isFileDrop = e.dataTransfer.types.includes("Files");
    if (!isFileDrop) {
      return;
    }
    e.preventDefault();
    if (!window.wails.flags.enableWailsDragAndDrop) {
      return;
    }
    if (!flags.useDropTarget) {
      return;
    }
    if (!e.target || !checkStyleDropTarget(getComputedStyle(e.target))) {
      return null;
    }
    if (flags.nextDeactivate)
      flags.nextDeactivate();
    flags.nextDeactivate = () => {
      Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach((el) => el.classList.remove(DROP_TARGET_ACTIVE));
      flags.nextDeactivate = null;
      if (flags.nextDeactivateTimeout) {
        clearTimeout(flags.nextDeactivateTimeout);
        flags.nextDeactivateTimeout = null;
      }
    };
    flags.nextDeactivateTimeout = setTimeout(() => {
      if (flags.nextDeactivate)
        flags.nextDeactivate();
    }, 50);
  }
  function onDrop(e) {
    const isFileDrop = e.dataTransfer.types.includes("Files");
    if (!isFileDrop) {
      return;
    }
    e.preventDefault();
    if (!window.wails.flags.enableWailsDragAndDrop) {
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
    if (!flags.useDropTarget) {
      return;
    }
    if (flags.nextDeactivate)
      flags.nextDeactivate();
    Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach((el) => el.classList.remove(DROP_TARGET_ACTIVE));
  }
  function CanResolveFilePaths() {
    return window.chrome?.webview?.postMessageWithAdditionalObjects != null;
  }
  function ResolveFilePaths(x, y, files) {
    if (window.chrome?.webview?.postMessageWithAdditionalObjects) {
      chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
    }
  }
  function OnFileDrop(callback, useDropTarget) {
    if (typeof callback !== "function") {
      console.error("DragAndDropCallback is not a function");
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
        const element = document.elementFromPoint(x, y);
        if (!element || !checkStyleDropTarget(getComputedStyle(element))) {
          return null;
        }
        callback(x, y, paths);
      };
    }
    EventsOn("wails:file-drop", cb);
  }
  function OnFileDropOff() {
    window.removeEventListener("dragover", onDragOver);
    window.removeEventListener("dragleave", onDragLeave);
    window.removeEventListener("drop", onDrop);
    EventsOff("wails:file-drop");
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
    EventsOff,
    EventsOffAll,
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
      enableWailsDragAndDrop: false
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL2RyYWdhbmRkcm9wLmpzIiwgImRlc2t0b3AvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9tYWluLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX19cclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG4vKipcclxuICogU2VuZHMgYSBsb2cgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB3aXRoIHRoZSBnaXZlbiBsZXZlbCArIG1lc3NhZ2VcclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGxldmVsXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5mdW5jdGlvbiBzZW5kTG9nTWVzc2FnZShsZXZlbCwgbWVzc2FnZSkge1xyXG5cclxuXHQvLyBMb2cgTWVzc2FnZSBmb3JtYXQ6XHJcblx0Ly8gbFt0eXBlXVttZXNzYWdlXVxyXG5cdHdpbmRvdy5XYWlsc0ludm9rZSgnTCcgKyBsZXZlbCArIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiB0cmFjZSBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ1RyYWNlKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnVCcsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ1ByaW50KG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnUCcsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBkZWJ1ZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ0RlYnVnKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnRCcsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBpbmZvIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nSW5mbyhtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0knLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gd2FybmluZyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ1dhcm5pbmcobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdXJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGVycm9yIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRXJyb3IobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdFJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIGZhdGFsIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nRmF0YWwobWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdGJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSBMb2cgbGV2ZWwgdG8gdGhlIGdpdmVuIGxvZyBsZXZlbFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBsb2dsZXZlbFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldExvZ0xldmVsKGxvZ2xldmVsKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1MnLCBsb2dsZXZlbCk7XHJcbn1cclxuXHJcbi8vIExvZyBsZXZlbHNcclxuZXhwb3J0IGNvbnN0IExvZ0xldmVsID0ge1xyXG5cdFRSQUNFOiAxLFxyXG5cdERFQlVHOiAyLFxyXG5cdElORk86IDMsXHJcblx0V0FSTklORzogNCxcclxuXHRFUlJPUjogNSxcclxufTtcclxuIiwgIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfX1xyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xyXG5cclxuLy8gRGVmaW5lcyBhIHNpbmdsZSBsaXN0ZW5lciB3aXRoIGEgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdG8gY2FsbGJhY2tcclxuXHJcbi8qKlxyXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcclxuICpcclxuICogQGNsYXNzIExpc3RlbmVyXHJcbiAqL1xyXG5jbGFzcyBMaXN0ZW5lciB7XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXHJcbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xyXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xyXG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXHJcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxyXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xyXG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcclxuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxyXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxyXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XHJcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XHJcblxyXG4gICAgLy8gR2V0IHRoZSBldmVudCBuYW1lXHJcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XHJcblxyXG4gICAgLy8gS2VlcCBhIGxpc3Qgb2YgbGlzdGVuZXIgaW5kZXhlcyB0byBkZXN0cm95XHJcbiAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0/LnNsaWNlKCkgfHwgW107XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgd2UgaGF2ZSBhbnkgbGlzdGVuZXJzIGZvciB0aGlzIGV2ZW50XHJcbiAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoKSB7XHJcblxyXG4gICAgICAgIC8vIEl0ZXJhdGUgbGlzdGVuZXJzXHJcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSBuZXdFdmVudExpc3RlbmVyTGlzdC5sZW5ndGggLSAxOyBjb3VudCA+PSAwOyBjb3VudCAtPSAxKSB7XHJcblxyXG4gICAgICAgICAgICAvLyBHZXQgbmV4dCBsaXN0ZW5lclxyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0W2NvdW50XTtcclxuXHJcbiAgICAgICAgICAgIGxldCBkYXRhID0gZXZlbnREYXRhLmRhdGE7XHJcblxyXG4gICAgICAgICAgICAvLyBEbyB0aGUgY2FsbGJhY2tcclxuICAgICAgICAgICAgY29uc3QgZGVzdHJveSA9IGxpc3RlbmVyLkNhbGxiYWNrKGRhdGEpO1xyXG4gICAgICAgICAgICBpZiAoZGVzdHJveSkge1xyXG4gICAgICAgICAgICAgICAgLy8gaWYgdGhlIGxpc3RlbmVyIGluZGljYXRlZCB0byBkZXN0cm95IGl0c2VsZiwgYWRkIGl0IHRvIHRoZSBkZXN0cm95IGxpc3RcclxuICAgICAgICAgICAgICAgIG5ld0V2ZW50TGlzdGVuZXJMaXN0LnNwbGljZShjb3VudCwgMSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIFVwZGF0ZSBjYWxsYmFja3Mgd2l0aCBuZXcgbGlzdCBvZiBsaXN0ZW5lcnNcclxuICAgICAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0O1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE5vdGlmeSBpbmZvcm1zIGZyb250ZW5kIGxpc3RlbmVycyB0aGF0IGFuIGV2ZW50IHdhcyBlbWl0dGVkIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbm90aWZ5TWVzc2FnZSAtIGVuY29kZWQgbm90aWZpY2F0aW9uIG1lc3NhZ2VcclxuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzTm90aWZ5KG5vdGlmeU1lc3NhZ2UpIHtcclxuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXHJcbiAgICBsZXQgbWVzc2FnZTtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc3QgZXJyb3IgPSAnSW52YWxpZCBKU09OIHBhc3NlZCB0byBOb3RpZnk6ICcgKyBub3RpZnlNZXNzYWdlO1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcbiAgICB9XHJcbiAgICBub3RpZnlMaXN0ZW5lcnMobWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzRW1pdChldmVudE5hbWUpIHtcclxuXHJcbiAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgIG5hbWU6IGV2ZW50TmFtZSxcclxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxyXG4gICAgfTtcclxuXHJcbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXHJcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSkge1xyXG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyc1xyXG4gICAgZGVsZXRlIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV07XHJcblxyXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFWCcgKyBldmVudE5hbWUpO1xyXG59XHJcblxyXG4vKipcclxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXHJcbiAqIG9wdGlvbmFsbHkgbXVsdGlwbGUgbGlzdGVuZXJlcyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XHJcbiAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXHJcblxyXG4gICAgaWYgKGFkZGl0aW9uYWxFdmVudE5hbWVzLmxlbmd0aCA+IDApIHtcclxuICAgICAgICBhZGRpdGlvbmFsRXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XHJcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcclxuICAgICAgICB9KVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogT2ZmIHVucmVnaXN0ZXJzIGFsbCBldmVudCBsaXN0ZW5lcnMgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cclxuICovXHJcbiBleHBvcnQgZnVuY3Rpb24gRXZlbnRzT2ZmQWxsKCkge1xyXG4gICAgY29uc3QgZXZlbnROYW1lcyA9IE9iamVjdC5rZXlzKGV2ZW50TGlzdGVuZXJzKTtcclxuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xyXG4gICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcclxuICAgIH0pXHJcbn1cclxuXHJcbi8qKlxyXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIEV2ZW50c09uXHJcbiAqXHJcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXHJcbiAqL1xyXG4gZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcclxuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcclxuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdID09PSB1bmRlZmluZWQpIHJldHVybjtcclxuXHJcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuXHJcbiAgICAvLyBDbGVhbiB1cCBpZiB0aGVyZSBhcmUgbm8gZXZlbnQgbGlzdGVuZXJzIGxlZnRcclxuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmxlbmd0aCA9PT0gMCkge1xyXG4gICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfICAgICAgIF9fICAgICAgXyBfX1xyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xyXG5cclxuZXhwb3J0IGNvbnN0IGNhbGxiYWNrcyA9IHt9O1xyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBudW1iZXIgZnJvbSB0aGUgbmF0aXZlIGJyb3dzZXIgcmFuZG9tIGZ1bmN0aW9uXHJcbiAqXHJcbiAqIEByZXR1cm5zIG51bWJlclxyXG4gKi9cclxuZnVuY3Rpb24gY3J5cHRvUmFuZG9tKCkge1xyXG5cdHZhciBhcnJheSA9IG5ldyBVaW50MzJBcnJheSgxKTtcclxuXHRyZXR1cm4gd2luZG93LmNyeXB0by5nZXRSYW5kb21WYWx1ZXMoYXJyYXkpWzBdO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIG51bWJlciB1c2luZyBkYSBvbGQtc2tvb2wgTWF0aC5SYW5kb21cclxuICogSSBsaWtlcyB0byBjYWxsIGl0IExPTFJhbmRvbVxyXG4gKlxyXG4gKiBAcmV0dXJucyBudW1iZXJcclxuICovXHJcbmZ1bmN0aW9uIGJhc2ljUmFuZG9tKCkge1xyXG5cdHJldHVybiBNYXRoLnJhbmRvbSgpICogOTAwNzE5OTI1NDc0MDk5MTtcclxufVxyXG5cclxuLy8gUGljayBhIHJhbmRvbSBudW1iZXIgZnVuY3Rpb24gYmFzZWQgb24gYnJvd3NlciBjYXBhYmlsaXR5XHJcbnZhciByYW5kb21GdW5jO1xyXG5pZiAod2luZG93LmNyeXB0bykge1xyXG5cdHJhbmRvbUZ1bmMgPSBjcnlwdG9SYW5kb207XHJcbn0gZWxzZSB7XHJcblx0cmFuZG9tRnVuYyA9IGJhc2ljUmFuZG9tO1xyXG59XHJcblxyXG5cclxuLyoqXHJcbiAqIENhbGwgc2VuZHMgYSBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHRvIGNhbGwgdGhlIGJpbmRpbmcgd2l0aCB0aGVcclxuICogZ2l2ZW4gZGF0YS4gQSBwcm9taXNlIGlzIHJldHVybmVkIGFuZCB3aWxsIGJlIGNvbXBsZXRlZCB3aGVuIHRoZVxyXG4gKiBiYWNrZW5kIHJlc3BvbmRzLiBUaGlzIHdpbGwgYmUgcmVzb2x2ZWQgd2hlbiB0aGUgY2FsbCB3YXMgc3VjY2Vzc2Z1bFxyXG4gKiBvciByZWplY3RlZCBpZiBhbiBlcnJvciBpcyBwYXNzZWQgYmFjay5cclxuICogVGhlcmUgaXMgYSB0aW1lb3V0IG1lY2hhbmlzbS4gSWYgdGhlIGNhbGwgZG9lc24ndCByZXNwb25kIGluIHRoZSBnaXZlblxyXG4gKiB0aW1lIChpbiBtaWxsaXNlY29uZHMpIHRoZW4gdGhlIHByb21pc2UgaXMgcmVqZWN0ZWQuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWVcclxuICogQHBhcmFtIHthbnk9fSBhcmdzXHJcbiAqIEBwYXJhbSB7bnVtYmVyPX0gdGltZW91dFxyXG4gKiBAcmV0dXJuc1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwobmFtZSwgYXJncywgdGltZW91dCkge1xyXG5cclxuXHQvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcclxuXHRpZiAodGltZW91dCA9PSBudWxsKSB7XHJcblx0XHR0aW1lb3V0ID0gMDtcclxuXHR9XHJcblxyXG5cdC8vIENyZWF0ZSBhIHByb21pc2VcclxuXHRyZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xyXG5cclxuXHRcdC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXHJcblx0XHR2YXIgY2FsbGJhY2tJRDtcclxuXHRcdGRvIHtcclxuXHRcdFx0Y2FsbGJhY2tJRCA9IG5hbWUgKyAnLScgKyByYW5kb21GdW5jKCk7XHJcblx0XHR9IHdoaWxlIChjYWxsYmFja3NbY2FsbGJhY2tJRF0pO1xyXG5cclxuXHRcdHZhciB0aW1lb3V0SGFuZGxlO1xyXG5cdFx0Ly8gU2V0IHRpbWVvdXRcclxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xyXG5cdFx0XHR0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XHJcblx0XHRcdFx0cmVqZWN0KEVycm9yKCdDYWxsIHRvICcgKyBuYW1lICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XHJcblx0XHRcdH0sIHRpbWVvdXQpO1xyXG5cdFx0fVxyXG5cclxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXHJcblx0XHRjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XHJcblx0XHRcdHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXHJcblx0XHRcdHJlamVjdDogcmVqZWN0LFxyXG5cdFx0XHRyZXNvbHZlOiByZXNvbHZlXHJcblx0XHR9O1xyXG5cclxuXHRcdHRyeSB7XHJcblx0XHRcdGNvbnN0IHBheWxvYWQgPSB7XHJcblx0XHRcdFx0bmFtZSxcclxuXHRcdFx0XHRhcmdzLFxyXG5cdFx0XHRcdGNhbGxiYWNrSUQsXHJcblx0XHRcdH07XHJcblxyXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXHJcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnQycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcclxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcclxuICAgICAgICB9XHJcbiAgICB9KTtcclxufVxyXG5cclxud2luZG93Lk9iZnVzY2F0ZWRDYWxsID0gKGlkLCBhcmdzLCB0aW1lb3V0KSA9PiB7XHJcblxyXG4gICAgLy8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XHJcbiAgICBpZiAodGltZW91dCA9PSBudWxsKSB7XHJcbiAgICAgICAgdGltZW91dCA9IDA7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ3JlYXRlIGEgcHJvbWlzZVxyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKGZ1bmN0aW9uIChyZXNvbHZlLCByZWplY3QpIHtcclxuXHJcbiAgICAgICAgLy8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcclxuICAgICAgICB2YXIgY2FsbGJhY2tJRDtcclxuICAgICAgICBkbyB7XHJcbiAgICAgICAgICAgIGNhbGxiYWNrSUQgPSBpZCArICctJyArIHJhbmRvbUZ1bmMoKTtcclxuICAgICAgICB9IHdoaWxlIChjYWxsYmFja3NbY2FsbGJhY2tJRF0pO1xyXG5cclxuICAgICAgICB2YXIgdGltZW91dEhhbmRsZTtcclxuICAgICAgICAvLyBTZXQgdGltZW91dFxyXG4gICAgICAgIGlmICh0aW1lb3V0ID4gMCkge1xyXG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgICAgICByZWplY3QoRXJyb3IoJ0NhbGwgdG8gbWV0aG9kICcgKyBpZCArICcgdGltZWQgb3V0LiBSZXF1ZXN0IElEOiAnICsgY2FsbGJhY2tJRCkpO1xyXG4gICAgICAgICAgICB9LCB0aW1lb3V0KTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIFN0b3JlIGNhbGxiYWNrXHJcbiAgICAgICAgY2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xyXG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlOiB0aW1lb3V0SGFuZGxlLFxyXG4gICAgICAgICAgICByZWplY3Q6IHJlamVjdCxcclxuICAgICAgICAgICAgcmVzb2x2ZTogcmVzb2x2ZVxyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XHJcblx0XHRcdFx0aWQsXHJcblx0XHRcdFx0YXJncyxcclxuXHRcdFx0XHRjYWxsYmFja0lELFxyXG5cdFx0XHR9O1xyXG5cclxuICAgICAgICAgICAgLy8gTWFrZSB0aGUgY2FsbFxyXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ2MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xyXG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXHJcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XHJcbiAgICAgICAgfVxyXG4gICAgfSk7XHJcbn07XHJcblxyXG5cclxuLyoqXHJcbiAqIENhbGxlZCBieSB0aGUgYmFja2VuZCB0byByZXR1cm4gZGF0YSB0byBhIHByZXZpb3VzbHkgY2FsbGVkXHJcbiAqIGJpbmRpbmcgaW52b2NhdGlvblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBpbmNvbWluZ01lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsYmFjayhpbmNvbWluZ01lc3NhZ2UpIHtcclxuXHQvLyBQYXJzZSB0aGUgbWVzc2FnZVxyXG5cdGxldCBtZXNzYWdlO1xyXG5cdHRyeSB7XHJcblx0XHRtZXNzYWdlID0gSlNPTi5wYXJzZShpbmNvbWluZ01lc3NhZ2UpO1xyXG5cdH0gY2F0Y2ggKGUpIHtcclxuXHRcdGNvbnN0IGVycm9yID0gYEludmFsaWQgSlNPTiBwYXNzZWQgdG8gY2FsbGJhY2s6ICR7ZS5tZXNzYWdlfS4gTWVzc2FnZTogJHtpbmNvbWluZ01lc3NhZ2V9YDtcclxuXHRcdHJ1bnRpbWUuTG9nRGVidWcoZXJyb3IpO1xyXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcclxuXHR9XHJcblx0bGV0IGNhbGxiYWNrSUQgPSBtZXNzYWdlLmNhbGxiYWNraWQ7XHJcblx0bGV0IGNhbGxiYWNrRGF0YSA9IGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcclxuXHRpZiAoIWNhbGxiYWNrRGF0YSkge1xyXG5cdFx0Y29uc3QgZXJyb3IgPSBgQ2FsbGJhY2sgJyR7Y2FsbGJhY2tJRH0nIG5vdCByZWdpc3RlcmVkISEhYDtcclxuXHRcdGNvbnNvbGUuZXJyb3IoZXJyb3IpOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lXHJcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xyXG5cdH1cclxuXHRjbGVhclRpbWVvdXQoY2FsbGJhY2tEYXRhLnRpbWVvdXRIYW5kbGUpO1xyXG5cclxuXHRkZWxldGUgY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xyXG5cclxuXHRpZiAobWVzc2FnZS5lcnJvcikge1xyXG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChtZXNzYWdlLmVycm9yKTtcclxuXHR9IGVsc2Uge1xyXG5cdFx0Y2FsbGJhY2tEYXRhLnJlc29sdmUobWVzc2FnZS5yZXN1bHQpO1xyXG5cdH1cclxufVxyXG4iLCAiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fICAgIFxyXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApIFxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy8gIFxyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xyXG5cclxuaW1wb3J0IHtDYWxsfSBmcm9tICcuL2NhbGxzJztcclxuXHJcbi8vIFRoaXMgaXMgd2hlcmUgd2UgYmluZCBnbyBtZXRob2Qgd3JhcHBlcnNcclxud2luZG93LmdvID0ge307XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gU2V0QmluZGluZ3MoYmluZGluZ3NNYXApIHtcclxuXHR0cnkge1xyXG5cdFx0YmluZGluZ3NNYXAgPSBKU09OLnBhcnNlKGJpbmRpbmdzTWFwKTtcclxuXHR9IGNhdGNoIChlKSB7XHJcblx0XHRjb25zb2xlLmVycm9yKGUpO1xyXG5cdH1cclxuXHJcblx0Ly8gSW5pdGlhbGlzZSB0aGUgYmluZGluZ3MgbWFwXHJcblx0d2luZG93LmdvID0gd2luZG93LmdvIHx8IHt9O1xyXG5cclxuXHQvLyBJdGVyYXRlIHBhY2thZ2UgbmFtZXNcclxuXHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcCkuZm9yRWFjaCgocGFja2FnZU5hbWUpID0+IHtcclxuXHJcblx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3RcclxuXHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdIHx8IHt9O1xyXG5cclxuXHRcdC8vIEl0ZXJhdGUgc3RydWN0IG5hbWVzXHJcblx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV0pLmZvckVhY2goKHN0cnVjdE5hbWUpID0+IHtcclxuXHJcblx0XHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxyXG5cdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSB8fCB7fTtcclxuXHJcblx0XHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSkuZm9yRWFjaCgobWV0aG9kTmFtZSkgPT4ge1xyXG5cclxuXHRcdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdW21ldGhvZE5hbWVdID0gZnVuY3Rpb24gKCkge1xyXG5cclxuXHRcdFx0XHRcdC8vIE5vIHRpbWVvdXQgYnkgZGVmYXVsdFxyXG5cdFx0XHRcdFx0bGV0IHRpbWVvdXQgPSAwO1xyXG5cclxuXHRcdFx0XHRcdC8vIEFjdHVhbCBmdW5jdGlvblxyXG5cdFx0XHRcdFx0ZnVuY3Rpb24gZHluYW1pYygpIHtcclxuXHRcdFx0XHRcdFx0Y29uc3QgYXJncyA9IFtdLnNsaWNlLmNhbGwoYXJndW1lbnRzKTtcclxuXHRcdFx0XHRcdFx0cmV0dXJuIENhbGwoW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXS5qb2luKCcuJyksIGFyZ3MsIHRpbWVvdXQpO1xyXG5cdFx0XHRcdFx0fVxyXG5cclxuXHRcdFx0XHRcdC8vIEFsbG93IHNldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxyXG5cdFx0XHRcdFx0ZHluYW1pYy5zZXRUaW1lb3V0ID0gZnVuY3Rpb24gKG5ld1RpbWVvdXQpIHtcclxuXHRcdFx0XHRcdFx0dGltZW91dCA9IG5ld1RpbWVvdXQ7XHJcblx0XHRcdFx0XHR9O1xyXG5cclxuXHRcdFx0XHRcdC8vIEFsbG93IGdldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxyXG5cdFx0XHRcdFx0ZHluYW1pYy5nZXRUaW1lb3V0ID0gZnVuY3Rpb24gKCkge1xyXG5cdFx0XHRcdFx0XHRyZXR1cm4gdGltZW91dDtcclxuXHRcdFx0XHRcdH07XHJcblxyXG5cdFx0XHRcdFx0cmV0dXJuIGR5bmFtaWM7XHJcblx0XHRcdFx0fSgpO1xyXG5cdFx0XHR9KTtcclxuXHRcdH0pO1xyXG5cdH0pO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5cclxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZCgpIHtcclxuICAgIHdpbmRvdy5sb2NhdGlvbi5yZWxvYWQoKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZEFwcCgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1InKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFN5c3RlbURlZmF1bHRUaGVtZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FTRFQnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldExpZ2h0VGhlbWUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBTFQnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldERhcmtUaGVtZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FEVCcpO1xyXG59XHJcblxyXG4vKipcclxuICogUGxhY2UgdGhlIHdpbmRvdyBpbiB0aGUgY2VudGVyIG9mIHRoZSBzY3JlZW5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0NlbnRlcigpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2MnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHdpbmRvdyB0aXRsZVxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdGl0bGVcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFRpdGxlKHRpdGxlKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dUJyArIHRpdGxlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1ha2VzIHRoZSB3aW5kb3cgZ28gZnVsbHNjcmVlblxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93RnVsbHNjcmVlbigpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0YnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldmVydHMgdGhlIHdpbmRvdyBmcm9tIGZ1bGxzY3JlZW5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VuZnVsbHNjcmVlbigpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2YnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIGluIGZ1bGwgc2NyZWVuIG1vZGUgb3Igbm90LlxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNGdWxsc2NyZWVuKCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNGdWxsc2NyZWVuXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFNpemUod2lkdGgsIGhlaWdodCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXczonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHt3OiBudW1iZXIsIGg6IG51bWJlcn0+fSBUaGUgc2l6ZSBvZiB0aGUgd2luZG93XHJcblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFNpemUoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dHZXRTaXplXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1daOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcclxufVxyXG5cclxuXHJcblxyXG4vKipcclxuICogU2V0IHRoZSB3aW5kb3cgQWx3YXlzT25Ub3Agb3Igbm90IG9uIHRvcFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QWx3YXlzT25Ub3AoYikge1xyXG5cclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FUUDonICsgKGIgPyAnMScgOiAnMCcpKTtcclxufVxyXG5cclxuXHJcblxyXG5cclxuLyoqXHJcbiAqIFNldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB5XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0UG9zaXRpb24oeCwgeSkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXcDonICsgeCArICc6JyArIHkpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTx7eDogbnVtYmVyLCB5OiBudW1iZXJ9Pn0gVGhlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRQb3NpdGlvbigpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0dldFBvc1wiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEhpZGUgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SGlkZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0gnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3cgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2hvdygpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1MnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1heGltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd01heGltaXNlKCkge1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXTScpO1xyXG59XHJcblxyXG4vKipcclxuICogVG9nZ2xlIHRoZSBNYXhpbWlzZSBvZiB0aGUgV2luZG93XHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dUb2dnbGVNYXhpbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3QnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFVubWF4aW1pc2UgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5tYXhpbWlzZSgpIHtcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1UnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1heGltaXNlZCBvciBub3QuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc01heGltaXNlZCgpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWF4aW1pc2VkXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogTWluaW1pc2UgdGhlIFdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWluaW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dtJyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBVbm1pbmltaXNlIHRoZSBXaW5kb3dcclxuICpcclxuICogQGV4cG9ydFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWluaW1pc2UoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d1Jyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBtaW5pbWlzZWQgb3Igbm90LlxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNNaW5pbWlzZWQoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc01pbmltaXNlZFwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG5vcm1hbCBvciBub3QuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc05vcm1hbCgpIHtcclxuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTm9ybWFsXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBSIFJlZFxyXG4gKiBAcGFyYW0ge251bWJlcn0gRyBHcmVlblxyXG4gKiBAcGFyYW0ge251bWJlcn0gQiBCbHVlXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBBIEFscGhhXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QmFja2dyb3VuZENvbG91cihSLCBHLCBCLCBBKSB7XHJcbiAgICBsZXQgcmdiYSA9IEpTT04uc3RyaW5naWZ5KHtyOiBSIHx8IDAsIGc6IEcgfHwgMCwgYjogQiB8fCAwLCBhOiBBIHx8IDI1NX0pO1xyXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXcjonICsgcmdiYSk7XHJcbn1cclxuXHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5cclxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xyXG5cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBhbGwgc2NyZWVucy4gQ2FsbCB0aGlzIGFuZXcgZWFjaCB0aW1lIHlvdSB3YW50IHRvIHJlZnJlc2ggZGF0YSBmcm9tIHRoZSB1bmRlcmx5aW5nIHdpbmRvd2luZyBzeXN0ZW0uXHJcbiAqIEBleHBvcnRcclxuICogQHR5cGVkZWYge2ltcG9ydCgnLi4vd3JhcHBlci9ydW50aW1lJykuU2NyZWVufSBTY3JlZW5cclxuICogQHJldHVybiB7UHJvbWlzZTx7U2NyZWVuW119Pn0gVGhlIHNjcmVlbnNcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTY3JlZW5HZXRBbGwoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpTY3JlZW5HZXRBbGxcIik7XHJcbn1cclxuIiwgIi8qKlxyXG4gKiBAZGVzY3JpcHRpb246IFVzZSB0aGUgc3lzdGVtIGRlZmF1bHQgYnJvd3NlciB0byBvcGVuIHRoZSB1cmxcclxuICogQHBhcmFtIHtzdHJpbmd9IHVybCBcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCcm93c2VyT3BlblVSTCh1cmwpIHtcclxuICB3aW5kb3cuV2FpbHNJbnZva2UoJ0JPOicgKyB1cmwpO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcclxuXHJcbi8qKlxyXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB0ZXh0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2xpcGJvYXJkU2V0VGV4dCh0ZXh0KSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRTZXRUZXh0XCIsIFt0ZXh0XSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIHRleHQgY29udGVudCBvZiB0aGUgY2xpcGJvYXJkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHJldHVybiB7UHJvbWlzZTx7c3RyaW5nfT59IFRleHQgY29udGVudCBvZiB0aGUgY2xpcGJvYXJkXHJcblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENsaXBib2FyZEdldFRleHQoKSB7XHJcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRHZXRUZXh0XCIpO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7RXZlbnRzT24sIEV2ZW50c09mZn0gZnJvbSBcIi4vZXZlbnRzXCI7XHJcblxyXG5jb25zdCBmbGFncyA9IHtcclxuICAgIHJlZ2lzdGVyZWQ6IGZhbHNlLFxyXG4gICAgZGVmYXVsdFVzZURyb3BUYXJnZXQ6IHRydWUsXHJcbiAgICB1c2VEcm9wVGFyZ2V0OiB0cnVlLFxyXG4gICAgbmV4dERlYWN0aXZhdGU6IG51bGwsXHJcbiAgICBuZXh0RGVhY3RpdmF0ZVRpbWVvdXQ6IG51bGwsXHJcbn07XHJcblxyXG5jb25zdCBEUk9QX1RBUkdFVF9BQ1RJVkUgPSBcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiO1xyXG5cclxuLyoqXHJcbiAqIGNoZWNrU3R5bGVEcm9wVGFyZ2V0IGNoZWNrcyBpZiB0aGUgc3R5bGUgaGFzIHRoZSBkcm9wIHRhcmdldCBhdHRyaWJ1dGVcclxuICogXHJcbiAqIEBwYXJhbSB7Q1NTU3R5bGVEZWNsYXJhdGlvbn0gc3R5bGUgXHJcbiAqIEByZXR1cm5zIFxyXG4gKi9cclxuZnVuY3Rpb24gY2hlY2tTdHlsZURyb3BUYXJnZXQoc3R5bGUpIHtcclxuICAgIGNvbnN0IGNzc0Ryb3BWYWx1ZSA9IHN0eWxlLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSkudHJpbSgpO1xyXG4gICAgaWYgKGNzc0Ryb3BWYWx1ZSkge1xyXG4gICAgICAgIGlmIChjc3NEcm9wVmFsdWUgPT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wVmFsdWUpIHtcclxuICAgICAgICAgICAgcmV0dXJuIHRydWU7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIC8vIGlmIHRoZSBlbGVtZW50IGhhcyB0aGUgZHJvcCB0YXJnZXQgYXR0cmlidXRlLCBidXQgXHJcbiAgICAgICAgLy8gdGhlIHZhbHVlIGlzIG5vdCBjb3JyZWN0LCB0ZXJtaW5hdGUgZmluZGluZyBwcm9jZXNzLlxyXG4gICAgICAgIC8vIFRoaXMgY2FuIGJlIHVzZWZ1bCB0byBibG9jayBzb21lIGNoaWxkIGVsZW1lbnRzIGZyb20gYmVpbmcgZHJvcCB0YXJnZXRzLlxyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgIH1cclxuICAgIHJldHVybiBmYWxzZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIG9uRHJhZ092ZXIgaXMgY2FsbGVkIHdoZW4gdGhlIGRyYWdvdmVyIGV2ZW50IGlzIGVtaXR0ZWQuXHJcbiAqIEBwYXJhbSB7RHJhZ0V2ZW50fSBlXHJcbiAqIEByZXR1cm5zXHJcbiAqL1xyXG5mdW5jdGlvbiBvbkRyYWdPdmVyKGUpIHtcclxuICAgIC8vIENoZWNrIGlmIHRoaXMgaXMgYW4gZXh0ZXJuYWwgZmlsZSBkcm9wIG9yIGludGVybmFsIEhUTUwgZHJhZ1xyXG4gICAgLy8gRXh0ZXJuYWwgZmlsZSBkcm9wcyB3aWxsIGhhdmUgXCJGaWxlc1wiIGluIHRoZSB0eXBlcyBhcnJheVxyXG4gICAgLy8gSW50ZXJuYWwgSFRNTCBkcmFncyB0eXBpY2FsbHkgaGF2ZSBcInRleHQvcGxhaW5cIiwgXCJ0ZXh0L2h0bWxcIiBvciBjdXN0b20gdHlwZXNcclxuICAgIGNvbnN0IGlzRmlsZURyb3AgPSBlLmRhdGFUcmFuc2Zlci50eXBlcy5pbmNsdWRlcyhcIkZpbGVzXCIpO1xyXG5cclxuICAgIC8vIE9ubHkgaGFuZGxlIGV4dGVybmFsIGZpbGUgZHJvcHMsIGxldCBpbnRlcm5hbCBIVE1MNSBkcmFnLWFuZC1kcm9wIHdvcmsgbm9ybWFsbHlcclxuICAgIGlmICghaXNGaWxlRHJvcCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBBTFdBWVMgcHJldmVudCBkZWZhdWx0IGZvciBmaWxlIGRyb3BzIHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICBlLnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICBlLmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xyXG5cclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3ApIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKCFmbGFncy51c2VEcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IGVsZW1lbnQgPSBlLnRhcmdldDtcclxuXHJcbiAgICAvLyBUcmlnZ2VyIGRlYm91bmNlIGZ1bmN0aW9uIHRvIGRlYWN0aXZhdGUgZHJvcCB0YXJnZXRzXHJcbiAgICBpZihmbGFncy5uZXh0RGVhY3RpdmF0ZSkgZmxhZ3MubmV4dERlYWN0aXZhdGUoKTtcclxuXHJcbiAgICAvLyBpZiB0aGUgZWxlbWVudCBpcyBudWxsIG9yIGVsZW1lbnQgaXMgbm90IGNoaWxkIG9mIGRyb3AgdGFyZ2V0IGVsZW1lbnRcclxuICAgIGlmICghZWxlbWVudCB8fCAhY2hlY2tTdHlsZURyb3BUYXJnZXQoZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KSkpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgbGV0IGN1cnJlbnRFbGVtZW50ID0gZWxlbWVudDtcclxuICAgIHdoaWxlIChjdXJyZW50RWxlbWVudCkge1xyXG4gICAgICAgIC8vIGNoZWNrIGlmIGN1cnJlbnRFbGVtZW50IGlzIGRyb3AgdGFyZ2V0IGVsZW1lbnRcclxuICAgICAgICBpZiAoY2hlY2tTdHlsZURyb3BUYXJnZXQoZ2V0Q29tcHV0ZWRTdHlsZShjdXJyZW50RWxlbWVudCkpKSB7XHJcbiAgICAgICAgICAgIGN1cnJlbnRFbGVtZW50LmNsYXNzTGlzdC5hZGQoRFJPUF9UQVJHRVRfQUNUSVZFKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgY3VycmVudEVsZW1lbnQgPSBjdXJyZW50RWxlbWVudC5wYXJlbnRFbGVtZW50O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogb25EcmFnTGVhdmUgaXMgY2FsbGVkIHdoZW4gdGhlIGRyYWdsZWF2ZSBldmVudCBpcyBlbWl0dGVkLlxyXG4gKiBAcGFyYW0ge0RyYWdFdmVudH0gZVxyXG4gKiBAcmV0dXJuc1xyXG4gKi9cclxuZnVuY3Rpb24gb25EcmFnTGVhdmUoZSkge1xyXG4gICAgLy8gQ2hlY2sgaWYgdGhpcyBpcyBhbiBleHRlcm5hbCBmaWxlIGRyb3Agb3IgaW50ZXJuYWwgSFRNTCBkcmFnXHJcbiAgICBjb25zdCBpc0ZpbGVEcm9wID0gZS5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoXCJGaWxlc1wiKTtcclxuXHJcbiAgICAvLyBPbmx5IGhhbmRsZSBleHRlcm5hbCBmaWxlIGRyb3BzLCBsZXQgaW50ZXJuYWwgSFRNTDUgZHJhZy1hbmQtZHJvcCB3b3JrIG5vcm1hbGx5XHJcbiAgICBpZiAoIWlzRmlsZURyb3ApIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQUxXQVlTIHByZXZlbnQgZGVmYXVsdCBmb3IgZmlsZSBkcm9wcyB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxyXG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG5cclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3ApIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKCFmbGFncy51c2VEcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIEZpbmQgdGhlIGNsb3NlIGRyb3AgdGFyZ2V0IGVsZW1lbnRcclxuICAgIGlmICghZS50YXJnZXQgfHwgIWNoZWNrU3R5bGVEcm9wVGFyZ2V0KGdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpKSkge1xyXG4gICAgICAgIHJldHVybiBudWxsO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFRyaWdnZXIgZGVib3VuY2UgZnVuY3Rpb24gdG8gZGVhY3RpdmF0ZSBkcm9wIHRhcmdldHNcclxuICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xyXG4gICAgXHJcbiAgICAvLyBVc2UgZGVib3VuY2UgdGVjaG5pcXVlIHRvIHRhY2xlIGRyYWdsZWF2ZSBldmVudHMgb24gb3ZlcmxhcHBpbmcgZWxlbWVudHMgYW5kIGRyb3AgdGFyZ2V0IGVsZW1lbnRzXHJcbiAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZSA9ICgpID0+IHtcclxuICAgICAgICAvLyBEZWFjdGl2YXRlIGFsbCBkcm9wIHRhcmdldHMsIG5ldyBkcm9wIHRhcmdldCB3aWxsIGJlIGFjdGl2YXRlZCBvbiBuZXh0IGRyYWdvdmVyIGV2ZW50XHJcbiAgICAgICAgQXJyYXkuZnJvbShkb2N1bWVudC5nZXRFbGVtZW50c0J5Q2xhc3NOYW1lKERST1BfVEFSR0VUX0FDVElWRSkpLmZvckVhY2goZWwgPT4gZWwuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkUpKTtcclxuICAgICAgICAvLyBSZXNldCBuZXh0RGVhY3RpdmF0ZVxyXG4gICAgICAgIGZsYWdzLm5leHREZWFjdGl2YXRlID0gbnVsbDtcclxuICAgICAgICAvLyBDbGVhciB0aW1lb3V0XHJcbiAgICAgICAgaWYgKGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCkge1xyXG4gICAgICAgICAgICBjbGVhclRpbWVvdXQoZmxhZ3MubmV4dERlYWN0aXZhdGVUaW1lb3V0KTtcclxuICAgICAgICAgICAgZmxhZ3MubmV4dERlYWN0aXZhdGVUaW1lb3V0ID0gbnVsbDtcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLy8gU2V0IHRpbWVvdXQgdG8gZGVhY3RpdmF0ZSBkcm9wIHRhcmdldHMgaWYgbm90IHRyaWdnZXJlZCBieSBuZXh0IGRyYWcgZXZlbnRcclxuICAgIGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCA9IHNldFRpbWVvdXQoKCkgPT4ge1xyXG4gICAgICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xyXG4gICAgfSwgNTApO1xyXG59XHJcblxyXG4vKipcclxuICogb25Ecm9wIGlzIGNhbGxlZCB3aGVuIHRoZSBkcm9wIGV2ZW50IGlzIGVtaXR0ZWQuXHJcbiAqIEBwYXJhbSB7RHJhZ0V2ZW50fSBlXHJcbiAqIEByZXR1cm5zXHJcbiAqL1xyXG5mdW5jdGlvbiBvbkRyb3AoZSkge1xyXG4gICAgLy8gQ2hlY2sgaWYgdGhpcyBpcyBhbiBleHRlcm5hbCBmaWxlIGRyb3Agb3IgaW50ZXJuYWwgSFRNTCBkcmFnXHJcbiAgICBjb25zdCBpc0ZpbGVEcm9wID0gZS5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoXCJGaWxlc1wiKTtcclxuXHJcbiAgICAvLyBPbmx5IGhhbmRsZSBleHRlcm5hbCBmaWxlIGRyb3BzLCBsZXQgaW50ZXJuYWwgSFRNTDUgZHJhZy1hbmQtZHJvcCB3b3JrIG5vcm1hbGx5XHJcbiAgICBpZiAoIWlzRmlsZURyb3ApIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQUxXQVlTIHByZXZlbnQgZGVmYXVsdCBmb3IgZmlsZSBkcm9wcyB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxyXG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG5cclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3ApIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKENhblJlc29sdmVGaWxlUGF0aHMoKSkge1xyXG4gICAgICAgIC8vIHByb2Nlc3MgZmlsZXNcclxuICAgICAgICBsZXQgZmlsZXMgPSBbXTtcclxuICAgICAgICBpZiAoZS5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcclxuICAgICAgICAgICAgZmlsZXMgPSBbLi4uZS5kYXRhVHJhbnNmZXIuaXRlbXNdLm1hcCgoaXRlbSwgaSkgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKGl0ZW0ua2luZCA9PT0gJ2ZpbGUnKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuIGl0ZW0uZ2V0QXNGaWxlKCk7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0pO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIGZpbGVzID0gWy4uLmUuZGF0YVRyYW5zZmVyLmZpbGVzXTtcclxuICAgICAgICB9XHJcbiAgICAgICAgd2luZG93LnJ1bnRpbWUuUmVzb2x2ZUZpbGVQYXRocyhlLngsIGUueSwgZmlsZXMpO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICghZmxhZ3MudXNlRHJvcFRhcmdldCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBUcmlnZ2VyIGRlYm91bmNlIGZ1bmN0aW9uIHRvIGRlYWN0aXZhdGUgZHJvcCB0YXJnZXRzXHJcbiAgICBpZihmbGFncy5uZXh0RGVhY3RpdmF0ZSkgZmxhZ3MubmV4dERlYWN0aXZhdGUoKTtcclxuXHJcbiAgICAvLyBEZWFjdGl2YXRlIGFsbCBkcm9wIHRhcmdldHNcclxuICAgIEFycmF5LmZyb20oZG9jdW1lbnQuZ2V0RWxlbWVudHNCeUNsYXNzTmFtZShEUk9QX1RBUkdFVF9BQ1RJVkUpKS5mb3JFYWNoKGVsID0+IGVsLmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFKSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyBjaGVja3MgdGhlIGJyb3dzZXIncyBjYXBhYmlsaXR5IG9mIHNlbmRpbmcgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHNcclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59XHJcbiAqIEBjb25zdHJ1Y3RvclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhblJlc29sdmVGaWxlUGF0aHMoKSB7XHJcbiAgICByZXR1cm4gd2luZG93LmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgIT0gbnVsbDtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlc29sdmVGaWxlUGF0aHMgc2VuZHMgZHJvcCBldmVudHMgdG8gdGhlIEdPIHNpZGUgdG8gcmVzb2x2ZSBmaWxlIHBhdGhzIG9uIHdpbmRvd3MuXHJcbiAqXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB5XHJcbiAqIEBwYXJhbSB7YW55W119IGZpbGVzXHJcbiAqIEBjb25zdHJ1Y3RvclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFJlc29sdmVGaWxlUGF0aHMoeCwgeSwgZmlsZXMpIHtcclxuICAgIC8vIE9ubHkgZm9yIHdpbmRvd3Mgd2VidmlldzIgPj0gMS4wLjE3NzQuMzBcclxuICAgIC8vIGh0dHBzOi8vbGVhcm4ubWljcm9zb2Z0LmNvbS9lbi11cy9taWNyb3NvZnQtZWRnZS93ZWJ2aWV3Mi9yZWZlcmVuY2Uvd2luMzIvaWNvcmV3ZWJ2aWV3MndlYm1lc3NhZ2VyZWNlaXZlZGV2ZW50YXJnczI/dmlldz13ZWJ2aWV3Mi0xLjAuMTgyMy4zMiNhcHBsaWVzLXRvXHJcbiAgICBpZiAod2luZG93LmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMpIHtcclxuICAgICAgICBjaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyhgZmlsZTpkcm9wOiR7eH06JHt5fWAsIGZpbGVzKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxiYWNrIGZvciBPbkZpbGVEcm9wIHJldHVybnMgYSBzbGljZSBvZiBmaWxlIHBhdGggc3RyaW5ncyB3aGVuIGEgZHJvcCBpcyBmaW5pc2hlZC5cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAY2FsbGJhY2sgT25GaWxlRHJvcENhbGxiYWNrXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB4IC0geCBjb29yZGluYXRlIG9mIHRoZSBkcm9wXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB5IC0geSBjb29yZGluYXRlIG9mIHRoZSBkcm9wXHJcbiAqIEBwYXJhbSB7c3RyaW5nW119IHBhdGhzIC0gQSBsaXN0IG9mIGZpbGUgcGF0aHMuXHJcbiAqL1xyXG5cclxuLyoqXHJcbiAqIE9uRmlsZURyb3AgbGlzdGVucyB0byBkcmFnIGFuZCBkcm9wIGV2ZW50cyBhbmQgY2FsbHMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGNvb3JkaW5hdGVzIG9mIHRoZSBkcm9wIGFuZCBhbiBhcnJheSBvZiBwYXRoIHN0cmluZ3MuXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtPbkZpbGVEcm9wQ2FsbGJhY2t9IGNhbGxiYWNrIC0gQ2FsbGJhY2sgZm9yIE9uRmlsZURyb3AgcmV0dXJucyBhIHNsaWNlIG9mIGZpbGUgcGF0aCBzdHJpbmdzIHdoZW4gYSBkcm9wIGlzIGZpbmlzaGVkLlxyXG4gKiBAcGFyYW0ge2Jvb2xlYW59IFt1c2VEcm9wVGFyZ2V0PXRydWVdIC0gT25seSBjYWxsIHRoZSBjYWxsYmFjayB3aGVuIHRoZSBkcm9wIGZpbmlzaGVkIG9uIGFuIGVsZW1lbnQgdGhhdCBoYXMgdGhlIGRyb3AgdGFyZ2V0IHN0eWxlLiAoLS13YWlscy1kcm9wLXRhcmdldClcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbkZpbGVEcm9wKGNhbGxiYWNrLCB1c2VEcm9wVGFyZ2V0KSB7XHJcbiAgICBpZiAodHlwZW9mIGNhbGxiYWNrICE9PSBcImZ1bmN0aW9uXCIpIHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKFwiRHJhZ0FuZERyb3BDYWxsYmFjayBpcyBub3QgYSBmdW5jdGlvblwiKTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGZsYWdzLnJlZ2lzdGVyZWQpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcbiAgICBmbGFncy5yZWdpc3RlcmVkID0gdHJ1ZTtcclxuXHJcbiAgICBjb25zdCB1RFRQVCA9IHR5cGVvZiB1c2VEcm9wVGFyZ2V0O1xyXG4gICAgZmxhZ3MudXNlRHJvcFRhcmdldCA9IHVEVFBUID09PSBcInVuZGVmaW5lZFwiIHx8IHVEVFBUICE9PSBcImJvb2xlYW5cIiA/IGZsYWdzLmRlZmF1bHRVc2VEcm9wVGFyZ2V0IDogdXNlRHJvcFRhcmdldDtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcmFnb3ZlcicsIG9uRHJhZ092ZXIpO1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdsZWF2ZScsIG9uRHJhZ0xlYXZlKTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgb25Ecm9wKTtcclxuXHJcbiAgICBsZXQgY2IgPSBjYWxsYmFjaztcclxuICAgIGlmIChmbGFncy51c2VEcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgY2IgPSBmdW5jdGlvbiAoeCwgeSwgcGF0aHMpIHtcclxuICAgICAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSlcclxuICAgICAgICAgICAgLy8gaWYgdGhlIGVsZW1lbnQgaXMgbnVsbCBvciBlbGVtZW50IGlzIG5vdCBjaGlsZCBvZiBkcm9wIHRhcmdldCBlbGVtZW50LCByZXR1cm4gbnVsbFxyXG4gICAgICAgICAgICBpZiAoIWVsZW1lbnQgfHwgIWNoZWNrU3R5bGVEcm9wVGFyZ2V0KGdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkpKSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gbnVsbDtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBjYWxsYmFjayh4LCB5LCBwYXRocyk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIEV2ZW50c09uKFwid2FpbHM6ZmlsZS1kcm9wXCIsIGNiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE9uRmlsZURyb3BPZmYgcmVtb3ZlcyB0aGUgZHJhZyBhbmQgZHJvcCBsaXN0ZW5lcnMgYW5kIGhhbmRsZXJzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uRmlsZURyb3BPZmYoKSB7XHJcbiAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBvbkRyYWdPdmVyKTtcclxuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCBvbkRyYWdMZWF2ZSk7XHJcbiAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignZHJvcCcsIG9uRHJvcCk7XHJcbiAgICBFdmVudHNPZmYoXCJ3YWlsczpmaWxlLWRyb3BcIik7XHJcbiAgICBmbGFncy5yZWdpc3RlcmVkID0gZmFsc2U7XHJcbn1cclxuIiwgIi8qXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG5cclxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcclxuKi9cclxuZXhwb3J0IGZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcclxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgIGNvbnN0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XHJcbiAgICBjb25zdCBjb21wdXRlZFN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XHJcbiAgICBjb25zdCBkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24gPSBjb21wdXRlZFN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpO1xyXG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcclxuICAgICAgICBjYXNlIFwic2hvd1wiOlxyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcclxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGRlZmF1bHQ6XHJcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXHJcbiAgICAgICAgICAgIGlmIChlbGVtZW50LmlzQ29udGVudEVkaXRhYmxlKSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgYW5kIGFjdGlvbiBpcyBvbiB0aGUgc2VsZWN0ZWQgZWxlbWVudHNcclxuICAgICAgICAgICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xyXG4gICAgICAgICAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSAoc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMClcclxuICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbikge1xyXG4gICAgICAgICAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcclxuICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0cyA9IHJhbmdlLmdldENsaWVudFJlY3RzKCk7XHJcbiAgICAgICAgICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHJlY3QubGVmdCwgcmVjdC50b3ApID09PSBlbGVtZW50KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGFnbmFtZSBpcyBpbnB1dCBvciB0ZXh0YXJlYVxyXG4gICAgICAgICAgICBpZiAoZWxlbWVudC50YWdOYW1lID09PSBcIklOUFVUXCIgfHwgZWxlbWVudC50YWdOYW1lID09PSBcIlRFWFRBUkVBXCIpIHtcclxuICAgICAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24gfHwgKCFlbGVtZW50LnJlYWRPbmx5ICYmICFlbGVtZW50LmRpc2FibGVkKSkge1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcclxuaW1wb3J0IHtcclxuICBldmVudExpc3RlbmVycyxcclxuICBFdmVudHNFbWl0LFxyXG4gIEV2ZW50c05vdGlmeSxcclxuICBFdmVudHNPZmYsXHJcbiAgRXZlbnRzT2ZmQWxsLFxyXG4gIEV2ZW50c09uLFxyXG4gIEV2ZW50c09uY2UsXHJcbiAgRXZlbnRzT25NdWx0aXBsZSxcclxufSBmcm9tIFwiLi9ldmVudHNcIjtcclxuaW1wb3J0IHsgQ2FsbCwgQ2FsbGJhY2ssIGNhbGxiYWNrcyB9IGZyb20gJy4vY2FsbHMnO1xyXG5pbXBvcnQgeyBTZXRCaW5kaW5ncyB9IGZyb20gXCIuL2JpbmRpbmdzXCI7XHJcbmltcG9ydCAqIGFzIFdpbmRvdyBmcm9tIFwiLi93aW5kb3dcIjtcclxuaW1wb3J0ICogYXMgU2NyZWVuIGZyb20gXCIuL3NjcmVlblwiO1xyXG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXJcIjtcclxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZFwiO1xyXG5pbXBvcnQgKiBhcyBEcmFnQW5kRHJvcCBmcm9tIFwiLi9kcmFnYW5kZHJvcFwiO1xyXG5pbXBvcnQgKiBhcyBDb250ZXh0TWVudSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1EnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1MnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0gnKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xyXG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6RW52aXJvbm1lbnRcIik7XHJcbn1cclxuXHJcbi8vIFRoZSBKUyBydW50aW1lXHJcbndpbmRvdy5ydW50aW1lID0ge1xyXG4gICAgLi4uTG9nLFxyXG4gICAgLi4uV2luZG93LFxyXG4gICAgLi4uQnJvd3NlcixcclxuICAgIC4uLlNjcmVlbixcclxuICAgIC4uLkNsaXBib2FyZCxcclxuICAgIC4uLkRyYWdBbmREcm9wLFxyXG4gICAgRXZlbnRzT24sXHJcbiAgICBFdmVudHNPbmNlLFxyXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcclxuICAgIEV2ZW50c0VtaXQsXHJcbiAgICBFdmVudHNPZmYsXHJcbiAgICBFdmVudHNPZmZBbGwsXHJcbiAgICBFbnZpcm9ubWVudCxcclxuICAgIFNob3csXHJcbiAgICBIaWRlLFxyXG4gICAgUXVpdFxyXG59O1xyXG5cclxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXHJcbndpbmRvdy53YWlscyA9IHtcclxuICAgIENhbGxiYWNrLFxyXG4gICAgRXZlbnRzTm90aWZ5LFxyXG4gICAgU2V0QmluZGluZ3MsXHJcbiAgICBldmVudExpc3RlbmVycyxcclxuICAgIGNhbGxiYWNrcyxcclxuICAgIGZsYWdzOiB7XHJcbiAgICAgICAgZGlzYWJsZVNjcm9sbGJhckRyYWc6IGZhbHNlLFxyXG4gICAgICAgIGRpc2FibGVEZWZhdWx0Q29udGV4dE1lbnU6IGZhbHNlLFxyXG4gICAgICAgIGVuYWJsZVJlc2l6ZTogZmFsc2UsXHJcbiAgICAgICAgZGVmYXVsdEN1cnNvcjogbnVsbCxcclxuICAgICAgICBib3JkZXJUaGlja25lc3M6IDYsXHJcbiAgICAgICAgc2hvdWxkRHJhZzogZmFsc2UsXHJcbiAgICAgICAgZGVmZXJEcmFnVG9Nb3VzZU1vdmU6IHRydWUsXHJcbiAgICAgICAgY3NzRHJhZ1Byb3BlcnR5OiBcIi0td2FpbHMtZHJhZ2dhYmxlXCIsXHJcbiAgICAgICAgY3NzRHJhZ1ZhbHVlOiBcImRyYWdcIixcclxuICAgICAgICBjc3NEcm9wUHJvcGVydHk6IFwiLS13YWlscy1kcm9wLXRhcmdldFwiLFxyXG4gICAgICAgIGNzc0Ryb3BWYWx1ZTogXCJkcm9wXCIsXHJcbiAgICAgICAgZW5hYmxlV2FpbHNEcmFnQW5kRHJvcDogZmFsc2UsXHJcbiAgICB9XHJcbn07XHJcblxyXG4vLyBTZXQgdGhlIGJpbmRpbmdzXHJcbmlmICh3aW5kb3cud2FpbHNiaW5kaW5ncykge1xyXG4gICAgd2luZG93LndhaWxzLlNldEJpbmRpbmdzKHdpbmRvdy53YWlsc2JpbmRpbmdzKTtcclxuICAgIGRlbGV0ZSB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3M7XHJcbn1cclxuXHJcbi8vIChib29sKSBUaGlzIGlzIGV2YWx1YXRlZCBhdCBidWlsZCB0aW1lIGluIHBhY2thZ2UuanNvblxyXG5pZiAoIURFQlVHKSB7XHJcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XHJcbn1cclxuXHJcbmxldCBkcmFnVGVzdCA9IGZ1bmN0aW9uKGUpIHtcclxuICAgIHZhciB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZSh3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1Byb3BlcnR5KTtcclxuICAgIGlmICh2YWwpIHtcclxuICAgICAgICB2YWwgPSB2YWwudHJpbSgpO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICh2YWwgIT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGUuYnV0dG9ucyAhPT0gMSkge1xyXG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBub3QgdGhlIHByaW1hcnkgYnV0dG9uIGhhcyBiZWVuIGNsaWNrZWQuXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChlLmRldGFpbCAhPT0gMSkge1xyXG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBtb3JlIHRoYW4gb25jZSBoYXMgYmVlbiBjbGlja2VkLCBlLmcuIHdoZW4gZG91YmxlIGNsaWNraW5nXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiB0cnVlO1xyXG59O1xyXG5cclxud2luZG93LndhaWxzLnNldENTU0RyYWdQcm9wZXJ0aWVzID0gZnVuY3Rpb24ocHJvcGVydHksIHZhbHVlKSB7XHJcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1Byb3BlcnR5ID0gcHJvcGVydHk7XHJcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1ZhbHVlID0gdmFsdWU7XHJcbn1cclxuXHJcbndpbmRvdy53YWlscy5zZXRDU1NEcm9wUHJvcGVydGllcyA9IGZ1bmN0aW9uKHByb3BlcnR5LCB2YWx1ZSkge1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSA9IHByb3BlcnR5O1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BWYWx1ZSA9IHZhbHVlO1xyXG59XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgKGUpID0+IHtcclxuICAgIC8vIENoZWNrIGZvciByZXNpemluZ1xyXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKSB7XHJcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwicmVzaXplOlwiICsgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpO1xyXG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGRyYWdUZXN0KGUpKSB7XHJcbiAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlU2Nyb2xsYmFyRHJhZykge1xyXG4gICAgICAgICAgICAvLyBUaGlzIGNoZWNrcyBmb3IgY2xpY2tzIG9uIHRoZSBzY3JvbGwgYmFyXHJcbiAgICAgICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgICAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRlZmVyRHJhZ1RvTW91c2VNb3ZlKSB7XHJcbiAgICAgICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gdHJ1ZTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBlLnByZXZlbnREZWZhdWx0KClcclxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG4gICAgfVxyXG59KTtcclxuXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgKCkgPT4ge1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcclxufSk7XHJcblxyXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XHJcbiAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gY3Vyc29yIHx8IHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yO1xyXG4gICAgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgPSBjdXJzb3I7XHJcbn1cclxuXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBmdW5jdGlvbihlKSB7XHJcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcpIHtcclxuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XHJcbiAgICAgICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcclxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVJlc2l6ZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XHJcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yO1xyXG4gICAgfVxyXG4gICAgaWYgKHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcyAmJiB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzKSB7XHJcbiAgICAgICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IFwic2UtcmVzaXplXCI7XHJcbiAgICB9XHJcbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XHJcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcclxuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xyXG5cclxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XHJcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9IGVsc2UgaWYgKHJpZ2h0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgdG9wQm9yZGVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIgJiYgcmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XHJcblxyXG59KTtcclxuXHJcbi8vIFNldHVwIGNvbnRleHQgbWVudSBob29rXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGZ1bmN0aW9uKGUpIHtcclxuICAgIC8vIGFsd2F5cyBzaG93IHRoZSBjb250ZXh0bWVudSBpbiBkZWJ1ZyAmIGRldlxyXG4gICAgaWYgKERFQlVHKSByZXR1cm47XHJcblxyXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlRGVmYXVsdENvbnRleHRNZW51KSB7XHJcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBDb250ZXh0TWVudS5wcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGUpO1xyXG4gICAgfVxyXG59KTtcclxuXHJcbndpbmRvdy5XYWlsc0ludm9rZShcInJ1bnRpbWU6cmVhZHlcIik7Il0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFrQkEsV0FBUyxlQUFlLE9BQU8sU0FBUztBQUl2QyxXQUFPLFlBQVksTUFBTSxRQUFRLE9BQU87QUFBQSxFQUN6QztBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsUUFBUSxTQUFTO0FBQ2hDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxXQUFXLFNBQVM7QUFDbkMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxZQUFZLFVBQVU7QUFDckMsbUJBQWUsS0FBSyxRQUFRO0FBQUEsRUFDN0I7QUFHTyxNQUFNLFdBQVc7QUFBQSxJQUN2QixPQUFPO0FBQUEsSUFDUCxPQUFPO0FBQUEsSUFDUCxNQUFNO0FBQUEsSUFDTixTQUFTO0FBQUEsSUFDVCxPQUFPO0FBQUEsRUFDUjs7O0FDOUZBLE1BQU0sV0FBTixNQUFlO0FBQUEsSUFRWCxZQUFZLFdBQVcsVUFBVSxjQUFjO0FBQzNDLFdBQUssWUFBWTtBQUVqQixXQUFLLGVBQWUsZ0JBQWdCO0FB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
