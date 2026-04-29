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
    const listeners = eventListeners[eventName]?.slice() || [];
    if (listeners.length) {
      const toRemove = /* @__PURE__ */ new Set();
      for (let count = listeners.length - 1; count >= 0; count -= 1) {
        const listener = listeners[count];
        let data = eventData.data;
        const destroy = listener.Callback(data);
        if (destroy) {
          toRemove.add(listener);
        }
      }
      if (toRemove.size > 0 && eventListeners[eventName]) {
        eventListeners[eventName] = eventListeners[eventName].filter((l) => !toRemove.has(l));
        if (eventListeners[eventName].length === 0) {
          removeListener(eventName);
        }
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

  // desktop/notifications.js
  var notifications_exports = {};
  __export(notifications_exports, {
    CheckNotificationAuthorization: () => CheckNotificationAuthorization,
    CleanupNotifications: () => CleanupNotifications,
    InitializeNotifications: () => InitializeNotifications,
    IsNotificationAvailable: () => IsNotificationAvailable,
    RegisterNotificationCategory: () => RegisterNotificationCategory,
    RemoveAllDeliveredNotifications: () => RemoveAllDeliveredNotifications,
    RemoveAllPendingNotifications: () => RemoveAllPendingNotifications,
    RemoveDeliveredNotification: () => RemoveDeliveredNotification,
    RemoveNotification: () => RemoveNotification,
    RemoveNotificationCategory: () => RemoveNotificationCategory,
    RemovePendingNotification: () => RemovePendingNotification,
    RequestNotificationAuthorization: () => RequestNotificationAuthorization,
    SendNotification: () => SendNotification,
    SendNotificationWithActions: () => SendNotificationWithActions
  });
  function InitializeNotifications() {
    return Call(":wails:InitializeNotifications");
  }
  function CleanupNotifications() {
    return Call(":wails:CleanupNotifications");
  }
  function IsNotificationAvailable() {
    return Call(":wails:IsNotificationAvailable");
  }
  function RequestNotificationAuthorization() {
    return Call(":wails:RequestNotificationAuthorization");
  }
  function CheckNotificationAuthorization() {
    return Call(":wails:CheckNotificationAuthorization");
  }
  function SendNotification(options) {
    return Call(":wails:SendNotification", [options]);
  }
  function SendNotificationWithActions(options) {
    return Call(":wails:SendNotificationWithActions", [options]);
  }
  function RegisterNotificationCategory(category) {
    return Call(":wails:RegisterNotificationCategory", [category]);
  }
  function RemoveNotificationCategory(categoryId) {
    return Call(":wails:RemoveNotificationCategory", [categoryId]);
  }
  function RemoveAllPendingNotifications() {
    return Call(":wails:RemoveAllPendingNotifications");
  }
  function RemovePendingNotification(identifier) {
    return Call(":wails:RemovePendingNotification", [identifier]);
  }
  function RemoveAllDeliveredNotifications() {
    return Call(":wails:RemoveAllDeliveredNotifications");
  }
  function RemoveDeliveredNotification(identifier) {
    return Call(":wails:RemoveDeliveredNotification", [identifier]);
  }
  function RemoveNotification(identifier) {
    return Call(":wails:RemoveNotification", [identifier]);
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
    ...notifications_exports,
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL2RyYWdhbmRkcm9wLmpzIiwgImRlc2t0b3AvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9ub3RpZmljYXRpb25zLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfICAgICAgIF9fICAgICAgXyBfX1xufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuLyoqXG4gKiBTZW5kcyBhIGxvZyBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHdpdGggdGhlIGdpdmVuIGxldmVsICsgbWVzc2FnZVxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBsZXZlbFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZnVuY3Rpb24gc2VuZExvZ01lc3NhZ2UobGV2ZWwsIG1lc3NhZ2UpIHtcblxuXHQvLyBMb2cgTWVzc2FnZSBmb3JtYXQ6XG5cdC8vIGxbdHlwZV1bbWVzc2FnZV1cblx0d2luZG93LldhaWxzSW52b2tlKCdMJyArIGxldmVsICsgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiB0cmFjZSBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nVHJhY2UobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnVCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1ByaW50KG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1AnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGRlYnVnIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dEZWJ1ZyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdEJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBpbmZvIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dJbmZvKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0knLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIHdhcm5pbmcgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1dhcm5pbmcobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnVycsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gZXJyb3IgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ0Vycm9yKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0UnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGZhdGFsIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dGYXRhbChtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdGJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgTG9nIGxldmVsIHRvIHRoZSBnaXZlbiBsb2cgbGV2ZWxcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gbG9nbGV2ZWxcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldExvZ0xldmVsKGxvZ2xldmVsKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdTJywgbG9nbGV2ZWwpO1xufVxuXG4vLyBMb2cgbGV2ZWxzXG5leHBvcnQgY29uc3QgTG9nTGV2ZWwgPSB7XG5cdFRSQUNFOiAxLFxuXHRERUJVRzogMixcblx0SU5GTzogMyxcblx0V0FSTklORzogNCxcblx0RVJST1I6IDUsXG59O1xuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XG5cbiAgICAvLyBHZXQgdGhlIGV2ZW50IG5hbWVcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XG5cbiAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lcnMgdG8gZGVzdHJveSAobWF4Q2FsbGJhY2tzIHJlYWNoZWQpXG4gICAgY29uc3QgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXT8uc2xpY2UoKSB8fCBbXTtcblxuICAgIC8vIENoZWNrIGlmIHdlIGhhdmUgYW55IGxpc3RlbmVycyBmb3IgdGhpcyBldmVudFxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoKSB7XG5cbiAgICAgICAgY29uc3QgdG9SZW1vdmUgPSBuZXcgU2V0KCk7XG5cbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSBsaXN0ZW5lcnMubGVuZ3RoIC0gMTsgY291bnQgPj0gMDsgY291bnQgLT0gMSkge1xuXG4gICAgICAgICAgICAvLyBHZXQgbmV4dCBsaXN0ZW5lclxuICAgICAgICAgICAgY29uc3QgbGlzdGVuZXIgPSBsaXN0ZW5lcnNbY291bnRdO1xuXG4gICAgICAgICAgICBsZXQgZGF0YSA9IGV2ZW50RGF0YS5kYXRhO1xuXG4gICAgICAgICAgICAvLyBEbyB0aGUgY2FsbGJhY2tcbiAgICAgICAgICAgIGNvbnN0IGRlc3Ryb3kgPSBsaXN0ZW5lci5DYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIGlmIChkZXN0cm95KSB7XG4gICAgICAgICAgICAgICAgdG9SZW1vdmUuYWRkKGxpc3RlbmVyKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFJlbW92ZSBkZXN0cm95ZWQgbGlzdGVuZXJzIGZyb20gdGhlIGN1cnJlbnQgbGl2ZSBsaXN0XG4gICAgICAgIC8vIFRoaXMgcmVzcGVjdHMgYW55IHJlbW92YWxzIHRoYXQgaGFwcGVuZWQgdmlhIGxpc3RlbmVyT2ZmIGR1cmluZyBjYWxsYmFja3NcbiAgICAgICAgaWYgKHRvUmVtb3ZlLnNpemUgPiAwICYmIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0pIHtcbiAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmZpbHRlcihsID0+ICF0b1JlbW92ZS5oYXMobCkpO1xuICAgICAgICAgICAgaWYgKGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0ubGVuZ3RoID09PSAwKSB7XG4gICAgICAgICAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBOb3RpZnkgaW5mb3JtcyBmcm9udGVuZCBsaXN0ZW5lcnMgdGhhdCBhbiBldmVudCB3YXMgZW1pdHRlZCB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5vdGlmeU1lc3NhZ2UgLSBlbmNvZGVkIG5vdGlmaWNhdGlvbiBtZXNzYWdlXG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c05vdGlmeShub3RpZnlNZXNzYWdlKSB7XG4gICAgLy8gUGFyc2UgdGhlIG1lc3NhZ2VcbiAgICBsZXQgbWVzc2FnZTtcbiAgICB0cnkge1xuICAgICAgICBtZXNzYWdlID0gSlNPTi5wYXJzZShub3RpZnlNZXNzYWdlKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGNvbnN0IGVycm9yID0gJ0ludmFsaWQgSlNPTiBwYXNzZWQgdG8gTm90aWZ5OiAnICsgbm90aWZ5TWVzc2FnZTtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGVycm9yKTtcbiAgICB9XG4gICAgbm90aWZ5TGlzdGVuZXJzKG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIEVtaXQgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c0VtaXQoZXZlbnROYW1lKSB7XG5cbiAgICBjb25zdCBwYXlsb2FkID0ge1xuICAgICAgICBuYW1lOiBldmVudE5hbWUsXG4gICAgICAgIGRhdGE6IFtdLnNsaWNlLmFwcGx5KGFyZ3VtZW50cykuc2xpY2UoMSksXG4gICAgfTtcblxuICAgIC8vIE5vdGlmeSBKUyBsaXN0ZW5lcnNcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XG5cbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG59XG5cbmZ1bmN0aW9uIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSkge1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lcnNcbiAgICBkZWxldGUgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXTtcblxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VYJyArIGV2ZW50TmFtZSk7XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVyZXMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKVxuXG4gICAgaWYgKGFkZGl0aW9uYWxFdmVudE5hbWVzLmxlbmd0aCA+IDApIHtcbiAgICAgICAgYWRkaXRpb25hbEV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKVxuICAgICAgICB9KVxuICAgIH1cbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYWxsIGV2ZW50IGxpc3RlbmVycyBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICovXG4gZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZkFsbCgpIHtcbiAgICBjb25zdCBldmVudE5hbWVzID0gT2JqZWN0LmtleXMoZXZlbnRMaXN0ZW5lcnMpO1xuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG4gICAgfSlcbn1cblxuLyoqXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIEV2ZW50c09uXG4gKlxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcbiAqL1xuIGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdID09PSB1bmRlZmluZWQpIHJldHVybjtcblxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcblxuICAgIC8vIENsZWFuIHVwIGlmIHRoZXJlIGFyZSBubyBldmVudCBsaXN0ZW5lcnMgbGVmdFxuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpO1xuICAgIH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuZXhwb3J0IGNvbnN0IGNhbGxiYWNrcyA9IHt9O1xuXG4vKipcbiAqIFJldHVybnMgYSBudW1iZXIgZnJvbSB0aGUgbmF0aXZlIGJyb3dzZXIgcmFuZG9tIGZ1bmN0aW9uXG4gKlxuICogQHJldHVybnMgbnVtYmVyXG4gKi9cbmZ1bmN0aW9uIGNyeXB0b1JhbmRvbSgpIHtcblx0dmFyIGFycmF5ID0gbmV3IFVpbnQzMkFycmF5KDEpO1xuXHRyZXR1cm4gd2luZG93LmNyeXB0by5nZXRSYW5kb21WYWx1ZXMoYXJyYXkpWzBdO1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBudW1iZXIgdXNpbmcgZGEgb2xkLXNrb29sIE1hdGguUmFuZG9tXG4gKiBJIGxpa2VzIHRvIGNhbGwgaXQgTE9MUmFuZG9tXG4gKlxuICogQHJldHVybnMgbnVtYmVyXG4gKi9cbmZ1bmN0aW9uIGJhc2ljUmFuZG9tKCkge1xuXHRyZXR1cm4gTWF0aC5yYW5kb20oKSAqIDkwMDcxOTkyNTQ3NDA5OTE7XG59XG5cbi8vIFBpY2sgYSByYW5kb20gbnVtYmVyIGZ1bmN0aW9uIGJhc2VkIG9uIGJyb3dzZXIgY2FwYWJpbGl0eVxudmFyIHJhbmRvbUZ1bmM7XG5pZiAod2luZG93LmNyeXB0bykge1xuXHRyYW5kb21GdW5jID0gY3J5cHRvUmFuZG9tO1xufSBlbHNlIHtcblx0cmFuZG9tRnVuYyA9IGJhc2ljUmFuZG9tO1xufVxuXG5cbi8qKlxuICogQ2FsbCBzZW5kcyBhIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgdG8gY2FsbCB0aGUgYmluZGluZyB3aXRoIHRoZVxuICogZ2l2ZW4gZGF0YS4gQSBwcm9taXNlIGlzIHJldHVybmVkIGFuZCB3aWxsIGJlIGNvbXBsZXRlZCB3aGVuIHRoZVxuICogYmFja2VuZCByZXNwb25kcy4gVGhpcyB3aWxsIGJlIHJlc29sdmVkIHdoZW4gdGhlIGNhbGwgd2FzIHN1Y2Nlc3NmdWxcbiAqIG9yIHJlamVjdGVkIGlmIGFuIGVycm9yIGlzIHBhc3NlZCBiYWNrLlxuICogVGhlcmUgaXMgYSB0aW1lb3V0IG1lY2hhbmlzbS4gSWYgdGhlIGNhbGwgZG9lc24ndCByZXNwb25kIGluIHRoZSBnaXZlblxuICogdGltZSAoaW4gbWlsbGlzZWNvbmRzKSB0aGVuIHRoZSBwcm9taXNlIGlzIHJlamVjdGVkLlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lXG4gKiBAcGFyYW0ge2FueT19IGFyZ3NcbiAqIEBwYXJhbSB7bnVtYmVyPX0gdGltZW91dFxuICogQHJldHVybnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGwobmFtZSwgYXJncywgdGltZW91dCkge1xuXG5cdC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxuXHRpZiAodGltZW91dCA9PSBudWxsKSB7XG5cdFx0dGltZW91dCA9IDA7XG5cdH1cblxuXHQvLyBDcmVhdGUgYSBwcm9taXNlXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cblx0XHQvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuXHRcdHZhciBjYWxsYmFja0lEO1xuXHRcdGRvIHtcblx0XHRcdGNhbGxiYWNrSUQgPSBuYW1lICsgJy0nICsgcmFuZG9tRnVuYygpO1xuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cblx0XHR2YXIgdGltZW91dEhhbmRsZTtcblx0XHQvLyBTZXQgdGltZW91dFxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xuXHRcdFx0dGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcblx0XHRcdH0sIHRpbWVvdXQpO1xuXHRcdH1cblxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXG5cdFx0Y2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcblx0XHRcdHJlamVjdDogcmVqZWN0LFxuXHRcdFx0cmVzb2x2ZTogcmVzb2x2ZVxuXHRcdH07XG5cblx0XHR0cnkge1xuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0bmFtZSxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdH07XG5cbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnQycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcbiAgICAgICAgfVxuICAgIH0pO1xufVxuXG53aW5kb3cuT2JmdXNjYXRlZENhbGwgPSAoaWQsIGFyZ3MsIHRpbWVvdXQpID0+IHtcblxuICAgIC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxuICAgIGlmICh0aW1lb3V0ID09IG51bGwpIHtcbiAgICAgICAgdGltZW91dCA9IDA7XG4gICAgfVxuXG4gICAgLy8gQ3JlYXRlIGEgcHJvbWlzZVxuICAgIHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cbiAgICAgICAgLy8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcbiAgICAgICAgdmFyIGNhbGxiYWNrSUQ7XG4gICAgICAgIGRvIHtcbiAgICAgICAgICAgIGNhbGxiYWNrSUQgPSBpZCArICctJyArIHJhbmRvbUZ1bmMoKTtcbiAgICAgICAgfSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcblxuICAgICAgICB2YXIgdGltZW91dEhhbmRsZTtcbiAgICAgICAgLy8gU2V0IHRpbWVvdXRcbiAgICAgICAgaWYgKHRpbWVvdXQgPiAwKSB7XG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKCdDYWxsIHRvIG1ldGhvZCAnICsgaWQgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcbiAgICAgICAgICAgIH0sIHRpbWVvdXQpO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gU3RvcmUgY2FsbGJhY2tcbiAgICAgICAgY2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuICAgICAgICAgICAgdGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcbiAgICAgICAgICAgIHJlamVjdDogcmVqZWN0LFxuICAgICAgICAgICAgcmVzb2x2ZTogcmVzb2x2ZVxuICAgICAgICB9O1xuXG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xuXHRcdFx0XHRpZCxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdH07XG5cbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnYycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcbiAgICAgICAgfVxuICAgIH0pO1xufTtcblxuXG4vKipcbiAqIENhbGxlZCBieSB0aGUgYmFja2VuZCB0byByZXR1cm4gZGF0YSB0byBhIHByZXZpb3VzbHkgY2FsbGVkXG4gKiBiaW5kaW5nIGludm9jYXRpb25cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gaW5jb21pbmdNZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsYmFjayhpbmNvbWluZ01lc3NhZ2UpIHtcblx0Ly8gUGFyc2UgdGhlIG1lc3NhZ2Vcblx0bGV0IG1lc3NhZ2U7XG5cdHRyeSB7XG5cdFx0bWVzc2FnZSA9IEpTT04ucGFyc2UoaW5jb21pbmdNZXNzYWdlKTtcblx0fSBjYXRjaCAoZSkge1xuXHRcdGNvbnN0IGVycm9yID0gYEludmFsaWQgSlNPTiBwYXNzZWQgdG8gY2FsbGJhY2s6ICR7ZS5tZXNzYWdlfS4gTWVzc2FnZTogJHtpbmNvbWluZ01lc3NhZ2V9YDtcblx0XHRydW50aW1lLkxvZ0RlYnVnKGVycm9yKTtcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuXHR9XG5cdGxldCBjYWxsYmFja0lEID0gbWVzc2FnZS5jYWxsYmFja2lkO1xuXHRsZXQgY2FsbGJhY2tEYXRhID0gY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xuXHRpZiAoIWNhbGxiYWNrRGF0YSkge1xuXHRcdGNvbnN0IGVycm9yID0gYENhbGxiYWNrICcke2NhbGxiYWNrSUR9JyBub3QgcmVnaXN0ZXJlZCEhIWA7XG5cdFx0Y29uc29sZS5lcnJvcihlcnJvcik7IC8vIGVzbGludC1kaXNhYmxlLWxpbmVcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuXHR9XG5cdGNsZWFyVGltZW91dChjYWxsYmFja0RhdGEudGltZW91dEhhbmRsZSk7XG5cblx0ZGVsZXRlIGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcblxuXHRpZiAobWVzc2FnZS5lcnJvcikge1xuXHRcdGNhbGxiYWNrRGF0YS5yZWplY3QobWVzc2FnZS5lcnJvcik7XG5cdH0gZWxzZSB7XG5cdFx0Y2FsbGJhY2tEYXRhLnJlc29sdmUobWVzc2FnZS5yZXN1bHQpO1xuXHR9XG59XG4iLCAiLypcbiBfICAgICAgIF9fICAgICAgXyBfXyAgICBcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKSBcbnxfXy98X18vXFxfXyxfL18vXy9fX19fLyAgXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5pbXBvcnQge0NhbGx9IGZyb20gJy4vY2FsbHMnO1xuXG4vLyBUaGlzIGlzIHdoZXJlIHdlIGJpbmQgZ28gbWV0aG9kIHdyYXBwZXJzXG53aW5kb3cuZ28gPSB7fTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldEJpbmRpbmdzKGJpbmRpbmdzTWFwKSB7XG5cdHRyeSB7XG5cdFx0YmluZGluZ3NNYXAgPSBKU09OLnBhcnNlKGJpbmRpbmdzTWFwKTtcblx0fSBjYXRjaCAoZSkge1xuXHRcdGNvbnNvbGUuZXJyb3IoZSk7XG5cdH1cblxuXHQvLyBJbml0aWFsaXNlIHRoZSBiaW5kaW5ncyBtYXBcblx0d2luZG93LmdvID0gd2luZG93LmdvIHx8IHt9O1xuXG5cdC8vIEl0ZXJhdGUgcGFja2FnZSBuYW1lc1xuXHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcCkuZm9yRWFjaCgocGFja2FnZU5hbWUpID0+IHtcblxuXHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxuXHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdIHx8IHt9O1xuXG5cdFx0Ly8gSXRlcmF0ZSBzdHJ1Y3QgbmFtZXNcblx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV0pLmZvckVhY2goKHN0cnVjdE5hbWUpID0+IHtcblxuXHRcdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XG5cdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSB8fCB7fTtcblxuXHRcdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdKS5mb3JFYWNoKChtZXRob2ROYW1lKSA9PiB7XG5cblx0XHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXVttZXRob2ROYW1lXSA9IGZ1bmN0aW9uICgpIHtcblxuXHRcdFx0XHRcdC8vIE5vIHRpbWVvdXQgYnkgZGVmYXVsdFxuXHRcdFx0XHRcdGxldCB0aW1lb3V0ID0gMDtcblxuXHRcdFx0XHRcdC8vIEFjdHVhbCBmdW5jdGlvblxuXHRcdFx0XHRcdGZ1bmN0aW9uIGR5bmFtaWMoKSB7XG5cdFx0XHRcdFx0XHRjb25zdCBhcmdzID0gW10uc2xpY2UuY2FsbChhcmd1bWVudHMpO1xuXHRcdFx0XHRcdFx0cmV0dXJuIENhbGwoW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXS5qb2luKCcuJyksIGFyZ3MsIHRpbWVvdXQpO1xuXHRcdFx0XHRcdH1cblxuXHRcdFx0XHRcdC8vIEFsbG93IHNldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxuXHRcdFx0XHRcdGR5bmFtaWMuc2V0VGltZW91dCA9IGZ1bmN0aW9uIChuZXdUaW1lb3V0KSB7XG5cdFx0XHRcdFx0XHR0aW1lb3V0ID0gbmV3VGltZW91dDtcblx0XHRcdFx0XHR9O1xuXG5cdFx0XHRcdFx0Ly8gQWxsb3cgZ2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZHluYW1pYy5nZXRUaW1lb3V0ID0gZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRcdFx0cmV0dXJuIHRpbWVvdXQ7XG5cdFx0XHRcdFx0fTtcblxuXHRcdFx0XHRcdHJldHVybiBkeW5hbWljO1xuXHRcdFx0XHR9KCk7XG5cdFx0XHR9KTtcblx0XHR9KTtcblx0fSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dSZWxvYWQoKSB7XG4gICAgd2luZG93LmxvY2F0aW9uLnJlbG9hZCgpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93UmVsb2FkQXBwKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1InKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFN5c3RlbURlZmF1bHRUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBU0RUJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRMaWdodFRoZW1lKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FMVCcpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0RGFya1RoZW1lKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV0FEVCcpO1xufVxuXG4vKipcbiAqIFBsYWNlIHRoZSB3aW5kb3cgaW4gdGhlIGNlbnRlciBvZiB0aGUgc2NyZWVuXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93Q2VudGVyKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV2MnKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSB3aW5kb3cgdGl0bGVcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gdGl0bGVcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFRpdGxlKHRpdGxlKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVCcgKyB0aXRsZSk7XG59XG5cbi8qKlxuICogTWFrZXMgdGhlIHdpbmRvdyBnbyBmdWxsc2NyZWVuXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93RnVsbHNjcmVlbigpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dGJyk7XG59XG5cbi8qKlxuICogUmV2ZXJ0cyB0aGUgd2luZG93IGZyb20gZnVsbHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VuZnVsbHNjcmVlbigpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dmJyk7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgaW4gZnVsbCBzY3JlZW4gbW9kZSBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzRnVsbHNjcmVlbigpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc0Z1bGxzY3JlZW5cIik7XG59XG5cbi8qKlxuICogU2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldFNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3M6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8e3c6IG51bWJlciwgaDogbnVtYmVyfT59IFRoZSBzaXplIG9mIHRoZSB3aW5kb3dcblxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0U2l6ZSgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dHZXRTaXplXCIpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldE1heFNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1o6JyArIHdpZHRoICsgJzonICsgaGVpZ2h0KTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d6OicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cblxuXG4vKipcbiAqIFNldCB0aGUgd2luZG93IEFsd2F5c09uVG9wIG9yIG5vdCBvbiB0b3BcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRBbHdheXNPblRvcChiKSB7XG5cbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBVFA6JyArIChiID8gJzEnIDogJzAnKSk7XG59XG5cblxuXG5cbi8qKlxuICogU2V0IHRoZSBQb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHhcbiAqIEBwYXJhbSB7bnVtYmVyfSB5XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRQb3NpdGlvbih4LCB5KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXcDonICsgeCArICc6JyArIHkpO1xufVxuXG4vKipcbiAqIEdldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8e3g6IG51bWJlciwgeTogbnVtYmVyfT59IFRoZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dHZXRQb3NpdGlvbigpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dHZXRQb3NcIik7XG59XG5cbi8qKlxuICogSGlkZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SGlkZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dIJyk7XG59XG5cbi8qKlxuICogU2hvdyB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2hvdygpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dTJyk7XG59XG5cbi8qKlxuICogTWF4aW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd01heGltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV00nKTtcbn1cblxuLyoqXG4gKiBUb2dnbGUgdGhlIE1heGltaXNlIG9mIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dUb2dnbGVNYXhpbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d0Jyk7XG59XG5cbi8qKlxuICogVW5tYXhpbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5tYXhpbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dVJyk7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNNYXhpbWlzZWQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNYXhpbWlzZWRcIik7XG59XG5cbi8qKlxuICogTWluaW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd01pbmltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV20nKTtcbn1cblxuLyoqXG4gKiBVbm1pbmltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbm1pbmltaXNlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3UnKTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBtaW5pbWlzZWQgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc01pbmltaXNlZCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc01pbmltaXNlZFwiKTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIHRoZSBzdGF0ZSBvZiB0aGUgd2luZG93LCBpLmUuIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyBub3JtYWwgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc05vcm1hbCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc05vcm1hbFwiKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IFIgUmVkXG4gKiBAcGFyYW0ge251bWJlcn0gRyBHcmVlblxuICogQHBhcmFtIHtudW1iZXJ9IEIgQmx1ZVxuICogQHBhcmFtIHtudW1iZXJ9IEEgQWxwaGFcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldEJhY2tncm91bmRDb2xvdXIoUiwgRywgQiwgQSkge1xuICAgIGxldCByZ2JhID0gSlNPTi5zdHJpbmdpZnkoe3I6IFIgfHwgMCwgZzogRyB8fCAwLCBiOiBCIHx8IDAsIGE6IEEgfHwgMjU1fSk7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXcjonICsgcmdiYSk7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuXG4vKipcbiAqIEdldHMgdGhlIGFsbCBzY3JlZW5zLiBDYWxsIHRoaXMgYW5ldyBlYWNoIHRpbWUgeW91IHdhbnQgdG8gcmVmcmVzaCBkYXRhIGZyb20gdGhlIHVuZGVybHlpbmcgd2luZG93aW5nIHN5c3RlbS5cbiAqIEBleHBvcnRcbiAqIEB0eXBlZGVmIHtpbXBvcnQoJy4uL3dyYXBwZXIvcnVudGltZScpLlNjcmVlbn0gU2NyZWVuXG4gKiBAcmV0dXJuIHtQcm9taXNlPHtTY3JlZW5bXX0+fSBUaGUgc2NyZWVuc1xuICovXG5leHBvcnQgZnVuY3Rpb24gU2NyZWVuR2V0QWxsKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlNjcmVlbkdldEFsbFwiKTtcbn1cbiIsICIvKipcbiAqIEBkZXNjcmlwdGlvbjogVXNlIHRoZSBzeXN0ZW0gZGVmYXVsdCBicm93c2VyIHRvIG9wZW4gdGhlIHVybFxuICogQHBhcmFtIHtzdHJpbmd9IHVybCBcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCcm93c2VyT3BlblVSTCh1cmwpIHtcbiAgd2luZG93LldhaWxzSW52b2tlKCdCTzonICsgdXJsKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG4vKipcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IHRleHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENsaXBib2FyZFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsaXBib2FyZFNldFRleHRcIiwgW3RleHRdKTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIHRleHQgY29udGVudCBvZiB0aGUgY2xpcGJvYXJkXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx7c3RyaW5nfT59IFRleHQgY29udGVudCBvZiB0aGUgY2xpcGJvYXJkXG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENsaXBib2FyZEdldFRleHQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6Q2xpcGJvYXJkR2V0VGV4dFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtFdmVudHNPbiwgRXZlbnRzT2ZmfSBmcm9tIFwiLi9ldmVudHNcIjtcblxuY29uc3QgZmxhZ3MgPSB7XG4gICAgcmVnaXN0ZXJlZDogZmFsc2UsXG4gICAgZGVmYXVsdFVzZURyb3BUYXJnZXQ6IHRydWUsXG4gICAgdXNlRHJvcFRhcmdldDogdHJ1ZSxcbiAgICBuZXh0RGVhY3RpdmF0ZTogbnVsbCxcbiAgICBuZXh0RGVhY3RpdmF0ZVRpbWVvdXQ6IG51bGwsXG59O1xuXG5jb25zdCBEUk9QX1RBUkdFVF9BQ1RJVkUgPSBcIndhaWxzLWRyb3AtdGFyZ2V0LWFjdGl2ZVwiO1xuXG4vKipcbiAqIGNoZWNrU3R5bGVEcm9wVGFyZ2V0IGNoZWNrcyBpZiB0aGUgc3R5bGUgaGFzIHRoZSBkcm9wIHRhcmdldCBhdHRyaWJ1dGVcbiAqIFxuICogQHBhcmFtIHtDU1NTdHlsZURlY2xhcmF0aW9ufSBzdHlsZSBcbiAqIEByZXR1cm5zIFxuICovXG5mdW5jdGlvbiBjaGVja1N0eWxlRHJvcFRhcmdldChzdHlsZSkge1xuICAgIGNvbnN0IGNzc0Ryb3BWYWx1ZSA9IHN0eWxlLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSkudHJpbSgpO1xuICAgIGlmIChjc3NEcm9wVmFsdWUpIHtcbiAgICAgICAgaWYgKGNzc0Ryb3BWYWx1ZSA9PT0gd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BWYWx1ZSkge1xuICAgICAgICAgICAgcmV0dXJuIHRydWU7XG4gICAgICAgIH1cbiAgICAgICAgLy8gaWYgdGhlIGVsZW1lbnQgaGFzIHRoZSBkcm9wIHRhcmdldCBhdHRyaWJ1dGUsIGJ1dCBcbiAgICAgICAgLy8gdGhlIHZhbHVlIGlzIG5vdCBjb3JyZWN0LCB0ZXJtaW5hdGUgZmluZGluZyBwcm9jZXNzLlxuICAgICAgICAvLyBUaGlzIGNhbiBiZSB1c2VmdWwgdG8gYmxvY2sgc29tZSBjaGlsZCBlbGVtZW50cyBmcm9tIGJlaW5nIGRyb3AgdGFyZ2V0cy5cbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbiAgICByZXR1cm4gZmFsc2U7XG59XG5cbi8qKlxuICogb25EcmFnT3ZlciBpcyBjYWxsZWQgd2hlbiB0aGUgZHJhZ292ZXIgZXZlbnQgaXMgZW1pdHRlZC5cbiAqIEBwYXJhbSB7RHJhZ0V2ZW50fSBlXG4gKiBAcmV0dXJuc1xuICovXG5mdW5jdGlvbiBvbkRyYWdPdmVyKGUpIHtcbiAgICAvLyBDaGVjayBpZiB0aGlzIGlzIGFuIGV4dGVybmFsIGZpbGUgZHJvcCBvciBpbnRlcm5hbCBIVE1MIGRyYWdcbiAgICAvLyBFeHRlcm5hbCBmaWxlIGRyb3BzIHdpbGwgaGF2ZSBcIkZpbGVzXCIgaW4gdGhlIHR5cGVzIGFycmF5XG4gICAgLy8gSW50ZXJuYWwgSFRNTCBkcmFncyB0eXBpY2FsbHkgaGF2ZSBcInRleHQvcGxhaW5cIiwgXCJ0ZXh0L2h0bWxcIiBvciBjdXN0b20gdHlwZXNcbiAgICBjb25zdCBpc0ZpbGVEcm9wID0gZS5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoXCJGaWxlc1wiKTtcblxuICAgIC8vIE9ubHkgaGFuZGxlIGV4dGVybmFsIGZpbGUgZHJvcHMsIGxldCBpbnRlcm5hbCBIVE1MNSBkcmFnLWFuZC1kcm9wIHdvcmsgbm9ybWFsbHlcbiAgICBpZiAoIWlzRmlsZURyb3ApIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIEFMV0FZUyBwcmV2ZW50IGRlZmF1bHQgZm9yIGZpbGUgZHJvcHMgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgZS5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdjb3B5JztcblxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3ApIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghZmxhZ3MudXNlRHJvcFRhcmdldCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgY29uc3QgZWxlbWVudCA9IGUudGFyZ2V0O1xuXG4gICAgLy8gVHJpZ2dlciBkZWJvdW5jZSBmdW5jdGlvbiB0byBkZWFjdGl2YXRlIGRyb3AgdGFyZ2V0c1xuICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xuXG4gICAgLy8gaWYgdGhlIGVsZW1lbnQgaXMgbnVsbCBvciBlbGVtZW50IGlzIG5vdCBjaGlsZCBvZiBkcm9wIHRhcmdldCBlbGVtZW50XG4gICAgaWYgKCFlbGVtZW50IHx8ICFjaGVja1N0eWxlRHJvcFRhcmdldChnZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgbGV0IGN1cnJlbnRFbGVtZW50ID0gZWxlbWVudDtcbiAgICB3aGlsZSAoY3VycmVudEVsZW1lbnQpIHtcbiAgICAgICAgLy8gY2hlY2sgaWYgY3VycmVudEVsZW1lbnQgaXMgZHJvcCB0YXJnZXQgZWxlbWVudFxuICAgICAgICBpZiAoY2hlY2tTdHlsZURyb3BUYXJnZXQoZ2V0Q29tcHV0ZWRTdHlsZShjdXJyZW50RWxlbWVudCkpKSB7XG4gICAgICAgICAgICBjdXJyZW50RWxlbWVudC5jbGFzc0xpc3QuYWRkKERST1BfVEFSR0VUX0FDVElWRSk7XG4gICAgICAgIH1cbiAgICAgICAgY3VycmVudEVsZW1lbnQgPSBjdXJyZW50RWxlbWVudC5wYXJlbnRFbGVtZW50O1xuICAgIH1cbn1cblxuLyoqXG4gKiBvbkRyYWdMZWF2ZSBpcyBjYWxsZWQgd2hlbiB0aGUgZHJhZ2xlYXZlIGV2ZW50IGlzIGVtaXR0ZWQuXG4gKiBAcGFyYW0ge0RyYWdFdmVudH0gZVxuICogQHJldHVybnNcbiAqL1xuZnVuY3Rpb24gb25EcmFnTGVhdmUoZSkge1xuICAgIC8vIENoZWNrIGlmIHRoaXMgaXMgYW4gZXh0ZXJuYWwgZmlsZSBkcm9wIG9yIGludGVybmFsIEhUTUwgZHJhZ1xuICAgIGNvbnN0IGlzRmlsZURyb3AgPSBlLmRhdGFUcmFuc2Zlci50eXBlcy5pbmNsdWRlcyhcIkZpbGVzXCIpO1xuXG4gICAgLy8gT25seSBoYW5kbGUgZXh0ZXJuYWwgZmlsZSBkcm9wcywgbGV0IGludGVybmFsIEhUTUw1IGRyYWctYW5kLWRyb3Agd29yayBub3JtYWxseVxuICAgIGlmICghaXNGaWxlRHJvcCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQUxXQVlTIHByZXZlbnQgZGVmYXVsdCBmb3IgZmlsZSBkcm9wcyB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgIGUucHJldmVudERlZmF1bHQoKTtcblxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVdhaWxzRHJhZ0FuZERyb3ApIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghZmxhZ3MudXNlRHJvcFRhcmdldCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gRmluZCB0aGUgY2xvc2UgZHJvcCB0YXJnZXQgZWxlbWVudFxuICAgIGlmICghZS50YXJnZXQgfHwgIWNoZWNrU3R5bGVEcm9wVGFyZ2V0KGdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpKSkge1xuICAgICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG5cbiAgICAvLyBUcmlnZ2VyIGRlYm91bmNlIGZ1bmN0aW9uIHRvIGRlYWN0aXZhdGUgZHJvcCB0YXJnZXRzXG4gICAgaWYoZmxhZ3MubmV4dERlYWN0aXZhdGUpIGZsYWdzLm5leHREZWFjdGl2YXRlKCk7XG4gICAgXG4gICAgLy8gVXNlIGRlYm91bmNlIHRlY2huaXF1ZSB0byB0YWNsZSBkcmFnbGVhdmUgZXZlbnRzIG9uIG92ZXJsYXBwaW5nIGVsZW1lbnRzIGFuZCBkcm9wIHRhcmdldCBlbGVtZW50c1xuICAgIGZsYWdzLm5leHREZWFjdGl2YXRlID0gKCkgPT4ge1xuICAgICAgICAvLyBEZWFjdGl2YXRlIGFsbCBkcm9wIHRhcmdldHMsIG5ldyBkcm9wIHRhcmdldCB3aWxsIGJlIGFjdGl2YXRlZCBvbiBuZXh0IGRyYWdvdmVyIGV2ZW50XG4gICAgICAgIEFycmF5LmZyb20oZG9jdW1lbnQuZ2V0RWxlbWVudHNCeUNsYXNzTmFtZShEUk9QX1RBUkdFVF9BQ1RJVkUpKS5mb3JFYWNoKGVsID0+IGVsLmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFKSk7XG4gICAgICAgIC8vIFJlc2V0IG5leHREZWFjdGl2YXRlXG4gICAgICAgIGZsYWdzLm5leHREZWFjdGl2YXRlID0gbnVsbDtcbiAgICAgICAgLy8gQ2xlYXIgdGltZW91dFxuICAgICAgICBpZiAoZmxhZ3MubmV4dERlYWN0aXZhdGVUaW1lb3V0KSB7XG4gICAgICAgICAgICBjbGVhclRpbWVvdXQoZmxhZ3MubmV4dERlYWN0aXZhdGVUaW1lb3V0KTtcbiAgICAgICAgICAgIGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCA9IG51bGw7XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBTZXQgdGltZW91dCB0byBkZWFjdGl2YXRlIGRyb3AgdGFyZ2V0cyBpZiBub3QgdHJpZ2dlcmVkIGJ5IG5leHQgZHJhZyBldmVudFxuICAgIGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCA9IHNldFRpbWVvdXQoKCkgPT4ge1xuICAgICAgICBpZihmbGFncy5uZXh0RGVhY3RpdmF0ZSkgZmxhZ3MubmV4dERlYWN0aXZhdGUoKTtcbiAgICB9LCA1MCk7XG59XG5cbi8qKlxuICogb25Ecm9wIGlzIGNhbGxlZCB3aGVuIHRoZSBkcm9wIGV2ZW50IGlzIGVtaXR0ZWQuXG4gKiBAcGFyYW0ge0RyYWdFdmVudH0gZVxuICogQHJldHVybnNcbiAqL1xuZnVuY3Rpb24gb25Ecm9wKGUpIHtcbiAgICAvLyBDaGVjayBpZiB0aGlzIGlzIGFuIGV4dGVybmFsIGZpbGUgZHJvcCBvciBpbnRlcm5hbCBIVE1MIGRyYWdcbiAgICBjb25zdCBpc0ZpbGVEcm9wID0gZS5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoXCJGaWxlc1wiKTtcblxuICAgIC8vIE9ubHkgaGFuZGxlIGV4dGVybmFsIGZpbGUgZHJvcHMsIGxldCBpbnRlcm5hbCBIVE1MNSBkcmFnLWFuZC1kcm9wIHdvcmsgbm9ybWFsbHlcbiAgICBpZiAoIWlzRmlsZURyb3ApIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIEFMV0FZUyBwcmV2ZW50IGRlZmF1bHQgZm9yIGZpbGUgZHJvcHMgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICBlLnByZXZlbnREZWZhdWx0KCk7XG5cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVXYWlsc0RyYWdBbmREcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoQ2FuUmVzb2x2ZUZpbGVQYXRocygpKSB7XG4gICAgICAgIC8vIHByb2Nlc3MgZmlsZXNcbiAgICAgICAgbGV0IGZpbGVzID0gW107XG4gICAgICAgIGlmIChlLmRhdGFUcmFuc2Zlci5pdGVtcykge1xuICAgICAgICAgICAgZmlsZXMgPSBbLi4uZS5kYXRhVHJhbnNmZXIuaXRlbXNdLm1hcCgoaXRlbSwgaSkgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChpdGVtLmtpbmQgPT09ICdmaWxlJykge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm4gaXRlbS5nZXRBc0ZpbGUoKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9KTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGZpbGVzID0gWy4uLmUuZGF0YVRyYW5zZmVyLmZpbGVzXTtcbiAgICAgICAgfVxuICAgICAgICB3aW5kb3cucnVudGltZS5SZXNvbHZlRmlsZVBhdGhzKGUueCwgZS55LCBmaWxlcyk7XG4gICAgfVxuXG4gICAgaWYgKCFmbGFncy51c2VEcm9wVGFyZ2V0KSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBUcmlnZ2VyIGRlYm91bmNlIGZ1bmN0aW9uIHRvIGRlYWN0aXZhdGUgZHJvcCB0YXJnZXRzXG4gICAgaWYoZmxhZ3MubmV4dERlYWN0aXZhdGUpIGZsYWdzLm5leHREZWFjdGl2YXRlKCk7XG5cbiAgICAvLyBEZWFjdGl2YXRlIGFsbCBkcm9wIHRhcmdldHNcbiAgICBBcnJheS5mcm9tKGRvY3VtZW50LmdldEVsZW1lbnRzQnlDbGFzc05hbWUoRFJPUF9UQVJHRVRfQUNUSVZFKSkuZm9yRWFjaChlbCA9PiBlbC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRSkpO1xufVxuXG4vKipcbiAqIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIGNoZWNrcyB0aGUgYnJvd3NlcidzIGNhcGFiaWxpdHkgb2Ygc2VuZGluZyBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0c1xuICpcbiAqIEByZXR1cm5zIHtib29sZWFufVxuICogQGNvbnN0cnVjdG9yXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYW5SZXNvbHZlRmlsZVBhdGhzKCkge1xuICAgIHJldHVybiB3aW5kb3cuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyAhPSBudWxsO1xufVxuXG4vKipcbiAqIFJlc29sdmVGaWxlUGF0aHMgc2VuZHMgZHJvcCBldmVudHMgdG8gdGhlIEdPIHNpZGUgdG8gcmVzb2x2ZSBmaWxlIHBhdGhzIG9uIHdpbmRvd3MuXG4gKlxuICogQHBhcmFtIHtudW1iZXJ9IHhcbiAqIEBwYXJhbSB7bnVtYmVyfSB5XG4gKiBAcGFyYW0ge2FueVtdfSBmaWxlc1xuICogQGNvbnN0cnVjdG9yXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZXNvbHZlRmlsZVBhdGhzKHgsIHksIGZpbGVzKSB7XG4gICAgLy8gT25seSBmb3Igd2luZG93cyB3ZWJ2aWV3MiA+PSAxLjAuMTc3NC4zMFxuICAgIC8vIGh0dHBzOi8vbGVhcm4ubWljcm9zb2Z0LmNvbS9lbi11cy9taWNyb3NvZnQtZWRnZS93ZWJ2aWV3Mi9yZWZlcmVuY2Uvd2luMzIvaWNvcmV3ZWJ2aWV3MndlYm1lc3NhZ2VyZWNlaXZlZGV2ZW50YXJnczI/dmlldz13ZWJ2aWV3Mi0xLjAuMTgyMy4zMiNhcHBsaWVzLXRvXG4gICAgaWYgKHdpbmRvdy5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKSB7XG4gICAgICAgIGNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKGBmaWxlOmRyb3A6JHt4fToke3l9YCwgZmlsZXMpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBDYWxsYmFjayBmb3IgT25GaWxlRHJvcCByZXR1cm5zIGEgc2xpY2Ugb2YgZmlsZSBwYXRoIHN0cmluZ3Mgd2hlbiBhIGRyb3AgaXMgZmluaXNoZWQuXG4gKlxuICogQGV4cG9ydFxuICogQGNhbGxiYWNrIE9uRmlsZURyb3BDYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IHggLSB4IGNvb3JkaW5hdGUgb2YgdGhlIGRyb3BcbiAqIEBwYXJhbSB7bnVtYmVyfSB5IC0geSBjb29yZGluYXRlIG9mIHRoZSBkcm9wXG4gKiBAcGFyYW0ge3N0cmluZ1tdfSBwYXRocyAtIEEgbGlzdCBvZiBmaWxlIHBhdGhzLlxuICovXG5cbi8qKlxuICogT25GaWxlRHJvcCBsaXN0ZW5zIHRvIGRyYWcgYW5kIGRyb3AgZXZlbnRzIGFuZCBjYWxscyB0aGUgY2FsbGJhY2sgd2l0aCB0aGUgY29vcmRpbmF0ZXMgb2YgdGhlIGRyb3AgYW5kIGFuIGFycmF5IG9mIHBhdGggc3RyaW5ncy5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge09uRmlsZURyb3BDYWxsYmFja30gY2FsbGJhY2sgLSBDYWxsYmFjayBmb3IgT25GaWxlRHJvcCByZXR1cm5zIGEgc2xpY2Ugb2YgZmlsZSBwYXRoIHN0cmluZ3Mgd2hlbiBhIGRyb3AgaXMgZmluaXNoZWQuXG4gKiBAcGFyYW0ge2Jvb2xlYW59IFt1c2VEcm9wVGFyZ2V0PXRydWVdIC0gT25seSBjYWxsIHRoZSBjYWxsYmFjayB3aGVuIHRoZSBkcm9wIGZpbmlzaGVkIG9uIGFuIGVsZW1lbnQgdGhhdCBoYXMgdGhlIGRyb3AgdGFyZ2V0IHN0eWxlLiAoLS13YWlscy1kcm9wLXRhcmdldClcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uRmlsZURyb3AoY2FsbGJhY2ssIHVzZURyb3BUYXJnZXQpIHtcbiAgICBpZiAodHlwZW9mIGNhbGxiYWNrICE9PSBcImZ1bmN0aW9uXCIpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihcIkRyYWdBbmREcm9wQ2FsbGJhY2sgaXMgbm90IGEgZnVuY3Rpb25cIik7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoZmxhZ3MucmVnaXN0ZXJlZCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGZsYWdzLnJlZ2lzdGVyZWQgPSB0cnVlO1xuXG4gICAgY29uc3QgdURUUFQgPSB0eXBlb2YgdXNlRHJvcFRhcmdldDtcbiAgICBmbGFncy51c2VEcm9wVGFyZ2V0ID0gdURUUFQgPT09IFwidW5kZWZpbmVkXCIgfHwgdURUUFQgIT09IFwiYm9vbGVhblwiID8gZmxhZ3MuZGVmYXVsdFVzZURyb3BUYXJnZXQgOiB1c2VEcm9wVGFyZ2V0O1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcmFnb3ZlcicsIG9uRHJhZ092ZXIpO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCBvbkRyYWdMZWF2ZSk7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCBvbkRyb3ApO1xuXG4gICAgbGV0IGNiID0gY2FsbGJhY2s7XG4gICAgaWYgKGZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgY2IgPSBmdW5jdGlvbiAoeCwgeSwgcGF0aHMpIHtcbiAgICAgICAgICAgIGNvbnN0IGVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHgsIHkpXG4gICAgICAgICAgICAvLyBpZiB0aGUgZWxlbWVudCBpcyBudWxsIG9yIGVsZW1lbnQgaXMgbm90IGNoaWxkIG9mIGRyb3AgdGFyZ2V0IGVsZW1lbnQsIHJldHVybiBudWxsXG4gICAgICAgICAgICBpZiAoIWVsZW1lbnQgfHwgIWNoZWNrU3R5bGVEcm9wVGFyZ2V0KGdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkpKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG51bGw7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBjYWxsYmFjayh4LCB5LCBwYXRocyk7XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICBFdmVudHNPbihcIndhaWxzOmZpbGUtZHJvcFwiLCBjYik7XG59XG5cbi8qKlxuICogT25GaWxlRHJvcE9mZiByZW1vdmVzIHRoZSBkcmFnIGFuZCBkcm9wIGxpc3RlbmVycyBhbmQgaGFuZGxlcnMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbkZpbGVEcm9wT2ZmKCkge1xuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcmFnb3ZlcicsIG9uRHJhZ092ZXIpO1xuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCBvbkRyYWdMZWF2ZSk7XG4gICAgd2luZG93LnJlbW92ZUV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCBvbkRyb3ApO1xuICAgIEV2ZW50c09mZihcIndhaWxzOmZpbGUtZHJvcFwiKTtcbiAgICBmbGFncy5yZWdpc3RlcmVkID0gZmFsc2U7XG59XG4iLCAiLypcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG5cblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXG4qL1xuZXhwb3J0IGZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgY29uc3QgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcbiAgICBjb25zdCBjb21wdXRlZFN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XG4gICAgY29uc3QgZGVmYXVsdENvbnRleHRNZW51QWN0aW9uID0gY29tcHV0ZWRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKTtcbiAgICBzd2l0Y2ggKGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbikge1xuICAgICAgICBjYXNlIFwic2hvd1wiOlxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBjYXNlIFwiaGlkZVwiOlxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgICAgICAgICBpZiAoZWxlbWVudC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBhbmQgYWN0aW9uIGlzIG9uIHRoZSBzZWxlY3RlZCBlbGVtZW50c1xuICAgICAgICAgICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgICAgICAgICAgY29uc3QgaGFzU2VsZWN0aW9uID0gKHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDApXG4gICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgICAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXG4gICAgICAgICAgICBpZiAoZWxlbWVudC50YWdOYW1lID09PSBcIklOUFVUXCIgfHwgZWxlbWVudC50YWdOYW1lID09PSBcIlRFWFRBUkVBXCIpIHtcbiAgICAgICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghZWxlbWVudC5yZWFkT25seSAmJiAhZWxlbWVudC5kaXNhYmxlZCkpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG59XG4iLCAiLypcbiBfICAgICAgIF9fICAgICAgXyBfX1xufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuLyoqXG4gKiBJbml0aWFsaXplIHRoZSBub3RpZmljYXRpb24gc2VydmljZSBmb3IgdGhlIGFwcGxpY2F0aW9uLlxuICogVGhpcyBtdXN0IGJlIGNhbGxlZCBiZWZvcmUgc2VuZGluZyBhbnkgbm90aWZpY2F0aW9ucy5cbiAqIE9uIG1hY09TLCB0aGlzIGFsc28gZW5zdXJlcyB0aGUgbm90aWZpY2F0aW9uIGRlbGVnYXRlIGlzIHByb3Blcmx5IGluaXRpYWxpemVkLlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbml0aWFsaXplTm90aWZpY2F0aW9ucygpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpJbml0aWFsaXplTm90aWZpY2F0aW9uc1wiKTtcbn1cblxuLyoqXG4gKiBDbGVhbiB1cCBub3RpZmljYXRpb24gcmVzb3VyY2VzIGFuZCByZWxlYXNlIGFueSBoZWxkIGNvbm5lY3Rpb25zLlxuICogVGhpcyBzaG91bGQgYmUgY2FsbGVkIHdoZW4gc2h1dHRpbmcgZG93biB0aGUgYXBwbGljYXRpb24gdG8gcHJvcGVybHkgcmVsZWFzZSByZXNvdXJjZXNcbiAqIChwcmltYXJpbHkgbmVlZGVkIG9uIExpbnV4IHRvIGNsb3NlIEQtQnVzIGNvbm5lY3Rpb25zKS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2xlYW51cE5vdGlmaWNhdGlvbnMoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6Q2xlYW51cE5vdGlmaWNhdGlvbnNcIik7XG59XG5cbi8qKlxuICogQ2hlY2sgaWYgbm90aWZpY2F0aW9ucyBhcmUgYXZhaWxhYmxlIG9uIHRoZSBjdXJyZW50IHBsYXRmb3JtLlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRydWUgaWYgbm90aWZpY2F0aW9ucyBhcmUgYXZhaWxhYmxlLCBmYWxzZSBvdGhlcndpc2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTm90aWZpY2F0aW9uQXZhaWxhYmxlKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOklzTm90aWZpY2F0aW9uQXZhaWxhYmxlXCIpO1xufVxuXG4vKipcbiAqIFJlcXVlc3Qgbm90aWZpY2F0aW9uIGF1dGhvcml6YXRpb24gZnJvbSB0aGUgdXNlci5cbiAqIE9uIG1hY09TLCB0aGlzIHByb21wdHMgdGhlIHVzZXIgdG8gYWxsb3cgbm90aWZpY2F0aW9ucy5cbiAqIE9uIG90aGVyIHBsYXRmb3JtcywgdGhpcyBhbHdheXMgcmV0dXJucyB0cnVlLlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRydWUgaWYgYXV0aG9yaXphdGlvbiB3YXMgZ3JhbnRlZCwgZmFsc2Ugb3RoZXJ3aXNlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZXF1ZXN0Tm90aWZpY2F0aW9uQXV0aG9yaXphdGlvbigpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZXF1ZXN0Tm90aWZpY2F0aW9uQXV0aG9yaXphdGlvblwiKTtcbn1cblxuLyoqXG4gKiBDaGVjayB0aGUgY3VycmVudCBub3RpZmljYXRpb24gYXV0aG9yaXphdGlvbiBzdGF0dXMuXG4gKiBPbiBtYWNPUywgdGhpcyBjaGVja3MgaWYgdGhlIGFwcCBoYXMgbm90aWZpY2F0aW9uIHBlcm1pc3Npb25zLlxuICogT24gb3RoZXIgcGxhdGZvcm1zLCB0aGlzIGFsd2F5cyByZXR1cm5zIHRydWUuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVHJ1ZSBpZiBhdXRob3JpemVkLCBmYWxzZSBvdGhlcndpc2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENoZWNrTm90aWZpY2F0aW9uQXV0aG9yaXphdGlvbigpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDaGVja05vdGlmaWNhdGlvbkF1dGhvcml6YXRpb25cIik7XG59XG5cbi8qKlxuICogU2VuZCBhIGJhc2ljIG5vdGlmaWNhdGlvbiB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogVGhlIG5vdGlmaWNhdGlvbiB3aWxsIGRpc3BsYXkgd2l0aCB0aGUgcHJvdmlkZWQgdGl0bGUsIHN1YnRpdGxlIChpZiBzdXBwb3J0ZWQpLCBhbmQgYm9keSB0ZXh0LlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7T2JqZWN0fSBvcHRpb25zIC0gTm90aWZpY2F0aW9uIG9wdGlvbnNcbiAqIEBwYXJhbSB7c3RyaW5nfSBvcHRpb25zLmlkIC0gVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBub3RpZmljYXRpb25cbiAqIEBwYXJhbSB7c3RyaW5nfSBvcHRpb25zLnRpdGxlIC0gTm90aWZpY2F0aW9uIHRpdGxlXG4gKiBAcGFyYW0ge3N0cmluZ30gW29wdGlvbnMuc3VidGl0bGVdIC0gTm90aWZpY2F0aW9uIHN1YnRpdGxlIChtYWNPUyBhbmQgTGludXggb25seSlcbiAqIEBwYXJhbSB7c3RyaW5nfSBbb3B0aW9ucy5ib2R5XSAtIE5vdGlmaWNhdGlvbiBib2R5IHRleHRcbiAqIEBwYXJhbSB7c3RyaW5nfSBbb3B0aW9ucy5jYXRlZ29yeUlkXSAtIENhdGVnb3J5IElEIGZvciBhY3Rpb24gYnV0dG9ucyAocmVxdWlyZXMgU2VuZE5vdGlmaWNhdGlvbldpdGhBY3Rpb25zKVxuICogQHBhcmFtIHtPYmplY3Q8c3RyaW5nLCBhbnk+fSBbb3B0aW9ucy5kYXRhXSAtIEFkZGl0aW9uYWwgdXNlciBkYXRhIHRvIGF0dGFjaCB0byB0aGUgbm90aWZpY2F0aW9uXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2VuZE5vdGlmaWNhdGlvbihvcHRpb25zKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6U2VuZE5vdGlmaWNhdGlvblwiLCBbb3B0aW9uc10pO1xufVxuXG4vKipcbiAqIFNlbmQgYSBub3RpZmljYXRpb24gd2l0aCBhY3Rpb24gYnV0dG9ucy5cbiAqIEEgTm90aWZpY2F0aW9uQ2F0ZWdvcnkgbXVzdCBiZSByZWdpc3RlcmVkIGZpcnN0IHVzaW5nIFJlZ2lzdGVyTm90aWZpY2F0aW9uQ2F0ZWdvcnkuXG4gKiBUaGUgb3B0aW9ucy5jYXRlZ29yeUlkIG11c3QgbWF0Y2ggYSBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgY2F0ZWdvcnkgSUQuXG4gKiBJZiB0aGUgY2F0ZWdvcnkgaXMgbm90IGZvdW5kLCBhIGJhc2ljIG5vdGlmaWNhdGlvbiB3aWxsIGJlIHNlbnQgaW5zdGVhZC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge09iamVjdH0gb3B0aW9ucyAtIE5vdGlmaWNhdGlvbiBvcHRpb25zXG4gKiBAcGFyYW0ge3N0cmluZ30gb3B0aW9ucy5pZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgbm90aWZpY2F0aW9uXG4gKiBAcGFyYW0ge3N0cmluZ30gb3B0aW9ucy50aXRsZSAtIE5vdGlmaWNhdGlvbiB0aXRsZVxuICogQHBhcmFtIHtzdHJpbmd9IFtvcHRpb25zLnN1YnRpdGxlXSAtIE5vdGlmaWNhdGlvbiBzdWJ0aXRsZSAobWFjT1MgYW5kIExpbnV4IG9ubHkpXG4gKiBAcGFyYW0ge3N0cmluZ30gW29wdGlvbnMuYm9keV0gLSBOb3RpZmljYXRpb24gYm9keSB0ZXh0XG4gKiBAcGFyYW0ge3N0cmluZ30gb3B0aW9ucy5jYXRlZ29yeUlkIC0gQ2F0ZWdvcnkgSUQgdGhhdCBtYXRjaGVzIGEgcmVnaXN0ZXJlZCBjYXRlZ29yeVxuICogQHBhcmFtIHtPYmplY3Q8c3RyaW5nLCBhbnk+fSBbb3B0aW9ucy5kYXRhXSAtIEFkZGl0aW9uYWwgdXNlciBkYXRhIHRvIGF0dGFjaCB0byB0aGUgbm90aWZpY2F0aW9uXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2VuZE5vdGlmaWNhdGlvbldpdGhBY3Rpb25zKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpTZW5kTm90aWZpY2F0aW9uV2l0aEFjdGlvbnNcIiwgW29wdGlvbnNdKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlciBhIG5vdGlmaWNhdGlvbiBjYXRlZ29yeSB0aGF0IGNhbiBiZSB1c2VkIHdpdGggU2VuZE5vdGlmaWNhdGlvbldpdGhBY3Rpb25zLlxuICogQ2F0ZWdvcmllcyBkZWZpbmUgdGhlIGFjdGlvbiBidXR0b25zIGFuZCBvcHRpb25hbCByZXBseSBmaWVsZHMgdGhhdCB3aWxsIGFwcGVhciBvbiBub3RpZmljYXRpb25zLlxuICogUmVnaXN0ZXJpbmcgYSBjYXRlZ29yeSB3aXRoIHRoZSBzYW1lIElEIGFzIGEgcHJldmlvdXNseSByZWdpc3RlcmVkIGNhdGVnb3J5IHdpbGwgb3ZlcnJpZGUgaXQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtPYmplY3R9IGNhdGVnb3J5IC0gTm90aWZpY2F0aW9uIGNhdGVnb3J5IGRlZmluaXRpb25cbiAqIEBwYXJhbSB7c3RyaW5nfSBjYXRlZ29yeS5pZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgY2F0ZWdvcnlcbiAqIEBwYXJhbSB7QXJyYXk8T2JqZWN0Pn0gW2NhdGVnb3J5LmFjdGlvbnNdIC0gQXJyYXkgb2YgYWN0aW9uIGJ1dHRvbnNcbiAqIEBwYXJhbSB7c3RyaW5nfSBjYXRlZ29yeS5hY3Rpb25zW10uaWQgLSBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIGFjdGlvblxuICogQHBhcmFtIHtzdHJpbmd9IGNhdGVnb3J5LmFjdGlvbnNbXS50aXRsZSAtIERpc3BsYXkgdGl0bGUgZm9yIHRoZSBhY3Rpb24gYnV0dG9uXG4gKiBAcGFyYW0ge2Jvb2xlYW59IFtjYXRlZ29yeS5hY3Rpb25zW10uZGVzdHJ1Y3RpdmVdIC0gV2hldGhlciB0aGUgYWN0aW9uIGlzIGRlc3RydWN0aXZlIChtYWNPUy1zcGVjaWZpYylcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gW2NhdGVnb3J5Lmhhc1JlcGx5RmllbGRdIC0gV2hldGhlciB0byBpbmNsdWRlIGEgdGV4dCBpbnB1dCBmaWVsZCBmb3IgcmVwbGllc1xuICogQHBhcmFtIHtzdHJpbmd9IFtjYXRlZ29yeS5yZXBseVBsYWNlaG9sZGVyXSAtIFBsYWNlaG9sZGVyIHRleHQgZm9yIHRoZSByZXBseSBmaWVsZCAocmVxdWlyZWQgaWYgaGFzUmVwbHlGaWVsZCBpcyB0cnVlKVxuICogQHBhcmFtIHtzdHJpbmd9IFtjYXRlZ29yeS5yZXBseUJ1dHRvblRpdGxlXSAtIFRpdGxlIGZvciB0aGUgcmVwbHkgYnV0dG9uIChyZXF1aXJlZCBpZiBoYXNSZXBseUZpZWxkIGlzIHRydWUpXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVnaXN0ZXJOb3RpZmljYXRpb25DYXRlZ29yeShjYXRlZ29yeSkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlZ2lzdGVyTm90aWZpY2F0aW9uQ2F0ZWdvcnlcIiwgW2NhdGVnb3J5XSk7XG59XG5cbi8qKlxuICogUmVtb3ZlIGEgcHJldmlvdXNseSByZWdpc3RlcmVkIG5vdGlmaWNhdGlvbiBjYXRlZ29yeS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gY2F0ZWdvcnlJZCAtIFRoZSBJRCBvZiB0aGUgY2F0ZWdvcnkgdG8gcmVtb3ZlXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVtb3ZlTm90aWZpY2F0aW9uQ2F0ZWdvcnkoY2F0ZWdvcnlJZCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlbW92ZU5vdGlmaWNhdGlvbkNhdGVnb3J5XCIsIFtjYXRlZ29yeUlkXSk7XG59XG5cbi8qKlxuICogUmVtb3ZlIGFsbCBwZW5kaW5nIG5vdGlmaWNhdGlvbnMgZnJvbSB0aGUgbm90aWZpY2F0aW9uIGNlbnRlci5cbiAqIE9uIFdpbmRvd3MsIHRoaXMgaXMgYSBuby1vcCBhcyB0aGUgcGxhdGZvcm0gbWFuYWdlcyBub3RpZmljYXRpb24gbGlmZWN5Y2xlIGF1dG9tYXRpY2FsbHkuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbW92ZUFsbFBlbmRpbmdOb3RpZmljYXRpb25zKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlbW92ZUFsbFBlbmRpbmdOb3RpZmljYXRpb25zXCIpO1xufVxuXG4vKipcbiAqIFJlbW92ZSBhIHNwZWNpZmljIHBlbmRpbmcgbm90aWZpY2F0aW9uIGJ5IGl0cyBpZGVudGlmaWVyLlxuICogT24gV2luZG93cywgdGhpcyBpcyBhIG5vLW9wIGFzIHRoZSBwbGF0Zm9ybSBtYW5hZ2VzIG5vdGlmaWNhdGlvbiBsaWZlY3ljbGUgYXV0b21hdGljYWxseS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gaWRlbnRpZmllciAtIFRoZSBJRCBvZiB0aGUgbm90aWZpY2F0aW9uIHRvIHJlbW92ZVxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbW92ZVBlbmRpbmdOb3RpZmljYXRpb24oaWRlbnRpZmllcikge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlbW92ZVBlbmRpbmdOb3RpZmljYXRpb25cIiwgW2lkZW50aWZpZXJdKTtcbn1cblxuLyoqXG4gKiBSZW1vdmUgYWxsIGRlbGl2ZXJlZCBub3RpZmljYXRpb25zIGZyb20gdGhlIG5vdGlmaWNhdGlvbiBjZW50ZXIuXG4gKiBPbiBXaW5kb3dzLCB0aGlzIGlzIGEgbm8tb3AgYXMgdGhlIHBsYXRmb3JtIG1hbmFnZXMgbm90aWZpY2F0aW9uIGxpZmVjeWNsZSBhdXRvbWF0aWNhbGx5LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZW1vdmVBbGxEZWxpdmVyZWROb3RpZmljYXRpb25zKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlbW92ZUFsbERlbGl2ZXJlZE5vdGlmaWNhdGlvbnNcIik7XG59XG5cbi8qKlxuICogUmVtb3ZlIGEgc3BlY2lmaWMgZGVsaXZlcmVkIG5vdGlmaWNhdGlvbiBieSBpdHMgaWRlbnRpZmllci5cbiAqIE9uIFdpbmRvd3MsIHRoaXMgaXMgYSBuby1vcCBhcyB0aGUgcGxhdGZvcm0gbWFuYWdlcyBub3RpZmljYXRpb24gbGlmZWN5Y2xlIGF1dG9tYXRpY2FsbHkuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGlkZW50aWZpZXIgLSBUaGUgSUQgb2YgdGhlIG5vdGlmaWNhdGlvbiB0byByZW1vdmVcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZW1vdmVEZWxpdmVyZWROb3RpZmljYXRpb24oaWRlbnRpZmllcikge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlbW92ZURlbGl2ZXJlZE5vdGlmaWNhdGlvblwiLCBbaWRlbnRpZmllcl0pO1xufVxuXG4vKipcbiAqIFJlbW92ZSBhIG5vdGlmaWNhdGlvbiBieSBpdHMgaWRlbnRpZmllci5cbiAqIFRoaXMgaXMgYSBjb252ZW5pZW5jZSBmdW5jdGlvbiB0aGF0IHdvcmtzIGFjcm9zcyBwbGF0Zm9ybXMuXG4gKiBPbiBtYWNPUywgdXNlIHRoZSBtb3JlIHNwZWNpZmljIFJlbW92ZVBlbmRpbmdOb3RpZmljYXRpb24gb3IgUmVtb3ZlRGVsaXZlcmVkTm90aWZpY2F0aW9uIGZ1bmN0aW9ucy5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gaWRlbnRpZmllciAtIFRoZSBJRCBvZiB0aGUgbm90aWZpY2F0aW9uIHRvIHJlbW92ZVxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbW92ZU5vdGlmaWNhdGlvbihpZGVudGlmaWVyKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVtb3ZlTm90aWZpY2F0aW9uXCIsIFtpZGVudGlmaWVyXSk7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XG5pbXBvcnQge1xuICBldmVudExpc3RlbmVycyxcbiAgRXZlbnRzRW1pdCxcbiAgRXZlbnRzTm90aWZ5LFxuICBFdmVudHNPZmYsXG4gIEV2ZW50c09mZkFsbCxcbiAgRXZlbnRzT24sXG4gIEV2ZW50c09uY2UsXG4gIEV2ZW50c09uTXVsdGlwbGUsXG59IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHsgQ2FsbCwgQ2FsbGJhY2ssIGNhbGxiYWNrcyB9IGZyb20gJy4vY2FsbHMnO1xuaW1wb3J0IHsgU2V0QmluZGluZ3MgfSBmcm9tIFwiLi9iaW5kaW5nc1wiO1xuaW1wb3J0ICogYXMgV2luZG93IGZyb20gXCIuL3dpbmRvd1wiO1xuaW1wb3J0ICogYXMgU2NyZWVuIGZyb20gXCIuL3NjcmVlblwiO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkXCI7XG5pbXBvcnQgKiBhcyBEcmFnQW5kRHJvcCBmcm9tIFwiLi9kcmFnYW5kZHJvcFwiO1xuaW1wb3J0ICogYXMgQ29udGV4dE1lbnUgZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCAqIGFzIE5vdGlmaWNhdGlvbnMgZnJvbSBcIi4vbm90aWZpY2F0aW9uc1wiO1xuXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1EnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdTJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnSCcpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6RW52aXJvbm1lbnRcIik7XG59XG5cbi8vIFRoZSBKUyBydW50aW1lXG53aW5kb3cucnVudGltZSA9IHtcbiAgICAuLi5Mb2csXG4gICAgLi4uV2luZG93LFxuICAgIC4uLkJyb3dzZXIsXG4gICAgLi4uU2NyZWVuLFxuICAgIC4uLkNsaXBib2FyZCxcbiAgICAuLi5EcmFnQW5kRHJvcCxcbiAgICAuLi5Ob3RpZmljYXRpb25zLFxuICAgIEV2ZW50c09uLFxuICAgIEV2ZW50c09uY2UsXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcbiAgICBFdmVudHNFbWl0LFxuICAgIEV2ZW50c09mZixcbiAgICBFdmVudHNPZmZBbGwsXG4gICAgRW52aXJvbm1lbnQsXG4gICAgU2hvdyxcbiAgICBIaWRlLFxuICAgIFF1aXRcbn07XG5cbi8vIEludGVybmFsIHdhaWxzIGVuZHBvaW50c1xud2luZG93LndhaWxzID0ge1xuICAgIENhbGxiYWNrLFxuICAgIEV2ZW50c05vdGlmeSxcbiAgICBTZXRCaW5kaW5ncyxcbiAgICBldmVudExpc3RlbmVycyxcbiAgICBjYWxsYmFja3MsXG4gICAgZmxhZ3M6IHtcbiAgICAgICAgZGlzYWJsZVNjcm9sbGJhckRyYWc6IGZhbHNlLFxuICAgICAgICBkaXNhYmxlRGVmYXVsdENvbnRleHRNZW51OiBmYWxzZSxcbiAgICAgICAgZW5hYmxlUmVzaXplOiBmYWxzZSxcbiAgICAgICAgZGVmYXVsdEN1cnNvcjogbnVsbCxcbiAgICAgICAgYm9yZGVyVGhpY2tuZXNzOiA2LFxuICAgICAgICBzaG91bGREcmFnOiBmYWxzZSxcbiAgICAgICAgZGVmZXJEcmFnVG9Nb3VzZU1vdmU6IHRydWUsXG4gICAgICAgIGNzc0RyYWdQcm9wZXJ0eTogXCItLXdhaWxzLWRyYWdnYWJsZVwiLFxuICAgICAgICBjc3NEcmFnVmFsdWU6IFwiZHJhZ1wiLFxuICAgICAgICBjc3NEcm9wUHJvcGVydHk6IFwiLS13YWlscy1kcm9wLXRhcmdldFwiLFxuICAgICAgICBjc3NEcm9wVmFsdWU6IFwiZHJvcFwiLFxuICAgICAgICBlbmFibGVXYWlsc0RyYWdBbmREcm9wOiBmYWxzZSxcbiAgICB9XG59O1xuXG4vLyBTZXQgdGhlIGJpbmRpbmdzXG5pZiAod2luZG93LndhaWxzYmluZGluZ3MpIHtcbiAgICB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3Mod2luZG93LndhaWxzYmluZGluZ3MpO1xuICAgIGRlbGV0ZSB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3M7XG59XG5cbi8vIChib29sKSBUaGlzIGlzIGV2YWx1YXRlZCBhdCBidWlsZCB0aW1lIGluIHBhY2thZ2UuanNvblxuaWYgKCFERUJVRykge1xuICAgIGRlbGV0ZSB3aW5kb3cud2FpbHNiaW5kaW5ncztcbn1cblxubGV0IGRyYWdUZXN0ID0gZnVuY3Rpb24oZSkge1xuICAgIHZhciB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZSh3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1Byb3BlcnR5KTtcbiAgICBpZiAodmFsKSB7XG4gICAgICAgIHZhbCA9IHZhbC50cmltKCk7XG4gICAgfVxuXG4gICAgaWYgKHZhbCAhPT0gd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdWYWx1ZSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgaWYgKGUuYnV0dG9ucyAhPT0gMSkge1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZ2dpbmcgaWYgbm90IHRoZSBwcmltYXJ5IGJ1dHRvbiBoYXMgYmVlbiBjbGlja2VkLlxuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgaWYgKGUuZGV0YWlsICE9PSAxKSB7XG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBtb3JlIHRoYW4gb25jZSBoYXMgYmVlbiBjbGlja2VkLCBlLmcuIHdoZW4gZG91YmxlIGNsaWNraW5nXG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICByZXR1cm4gdHJ1ZTtcbn07XG5cbndpbmRvdy53YWlscy5zZXRDU1NEcmFnUHJvcGVydGllcyA9IGZ1bmN0aW9uKHByb3BlcnR5LCB2YWx1ZSkge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnUHJvcGVydHkgPSBwcm9wZXJ0eTtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1ZhbHVlID0gdmFsdWU7XG59XG5cbndpbmRvdy53YWlscy5zZXRDU1NEcm9wUHJvcGVydGllcyA9IGZ1bmN0aW9uKHByb3BlcnR5LCB2YWx1ZSkge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkgPSBwcm9wZXJ0eTtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFZhbHVlID0gdmFsdWU7XG59XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xuICAgIC8vIENoZWNrIGZvciByZXNpemluZ1xuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSkge1xuICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJyZXNpemU6XCIgKyB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSk7XG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmIChkcmFnVGVzdChlKSkge1xuICAgICAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRpc2FibGVTY3JvbGxiYXJEcmFnKSB7XG4gICAgICAgICAgICAvLyBUaGlzIGNoZWNrcyBmb3IgY2xpY2tzIG9uIHRoZSBzY3JvbGwgYmFyXG4gICAgICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmZXJEcmFnVG9Nb3VzZU1vdmUpIHtcbiAgICAgICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gdHJ1ZTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGUucHJldmVudERlZmF1bHQoKVxuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm47XG4gICAgfSBlbHNlIHtcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcbiAgICB9XG59KTtcblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCAoKSA9PiB7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcbn0pO1xuXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvcjtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSA9IGN1cnNvcjtcbn1cblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIGZ1bmN0aW9uKGUpIHtcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcpIHtcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSBmYWxzZTtcbiAgICAgICAgbGV0IG1vdXNlUHJlc3NlZCA9IGUuYnV0dG9ucyAhPT0gdW5kZWZpbmVkID8gZS5idXR0b25zIDogZS53aGljaDtcbiAgICAgICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlUmVzaXplKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yID09IG51bGwpIHtcbiAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yO1xuICAgIH1cbiAgICBpZiAod2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzICYmIHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MpIHtcbiAgICAgICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IFwic2UtcmVzaXplXCI7XG4gICAgfVxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG4gICAgbGV0IHRvcEJvcmRlciA9IGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG5cbiAgICAvLyBJZiB3ZSBhcmVuJ3Qgb24gYW4gZWRnZSwgYnV0IHdlcmUsIHJlc2V0IHRoZSBjdXJzb3IgdG8gZGVmYXVsdFxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfSBlbHNlIGlmIChyaWdodEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgdG9wQm9yZGVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyICYmIHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcblxufSk7XG5cbi8vIFNldHVwIGNvbnRleHQgbWVudSBob29rXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBmdW5jdGlvbihlKSB7XG4gICAgLy8gYWx3YXlzIHNob3cgdGhlIGNvbnRleHRtZW51IGluIGRlYnVnICYgZGV2XG4gICAgaWYgKERFQlVHKSByZXR1cm47XG5cbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRpc2FibGVEZWZhdWx0Q29udGV4dE1lbnUpIHtcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIENvbnRleHRNZW51LnByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZSk7XG4gICAgfVxufSk7XG5cbndpbmRvdy5XYWlsc0ludm9rZShcInJ1bnRpbWU6cmVhZHlcIik7Il0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFrQkEsV0FBUyxlQUFlLE9BQU8sU0FBUztBQUl2QyxXQUFPLFlBQVksTUFBTSxRQUFRLE9BQU87QUFBQSxFQUN6QztBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsUUFBUSxTQUFTO0FBQ2hDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxXQUFXLFNBQVM7QUFDbkMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxZQUFZLFVBQVU7QUFDckMsbUJBQWUsS0FBSyxRQUFRO0FBQUEsRUFDN0I7QUFHTyxNQUFNLFdBQVc7QUFBQSxJQUN2QixPQUFPO0FBQUEsSUFDUCxPQUFPO0FBQUEsSUFDUCxNQUFNO0FBQUEsSUFDTixTQUFTO0FBQUEsSUFDVCxPQUFPO0FBQUEsRUFDUjs7O0FDOUZBLE1BQU0sV0FBTixNQUFlO0FBQUEsSUFRWCxZQUFZLFdBQVcsVUFBVSxjQUFjO0FBQzNDLFdBQUssWUFBWTtBQUVqQixXQUFLLGVBQWUsZ0JBQWdCO0FBR3BDLFdBQUssV0FBVyxDQUFDLFNBQVM7QUFDdEIsaUJBQVMsTUFBTSxNQUFNLElBQUk7QUFFekIsWUFBSSxLQUFLLGlCQUFpQixJQUFJO0FBQzFCLGlCQUFPO0FBQUEsUUFDWDtBQUVBLGFBQUssZ0JBQWdCO0FBQ3JCLGVBQU8sS0FBSyxpQkFBaUI7QUFBQSxNQUNqQztBQUFBLElBQ0o7QUFBQSxFQUNKO0FBRU8sTUFBTSxpQkFBaUIsQ0FBQztBQVd4QixXQUFTLGlCQUFpQixXQUFXLFVBQVUsY0FBYztBQUNoRSxtQkFBZSxhQUFhLGVBQWUsY0FBYyxDQUFDO0FBQzFELFVBQU0sZUFBZSxJQUFJLFNBQVMsV0FBVyxVQUFVLFlBQVk7QUFDbkUsbUJBQWUsV0FBVyxLQUFLLFlBQVk7QUFDM0MsV0FBTyxNQUFNLFlBQVksWUFBWTtBQUFBLEVBQ3pDO0FBVU8sV0FBUyxTQUFTLFdBQVcsVUFBVTtBQUMxQyxXQUFPLGlCQUFpQixXQUFXLFVBQVUsRUFBRTtBQUFBLEVBQ25EO0FBVU8sV0FBUyxXQUFXLFdBQVcsVUFBVTtBQUM1QyxXQUFPLGlCQUFpQixXQUFXLFVBQVUsQ0FBQztBQUFBLEVBQ2xEO0FBRUEsV0FBUyxnQkFBZ0IsV0FBVztBQUdoQyxRQUFJLFlBQVksVUFBVTtBQUcxQixVQUFNLFlBQVksZUFBZSxZQUFZLE1BQU0sS0FBSyxDQUFDO0FBR3pELFFBQUksVUFBVSxRQUFRO0FBRWxCLFlBQU0sV0FBVyxvQkFBSSxJQUFJO0FBR3pCLGVBQVMsUUFBUSxVQUFVLFNBQVMsR0FBRyxTQUFTLEdBQUcsU0FBUyxHQUFHO0FBRzNELGNBQU0sV0FBVyxVQUFVO0FBRTNCLFlBQUksT0FBTyxVQUFVO0FBR3JCLGNBQU0sVUFBVSxTQUFTLFNBQVMsSUFBSTtBQUN0QyxZQUFJLFNBQVM7QUFDVCxtQkFBUyxJQUFJLFFBQVE7QUFBQSxRQUN6QjtBQUFBLE1BQ0o7QUFJQSxVQUFJLFNBQVMsT0FBTyxLQUFLLGVBQWUsWUFBWTtBQUNoRCx1QkFBZSxhQUFhLGVBQWUsV0FBVyxPQUFPLE9BQUssQ0FBQyxTQUFTLElBQUksQ0FBQyxDQUFDO0FBQ2xGLFlBQUksZUFBZSxXQUFXLFdBQVcsR0FBRztBQUN4Qyx5QkFBZSxTQUFTO0FBQUEsUUFDNUI7QUFBQSxNQUNKO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFTTyxXQUFTLGFBQWEsZUFBZTtBQUV4QyxRQUFJO0FBQ0osUUFBSTtBQUNBLGdCQUFVLEtBQUssTUFBTSxhQUFhO0FBQUEsSUFDdEMsU0FBUyxHQUFQO0FBQ0UsWUFBTSxRQUFRLG9DQUFvQztBQUNsRCxZQUFNLElBQUksTUFBTSxLQUFLO0FBQUEsSUFDekI7QUFDQSxvQkFBZ0IsT0FBTztBQUFBLEVBQzNCO0FBUU8sV0FBUyxXQUFXLFdBQVc7QUFFbEMsVUFBTSxVQUFVO0FBQUEsTUFDWixNQUFNO0FBQUEsTUFDTixNQUFNLENBQUMsRUFBRSxNQUFNLE1BQU0sU0FBUyxFQUFFLE1BQU0sQ0FBQztBQUFBLElBQzNDO0FBR0Esb0JBQWdCLE9BQU87QUFHdkIsV0FBTyxZQUFZLE9BQU8sS0FBSyxVQUFVLE9BQU8sQ0FBQztBQUFBLEVBQ3JEO0FBRUEsV0FBUyxlQUFlLFdBQVc7QUFFL0IsV0FBTyxlQUFlO0FBR3RCLFdBQU8sWUFBWSxPQUFPLFNBQVM7QUFBQSxFQUN2QztBQVNPLFdBQVMsVUFBVSxjQUFjLHNCQUFzQjtBQUMxRCxtQkFBZSxTQUFTO0FBRXhCLFFBQUkscUJBQXFCLFNBQVMsR0FBRztBQUNqQywyQkFBcUIsUUFBUSxDQUFBQSxlQUFhO0FBQ3RDLHVCQUFlQSxVQUFTO0FBQUEsTUFDNUIsQ0FBQztBQUFBLElBQ0w7QUFBQSxFQUNKO0FBS1EsV0FBUyxlQUFlO0FBQzVCLFVBQU0sYUFBYSxPQUFPLEtBQUssY0FBYztBQUM3QyxlQUFXLFFBQVEsZUFBYTtBQUM1QixxQkFBZSxTQUFTO0FBQUEsSUFDNUIsQ0FBQztBQUFBLEVBQ0w7QUFPQyxXQUFTLFlBQVksVUFBVTtBQUM1QixVQUFNLFlBQVksU0FBUztBQUMzQixRQUFJLGVBQWUsZUFBZTtBQUFXO0FBRzdDLG1CQUFlLGFBQWEsZUFBZSxXQUFXLE9BQU8sT0FBSyxNQUFNLFFBQVE7QUFHaEYsUUFBSSxlQUFlLFdBQVcsV0FBVyxHQUFHO0FBQ3hDLHFCQUFlLFNBQVM7QUFBQSxJQUM1QjtBQUFBLEVBQ0o7OztBQzdNTyxNQUFNLFlBQVksQ0FBQztBQU8xQixXQUFTLGVBQWU7QUFDdkIsUUFBSSxRQUFRLElBQUksWUFBWSxDQUFDO0FBQzdCLFdBQU8sT0FBTyxPQUFPLGdCQUFnQixLQUFLLEVBQUU7QUFBQSxFQUM3QztBQVFBLFdBQVMsY0FBYztBQUN0QixXQUFPLEtBQUssT0FBTyxJQUFJO0FBQUEsRUFDeEI7QUFHQSxNQUFJO0FBQ0osTUFBSSxPQUFPLFFBQVE7QUFDbEIsaUJBQWE7QUFBQSxFQUNkLE9BQU87QUFDTixpQkFBYTtBQUFBLEVBQ2Q7QUFpQk8sV0FBUyxLQUFLLE1BQU0sTUFBTSxTQUFTO0FBR3pDLFFBQUksV0FBVyxNQUFNO0FBQ3BCLGdCQUFVO0FBQUEsSUFDWDtBQUdBLFdBQU8sSUFBSSxRQUFRLFNBQVUsU0FBUyxRQUFRO0FBRzdDLFVBQUk7QUFDSixTQUFHO0FBQ0YscUJBQWEsT0FBTyxNQUFNLFdBQVc7QUFBQSxNQUN0QyxTQUFTLFVBQVU7QUFFbkIsVUFBSTtBQUVKLFVBQUksVUFBVSxHQUFHO0FBQ2hCLHdCQUFnQixXQUFXLFdBQVk7QUFDdEMsaUJBQU8sTUFBTSxhQUFhLE9BQU8sNkJBQTZCLFVBQVUsQ0FBQztBQUFBLFFBQzFFLEdBQUcsT0FBTztBQUFBLE1BQ1g7QUFHQSxnQkFBVSxjQUFjO0FBQUEsUUFDdkI7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0Q7QUFFQSxVQUFJO0FBQ0gsY0FBTSxVQUFVO0FBQUEsVUFDZjtBQUFBLFVBQ0E7QUFBQSxVQUNBO0FBQUEsUUFDRDtBQUdTLGVBQU8sWUFBWSxNQUFNLEtBQUssVUFBVSxPQUFPLENBQUM7QUFBQSxNQUNwRCxTQUFTLEdBQVA7QUFFRSxnQkFBUSxNQUFNLENBQUM7QUFBQSxNQUNuQjtBQUFBLElBQ0osQ0FBQztBQUFBLEVBQ0w7QUFFQSxTQUFPLGlCQUFpQixDQUFDLElBQUksTUFBTSxZQUFZO0FBRzNDLFFBQUksV0FBVyxNQUFNO0FBQ2pCLGdCQUFVO0FBQUEsSUFDZDtBQUdBLFdBQU8sSUFBSSxRQUFRLFNBQVUsU0FBUyxRQUFRO0FBRzFDLFVBQUk7QUFDSixTQUFHO0FBQ0MscUJBQWEsS0FBSyxNQUFNLFdBQVc7QUFBQSxNQUN2QyxTQUFTLFVBQVU7QUFFbkIsVUFBSTtBQUVKLFVBQUksVUFBVSxHQUFHO0FBQ2Isd0JBQWdCLFdBQVcsV0FBWTtBQUNuQyxpQkFBTyxNQUFNLG9CQUFvQixLQUFLLDZCQUE2QixVQUFVLENBQUM7QUFBQSxRQUNsRixHQUFHLE9BQU87QUFBQSxNQUNkO0FBR0EsZ0JBQVUsY0FBYztBQUFBLFFBQ3BCO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBRUEsVUFBSTtBQUNBLGNBQU0sVUFBVTtBQUFBLFVBQ3hCO0FBQUEsVUFDQTtBQUFBLFVBQ0E7QUFBQSxRQUNEO0FBR1MsZUFBTyxZQUFZLE1BQU0sS0FBSyxVQUFVLE9BQU8sQ0FBQztBQUFBLE1BQ3BELFNBQVMsR0FBUDtBQUVFLGdCQUFRLE1BQU0sQ0FBQztBQUFBLE1BQ25CO0FBQUEsSUFDSixDQUFDO0FBQUEsRUFDTDtBQVVPLFdBQVMsU0FBUyxpQkFBaUI7QUFFekMsUUFBSTtBQUNKLFFBQUk7QUFDSCxnQkFBVSxLQUFLLE1BQU0sZUFBZTtBQUFBLElBQ3JDLFNBQVMsR0FBUDtBQUNELFlBQU0sUUFBUSxvQ0FBb0MsRUFBRSxxQkFBcUI7QUFDekUsY0FBUSxTQUFTLEtBQUs7QUFDdEIsWUFBTSxJQUFJLE1BQU0sS0FBSztBQUFBLElBQ3RCO0FBQ0EsUUFBSSxhQUFhLFFBQVE7QUFDekIsUUFBSSxlQUFlLFVBQVU7QUFDN0IsUUFBSSxDQUFDLGNBQWM7QUFDbEIsWUFBTSxRQUFRLGFBQWE7QUFDM0IsY0FBUSxNQUFNLEtBQUs7QUFDbkIsWUFBTSxJQUFJLE1BQU0sS0FBSztBQUFBLElBQ3RCO0FBQ0EsaUJBQWEsYUFBYSxhQUFhO0FBRXZDLFdBQU8sVUFBVTtBQUVqQixRQUFJLFFBQVEsT0FBTztBQUNsQixtQkFBYSxPQUFPLFFBQVEsS0FBSztBQUFBLElBQ2xDLE9BQU87QUFDTixtQkFBYSxRQUFRLFFBQVEsTUFBTTtBQUFBLElBQ3BDO0FBQUEsRUFDRDs7O0FDMUtBLFNBQU8sS0FBSyxDQUFDO0FBRU4sV0FBUyxZQUFZLGFBQWE7QUFDeEMsUUFBSTtBQUNILG9CQUFjLEtBQUssTUFBTSxXQUFXO0FBQUEsSUFDckMsU0FBUyxHQUFQO0FBQ0QsY0FBUSxNQUFNLENBQUM7QUFBQSxJQUNoQjtBQUdBLFdBQU8sS0FBSyxPQUFPLE1BQU0sQ0FBQztBQUcxQixXQUFPLEtBQUssV0FBVyxFQUFFLFFBQVEsQ0FBQyxnQkFBZ0I7QUFHakQsYUFBTyxHQUFHLGVBQWUsT0FBTyxHQUFHLGdCQUFnQixDQUFDO0FBR3BELGFBQU8sS0FBSyxZQUFZLFlBQVksRUFBRSxRQUFRLENBQUMsZUFBZTtBQUc3RCxlQUFPLEdBQUcsYUFBYSxjQUFjLE9BQU8sR0FBRyxhQUFhLGVBQWUsQ0FBQztBQUU1RSxlQUFPLEtBQUssWUFBWSxhQUFhLFdBQVcsRUFBRSxRQUFRLENBQUMsZUFBZTtBQUV6RSxpQkFBTyxHQUFHLGFBQWEsWUFBWSxjQUFjLFdBQVk7QUFHNUQsZ0JBQUksVUFBVTtBQUdkLHFCQUFTLFVBQVU7QUFDbEIsb0JBQU0sT0FBTyxDQUFDLEVBQUUsTUFBTSxLQUFLLFNBQVM7QUFDcEMscUJBQU8sS0FBSyxDQUFDLGFBQWEsWUFBWSxVQUFVLEVBQUUsS0FBSyxHQUFHLEdBQUcsTUFBTSxPQUFPO0FBQUEsWUFDM0U7QUFHQSxvQkFBUSxhQUFhLFNBQVUsWUFBWTtBQUMxQyx3QkFBVTtBQUFBLFlBQ1g7QUFHQSxvQkFBUSxhQUFhLFdBQVk7QUFDaEMscUJBQU87QUFBQSxZQUNSO0FBRUEsbUJBQU87QUFBQSxVQUNSLEVBQUU7QUFBQSxRQUNILENBQUM7QUFBQSxNQUNGLENBQUM7QUFBQSxJQUNGLENBQUM7QUFBQSxFQUNGOzs7QUNsRUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFlTyxXQUFTLGVBQWU7QUFDM0IsV0FBTyxTQUFTLE9BQU87QUFBQSxFQUMzQjtBQUVPLFdBQVMsa0JBQWtCO0FBQzlCLFdBQU8sWUFBWSxJQUFJO0FBQUEsRUFDM0I7QUFFTyxXQUFTLDhCQUE4QjtBQUMxQyxXQUFPLFlBQVksT0FBTztBQUFBLEVBQzlCO0FBRU8sV0FBUyxzQkFBc0I7QUFDbEMsV0FBTyxZQUFZLE1BQU07QUFBQSxFQUM3QjtBQUVPLFdBQVMscUJBQXFCO0FBQ2pDLFdBQU8sWUFBWSxNQUFNO0FBQUEsRUFDN0I7QUFPTyxXQUFTLGVBQWU7QUFDM0IsV0FBTyxZQUFZLElBQUk7QUFBQSxFQUMzQjtBQVFPLFdBQVMsZUFBZSxPQUFPO0FBQ2xDLFdBQU8sWUFBWSxPQUFPLEtBQUs7QUFBQSxFQUNuQztBQU9PLFdBQVMsbUJBQW1CO0FBQy9CLFdBQU8sWUFBWSxJQUFJO0FBQUEsRUFDM0I7QUFPTyxXQUFTLHFCQUFxQjtBQUNqQyxXQUFPLFlBQVksSUFBSTtBQUFBLEVBQzNCO0FBUU8sV0FBUyxxQkFBcUI7QUFDakMsV0FBTyxLQUFLLDJCQUEyQjtBQUFBLEVBQzNDO0FBU08sV0FBUyxjQUFjLE9BQU8sUUFBUTtBQUN6QyxXQUFPLFlBQVksUUFBUSxRQUFRLE1BQU0sTUFBTTtBQUFBLEVBQ25EO0FBU08sV0FBUyxnQkFBZ0I7QUFDNUIsV0FBTyxLQUFLLHNCQUFzQjtBQUFBLEVBQ3RDO0FBU08sV0FBUyxpQkFBaUIsT0FBTyxRQUFRO0FBQzVDLFdBQU8sWUFBWSxRQUFRLFFBQVEsTUFBTSxNQUFNO0FBQUEsRUFDbkQ7QUFTTyxXQUFTLGlCQUFpQixPQUFPLFFBQVE7QUFDNUMsV0FBTyxZQUFZLFFBQVEsUUFBUSxNQUFNLE1BQU07QUFBQSxFQUNuRDtBQVNPLFdBQVMscUJBQXFCLEdBQUc7QUFFcEMsV0FBTyxZQUFZLFdBQVcsSUFBSSxNQUFNLElBQUk7QUFBQSxFQUNoRDtBQVlPLFdBQVMsa0JBQWtCLEdBQUcsR0FBRztBQUNwQyxXQUFPLFlBQVksUUFBUSxJQUFJLE1BQU0sQ0FBQztBQUFBLEVBQzFDO0FBUU8sV0FBUyxvQkFBb0I7QUFDaEMsV0FBTyxLQUFLLHFCQUFxQjtBQUFBLEVBQ3JDO0FBT08sV0FBUyxhQUFhO0FBQ3pCLFdBQU8sWUFBWSxJQUFJO0FBQUEsRUFDM0I7QUFPTyxXQUFTLGFBQWE7QUFDekIsV0FBTyxZQUFZLElBQUk7QUFBQSxFQUMzQjtBQU9PLFdBQVMsaUJBQWlCO0FBQzdCLFdBQU8sWUFBWSxJQUFJO0FBQUEsRUFDM0I7QUFPTyxXQUFTLHVCQUF1QjtBQUNuQyxXQUFPLFlBQVksSUFBSTtBQUFBLEVBQzNCO0FBT08sV0FBUyxtQkFBbUI7QUFDL0IsV0FBTyxZQUFZLElBQUk7QUFBQSxFQUMzQjtBQVFPLFdBQVMsb0JBQW9CO0FBQ2hDLFdBQU8sS0FBSywwQkFBMEI7QUFBQSxFQUMxQztBQU9PLFdBQVMsaUJBQWlCO0FBQzdCLFdBQU8sWUFBWSxJQUFJO0FBQUEsRUFDM0I7QUFPTyxXQUFTLG1CQUFtQjtBQUMvQixXQUFPLFlBQVksSUFBSTtBQUFBLEVBQzNCO0FBUU8sV0FBUyxvQkFBb0I7QUFDaEMsV0FBTyxLQUFLLDBCQUEwQjtBQUFBLEVBQzFDO0FBUU8sV0FBUyxpQkFBaUI7QUFDN0IsV0FBTyxLQUFLLHVCQUF1QjtBQUFBLEVBQ3ZDO0FBV08sV0FBUywwQkFBMEIsR0FBRyxHQUFHLEdBQUcsR0FBRztBQUNsRCxRQUFJLE9BQU8sS0FBSyxVQUFVLEVBQUMsR0FBRyxLQUFLLEdBQUcsR0FBRyxLQUFLLEdBQUcsR0FBRyxLQUFLLEdBQUcsR0FBRyxLQUFLLElBQUcsQ0FBQztBQUN4RSxXQUFPLFlBQVksUUFBUSxJQUFJO0FBQUEsRUFDbkM7OztBQzNRQTtBQUFBO0FBQUE7QUFBQTtBQXNCTyxXQUFTLGVBQWU7QUFDM0IsV0FBTyxLQUFLLHFCQUFxQjtBQUFBLEVBQ3JDOzs7QUN4QkE7QUFBQTtBQUFBO0FBQUE7QUFLTyxXQUFTLGVBQWUsS0FBSztBQUNsQyxXQUFPLFlBQVksUUFBUSxHQUFHO0FBQUEsRUFDaEM7OztBQ1BBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFvQk8sV0FBUyxpQkFBaUIsTUFBTTtBQUNuQyxXQUFPLEtBQUssMkJBQTJCLENBQUMsSUFBSSxDQUFDO0FBQUEsRUFDakQ7QUFTTyxXQUFTLG1CQUFtQjtBQUMvQixXQUFPLEtBQUsseUJBQXlCO0FBQUEsRUFDekM7OztBQ2pDQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWNBLE1BQU0sUUFBUTtBQUFBLElBQ1YsWUFBWTtBQUFBLElBQ1osc0JBQXNCO0FBQUEsSUFDdEIsZUFBZTtBQUFBLElBQ2YsZ0JBQWdCO0FBQUEsSUFDaEIsdUJBQXVCO0FBQUEsRUFDM0I7QUFFQSxNQUFNLHFCQUFxQjtBQVEzQixXQUFTLHFCQUFxQixPQUFPO0FBQ2pDLFVBQU0sZUFBZSxNQUFNLGlCQUFpQixPQUFPLE1BQU0sTUFBTSxlQUFlLEVBQUUsS0FBSztBQUNyRixRQUFJLGNBQWM7QUFDZCxVQUFJLGlCQUFpQixPQUFPLE1BQU0sTUFBTSxjQUFjO0FBQ2xELGVBQU87QUFBQSxNQUNYO0FBSUEsYUFBTztBQUFBLElBQ1g7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQU9BLFdBQVMsV0FBVyxHQUFHO0FBSW5CLFVBQU0sYUFBYSxFQUFFLGFBQWEsTUFBTSxTQUFTLE9BQU87QUFHeEQsUUFBSSxDQUFDLFlBQVk7QUFDYjtBQUFBLElBQ0o7QUFHQSxNQUFFLGVBQWU7QUFDakIsTUFBRSxhQUFhLGFBQWE7QUFFNUIsUUFBSSxDQUFDLE9BQU8sTUFBTSxNQUFNLHdCQUF3QjtBQUM1QztBQUFBLElBQ0o7QUFFQSxRQUFJLENBQUMsTUFBTSxlQUFlO0FBQ3RCO0FBQUEsSUFDSjtBQUVBLFVBQU0sVUFBVSxFQUFFO0FBR2xCLFFBQUcsTUFBTTtBQUFnQixZQUFNLGVBQWU7QUFHOUMsUUFBSSxDQUFDLFdBQVcsQ0FBQyxxQkFBcUIsaUJBQWlCLE9BQU8sQ0FBQyxHQUFHO0FBQzlEO0FBQUEsSUFDSjtBQUVBLFFBQUksaUJBQWlCO0FBQ3JCLFdBQU8sZ0JBQWdCO0FBRW5CLFVBQUkscUJBQXFCLGlCQUFpQixjQUFjLENBQUMsR0FBRztBQUN4RCx1QkFBZSxVQUFVLElBQUksa0JBQWtCO0FBQUEsTUFDbkQ7QUFDQSx1QkFBaUIsZUFBZTtBQUFBLElBQ3BDO0FBQUEsRUFDSjtBQU9BLFdBQVMsWUFBWSxHQUFHO0FBRXBCLFVBQU0sYUFBYSxFQUFFLGFBQWEsTUFBTSxTQUFTLE9BQU87QUFHeEQsUUFBSSxDQUFDLFlBQVk7QUFDYjtBQUFBLElBQ0o7QUFHQSxNQUFFLGVBQWU7QUFFakIsUUFBSSxDQUFDLE9BQU8sTUFBTSxNQUFNLHdCQUF3QjtBQUM1QztBQUFBLElBQ0o7QUFFQSxRQUFJLENBQUMsTUFBTSxlQUFlO0FBQ3RCO0FBQUEsSUFDSjtBQUdBLFFBQUksQ0FBQyxFQUFFLFVBQVUsQ0FBQyxxQkFBcUIsaUJBQWlCLEVBQUUsTUFBTSxDQUFDLEdBQUc7QUFDaEUsYUFBTztBQUFBLElBQ1g7QUFHQSxRQUFHLE1BQU07QUFBZ0IsWUFBTSxlQUFlO0FBRzlDLFVBQU0saUJBQWlCLE1BQU07QUFFekIsWUFBTSxLQUFLLFNBQVMsdUJBQXVCLGtCQUFrQixDQUFDLEVBQUUsUUFBUSxRQUFNLEdBQUcsVUFBVSxPQUFPLGtCQUFrQixDQUFDO0FBRXJILFlBQU0saUJBQWlCO0FBRXZCLFVBQUksTUFBTSx1QkFBdUI7QUFDN0IscUJBQWEsTUFBTSxxQkFBcUI7QUFDeEMsY0FBTSx3QkFBd0I7QUFBQSxNQUNsQztBQUFBLElBQ0o7QUFHQSxVQUFNLHdCQUF3QixXQUFXLE1BQU07QUFDM0MsVUFBRyxNQUFNO0FBQWdCLGNBQU0sZUFBZTtBQUFBLElBQ2xELEdBQUcsRUFBRTtBQUFBLEVBQ1Q7QUFPQSxXQUFTLE9BQU8sR0FBRztBQUVmLFVBQU0sYUFBYSxFQUFFLGFBQWEsTUFBTSxTQUFTLE9BQU87QUFHeEQsUUFBSSxDQUFDLFlBQVk7QUFDYjtBQUFBLElBQ0o7QUFHQSxNQUFFLGVBQWU7QUFFakIsUUFBSSxDQUFDLE9BQU8sTUFBTSxNQUFNLHdCQUF3QjtBQUM1QztBQUFBLElBQ0o7QUFFQSxRQUFJLG9CQUFvQixHQUFHO0FBRXZCLFVBQUksUUFBUSxDQUFDO0FBQ2IsVUFBSSxFQUFFLGFBQWEsT0FBTztBQUN0QixnQkFBUSxDQUFDLEdBQUcsRUFBRSxhQUFhLEtBQUssRUFBRSxJQUFJLENBQUMsTUFBTSxNQUFNO0FBQy9DLGNBQUksS0FBSyxTQUFTLFFBQVE7QUFDdEIsbUJBQU8sS0FBSyxVQUFVO0FBQUEsVUFDMUI7QUFBQSxRQUNKLENBQUM7QUFBQSxNQUNMLE9BQU87QUFDSCxnQkFBUSxDQUFDLEdBQUcsRUFBRSxhQUFhLEtBQUs7QUFBQSxNQUNwQztBQUNBLGFBQU8sUUFBUSxpQkFBaUIsRUFBRSxHQUFHLEVBQUUsR0FBRyxLQUFLO0FBQUEsSUFDbkQ7QUFFQSxRQUFJLENBQUMsTUFBTSxlQUFlO0FBQ3RCO0FBQUEsSUFDSjtBQUdBLFFBQUcsTUFBTTtBQUFnQixZQUFNLGVBQWU7QUFHOUMsVUFBTSxLQUFLLFNBQVMsdUJBQXVCLGtCQUFrQixDQUFDLEVBQUUsUUFBUSxRQUFNLEdBQUcsVUFBVSxPQUFPLGtCQUFrQixDQUFDO0FBQUEsRUFDekg7QUFRTyxXQUFTLHNCQUFzQjtBQUNsQyxXQUFPLE9BQU8sUUFBUSxTQUFTLG9DQUFvQztBQUFBLEVBQ3ZFO0FBVU8sV0FBUyxpQkFBaUIsR0FBRyxHQUFHLE9BQU87QUFHMUMsUUFBSSxPQUFPLFFBQVEsU0FBUyxrQ0FBa0M7QUFDMUQsYUFBTyxRQUFRLGlDQUFpQyxhQUFhLEtBQUssS0FBSyxLQUFLO0FBQUEsSUFDaEY7QUFBQSxFQUNKO0FBbUJPLFdBQVMsV0FBVyxVQUFVLGVBQWU7QUFDaEQsUUFBSSxPQUFPLGFBQWEsWUFBWTtBQUNoQyxjQUFRLE1BQU0sdUNBQXVDO0FBQ3JEO0FBQUEsSUFDSjtBQUVBLFFBQUksTUFBTSxZQUFZO0FBQ2xCO0FBQUEsSUFDSjtBQUNBLFVBQU0sYUFBYTtBQUVuQixVQUFNLFFBQVEsT0FBTztBQUNyQixVQUFNLGdCQUFnQixVQUFVLGVBQWUsVUFBVSxZQUFZLE1BQU0sdUJBQXVCO0FBQ2xHLFdBQU8saUJBQWlCLFlBQVksVUFBVTtBQUM5QyxXQUFPLGlCQUFpQixhQUFhLFdBQVc7QUFDaEQsV0FBTyxpQkFBaUIsUUFBUSxNQUFNO0FBRXRDLFFBQUksS0FBSztBQUNULFFBQUksTUFBTSxlQUFlO0FBQ3JCLFdBQUssU0FBVSxHQUFHLEdBQUcsT0FBTztBQUN4QixjQUFNLFVBQVUsU0FBUyxpQkFBaUIsR0FBRyxDQUFDO0FBRTlDLFlBQUksQ0FBQyxXQUFXLENBQUMscUJBQXFCLGlCQUFpQixPQUFPLENBQUMsR0FBRztBQUM5RCxpQkFBTztBQUFBLFFBQ1g7QUFDQSxpQkFBUyxHQUFHLEdBQUcsS0FBSztBQUFBLE1BQ3hCO0FBQUEsSUFDSjtBQUVBLGFBQVMsbUJBQW1CLEVBQUU7QUFBQSxFQUNsQztBQUtPLFdBQVMsZ0JBQWdCO0FBQzVCLFdBQU8sb0JBQW9CLFlBQVksVUFBVTtBQUNqRCxXQUFPLG9CQUFvQixhQUFhLFdBQVc7QUFDbkQsV0FBTyxvQkFBb0IsUUFBUSxNQUFNO0FBQ3pDLGNBQVUsaUJBQWlCO0FBQzNCLFVBQU0sYUFBYTtBQUFBLEVBQ3ZCOzs7QUM1UU8sV0FBUywwQkFBMEIsT0FBTztBQUU3QyxVQUFNLFVBQVUsTUFBTTtBQUN0QixVQUFNLGdCQUFnQixPQUFPLGlCQUFpQixPQUFPO0FBQ3JELFVBQU0sMkJBQTJCLGNBQWMsaUJBQWlCLHVCQUF1QixFQUFFLEtBQUs7QUFDOUYsWUFBUSwwQkFBMEI7QUFBQSxNQUM5QixLQUFLO0FBQ0Q7QUFBQSxNQUNKLEtBQUs7QUFDRCxjQUFNLGVBQWU7QUFDckI7QUFBQSxNQUNKO0FBRUksWUFBSSxRQUFRLG1CQUFtQjtBQUMzQjtBQUFBLFFBQ0o7QUFHQSxjQUFNLFlBQVksT0FBTyxhQUFhO0FBQ3RDLGNBQU0sZUFBZ0IsVUFBVSxTQUFTLEVBQUUsU0FBUztBQUNwRCxZQUFJLGNBQWM7QUFDZCxtQkFBUyxJQUFJLEdBQUcsSUFBSSxVQUFVLFlBQVksS0FBSztBQUMzQyxrQkFBTSxRQUFRLFVBQVUsV0FBVyxDQUFDO0FBQ3BDLGtCQUFNLFFBQVEsTUFBTSxlQUFlO0FBQ25DLHFCQUFTLElBQUksR0FBRyxJQUFJLE1BQU0sUUFBUSxLQUFLO0FBQ25DLG9CQUFNLE9BQU8sTUFBTTtBQUNuQixrQkFBSSxTQUFTLGlCQUFpQixLQUFLLE1BQU0sS0FBSyxHQUFHLE1BQU0sU0FBUztBQUM1RDtBQUFBLGNBQ0o7QUFBQSxZQUNKO0FBQUEsVUFDSjtBQUFBLFFBQ0o7QUFFQSxZQUFJLFFBQVEsWUFBWSxXQUFXLFFBQVEsWUFBWSxZQUFZO0FBQy9ELGNBQUksZ0JBQWlCLENBQUMsUUFBUSxZQUFZLENBQUMsUUFBUSxVQUFXO0FBQzFEO0FBQUEsVUFDSjtBQUFBLFFBQ0o7QUFHQSxjQUFNLGVBQWU7QUFBQSxJQUM3QjtBQUFBLEVBQ0o7OztBQ2pEQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBcUJPLFdBQVMsMEJBQTBCO0FBQ3RDLFdBQU8sS0FBSyxnQ0FBZ0M7QUFBQSxFQUNoRDtBQVVPLFdBQVMsdUJBQXVCO0FBQ25DLFdBQU8sS0FBSyw2QkFBNkI7QUFBQSxFQUM3QztBQVFPLFdBQVMsMEJBQTBCO0FBQ3RDLFdBQU8sS0FBSyxnQ0FBZ0M7QUFBQSxFQUNoRDtBQVVPLFdBQVMsbUNBQW1DO0FBQy9DLFdBQU8sS0FBSyx5Q0FBeUM7QUFBQSxFQUN6RDtBQVVPLFdBQVMsaUNBQWlDO0FBQzdDLFdBQU8sS0FBSyx1Q0FBdUM7QUFBQSxFQUN2RDtBQWdCTyxXQUFTLGlCQUFpQixTQUFTO0FBQ3RDLFdBQU8sS0FBSywyQkFBMkIsQ0FBQyxPQUFPLENBQUM7QUFBQSxFQUNwRDtBQWtCTyxXQUFTLDRCQUE0QixTQUFTO0FBQ2pELFdBQU8sS0FBSyxzQ0FBc0MsQ0FBQyxPQUFPLENBQUM7QUFBQSxFQUMvRDtBQW1CTyxXQUFTLDZCQUE2QixVQUFVO0FBQ25ELFdBQU8sS0FBSyx1Q0FBdUMsQ0FBQyxRQUFRLENBQUM7QUFBQSxFQUNqRTtBQVNPLFdBQVMsMkJBQTJCLFlBQVk7QUFDbkQsV0FBTyxLQUFLLHFDQUFxQyxDQUFDLFVBQVUsQ0FBQztBQUFBLEVBQ2pFO0FBU08sV0FBUyxnQ0FBZ0M7QUFDNUMsV0FBTyxLQUFLLHNDQUFzQztBQUFBLEVBQ3REO0FBVU8sV0FBUywwQkFBMEIsWUFBWTtBQUNsRCxXQUFPLEtBQUssb0NBQW9DLENBQUMsVUFBVSxDQUFDO0FBQUEsRUFDaEU7QUFTTyxXQUFTLGtDQUFrQztBQUM5QyxXQUFPLEtBQUssd0NBQXdDO0FBQUEsRUFDeEQ7QUFVTyxXQUFTLDRCQUE0QixZQUFZO0FBQ3BELFdBQU8sS0FBSyxzQ0FBc0MsQ0FBQyxVQUFVLENBQUM7QUFBQSxFQUNsRTtBQVdPLFdBQVMsbUJBQW1CLFlBQVk7QUFDM0MsV0FBTyxLQUFLLDZCQUE2QixDQUFDLFVBQVUsQ0FBQztBQUFBLEVBQ3pEOzs7QUN2S08sV0FBUyxPQUFPO0FBQ25CLFdBQU8sWUFBWSxHQUFHO0FBQUEsRUFDMUI7QUFFTyxXQUFTLE9BQU87QUFDbkIsV0FBTyxZQUFZLEdBQUc7QUFBQSxFQUMxQjtBQUVPLFdBQVMsT0FBTztBQUNuQixXQUFPLFlBQVksR0FBRztBQUFBLEVBQzFCO0FBRU8sV0FBUyxjQUFjO0FBQzFCLFdBQU8sS0FBSyxvQkFBb0I7QUFBQSxFQUNwQztBQUdBLFNBQU8sVUFBVTtBQUFBLElBQ2IsR0FBRztBQUFBLElBQ0gsR0FBRztBQUFBLElBQ0gsR0FBRztBQUFBLElBQ0gsR0FBRztBQUFBLElBQ0gsR0FBRztBQUFBLElBQ0gsR0FBRztBQUFBLElBQ0gsR0FBRztBQUFBLElBQ0g7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxFQUNKO0FBR0EsU0FBTyxRQUFRO0FBQUEsSUFDWDtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBLE9BQU87QUFBQSxNQUNILHNCQUFzQjtBQUFBLE1BQ3RCLDJCQUEyQjtBQUFBLE1BQzNCLGNBQWM7QUFBQSxNQUNkLGVBQWU7QUFBQSxNQUNmLGlCQUFpQjtBQUFBLE1BQ2pCLFlBQVk7QUFBQSxNQUNaLHNCQUFzQjtBQUFBLE1BQ3RCLGlCQUFpQjtBQUFBLE1BQ2pCLGNBQWM7QUFBQSxNQUNkLGlCQUFpQjtBQUFBLE1BQ2pCLGNBQWM7QUFBQSxNQUNkLHdCQUF3QjtBQUFBLElBQzVCO0FBQUEsRUFDSjtBQUdBLE1BQUksT0FBTyxlQUFlO0FBQ3RCLFdBQU8sTUFBTSxZQUFZLE9BQU8sYUFBYTtBQUM3QyxXQUFPLE9BQU8sTUFBTTtBQUFBLEVBQ3hCO0FBR0EsTUFBSSxPQUFRO0FBQ1IsV0FBTyxPQUFPO0FBQUEsRUFDbEI7QUFFQSxNQUFJLFdBQVcsU0FBUyxHQUFHO0FBQ3ZCLFFBQUksTUFBTSxPQUFPLGlCQUFpQixFQUFFLE1BQU0sRUFBRSxpQkFBaUIsT0FBTyxNQUFNLE1BQU0sZUFBZTtBQUMvRixRQUFJLEtBQUs7QUFDTCxZQUFNLElBQUksS0FBSztBQUFBLElBQ25CO0FBRUEsUUFBSSxRQUFRLE9BQU8sTUFBTSxNQUFNLGNBQWM7QUFDekMsYUFBTztBQUFBLElBQ1g7QUFFQSxRQUFJLEVBQUUsWUFBWSxHQUFHO0FBRWpCLGFBQU87QUFBQSxJQUNYO0FBRUEsUUFBSSxFQUFFLFdBQVcsR0FBRztBQUVoQixhQUFPO0FBQUEsSUFDWDtBQUVBLFdBQU87QUFBQSxFQUNYO0FBRUEsU0FBTyxNQUFNLHVCQUF1QixTQUFTLFVBQVUsT0FBTztBQUMxRCxXQUFPLE1BQU0sTUFBTSxrQkFBa0I7QUFDckMsV0FBTyxNQUFNLE1BQU0sZUFBZTtBQUFBLEVBQ3RDO0FBRUEsU0FBTyxNQUFNLHVCQUF1QixTQUFTLFVBQVUsT0FBTztBQUMxRCxXQUFPLE1BQU0sTUFBTSxrQkFBa0I7QUFDckMsV0FBTyxNQUFNLE1BQU0sZUFBZTtBQUFBLEVBQ3RDO0FBRUEsU0FBTyxpQkFBaUIsYUFBYSxDQUFDLE1BQU07QUFFeEMsUUFBSSxPQUFPLE1BQU0sTUFBTSxZQUFZO0FBQy9CLGFBQU8sWUFBWSxZQUFZLE9BQU8sTUFBTSxNQUFNLFVBQVU7QUFDNUQsUUFBRSxlQUFlO0FBQ2pCO0FBQUEsSUFDSjtBQUVBLFFBQUksU0FBUyxDQUFDLEdBQUc7QUFDYixVQUFJLE9BQU8sTUFBTSxNQUFNLHNCQUFzQjtBQUV6QyxZQUFJLEVBQUUsVUFBVSxFQUFFLE9BQU8sZUFBZSxFQUFFLFVBQVUsRUFBRSxPQUFPLGNBQWM7QUFDdkU7QUFBQSxRQUNKO0FBQUEsTUFDSjtBQUNBLFVBQUksT0FBTyxNQUFNLE1BQU0sc0JBQXNCO0FBQ3pDLGVBQU8sTUFBTSxNQUFNLGFBQWE7QUFBQSxNQUNwQyxPQUFPO0FBQ0gsVUFBRSxlQUFlO0FBQ2pCLGVBQU8sWUFBWSxNQUFNO0FBQUEsTUFDN0I7QUFDQTtBQUFBLElBQ0osT0FBTztBQUNILGFBQU8sTUFBTSxNQUFNLGFBQWE7QUFBQSxJQUNwQztBQUFBLEVBQ0osQ0FBQztBQUVELFNBQU8saUJBQWlCLFdBQVcsTUFBTTtBQUNyQyxXQUFPLE1BQU0sTUFBTSxhQUFhO0FBQUEsRUFDcEMsQ0FBQztBQUVELFdBQVMsVUFBVSxRQUFRO0FBQ3ZCLGFBQVMsZ0JBQWdCLE1BQU0sU0FBUyxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBQ3JFLFdBQU8sTUFBTSxNQUFNLGFBQWE7QUFBQSxFQUNwQztBQUVBLFNBQU8saUJBQWlCLGFBQWEsU0FBUyxHQUFHO0FBQzdDLFFBQUksT0FBTyxNQUFNLE1BQU0sWUFBWTtBQUMvQixhQUFPLE1BQU0sTUFBTSxhQUFhO0FBQ2hDLFVBQUksZUFBZSxFQUFFLFlBQVksU0FBWSxFQUFFLFVBQVUsRUFBRTtBQUMzRCxVQUFJLGVBQWUsR0FBRztBQUNsQixlQUFPLFlBQVksTUFBTTtBQUN6QjtBQUFBLE1BQ0o7QUFBQSxJQUNKO0FBQ0EsUUFBSSxDQUFDLE9BQU8sTUFBTSxNQUFNLGNBQWM7QUFDbEM7QUFBQSxJQUNKO0FBQ0EsUUFBSSxPQUFPLE1BQU0sTUFBTSxpQkFBaUIsTUFBTTtBQUMxQyxhQUFPLE1BQU0sTUFBTSxnQkFBZ0IsU0FBUyxnQkFBZ0IsTUFBTTtBQUFBLElBQ3RFO0FBQ0EsUUFBSSxPQUFPLGFBQWEsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNLG1CQUFtQixPQUFPLGNBQWMsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNLGlCQUFpQjtBQUMzSSxlQUFTLGdCQUFnQixNQUFNLFNBQVM7QUFBQSxJQUM1QztBQUNBLFFBQUksY0FBYyxPQUFPLGFBQWEsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBQ3JFLFFBQUksYUFBYSxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU07QUFDaEQsUUFBSSxZQUFZLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUMvQyxRQUFJLGVBQWUsT0FBTyxjQUFjLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUd2RSxRQUFJLENBQUMsY0FBYyxDQUFDLGVBQWUsQ0FBQyxhQUFhLENBQUMsZ0JBQWdCLE9BQU8sTUFBTSxNQUFNLGVBQWUsUUFBVztBQUMzRyxnQkFBVTtBQUFBLElBQ2QsV0FBVyxlQUFlO0FBQWMsZ0JBQVUsV0FBVztBQUFBLGFBQ3BELGNBQWM7QUFBYyxnQkFBVSxXQUFXO0FBQUEsYUFDakQsY0FBYztBQUFXLGdCQUFVLFdBQVc7QUFBQSxhQUM5QyxhQUFhO0FBQWEsZ0JBQVUsV0FBVztBQUFBLGFBQy9DO0FBQVksZ0JBQVUsVUFBVTtBQUFBLGFBQ2hDO0FBQVcsZ0JBQVUsVUFBVTtBQUFBLGFBQy9CO0FBQWMsZ0JBQVUsVUFBVTtBQUFBLGFBQ2xDO0FBQWEsZ0JBQVUsVUFBVTtBQUFBLEVBRTlDLENBQUM7QUFHRCxTQUFPLGlCQUFpQixlQUFlLFNBQVMsR0FBRztBQUUvQyxRQUFJO0FBQU87QUFFWCxRQUFJLE9BQU8sTUFBTSxNQUFNLDJCQUEyQjtBQUM5QyxRQUFFLGVBQWU7QUFBQSxJQUNyQixPQUFPO0FBQ0gsTUFBWSwwQkFBMEIsQ0FBQztBQUFBLElBQzNDO0FBQUEsRUFDSixDQUFDO0FBRUQsU0FBTyxZQUFZLGVBQWU7IiwKICAibmFtZXMiOiBbImV2ZW50TmFtZSJdCn0K
