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
      const err = message.error instanceof Error ? message.error : new Error(message.error);
      callbackData.reject(err);
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL2RyYWdhbmRkcm9wLmpzIiwgImRlc2t0b3AvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9ub3RpZmljYXRpb25zLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfICAgICAgIF9fICAgICAgXyBfX1xufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuLyoqXG4gKiBTZW5kcyBhIGxvZyBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHdpdGggdGhlIGdpdmVuIGxldmVsICsgbWVzc2FnZVxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBsZXZlbFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZnVuY3Rpb24gc2VuZExvZ01lc3NhZ2UobGV2ZWwsIG1lc3NhZ2UpIHtcblxuXHQvLyBMb2cgTWVzc2FnZSBmb3JtYXQ6XG5cdC8vIGxbdHlwZV1bbWVzc2FnZV1cblx0d2luZG93LldhaWxzSW52b2tlKCdMJyArIGxldmVsICsgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiB0cmFjZSBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nVHJhY2UobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnVCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1ByaW50KG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1AnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGRlYnVnIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dEZWJ1ZyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdEJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBpbmZvIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dJbmZvKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0knLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIHdhcm5pbmcgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1dhcm5pbmcobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnVycsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gZXJyb3IgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ0Vycm9yKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0UnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGZhdGFsIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dGYXRhbChtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdGJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgTG9nIGxldmVsIHRvIHRoZSBnaXZlbiBsb2cgbGV2ZWxcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gbG9nbGV2ZWxcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldExvZ0xldmVsKGxvZ2xldmVsKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdTJywgbG9nbGV2ZWwpO1xufVxuXG4vLyBMb2cgbGV2ZWxzXG5leHBvcnQgY29uc3QgTG9nTGV2ZWwgPSB7XG5cdFRSQUNFOiAxLFxuXHRERUJVRzogMixcblx0SU5GTzogMyxcblx0V0FSTklORzogNCxcblx0RVJST1I6IDUsXG59O1xuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XG5cbiAgICAvLyBHZXQgdGhlIGV2ZW50IG5hbWVcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XG5cbiAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lciBpbmRleGVzIHRvIGRlc3Ryb3lcbiAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0/LnNsaWNlKCkgfHwgW107XG5cbiAgICAvLyBDaGVjayBpZiB3ZSBoYXZlIGFueSBsaXN0ZW5lcnMgZm9yIHRoaXMgZXZlbnRcbiAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoKSB7XG5cbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSBuZXdFdmVudExpc3RlbmVyTGlzdC5sZW5ndGggLSAxOyBjb3VudCA+PSAwOyBjb3VudCAtPSAxKSB7XG5cbiAgICAgICAgICAgIC8vIEdldCBuZXh0IGxpc3RlbmVyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0W2NvdW50XTtcblxuICAgICAgICAgICAgbGV0IGRhdGEgPSBldmVudERhdGEuZGF0YTtcblxuICAgICAgICAgICAgLy8gRG8gdGhlIGNhbGxiYWNrXG4gICAgICAgICAgICBjb25zdCBkZXN0cm95ID0gbGlzdGVuZXIuQ2FsbGJhY2soZGF0YSk7XG4gICAgICAgICAgICBpZiAoZGVzdHJveSkge1xuICAgICAgICAgICAgICAgIC8vIGlmIHRoZSBsaXN0ZW5lciBpbmRpY2F0ZWQgdG8gZGVzdHJveSBpdHNlbGYsIGFkZCBpdCB0byB0aGUgZGVzdHJveSBsaXN0XG4gICAgICAgICAgICAgICAgbmV3RXZlbnRMaXN0ZW5lckxpc3Quc3BsaWNlKGNvdW50LCAxKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFVwZGF0ZSBjYWxsYmFja3Mgd2l0aCBuZXcgbGlzdCBvZiBsaXN0ZW5lcnNcbiAgICAgICAgaWYgKG5ld0V2ZW50TGlzdGVuZXJMaXN0Lmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBuZXdFdmVudExpc3RlbmVyTGlzdDtcbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBOb3RpZnkgaW5mb3JtcyBmcm9udGVuZCBsaXN0ZW5lcnMgdGhhdCBhbiBldmVudCB3YXMgZW1pdHRlZCB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5vdGlmeU1lc3NhZ2UgLSBlbmNvZGVkIG5vdGlmaWNhdGlvbiBtZXNzYWdlXG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c05vdGlmeShub3RpZnlNZXNzYWdlKSB7XG4gICAgLy8gUGFyc2UgdGhlIG1lc3NhZ2VcbiAgICBsZXQgbWVzc2FnZTtcbiAgICB0cnkge1xuICAgICAgICBtZXNzYWdlID0gSlNPTi5wYXJzZShub3RpZnlNZXNzYWdlKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGNvbnN0IGVycm9yID0gJ0ludmFsaWQgSlNPTiBwYXNzZWQgdG8gTm90aWZ5OiAnICsgbm90aWZ5TWVzc2FnZTtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGVycm9yKTtcbiAgICB9XG4gICAgbm90aWZ5TGlzdGVuZXJzKG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIEVtaXQgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c0VtaXQoZXZlbnROYW1lKSB7XG5cbiAgICBjb25zdCBwYXlsb2FkID0ge1xuICAgICAgICBuYW1lOiBldmVudE5hbWUsXG4gICAgICAgIGRhdGE6IFtdLnNsaWNlLmFwcGx5KGFyZ3VtZW50cykuc2xpY2UoMSksXG4gICAgfTtcblxuICAgIC8vIE5vdGlmeSBKUyBsaXN0ZW5lcnNcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XG5cbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG59XG5cbmZ1bmN0aW9uIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSkge1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lcnNcbiAgICBkZWxldGUgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXTtcblxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VYJyArIGV2ZW50TmFtZSk7XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVyZXMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKVxuXG4gICAgaWYgKGFkZGl0aW9uYWxFdmVudE5hbWVzLmxlbmd0aCA+IDApIHtcbiAgICAgICAgYWRkaXRpb25hbEV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKVxuICAgICAgICB9KVxuICAgIH1cbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYWxsIGV2ZW50IGxpc3RlbmVycyBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICovXG4gZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZkFsbCgpIHtcbiAgICBjb25zdCBldmVudE5hbWVzID0gT2JqZWN0LmtleXMoZXZlbnRMaXN0ZW5lcnMpO1xuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG4gICAgfSlcbn1cblxuLyoqXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIEV2ZW50c09uXG4gKlxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcbiAqL1xuIGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdID09PSB1bmRlZmluZWQpIHJldHVybjtcblxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcblxuICAgIC8vIENsZWFuIHVwIGlmIHRoZXJlIGFyZSBubyBldmVudCBsaXN0ZW5lcnMgbGVmdFxuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpO1xuICAgIH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuZXhwb3J0IGNvbnN0IGNhbGxiYWNrcyA9IHt9O1xuXG4vKipcbiAqIFJldHVybnMgYSBudW1iZXIgZnJvbSB0aGUgbmF0aXZlIGJyb3dzZXIgcmFuZG9tIGZ1bmN0aW9uXG4gKlxuICogQHJldHVybnMgbnVtYmVyXG4gKi9cbmZ1bmN0aW9uIGNyeXB0b1JhbmRvbSgpIHtcblx0dmFyIGFycmF5ID0gbmV3IFVpbnQzMkFycmF5KDEpO1xuXHRyZXR1cm4gd2luZG93LmNyeXB0by5nZXRSYW5kb21WYWx1ZXMoYXJyYXkpWzBdO1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBudW1iZXIgdXNpbmcgZGEgb2xkLXNrb29sIE1hdGguUmFuZG9tXG4gKiBJIGxpa2VzIHRvIGNhbGwgaXQgTE9MUmFuZG9tXG4gKlxuICogQHJldHVybnMgbnVtYmVyXG4gKi9cbmZ1bmN0aW9uIGJhc2ljUmFuZG9tKCkge1xuXHRyZXR1cm4gTWF0aC5yYW5kb20oKSAqIDkwMDcxOTkyNTQ3NDA5OTE7XG59XG5cbi8vIFBpY2sgYSByYW5kb20gbnVtYmVyIGZ1bmN0aW9uIGJhc2VkIG9uIGJyb3dzZXIgY2FwYWJpbGl0eVxudmFyIHJhbmRvbUZ1bmM7XG5pZiAod2luZG93LmNyeXB0bykge1xuXHRyYW5kb21GdW5jID0gY3J5cHRvUmFuZG9tO1xufSBlbHNlIHtcblx0cmFuZG9tRnVuYyA9IGJhc2ljUmFuZG9tO1xufVxuXG5cbi8qKlxuICogQ2FsbCBzZW5kcyBhIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgdG8gY2FsbCB0aGUgYmluZGluZyB3aXRoIHRoZVxuICogZ2l2ZW4gZGF0YS4gQSBwcm9taXNlIGlzIHJldHVybmVkIGFuZCB3aWxsIGJlIGNvbXBsZXRlZCB3aGVuIHRoZVxuICogYmFja2VuZCByZXNwb25kcy4gVGhpcyB3aWxsIGJlIHJlc29sdmVkIHdoZW4gdGhlIGNhbGwgd2FzIHN1Y2Nlc3NmdWxcbiAqIG9yIHJlamVjdGVkIGlmIGFuIGVycm9yIGlzIHBhc3NlZCBiYWNrLlxuICogVGhlcmUgaXMgYSB0aW1lb3V0IG1lY2hhbmlzbS4gSWYgdGhlIGNhbGwgZG9lc24ndCByZXNwb25kIGluIHRoZSBnaXZlblxuICogdGltZSAoaW4gbWlsbGlzZWNvbmRzKSB0aGVuIHRoZSBwcm9taXNlIGlzIHJlamVjdGVkLlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lXG4gKiBAcGFyYW0ge2FueT19IGFyZ3NcbiAqIEBwYXJhbSB7bnVtYmVyPX0gdGltZW91dFxuICogQHJldHVybnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGwobmFtZSwgYXJncywgdGltZW91dCkge1xuXG5cdC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxuXHRpZiAodGltZW91dCA9PSBudWxsKSB7XG5cdFx0dGltZW91dCA9IDA7XG5cdH1cblxuXHQvLyBDcmVhdGUgYSBwcm9taXNlXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cblx0XHQvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuXHRcdHZhciBjYWxsYmFja0lEO1xuXHRcdGRvIHtcblx0XHRcdGNhbGxiYWNrSUQgPSBuYW1lICsgJy0nICsgcmFuZG9tRnVuYygpO1xuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cblx0XHR2YXIgdGltZW91dEhhbmRsZTtcblx0XHQvLyBTZXQgdGltZW91dFxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xuXHRcdFx0dGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcblx0XHRcdH0sIHRpbWVvdXQpO1xuXHRcdH1cblxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXG5cdFx0Y2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcblx0XHRcdHJlamVjdDogcmVqZWN0LFxuXHRcdFx0cmVzb2x2ZTogcmVzb2x2ZVxuXHRcdH07XG5cblx0XHR0cnkge1xuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0bmFtZSxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdH07XG5cbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnQycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcbiAgICAgICAgfVxuICAgIH0pO1xufVxuXG53aW5kb3cuT2JmdXNjYXRlZENhbGwgPSAoaWQsIGFyZ3MsIHRpbWVvdXQpID0+IHtcblxuICAgIC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxuICAgIGlmICh0aW1lb3V0ID09IG51bGwpIHtcbiAgICAgICAgdGltZW91dCA9IDA7XG4gICAgfVxuXG4gICAgLy8gQ3JlYXRlIGEgcHJvbWlzZVxuICAgIHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cbiAgICAgICAgLy8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcbiAgICAgICAgdmFyIGNhbGxiYWNrSUQ7XG4gICAgICAgIGRvIHtcbiAgICAgICAgICAgIGNhbGxiYWNrSUQgPSBpZCArICctJyArIHJhbmRvbUZ1bmMoKTtcbiAgICAgICAgfSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcblxuICAgICAgICB2YXIgdGltZW91dEhhbmRsZTtcbiAgICAgICAgLy8gU2V0IHRpbWVvdXRcbiAgICAgICAgaWYgKHRpbWVvdXQgPiAwKSB7XG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKCdDYWxsIHRvIG1ldGhvZCAnICsgaWQgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcbiAgICAgICAgICAgIH0sIHRpbWVvdXQpO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gU3RvcmUgY2FsbGJhY2tcbiAgICAgICAgY2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuICAgICAgICAgICAgdGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcbiAgICAgICAgICAgIHJlamVjdDogcmVqZWN0LFxuICAgICAgICAgICAgcmVzb2x2ZTogcmVzb2x2ZVxuICAgICAgICB9O1xuXG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xuXHRcdFx0XHRpZCxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdH07XG5cbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnYycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcbiAgICAgICAgfVxuICAgIH0pO1xufTtcblxuXG4vKipcbiAqIENhbGxlZCBieSB0aGUgYmFja2VuZCB0byByZXR1cm4gZGF0YSB0byBhIHByZXZpb3VzbHkgY2FsbGVkXG4gKiBiaW5kaW5nIGludm9jYXRpb25cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gaW5jb21pbmdNZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsYmFjayhpbmNvbWluZ01lc3NhZ2UpIHtcblx0Ly8gUGFyc2UgdGhlIG1lc3NhZ2Vcblx0bGV0IG1lc3NhZ2U7XG5cdHRyeSB7XG5cdFx0bWVzc2FnZSA9IEpTT04ucGFyc2UoaW5jb21pbmdNZXNzYWdlKTtcblx0fSBjYXRjaCAoZSkge1xuXHRcdGNvbnN0IGVycm9yID0gYEludmFsaWQgSlNPTiBwYXNzZWQgdG8gY2FsbGJhY2s6ICR7ZS5tZXNzYWdlfS4gTWVzc2FnZTogJHtpbmNvbWluZ01lc3NhZ2V9YDtcblx0XHRydW50aW1lLkxvZ0RlYnVnKGVycm9yKTtcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuXHR9XG5cdGxldCBjYWxsYmFja0lEID0gbWVzc2FnZS5jYWxsYmFja2lkO1xuXHRsZXQgY2FsbGJhY2tEYXRhID0gY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xuXHRpZiAoIWNhbGxiYWNrRGF0YSkge1xuXHRcdGNvbnN0IGVycm9yID0gYENhbGxiYWNrICcke2NhbGxiYWNrSUR9JyBub3QgcmVnaXN0ZXJlZCEhIWA7XG5cdFx0Y29uc29sZS5lcnJvcihlcnJvcik7IC8vIGVzbGludC1kaXNhYmxlLWxpbmVcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuXHR9XG5cdGNsZWFyVGltZW91dChjYWxsYmFja0RhdGEudGltZW91dEhhbmRsZSk7XG5cblx0ZGVsZXRlIGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcblxuXHRpZiAobWVzc2FnZS5lcnJvcikge1xuXHRcdGNvbnN0IGVyciA9IG1lc3NhZ2UuZXJyb3IgaW5zdGFuY2VvZiBFcnJvciA/IG1lc3NhZ2UuZXJyb3IgOiBuZXcgRXJyb3IobWVzc2FnZS5lcnJvcik7XG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChlcnIpO1xuXHR9IGVsc2Uge1xuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX18gICAgXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gICkgXG58X18vfF9fL1xcX18sXy9fL18vX19fXy8gIFxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuaW1wb3J0IHtDYWxsfSBmcm9tICcuL2NhbGxzJztcblxuLy8gVGhpcyBpcyB3aGVyZSB3ZSBiaW5kIGdvIG1ldGhvZCB3cmFwcGVyc1xud2luZG93LmdvID0ge307XG5cbmV4cG9ydCBmdW5jdGlvbiBTZXRCaW5kaW5ncyhiaW5kaW5nc01hcCkge1xuXHR0cnkge1xuXHRcdGJpbmRpbmdzTWFwID0gSlNPTi5wYXJzZShiaW5kaW5nc01hcCk7XG5cdH0gY2F0Y2ggKGUpIHtcblx0XHRjb25zb2xlLmVycm9yKGUpO1xuXHR9XG5cblx0Ly8gSW5pdGlhbGlzZSB0aGUgYmluZGluZ3MgbWFwXG5cdHdpbmRvdy5nbyA9IHdpbmRvdy5nbyB8fCB7fTtcblxuXHQvLyBJdGVyYXRlIHBhY2thZ2UgbmFtZXNcblx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXApLmZvckVhY2goKHBhY2thZ2VOYW1lKSA9PiB7XG5cblx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3Rcblx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXSB8fCB7fTtcblxuXHRcdC8vIEl0ZXJhdGUgc3RydWN0IG5hbWVzXG5cdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdKS5mb3JFYWNoKChzdHJ1Y3ROYW1lKSA9PiB7XG5cblx0XHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxuXHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gfHwge307XG5cblx0XHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSkuZm9yRWFjaCgobWV0aG9kTmFtZSkgPT4ge1xuXG5cdFx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV1bbWV0aG9kTmFtZV0gPSBmdW5jdGlvbiAoKSB7XG5cblx0XHRcdFx0XHQvLyBObyB0aW1lb3V0IGJ5IGRlZmF1bHRcblx0XHRcdFx0XHRsZXQgdGltZW91dCA9IDA7XG5cblx0XHRcdFx0XHQvLyBBY3R1YWwgZnVuY3Rpb25cblx0XHRcdFx0XHRmdW5jdGlvbiBkeW5hbWljKCkge1xuXHRcdFx0XHRcdFx0Y29uc3QgYXJncyA9IFtdLnNsaWNlLmNhbGwoYXJndW1lbnRzKTtcblx0XHRcdFx0XHRcdHJldHVybiBDYWxsKFtwYWNrYWdlTmFtZSwgc3RydWN0TmFtZSwgbWV0aG9kTmFtZV0uam9pbignLicpLCBhcmdzLCB0aW1lb3V0KTtcblx0XHRcdFx0XHR9XG5cblx0XHRcdFx0XHQvLyBBbGxvdyBzZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cblx0XHRcdFx0XHRkeW5hbWljLnNldFRpbWVvdXQgPSBmdW5jdGlvbiAobmV3VGltZW91dCkge1xuXHRcdFx0XHRcdFx0dGltZW91dCA9IG5ld1RpbWVvdXQ7XG5cdFx0XHRcdFx0fTtcblxuXHRcdFx0XHRcdC8vIEFsbG93IGdldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxuXHRcdFx0XHRcdGR5bmFtaWMuZ2V0VGltZW91dCA9IGZ1bmN0aW9uICgpIHtcblx0XHRcdFx0XHRcdHJldHVybiB0aW1lb3V0O1xuXHRcdFx0XHRcdH07XG5cblx0XHRcdFx0XHRyZXR1cm4gZHluYW1pYztcblx0XHRcdFx0fSgpO1xuXHRcdFx0fSk7XG5cdFx0fSk7XG5cdH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93UmVsb2FkKCkge1xuICAgIHdpbmRvdy5sb2NhdGlvbi5yZWxvYWQoKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZEFwcCgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dSJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTeXN0ZW1EZWZhdWx0VGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVNEVCcpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TGlnaHRUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBTFQnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldERhcmtUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBRFQnKTtcbn1cblxuLyoqXG4gKiBQbGFjZSB0aGUgd2luZG93IGluIHRoZSBjZW50ZXIgb2YgdGhlIHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0NlbnRlcigpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1djJyk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgd2luZG93IHRpdGxlXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHRpdGxlXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRUaXRsZSh0aXRsZSkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1QnICsgdGl0bGUpO1xufVxuXG4vKipcbiAqIE1ha2VzIHRoZSB3aW5kb3cgZ28gZnVsbHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0Z1bGxzY3JlZW4oKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXRicpO1xufVxuXG4vKipcbiAqIFJldmVydHMgdGhlIHdpbmRvdyBmcm9tIGZ1bGxzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbmZ1bGxzY3JlZW4oKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXZicpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIGluIGZ1bGwgc2NyZWVuIG1vZGUgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc0Z1bGxzY3JlZW4oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNGdWxsc2NyZWVuXCIpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dzOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHt3OiBudW1iZXIsIGg6IG51bWJlcn0+fSBUaGUgc2l6ZSBvZiB0aGUgd2luZG93XG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFNpemUoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0U2l6ZVwiKTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNYXhTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1daOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cbi8qKlxuICogU2V0IHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXejonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG5cblxuLyoqXG4gKiBTZXQgdGhlIHdpbmRvdyBBbHdheXNPblRvcCBvciBub3Qgb24gdG9wXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QWx3YXlzT25Ub3AoYikge1xuXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVRQOicgKyAoYiA/ICcxJyA6ICcwJykpO1xufVxuXG5cblxuXG4vKipcbiAqIFNldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gKiBAcGFyYW0ge251bWJlcn0geVxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0UG9zaXRpb24oeCwgeSkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3A6JyArIHggKyAnOicgKyB5KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHt4OiBudW1iZXIsIHk6IG51bWJlcn0+fSBUaGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0UG9zaXRpb24oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0UG9zXCIpO1xufVxuXG4vKipcbiAqIEhpZGUgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0hpZGUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXSCcpO1xufVxuXG4vKipcbiAqIFNob3cgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1Nob3coKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUycpO1xufVxuXG4vKipcbiAqIE1heGltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNYXhpbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dNJyk7XG59XG5cbi8qKlxuICogVG9nZ2xlIHRoZSBNYXhpbWlzZSBvZiB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VG9nZ2xlTWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdCcpO1xufVxuXG4vKipcbiAqIFVubWF4aW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVScpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1heGltaXNlZCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWF4aW1pc2VkKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWF4aW1pc2VkXCIpO1xufVxuXG4vKipcbiAqIE1pbmltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNaW5pbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dtJyk7XG59XG5cbi8qKlxuICogVW5taW5pbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5taW5pbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d1Jyk7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNNaW5pbWlzZWQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNaW5pbWlzZWRcIik7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbm9ybWFsIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNOb3JtYWwoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNOb3JtYWxcIik7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSBSIFJlZFxuICogQHBhcmFtIHtudW1iZXJ9IEcgR3JlZW5cbiAqIEBwYXJhbSB7bnVtYmVyfSBCIEJsdWVcbiAqIEBwYXJhbSB7bnVtYmVyfSBBIEFscGhhXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRCYWNrZ3JvdW5kQ29sb3VyKFIsIEcsIEIsIEEpIHtcbiAgICBsZXQgcmdiYSA9IEpTT04uc3RyaW5naWZ5KHtyOiBSIHx8IDAsIGc6IEcgfHwgMCwgYjogQiB8fCAwLCBhOiBBIHx8IDI1NX0pO1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cblxuLyoqXG4gKiBHZXRzIHRoZSBhbGwgc2NyZWVucy4gQ2FsbCB0aGlzIGFuZXcgZWFjaCB0aW1lIHlvdSB3YW50IHRvIHJlZnJlc2ggZGF0YSBmcm9tIHRoZSB1bmRlcmx5aW5nIHdpbmRvd2luZyBzeXN0ZW0uXG4gKiBAZXhwb3J0XG4gKiBAdHlwZWRlZiB7aW1wb3J0KCcuLi93cmFwcGVyL3J1bnRpbWUnKS5TY3JlZW59IFNjcmVlblxuICogQHJldHVybiB7UHJvbWlzZTx7U2NyZWVuW119Pn0gVGhlIHNjcmVlbnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNjcmVlbkdldEFsbCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpTY3JlZW5HZXRBbGxcIik7XG59XG4iLCAiLyoqXG4gKiBAZGVzY3JpcHRpb246IFVzZSB0aGUgc3lzdGVtIGRlZmF1bHQgYnJvd3NlciB0byBvcGVuIHRoZSB1cmxcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gQnJvd3Nlck9wZW5VUkwodXJsKSB7XG4gIHdpbmRvdy5XYWlsc0ludm9rZSgnQk86JyArIHVybCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuLyoqXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRTZXRUZXh0XCIsIFt0ZXh0XSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSB0ZXh0IGNvbnRlbnQgb2YgdGhlIGNsaXBib2FyZFxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8e3N0cmluZ30+fSBUZXh0IGNvbnRlbnQgb2YgdGhlIGNsaXBib2FyZFxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRHZXRUZXh0KCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsaXBib2FyZEdldFRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7RXZlbnRzT24sIEV2ZW50c09mZn0gZnJvbSBcIi4vZXZlbnRzXCI7XG5cbmNvbnN0IGZsYWdzID0ge1xuICAgIHJlZ2lzdGVyZWQ6IGZhbHNlLFxuICAgIGRlZmF1bHRVc2VEcm9wVGFyZ2V0OiB0cnVlLFxuICAgIHVzZURyb3BUYXJnZXQ6IHRydWUsXG4gICAgbmV4dERlYWN0aXZhdGU6IG51bGwsXG4gICAgbmV4dERlYWN0aXZhdGVUaW1lb3V0OiBudWxsLFxufTtcblxuY29uc3QgRFJPUF9UQVJHRVRfQUNUSVZFID0gXCJ3YWlscy1kcm9wLXRhcmdldC1hY3RpdmVcIjtcblxuLyoqXG4gKiBjaGVja1N0eWxlRHJvcFRhcmdldCBjaGVja3MgaWYgdGhlIHN0eWxlIGhhcyB0aGUgZHJvcCB0YXJnZXQgYXR0cmlidXRlXG4gKiBcbiAqIEBwYXJhbSB7Q1NTU3R5bGVEZWNsYXJhdGlvbn0gc3R5bGUgXG4gKiBAcmV0dXJucyBcbiAqL1xuZnVuY3Rpb24gY2hlY2tTdHlsZURyb3BUYXJnZXQoc3R5bGUpIHtcbiAgICBjb25zdCBjc3NEcm9wVmFsdWUgPSBzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkpLnRyaW0oKTtcbiAgICBpZiAoY3NzRHJvcFZhbHVlKSB7XG4gICAgICAgIGlmIChjc3NEcm9wVmFsdWUgPT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wVmFsdWUpIHtcbiAgICAgICAgICAgIHJldHVybiB0cnVlO1xuICAgICAgICB9XG4gICAgICAgIC8vIGlmIHRoZSBlbGVtZW50IGhhcyB0aGUgZHJvcCB0YXJnZXQgYXR0cmlidXRlLCBidXQgXG4gICAgICAgIC8vIHRoZSB2YWx1ZSBpcyBub3QgY29ycmVjdCwgdGVybWluYXRlIGZpbmRpbmcgcHJvY2Vzcy5cbiAgICAgICAgLy8gVGhpcyBjYW4gYmUgdXNlZnVsIHRvIGJsb2NrIHNvbWUgY2hpbGQgZWxlbWVudHMgZnJvbSBiZWluZyBkcm9wIHRhcmdldHMuXG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG4gICAgcmV0dXJuIGZhbHNlO1xufVxuXG4vKipcbiAqIG9uRHJhZ092ZXIgaXMgY2FsbGVkIHdoZW4gdGhlIGRyYWdvdmVyIGV2ZW50IGlzIGVtaXR0ZWQuXG4gKiBAcGFyYW0ge0RyYWdFdmVudH0gZVxuICogQHJldHVybnNcbiAqL1xuZnVuY3Rpb24gb25EcmFnT3ZlcihlKSB7XG4gICAgLy8gQ2hlY2sgaWYgdGhpcyBpcyBhbiBleHRlcm5hbCBmaWxlIGRyb3Agb3IgaW50ZXJuYWwgSFRNTCBkcmFnXG4gICAgLy8gRXh0ZXJuYWwgZmlsZSBkcm9wcyB3aWxsIGhhdmUgXCJGaWxlc1wiIGluIHRoZSB0eXBlcyBhcnJheVxuICAgIC8vIEludGVybmFsIEhUTUwgZHJhZ3MgdHlwaWNhbGx5IGhhdmUgXCJ0ZXh0L3BsYWluXCIsIFwidGV4dC9odG1sXCIgb3IgY3VzdG9tIHR5cGVzXG4gICAgY29uc3QgaXNGaWxlRHJvcCA9IGUuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKFwiRmlsZXNcIik7XG5cbiAgICAvLyBPbmx5IGhhbmRsZSBleHRlcm5hbCBmaWxlIGRyb3BzLCBsZXQgaW50ZXJuYWwgSFRNTDUgZHJhZy1hbmQtZHJvcCB3b3JrIG5vcm1hbGx5XG4gICAgaWYgKCFpc0ZpbGVEcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBBTFdBWVMgcHJldmVudCBkZWZhdWx0IGZvciBmaWxlIGRyb3BzIHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIGUuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG5cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVXYWlsc0RyYWdBbmREcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGNvbnN0IGVsZW1lbnQgPSBlLnRhcmdldDtcblxuICAgIC8vIFRyaWdnZXIgZGVib3VuY2UgZnVuY3Rpb24gdG8gZGVhY3RpdmF0ZSBkcm9wIHRhcmdldHNcbiAgICBpZihmbGFncy5uZXh0RGVhY3RpdmF0ZSkgZmxhZ3MubmV4dERlYWN0aXZhdGUoKTtcblxuICAgIC8vIGlmIHRoZSBlbGVtZW50IGlzIG51bGwgb3IgZWxlbWVudCBpcyBub3QgY2hpbGQgb2YgZHJvcCB0YXJnZXQgZWxlbWVudFxuICAgIGlmICghZWxlbWVudCB8fCAhY2hlY2tTdHlsZURyb3BUYXJnZXQoZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KSkpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCBjdXJyZW50RWxlbWVudCA9IGVsZW1lbnQ7XG4gICAgd2hpbGUgKGN1cnJlbnRFbGVtZW50KSB7XG4gICAgICAgIC8vIGNoZWNrIGlmIGN1cnJlbnRFbGVtZW50IGlzIGRyb3AgdGFyZ2V0IGVsZW1lbnRcbiAgICAgICAgaWYgKGNoZWNrU3R5bGVEcm9wVGFyZ2V0KGdldENvbXB1dGVkU3R5bGUoY3VycmVudEVsZW1lbnQpKSkge1xuICAgICAgICAgICAgY3VycmVudEVsZW1lbnQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkUpO1xuICAgICAgICB9XG4gICAgICAgIGN1cnJlbnRFbGVtZW50ID0gY3VycmVudEVsZW1lbnQucGFyZW50RWxlbWVudDtcbiAgICB9XG59XG5cbi8qKlxuICogb25EcmFnTGVhdmUgaXMgY2FsbGVkIHdoZW4gdGhlIGRyYWdsZWF2ZSBldmVudCBpcyBlbWl0dGVkLlxuICogQHBhcmFtIHtEcmFnRXZlbnR9IGVcbiAqIEByZXR1cm5zXG4gKi9cbmZ1bmN0aW9uIG9uRHJhZ0xlYXZlKGUpIHtcbiAgICAvLyBDaGVjayBpZiB0aGlzIGlzIGFuIGV4dGVybmFsIGZpbGUgZHJvcCBvciBpbnRlcm5hbCBIVE1MIGRyYWdcbiAgICBjb25zdCBpc0ZpbGVEcm9wID0gZS5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoXCJGaWxlc1wiKTtcblxuICAgIC8vIE9ubHkgaGFuZGxlIGV4dGVybmFsIGZpbGUgZHJvcHMsIGxldCBpbnRlcm5hbCBIVE1MNSBkcmFnLWFuZC1kcm9wIHdvcmsgbm9ybWFsbHlcbiAgICBpZiAoIWlzRmlsZURyb3ApIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIEFMV0FZUyBwcmV2ZW50IGRlZmF1bHQgZm9yIGZpbGUgZHJvcHMgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICBlLnByZXZlbnREZWZhdWx0KCk7XG5cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVXYWlsc0RyYWdBbmREcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIEZpbmQgdGhlIGNsb3NlIGRyb3AgdGFyZ2V0IGVsZW1lbnRcbiAgICBpZiAoIWUudGFyZ2V0IHx8ICFjaGVja1N0eWxlRHJvcFRhcmdldChnZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KSkpIHtcbiAgICAgICAgcmV0dXJuIG51bGw7XG4gICAgfVxuXG4gICAgLy8gVHJpZ2dlciBkZWJvdW5jZSBmdW5jdGlvbiB0byBkZWFjdGl2YXRlIGRyb3AgdGFyZ2V0c1xuICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xuICAgIFxuICAgIC8vIFVzZSBkZWJvdW5jZSB0ZWNobmlxdWUgdG8gdGFjbGUgZHJhZ2xlYXZlIGV2ZW50cyBvbiBvdmVybGFwcGluZyBlbGVtZW50cyBhbmQgZHJvcCB0YXJnZXQgZWxlbWVudHNcbiAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZSA9ICgpID0+IHtcbiAgICAgICAgLy8gRGVhY3RpdmF0ZSBhbGwgZHJvcCB0YXJnZXRzLCBuZXcgZHJvcCB0YXJnZXQgd2lsbCBiZSBhY3RpdmF0ZWQgb24gbmV4dCBkcmFnb3ZlciBldmVudFxuICAgICAgICBBcnJheS5mcm9tKGRvY3VtZW50LmdldEVsZW1lbnRzQnlDbGFzc05hbWUoRFJPUF9UQVJHRVRfQUNUSVZFKSkuZm9yRWFjaChlbCA9PiBlbC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRSkpO1xuICAgICAgICAvLyBSZXNldCBuZXh0RGVhY3RpdmF0ZVxuICAgICAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZSA9IG51bGw7XG4gICAgICAgIC8vIENsZWFyIHRpbWVvdXRcbiAgICAgICAgaWYgKGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCkge1xuICAgICAgICAgICAgY2xlYXJUaW1lb3V0KGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCk7XG4gICAgICAgICAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZVRpbWVvdXQgPSBudWxsO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gU2V0IHRpbWVvdXQgdG8gZGVhY3RpdmF0ZSBkcm9wIHRhcmdldHMgaWYgbm90IHRyaWdnZXJlZCBieSBuZXh0IGRyYWcgZXZlbnRcbiAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZVRpbWVvdXQgPSBzZXRUaW1lb3V0KCgpID0+IHtcbiAgICAgICAgaWYoZmxhZ3MubmV4dERlYWN0aXZhdGUpIGZsYWdzLm5leHREZWFjdGl2YXRlKCk7XG4gICAgfSwgNTApO1xufVxuXG4vKipcbiAqIG9uRHJvcCBpcyBjYWxsZWQgd2hlbiB0aGUgZHJvcCBldmVudCBpcyBlbWl0dGVkLlxuICogQHBhcmFtIHtEcmFnRXZlbnR9IGVcbiAqIEByZXR1cm5zXG4gKi9cbmZ1bmN0aW9uIG9uRHJvcChlKSB7XG4gICAgLy8gQ2hlY2sgaWYgdGhpcyBpcyBhbiBleHRlcm5hbCBmaWxlIGRyb3Agb3IgaW50ZXJuYWwgSFRNTCBkcmFnXG4gICAgY29uc3QgaXNGaWxlRHJvcCA9IGUuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKFwiRmlsZXNcIik7XG5cbiAgICAvLyBPbmx5IGhhbmRsZSBleHRlcm5hbCBmaWxlIGRyb3BzLCBsZXQgaW50ZXJuYWwgSFRNTDUgZHJhZy1hbmQtZHJvcCB3b3JrIG5vcm1hbGx5XG4gICAgaWYgKCFpc0ZpbGVEcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBBTFdBWVMgcHJldmVudCBkZWZhdWx0IGZvciBmaWxlIGRyb3BzIHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuXG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlV2FpbHNEcmFnQW5kRHJvcCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKENhblJlc29sdmVGaWxlUGF0aHMoKSkge1xuICAgICAgICAvLyBwcm9jZXNzIGZpbGVzXG4gICAgICAgIGxldCBmaWxlcyA9IFtdO1xuICAgICAgICBpZiAoZS5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcbiAgICAgICAgICAgIGZpbGVzID0gWy4uLmUuZGF0YVRyYW5zZmVyLml0ZW1zXS5tYXAoKGl0ZW0sIGkpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuIGl0ZW0uZ2V0QXNGaWxlKCk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBmaWxlcyA9IFsuLi5lLmRhdGFUcmFuc2Zlci5maWxlc107XG4gICAgICAgIH1cbiAgICAgICAgd2luZG93LnJ1bnRpbWUuUmVzb2x2ZUZpbGVQYXRocyhlLngsIGUueSwgZmlsZXMpO1xuICAgIH1cblxuICAgIGlmICghZmxhZ3MudXNlRHJvcFRhcmdldCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gVHJpZ2dlciBkZWJvdW5jZSBmdW5jdGlvbiB0byBkZWFjdGl2YXRlIGRyb3AgdGFyZ2V0c1xuICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xuXG4gICAgLy8gRGVhY3RpdmF0ZSBhbGwgZHJvcCB0YXJnZXRzXG4gICAgQXJyYXkuZnJvbShkb2N1bWVudC5nZXRFbGVtZW50c0J5Q2xhc3NOYW1lKERST1BfVEFSR0VUX0FDVElWRSkpLmZvckVhY2goZWwgPT4gZWwuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkUpKTtcbn1cblxuLyoqXG4gKiBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyBjaGVja3MgdGhlIGJyb3dzZXIncyBjYXBhYmlsaXR5IG9mIHNlbmRpbmcgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHNcbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn1cbiAqIEBjb25zdHJ1Y3RvclxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FuUmVzb2x2ZUZpbGVQYXRocygpIHtcbiAgICByZXR1cm4gd2luZG93LmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgIT0gbnVsbDtcbn1cblxuLyoqXG4gKiBSZXNvbHZlRmlsZVBhdGhzIHNlbmRzIGRyb3AgZXZlbnRzIHRvIHRoZSBHTyBzaWRlIHRvIHJlc29sdmUgZmlsZSBwYXRocyBvbiB3aW5kb3dzLlxuICpcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gKiBAcGFyYW0ge251bWJlcn0geVxuICogQHBhcmFtIHthbnlbXX0gZmlsZXNcbiAqIEBjb25zdHJ1Y3RvclxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVzb2x2ZUZpbGVQYXRocyh4LCB5LCBmaWxlcykge1xuICAgIC8vIE9ubHkgZm9yIHdpbmRvd3Mgd2VidmlldzIgPj0gMS4wLjE3NzQuMzBcbiAgICAvLyBodHRwczovL2xlYXJuLm1pY3Jvc29mdC5jb20vZW4tdXMvbWljcm9zb2Z0LWVkZ2Uvd2VidmlldzIvcmVmZXJlbmNlL3dpbjMyL2ljb3Jld2VidmlldzJ3ZWJtZXNzYWdlcmVjZWl2ZWRldmVudGFyZ3MyP3ZpZXc9d2VidmlldzItMS4wLjE4MjMuMzIjYXBwbGllcy10b1xuICAgIGlmICh3aW5kb3cuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cykge1xuICAgICAgICBjaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyhgZmlsZTpkcm9wOiR7eH06JHt5fWAsIGZpbGVzKTtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbGJhY2sgZm9yIE9uRmlsZURyb3AgcmV0dXJucyBhIHNsaWNlIG9mIGZpbGUgcGF0aCBzdHJpbmdzIHdoZW4gYSBkcm9wIGlzIGZpbmlzaGVkLlxuICpcbiAqIEBleHBvcnRcbiAqIEBjYWxsYmFjayBPbkZpbGVEcm9wQ2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSB4IC0geCBjb29yZGluYXRlIG9mIHRoZSBkcm9wXG4gKiBAcGFyYW0ge251bWJlcn0geSAtIHkgY29vcmRpbmF0ZSBvZiB0aGUgZHJvcFxuICogQHBhcmFtIHtzdHJpbmdbXX0gcGF0aHMgLSBBIGxpc3Qgb2YgZmlsZSBwYXRocy5cbiAqL1xuXG4vKipcbiAqIE9uRmlsZURyb3AgbGlzdGVucyB0byBkcmFnIGFuZCBkcm9wIGV2ZW50cyBhbmQgY2FsbHMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGNvb3JkaW5hdGVzIG9mIHRoZSBkcm9wIGFuZCBhbiBhcnJheSBvZiBwYXRoIHN0cmluZ3MuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtPbkZpbGVEcm9wQ2FsbGJhY2t9IGNhbGxiYWNrIC0gQ2FsbGJhY2sgZm9yIE9uRmlsZURyb3AgcmV0dXJucyBhIHNsaWNlIG9mIGZpbGUgcGF0aCBzdHJpbmdzIHdoZW4gYSBkcm9wIGlzIGZpbmlzaGVkLlxuICogQHBhcmFtIHtib29sZWFufSBbdXNlRHJvcFRhcmdldD10cnVlXSAtIE9ubHkgY2FsbCB0aGUgY2FsbGJhY2sgd2hlbiB0aGUgZHJvcCBmaW5pc2hlZCBvbiBhbiBlbGVtZW50IHRoYXQgaGFzIHRoZSBkcm9wIHRhcmdldCBzdHlsZS4gKC0td2FpbHMtZHJvcC10YXJnZXQpXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbkZpbGVEcm9wKGNhbGxiYWNrLCB1c2VEcm9wVGFyZ2V0KSB7XG4gICAgaWYgKHR5cGVvZiBjYWxsYmFjayAhPT0gXCJmdW5jdGlvblwiKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoXCJEcmFnQW5kRHJvcENhbGxiYWNrIGlzIG5vdCBhIGZ1bmN0aW9uXCIpO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKGZsYWdzLnJlZ2lzdGVyZWQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBmbGFncy5yZWdpc3RlcmVkID0gdHJ1ZTtcblxuICAgIGNvbnN0IHVEVFBUID0gdHlwZW9mIHVzZURyb3BUYXJnZXQ7XG4gICAgZmxhZ3MudXNlRHJvcFRhcmdldCA9IHVEVFBUID09PSBcInVuZGVmaW5lZFwiIHx8IHVEVFBUICE9PSBcImJvb2xlYW5cIiA/IGZsYWdzLmRlZmF1bHRVc2VEcm9wVGFyZ2V0IDogdXNlRHJvcFRhcmdldDtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBvbkRyYWdPdmVyKTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgb25EcmFnTGVhdmUpO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgb25Ecm9wKTtcblxuICAgIGxldCBjYiA9IGNhbGxiYWNrO1xuICAgIGlmIChmbGFncy51c2VEcm9wVGFyZ2V0KSB7XG4gICAgICAgIGNiID0gZnVuY3Rpb24gKHgsIHksIHBhdGhzKSB7XG4gICAgICAgICAgICBjb25zdCBlbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KVxuICAgICAgICAgICAgLy8gaWYgdGhlIGVsZW1lbnQgaXMgbnVsbCBvciBlbGVtZW50IGlzIG5vdCBjaGlsZCBvZiBkcm9wIHRhcmdldCBlbGVtZW50LCByZXR1cm4gbnVsbFxuICAgICAgICAgICAgaWYgKCFlbGVtZW50IHx8ICFjaGVja1N0eWxlRHJvcFRhcmdldChnZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpKSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBudWxsO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbGJhY2soeCwgeSwgcGF0aHMpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgRXZlbnRzT24oXCJ3YWlsczpmaWxlLWRyb3BcIiwgY2IpO1xufVxuXG4vKipcbiAqIE9uRmlsZURyb3BPZmYgcmVtb3ZlcyB0aGUgZHJhZyBhbmQgZHJvcCBsaXN0ZW5lcnMgYW5kIGhhbmRsZXJzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25GaWxlRHJvcE9mZigpIHtcbiAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBvbkRyYWdPdmVyKTtcbiAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgb25EcmFnTGVhdmUpO1xuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcm9wJywgb25Ecm9wKTtcbiAgICBFdmVudHNPZmYoXCJ3YWlsczpmaWxlLWRyb3BcIik7XG4gICAgZmxhZ3MucmVnaXN0ZXJlZCA9IGZhbHNlO1xufVxuIiwgIi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmV4cG9ydCBmdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KSB7XG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcbiAgICAgICAgY2FzZSBcInNob3dcIjpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQuaXNDb250ZW50RWRpdGFibGUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgYW5kIGFjdGlvbiBpcyBvbiB0aGUgc2VsZWN0ZWQgZWxlbWVudHNcbiAgICAgICAgICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICAgICAgICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IChzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKVxuICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbikge1xuICAgICAgICAgICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0cyA9IHJhbmdlLmdldENsaWVudFJlY3RzKCk7XG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHJlY3QubGVmdCwgcmVjdC50b3ApID09PSBlbGVtZW50KSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGFnbmFtZSBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQudGFnTmFtZSA9PT0gXCJJTlBVVFwiIHx8IGVsZW1lbnQudGFnTmFtZSA9PT0gXCJURVhUQVJFQVwiKSB7XG4gICAgICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIWVsZW1lbnQucmVhZE9ubHkgJiYgIWVsZW1lbnQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cbi8qKlxuICogSW5pdGlhbGl6ZSB0aGUgbm90aWZpY2F0aW9uIHNlcnZpY2UgZm9yIHRoZSBhcHBsaWNhdGlvbi5cbiAqIFRoaXMgbXVzdCBiZSBjYWxsZWQgYmVmb3JlIHNlbmRpbmcgYW55IG5vdGlmaWNhdGlvbnMuXG4gKiBPbiBtYWNPUywgdGhpcyBhbHNvIGVuc3VyZXMgdGhlIG5vdGlmaWNhdGlvbiBkZWxlZ2F0ZSBpcyBwcm9wZXJseSBpbml0aWFsaXplZC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gSW5pdGlhbGl6ZU5vdGlmaWNhdGlvbnMoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6SW5pdGlhbGl6ZU5vdGlmaWNhdGlvbnNcIik7XG59XG5cbi8qKlxuICogQ2xlYW4gdXAgbm90aWZpY2F0aW9uIHJlc291cmNlcyBhbmQgcmVsZWFzZSBhbnkgaGVsZCBjb25uZWN0aW9ucy5cbiAqIFRoaXMgc2hvdWxkIGJlIGNhbGxlZCB3aGVuIHNodXR0aW5nIGRvd24gdGhlIGFwcGxpY2F0aW9uIHRvIHByb3Blcmx5IHJlbGVhc2UgcmVzb3VyY2VzXG4gKiAocHJpbWFyaWx5IG5lZWRlZCBvbiBMaW51eCB0byBjbG9zZSBELUJ1cyBjb25uZWN0aW9ucykuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENsZWFudXBOb3RpZmljYXRpb25zKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsZWFudXBOb3RpZmljYXRpb25zXCIpO1xufVxuXG4vKipcbiAqIENoZWNrIGlmIG5vdGlmaWNhdGlvbnMgYXJlIGF2YWlsYWJsZSBvbiB0aGUgY3VycmVudCBwbGF0Zm9ybS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUcnVlIGlmIG5vdGlmaWNhdGlvbnMgYXJlIGF2YWlsYWJsZSwgZmFsc2Ugb3RoZXJ3aXNlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc05vdGlmaWNhdGlvbkF2YWlsYWJsZSgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpJc05vdGlmaWNhdGlvbkF2YWlsYWJsZVwiKTtcbn1cblxuLyoqXG4gKiBSZXF1ZXN0IG5vdGlmaWNhdGlvbiBhdXRob3JpemF0aW9uIGZyb20gdGhlIHVzZXIuXG4gKiBPbiBtYWNPUywgdGhpcyBwcm9tcHRzIHRoZSB1c2VyIHRvIGFsbG93IG5vdGlmaWNhdGlvbnMuXG4gKiBPbiBvdGhlciBwbGF0Zm9ybXMsIHRoaXMgYWx3YXlzIHJldHVybnMgdHJ1ZS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUcnVlIGlmIGF1dGhvcml6YXRpb24gd2FzIGdyYW50ZWQsIGZhbHNlIG90aGVyd2lzZVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVxdWVzdE5vdGlmaWNhdGlvbkF1dGhvcml6YXRpb24oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVxdWVzdE5vdGlmaWNhdGlvbkF1dGhvcml6YXRpb25cIik7XG59XG5cbi8qKlxuICogQ2hlY2sgdGhlIGN1cnJlbnQgbm90aWZpY2F0aW9uIGF1dGhvcml6YXRpb24gc3RhdHVzLlxuICogT24gbWFjT1MsIHRoaXMgY2hlY2tzIGlmIHRoZSBhcHAgaGFzIG5vdGlmaWNhdGlvbiBwZXJtaXNzaW9ucy5cbiAqIE9uIG90aGVyIHBsYXRmb3JtcywgdGhpcyBhbHdheXMgcmV0dXJucyB0cnVlLlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRydWUgaWYgYXV0aG9yaXplZCwgZmFsc2Ugb3RoZXJ3aXNlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDaGVja05vdGlmaWNhdGlvbkF1dGhvcml6YXRpb24oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6Q2hlY2tOb3RpZmljYXRpb25BdXRob3JpemF0aW9uXCIpO1xufVxuXG4vKipcbiAqIFNlbmQgYSBiYXNpYyBub3RpZmljYXRpb24gd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIFRoZSBub3RpZmljYXRpb24gd2lsbCBkaXNwbGF5IHdpdGggdGhlIHByb3ZpZGVkIHRpdGxlLCBzdWJ0aXRsZSAoaWYgc3VwcG9ydGVkKSwgYW5kIGJvZHkgdGV4dC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge09iamVjdH0gb3B0aW9ucyAtIE5vdGlmaWNhdGlvbiBvcHRpb25zXG4gKiBAcGFyYW0ge3N0cmluZ30gb3B0aW9ucy5pZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgbm90aWZpY2F0aW9uXG4gKiBAcGFyYW0ge3N0cmluZ30gb3B0aW9ucy50aXRsZSAtIE5vdGlmaWNhdGlvbiB0aXRsZVxuICogQHBhcmFtIHtzdHJpbmd9IFtvcHRpb25zLnN1YnRpdGxlXSAtIE5vdGlmaWNhdGlvbiBzdWJ0aXRsZSAobWFjT1MgYW5kIExpbnV4IG9ubHkpXG4gKiBAcGFyYW0ge3N0cmluZ30gW29wdGlvbnMuYm9keV0gLSBOb3RpZmljYXRpb24gYm9keSB0ZXh0XG4gKiBAcGFyYW0ge3N0cmluZ30gW29wdGlvbnMuY2F0ZWdvcnlJZF0gLSBDYXRlZ29yeSBJRCBmb3IgYWN0aW9uIGJ1dHRvbnMgKHJlcXVpcmVzIFNlbmROb3RpZmljYXRpb25XaXRoQWN0aW9ucylcbiAqIEBwYXJhbSB7T2JqZWN0PHN0cmluZywgYW55Pn0gW29wdGlvbnMuZGF0YV0gLSBBZGRpdGlvbmFsIHVzZXIgZGF0YSB0byBhdHRhY2ggdG8gdGhlIG5vdGlmaWNhdGlvblxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNlbmROb3RpZmljYXRpb24ob3B0aW9ucykge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlNlbmROb3RpZmljYXRpb25cIiwgW29wdGlvbnNdKTtcbn1cblxuLyoqXG4gKiBTZW5kIGEgbm90aWZpY2F0aW9uIHdpdGggYWN0aW9uIGJ1dHRvbnMuXG4gKiBBIE5vdGlmaWNhdGlvbkNhdGVnb3J5IG11c3QgYmUgcmVnaXN0ZXJlZCBmaXJzdCB1c2luZyBSZWdpc3Rlck5vdGlmaWNhdGlvbkNhdGVnb3J5LlxuICogVGhlIG9wdGlvbnMuY2F0ZWdvcnlJZCBtdXN0IG1hdGNoIGEgcHJldmlvdXNseSByZWdpc3RlcmVkIGNhdGVnb3J5IElELlxuICogSWYgdGhlIGNhdGVnb3J5IGlzIG5vdCBmb3VuZCwgYSBiYXNpYyBub3RpZmljYXRpb24gd2lsbCBiZSBzZW50IGluc3RlYWQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtPYmplY3R9IG9wdGlvbnMgLSBOb3RpZmljYXRpb24gb3B0aW9uc1xuICogQHBhcmFtIHtzdHJpbmd9IG9wdGlvbnMuaWQgLSBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIG5vdGlmaWNhdGlvblxuICogQHBhcmFtIHtzdHJpbmd9IG9wdGlvbnMudGl0bGUgLSBOb3RpZmljYXRpb24gdGl0bGVcbiAqIEBwYXJhbSB7c3RyaW5nfSBbb3B0aW9ucy5zdWJ0aXRsZV0gLSBOb3RpZmljYXRpb24gc3VidGl0bGUgKG1hY09TIGFuZCBMaW51eCBvbmx5KVxuICogQHBhcmFtIHtzdHJpbmd9IFtvcHRpb25zLmJvZHldIC0gTm90aWZpY2F0aW9uIGJvZHkgdGV4dFxuICogQHBhcmFtIHtzdHJpbmd9IG9wdGlvbnMuY2F0ZWdvcnlJZCAtIENhdGVnb3J5IElEIHRoYXQgbWF0Y2hlcyBhIHJlZ2lzdGVyZWQgY2F0ZWdvcnlcbiAqIEBwYXJhbSB7T2JqZWN0PHN0cmluZywgYW55Pn0gW29wdGlvbnMuZGF0YV0gLSBBZGRpdGlvbmFsIHVzZXIgZGF0YSB0byBhdHRhY2ggdG8gdGhlIG5vdGlmaWNhdGlvblxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNlbmROb3RpZmljYXRpb25XaXRoQWN0aW9ucyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6U2VuZE5vdGlmaWNhdGlvbldpdGhBY3Rpb25zXCIsIFtvcHRpb25zXSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBub3RpZmljYXRpb24gY2F0ZWdvcnkgdGhhdCBjYW4gYmUgdXNlZCB3aXRoIFNlbmROb3RpZmljYXRpb25XaXRoQWN0aW9ucy5cbiAqIENhdGVnb3JpZXMgZGVmaW5lIHRoZSBhY3Rpb24gYnV0dG9ucyBhbmQgb3B0aW9uYWwgcmVwbHkgZmllbGRzIHRoYXQgd2lsbCBhcHBlYXIgb24gbm90aWZpY2F0aW9ucy5cbiAqIFJlZ2lzdGVyaW5nIGEgY2F0ZWdvcnkgd2l0aCB0aGUgc2FtZSBJRCBhcyBhIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCBjYXRlZ29yeSB3aWxsIG92ZXJyaWRlIGl0LlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7T2JqZWN0fSBjYXRlZ29yeSAtIE5vdGlmaWNhdGlvbiBjYXRlZ29yeSBkZWZpbml0aW9uXG4gKiBAcGFyYW0ge3N0cmluZ30gY2F0ZWdvcnkuaWQgLSBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIGNhdGVnb3J5XG4gKiBAcGFyYW0ge0FycmF5PE9iamVjdD59IFtjYXRlZ29yeS5hY3Rpb25zXSAtIEFycmF5IG9mIGFjdGlvbiBidXR0b25zXG4gKiBAcGFyYW0ge3N0cmluZ30gY2F0ZWdvcnkuYWN0aW9uc1tdLmlkIC0gVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBhY3Rpb25cbiAqIEBwYXJhbSB7c3RyaW5nfSBjYXRlZ29yeS5hY3Rpb25zW10udGl0bGUgLSBEaXNwbGF5IHRpdGxlIGZvciB0aGUgYWN0aW9uIGJ1dHRvblxuICogQHBhcmFtIHtib29sZWFufSBbY2F0ZWdvcnkuYWN0aW9uc1tdLmRlc3RydWN0aXZlXSAtIFdoZXRoZXIgdGhlIGFjdGlvbiBpcyBkZXN0cnVjdGl2ZSAobWFjT1Mtc3BlY2lmaWMpXG4gKiBAcGFyYW0ge2Jvb2xlYW59IFtjYXRlZ29yeS5oYXNSZXBseUZpZWxkXSAtIFdoZXRoZXIgdG8gaW5jbHVkZSBhIHRleHQgaW5wdXQgZmllbGQgZm9yIHJlcGxpZXNcbiAqIEBwYXJhbSB7c3RyaW5nfSBbY2F0ZWdvcnkucmVwbHlQbGFjZWhvbGRlcl0gLSBQbGFjZWhvbGRlciB0ZXh0IGZvciB0aGUgcmVwbHkgZmllbGQgKHJlcXVpcmVkIGlmIGhhc1JlcGx5RmllbGQgaXMgdHJ1ZSlcbiAqIEBwYXJhbSB7c3RyaW5nfSBbY2F0ZWdvcnkucmVwbHlCdXR0b25UaXRsZV0gLSBUaXRsZSBmb3IgdGhlIHJlcGx5IGJ1dHRvbiAocmVxdWlyZWQgaWYgaGFzUmVwbHlGaWVsZCBpcyB0cnVlKVxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlZ2lzdGVyTm90aWZpY2F0aW9uQ2F0ZWdvcnkoY2F0ZWdvcnkpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZWdpc3Rlck5vdGlmaWNhdGlvbkNhdGVnb3J5XCIsIFtjYXRlZ29yeV0pO1xufVxuXG4vKipcbiAqIFJlbW92ZSBhIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCBub3RpZmljYXRpb24gY2F0ZWdvcnkuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGNhdGVnb3J5SWQgLSBUaGUgSUQgb2YgdGhlIGNhdGVnb3J5IHRvIHJlbW92ZVxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbW92ZU5vdGlmaWNhdGlvbkNhdGVnb3J5KGNhdGVnb3J5SWQpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZW1vdmVOb3RpZmljYXRpb25DYXRlZ29yeVwiLCBbY2F0ZWdvcnlJZF0pO1xufVxuXG4vKipcbiAqIFJlbW92ZSBhbGwgcGVuZGluZyBub3RpZmljYXRpb25zIGZyb20gdGhlIG5vdGlmaWNhdGlvbiBjZW50ZXIuXG4gKiBPbiBXaW5kb3dzLCB0aGlzIGlzIGEgbm8tb3AgYXMgdGhlIHBsYXRmb3JtIG1hbmFnZXMgbm90aWZpY2F0aW9uIGxpZmVjeWNsZSBhdXRvbWF0aWNhbGx5LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZW1vdmVBbGxQZW5kaW5nTm90aWZpY2F0aW9ucygpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZW1vdmVBbGxQZW5kaW5nTm90aWZpY2F0aW9uc1wiKTtcbn1cblxuLyoqXG4gKiBSZW1vdmUgYSBzcGVjaWZpYyBwZW5kaW5nIG5vdGlmaWNhdGlvbiBieSBpdHMgaWRlbnRpZmllci5cbiAqIE9uIFdpbmRvd3MsIHRoaXMgaXMgYSBuby1vcCBhcyB0aGUgcGxhdGZvcm0gbWFuYWdlcyBub3RpZmljYXRpb24gbGlmZWN5Y2xlIGF1dG9tYXRpY2FsbHkuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGlkZW50aWZpZXIgLSBUaGUgSUQgb2YgdGhlIG5vdGlmaWNhdGlvbiB0byByZW1vdmVcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZW1vdmVQZW5kaW5nTm90aWZpY2F0aW9uKGlkZW50aWZpZXIpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZW1vdmVQZW5kaW5nTm90aWZpY2F0aW9uXCIsIFtpZGVudGlmaWVyXSk7XG59XG5cbi8qKlxuICogUmVtb3ZlIGFsbCBkZWxpdmVyZWQgbm90aWZpY2F0aW9ucyBmcm9tIHRoZSBub3RpZmljYXRpb24gY2VudGVyLlxuICogT24gV2luZG93cywgdGhpcyBpcyBhIG5vLW9wIGFzIHRoZSBwbGF0Zm9ybSBtYW5hZ2VzIG5vdGlmaWNhdGlvbiBsaWZlY3ljbGUgYXV0b21hdGljYWxseS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVtb3ZlQWxsRGVsaXZlcmVkTm90aWZpY2F0aW9ucygpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZW1vdmVBbGxEZWxpdmVyZWROb3RpZmljYXRpb25zXCIpO1xufVxuXG4vKipcbiAqIFJlbW92ZSBhIHNwZWNpZmljIGRlbGl2ZXJlZCBub3RpZmljYXRpb24gYnkgaXRzIGlkZW50aWZpZXIuXG4gKiBPbiBXaW5kb3dzLCB0aGlzIGlzIGEgbm8tb3AgYXMgdGhlIHBsYXRmb3JtIG1hbmFnZXMgbm90aWZpY2F0aW9uIGxpZmVjeWNsZSBhdXRvbWF0aWNhbGx5LlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZGVudGlmaWVyIC0gVGhlIElEIG9mIHRoZSBub3RpZmljYXRpb24gdG8gcmVtb3ZlXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVtb3ZlRGVsaXZlcmVkTm90aWZpY2F0aW9uKGlkZW50aWZpZXIpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZW1vdmVEZWxpdmVyZWROb3RpZmljYXRpb25cIiwgW2lkZW50aWZpZXJdKTtcbn1cblxuLyoqXG4gKiBSZW1vdmUgYSBub3RpZmljYXRpb24gYnkgaXRzIGlkZW50aWZpZXIuXG4gKiBUaGlzIGlzIGEgY29udmVuaWVuY2UgZnVuY3Rpb24gdGhhdCB3b3JrcyBhY3Jvc3MgcGxhdGZvcm1zLlxuICogT24gbWFjT1MsIHVzZSB0aGUgbW9yZSBzcGVjaWZpYyBSZW1vdmVQZW5kaW5nTm90aWZpY2F0aW9uIG9yIFJlbW92ZURlbGl2ZXJlZE5vdGlmaWNhdGlvbiBmdW5jdGlvbnMuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGlkZW50aWZpZXIgLSBUaGUgSUQgb2YgdGhlIG5vdGlmaWNhdGlvbiB0byByZW1vdmVcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZW1vdmVOb3RpZmljYXRpb24oaWRlbnRpZmllcikge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlbW92ZU5vdGlmaWNhdGlvblwiLCBbaWRlbnRpZmllcl0pO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xuaW1wb3J0IHtcbiAgZXZlbnRMaXN0ZW5lcnMsXG4gIEV2ZW50c0VtaXQsXG4gIEV2ZW50c05vdGlmeSxcbiAgRXZlbnRzT2ZmLFxuICBFdmVudHNPZmZBbGwsXG4gIEV2ZW50c09uLFxuICBFdmVudHNPbmNlLFxuICBFdmVudHNPbk11bHRpcGxlLFxufSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7IENhbGwsIENhbGxiYWNrLCBjYWxsYmFja3MgfSBmcm9tICcuL2NhbGxzJztcbmltcG9ydCB7IFNldEJpbmRpbmdzIH0gZnJvbSBcIi4vYmluZGluZ3NcIjtcbmltcG9ydCAqIGFzIFdpbmRvdyBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCAqIGFzIFNjcmVlbiBmcm9tIFwiLi9zY3JlZW5cIjtcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3NlclwiO1xuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZFwiO1xuaW1wb3J0ICogYXMgRHJhZ0FuZERyb3AgZnJvbSBcIi4vZHJhZ2FuZGRyb3BcIjtcbmltcG9ydCAqIGFzIENvbnRleHRNZW51IGZyb20gXCIuL2NvbnRleHRtZW51XCI7XG5pbXBvcnQgKiBhcyBOb3RpZmljYXRpb25zIGZyb20gXCIuL25vdGlmaWNhdGlvbnNcIjtcblxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdRJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUycpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0gnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkVudmlyb25tZW50XCIpO1xufVxuXG4vLyBUaGUgSlMgcnVudGltZVxud2luZG93LnJ1bnRpbWUgPSB7XG4gICAgLi4uTG9nLFxuICAgIC4uLldpbmRvdyxcbiAgICAuLi5Ccm93c2VyLFxuICAgIC4uLlNjcmVlbixcbiAgICAuLi5DbGlwYm9hcmQsXG4gICAgLi4uRHJhZ0FuZERyb3AsXG4gICAgLi4uTm90aWZpY2F0aW9ucyxcbiAgICBFdmVudHNPbixcbiAgICBFdmVudHNPbmNlLFxuICAgIEV2ZW50c09uTXVsdGlwbGUsXG4gICAgRXZlbnRzRW1pdCxcbiAgICBFdmVudHNPZmYsXG4gICAgRXZlbnRzT2ZmQWxsLFxuICAgIEVudmlyb25tZW50LFxuICAgIFNob3csXG4gICAgSGlkZSxcbiAgICBRdWl0XG59O1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy53YWlscyA9IHtcbiAgICBDYWxsYmFjayxcbiAgICBFdmVudHNOb3RpZnksXG4gICAgU2V0QmluZGluZ3MsXG4gICAgZXZlbnRMaXN0ZW5lcnMsXG4gICAgY2FsbGJhY2tzLFxuICAgIGZsYWdzOiB7XG4gICAgICAgIGRpc2FibGVTY3JvbGxiYXJEcmFnOiBmYWxzZSxcbiAgICAgICAgZGlzYWJsZURlZmF1bHRDb250ZXh0TWVudTogZmFsc2UsXG4gICAgICAgIGVuYWJsZVJlc2l6ZTogZmFsc2UsXG4gICAgICAgIGRlZmF1bHRDdXJzb3I6IG51bGwsXG4gICAgICAgIGJvcmRlclRoaWNrbmVzczogNixcbiAgICAgICAgc2hvdWxkRHJhZzogZmFsc2UsXG4gICAgICAgIGRlZmVyRHJhZ1RvTW91c2VNb3ZlOiB0cnVlLFxuICAgICAgICBjc3NEcmFnUHJvcGVydHk6IFwiLS13YWlscy1kcmFnZ2FibGVcIixcbiAgICAgICAgY3NzRHJhZ1ZhbHVlOiBcImRyYWdcIixcbiAgICAgICAgY3NzRHJvcFByb3BlcnR5OiBcIi0td2FpbHMtZHJvcC10YXJnZXRcIixcbiAgICAgICAgY3NzRHJvcFZhbHVlOiBcImRyb3BcIixcbiAgICAgICAgZW5hYmxlV2FpbHNEcmFnQW5kRHJvcDogZmFsc2UsXG4gICAgfVxufTtcblxuLy8gU2V0IHRoZSBiaW5kaW5nc1xuaWYgKHdpbmRvdy53YWlsc2JpbmRpbmdzKSB7XG4gICAgd2luZG93LndhaWxzLlNldEJpbmRpbmdzKHdpbmRvdy53YWlsc2JpbmRpbmdzKTtcbiAgICBkZWxldGUgd2luZG93LndhaWxzLlNldEJpbmRpbmdzO1xufVxuXG4vLyAoYm9vbCkgVGhpcyBpcyBldmFsdWF0ZWQgYXQgYnVpbGQgdGltZSBpbiBwYWNrYWdlLmpzb25cbmlmICghREVCVUcpIHtcbiAgICBkZWxldGUgd2luZG93LndhaWxzYmluZGluZ3M7XG59XG5cbmxldCBkcmFnVGVzdCA9IGZ1bmN0aW9uKGUpIHtcbiAgICB2YXIgdmFsID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpLmdldFByb3BlcnR5VmFsdWUod2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdQcm9wZXJ0eSk7XG4gICAgaWYgKHZhbCkge1xuICAgICAgICB2YWwgPSB2YWwudHJpbSgpO1xuICAgIH1cblxuICAgIGlmICh2YWwgIT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIGlmIChlLmJ1dHRvbnMgIT09IDEpIHtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWdnaW5nIGlmIG5vdCB0aGUgcHJpbWFyeSBidXR0b24gaGFzIGJlZW4gY2xpY2tlZC5cbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIGlmIChlLmRldGFpbCAhPT0gMSkge1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZ2dpbmcgaWYgbW9yZSB0aGFuIG9uY2UgaGFzIGJlZW4gY2xpY2tlZCwgZS5nLiB3aGVuIGRvdWJsZSBjbGlja2luZ1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgcmV0dXJuIHRydWU7XG59O1xuXG53aW5kb3cud2FpbHMuc2V0Q1NTRHJhZ1Byb3BlcnRpZXMgPSBmdW5jdGlvbihwcm9wZXJ0eSwgdmFsdWUpIHtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1Byb3BlcnR5ID0gcHJvcGVydHk7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdWYWx1ZSA9IHZhbHVlO1xufVxuXG53aW5kb3cud2FpbHMuc2V0Q1NTRHJvcFByb3BlcnRpZXMgPSBmdW5jdGlvbihwcm9wZXJ0eSwgdmFsdWUpIHtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJvcFByb3BlcnR5ID0gcHJvcGVydHk7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BWYWx1ZSA9IHZhbHVlO1xufVxuXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgKGUpID0+IHtcbiAgICAvLyBDaGVjayBmb3IgcmVzaXppbmdcbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwicmVzaXplOlwiICsgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UpO1xuICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoZHJhZ1Rlc3QoZSkpIHtcbiAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlU2Nyb2xsYmFyRHJhZykge1xuICAgICAgICAgICAgLy8gVGhpcyBjaGVja3MgZm9yIGNsaWNrcyBvbiB0aGUgc2Nyb2xsIGJhclxuICAgICAgICAgICAgaWYgKGUub2Zmc2V0WCA+IGUudGFyZ2V0LmNsaWVudFdpZHRoIHx8IGUub2Zmc2V0WSA+IGUudGFyZ2V0LmNsaWVudEhlaWdodCkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRlZmVyRHJhZ1RvTW91c2VNb3ZlKSB7XG4gICAgICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IHRydWU7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBlLnByZXZlbnREZWZhdWx0KClcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcImRyYWdcIik7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufSk7XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgKCkgPT4ge1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG59KTtcblxuZnVuY3Rpb24gc2V0UmVzaXplKGN1cnNvcikge1xuICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBjdXJzb3IgfHwgd2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3I7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBmdW5jdGlvbihlKSB7XG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnKSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5zaG91bGREcmFnID0gZmFsc2U7XG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJkcmFnXCIpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmICghd2luZG93LndhaWxzLmZsYWdzLmVuYWJsZVJlc2l6ZSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XG4gICAgICAgIHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yID0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvcjtcbiAgICB9XG4gICAgaWYgKHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcyAmJiB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzKSB7XG4gICAgICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBcInNlLXJlc2l6ZVwiO1xuICAgIH1cbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3M7XG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiB3aW5kb3cud2FpbHMuZmxhZ3MucmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH0gZWxzZSBpZiAocmlnaHRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiBib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyICYmIHRvcEJvcmRlcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlciAmJiByaWdodEJvcmRlcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XG5cbn0pO1xuXG4vLyBTZXR1cCBjb250ZXh0IG1lbnUgaG9va1xud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgZnVuY3Rpb24oZSkge1xuICAgIC8vIGFsd2F5cyBzaG93IHRoZSBjb250ZXh0bWVudSBpbiBkZWJ1ZyAmIGRldlxuICAgIGlmIChERUJVRykgcmV0dXJuO1xuXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kaXNhYmxlRGVmYXVsdENvbnRleHRNZW51KSB7XG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBDb250ZXh0TWVudS5wcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGUpO1xuICAgIH1cbn0pO1xuXG53aW5kb3cuV2FpbHNJbnZva2UoXCJydW50aW1lOnJlYWR5XCIpOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBa0JBLFdBQVMsZUFBZSxPQUFPLFNBQVM7QUFJdkMsV0FBTyxZQUFZLE1BQU0sUUFBUSxPQUFPO0FBQUEsRUFDekM7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFFBQVEsU0FBUztBQUNoQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsV0FBVyxTQUFTO0FBQ25DLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsWUFBWSxVQUFVO0FBQ3JDLG1CQUFlLEtB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T0FBTyxNQUFNLE1BQU0sWUFBWTtBQUMvQixhQUFPLE1BQU0sTUFBTSxhQUFhO0FBQ2hDLFVBQUksZUFBZSxFQUFFLFlBQVksU0FBWSxFQUFFLFVBQVUsRUFBRTtBQUMzRCxVQUFJLGVBQWUsR0FBRztBQUNsQixlQUFPLFlBQVksTUFBTTtBQUN6QjtBQUFBLE1BQ0o7QUFBQSxJQUNKO0FBQ0EsUUFBSSxDQUFDLE9BQU8sTUFBTSxNQUFNLGNBQWM7QUFDbEM7QUFBQSxJQUNKO0FBQ0EsUUFBSSxPQUFPLE1BQU0sTUFBTSxpQkFBaUIsTUFBTTtBQUMxQyxhQUFPLE1BQU0sTUFBTSxnQkFBZ0IsU0FBUyxnQkFBZ0IsTUFBTTtBQUFBLElBQ3RFO0FBQ0EsUUFBSSxPQUFPLGFBQWEsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNLG1CQUFtQixPQUFPLGNBQWMsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNLGlCQUFpQjtBQUMzSSxlQUFTLGdCQUFnQixNQUFNLFNBQVM7QUFBQSxJQUM1QztBQUNBLFFBQUksY0FBYyxPQUFPLGFBQWEsRUFBRSxVQUFVLE9BQU8sTUFBTSxNQUFNO0FBQ3JFLFFBQUksYUFBYSxFQUFFLFVBQVUsT0FBTyxNQUFNLE1BQU07QUFDaEQsUUFBSSxZQUFZLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUMvQyxRQUFJLGVBQWUsT0FBTyxjQUFjLEVBQUUsVUFBVSxPQUFPLE1BQU0sTUFBTTtBQUd2RSxRQUFJLENBQUMsY0FBYyxDQUFDLGVBQWUsQ0FBQyxhQUFhLENBQUMsZ0JBQWdCLE9BQU8sTUFBTSxNQUFNLGVBQWUsUUFBVztBQUMzRyxnQkFBVTtBQUFBLElBQ2QsV0FBVyxlQUFlO0FBQWMsZ0JBQVUsV0FBVztBQUFBLGFBQ3BELGNBQWM7QUFBYyxnQkFBVSxXQUFXO0FBQUEsYUFDakQsY0FBYztBQUFXLGdCQUFVLFdBQVc7QUFBQSxhQUM5QyxhQUFhO0FBQWEsZ0JBQVUsV0FBVztBQUFBLGFBQy9DO0FBQVksZ0JBQVUsVUFBVTtBQUFBLGFBQ2hDO0FBQVcsZ0JBQVUsVUFBVTtBQUFBLGFBQy9CO0FBQWMsZ0JBQVUsVUFBVTtBQUFBLGFBQ2xDO0FBQWEsZ0JBQVUsVUFBVTtBQUFBLEVBRTlDLENBQUM7QUFHRCxTQUFPLGlCQUFpQixlQUFlLFNBQVMsR0FBRztBQUUvQyxRQUFJO0FBQU87QUFFWCxRQUFJLE9BQU8sTUFBTSxNQUFNLDJCQUEyQjtBQUM5QyxRQUFFLGVBQWU7QUFBQSxJQUNyQixPQUFPO0FBQ0gsTUFBWSwwQkFBMEIsQ0FBQztBQUFBLElBQzNDO0FBQUEsRUFDSixDQUFDO0FBRUQsU0FBTyxZQUFZLGVBQWU7IiwKICAibmFtZXMiOiBbImV2ZW50TmFtZSJdCn0K
