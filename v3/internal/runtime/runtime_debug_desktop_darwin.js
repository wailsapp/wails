(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

  // desktop/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    SetText: () => SetText,
    Text: () => Text
  });

  // node_modules/nanoid/non-secure/index.js
  var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
  var nanoid = (size = 21) => {
    let id = "";
    let i = size;
    while (i--) {
      id += urlAlphabet[Math.random() * 64 | 0];
    }
    return id;
  };

  // desktop/runtime.js
  var runtimeURL = window.location.origin + "/wails/runtime";
  var objectNames = {
    Call: 0,
    Clipboard: 1,
    Application: 2,
    Events: 3,
    ContextMenu: 4,
    Dialog: 5,
    Window: 6,
    Screens: 7,
    System: 8,
    Browser: 9
  };
  var clientId = nanoid();
  function runtimeCallWithID(objectID, method, windowName, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID);
    url.searchParams.append("method", method);
    let fetchOptions = {
      headers: {}
    };
    if (windowName) {
      fetchOptions.headers["x-wails-window-name"] = windowName;
    }
    if (args) {
      url.searchParams.append("args", JSON.stringify(args));
    }
    fetchOptions.headers["x-wails-client-id"] = clientId;
    return new Promise((resolve, reject) => {
      fetch(url, fetchOptions).then((response) => {
        if (response.ok) {
          if (response.headers.get("Content-Type") && response.headers.get("Content-Type").indexOf("application/json") !== -1) {
            return response.json();
          } else {
            return response.text();
          }
        }
        reject(Error(response.statusText));
      }).then((data) => resolve(data)).catch((error) => reject(error));
    });
  }
  function newRuntimeCallerWithID(object, windowName) {
    return function(method, args = null) {
      return runtimeCallWithID(object, method, windowName, args);
    };
  }

  // desktop/clipboard.js
  var call = newRuntimeCallerWithID(objectNames.Clipboard);
  var ClipboardSetText = 0;
  var ClipboardText = 1;
  function SetText(text) {
    void call(ClipboardSetText, { text });
  }
  function Text() {
    return call(ClipboardText);
  }

  // desktop/application.js
  var application_exports = {};
  __export(application_exports, {
    Hide: () => Hide,
    Quit: () => Quit,
    Show: () => Show
  });
  var call2 = newRuntimeCallerWithID(objectNames.Application);
  var methods = {
    Hide: 0,
    Show: 1,
    Quit: 2
  };
  function Hide() {
    void call2(methods.Hide);
  }
  function Show() {
    void call2(methods.Show);
  }
  function Quit() {
    void call2(methods.Quit);
  }

  // desktop/screens.js
  var screens_exports = {};
  __export(screens_exports, {
    GetAll: () => GetAll,
    GetCurrent: () => GetCurrent,
    GetPrimary: () => GetPrimary
  });
  var call3 = newRuntimeCallerWithID(objectNames.Screens);
  var ScreensGetAll = 0;
  var ScreensGetPrimary = 1;
  var ScreensGetCurrent = 2;
  function GetAll() {
    return call3(ScreensGetAll);
  }
  function GetPrimary() {
    return call3(ScreensGetPrimary);
  }
  function GetCurrent() {
    return call3(ScreensGetCurrent);
  }

  // desktop/system.js
  var system_exports = {};
  __export(system_exports, {
    IsDarkMode: () => IsDarkMode
  });
  var call4 = newRuntimeCallerWithID(objectNames.System);
  var SystemIsDarkMode = 0;
  function IsDarkMode() {
    return call4(SystemIsDarkMode);
  }

  // desktop/browser.js
  var browser_exports = {};
  __export(browser_exports, {
    OpenURL: () => OpenURL
  });
  var call5 = newRuntimeCallerWithID(objectNames.Browser);
  var BrowserOpenURL = 0;
  function OpenURL(url) {
    void call5(BrowserOpenURL, { url });
  }

  // desktop/calls.js
  var call6 = newRuntimeCallerWithID(objectNames.Call);
  var CallBinding = 0;
  var callResponses = /* @__PURE__ */ new Map();
  function generateID() {
    let result;
    do {
      result = nanoid();
    } while (callResponses.has(result));
    return result;
  }
  function callCallback(id, data, isJSON) {
    let p = callResponses.get(id);
    if (p) {
      if (isJSON) {
        p.resolve(JSON.parse(data));
      } else {
        p.resolve(data);
      }
      callResponses.delete(id);
    }
  }
  function callErrorCallback(id, message) {
    let p = callResponses.get(id);
    if (p) {
      p.reject(message);
      callResponses.delete(id);
    }
  }
  function callBinding(type, options) {
    return new Promise((resolve, reject) => {
      let id = generateID();
      options = options || {};
      options["call-id"] = id;
      callResponses.set(id, { resolve, reject });
      call6(type, options).catch((error) => {
        reject(error);
        callResponses.delete(id);
      });
    });
  }
  function Call(options) {
    return callBinding(CallBinding, options);
  }
  function CallByName(name, ...args) {
    if (typeof name !== "string" || name.split(".").length !== 3) {
      throw new Error("CallByName requires a string in the format 'package.struct.method'");
    }
    let parts = name.split(".");
    return callBinding(CallBinding, {
      packageName: parts[0],
      structName: parts[1],
      methodName: parts[2],
      args
    });
  }
  function CallByID(methodID, ...args) {
    return callBinding(CallBinding, {
      methodID,
      args
    });
  }
  function Plugin(pluginName, methodName, ...args) {
    return callBinding(CallBinding, {
      packageName: "wails-plugins",
      structName: pluginName,
      methodName,
      args
    });
  }

  // desktop/window.js
  var WindowCenter = 0;
  var WindowSetTitle = 1;
  var WindowFullscreen = 2;
  var WindowUnFullscreen = 3;
  var WindowSetSize = 4;
  var WindowSize = 5;
  var WindowSetMaxSize = 6;
  var WindowSetMinSize = 7;
  var WindowSetAlwaysOnTop = 8;
  var WindowSetRelativePosition = 9;
  var WindowRelativePosition = 10;
  var WindowScreen = 11;
  var WindowHide = 12;
  var WindowMaximise = 13;
  var WindowUnMaximise = 14;
  var WindowToggleMaximise = 15;
  var WindowMinimise = 16;
  var WindowUnMinimise = 17;
  var WindowRestore = 18;
  var WindowShow = 19;
  var WindowClose = 20;
  var WindowSetBackgroundColour = 21;
  var WindowSetResizable = 22;
  var WindowWidth = 23;
  var WindowHeight = 24;
  var WindowZoomIn = 25;
  var WindowZoomOut = 26;
  var WindowZoomReset = 27;
  var WindowGetZoomLevel = 28;
  var WindowSetZoomLevel = 29;
  function newWindow(windowName) {
    let call10 = newRuntimeCallerWithID(objectNames.Window, windowName);
    return {
      /**
       * Centers the window.
       */
      Center: () => void call10(WindowCenter),
      /**
       * Set the window title.
       * @param title
       */
      SetTitle: (title) => void call10(WindowSetTitle, { title }),
      /**
       * Makes the window fullscreen.
       */
      Fullscreen: () => void call10(WindowFullscreen),
      /**
       * Unfullscreen the window.
       */
      UnFullscreen: () => void call10(WindowUnFullscreen),
      /**
       * Set the window size.
       * @param {number} width The window width
       * @param {number} height The window height
       */
      SetSize: (width, height) => call10(WindowSetSize, { width, height }),
      /**
       * Get the window size.
       * @returns {Promise<Size>} The window size
       */
      Size: () => {
        return call10(WindowSize);
      },
      /**
       * Set the window maximum size.
       * @param {number} width
       * @param {number} height
       */
      SetMaxSize: (width, height) => void call10(WindowSetMaxSize, { width, height }),
      /**
       * Set the window minimum size.
       * @param {number} width
       * @param {number} height
       */
      SetMinSize: (width, height) => void call10(WindowSetMinSize, { width, height }),
      /**
       * Set window to be always on top.
       * @param {boolean} onTop Whether the window should be always on top
       */
      SetAlwaysOnTop: (onTop) => void call10(WindowSetAlwaysOnTop, { alwaysOnTop: onTop }),
      /**
       * Set the window relative position.
       * @param {number} x
       * @param {number} y
       */
      SetRelativePosition: (x, y) => call10(WindowSetRelativePosition, { x, y }),
      /**
       * Get the window position.
       * @returns {Promise<Position>} The window position
       */
      RelativePosition: () => {
        return call10(WindowRelativePosition);
      },
      /**
       * Get the screen the window is on.
       * @returns {Promise<Screen>}
       */
      Screen: () => {
        return call10(WindowScreen);
      },
      /**
       * Hide the window
       */
      Hide: () => void call10(WindowHide),
      /**
       * Maximise the window
       */
      Maximise: () => void call10(WindowMaximise),
      /**
       * Show the window
       */
      Show: () => void call10(WindowShow),
      /**
       * Close the window
       */
      Close: () => void call10(WindowClose),
      /**
       * Toggle the window maximise state
       */
      ToggleMaximise: () => void call10(WindowToggleMaximise),
      /**
       * Unmaximise the window
       */
      UnMaximise: () => void call10(WindowUnMaximise),
      /**
       * Minimise the window
       */
      Minimise: () => void call10(WindowMinimise),
      /**
       * Unminimise the window
       */
      UnMinimise: () => void call10(WindowUnMinimise),
      /**
       * Restore the window
       */
      Restore: () => void call10(WindowRestore),
      /**
       * Set the background colour of the window.
       * @param {number} r - A value between 0 and 255
       * @param {number} g - A value between 0 and 255
       * @param {number} b - A value between 0 and 255
       * @param {number} a - A value between 0 and 255
       */
      SetBackgroundColour: (r, g, b, a) => void call10(WindowSetBackgroundColour, { r, g, b, a }),
      /**
       * Set whether the window can be resized or not
       * @param {boolean} resizable
       */
      SetResizable: (resizable) => void call10(WindowSetResizable, { resizable }),
      /**
       * Get the window width
       * @returns {Promise<number>}
       */
      Width: () => {
        return call10(WindowWidth);
      },
      /**
       * Get the window height
       * @returns {Promise<number>}
       */
      Height: () => {
        return call10(WindowHeight);
      },
      /**
       * Zoom in the window
       */
      ZoomIn: () => void call10(WindowZoomIn),
      /**
       * Zoom out the window
       */
      ZoomOut: () => void call10(WindowZoomOut),
      /**
       * Reset the window zoom
       */
      ZoomReset: () => void call10(WindowZoomReset),
      /**
       * Get the window zoom
       * @returns {Promise<number>}
       */
      GetZoomLevel: () => {
        return call10(WindowGetZoomLevel);
      },
      /**
       * Set the window zoom level
       * @param {number} zoomLevel
       */
      SetZoomLevel: (zoomLevel) => void call10(WindowSetZoomLevel, { zoomLevel })
    };
  }

  // desktop/events.js
  var call7 = newRuntimeCallerWithID(objectNames.Events);
  var EventEmit = 0;
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
        callback(data);
        if (this.maxCallbacks === -1) {
          return false;
        }
        this.maxCallbacks -= 1;
        return this.maxCallbacks === 0;
      };
    }
  };
  var WailsEvent = class {
    /**
     * Creates an instance of WailsEvent.
     * @param {string} name - Name of the event
     * @param {any=null} data - Data associated with the event
     * @memberof WailsEvent
     */
    constructor(name, data = null) {
      this.name = name;
      this.data = data;
    }
  };
  var eventListeners = /* @__PURE__ */ new Map();
  function OnMultiple(eventName, callback, maxCallbacks) {
    let listeners = eventListeners.get(eventName) || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    listeners.push(thisListener);
    eventListeners.set(eventName, listeners);
    return () => listenerOff(thisListener);
  }
  function On(eventName, callback) {
    return OnMultiple(eventName, callback, -1);
  }
  function Once(eventName, callback) {
    return OnMultiple(eventName, callback, 1);
  }
  function listenerOff(listener) {
    const eventName = listener.eventName;
    let listeners = eventListeners.get(eventName).filter((l) => l !== listener);
    if (listeners.length === 0) {
      eventListeners.delete(eventName);
    } else {
      eventListeners.set(eventName, listeners);
    }
  }
  function dispatchWailsEvent(event) {
    let listeners = eventListeners.get(event.name);
    if (listeners) {
      let toRemove = [];
      listeners.forEach((listener) => {
        let remove = listener.Callback(event);
        if (remove) {
          toRemove.push(listener);
        }
      });
      if (toRemove.length > 0) {
        listeners = listeners.filter((l) => !toRemove.includes(l));
        if (listeners.length === 0) {
          eventListeners.delete(event.name);
        } else {
          eventListeners.set(event.name, listeners);
        }
      }
    }
  }
  function Off(eventName, ...additionalEventNames) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach((eventName2) => {
      eventListeners.delete(eventName2);
    });
  }
  function OffAll() {
    eventListeners.clear();
  }
  function Emit(event) {
    void call7(EventEmit, event);
  }

  // desktop/dialogs.js
  var call8 = newRuntimeCallerWithID(objectNames.Dialog);
  var DialogInfo = 0;
  var DialogWarning = 1;
  var DialogError = 2;
  var DialogQuestion = 3;
  var DialogOpenFile = 4;
  var DialogSaveFile = 5;
  var dialogResponses = /* @__PURE__ */ new Map();
  function generateID2() {
    let result;
    do {
      result = nanoid();
    } while (dialogResponses.has(result));
    return result;
  }
  function dialogCallback(id, data, isJSON) {
    let p = dialogResponses.get(id);
    if (p) {
      if (isJSON) {
        p.resolve(JSON.parse(data));
      } else {
        p.resolve(data);
      }
      dialogResponses.delete(id);
    }
  }
  function dialogErrorCallback(id, message) {
    let p = dialogResponses.get(id);
    if (p) {
      p.reject(message);
      dialogResponses.delete(id);
    }
  }
  function dialog(type, options) {
    return new Promise((resolve, reject) => {
      let id = generateID2();
      options = options || {};
      options["dialog-id"] = id;
      dialogResponses.set(id, { resolve, reject });
      call8(type, options).catch((error) => {
        reject(error);
        dialogResponses.delete(id);
      });
    });
  }
  function Info(options) {
    return dialog(DialogInfo, options);
  }
  function Warning(options) {
    return dialog(DialogWarning, options);
  }
  function Error2(options) {
    return dialog(DialogError, options);
  }
  function Question(options) {
    return dialog(DialogQuestion, options);
  }
  function OpenFile(options) {
    return dialog(DialogOpenFile, options);
  }
  function SaveFile(options) {
    return dialog(DialogSaveFile, options);
  }

  // desktop/contextmenu.js
  var call9 = newRuntimeCallerWithID(objectNames.ContextMenu);
  var ContextMenuOpen = 0;
  function openContextMenu(id, x, y, data) {
    void call9(ContextMenuOpen, { id, x, y, data });
  }
  function setupContextMenus() {
    window.addEventListener("contextmenu", contextMenuHandler);
  }
  function contextMenuHandler(event) {
    let element = event.target;
    let customContextMenu = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu");
    customContextMenu = customContextMenu ? customContextMenu.trim() : "";
    if (customContextMenu) {
      event.preventDefault();
      let customContextMenuData = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu-data");
      openContextMenu(customContextMenu, event.clientX, event.clientY, customContextMenuData);
      return;
    }
    processDefaultContextMenu(event);
  }
  function processDefaultContextMenu(event) {
    if (true) {
      return;
    }
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

  // desktop/wml.js
  function sendEvent(eventName, data = null) {
    let event = new WailsEvent(eventName, data);
    Emit(event);
  }
  function addWMLEventListeners() {
    const elements = document.querySelectorAll("[wml-event]");
    elements.forEach(function(element) {
      const eventType = element.getAttribute("wml-event");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Detached: false, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              sendEvent(eventType);
            }
          });
          return;
        }
        sendEvent(eventType);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function callWindowMethod(method) {
    if (wails.Window[method] === void 0) {
      console.log("Window method " + method + " not found");
    }
    wails.Window[method]();
  }
  function addWMLWindowListeners() {
    const elements = document.querySelectorAll("[wml-window]");
    elements.forEach(function(element) {
      const windowMethod = element.getAttribute("wml-window");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              callWindowMethod(windowMethod);
            }
          });
          return;
        }
        callWindowMethod(windowMethod);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function addWMLOpenBrowserListener() {
    const elements = document.querySelectorAll("[wml-openurl]");
    elements.forEach(function(element) {
      const url = element.getAttribute("wml-openurl");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              void wails.Browser.OpenURL(url);
            }
          });
          return;
        }
        void wails.Browser.OpenURL(url);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function reloadWML() {
    addWMLEventListeners();
    addWMLWindowListeners();
    addWMLOpenBrowserListener();
  }

  // desktop/invoke.js
  var invoke = function(input) {
    if (false) {
      chrome.webview.postMessage(input);
    } else {
      webkit.messageHandlers.external.postMessage(input);
    }
  };

  // desktop/flags.js
  var flags = /* @__PURE__ */ new Map();
  function convertToMap(obj) {
    const map = /* @__PURE__ */ new Map();
    for (const [key, value] of Object.entries(obj)) {
      if (typeof value === "object" && value !== null) {
        map.set(key, convertToMap(value));
      } else {
        map.set(key, value);
      }
    }
    return map;
  }
  fetch("/wails/flags").then((response) => {
    response.json().then((data) => {
      flags = convertToMap(data);
    });
  });

  // desktop/drag.js
  var shouldDrag = false;
  function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--webkit-app-region");
    if (val) {
      val = val.trim();
    }
    if (val !== "drag") {
      return false;
    }
    if (e.buttons !== 1) {
      return false;
    }
    return e.detail === 1;
  }
  function setupDrag() {
    window.addEventListener("mousedown", onMouseDown);
    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
  }
  function onMouseDown(e) {
    if (false) {
      if (testResize()) {
        return;
      }
    }
    if (dragTest(e)) {
      if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
        return;
      }
      shouldDrag = true;
    } else {
      shouldDrag = false;
    }
  }
  function onMouseUp(e) {
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (mousePressed > 0) {
      endDrag();
    }
  }
  function endDrag() {
    document.body.style.cursor = "default";
    shouldDrag = false;
  }
  function onMouseMove(e) {
    if (shouldDrag) {
      shouldDrag = false;
      let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
      if (mousePressed > 0) {
        invoke("drag");
      }
      return;
    }
    if (false) {
      handleResize(e);
    }
  }

  // desktop/main.js
  window.wails = {
    ...newRuntime(null),
    Capabilities: {},
    clientId
  };
  fetch("/wails/capabilities").then((response) => {
    response.json().then((data) => {
      window.wails.Capabilities = data;
    });
  });
  window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchWailsEvent,
    callCallback,
    callErrorCallback,
    endDrag
  };
  function newRuntime(windowName) {
    return {
      Clipboard: {
        ...clipboard_exports
      },
      Application: {
        ...application_exports,
        GetWindowByName(windowName2) {
          return newRuntime(windowName2);
        }
      },
      System: system_exports,
      Screens: screens_exports,
      Browser: browser_exports,
      Call,
      CallByID,
      CallByName,
      Plugin,
      WML: {
        Reload: reloadWML
      },
      Dialog: {
        Info,
        Warning,
        Error: Error2,
        Question,
        OpenFile,
        SaveFile
      },
      Events: {
        Emit,
        On,
        Once,
        OnMultiple,
        Off,
        OffAll
      },
      Window: newWindow(windowName)
    };
  }
  if (true) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
  }
  setupContextMenus();
  setupDrag();
  document.addEventListener("DOMContentLoaded", function() {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAibm9kZV9tb2R1bGVzL25hbm9pZC9ub24tc2VjdXJlL2luZGV4LmpzIiwgImRlc2t0b3AvcnVudGltZS5qcyIsICJkZXNrdG9wL2FwcGxpY2F0aW9uLmpzIiwgImRlc2t0b3Avc2NyZWVucy5qcyIsICJkZXNrdG9wL3N5c3RlbS5qcyIsICJkZXNrdG9wL2Jyb3dzZXIuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL3dpbmRvdy5qcyIsICJkZXNrdG9wL2V2ZW50cy5qcyIsICJkZXNrdG9wL2RpYWxvZ3MuanMiLCAiZGVza3RvcC9jb250ZXh0bWVudS5qcyIsICJkZXNrdG9wL3dtbC5qcyIsICJkZXNrdG9wL2ludm9rZS5qcyIsICJkZXNrdG9wL2ZsYWdzLmpzIiwgImRlc2t0b3AvZHJhZy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DbGlwYm9hcmQpO1xyXG5cclxubGV0IENsaXBib2FyZFNldFRleHQgPSAwO1xyXG5sZXQgQ2xpcGJvYXJkVGV4dCA9IDE7XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xyXG4gICAgdm9pZCBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xyXG59IiwgImxldCB1cmxBbHBoYWJldCA9XG4gICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuZXhwb3J0IGxldCBjdXN0b21BbHBoYWJldCA9IChhbHBoYWJldCwgZGVmYXVsdFNpemUgPSAyMSkgPT4ge1xuICByZXR1cm4gKHNpemUgPSBkZWZhdWx0U2l6ZSkgPT4ge1xuICAgIGxldCBpZCA9ICcnXG4gICAgbGV0IGkgPSBzaXplXG4gICAgd2hpbGUgKGktLSkge1xuICAgICAgaWQgKz0gYWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiBhbHBoYWJldC5sZW5ndGgpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG4gIH1cbn1cbmV4cG9ydCBsZXQgbmFub2lkID0gKHNpemUgPSAyMSkgPT4ge1xuICBsZXQgaWQgPSAnJ1xuICBsZXQgaSA9IHNpemVcbiAgd2hpbGUgKGktLSkge1xuICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgfVxuICByZXR1cm4gaWRcbn1cbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xyXG4vLyBPYmplY3QgTmFtZXNcclxuZXhwb3J0IGNvbnN0IG9iamVjdE5hbWVzID0ge1xyXG4gICAgQ2FsbDogMCxcclxuICAgIENsaXBib2FyZDogMSxcclxuICAgIEFwcGxpY2F0aW9uOiAyLFxyXG4gICAgRXZlbnRzOiAzLFxyXG4gICAgQ29udGV4dE1lbnU6IDQsXHJcbiAgICBEaWFsb2c6IDUsXHJcbiAgICBXaW5kb3c6IDYsXHJcbiAgICBTY3JlZW5zOiA3LFxyXG4gICAgU3lzdGVtOiA4LFxyXG4gICAgQnJvd3NlcjogOSxcclxufVxyXG5leHBvcnQgbGV0IGNsaWVudElkID0gbmFub2lkKCk7XHJcblxyXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgaWYoIG1ldGhvZCApIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xyXG4gICAgfVxyXG4gICAgbGV0IGZldGNoT3B0aW9ucyA9IHtcclxuICAgICAgICBoZWFkZXJzOiB7fSxcclxuICAgIH07XHJcbiAgICBpZiAod2luZG93TmFtZSkge1xyXG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XHJcbiAgICB9XHJcbiAgICBpZiAoYXJncykge1xyXG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XHJcbiAgICB9XHJcbiAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtY2xpZW50LWlkXCJdID0gY2xpZW50SWQ7XHJcblxyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0LCB3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0LCB3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5BcHBsaWNhdGlvbik7XHJcblxyXG5sZXQgbWV0aG9kcyA9IHtcclxuICAgIEhpZGU6IDAsXHJcbiAgICBTaG93OiAxLFxyXG4gICAgUXVpdDogMixcclxufVxyXG5cclxuLyoqXHJcbiAqIEhpZGUgdGhlIGFwcGxpY2F0aW9uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcclxuICAgIHZvaWQgY2FsbChtZXRob2RzLkhpZGUpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xyXG4gICAgdm9pZCBjYWxsKG1ldGhvZHMuU2hvdyk7XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogUXVpdCB0aGUgYXBwbGljYXRpb25cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xyXG4gICAgdm9pZCBjYWxsKG1ldGhvZHMuUXVpdCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5TY3JlZW59IFNjcmVlblxyXG4gKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TY3JlZW5zKTtcclxuXHJcbmxldCBTY3JlZW5zR2V0QWxsID0gMDtcclxubGV0IFNjcmVlbnNHZXRQcmltYXJ5ID0gMTtcclxubGV0IFNjcmVlbnNHZXRDdXJyZW50ID0gMjtcclxuXHJcbi8qKlxyXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW5bXT59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoU2NyZWVuc0dldEFsbCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cclxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoU2NyZWVuc0dldFByaW1hcnkpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxyXG4gKiBAY29uc3RydWN0b3JcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoU2NyZWVuc0dldEN1cnJlbnQpO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TeXN0ZW0pO1xyXG5cclxubGV0IFN5c3RlbUlzRGFya01vZGUgPSAwO1xyXG5cclxuLyoqXHJcbiAqIERldGVybWluZXMgaWYgdGhlIHN5c3RlbSBpcyBjdXJyZW50bHkgdXNpbmcgZGFyayBtb2RlXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPGJvb2xlYW4+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKSB7XHJcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1Jc0RhcmtNb2RlKTtcclxufSIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQnJvd3Nlcik7XHJcblxyXG5sZXQgQnJvd3Nlck9wZW5VUkwgPSAwO1xyXG5cclxuLyoqXHJcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybCkge1xyXG4gICAgdm9pZCBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsfSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYWxsKTtcclxuXHJcbmxldCBDYWxsQmluZGluZyA9IDA7XHJcblxyXG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcclxuXHJcbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XHJcbiAgICBsZXQgcmVzdWx0O1xyXG4gICAgZG8ge1xyXG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xyXG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gY2FsbENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcclxuICAgIGxldCBwID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgaWYgKHApIHtcclxuICAgICAgICBpZiAoaXNKU09OKSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XHJcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGlmIChwKSB7XHJcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBjYWxsQmluZGluZyh0eXBlLCBvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcclxuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xyXG5cclxuICAgICAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIG9wdGlvbnMpO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbEJ5TmFtZShuYW1lLCAuLi5hcmdzKSB7XHJcblxyXG4gICAgLy8gRW5zdXJlIGZpcnN0IGFyZ3VtZW50IGlzIGEgc3RyaW5nIGFuZCBoYXMgMiBkb3RzXHJcbiAgICBpZiAodHlwZW9mIG5hbWUgIT09IFwic3RyaW5nXCIgfHwgbmFtZS5zcGxpdChcIi5cIikubGVuZ3RoICE9PSAzKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiQ2FsbEJ5TmFtZSByZXF1aXJlcyBhIHN0cmluZyBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnXCIpO1xyXG4gICAgfVxyXG4gICAgLy8gU3BsaXQgaW5wdXRzXHJcbiAgICBsZXQgcGFydHMgPSBuYW1lLnNwbGl0KFwiLlwiKTtcclxuXHJcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIHtcclxuICAgICAgICBwYWNrYWdlTmFtZTogcGFydHNbMF0sXHJcbiAgICAgICAgc3RydWN0TmFtZTogcGFydHNbMV0sXHJcbiAgICAgICAgbWV0aG9kTmFtZTogcGFydHNbMl0sXHJcbiAgICAgICAgYXJnczogYXJncyxcclxuICAgIH0pO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbEJ5SUQobWV0aG9kSUQsIC4uLmFyZ3MpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xyXG4gICAgICAgIG1ldGhvZElEOiBtZXRob2RJRCxcclxuICAgICAgICBhcmdzOiBhcmdzLFxyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsIGEgcGx1Z2luIG1ldGhvZFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gcGx1Z2luTmFtZSAtIG5hbWUgb2YgdGhlIHBsdWdpblxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWV0aG9kTmFtZSAtIG5hbWUgb2YgdGhlIG1ldGhvZFxyXG4gKiBAcGFyYW0gey4uLmFueX0gYXJncyAtIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2RcclxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gLSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzdWx0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUGx1Z2luKHBsdWdpbk5hbWUsIG1ldGhvZE5hbWUsIC4uLmFyZ3MpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xyXG4gICAgICAgIHBhY2thZ2VOYW1lOiBcIndhaWxzLXBsdWdpbnNcIixcclxuICAgICAgICBzdHJ1Y3ROYW1lOiBwbHVnaW5OYW1lLFxyXG4gICAgICAgIG1ldGhvZE5hbWU6IG1ldGhvZE5hbWUsXHJcbiAgICAgICAgYXJnczogYXJncyxcclxuICAgIH0pO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNpemV9IFNpemVcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5Qb3NpdGlvbn0gUG9zaXRpb25cclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5TY3JlZW59IFNjcmVlblxyXG4gKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmxldCBXaW5kb3dDZW50ZXIgPSAwO1xyXG5sZXQgV2luZG93U2V0VGl0bGUgPSAxO1xyXG5sZXQgV2luZG93RnVsbHNjcmVlbiA9IDI7XHJcbmxldCBXaW5kb3dVbkZ1bGxzY3JlZW4gPSAzO1xyXG5sZXQgV2luZG93U2V0U2l6ZSA9IDQ7XHJcbmxldCBXaW5kb3dTaXplID0gNTtcclxubGV0IFdpbmRvd1NldE1heFNpemUgPSA2O1xyXG5sZXQgV2luZG93U2V0TWluU2l6ZSA9IDc7XHJcbmxldCBXaW5kb3dTZXRBbHdheXNPblRvcCA9IDg7XHJcbmxldCBXaW5kb3dTZXRSZWxhdGl2ZVBvc2l0aW9uID0gOTtcclxubGV0IFdpbmRvd1JlbGF0aXZlUG9zaXRpb24gPSAxMDtcclxubGV0IFdpbmRvd1NjcmVlbiA9IDExO1xyXG5sZXQgV2luZG93SGlkZSA9IDEyO1xyXG5sZXQgV2luZG93TWF4aW1pc2UgPSAxMztcclxubGV0IFdpbmRvd1VuTWF4aW1pc2UgPSAxNDtcclxubGV0IFdpbmRvd1RvZ2dsZU1heGltaXNlID0gMTU7XHJcbmxldCBXaW5kb3dNaW5pbWlzZSA9IDE2O1xyXG5sZXQgV2luZG93VW5NaW5pbWlzZSA9IDE3O1xyXG5sZXQgV2luZG93UmVzdG9yZSA9IDE4O1xyXG5sZXQgV2luZG93U2hvdyA9IDE5O1xyXG5sZXQgV2luZG93Q2xvc2UgPSAyMDtcclxubGV0IFdpbmRvd1NldEJhY2tncm91bmRDb2xvdXIgPSAyMTtcclxubGV0IFdpbmRvd1NldFJlc2l6YWJsZSA9IDIyO1xyXG5sZXQgV2luZG93V2lkdGggPSAyMztcclxubGV0IFdpbmRvd0hlaWdodCA9IDI0O1xyXG5sZXQgV2luZG93Wm9vbUluID0gMjU7XHJcbmxldCBXaW5kb3dab29tT3V0ID0gMjY7XHJcbmxldCBXaW5kb3dab29tUmVzZXQgPSAyNztcclxubGV0IFdpbmRvd0dldFpvb21MZXZlbCA9IDI4O1xyXG5sZXQgV2luZG93U2V0Wm9vbUxldmVsID0gMjk7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gbmV3V2luZG93KHdpbmRvd05hbWUpIHtcclxuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIHdpbmRvd05hbWUpO1xyXG4gICAgcmV0dXJuIHtcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogQ2VudGVycyB0aGUgd2luZG93LlxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIENlbnRlcjogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd0NlbnRlciksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHRpdGxlLlxyXG4gICAgICAgICAqIEBwYXJhbSB0aXRsZVxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IHZvaWQgY2FsbChXaW5kb3dTZXRUaXRsZSwge3RpdGxlfSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIE1ha2VzIHRoZSB3aW5kb3cgZnVsbHNjcmVlbi5cclxuICAgICAgICAgKi9cclxuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoV2luZG93RnVsbHNjcmVlbiksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFVuZnVsbHNjcmVlbiB0aGUgd2luZG93LlxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFVuRnVsbHNjcmVlbjogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1VuRnVsbHNjcmVlbiksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHNpemUuXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIFRoZSB3aW5kb3cgd2lkdGhcclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IFRoZSB3aW5kb3cgaGVpZ2h0XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoV2luZG93U2V0U2l6ZSwge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBzaXplLlxyXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNpemU+fSBUaGUgd2luZG93IHNpemVcclxuICAgICAgICAgKi9cclxuICAgICAgICBTaXplOiAoKSA9PiB7IHJldHVybiBjYWxsKFdpbmRvd1NpemUpOyB9LFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBtYXhpbXVtIHNpemUuXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoV2luZG93U2V0TWF4U2l6ZSwge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBtaW5pbXVtIHNpemUuXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNldE1pblNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoV2luZG93U2V0TWluU2l6ZSwge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXHJcbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSBvblRvcCBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGFsd2F5cyBvbiB0b3BcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRBbHdheXNPblRvcDogKG9uVG9wKSA9PiB2b2lkIGNhbGwoV2luZG93U2V0QWx3YXlzT25Ub3AsIHthbHdheXNPblRvcDpvblRvcH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyByZWxhdGl2ZSBwb3NpdGlvbi5cclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0geFxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB5XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0UmVsYXRpdmVQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoV2luZG93U2V0UmVsYXRpdmVQb3NpdGlvbiwge3gseX0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBwb3NpdGlvbi5cclxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxQb3NpdGlvbj59IFRoZSB3aW5kb3cgcG9zaXRpb25cclxuICAgICAgICAgKi9cclxuICAgICAgICBSZWxhdGl2ZVBvc2l0aW9uOiAoKSA9PiB7IHJldHVybiBjYWxsKFdpbmRvd1JlbGF0aXZlUG9zaXRpb24pOyB9LFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIG9uLlxyXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2NyZWVuOiAoKSA9PiB7IHJldHVybiBjYWxsKFdpbmRvd1NjcmVlbik7IH0sXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIEhpZGUgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIEhpZGU6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dIaWRlKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogTWF4aW1pc2UgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIE1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoV2luZG93TWF4aW1pc2UpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTaG93IHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBTaG93OiAoKSA9PiB2b2lkIGNhbGwoV2luZG93U2hvdyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIENsb3NlIHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBDbG9zZTogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd0Nsb3NlKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogVG9nZ2xlIHRoZSB3aW5kb3cgbWF4aW1pc2Ugc3RhdGVcclxuICAgICAgICAgKi9cclxuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1RvZ2dsZU1heGltaXNlKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogVW5tYXhpbWlzZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgVW5NYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1VuTWF4aW1pc2UpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBNaW5pbWlzZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dNaW5pbWlzZSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFVubWluaW1pc2UgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dVbk1pbmltaXNlKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogUmVzdG9yZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgUmVzdG9yZTogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1Jlc3RvcmUpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGcgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGEgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IHZvaWQgY2FsbChXaW5kb3dTZXRCYWNrZ3JvdW5kQ29sb3VyLCB7ciwgZywgYiwgYX0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgd2hldGhlciB0aGUgd2luZG93IGNhbiBiZSByZXNpemVkIG9yIG5vdFxyXG4gICAgICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gcmVzaXphYmxlXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0UmVzaXphYmxlOiAocmVzaXphYmxlKSA9PiB2b2lkIGNhbGwoV2luZG93U2V0UmVzaXphYmxlLCB7cmVzaXphYmxlfSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IHdpZHRoXHJcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8bnVtYmVyPn1cclxuICAgICAgICAgKi9cclxuICAgICAgICBXaWR0aDogKCkgPT4geyByZXR1cm4gY2FsbChXaW5kb3dXaWR0aCk7IH0sXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IGhlaWdodFxyXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPG51bWJlcj59XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgSGVpZ2h0OiAoKSA9PiB7IHJldHVybiBjYWxsKFdpbmRvd0hlaWdodCk7IH0sXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFpvb20gaW4gdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFpvb21JbjogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1pvb21JbiksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFpvb20gb3V0IHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBab29tT3V0OiAoKSA9PiB2b2lkIGNhbGwoV2luZG93Wm9vbU91dCksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFJlc2V0IHRoZSB3aW5kb3cgem9vbVxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFpvb21SZXNldDogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1pvb21SZXNldCksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IHpvb21cclxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxudW1iZXI+fVxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIEdldFpvb21MZXZlbDogKCkgPT4geyByZXR1cm4gY2FsbChXaW5kb3dHZXRab29tTGV2ZWwpOyB9LFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyB6b29tIGxldmVsXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHpvb21MZXZlbFxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNldFpvb21MZXZlbDogKHpvb21MZXZlbCkgPT4gdm9pZCBjYWxsKFdpbmRvd1NldFpvb21MZXZlbCwge3pvb21MZXZlbH0pLFxyXG4gICAgfTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5XYWlsc0V2ZW50fSBXYWlsc0V2ZW50XHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkV2ZW50cyk7XHJcbmxldCBFdmVudEVtaXQgPSAwO1xyXG5cclxuLyoqXHJcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxyXG4gKlxyXG4gKiBAY2xhc3MgTGlzdGVuZXJcclxuICovXHJcbmNsYXNzIExpc3RlbmVyIHtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cclxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XHJcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XHJcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxyXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xyXG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXHJcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XHJcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xyXG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXHJcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXHJcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XHJcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcclxuICAgICAgICB9O1xyXG4gICAgfVxyXG59XHJcblxyXG5cclxuLyoqXHJcbiAqIFdhaWxzRXZlbnQgZGVmaW5lcyBhIGN1c3RvbSBldmVudC4gSXQgaXMgcGFzc2VkIHRvIGV2ZW50IGxpc3RlbmVycy5cclxuICpcclxuICogQGNsYXNzIFdhaWxzRXZlbnRcclxuICogQHByb3BlcnR5IHtzdHJpbmd9IG5hbWUgLSBOYW1lIG9mIHRoZSBldmVudFxyXG4gKiBAcHJvcGVydHkge2FueX0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQge1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIFdhaWxzRXZlbnQuXHJcbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XHJcbiAgICAgKiBAcGFyYW0ge2FueT1udWxsfSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XHJcbiAgICAgKiBAbWVtYmVyb2YgV2FpbHNFdmVudFxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhID0gbnVsbCkge1xyXG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XHJcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xyXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcclxuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xyXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xyXG4gQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xyXG59XHJcblxyXG4vKipcclxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxyXG4gKlxyXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxyXG4gKi9cclxuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcclxuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcclxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBkaXNwYXRjaGVzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50KSB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xyXG4gICAgaWYgKGxpc3RlbmVycykge1xyXG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxyXG4gICAgICAgIGxldCB0b1JlbW92ZSA9IFtdO1xyXG4gICAgICAgIGxpc3RlbmVycy5mb3JFYWNoKGxpc3RlbmVyID0+IHtcclxuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcclxuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xyXG4gICAgICAgICAgICAgICAgdG9SZW1vdmUucHVzaChsaXN0ZW5lcik7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgICAgICAvLyByZW1vdmUgbGlzdGVuZXJzXHJcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XHJcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XHJcbiAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxyXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXHJcbiAqXHJcbiBbdjMgQ0hBTkdFXSBPZmYgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xyXG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xyXG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBPZmZBbGwgdW5yZWdpc3RlcnMgYWxsIGxpc3RlbmVyc1xyXG4gKiBbdjMgQ0hBTkdFXSBPZmZBbGwgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xyXG4gKlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50XHJcbiAqIEBwYXJhbSB7V2FpbHNFdmVudH0gZXZlbnQgVGhlIGV2ZW50IHRvIGVtaXRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XHJcbiAgICB2b2lkIGNhbGwoRXZlbnRFbWl0LCBldmVudCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5NZXNzYWdlRGlhbG9nT3B0aW9uc30gTWVzc2FnZURpYWxvZ09wdGlvbnNcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk9wZW5EaWFsb2dPcHRpb25zfSBPcGVuRGlhbG9nT3B0aW9uc1xyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuU2F2ZURpYWxvZ09wdGlvbnN9IFNhdmVEaWFsb2dPcHRpb25zXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZyk7XHJcblxyXG5sZXQgRGlhbG9nSW5mbyA9IDA7XHJcbmxldCBEaWFsb2dXYXJuaW5nID0gMTtcclxubGV0IERpYWxvZ0Vycm9yID0gMjtcclxubGV0IERpYWxvZ1F1ZXN0aW9uID0gMztcclxubGV0IERpYWxvZ09wZW5GaWxlID0gNDtcclxubGV0IERpYWxvZ1NhdmVGaWxlID0gNTtcclxuXHJcblxyXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0NhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xyXG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucykge1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XHJcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XHJcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xyXG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcclxuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xyXG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xyXG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGFuIEluZm8gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGEgV2FybmluZyBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgYW4gRXJyb3IgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgYSBRdWVzdGlvbiBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhbiBPcGVuIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZ1tdfHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUgb3IgYW4gYXJyYXkgb2Ygc2VsZWN0ZWQgZmlsZXMgaWYgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24gaXMgdHJ1ZS4gQSBibGFuayBzdHJpbmcgaXMgcmV0dXJuZWQgaWYgbm8gZmlsZSB3YXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhEaWFsb2dPcGVuRmlsZSwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhIFNhdmUgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7U2F2ZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZS4gQSBibGFuayBzdHJpbmcgaXMgcmV0dXJuZWQgaWYgbm8gZmlsZSB3YXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbiIsICJpbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xyXG5cclxubGV0IENvbnRleHRNZW51T3BlbiA9IDA7XHJcblxyXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQsIHgsIHksIGRhdGEpIHtcclxuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gc2V0dXBDb250ZXh0TWVudXMoKSB7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcclxuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XHJcbiAgICBsZXQgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcclxuICAgIGxldCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudVwiKTtcclxuICAgIGN1c3RvbUNvbnRleHRNZW51ID0gY3VzdG9tQ29udGV4dE1lbnUgPyBjdXN0b21Db250ZXh0TWVudS50cmltKCkgOiBcIlwiO1xyXG4gICAgaWYgKGN1c3RvbUNvbnRleHRNZW51KSB7XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnVEYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XHJcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBjdXN0b21Db250ZXh0TWVudURhdGEpO1xyXG4gICAgICAgIHJldHVyblxyXG4gICAgfVxyXG5cclxuICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpO1xyXG59XHJcblxyXG5cclxuLypcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBhdXRvOyAoZGVmYXVsdCkgd2lsbCBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudSBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZSBPUiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkIE9SIGVsZW1lbnQgaXMgaW5wdXQgb3IgdGV4dGFyZWFcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcblxyXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxyXG4qL1xyXG5mdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KSB7XHJcbiAgICAvLyBEZWJ1ZyBidWlsZHMgYWx3YXlzIHNob3cgdGhlIG1lbnVcclxuICAgIGlmIChERUJVRykge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICBjb25zdCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xyXG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xyXG4gICAgY29uc3QgZGVmYXVsdENvbnRleHRNZW51QWN0aW9uID0gY29tcHV0ZWRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKTtcclxuICAgIHN3aXRjaCAoZGVmYXVsdENvbnRleHRNZW51QWN0aW9uKSB7XHJcbiAgICAgICAgY2FzZSBcInNob3dcIjpcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGNhc2UgXCJoaWRlXCI6XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICBkZWZhdWx0OlxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxyXG4gICAgICAgICAgICBpZiAoZWxlbWVudC5pc0NvbnRlbnRFZGl0YWJsZSkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXHJcbiAgICAgICAgICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcclxuICAgICAgICAgICAgY29uc3QgaGFzU2VsZWN0aW9uID0gKHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDApXHJcbiAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcclxuICAgICAgICAgICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xyXG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xyXG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gZWxlbWVudCkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRhZ25hbWUgaXMgaW5wdXQgb3IgdGV4dGFyZWFcclxuICAgICAgICAgICAgaWYgKGVsZW1lbnQudGFnTmFtZSA9PT0gXCJJTlBVVFwiIHx8IGVsZW1lbnQudGFnTmFtZSA9PT0gXCJURVhUQVJFQVwiKSB7XHJcbiAgICAgICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghZWxlbWVudC5yZWFkT25seSAmJiAhZWxlbWVudC5kaXNhYmxlZCkpIHtcclxuICAgICAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgIH1cclxufVxyXG4iLCAiXHJcbmltcG9ydCB7RW1pdCwgV2FpbHNFdmVudH0gZnJvbSBcIi4vZXZlbnRzXCI7XHJcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcclxuXHJcbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xyXG4gICAgbGV0IGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKTtcclxuICAgIEVtaXQoZXZlbnQpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcclxuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0nKTtcclxuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcclxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50Jyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgRGV0YWNoZWQ6IGZhbHNlLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICB9O1xyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuXHJcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuXHJcbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcclxuICAgIGlmICh3YWlscy5XaW5kb3dbbWV0aG9kXSA9PT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xyXG4gICAgfVxyXG4gICAgd2FpbHMuV2luZG93W21ldGhvZF0oKTtcclxufVxyXG5cclxuZnVuY3Rpb24gYWRkV01MV2luZG93TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLXdpbmRvd10nKTtcclxuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcclxuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpO1xyXG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcclxuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xyXG5cclxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XHJcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xyXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcblxyXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGFkZFdNTE9wZW5Ccm93c2VyTGlzdGVuZXIoKSB7XHJcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtYnJvd3Nlci1vcGVudXJsXScpO1xyXG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xyXG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtYnJvd3Nlci1vcGVudXJsJyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCB3YWlscy5Ccm93c2VyLk9wZW5VUkwodXJsKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB2b2lkIHdhaWxzLkJyb3dzZXIuT3BlblVSTCh1cmwpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gcmVsb2FkV01MKCkge1xyXG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcclxuICAgIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpO1xyXG4gICAgYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vLyBkZWZpbmVkIGluIHRoZSBUYXNrZmlsZVxyXG5leHBvcnQgbGV0IGludm9rZSA9IGZ1bmN0aW9uKGlucHV0KSB7XHJcbiAgICBpZihXSU5ET1dTKSB7XHJcbiAgICAgICAgY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UoaW5wdXQpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICB3ZWJraXQubWVzc2FnZUhhbmRsZXJzLmV4dGVybmFsLnBvc3RNZXNzYWdlKGlucHV0KTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxubGV0IGZsYWdzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gY29udmVydFRvTWFwKG9iaikge1xyXG4gICAgY29uc3QgbWFwID0gbmV3IE1hcCgpO1xyXG5cclxuICAgIGZvciAoY29uc3QgW2tleSwgdmFsdWVdIG9mIE9iamVjdC5lbnRyaWVzKG9iaikpIHtcclxuICAgICAgICBpZiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyAmJiB2YWx1ZSAhPT0gbnVsbCkge1xyXG4gICAgICAgICAgICBtYXAuc2V0KGtleSwgY29udmVydFRvTWFwKHZhbHVlKSk7IC8vIFJlY3Vyc2l2ZWx5IGNvbnZlcnQgbmVzdGVkIG9iamVjdFxyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIG1hcC5zZXQoa2V5LCB2YWx1ZSk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBtYXA7XHJcbn1cclxuXHJcbmZldGNoKFwiL3dhaWxzL2ZsYWdzXCIpLnRoZW4oKHJlc3BvbnNlKSA9PiB7XHJcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xyXG4gICAgICAgIGZsYWdzID0gY29udmVydFRvTWFwKGRhdGEpO1xyXG4gICAgfSk7XHJcbn0pO1xyXG5cclxuXHJcbmZ1bmN0aW9uIGdldFZhbHVlRnJvbU1hcChrZXlTdHJpbmcpIHtcclxuICAgIGNvbnN0IGtleXMgPSBrZXlTdHJpbmcuc3BsaXQoJy4nKTtcclxuICAgIGxldCB2YWx1ZSA9IGZsYWdzO1xyXG5cclxuICAgIGZvciAoY29uc3Qga2V5IG9mIGtleXMpIHtcclxuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBNYXApIHtcclxuICAgICAgICAgICAgdmFsdWUgPSB2YWx1ZS5nZXQoa2V5KTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICB2YWx1ZSA9IHZhbHVlW2tleV07XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAodmFsdWUgPT09IHVuZGVmaW5lZCkge1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIHZhbHVlO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXlTdHJpbmcpIHtcclxuICAgIHJldHVybiBnZXRWYWx1ZUZyb21NYXAoa2V5U3RyaW5nKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtpbnZva2V9IGZyb20gXCIuL2ludm9rZVwiO1xyXG5pbXBvcnQge0dldEZsYWd9IGZyb20gXCIuL2ZsYWdzXCI7XHJcblxyXG5sZXQgc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGRyYWdUZXN0KGUpIHtcclxuICAgIGxldCB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2Via2l0LWFwcC1yZWdpb25cIik7XHJcbiAgICBpZiAodmFsKSB7XHJcbiAgICAgICAgdmFsID0gdmFsLnRyaW0oKTtcclxuICAgIH1cclxuXHJcbiAgICBpZiAodmFsICE9PSBcImRyYWdcIikge1xyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBPbmx5IHByb2Nlc3MgdGhlIHByaW1hcnkgYnV0dG9uXHJcbiAgICBpZiAoZS5idXR0b25zICE9PSAxKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHNldHVwRHJhZygpIHtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCBvbk1vdXNlRG93bik7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgb25Nb3VzZU1vdmUpO1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCBvbk1vdXNlVXApO1xyXG59XHJcblxyXG5sZXQgcmVzaXplRWRnZSA9IG51bGw7XHJcblxyXG5mdW5jdGlvbiB0ZXN0UmVzaXplKGUpIHtcclxuICAgIGlmKCByZXNpemVFZGdlICkge1xyXG4gICAgICAgIGludm9rZShcInJlc2l6ZTpcIiArIHJlc2l6ZUVkZ2UpO1xyXG4gICAgICAgIHJldHVybiB0cnVlXHJcbiAgICB9XHJcbiAgICByZXR1cm4gZmFsc2U7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VEb3duKGUpIHtcclxuXHJcbiAgICAvLyBDaGVjayBmb3IgcmVzaXppbmcgb24gV2luZG93c1xyXG4gICAgaWYoIFdJTkRPV1MgKSB7XHJcbiAgICAgICAgaWYgKHRlc3RSZXNpemUoKSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG4gICAgaWYgKGRyYWdUZXN0KGUpKSB7XHJcbiAgICAgICAgLy8gSWdub3JlIGRyYWcgb24gc2Nyb2xsYmFyc1xyXG4gICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBzaG91bGREcmFnID0gdHJ1ZTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBvbk1vdXNlVXAoZSkge1xyXG4gICAgbGV0IG1vdXNlUHJlc3NlZCA9IGUuYnV0dG9ucyAhPT0gdW5kZWZpbmVkID8gZS5idXR0b25zIDogZS53aGljaDtcclxuICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XHJcbiAgICAgICAgZW5kRHJhZygpO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gZW5kRHJhZygpIHtcclxuICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gJ2RlZmF1bHQnO1xyXG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG59XHJcblxyXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XHJcbiAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gY3Vyc29yIHx8IGRlZmF1bHRDdXJzb3I7XHJcbiAgICByZXNpemVFZGdlID0gY3Vyc29yO1xyXG59XHJcblxyXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShlKSB7XHJcbiAgICBpZiAoc2hvdWxkRHJhZykge1xyXG4gICAgICAgIHNob3VsZERyYWcgPSBmYWxzZTtcclxuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xyXG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XHJcbiAgICAgICAgICAgIGludm9rZShcImRyYWdcIik7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoV0lORE9XUykge1xyXG4gICAgICAgIGhhbmRsZVJlc2l6ZShlKTtcclxuICAgIH1cclxufVxyXG5cclxubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcclxuXHJcbmZ1bmN0aW9uIGhhbmRsZVJlc2l6ZShlKSB7XHJcbiAgICBsZXQgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcclxuICAgIGxldCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcclxuXHJcbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXNcclxuICAgIGxldCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcclxuXHJcbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xyXG5cclxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyc1xyXG4gICAgbGV0IHJpZ2h0Q29ybmVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBsZXQgbGVmdENvcm5lciA9IGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGxldCB0b3BDb3JuZXIgPSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG4gICAgbGV0IGJvdHRvbUNvcm5lciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XHJcblxyXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcclxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xyXG4gICAgICAgIHNldFJlc2l6ZSgpO1xyXG4gICAgfVxyXG4gICAgLy8gQWRqdXN0ZWQgZm9yIGNvcm5lciBhcmVhc1xyXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcclxufVxyXG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5cclxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gJy4vY2xpcGJvYXJkJztcclxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSAnLi9hcHBsaWNhdGlvbic7XHJcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSAnLi9zY3JlZW5zJztcclxuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gJy4vc3lzdGVtJztcclxuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tICcuL2Jyb3dzZXInO1xyXG5pbXBvcnQge1BsdWdpbiwgQ2FsbCwgY2FsbEVycm9yQ2FsbGJhY2ssIGNhbGxDYWxsYmFjaywgQ2FsbEJ5SUQsIENhbGxCeU5hbWV9IGZyb20gXCIuL2NhbGxzXCI7XHJcbmltcG9ydCB7Y2xpZW50SWR9IGZyb20gJy4vcnVudGltZSc7XHJcbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcclxuaW1wb3J0IHtkaXNwYXRjaFdhaWxzRXZlbnQsIEVtaXQsIE9mZiwgT2ZmQWxsLCBPbiwgT25jZSwgT25NdWx0aXBsZX0gZnJvbSBcIi4vZXZlbnRzXCI7XHJcbmltcG9ydCB7ZGlhbG9nQ2FsbGJhY2ssIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssIEVycm9yLCBJbmZvLCBPcGVuRmlsZSwgUXVlc3Rpb24sIFNhdmVGaWxlLCBXYXJuaW5nLH0gZnJvbSBcIi4vZGlhbG9nc1wiO1xyXG5pbXBvcnQge3NldHVwQ29udGV4dE1lbnVzfSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xyXG5pbXBvcnQge3JlbG9hZFdNTH0gZnJvbSBcIi4vd21sXCI7XHJcbmltcG9ydCB7c2V0dXBEcmFnLCBlbmREcmFnfSBmcm9tIFwiLi9kcmFnXCI7XHJcblxyXG53aW5kb3cud2FpbHMgPSB7XHJcbiAgICAuLi5uZXdSdW50aW1lKG51bGwpLFxyXG4gICAgQ2FwYWJpbGl0aWVzOiB7fSxcclxuICAgIGNsaWVudElkOiBjbGllbnRJZCxcclxufTtcclxuXHJcbmZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKS50aGVuKChyZXNwb25zZSkgPT4ge1xyXG4gICAgcmVzcG9uc2UuanNvbigpLnRoZW4oKGRhdGEpID0+IHtcclxuICAgICAgICB3aW5kb3cud2FpbHMuQ2FwYWJpbGl0aWVzID0gZGF0YTtcclxuICAgIH0pO1xyXG59KTtcclxuXHJcbi8vIEludGVybmFsIHdhaWxzIGVuZHBvaW50c1xyXG53aW5kb3cuX3dhaWxzID0ge1xyXG4gICAgZGlhbG9nQ2FsbGJhY2ssXHJcbiAgICBkaWFsb2dFcnJvckNhbGxiYWNrLFxyXG4gICAgZGlzcGF0Y2hXYWlsc0V2ZW50LFxyXG4gICAgY2FsbENhbGxiYWNrLFxyXG4gICAgY2FsbEVycm9yQ2FsbGJhY2ssXHJcbiAgICBlbmREcmFnLFxyXG59O1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWUod2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIHtcclxuICAgICAgICBDbGlwYm9hcmQ6IHtcclxuICAgICAgICAgICAgLi4uQ2xpcGJvYXJkXHJcbiAgICAgICAgfSxcclxuICAgICAgICBBcHBsaWNhdGlvbjoge1xyXG4gICAgICAgICAgICAuLi5BcHBsaWNhdGlvbixcclxuICAgICAgICAgICAgR2V0V2luZG93QnlOYW1lKHdpbmRvd05hbWUpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBuZXdSdW50aW1lKHdpbmRvd05hbWUpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSxcclxuICAgICAgICBTeXN0ZW0sXHJcbiAgICAgICAgU2NyZWVucyxcclxuICAgICAgICBCcm93c2VyLFxyXG4gICAgICAgIENhbGwsXHJcbiAgICAgICAgQ2FsbEJ5SUQsXHJcbiAgICAgICAgQ2FsbEJ5TmFtZSxcclxuICAgICAgICBQbHVnaW4sXHJcbiAgICAgICAgV01MOiB7XHJcbiAgICAgICAgICAgIFJlbG9hZDogcmVsb2FkV01MLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgRGlhbG9nOiB7XHJcbiAgICAgICAgICAgIEluZm8sXHJcbiAgICAgICAgICAgIFdhcm5pbmcsXHJcbiAgICAgICAgICAgIEVycm9yLFxyXG4gICAgICAgICAgICBRdWVzdGlvbixcclxuICAgICAgICAgICAgT3BlbkZpbGUsXHJcbiAgICAgICAgICAgIFNhdmVGaWxlLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgRXZlbnRzOiB7XHJcbiAgICAgICAgICAgIEVtaXQsXHJcbiAgICAgICAgICAgIE9uLFxyXG4gICAgICAgICAgICBPbmNlLFxyXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxyXG4gICAgICAgICAgICBPZmYsXHJcbiAgICAgICAgICAgIE9mZkFsbCxcclxuICAgICAgICB9LFxyXG4gICAgICAgIFdpbmRvdzogbmV3V2luZG93KHdpbmRvd05hbWUpLFxyXG4gICAgfTtcclxufVxyXG5cclxuaWYgKERFQlVHKSB7XHJcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XHJcbn1cclxuXHJcbnNldHVwQ29udGV4dE1lbnVzKCk7XHJcbnNldHVwRHJhZygpO1xyXG5cclxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcihcIkRPTUNvbnRlbnRMb2FkZWRcIiwgZnVuY3Rpb24oKSB7XHJcbiAgICByZWxvYWRXTUwoKTtcclxufSk7XHJcbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUEsTUFBSSxjQUNGO0FBV0ssTUFBSSxTQUFTLENBQUMsT0FBTyxPQUFPO0FBQ2pDLFFBQUksS0FBSztBQUNULFFBQUksSUFBSTtBQUNSLFdBQU8sS0FB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V0EsV0FBVTtBQUFBLFFBQ2hDO0FBQUEsTUFDSjtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFFBQVE7QUFBQSxNQUNaO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxVQUFVO0FBQUEsSUFDaEM7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEsb0JBQWtCO0FBQ2xCLFlBQVU7QUFFVixXQUFTLGlCQUFpQixvQkFBb0IsV0FBVztBQUNyRCxjQUFVO0FBQUEsRUFDZCxDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiRXJyb3IiLCAiY2FsbCIsICJ3aW5kb3dOYW1lIiwgIkVycm9yIl0KfQo=
