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
    if (eventListeners[eventName] === void 0) return;
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
          window.go[packageName][structName][methodName] = (function() {
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
          })();
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
    if (flags.nextDeactivate) flags.nextDeactivate();
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
    if (flags.nextDeactivate) flags.nextDeactivate();
    flags.nextDeactivate = () => {
      Array.from(document.getElementsByClassName(DROP_TARGET_ACTIVE)).forEach((el) => el.classList.remove(DROP_TARGET_ACTIVE));
      flags.nextDeactivate = null;
      if (flags.nextDeactivateTimeout) {
        clearTimeout(flags.nextDeactivateTimeout);
        flags.nextDeactivateTimeout = null;
      }
    };
    flags.nextDeactivateTimeout = setTimeout(() => {
      if (flags.nextDeactivate) flags.nextDeactivate();
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
    if (flags.nextDeactivate) flags.nextDeactivate();
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
      default: {
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
    } else if (rightBorder && bottomBorder) setResize("se-resize");
    else if (leftBorder && bottomBorder) setResize("sw-resize");
    else if (leftBorder && topBorder) setResize("nw-resize");
    else if (topBorder && rightBorder) setResize("ne-resize");
    else if (leftBorder) setResize("w-resize");
    else if (topBorder) setResize("n-resize");
    else if (bottomBorder) setResize("s-resize");
    else if (rightBorder) setResize("e-resize");
  });
  window.addEventListener("contextmenu", function(e) {
    if (true) return;
    if (window.wails.flags.disableDefaultContextMenu) {
      e.preventDefault();
    } else {
      processDefaultContextMenu(e);
    }
  });
  window.WailsInvoke("runtime:ready");
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3Avc2NyZWVuLmpzIiwgImRlc2t0b3AvYnJvd3Nlci5qcyIsICJkZXNrdG9wL2NsaXBib2FyZC5qcyIsICJkZXNrdG9wL2RyYWdhbmRkcm9wLmpzIiwgImRlc2t0b3AvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9ub3RpZmljYXRpb25zLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfICAgICAgIF9fICAgICAgXyBfX1xufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuLyoqXG4gKiBTZW5kcyBhIGxvZyBtZXNzYWdlIHRvIHRoZSBiYWNrZW5kIHdpdGggdGhlIGdpdmVuIGxldmVsICsgbWVzc2FnZVxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBsZXZlbFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZnVuY3Rpb24gc2VuZExvZ01lc3NhZ2UobGV2ZWwsIG1lc3NhZ2UpIHtcblxuXHQvLyBMb2cgTWVzc2FnZSBmb3JtYXQ6XG5cdC8vIGxbdHlwZV1bbWVzc2FnZV1cblx0d2luZG93LldhaWxzSW52b2tlKCdMJyArIGxldmVsICsgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiB0cmFjZSBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gTG9nVHJhY2UobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnVCcsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1ByaW50KG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ1AnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGRlYnVnIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dEZWJ1ZyhtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdEJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogTG9nIHRoZSBnaXZlbiBpbmZvIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dJbmZvKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0knLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIHdhcm5pbmcgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ1dhcm5pbmcobWVzc2FnZSkge1xuXHRzZW5kTG9nTWVzc2FnZSgnVycsIG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIExvZyB0aGUgZ2l2ZW4gZXJyb3IgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZ0Vycm9yKG1lc3NhZ2UpIHtcblx0c2VuZExvZ01lc3NhZ2UoJ0UnLCBtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBMb2cgdGhlIGdpdmVuIGZhdGFsIG1lc3NhZ2Ugd2l0aCB0aGUgYmFja2VuZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBMb2dGYXRhbChtZXNzYWdlKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdGJywgbWVzc2FnZSk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgTG9nIGxldmVsIHRvIHRoZSBnaXZlbiBsb2cgbGV2ZWxcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge251bWJlcn0gbG9nbGV2ZWxcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldExvZ0xldmVsKGxvZ2xldmVsKSB7XG5cdHNlbmRMb2dNZXNzYWdlKCdTJywgbG9nbGV2ZWwpO1xufVxuXG4vLyBMb2cgbGV2ZWxzXG5leHBvcnQgY29uc3QgTG9nTGV2ZWwgPSB7XG5cdFRSQUNFOiAxLFxuXHRERUJVRzogMixcblx0SU5GTzogMyxcblx0V0FSTklORzogNCxcblx0RVJST1I6IDUsXG59O1xuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XG5cbiAgICAvLyBHZXQgdGhlIGV2ZW50IG5hbWVcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XG5cbiAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lciBpbmRleGVzIHRvIGRlc3Ryb3lcbiAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0/LnNsaWNlKCkgfHwgW107XG5cbiAgICAvLyBDaGVjayBpZiB3ZSBoYXZlIGFueSBsaXN0ZW5lcnMgZm9yIHRoaXMgZXZlbnRcbiAgICBpZiAobmV3RXZlbnRMaXN0ZW5lckxpc3QubGVuZ3RoKSB7XG5cbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcbiAgICAgICAgZm9yIChsZXQgY291bnQgPSBuZXdFdmVudExpc3RlbmVyTGlzdC5sZW5ndGggLSAxOyBjb3VudCA+PSAwOyBjb3VudCAtPSAxKSB7XG5cbiAgICAgICAgICAgIC8vIEdldCBuZXh0IGxpc3RlbmVyXG4gICAgICAgICAgICBjb25zdCBsaXN0ZW5lciA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0W2NvdW50XTtcblxuICAgICAgICAgICAgbGV0IGRhdGEgPSBldmVudERhdGEuZGF0YTtcblxuICAgICAgICAgICAgLy8gRG8gdGhlIGNhbGxiYWNrXG4gICAgICAgICAgICBjb25zdCBkZXN0cm95ID0gbGlzdGVuZXIuQ2FsbGJhY2soZGF0YSk7XG4gICAgICAgICAgICBpZiAoZGVzdHJveSkge1xuICAgICAgICAgICAgICAgIC8vIGlmIHRoZSBsaXN0ZW5lciBpbmRpY2F0ZWQgdG8gZGVzdHJveSBpdHNlbGYsIGFkZCBpdCB0byB0aGUgZGVzdHJveSBsaXN0XG4gICAgICAgICAgICAgICAgbmV3RXZlbnRMaXN0ZW5lckxpc3Quc3BsaWNlKGNvdW50LCAxKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFVwZGF0ZSBjYWxsYmFja3Mgd2l0aCBuZXcgbGlzdCBvZiBsaXN0ZW5lcnNcbiAgICAgICAgaWYgKG5ld0V2ZW50TGlzdGVuZXJMaXN0Lmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBuZXdFdmVudExpc3RlbmVyTGlzdDtcbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBOb3RpZnkgaW5mb3JtcyBmcm9udGVuZCBsaXN0ZW5lcnMgdGhhdCBhbiBldmVudCB3YXMgZW1pdHRlZCB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5vdGlmeU1lc3NhZ2UgLSBlbmNvZGVkIG5vdGlmaWNhdGlvbiBtZXNzYWdlXG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c05vdGlmeShub3RpZnlNZXNzYWdlKSB7XG4gICAgLy8gUGFyc2UgdGhlIG1lc3NhZ2VcbiAgICBsZXQgbWVzc2FnZTtcbiAgICB0cnkge1xuICAgICAgICBtZXNzYWdlID0gSlNPTi5wYXJzZShub3RpZnlNZXNzYWdlKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGNvbnN0IGVycm9yID0gJ0ludmFsaWQgSlNPTiBwYXNzZWQgdG8gTm90aWZ5OiAnICsgbm90aWZ5TWVzc2FnZTtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGVycm9yKTtcbiAgICB9XG4gICAgbm90aWZ5TGlzdGVuZXJzKG1lc3NhZ2UpO1xufVxuXG4vKipcbiAqIEVtaXQgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c0VtaXQoZXZlbnROYW1lKSB7XG5cbiAgICBjb25zdCBwYXlsb2FkID0ge1xuICAgICAgICBuYW1lOiBldmVudE5hbWUsXG4gICAgICAgIGRhdGE6IFtdLnNsaWNlLmFwcGx5KGFyZ3VtZW50cykuc2xpY2UoMSksXG4gICAgfTtcblxuICAgIC8vIE5vdGlmeSBKUyBsaXN0ZW5lcnNcbiAgICBub3RpZnlMaXN0ZW5lcnMocGF5bG9hZCk7XG5cbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFRScgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG59XG5cbmZ1bmN0aW9uIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSkge1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lcnNcbiAgICBkZWxldGUgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXTtcblxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VYJyArIGV2ZW50TmFtZSk7XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVyZXMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKVxuXG4gICAgaWYgKGFkZGl0aW9uYWxFdmVudE5hbWVzLmxlbmd0aCA+IDApIHtcbiAgICAgICAgYWRkaXRpb25hbEV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKVxuICAgICAgICB9KVxuICAgIH1cbn1cblxuLyoqXG4gKiBPZmYgdW5yZWdpc3RlcnMgYWxsIGV2ZW50IGxpc3RlbmVycyBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICovXG4gZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZkFsbCgpIHtcbiAgICBjb25zdCBldmVudE5hbWVzID0gT2JqZWN0LmtleXMoZXZlbnRMaXN0ZW5lcnMpO1xuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpXG4gICAgfSlcbn1cblxuLyoqXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIEV2ZW50c09uXG4gKlxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcbiAqL1xuIGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdID09PSB1bmRlZmluZWQpIHJldHVybjtcblxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gPSBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcblxuICAgIC8vIENsZWFuIHVwIGlmIHRoZXJlIGFyZSBubyBldmVudCBsaXN0ZW5lcnMgbGVmdFxuICAgIGlmIChldmVudExpc3RlbmVyc1tldmVudE5hbWVdLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpO1xuICAgIH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuZXhwb3J0IGNvbnN0IGNhbGxiYWNrcyA9IHt9O1xuXG4vKipcbiAqIFJldHVybnMgYSBudW1iZXIgZnJvbSB0aGUgbmF0aXZlIGJyb3dzZXIgcmFuZG9tIGZ1bmN0aW9uXG4gKlxuICogQHJldHVybnMgbnVtYmVyXG4gKi9cbmZ1bmN0aW9uIGNyeXB0b1JhbmRvbSgpIHtcblx0dmFyIGFycmF5ID0gbmV3IFVpbnQzMkFycmF5KDEpO1xuXHRyZXR1cm4gd2luZG93LmNyeXB0by5nZXRSYW5kb21WYWx1ZXMoYXJyYXkpWzBdO1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBudW1iZXIgdXNpbmcgZGEgb2xkLXNrb29sIE1hdGguUmFuZG9tXG4gKiBJIGxpa2VzIHRvIGNhbGwgaXQgTE9MUmFuZG9tXG4gKlxuICogQHJldHVybnMgbnVtYmVyXG4gKi9cbmZ1bmN0aW9uIGJhc2ljUmFuZG9tKCkge1xuXHRyZXR1cm4gTWF0aC5yYW5kb20oKSAqIDkwMDcxOTkyNTQ3NDA5OTE7XG59XG5cbi8vIFBpY2sgYSByYW5kb20gbnVtYmVyIGZ1bmN0aW9uIGJhc2VkIG9uIGJyb3dzZXIgY2FwYWJpbGl0eVxudmFyIHJhbmRvbUZ1bmM7XG5pZiAod2luZG93LmNyeXB0bykge1xuXHRyYW5kb21GdW5jID0gY3J5cHRvUmFuZG9tO1xufSBlbHNlIHtcblx0cmFuZG9tRnVuYyA9IGJhc2ljUmFuZG9tO1xufVxuXG5cbi8qKlxuICogQ2FsbCBzZW5kcyBhIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgdG8gY2FsbCB0aGUgYmluZGluZyB3aXRoIHRoZVxuICogZ2l2ZW4gZGF0YS4gQSBwcm9taXNlIGlzIHJldHVybmVkIGFuZCB3aWxsIGJlIGNvbXBsZXRlZCB3aGVuIHRoZVxuICogYmFja2VuZCByZXNwb25kcy4gVGhpcyB3aWxsIGJlIHJlc29sdmVkIHdoZW4gdGhlIGNhbGwgd2FzIHN1Y2Nlc3NmdWxcbiAqIG9yIHJlamVjdGVkIGlmIGFuIGVycm9yIGlzIHBhc3NlZCBiYWNrLlxuICogVGhlcmUgaXMgYSB0aW1lb3V0IG1lY2hhbmlzbS4gSWYgdGhlIGNhbGwgZG9lc24ndCByZXNwb25kIGluIHRoZSBnaXZlblxuICogdGltZSAoaW4gbWlsbGlzZWNvbmRzKSB0aGVuIHRoZSBwcm9taXNlIGlzIHJlamVjdGVkLlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lXG4gKiBAcGFyYW0ge2FueT19IGFyZ3NcbiAqIEBwYXJhbSB7bnVtYmVyPX0gdGltZW91dFxuICogQHJldHVybnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGwobmFtZSwgYXJncywgdGltZW91dCkge1xuXG5cdC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxuXHRpZiAodGltZW91dCA9PSBudWxsKSB7XG5cdFx0dGltZW91dCA9IDA7XG5cdH1cblxuXHQvLyBDcmVhdGUgYSBwcm9taXNlXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cblx0XHQvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuXHRcdHZhciBjYWxsYmFja0lEO1xuXHRcdGRvIHtcblx0XHRcdGNhbGxiYWNrSUQgPSBuYW1lICsgJy0nICsgcmFuZG9tRnVuYygpO1xuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cblx0XHR2YXIgdGltZW91dEhhbmRsZTtcblx0XHQvLyBTZXQgdGltZW91dFxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xuXHRcdFx0dGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcblx0XHRcdH0sIHRpbWVvdXQpO1xuXHRcdH1cblxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXG5cdFx0Y2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcblx0XHRcdHJlamVjdDogcmVqZWN0LFxuXHRcdFx0cmVzb2x2ZTogcmVzb2x2ZVxuXHRcdH07XG5cblx0XHR0cnkge1xuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0bmFtZSxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdH07XG5cbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnQycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcbiAgICAgICAgfVxuICAgIH0pO1xufVxuXG53aW5kb3cuT2JmdXNjYXRlZENhbGwgPSAoaWQsIGFyZ3MsIHRpbWVvdXQpID0+IHtcblxuICAgIC8vIFRpbWVvdXQgaW5maW5pdGUgYnkgZGVmYXVsdFxuICAgIGlmICh0aW1lb3V0ID09IG51bGwpIHtcbiAgICAgICAgdGltZW91dCA9IDA7XG4gICAgfVxuXG4gICAgLy8gQ3JlYXRlIGEgcHJvbWlzZVxuICAgIHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cbiAgICAgICAgLy8gQ3JlYXRlIGEgdW5pcXVlIGNhbGxiYWNrSURcbiAgICAgICAgdmFyIGNhbGxiYWNrSUQ7XG4gICAgICAgIGRvIHtcbiAgICAgICAgICAgIGNhbGxiYWNrSUQgPSBpZCArICctJyArIHJhbmRvbUZ1bmMoKTtcbiAgICAgICAgfSB3aGlsZSAoY2FsbGJhY2tzW2NhbGxiYWNrSURdKTtcblxuICAgICAgICB2YXIgdGltZW91dEhhbmRsZTtcbiAgICAgICAgLy8gU2V0IHRpbWVvdXRcbiAgICAgICAgaWYgKHRpbWVvdXQgPiAwKSB7XG4gICAgICAgICAgICB0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKCdDYWxsIHRvIG1ldGhvZCAnICsgaWQgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcbiAgICAgICAgICAgIH0sIHRpbWVvdXQpO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gU3RvcmUgY2FsbGJhY2tcbiAgICAgICAgY2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuICAgICAgICAgICAgdGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcbiAgICAgICAgICAgIHJlamVjdDogcmVqZWN0LFxuICAgICAgICAgICAgcmVzb2x2ZTogcmVzb2x2ZVxuICAgICAgICB9O1xuXG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xuXHRcdFx0XHRpZCxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdH07XG5cbiAgICAgICAgICAgIC8vIE1ha2UgdGhlIGNhbGxcbiAgICAgICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnYycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XG4gICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlKTtcbiAgICAgICAgfVxuICAgIH0pO1xufTtcblxuXG4vKipcbiAqIENhbGxlZCBieSB0aGUgYmFja2VuZCB0byByZXR1cm4gZGF0YSB0byBhIHByZXZpb3VzbHkgY2FsbGVkXG4gKiBiaW5kaW5nIGludm9jYXRpb25cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gaW5jb21pbmdNZXNzYWdlXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsYmFjayhpbmNvbWluZ01lc3NhZ2UpIHtcblx0Ly8gUGFyc2UgdGhlIG1lc3NhZ2Vcblx0bGV0IG1lc3NhZ2U7XG5cdHRyeSB7XG5cdFx0bWVzc2FnZSA9IEpTT04ucGFyc2UoaW5jb21pbmdNZXNzYWdlKTtcblx0fSBjYXRjaCAoZSkge1xuXHRcdGNvbnN0IGVycm9yID0gYEludmFsaWQgSlNPTiBwYXNzZWQgdG8gY2FsbGJhY2s6ICR7ZS5tZXNzYWdlfS4gTWVzc2FnZTogJHtpbmNvbWluZ01lc3NhZ2V9YDtcblx0XHRydW50aW1lLkxvZ0RlYnVnKGVycm9yKTtcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuXHR9XG5cdGxldCBjYWxsYmFja0lEID0gbWVzc2FnZS5jYWxsYmFja2lkO1xuXHRsZXQgY2FsbGJhY2tEYXRhID0gY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xuXHRpZiAoIWNhbGxiYWNrRGF0YSkge1xuXHRcdGNvbnN0IGVycm9yID0gYENhbGxiYWNrICcke2NhbGxiYWNrSUR9JyBub3QgcmVnaXN0ZXJlZCEhIWA7XG5cdFx0Y29uc29sZS5lcnJvcihlcnJvcik7IC8vIGVzbGludC1kaXNhYmxlLWxpbmVcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xuXHR9XG5cdGNsZWFyVGltZW91dChjYWxsYmFja0RhdGEudGltZW91dEhhbmRsZSk7XG5cblx0ZGVsZXRlIGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcblxuXHRpZiAobWVzc2FnZS5lcnJvcikge1xuXHRcdGNvbnN0IGVyciA9IG1lc3NhZ2UuZXJyb3IgaW5zdGFuY2VvZiBFcnJvciA/IG1lc3NhZ2UuZXJyb3IgOiBuZXcgRXJyb3IobWVzc2FnZS5lcnJvcik7XG5cdFx0Y2FsbGJhY2tEYXRhLnJlamVjdChlcnIpO1xuXHR9IGVsc2Uge1xuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX18gICAgXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gICkgXG58X18vfF9fL1xcX18sXy9fL18vX19fXy8gIFxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cblxuaW1wb3J0IHtDYWxsfSBmcm9tICcuL2NhbGxzJztcblxuLy8gVGhpcyBpcyB3aGVyZSB3ZSBiaW5kIGdvIG1ldGhvZCB3cmFwcGVyc1xud2luZG93LmdvID0ge307XG5cbmV4cG9ydCBmdW5jdGlvbiBTZXRCaW5kaW5ncyhiaW5kaW5nc01hcCkge1xuXHR0cnkge1xuXHRcdGJpbmRpbmdzTWFwID0gSlNPTi5wYXJzZShiaW5kaW5nc01hcCk7XG5cdH0gY2F0Y2ggKGUpIHtcblx0XHRjb25zb2xlLmVycm9yKGUpO1xuXHR9XG5cblx0Ly8gSW5pdGlhbGlzZSB0aGUgYmluZGluZ3MgbWFwXG5cdHdpbmRvdy5nbyA9IHdpbmRvdy5nbyB8fCB7fTtcblxuXHQvLyBJdGVyYXRlIHBhY2thZ2UgbmFtZXNcblx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXApLmZvckVhY2goKHBhY2thZ2VOYW1lKSA9PiB7XG5cblx0XHQvLyBDcmVhdGUgaW5uZXIgbWFwIGlmIGl0IGRvZXNuJ3QgZXhpc3Rcblx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXSB8fCB7fTtcblxuXHRcdC8vIEl0ZXJhdGUgc3RydWN0IG5hbWVzXG5cdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdKS5mb3JFYWNoKChzdHJ1Y3ROYW1lKSA9PiB7XG5cblx0XHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxuXHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSA9IHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV0gfHwge307XG5cblx0XHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSkuZm9yRWFjaCgobWV0aG9kTmFtZSkgPT4ge1xuXG5cdFx0XHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV1bc3RydWN0TmFtZV1bbWV0aG9kTmFtZV0gPSBmdW5jdGlvbiAoKSB7XG5cblx0XHRcdFx0XHQvLyBObyB0aW1lb3V0IGJ5IGRlZmF1bHRcblx0XHRcdFx0XHRsZXQgdGltZW91dCA9IDA7XG5cblx0XHRcdFx0XHQvLyBBY3R1YWwgZnVuY3Rpb25cblx0XHRcdFx0XHRmdW5jdGlvbiBkeW5hbWljKCkge1xuXHRcdFx0XHRcdFx0Y29uc3QgYXJncyA9IFtdLnNsaWNlLmNhbGwoYXJndW1lbnRzKTtcblx0XHRcdFx0XHRcdHJldHVybiBDYWxsKFtwYWNrYWdlTmFtZSwgc3RydWN0TmFtZSwgbWV0aG9kTmFtZV0uam9pbignLicpLCBhcmdzLCB0aW1lb3V0KTtcblx0XHRcdFx0XHR9XG5cblx0XHRcdFx0XHQvLyBBbGxvdyBzZXR0aW5nIHRpbWVvdXQgdG8gZnVuY3Rpb25cblx0XHRcdFx0XHRkeW5hbWljLnNldFRpbWVvdXQgPSBmdW5jdGlvbiAobmV3VGltZW91dCkge1xuXHRcdFx0XHRcdFx0dGltZW91dCA9IG5ld1RpbWVvdXQ7XG5cdFx0XHRcdFx0fTtcblxuXHRcdFx0XHRcdC8vIEFsbG93IGdldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxuXHRcdFx0XHRcdGR5bmFtaWMuZ2V0VGltZW91dCA9IGZ1bmN0aW9uICgpIHtcblx0XHRcdFx0XHRcdHJldHVybiB0aW1lb3V0O1xuXHRcdFx0XHRcdH07XG5cblx0XHRcdFx0XHRyZXR1cm4gZHluYW1pYztcblx0XHRcdFx0fSgpO1xuXHRcdFx0fSk7XG5cdFx0fSk7XG5cdH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93UmVsb2FkKCkge1xuICAgIHdpbmRvdy5sb2NhdGlvbi5yZWxvYWQoKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1JlbG9hZEFwcCgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dSJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTeXN0ZW1EZWZhdWx0VGhlbWUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVNEVCcpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TGlnaHRUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBTFQnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1NldERhcmtUaGVtZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dBRFQnKTtcbn1cblxuLyoqXG4gKiBQbGFjZSB0aGUgd2luZG93IGluIHRoZSBjZW50ZXIgb2YgdGhlIHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0NlbnRlcigpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1djJyk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgd2luZG93IHRpdGxlXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHRpdGxlXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRUaXRsZSh0aXRsZSkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV1QnICsgdGl0bGUpO1xufVxuXG4vKipcbiAqIE1ha2VzIHRoZSB3aW5kb3cgZ28gZnVsbHNjcmVlblxuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0Z1bGxzY3JlZW4oKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXRicpO1xufVxuXG4vKipcbiAqIFJldmVydHMgdGhlIHdpbmRvdyBmcm9tIGZ1bGxzY3JlZW5cbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dVbmZ1bGxzY3JlZW4oKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXZicpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIGluIGZ1bGwgc2NyZWVuIG1vZGUgb3Igbm90LlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IFRoZSBzdGF0ZSBvZiB0aGUgd2luZG93XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dJc0Z1bGxzY3JlZW4oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNGdWxsc2NyZWVuXCIpO1xufVxuXG4vKipcbiAqIFNldCB0aGUgU2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dzOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBTaXplIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHt3OiBudW1iZXIsIGg6IG51bWJlcn0+fSBUaGUgc2l6ZSBvZiB0aGUgd2luZG93XG5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0dldFNpemUoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0U2l6ZVwiKTtcbn1cblxuLyoqXG4gKiBTZXQgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93XG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRNYXhTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1daOicgKyB3aWR0aCArICc6JyArIGhlaWdodCk7XG59XG5cbi8qKlxuICogU2V0IHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXejonICsgd2lkdGggKyAnOicgKyBoZWlnaHQpO1xufVxuXG5cblxuLyoqXG4gKiBTZXQgdGhlIHdpbmRvdyBBbHdheXNPblRvcCBvciBub3Qgb24gdG9wXG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0QWx3YXlzT25Ub3AoYikge1xuXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXQVRQOicgKyAoYiA/ICcxJyA6ICcwJykpO1xufVxuXG5cblxuXG4vKipcbiAqIFNldCB0aGUgUG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gKiBAcGFyYW0ge251bWJlcn0geVxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93U2V0UG9zaXRpb24oeCwgeSkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3A6JyArIHggKyAnOicgKyB5KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIFBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHt4OiBudW1iZXIsIHk6IG51bWJlcn0+fSBUaGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93R2V0UG9zaXRpb24oKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93R2V0UG9zXCIpO1xufVxuXG4vKipcbiAqIEhpZGUgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0hpZGUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXSCcpO1xufVxuXG4vKipcbiAqIFNob3cgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1Nob3coKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXUycpO1xufVxuXG4vKipcbiAqIE1heGltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNYXhpbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dNJyk7XG59XG5cbi8qKlxuICogVG9nZ2xlIHRoZSBNYXhpbWlzZSBvZiB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VG9nZ2xlTWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXdCcpO1xufVxuXG4vKipcbiAqIFVubWF4aW1pc2UgdGhlIFdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd1VubWF4aW1pc2UoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdXVScpO1xufVxuXG4vKipcbiAqIFJldHVybnMgdGhlIHN0YXRlIG9mIHRoZSB3aW5kb3csIGkuZS4gd2hldGhlciB0aGUgd2luZG93IGlzIG1heGltaXNlZCBvciBub3QuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVGhlIHN0YXRlIG9mIHRoZSB3aW5kb3dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdpbmRvd0lzTWF4aW1pc2VkKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWF4aW1pc2VkXCIpO1xufVxuXG4vKipcbiAqIE1pbmltaXNlIHRoZSBXaW5kb3dcbiAqXG4gKiBAZXhwb3J0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dNaW5pbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1dtJyk7XG59XG5cbi8qKlxuICogVW5taW5pbWlzZSB0aGUgV2luZG93XG4gKlxuICogQGV4cG9ydFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93VW5taW5pbWlzZSgpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1d1Jyk7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNNaW5pbWlzZWQoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNaW5pbWlzZWRcIik7XG59XG5cbi8qKlxuICogUmV0dXJucyB0aGUgc3RhdGUgb2YgdGhlIHdpbmRvdywgaS5lLiB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgbm9ybWFsIG9yIG5vdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUaGUgc3RhdGUgb2YgdGhlIHdpbmRvd1xuICovXG5leHBvcnQgZnVuY3Rpb24gV2luZG93SXNOb3JtYWwoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNOb3JtYWxcIik7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7bnVtYmVyfSBSIFJlZFxuICogQHBhcmFtIHtudW1iZXJ9IEcgR3JlZW5cbiAqIEBwYXJhbSB7bnVtYmVyfSBCIEJsdWVcbiAqIEBwYXJhbSB7bnVtYmVyfSBBIEFscGhhXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaW5kb3dTZXRCYWNrZ3JvdW5kQ29sb3VyKFIsIEcsIEIsIEEpIHtcbiAgICBsZXQgcmdiYSA9IEpTT04uc3RyaW5naWZ5KHtyOiBSIHx8IDAsIGc6IEcgfHwgMCwgYjogQiB8fCAwLCBhOiBBIHx8IDI1NX0pO1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnV3I6JyArIHJnYmEpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQge0NhbGx9IGZyb20gXCIuL2NhbGxzXCI7XG5cblxuLyoqXG4gKiBHZXRzIHRoZSBhbGwgc2NyZWVucy4gQ2FsbCB0aGlzIGFuZXcgZWFjaCB0aW1lIHlvdSB3YW50IHRvIHJlZnJlc2ggZGF0YSBmcm9tIHRoZSB1bmRlcmx5aW5nIHdpbmRvd2luZyBzeXN0ZW0uXG4gKiBAZXhwb3J0XG4gKiBAdHlwZWRlZiB7aW1wb3J0KCcuLi93cmFwcGVyL3J1bnRpbWUnKS5TY3JlZW59IFNjcmVlblxuICogQHJldHVybiB7UHJvbWlzZTx7U2NyZWVuW119Pn0gVGhlIHNjcmVlbnNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNjcmVlbkdldEFsbCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpTY3JlZW5HZXRBbGxcIik7XG59XG4iLCAiLyoqXG4gKiBAZGVzY3JpcHRpb246IFVzZSB0aGUgc3lzdGVtIGRlZmF1bHQgYnJvd3NlciB0byBvcGVuIHRoZSB1cmxcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gQnJvd3Nlck9wZW5VUkwodXJsKSB7XG4gIHdpbmRvdy5XYWlsc0ludm9rZSgnQk86JyArIHVybCk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7Q2FsbH0gZnJvbSBcIi4vY2FsbHNcIjtcblxuLyoqXG4gKiBTZXQgdGhlIFNpemUgb2YgdGhlIHdpbmRvd1xuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGlwYm9hcmRTZXRUZXh0XCIsIFt0ZXh0XSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSB0ZXh0IGNvbnRlbnQgb2YgdGhlIGNsaXBib2FyZFxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8e3N0cmluZ30+fSBUZXh0IGNvbnRlbnQgb2YgdGhlIGNsaXBib2FyZFxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGlwYm9hcmRHZXRUZXh0KCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNsaXBib2FyZEdldFRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7RXZlbnRzT24sIEV2ZW50c09mZn0gZnJvbSBcIi4vZXZlbnRzXCI7XG5cbmNvbnN0IGZsYWdzID0ge1xuICAgIHJlZ2lzdGVyZWQ6IGZhbHNlLFxuICAgIGRlZmF1bHRVc2VEcm9wVGFyZ2V0OiB0cnVlLFxuICAgIHVzZURyb3BUYXJnZXQ6IHRydWUsXG4gICAgbmV4dERlYWN0aXZhdGU6IG51bGwsXG4gICAgbmV4dERlYWN0aXZhdGVUaW1lb3V0OiBudWxsLFxufTtcblxuY29uc3QgRFJPUF9UQVJHRVRfQUNUSVZFID0gXCJ3YWlscy1kcm9wLXRhcmdldC1hY3RpdmVcIjtcblxuLyoqXG4gKiBjaGVja1N0eWxlRHJvcFRhcmdldCBjaGVja3MgaWYgdGhlIHN0eWxlIGhhcyB0aGUgZHJvcCB0YXJnZXQgYXR0cmlidXRlXG4gKiBcbiAqIEBwYXJhbSB7Q1NTU3R5bGVEZWNsYXJhdGlvbn0gc3R5bGUgXG4gKiBAcmV0dXJucyBcbiAqL1xuZnVuY3Rpb24gY2hlY2tTdHlsZURyb3BUYXJnZXQoc3R5bGUpIHtcbiAgICBjb25zdCBjc3NEcm9wVmFsdWUgPSBzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wUHJvcGVydHkpLnRyaW0oKTtcbiAgICBpZiAoY3NzRHJvcFZhbHVlKSB7XG4gICAgICAgIGlmIChjc3NEcm9wVmFsdWUgPT09IHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wVmFsdWUpIHtcbiAgICAgICAgICAgIHJldHVybiB0cnVlO1xuICAgICAgICB9XG4gICAgICAgIC8vIGlmIHRoZSBlbGVtZW50IGhhcyB0aGUgZHJvcCB0YXJnZXQgYXR0cmlidXRlLCBidXQgXG4gICAgICAgIC8vIHRoZSB2YWx1ZSBpcyBub3QgY29ycmVjdCwgdGVybWluYXRlIGZpbmRpbmcgcHJvY2Vzcy5cbiAgICAgICAgLy8gVGhpcyBjYW4gYmUgdXNlZnVsIHRvIGJsb2NrIHNvbWUgY2hpbGQgZWxlbWVudHMgZnJvbSBiZWluZyBkcm9wIHRhcmdldHMuXG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG4gICAgcmV0dXJuIGZhbHNlO1xufVxuXG4vKipcbiAqIG9uRHJhZ092ZXIgaXMgY2FsbGVkIHdoZW4gdGhlIGRyYWdvdmVyIGV2ZW50IGlzIGVtaXR0ZWQuXG4gKiBAcGFyYW0ge0RyYWdFdmVudH0gZVxuICogQHJldHVybnNcbiAqL1xuZnVuY3Rpb24gb25EcmFnT3ZlcihlKSB7XG4gICAgLy8gQ2hlY2sgaWYgdGhpcyBpcyBhbiBleHRlcm5hbCBmaWxlIGRyb3Agb3IgaW50ZXJuYWwgSFRNTCBkcmFnXG4gICAgLy8gRXh0ZXJuYWwgZmlsZSBkcm9wcyB3aWxsIGhhdmUgXCJGaWxlc1wiIGluIHRoZSB0eXBlcyBhcnJheVxuICAgIC8vIEludGVybmFsIEhUTUwgZHJhZ3MgdHlwaWNhbGx5IGhhdmUgXCJ0ZXh0L3BsYWluXCIsIFwidGV4dC9odG1sXCIgb3IgY3VzdG9tIHR5cGVzXG4gICAgY29uc3QgaXNGaWxlRHJvcCA9IGUuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKFwiRmlsZXNcIik7XG5cbiAgICAvLyBPbmx5IGhhbmRsZSBleHRlcm5hbCBmaWxlIGRyb3BzLCBsZXQgaW50ZXJuYWwgSFRNTDUgZHJhZy1hbmQtZHJvcCB3b3JrIG5vcm1hbGx5XG4gICAgaWYgKCFpc0ZpbGVEcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBBTFdBWVMgcHJldmVudCBkZWZhdWx0IGZvciBmaWxlIGRyb3BzIHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIGUuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG5cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVXYWlsc0RyYWdBbmREcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGNvbnN0IGVsZW1lbnQgPSBlLnRhcmdldDtcblxuICAgIC8vIFRyaWdnZXIgZGVib3VuY2UgZnVuY3Rpb24gdG8gZGVhY3RpdmF0ZSBkcm9wIHRhcmdldHNcbiAgICBpZihmbGFncy5uZXh0RGVhY3RpdmF0ZSkgZmxhZ3MubmV4dERlYWN0aXZhdGUoKTtcblxuICAgIC8vIGlmIHRoZSBlbGVtZW50IGlzIG51bGwgb3IgZWxlbWVudCBpcyBub3QgY2hpbGQgb2YgZHJvcCB0YXJnZXQgZWxlbWVudFxuICAgIGlmICghZWxlbWVudCB8fCAhY2hlY2tTdHlsZURyb3BUYXJnZXQoZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KSkpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCBjdXJyZW50RWxlbWVudCA9IGVsZW1lbnQ7XG4gICAgd2hpbGUgKGN1cnJlbnRFbGVtZW50KSB7XG4gICAgICAgIC8vIGNoZWNrIGlmIGN1cnJlbnRFbGVtZW50IGlzIGRyb3AgdGFyZ2V0IGVsZW1lbnRcbiAgICAgICAgaWYgKGNoZWNrU3R5bGVEcm9wVGFyZ2V0KGdldENvbXB1dGVkU3R5bGUoY3VycmVudEVsZW1lbnQpKSkge1xuICAgICAgICAgICAgY3VycmVudEVsZW1lbnQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkUpO1xuICAgICAgICB9XG4gICAgICAgIGN1cnJlbnRFbGVtZW50ID0gY3VycmVudEVsZW1lbnQucGFyZW50RWxlbWVudDtcbiAgICB9XG59XG5cbi8qKlxuICogb25EcmFnTGVhdmUgaXMgY2FsbGVkIHdoZW4gdGhlIGRyYWdsZWF2ZSBldmVudCBpcyBlbWl0dGVkLlxuICogQHBhcmFtIHtEcmFnRXZlbnR9IGVcbiAqIEByZXR1cm5zXG4gKi9cbmZ1bmN0aW9uIG9uRHJhZ0xlYXZlKGUpIHtcbiAgICAvLyBDaGVjayBpZiB0aGlzIGlzIGFuIGV4dGVybmFsIGZpbGUgZHJvcCBvciBpbnRlcm5hbCBIVE1MIGRyYWdcbiAgICBjb25zdCBpc0ZpbGVEcm9wID0gZS5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoXCJGaWxlc1wiKTtcblxuICAgIC8vIE9ubHkgaGFuZGxlIGV4dGVybmFsIGZpbGUgZHJvcHMsIGxldCBpbnRlcm5hbCBIVE1MNSBkcmFnLWFuZC1kcm9wIHdvcmsgbm9ybWFsbHlcbiAgICBpZiAoIWlzRmlsZURyb3ApIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIEFMV0FZUyBwcmV2ZW50IGRlZmF1bHQgZm9yIGZpbGUgZHJvcHMgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICBlLnByZXZlbnREZWZhdWx0KCk7XG5cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVXYWlsc0RyYWdBbmREcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIWZsYWdzLnVzZURyb3BUYXJnZXQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIEZpbmQgdGhlIGNsb3NlIGRyb3AgdGFyZ2V0IGVsZW1lbnRcbiAgICBpZiAoIWUudGFyZ2V0IHx8ICFjaGVja1N0eWxlRHJvcFRhcmdldChnZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KSkpIHtcbiAgICAgICAgcmV0dXJuIG51bGw7XG4gICAgfVxuXG4gICAgLy8gVHJpZ2dlciBkZWJvdW5jZSBmdW5jdGlvbiB0byBkZWFjdGl2YXRlIGRyb3AgdGFyZ2V0c1xuICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xuICAgIFxuICAgIC8vIFVzZSBkZWJvdW5jZSB0ZWNobmlxdWUgdG8gdGFjbGUgZHJhZ2xlYXZlIGV2ZW50cyBvbiBvdmVybGFwcGluZyBlbGVtZW50cyBhbmQgZHJvcCB0YXJnZXQgZWxlbWVudHNcbiAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZSA9ICgpID0+IHtcbiAgICAgICAgLy8gRGVhY3RpdmF0ZSBhbGwgZHJvcCB0YXJnZXRzLCBuZXcgZHJvcCB0YXJnZXQgd2lsbCBiZSBhY3RpdmF0ZWQgb24gbmV4dCBkcmFnb3ZlciBldmVudFxuICAgICAgICBBcnJheS5mcm9tKGRvY3VtZW50LmdldEVsZW1lbnRzQnlDbGFzc05hbWUoRFJPUF9UQVJHRVRfQUNUSVZFKSkuZm9yRWFjaChlbCA9PiBlbC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRSkpO1xuICAgICAgICAvLyBSZXNldCBuZXh0RGVhY3RpdmF0ZVxuICAgICAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZSA9IG51bGw7XG4gICAgICAgIC8vIENsZWFyIHRpbWVvdXRcbiAgICAgICAgaWYgKGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCkge1xuICAgICAgICAgICAgY2xlYXJUaW1lb3V0KGZsYWdzLm5leHREZWFjdGl2YXRlVGltZW91dCk7XG4gICAgICAgICAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZVRpbWVvdXQgPSBudWxsO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gU2V0IHRpbWVvdXQgdG8gZGVhY3RpdmF0ZSBkcm9wIHRhcmdldHMgaWYgbm90IHRyaWdnZXJlZCBieSBuZXh0IGRyYWcgZXZlbnRcbiAgICBmbGFncy5uZXh0RGVhY3RpdmF0ZVRpbWVvdXQgPSBzZXRUaW1lb3V0KCgpID0+IHtcbiAgICAgICAgaWYoZmxhZ3MubmV4dERlYWN0aXZhdGUpIGZsYWdzLm5leHREZWFjdGl2YXRlKCk7XG4gICAgfSwgNTApO1xufVxuXG4vKipcbiAqIG9uRHJvcCBpcyBjYWxsZWQgd2hlbiB0aGUgZHJvcCBldmVudCBpcyBlbWl0dGVkLlxuICogQHBhcmFtIHtEcmFnRXZlbnR9IGVcbiAqIEByZXR1cm5zXG4gKi9cbmZ1bmN0aW9uIG9uRHJvcChlKSB7XG4gICAgLy8gQ2hlY2sgaWYgdGhpcyBpcyBhbiBleHRlcm5hbCBmaWxlIGRyb3Agb3IgaW50ZXJuYWwgSFRNTCBkcmFnXG4gICAgY29uc3QgaXNGaWxlRHJvcCA9IGUuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKFwiRmlsZXNcIik7XG5cbiAgICAvLyBPbmx5IGhhbmRsZSBleHRlcm5hbCBmaWxlIGRyb3BzLCBsZXQgaW50ZXJuYWwgSFRNTDUgZHJhZy1hbmQtZHJvcCB3b3JrIG5vcm1hbGx5XG4gICAgaWYgKCFpc0ZpbGVEcm9wKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBBTFdBWVMgcHJldmVudCBkZWZhdWx0IGZvciBmaWxlIGRyb3BzIHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuXG4gICAgaWYgKCF3aW5kb3cud2FpbHMuZmxhZ3MuZW5hYmxlV2FpbHNEcmFnQW5kRHJvcCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKENhblJlc29sdmVGaWxlUGF0aHMoKSkge1xuICAgICAgICAvLyBwcm9jZXNzIGZpbGVzXG4gICAgICAgIGxldCBmaWxlcyA9IFtdO1xuICAgICAgICBpZiAoZS5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcbiAgICAgICAgICAgIGZpbGVzID0gWy4uLmUuZGF0YVRyYW5zZmVyLml0ZW1zXS5tYXAoKGl0ZW0sIGkpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuIGl0ZW0uZ2V0QXNGaWxlKCk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBmaWxlcyA9IFsuLi5lLmRhdGFUcmFuc2Zlci5maWxlc107XG4gICAgICAgIH1cbiAgICAgICAgd2luZG93LnJ1bnRpbWUuUmVzb2x2ZUZpbGVQYXRocyhlLngsIGUueSwgZmlsZXMpO1xuICAgIH1cblxuICAgIGlmICghZmxhZ3MudXNlRHJvcFRhcmdldCkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gVHJpZ2dlciBkZWJvdW5jZSBmdW5jdGlvbiB0byBkZWFjdGl2YXRlIGRyb3AgdGFyZ2V0c1xuICAgIGlmKGZsYWdzLm5leHREZWFjdGl2YXRlKSBmbGFncy5uZXh0RGVhY3RpdmF0ZSgpO1xuXG4gICAgLy8gRGVhY3RpdmF0ZSBhbGwgZHJvcCB0YXJnZXRzXG4gICAgQXJyYXkuZnJvbShkb2N1bWVudC5nZXRFbGVtZW50c0J5Q2xhc3NOYW1lKERST1BfVEFSR0VUX0FDVElWRSkpLmZvckVhY2goZWwgPT4gZWwuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkUpKTtcbn1cblxuLyoqXG4gKiBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyBjaGVja3MgdGhlIGJyb3dzZXIncyBjYXBhYmlsaXR5IG9mIHNlbmRpbmcgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHNcbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn1cbiAqIEBjb25zdHJ1Y3RvclxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FuUmVzb2x2ZUZpbGVQYXRocygpIHtcbiAgICByZXR1cm4gd2luZG93LmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgIT0gbnVsbDtcbn1cblxuLyoqXG4gKiBSZXNvbHZlRmlsZVBhdGhzIHNlbmRzIGRyb3AgZXZlbnRzIHRvIHRoZSBHTyBzaWRlIHRvIHJlc29sdmUgZmlsZSBwYXRocyBvbiB3aW5kb3dzLlxuICpcbiAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gKiBAcGFyYW0ge251bWJlcn0geVxuICogQHBhcmFtIHthbnlbXX0gZmlsZXNcbiAqIEBjb25zdHJ1Y3RvclxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVzb2x2ZUZpbGVQYXRocyh4LCB5LCBmaWxlcykge1xuICAgIC8vIE9ubHkgZm9yIHdpbmRvd3Mgd2VidmlldzIgPj0gMS4wLjE3NzQuMzBcbiAgICAvLyBodHRwczovL2xlYXJuLm1pY3Jvc29mdC5jb20vZW4tdXMvbWljcm9zb2Z0LWVkZ2Uvd2VidmlldzIvcmVmZXJlbmNlL3dpbjMyL2ljb3Jld2VidmlldzJ3ZWJtZXNzYWdlcmVjZWl2ZWRldmVudGFyZ3MyP3ZpZXc9d2VidmlldzItMS4wLjE4MjMuMzIjYXBwbGllcy10b1xuICAgIGlmICh3aW5kb3cuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cykge1xuICAgICAgICBjaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyhgZmlsZTpkcm9wOiR7eH06JHt5fWAsIGZpbGVzKTtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbGJhY2sgZm9yIE9uRmlsZURyb3AgcmV0dXJucyBhIHNsaWNlIG9mIGZpbGUgcGF0aCBzdHJpbmdzIHdoZW4gYSBkcm9wIGlzIGZpbmlzaGVkLlxuICpcbiAqIEBleHBvcnRcbiAqIEBjYWxsYmFjayBPbkZpbGVEcm9wQ2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSB4IC0geCBjb29yZGluYXRlIG9mIHRoZSBkcm9wXG4gKiBAcGFyYW0ge251bWJlcn0geSAtIHkgY29vcmRpbmF0ZSBvZiB0aGUgZHJvcFxuICogQHBhcmFtIHtzdHJpbmdbXX0gcGF0aHMgLSBBIGxpc3Qgb2YgZmlsZSBwYXRocy5cbiAqL1xuXG4vKipcbiAqIE9uRmlsZURyb3AgbGlzdGVucyB0byBkcmFnIGFuZCBkcm9wIGV2ZW50cyBhbmQgY2FsbHMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGNvb3JkaW5hdGVzIG9mIHRoZSBkcm9wIGFuZCBhbiBhcnJheSBvZiBwYXRoIHN0cmluZ3MuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtPbkZpbGVEcm9wQ2FsbGJhY2t9IGNhbGxiYWNrIC0gQ2FsbGJhY2sgZm9yIE9uRmlsZURyb3AgcmV0dXJucyBhIHNsaWNlIG9mIGZpbGUgcGF0aCBzdHJpbmdzIHdoZW4gYSBkcm9wIGlzIGZpbmlzaGVkLlxuICogQHBhcmFtIHtib29sZWFufSBbdXNlRHJvcFRhcmdldD10cnVlXSAtIE9ubHkgY2FsbCB0aGUgY2FsbGJhY2sgd2hlbiB0aGUgZHJvcCBmaW5pc2hlZCBvbiBhbiBlbGVtZW50IHRoYXQgaGFzIHRoZSBkcm9wIHRhcmdldCBzdHlsZS4gKC0td2FpbHMtZHJvcC10YXJnZXQpXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbkZpbGVEcm9wKGNhbGxiYWNrLCB1c2VEcm9wVGFyZ2V0KSB7XG4gICAgaWYgKHR5cGVvZiBjYWxsYmFjayAhPT0gXCJmdW5jdGlvblwiKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoXCJEcmFnQW5kRHJvcENhbGxiYWNrIGlzIG5vdCBhIGZ1bmN0aW9uXCIpO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKGZsYWdzLnJlZ2lzdGVyZWQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBmbGFncy5yZWdpc3RlcmVkID0gdHJ1ZTtcblxuICAgIGNvbnN0IHVEVFBUID0gdHlwZW9mIHVzZURyb3BUYXJnZXQ7XG4gICAgZmxhZ3MudXNlRHJvcFRhcmdldCA9IHVEVFBUID09PSBcInVuZGVmaW5lZFwiIHx8IHVEVFBUICE9PSBcImJvb2xlYW5cIiA/IGZsYWdzLmRlZmF1bHRVc2VEcm9wVGFyZ2V0IDogdXNlRHJvcFRhcmdldDtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBvbkRyYWdPdmVyKTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgb25EcmFnTGVhdmUpO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgb25Ecm9wKTtcblxuICAgIGxldCBjYiA9IGNhbGxiYWNrO1xuICAgIGlmIChmbGFncy51c2VEcm9wVGFyZ2V0KSB7XG4gICAgICAgIGNiID0gZnVuY3Rpb24gKHgsIHksIHBhdGhzKSB7XG4gICAgICAgICAgICBjb25zdCBlbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KVxuICAgICAgICAgICAgLy8gaWYgdGhlIGVsZW1lbnQgaXMgbnVsbCBvciBlbGVtZW50IGlzIG5vdCBjaGlsZCBvZiBkcm9wIHRhcmdldCBlbGVtZW50LCByZXR1cm4gbnVsbFxuICAgICAgICAgICAgaWYgKCFlbGVtZW50IHx8ICFjaGVja1N0eWxlRHJvcFRhcmdldChnZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpKSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBudWxsO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbGJhY2soeCwgeSwgcGF0aHMpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgRXZlbnRzT24oXCJ3YWlsczpmaWxlLWRyb3BcIiwgY2IpO1xufVxuXG4vKipcbiAqIE9uRmlsZURyb3BPZmYgcmVtb3ZlcyB0aGUgZHJhZyBhbmQgZHJvcCBsaXN0ZW5lcnMgYW5kIGhhbmRsZXJzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25GaWxlRHJvcE9mZigpIHtcbiAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCBvbkRyYWdPdmVyKTtcbiAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgb25EcmFnTGVhdmUpO1xuICAgIHdpbmRvdy5yZW1vdmVFdmVudExpc3RlbmVyKCdkcm9wJywgb25Ecm9wKTtcbiAgICBFdmVudHNPZmYoXCJ3YWlsczpmaWxlLWRyb3BcIik7XG4gICAgZmxhZ3MucmVnaXN0ZXJlZCA9IGZhbHNlO1xufVxuIiwgIi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmV4cG9ydCBmdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KSB7XG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcbiAgICAgICAgY2FzZSBcInNob3dcIjpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGRlZmF1bHQ6IHtcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgICAgICAgICBpZiAoZWxlbWVudC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBhbmQgYWN0aW9uIGlzIG9uIHRoZSBzZWxlY3RlZCBlbGVtZW50c1xuICAgICAgICAgICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgICAgICAgICAgY29uc3QgaGFzU2VsZWN0aW9uID0gKHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDApXG4gICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgICAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXG4gICAgICAgICAgICBpZiAoZWxlbWVudC50YWdOYW1lID09PSBcIklOUFVUXCIgfHwgZWxlbWVudC50YWdOYW1lID09PSBcIlRFWFRBUkVBXCIpIHtcbiAgICAgICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghZWxlbWVudC5yZWFkT25seSAmJiAhZWxlbWVudC5kaXNhYmxlZCkpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgfVxuICAgIH1cbn1cbiIsICIvKlxuIF8gICAgICAgX18gICAgICBfIF9fXG58IHwgICAgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtDYWxsfSBmcm9tIFwiLi9jYWxsc1wiO1xuXG4vKipcbiAqIEluaXRpYWxpemUgdGhlIG5vdGlmaWNhdGlvbiBzZXJ2aWNlIGZvciB0aGUgYXBwbGljYXRpb24uXG4gKiBUaGlzIG11c3QgYmUgY2FsbGVkIGJlZm9yZSBzZW5kaW5nIGFueSBub3RpZmljYXRpb25zLlxuICogT24gbWFjT1MsIHRoaXMgYWxzbyBlbnN1cmVzIHRoZSBub3RpZmljYXRpb24gZGVsZWdhdGUgaXMgcHJvcGVybHkgaW5pdGlhbGl6ZWQuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEluaXRpYWxpemVOb3RpZmljYXRpb25zKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkluaXRpYWxpemVOb3RpZmljYXRpb25zXCIpO1xufVxuXG4vKipcbiAqIENsZWFuIHVwIG5vdGlmaWNhdGlvbiByZXNvdXJjZXMgYW5kIHJlbGVhc2UgYW55IGhlbGQgY29ubmVjdGlvbnMuXG4gKiBUaGlzIHNob3VsZCBiZSBjYWxsZWQgd2hlbiBzaHV0dGluZyBkb3duIHRoZSBhcHBsaWNhdGlvbiB0byBwcm9wZXJseSByZWxlYXNlIHJlc291cmNlc1xuICogKHByaW1hcmlseSBuZWVkZWQgb24gTGludXggdG8gY2xvc2UgRC1CdXMgY29ubmVjdGlvbnMpLlxuICpcbiAqIEBleHBvcnRcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDbGVhbnVwTm90aWZpY2F0aW9ucygpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpDbGVhbnVwTm90aWZpY2F0aW9uc1wiKTtcbn1cblxuLyoqXG4gKiBDaGVjayBpZiBub3RpZmljYXRpb25zIGFyZSBhdmFpbGFibGUgb24gdGhlIGN1cnJlbnQgcGxhdGZvcm0uXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVHJ1ZSBpZiBub3RpZmljYXRpb25zIGFyZSBhdmFpbGFibGUsIGZhbHNlIG90aGVyd2lzZVxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNOb3RpZmljYXRpb25BdmFpbGFibGUoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6SXNOb3RpZmljYXRpb25BdmFpbGFibGVcIik7XG59XG5cbi8qKlxuICogUmVxdWVzdCBub3RpZmljYXRpb24gYXV0aG9yaXphdGlvbiBmcm9tIHRoZSB1c2VyLlxuICogT24gbWFjT1MsIHRoaXMgcHJvbXB0cyB0aGUgdXNlciB0byBhbGxvdyBub3RpZmljYXRpb25zLlxuICogT24gb3RoZXIgcGxhdGZvcm1zLCB0aGlzIGFsd2F5cyByZXR1cm5zIHRydWUuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gVHJ1ZSBpZiBhdXRob3JpemF0aW9uIHdhcyBncmFudGVkLCBmYWxzZSBvdGhlcndpc2VcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlcXVlc3ROb3RpZmljYXRpb25BdXRob3JpemF0aW9uKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlJlcXVlc3ROb3RpZmljYXRpb25BdXRob3JpemF0aW9uXCIpO1xufVxuXG4vKipcbiAqIENoZWNrIHRoZSBjdXJyZW50IG5vdGlmaWNhdGlvbiBhdXRob3JpemF0aW9uIHN0YXR1cy5cbiAqIE9uIG1hY09TLCB0aGlzIGNoZWNrcyBpZiB0aGUgYXBwIGhhcyBub3RpZmljYXRpb24gcGVybWlzc2lvbnMuXG4gKiBPbiBvdGhlciBwbGF0Zm9ybXMsIHRoaXMgYWx3YXlzIHJldHVybnMgdHJ1ZS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSBUcnVlIGlmIGF1dGhvcml6ZWQsIGZhbHNlIG90aGVyd2lzZVxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2hlY2tOb3RpZmljYXRpb25BdXRob3JpemF0aW9uKCkge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkNoZWNrTm90aWZpY2F0aW9uQXV0aG9yaXphdGlvblwiKTtcbn1cblxuLyoqXG4gKiBTZW5kIGEgYmFzaWMgbm90aWZpY2F0aW9uIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBUaGUgbm90aWZpY2F0aW9uIHdpbGwgZGlzcGxheSB3aXRoIHRoZSBwcm92aWRlZCB0aXRsZSwgc3VidGl0bGUgKGlmIHN1cHBvcnRlZCksIGFuZCBib2R5IHRleHQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtPYmplY3R9IG9wdGlvbnMgLSBOb3RpZmljYXRpb24gb3B0aW9uc1xuICogQHBhcmFtIHtzdHJpbmd9IG9wdGlvbnMuaWQgLSBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIG5vdGlmaWNhdGlvblxuICogQHBhcmFtIHtzdHJpbmd9IG9wdGlvbnMudGl0bGUgLSBOb3RpZmljYXRpb24gdGl0bGVcbiAqIEBwYXJhbSB7c3RyaW5nfSBbb3B0aW9ucy5zdWJ0aXRsZV0gLSBOb3RpZmljYXRpb24gc3VidGl0bGUgKG1hY09TIGFuZCBMaW51eCBvbmx5KVxuICogQHBhcmFtIHtzdHJpbmd9IFtvcHRpb25zLmJvZHldIC0gTm90aWZpY2F0aW9uIGJvZHkgdGV4dFxuICogQHBhcmFtIHtzdHJpbmd9IFtvcHRpb25zLmNhdGVnb3J5SWRdIC0gQ2F0ZWdvcnkgSUQgZm9yIGFjdGlvbiBidXR0b25zIChyZXF1aXJlcyBTZW5kTm90aWZpY2F0aW9uV2l0aEFjdGlvbnMpXG4gKiBAcGFyYW0ge09iamVjdDxzdHJpbmcsIGFueT59IFtvcHRpb25zLmRhdGFdIC0gQWRkaXRpb25hbCB1c2VyIGRhdGEgdG8gYXR0YWNoIHRvIHRoZSBub3RpZmljYXRpb25cbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZW5kTm90aWZpY2F0aW9uKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpTZW5kTm90aWZpY2F0aW9uXCIsIFtvcHRpb25zXSk7XG59XG5cbi8qKlxuICogU2VuZCBhIG5vdGlmaWNhdGlvbiB3aXRoIGFjdGlvbiBidXR0b25zLlxuICogQSBOb3RpZmljYXRpb25DYXRlZ29yeSBtdXN0IGJlIHJlZ2lzdGVyZWQgZmlyc3QgdXNpbmcgUmVnaXN0ZXJOb3RpZmljYXRpb25DYXRlZ29yeS5cbiAqIFRoZSBvcHRpb25zLmNhdGVnb3J5SWQgbXVzdCBtYXRjaCBhIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCBjYXRlZ29yeSBJRC5cbiAqIElmIHRoZSBjYXRlZ29yeSBpcyBub3QgZm91bmQsIGEgYmFzaWMgbm90aWZpY2F0aW9uIHdpbGwgYmUgc2VudCBpbnN0ZWFkLlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7T2JqZWN0fSBvcHRpb25zIC0gTm90aWZpY2F0aW9uIG9wdGlvbnNcbiAqIEBwYXJhbSB7c3RyaW5nfSBvcHRpb25zLmlkIC0gVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBub3RpZmljYXRpb25cbiAqIEBwYXJhbSB7c3RyaW5nfSBvcHRpb25zLnRpdGxlIC0gTm90aWZpY2F0aW9uIHRpdGxlXG4gKiBAcGFyYW0ge3N0cmluZ30gW29wdGlvbnMuc3VidGl0bGVdIC0gTm90aWZpY2F0aW9uIHN1YnRpdGxlIChtYWNPUyBhbmQgTGludXggb25seSlcbiAqIEBwYXJhbSB7c3RyaW5nfSBbb3B0aW9ucy5ib2R5XSAtIE5vdGlmaWNhdGlvbiBib2R5IHRleHRcbiAqIEBwYXJhbSB7c3RyaW5nfSBvcHRpb25zLmNhdGVnb3J5SWQgLSBDYXRlZ29yeSBJRCB0aGF0IG1hdGNoZXMgYSByZWdpc3RlcmVkIGNhdGVnb3J5XG4gKiBAcGFyYW0ge09iamVjdDxzdHJpbmcsIGFueT59IFtvcHRpb25zLmRhdGFdIC0gQWRkaXRpb25hbCB1c2VyIGRhdGEgdG8gYXR0YWNoIHRvIHRoZSBub3RpZmljYXRpb25cbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZW5kTm90aWZpY2F0aW9uV2l0aEFjdGlvbnMob3B0aW9ucykge1xuICAgIHJldHVybiBDYWxsKFwiOndhaWxzOlNlbmROb3RpZmljYXRpb25XaXRoQWN0aW9uc1wiLCBbb3B0aW9uc10pO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVyIGEgbm90aWZpY2F0aW9uIGNhdGVnb3J5IHRoYXQgY2FuIGJlIHVzZWQgd2l0aCBTZW5kTm90aWZpY2F0aW9uV2l0aEFjdGlvbnMuXG4gKiBDYXRlZ29yaWVzIGRlZmluZSB0aGUgYWN0aW9uIGJ1dHRvbnMgYW5kIG9wdGlvbmFsIHJlcGx5IGZpZWxkcyB0aGF0IHdpbGwgYXBwZWFyIG9uIG5vdGlmaWNhdGlvbnMuXG4gKiBSZWdpc3RlcmluZyBhIGNhdGVnb3J5IHdpdGggdGhlIHNhbWUgSUQgYXMgYSBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgY2F0ZWdvcnkgd2lsbCBvdmVycmlkZSBpdC5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge09iamVjdH0gY2F0ZWdvcnkgLSBOb3RpZmljYXRpb24gY2F0ZWdvcnkgZGVmaW5pdGlvblxuICogQHBhcmFtIHtzdHJpbmd9IGNhdGVnb3J5LmlkIC0gVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBjYXRlZ29yeVxuICogQHBhcmFtIHtBcnJheTxPYmplY3Q+fSBbY2F0ZWdvcnkuYWN0aW9uc10gLSBBcnJheSBvZiBhY3Rpb24gYnV0dG9uc1xuICogQHBhcmFtIHtzdHJpbmd9IGNhdGVnb3J5LmFjdGlvbnNbXS5pZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgYWN0aW9uXG4gKiBAcGFyYW0ge3N0cmluZ30gY2F0ZWdvcnkuYWN0aW9uc1tdLnRpdGxlIC0gRGlzcGxheSB0aXRsZSBmb3IgdGhlIGFjdGlvbiBidXR0b25cbiAqIEBwYXJhbSB7Ym9vbGVhbn0gW2NhdGVnb3J5LmFjdGlvbnNbXS5kZXN0cnVjdGl2ZV0gLSBXaGV0aGVyIHRoZSBhY3Rpb24gaXMgZGVzdHJ1Y3RpdmUgKG1hY09TLXNwZWNpZmljKVxuICogQHBhcmFtIHtib29sZWFufSBbY2F0ZWdvcnkuaGFzUmVwbHlGaWVsZF0gLSBXaGV0aGVyIHRvIGluY2x1ZGUgYSB0ZXh0IGlucHV0IGZpZWxkIGZvciByZXBsaWVzXG4gKiBAcGFyYW0ge3N0cmluZ30gW2NhdGVnb3J5LnJlcGx5UGxhY2Vob2xkZXJdIC0gUGxhY2Vob2xkZXIgdGV4dCBmb3IgdGhlIHJlcGx5IGZpZWxkIChyZXF1aXJlZCBpZiBoYXNSZXBseUZpZWxkIGlzIHRydWUpXG4gKiBAcGFyYW0ge3N0cmluZ30gW2NhdGVnb3J5LnJlcGx5QnV0dG9uVGl0bGVdIC0gVGl0bGUgZm9yIHRoZSByZXBseSBidXR0b24gKHJlcXVpcmVkIGlmIGhhc1JlcGx5RmllbGQgaXMgdHJ1ZSlcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZWdpc3Rlck5vdGlmaWNhdGlvbkNhdGVnb3J5KGNhdGVnb3J5KSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVnaXN0ZXJOb3RpZmljYXRpb25DYXRlZ29yeVwiLCBbY2F0ZWdvcnldKTtcbn1cblxuLyoqXG4gKiBSZW1vdmUgYSBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgbm90aWZpY2F0aW9uIGNhdGVnb3J5LlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBjYXRlZ29yeUlkIC0gVGhlIElEIG9mIHRoZSBjYXRlZ29yeSB0byByZW1vdmVcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZW1vdmVOb3RpZmljYXRpb25DYXRlZ29yeShjYXRlZ29yeUlkKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVtb3ZlTm90aWZpY2F0aW9uQ2F0ZWdvcnlcIiwgW2NhdGVnb3J5SWRdKTtcbn1cblxuLyoqXG4gKiBSZW1vdmUgYWxsIHBlbmRpbmcgbm90aWZpY2F0aW9ucyBmcm9tIHRoZSBub3RpZmljYXRpb24gY2VudGVyLlxuICogT24gV2luZG93cywgdGhpcyBpcyBhIG5vLW9wIGFzIHRoZSBwbGF0Zm9ybSBtYW5hZ2VzIG5vdGlmaWNhdGlvbiBsaWZlY3ljbGUgYXV0b21hdGljYWxseS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVtb3ZlQWxsUGVuZGluZ05vdGlmaWNhdGlvbnMoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVtb3ZlQWxsUGVuZGluZ05vdGlmaWNhdGlvbnNcIik7XG59XG5cbi8qKlxuICogUmVtb3ZlIGEgc3BlY2lmaWMgcGVuZGluZyBub3RpZmljYXRpb24gYnkgaXRzIGlkZW50aWZpZXIuXG4gKiBPbiBXaW5kb3dzLCB0aGlzIGlzIGEgbm8tb3AgYXMgdGhlIHBsYXRmb3JtIG1hbmFnZXMgbm90aWZpY2F0aW9uIGxpZmVjeWNsZSBhdXRvbWF0aWNhbGx5LlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZGVudGlmaWVyIC0gVGhlIElEIG9mIHRoZSBub3RpZmljYXRpb24gdG8gcmVtb3ZlXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVtb3ZlUGVuZGluZ05vdGlmaWNhdGlvbihpZGVudGlmaWVyKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVtb3ZlUGVuZGluZ05vdGlmaWNhdGlvblwiLCBbaWRlbnRpZmllcl0pO1xufVxuXG4vKipcbiAqIFJlbW92ZSBhbGwgZGVsaXZlcmVkIG5vdGlmaWNhdGlvbnMgZnJvbSB0aGUgbm90aWZpY2F0aW9uIGNlbnRlci5cbiAqIE9uIFdpbmRvd3MsIHRoaXMgaXMgYSBuby1vcCBhcyB0aGUgcGxhdGZvcm0gbWFuYWdlcyBub3RpZmljYXRpb24gbGlmZWN5Y2xlIGF1dG9tYXRpY2FsbHkuXG4gKlxuICogQGV4cG9ydFxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbW92ZUFsbERlbGl2ZXJlZE5vdGlmaWNhdGlvbnMoKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVtb3ZlQWxsRGVsaXZlcmVkTm90aWZpY2F0aW9uc1wiKTtcbn1cblxuLyoqXG4gKiBSZW1vdmUgYSBzcGVjaWZpYyBkZWxpdmVyZWQgbm90aWZpY2F0aW9uIGJ5IGl0cyBpZGVudGlmaWVyLlxuICogT24gV2luZG93cywgdGhpcyBpcyBhIG5vLW9wIGFzIHRoZSBwbGF0Zm9ybSBtYW5hZ2VzIG5vdGlmaWNhdGlvbiBsaWZlY3ljbGUgYXV0b21hdGljYWxseS5cbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gaWRlbnRpZmllciAtIFRoZSBJRCBvZiB0aGUgbm90aWZpY2F0aW9uIHRvIHJlbW92ZVxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbW92ZURlbGl2ZXJlZE5vdGlmaWNhdGlvbihpZGVudGlmaWVyKSB7XG4gICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6UmVtb3ZlRGVsaXZlcmVkTm90aWZpY2F0aW9uXCIsIFtpZGVudGlmaWVyXSk7XG59XG5cbi8qKlxuICogUmVtb3ZlIGEgbm90aWZpY2F0aW9uIGJ5IGl0cyBpZGVudGlmaWVyLlxuICogVGhpcyBpcyBhIGNvbnZlbmllbmNlIGZ1bmN0aW9uIHRoYXQgd29ya3MgYWNyb3NzIHBsYXRmb3Jtcy5cbiAqIE9uIG1hY09TLCB1c2UgdGhlIG1vcmUgc3BlY2lmaWMgUmVtb3ZlUGVuZGluZ05vdGlmaWNhdGlvbiBvciBSZW1vdmVEZWxpdmVyZWROb3RpZmljYXRpb24gZnVuY3Rpb25zLlxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZGVudGlmaWVyIC0gVGhlIElEIG9mIHRoZSBub3RpZmljYXRpb24gdG8gcmVtb3ZlXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVtb3ZlTm90aWZpY2F0aW9uKGlkZW50aWZpZXIpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpSZW1vdmVOb3RpZmljYXRpb25cIiwgW2lkZW50aWZpZXJdKTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuaW1wb3J0ICogYXMgTG9nIGZyb20gJy4vbG9nJztcbmltcG9ydCB7XG4gIGV2ZW50TGlzdGVuZXJzLFxuICBFdmVudHNFbWl0LFxuICBFdmVudHNOb3RpZnksXG4gIEV2ZW50c09mZixcbiAgRXZlbnRzT2ZmQWxsLFxuICBFdmVudHNPbixcbiAgRXZlbnRzT25jZSxcbiAgRXZlbnRzT25NdWx0aXBsZSxcbn0gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQgeyBDYWxsLCBDYWxsYmFjaywgY2FsbGJhY2tzIH0gZnJvbSAnLi9jYWxscyc7XG5pbXBvcnQgeyBTZXRCaW5kaW5ncyB9IGZyb20gXCIuL2JpbmRpbmdzXCI7XG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQgKiBhcyBTY3JlZW4gZnJvbSBcIi4vc2NyZWVuXCI7XG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXJcIjtcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmRcIjtcbmltcG9ydCAqIGFzIERyYWdBbmREcm9wIGZyb20gXCIuL2RyYWdhbmRkcm9wXCI7XG5pbXBvcnQgKiBhcyBDb250ZXh0TWVudSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xuaW1wb3J0ICogYXMgTm90aWZpY2F0aW9ucyBmcm9tIFwiLi9ub3RpZmljYXRpb25zXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnUScpO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ1MnKTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgd2luZG93LldhaWxzSW52b2tlKCdIJyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpIHtcbiAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpFbnZpcm9ubWVudFwiKTtcbn1cblxuLy8gVGhlIEpTIHJ1bnRpbWVcbndpbmRvdy5ydW50aW1lID0ge1xuICAgIC4uLkxvZyxcbiAgICAuLi5XaW5kb3csXG4gICAgLi4uQnJvd3NlcixcbiAgICAuLi5TY3JlZW4sXG4gICAgLi4uQ2xpcGJvYXJkLFxuICAgIC4uLkRyYWdBbmREcm9wLFxuICAgIC4uLk5vdGlmaWNhdGlvbnMsXG4gICAgRXZlbnRzT24sXG4gICAgRXZlbnRzT25jZSxcbiAgICBFdmVudHNPbk11bHRpcGxlLFxuICAgIEV2ZW50c0VtaXQsXG4gICAgRXZlbnRzT2ZmLFxuICAgIEV2ZW50c09mZkFsbCxcbiAgICBFbnZpcm9ubWVudCxcbiAgICBTaG93LFxuICAgIEhpZGUsXG4gICAgUXVpdFxufTtcblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cud2FpbHMgPSB7XG4gICAgQ2FsbGJhY2ssXG4gICAgRXZlbnRzTm90aWZ5LFxuICAgIFNldEJpbmRpbmdzLFxuICAgIGV2ZW50TGlzdGVuZXJzLFxuICAgIGNhbGxiYWNrcyxcbiAgICBmbGFnczoge1xuICAgICAgICBkaXNhYmxlU2Nyb2xsYmFyRHJhZzogZmFsc2UsXG4gICAgICAgIGRpc2FibGVEZWZhdWx0Q29udGV4dE1lbnU6IGZhbHNlLFxuICAgICAgICBlbmFibGVSZXNpemU6IGZhbHNlLFxuICAgICAgICBkZWZhdWx0Q3Vyc29yOiBudWxsLFxuICAgICAgICBib3JkZXJUaGlja25lc3M6IDYsXG4gICAgICAgIHNob3VsZERyYWc6IGZhbHNlLFxuICAgICAgICBkZWZlckRyYWdUb01vdXNlTW92ZTogdHJ1ZSxcbiAgICAgICAgY3NzRHJhZ1Byb3BlcnR5OiBcIi0td2FpbHMtZHJhZ2dhYmxlXCIsXG4gICAgICAgIGNzc0RyYWdWYWx1ZTogXCJkcmFnXCIsXG4gICAgICAgIGNzc0Ryb3BQcm9wZXJ0eTogXCItLXdhaWxzLWRyb3AtdGFyZ2V0XCIsXG4gICAgICAgIGNzc0Ryb3BWYWx1ZTogXCJkcm9wXCIsXG4gICAgICAgIGVuYWJsZVdhaWxzRHJhZ0FuZERyb3A6IGZhbHNlLFxuICAgIH1cbn07XG5cbi8vIFNldCB0aGUgYmluZGluZ3NcbmlmICh3aW5kb3cud2FpbHNiaW5kaW5ncykge1xuICAgIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncyh3aW5kb3cud2FpbHNiaW5kaW5ncyk7XG4gICAgZGVsZXRlIHdpbmRvdy53YWlscy5TZXRCaW5kaW5ncztcbn1cblxuLy8gKGJvb2wpIFRoaXMgaXMgZXZhbHVhdGVkIGF0IGJ1aWxkIHRpbWUgaW4gcGFja2FnZS5qc29uXG5pZiAoIURFQlVHKSB7XG4gICAgZGVsZXRlIHdpbmRvdy53YWlsc2JpbmRpbmdzO1xufVxuXG5sZXQgZHJhZ1Rlc3QgPSBmdW5jdGlvbihlKSB7XG4gICAgdmFyIHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnUHJvcGVydHkpO1xuICAgIGlmICh2YWwpIHtcbiAgICAgICAgdmFsID0gdmFsLnRyaW0oKTtcbiAgICB9XG5cbiAgICBpZiAodmFsICE9PSB3aW5kb3cud2FpbHMuZmxhZ3MuY3NzRHJhZ1ZhbHVlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICBpZiAoZS5idXR0b25zICE9PSAxKSB7XG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnZ2luZyBpZiBub3QgdGhlIHByaW1hcnkgYnV0dG9uIGhhcyBiZWVuIGNsaWNrZWQuXG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICBpZiAoZS5kZXRhaWwgIT09IDEpIHtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWdnaW5nIGlmIG1vcmUgdGhhbiBvbmNlIGhhcyBiZWVuIGNsaWNrZWQsIGUuZy4gd2hlbiBkb3VibGUgY2xpY2tpbmdcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIHJldHVybiB0cnVlO1xufTtcblxud2luZG93LndhaWxzLnNldENTU0RyYWdQcm9wZXJ0aWVzID0gZnVuY3Rpb24ocHJvcGVydHksIHZhbHVlKSB7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0RyYWdQcm9wZXJ0eSA9IHByb3BlcnR5O1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcmFnVmFsdWUgPSB2YWx1ZTtcbn1cblxud2luZG93LndhaWxzLnNldENTU0Ryb3BQcm9wZXJ0aWVzID0gZnVuY3Rpb24ocHJvcGVydHksIHZhbHVlKSB7XG4gICAgd2luZG93LndhaWxzLmZsYWdzLmNzc0Ryb3BQcm9wZXJ0eSA9IHByb3BlcnR5O1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5jc3NEcm9wVmFsdWUgPSB2YWx1ZTtcbn1cblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIChlKSA9PiB7XG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nXG4gICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKSB7XG4gICAgICAgIHdpbmRvdy5XYWlsc0ludm9rZShcInJlc2l6ZTpcIiArIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlKTtcbiAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKGRyYWdUZXN0KGUpKSB7XG4gICAgICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZVNjcm9sbGJhckRyYWcpIHtcbiAgICAgICAgICAgIC8vIFRoaXMgY2hlY2tzIGZvciBjbGlja3Mgb24gdGhlIHNjcm9sbCBiYXJcbiAgICAgICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICAgICAgaWYgKHdpbmRvdy53YWlscy5mbGFncy5kZWZlckRyYWdUb01vdXNlTW92ZSkge1xuICAgICAgICAgICAgd2luZG93LndhaWxzLmZsYWdzLnNob3VsZERyYWcgPSB0cnVlO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgZS5wcmV2ZW50RGVmYXVsdCgpXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJkcmFnXCIpO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybjtcbiAgICB9IGVsc2Uge1xuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgIH1cbn0pO1xuXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsICgpID0+IHtcbiAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xufSk7XG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IpIHtcbiAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gY3Vyc29yIHx8IHdpbmRvdy53YWlscy5mbGFncy5kZWZhdWx0Q3Vyc29yO1xuICAgIHdpbmRvdy53YWlscy5mbGFncy5yZXNpemVFZGdlID0gY3Vyc29yO1xufVxuXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgZnVuY3Rpb24oZSkge1xuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZykge1xuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3Muc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgICAgICBpZiAobW91c2VQcmVzc2VkID4gMCkge1xuICAgICAgICAgICAgd2luZG93LldhaWxzSW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoIXdpbmRvdy53YWlscy5mbGFncy5lbmFibGVSZXNpemUpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBpZiAod2luZG93LndhaWxzLmZsYWdzLmRlZmF1bHRDdXJzb3IgPT0gbnVsbCkge1xuICAgICAgICB3aW5kb3cud2FpbHMuZmxhZ3MuZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3I7XG4gICAgfVxuICAgIGlmICh3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHdpbmRvdy53YWlscy5mbGFncy5ib3JkZXJUaGlja25lc3MgJiYgd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcykge1xuICAgICAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gXCJzZS1yZXNpemVcIjtcbiAgICB9XG4gICAgbGV0IHJpZ2h0Qm9yZGVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCB3aW5kb3cud2FpbHMuZmxhZ3MuYm9yZGVyVGhpY2tuZXNzO1xuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgd2luZG93LndhaWxzLmZsYWdzLmJvcmRlclRoaWNrbmVzcztcblxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgd2luZG93LndhaWxzLmZsYWdzLnJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9IGVsc2UgaWYgKHJpZ2h0Qm9yZGVyICYmIGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIgJiYgYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlciAmJiB0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIgJiYgcmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xuXG59KTtcblxuLy8gU2V0dXAgY29udGV4dCBtZW51IGhvb2tcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGZ1bmN0aW9uKGUpIHtcbiAgICAvLyBhbHdheXMgc2hvdyB0aGUgY29udGV4dG1lbnUgaW4gZGVidWcgJiBkZXZcbiAgICBpZiAoREVCVUcpIHJldHVybjtcblxuICAgIGlmICh3aW5kb3cud2FpbHMuZmxhZ3MuZGlzYWJsZURlZmF1bHRDb250ZXh0TWVudSkge1xuICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgQ29udGV4dE1lbnUucHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShlKTtcbiAgICB9XG59KTtcblxud2luZG93LldhaWxzSW52b2tlKFwicnVudGltZTpyZWFkeVwiKTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWtCQSxXQUFTLGVBQWUsT0FBTyxTQUFTO0FBSXZDLFdBQU8sWUFBWSxNQUFNLFFBQVEsT0FBTztBQUFBLEVBQ3pDO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFNBQVMsU0FBUztBQUNqQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxRQUFRLFNBQVM7QUFDaEMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFdBQVcsU0FBUztBQUNuQyxtQkFBZSxLQUFLLE9BQU87QUFBQSxFQUM1QjtBQVFPLFdBQVMsU0FBUyxTQUFTO0FBQ2pDLG1CQUFlLEtBQUssT0FBTztBQUFBLEVBQzVCO0FBUU8sV0FBUyxTQUFTLFNBQVM7QUFDakMsbUJBQWUsS0FBSyxPQUFPO0FBQUEsRUFDNUI7QUFRTyxXQUFTLFlBQVksVUFBVTtB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
