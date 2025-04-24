(function polyfill() {
  const relList = document.createElement("link").relList;
  if (relList && relList.supports && relList.supports("modulepreload")) {
    return;
  }
  for (const link of document.querySelectorAll('link[rel="modulepreload"]')) {
    processPreload(link);
  }
  new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.type !== "childList") {
        continue;
      }
      for (const node of mutation.addedNodes) {
        if (node.tagName === "LINK" && node.rel === "modulepreload")
          processPreload(node);
      }
    }
  }).observe(document, { childList: true, subtree: true });
  function getFetchOpts(link) {
    const fetchOpts = {};
    if (link.integrity) fetchOpts.integrity = link.integrity;
    if (link.referrerPolicy) fetchOpts.referrerPolicy = link.referrerPolicy;
    if (link.crossOrigin === "use-credentials")
      fetchOpts.credentials = "include";
    else if (link.crossOrigin === "anonymous") fetchOpts.credentials = "omit";
    else fetchOpts.credentials = "same-origin";
    return fetchOpts;
  }
  function processPreload(link) {
    if (link.ep)
      return;
    link.ep = true;
    const fetchOpts = getFetchOpts(link);
    fetch(link.href, fetchOpts);
  }
})();
let urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
let nanoid = (size = 21) => {
  let id = "";
  let i = size;
  while (i--) {
    id += urlAlphabet[Math.random() * 64 | 0];
  }
  return id;
};
const runtimeURL = window.location.origin + "/wails/runtime";
const objectNames = {
  Call: 0,
  Clipboard: 1,
  Application: 2,
  Events: 3,
  ContextMenu: 4,
  Dialog: 5,
  Window: 6,
  Screens: 7,
  System: 8,
  Browser: 9,
  CancelCall: 10
};
let clientId = nanoid();
function newRuntimeCallerWithID(object, windowName) {
  return function(method, args = null) {
    return runtimeCallWithID(object, method, windowName, args);
  };
}
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
function invoke(msg) {
  if (window.chrome) {
    return window.chrome.webview.postMessage(msg);
  }
  return window.webkit.messageHandlers.external.postMessage(msg);
}
function IsWindows() {
  return window._wails.environment.OS === "windows";
}
function IsDebug() {
  return window._wails.environment.Debug === true;
}
window.addEventListener("contextmenu", contextMenuHandler);
const call$4 = newRuntimeCallerWithID(objectNames.ContextMenu, "");
const ContextMenuOpen = 0;
function openContextMenu(id, x, y, data) {
  void call$4(ContextMenuOpen, { id, x, y, data });
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
  if (IsDebug()) {
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
function GetFlag(keyString) {
  try {
    return window._wails.flags[keyString];
  } catch (e) {
    throw new Error("Unable to retrieve flag '" + keyString + "': " + e);
  }
}
let shouldDrag = false;
let resizable = false;
let resizeEdge = null;
let defaultCursor = "auto";
window._wails = window._wails || {};
window._wails.setResizable = function(value) {
  resizable = value;
};
window._wails.endDrag = function() {
  document.body.style.cursor = "default";
  shouldDrag = false;
};
window.addEventListener("mousedown", onMouseDown);
window.addEventListener("mousemove", onMouseMove);
window.addEventListener("mouseup", onMouseUp);
function dragTest(e) {
  let val = window.getComputedStyle(e.target).getPropertyValue("--wails-draggable");
  let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
  if (!val || val === "" || val.trim() !== "drag" || mousePressed === 0) {
    return false;
  }
  return e.detail === 1;
}
function onMouseDown(e) {
  if (resizeEdge) {
    invoke("wails:resize:" + resizeEdge);
    e.preventDefault();
    return;
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
function onMouseUp() {
  shouldDrag = false;
}
function setResize(cursor) {
  document.documentElement.style.cursor = cursor || defaultCursor;
  resizeEdge = cursor;
}
function onMouseMove(e) {
  if (shouldDrag) {
    shouldDrag = false;
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (mousePressed > 0) {
      invoke("wails:drag");
      return;
    }
  }
  if (!resizable || !IsWindows()) {
    return;
  }
  let resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
  let resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
  let cornerExtra = GetFlag("resizeCornerExtra") || 10;
  let rightBorder = window.outerWidth - e.clientX < resizeHandleWidth;
  let leftBorder = e.clientX < resizeHandleWidth;
  let topBorder = e.clientY < resizeHandleHeight;
  let bottomBorder = window.outerHeight - e.clientY < resizeHandleHeight;
  let rightCorner = window.outerWidth - e.clientX < resizeHandleWidth + cornerExtra;
  let leftCorner = e.clientX < resizeHandleWidth + cornerExtra;
  let topCorner = e.clientY < resizeHandleHeight + cornerExtra;
  let bottomCorner = window.outerHeight - e.clientY < resizeHandleHeight + cornerExtra;
  if (!leftBorder && !rightBorder && !topBorder && !bottomBorder && resizeEdge !== void 0) {
    setResize();
  } else if (rightCorner && bottomCorner) setResize("se-resize");
  else if (leftCorner && bottomCorner) setResize("sw-resize");
  else if (leftCorner && topCorner) setResize("nw-resize");
  else if (topCorner && rightCorner) setResize("ne-resize");
  else if (leftBorder) setResize("w-resize");
  else if (topBorder) setResize("n-resize");
  else if (bottomBorder) setResize("s-resize");
  else if (rightBorder) setResize("e-resize");
}
const call$3 = newRuntimeCallerWithID(objectNames.Browser, "");
const BrowserOpenURL = 0;
function OpenURL(url) {
  return call$3(BrowserOpenURL, { url });
}
window._wails = window._wails || {};
window._wails.callResultHandler = resultHandler;
window._wails.callErrorHandler = errorHandler;
const CallBinding = 0;
const call$2 = newRuntimeCallerWithID(objectNames.Call, "");
const cancelCall = newRuntimeCallerWithID(objectNames.CancelCall, "");
let callResponses = /* @__PURE__ */ new Map();
function generateID$1() {
  let result;
  do {
    result = nanoid();
  } while (callResponses.has(result));
  return result;
}
function resultHandler(id, data, isJSON) {
  const promiseHandler = getAndDeleteResponse(id);
  if (promiseHandler) {
    promiseHandler.resolve(isJSON ? JSON.parse(data) : data);
  }
}
function errorHandler(id, message) {
  const promiseHandler = getAndDeleteResponse(id);
  if (promiseHandler) {
    promiseHandler.reject(message);
  }
}
function getAndDeleteResponse(id) {
  const response = callResponses.get(id);
  callResponses.delete(id);
  return response;
}
function callBinding(type, options = {}) {
  const id = generateID$1();
  const doCancel = () => {
    return cancelCall(type, { "call-id": id });
  };
  let queuedCancel = false, callRunning = false;
  let p = new Promise((resolve, reject) => {
    options["call-id"] = id;
    callResponses.set(id, { resolve, reject });
    call$2(type, options).then((_) => {
      callRunning = true;
      if (queuedCancel) {
        return doCancel();
      }
    }).catch((error) => {
      reject(error);
      callResponses.delete(id);
    });
  });
  p.cancel = () => {
    if (callRunning) {
      return doCancel();
    } else {
      queuedCancel = true;
    }
  };
  return p;
}
function ByID(methodID, ...args) {
  return callBinding(CallBinding, {
    methodID,
    args
  });
}
window._wails = window._wails || {};
window._wails.dialogErrorCallback = dialogErrorCallback;
window._wails.dialogResultCallback = dialogResultCallback;
const DialogQuestion = 3;
const call$1 = newRuntimeCallerWithID(objectNames.Dialog, "");
const dialogResponses = /* @__PURE__ */ new Map();
function generateID() {
  let result;
  do {
    result = nanoid();
  } while (dialogResponses.has(result));
  return result;
}
function dialog(type, options = {}) {
  const id = generateID();
  options["dialog-id"] = id;
  return new Promise((resolve, reject) => {
    dialogResponses.set(id, { resolve, reject });
    call$1(type, options).catch((error) => {
      reject(error);
      dialogResponses.delete(id);
    });
  });
}
function dialogResultCallback(id, data, isJSON) {
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
const Question = (options) => dialog(DialogQuestion, options);
window._wails = window._wails || {};
window._wails.dispatchWailsEvent = dispatchWailsEvent;
const call = newRuntimeCallerWithID(objectNames.Events, "");
const EmitMethod = 0;
const eventListeners = /* @__PURE__ */ new Map();
class Listener {
  constructor(eventName, callback, maxCallbacks) {
    this.eventName = eventName;
    this.maxCallbacks = maxCallbacks || -1;
    this.Callback = (data) => {
      callback(data);
      if (this.maxCallbacks === -1) return false;
      this.maxCallbacks -= 1;
      return this.maxCallbacks === 0;
    };
  }
}
class WailsEvent {
  constructor(name, data = null) {
    this.name = name;
    this.data = data;
  }
}
function dispatchWailsEvent(event) {
  let listeners = eventListeners.get(event.name);
  if (listeners) {
    let toRemove = listeners.filter((listener) => {
      let remove = listener.Callback(event);
      if (remove) return true;
    });
    if (toRemove.length > 0) {
      listeners = listeners.filter((l) => !toRemove.includes(l));
      if (listeners.length === 0) eventListeners.delete(event.name);
      else eventListeners.set(event.name, listeners);
    }
  }
}
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
function listenerOff(listener) {
  const eventName = listener.eventName;
  let listeners = eventListeners.get(eventName).filter((l) => l !== listener);
  if (listeners.length === 0) eventListeners.delete(eventName);
  else eventListeners.set(eventName, listeners);
}
function Emit(event) {
  return call(EmitMethod, event);
}
const PositionMethod = 0;
const CenterMethod = 1;
const CloseMethod = 2;
const DisableSizeConstraintsMethod = 3;
const EnableSizeConstraintsMethod = 4;
const FocusMethod = 5;
const ForceReloadMethod = 6;
const FullscreenMethod = 7;
const GetScreenMethod = 8;
const GetZoomMethod = 9;
const HeightMethod = 10;
const HideMethod = 11;
const IsFocusedMethod = 12;
const IsFullscreenMethod = 13;
const IsMaximisedMethod = 14;
const IsMinimisedMethod = 15;
const MaximiseMethod = 16;
const MinimiseMethod = 17;
const NameMethod = 18;
const OpenDevToolsMethod = 19;
const RelativePositionMethod = 20;
const ReloadMethod = 21;
const ResizableMethod = 22;
const RestoreMethod = 23;
const SetPositionMethod = 24;
const SetAlwaysOnTopMethod = 25;
const SetBackgroundColourMethod = 26;
const SetFramelessMethod = 27;
const SetFullscreenButtonEnabledMethod = 28;
const SetMaxSizeMethod = 29;
const SetMinSizeMethod = 30;
const SetRelativePositionMethod = 31;
const SetResizableMethod = 32;
const SetSizeMethod = 33;
const SetTitleMethod = 34;
const SetZoomMethod = 35;
const ShowMethod = 36;
const SizeMethod = 37;
const ToggleFullscreenMethod = 38;
const ToggleMaximiseMethod = 39;
const UnFullscreenMethod = 40;
const UnMaximiseMethod = 41;
const UnMinimiseMethod = 42;
const WidthMethod = 43;
const ZoomMethod = 44;
const ZoomInMethod = 45;
const ZoomOutMethod = 46;
const ZoomResetMethod = 47;
const caller = Symbol();
class Window {
  /**
   * Initialises a window object with the specified name.
   *
   * @private
   * @param {string} name - The name of the target window.
   */
  constructor(name = "") {
    this[caller] = newRuntimeCallerWithID(objectNames.Window, name);
    for (const method of Object.getOwnPropertyNames(Window.prototype)) {
      if (method !== "constructor" && typeof this[method] === "function") {
        this[method] = this[method].bind(this);
      }
    }
  }
  /**
   * Gets the specified window.
   *
   * @public
   * @param {string} name - The name of the window to get.
   * @return {Window} - The corresponding window object.
   */
  Get(name) {
    return new Window(name);
  }
  /**
   * Returns the absolute position of the window.
   *
   * @public
   * @return {Promise<Position>} - The current absolute position of the window.
   */
  Position() {
    return this[caller](PositionMethod);
  }
  /**
   * Centers the window on the screen.
   *
   * @public
   * @return {Promise<void>}
   */
  Center() {
    return this[caller](CenterMethod);
  }
  /**
   * Closes the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Close() {
    return this[caller](CloseMethod);
  }
  /**
   * Disables min/max size constraints.
   *
   * @public
   * @return {Promise<void>}
   */
  DisableSizeConstraints() {
    return this[caller](DisableSizeConstraintsMethod);
  }
  /**
   * Enables min/max size constraints.
   *
   * @public
   * @return {Promise<void>}
   */
  EnableSizeConstraints() {
    return this[caller](EnableSizeConstraintsMethod);
  }
  /**
   * Focuses the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Focus() {
    return this[caller](FocusMethod);
  }
  /**
   * Forces the window to reload the page assets.
   *
   * @public
   * @return {Promise<void>}
   */
  ForceReload() {
    return this[caller](ForceReloadMethod);
  }
  /**
   * Doc.
   *
   * @public
   * @return {Promise<void>}
   */
  Fullscreen() {
    return this[caller](FullscreenMethod);
  }
  /**
   * Returns the screen that the window is on.
   *
   * @public
   * @return {Promise<Screen>} - The screen the window is currently on
   */
  GetScreen() {
    return this[caller](GetScreenMethod);
  }
  /**
   * Returns the current zoom level of the window.
   *
   * @public
   * @return {Promise<number>} - The current zoom level
   */
  GetZoom() {
    return this[caller](GetZoomMethod);
  }
  /**
   * Returns the height of the window.
   *
   * @public
   * @return {Promise<number>} - The current height of the window
   */
  Height() {
    return this[caller](HeightMethod);
  }
  /**
   * Hides the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Hide() {
    return this[caller](HideMethod);
  }
  /**
   * Returns true if the window is focused.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently focused
   */
  IsFocused() {
    return this[caller](IsFocusedMethod);
  }
  /**
   * Returns true if the window is fullscreen.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently fullscreen
   */
  IsFullscreen() {
    return this[caller](IsFullscreenMethod);
  }
  /**
   * Returns true if the window is maximised.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently maximised
   */
  IsMaximised() {
    return this[caller](IsMaximisedMethod);
  }
  /**
   * Returns true if the window is minimised.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently minimised
   */
  IsMinimised() {
    return this[caller](IsMinimisedMethod);
  }
  /**
   * Maximises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Maximise() {
    return this[caller](MaximiseMethod);
  }
  /**
   * Minimises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Minimise() {
    return this[caller](MinimiseMethod);
  }
  /**
   * Returns the name of the window.
   *
   * @public
   * @return {Promise<string>} - The name of the window
   */
  Name() {
    return this[caller](NameMethod);
  }
  /**
   * Opens the development tools pane.
   *
   * @public
   * @return {Promise<void>}
   */
  OpenDevTools() {
    return this[caller](OpenDevToolsMethod);
  }
  /**
   * Returns the relative position of the window to the screen.
   *
   * @public
   * @return {Promise<Position>} - The current relative position of the window
   */
  RelativePosition() {
    return this[caller](RelativePositionMethod);
  }
  /**
   * Reloads the page assets.
   *
   * @public
   * @return {Promise<void>}
   */
  Reload() {
    return this[caller](ReloadMethod);
  }
  /**
   * Returns true if the window is resizable.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently resizable
   */
  Resizable() {
    return this[caller](ResizableMethod);
  }
  /**
   * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
   *
   * @public
   * @return {Promise<void>}
   */
  Restore() {
    return this[caller](RestoreMethod);
  }
  /**
   * Sets the absolute position of the window.
   *
   * @public
   * @param {number} x - The desired horizontal absolute position of the window
   * @param {number} y - The desired vertical absolute position of the window
   * @return {Promise<void>}
   */
  SetPosition(x, y) {
    return this[caller](SetPositionMethod, { x, y });
  }
  /**
   * Sets the window to be always on top.
   *
   * @public
   * @param {boolean} alwaysOnTop - Whether the window should stay on top
   * @return {Promise<void>}
   */
  SetAlwaysOnTop(alwaysOnTop) {
    return this[caller](SetAlwaysOnTopMethod, { alwaysOnTop });
  }
  /**
   * Sets the background colour of the window.
   *
   * @public
   * @param {number} r - The desired red component of the window background
   * @param {number} g - The desired green component of the window background
   * @param {number} b - The desired blue component of the window background
   * @param {number} a - The desired alpha component of the window background
   * @return {Promise<void>}
   */
  SetBackgroundColour(r, g, b, a) {
    return this[caller](SetBackgroundColourMethod, { r, g, b, a });
  }
  /**
   * Removes the window frame and title bar.
   *
   * @public
   * @param {boolean} frameless - Whether the window should be frameless
   * @return {Promise<void>}
   */
  SetFrameless(frameless) {
    return this[caller](SetFramelessMethod, { frameless });
  }
  /**
   * Disables the system fullscreen button.
   *
   * @public
   * @param {boolean} enabled - Whether the fullscreen button should be enabled
   * @return {Promise<void>}
   */
  SetFullscreenButtonEnabled(enabled) {
    return this[caller](SetFullscreenButtonEnabledMethod, { enabled });
  }
  /**
   * Sets the maximum size of the window.
   *
   * @public
   * @param {number} width - The desired maximum width of the window
   * @param {number} height - The desired maximum height of the window
   * @return {Promise<void>}
   */
  SetMaxSize(width, height) {
    return this[caller](SetMaxSizeMethod, { width, height });
  }
  /**
   * Sets the minimum size of the window.
   *
   * @public
   * @param {number} width - The desired minimum width of the window
   * @param {number} height - The desired minimum height of the window
   * @return {Promise<void>}
   */
  SetMinSize(width, height) {
    return this[caller](SetMinSizeMethod, { width, height });
  }
  /**
   * Sets the relative position of the window to the screen.
   *
   * @public
   * @param {number} x - The desired horizontal relative position of the window
   * @param {number} y - The desired vertical relative position of the window
   * @return {Promise<void>}
   */
  SetRelativePosition(x, y) {
    return this[caller](SetRelativePositionMethod, { x, y });
  }
  /**
   * Sets whether the window is resizable.
   *
   * @public
   * @param {boolean} resizable - Whether the window should be resizable
   * @return {Promise<void>}
   */
  SetResizable(resizable2) {
    return this[caller](SetResizableMethod, { resizable: resizable2 });
  }
  /**
   * Sets the size of the window.
   *
   * @public
   * @param {number} width - The desired width of the window
   * @param {number} height - The desired height of the window
   * @return {Promise<void>}
   */
  SetSize(width, height) {
    return this[caller](SetSizeMethod, { width, height });
  }
  /**
   * Sets the title of the window.
   *
   * @public
   * @param {string} title - The desired title of the window
   * @return {Promise<void>}
   */
  SetTitle(title) {
    return this[caller](SetTitleMethod, { title });
  }
  /**
   * Sets the zoom level of the window.
   *
   * @public
   * @param {number} zoom - The desired zoom level
   * @return {Promise<void>}
   */
  SetZoom(zoom) {
    return this[caller](SetZoomMethod, { zoom });
  }
  /**
   * Shows the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Show() {
    return this[caller](ShowMethod);
  }
  /**
   * Returns the size of the window.
   *
   * @public
   * @return {Promise<Size>} - The current size of the window
   */
  Size() {
    return this[caller](SizeMethod);
  }
  /**
   * Toggles the window between fullscreen and normal.
   *
   * @public
   * @return {Promise<void>}
   */
  ToggleFullscreen() {
    return this[caller](ToggleFullscreenMethod);
  }
  /**
   * Toggles the window between maximised and normal.
   *
   * @public
   * @return {Promise<void>}
   */
  ToggleMaximise() {
    return this[caller](ToggleMaximiseMethod);
  }
  /**
   * Un-fullscreens the window.
   *
   * @public
   * @return {Promise<void>}
   */
  UnFullscreen() {
    return this[caller](UnFullscreenMethod);
  }
  /**
   * Un-maximises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  UnMaximise() {
    return this[caller](UnMaximiseMethod);
  }
  /**
   * Un-minimises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  UnMinimise() {
    return this[caller](UnMinimiseMethod);
  }
  /**
   * Returns the width of the window.
   *
   * @public
   * @return {Promise<number>} - The current width of the window
   */
  Width() {
    return this[caller](WidthMethod);
  }
  /**
   * Zooms the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Zoom() {
    return this[caller](ZoomMethod);
  }
  /**
   * Increases the zoom level of the webview content.
   *
   * @public
   * @return {Promise<void>}
   */
  ZoomIn() {
    return this[caller](ZoomInMethod);
  }
  /**
   * Decreases the zoom level of the webview content.
   *
   * @public
   * @return {Promise<void>}
   */
  ZoomOut() {
    return this[caller](ZoomOutMethod);
  }
  /**
   * Resets the zoom level of the webview content.
   *
   * @public
   * @return {Promise<void>}
   */
  ZoomReset() {
    return this[caller](ZoomResetMethod);
  }
}
const thisWindow = new Window("");
function canAbortListeners() {
  if (!EventTarget || !AbortSignal || !AbortController)
    return false;
  let result = true;
  const target = new EventTarget();
  const controller2 = new AbortController();
  target.addEventListener("test", () => {
    result = false;
  }, { signal: controller2.signal });
  controller2.abort();
  target.dispatchEvent(new CustomEvent("test"));
  return result;
}
document.addEventListener("DOMContentLoaded", () => true);
function sendEvent(eventName, data = null) {
  Emit(new WailsEvent(eventName, data));
}
function callWindowMethod(windowName, methodName) {
  const targetWindow = thisWindow.Get(windowName);
  const method = targetWindow[methodName];
  if (typeof method !== "function") {
    console.error(`Window method '${methodName}' not found`);
    return;
  }
  try {
    method.call(targetWindow);
  } catch (e) {
    console.error(`Error calling window method '${methodName}': `, e);
  }
}
function onWMLTriggered(ev) {
  const element = ev.currentTarget;
  function runEffect(choice = "Yes") {
    if (choice !== "Yes")
      return;
    const eventType = element.getAttribute("wml-event");
    const targetWindow = element.getAttribute("wml-target-window") || "";
    const windowMethod = element.getAttribute("wml-window");
    const url = element.getAttribute("wml-openurl");
    if (eventType !== null)
      sendEvent(eventType);
    if (windowMethod !== null)
      callWindowMethod(targetWindow, windowMethod);
    if (url !== null)
      void OpenURL(url);
  }
  const confirm = element.getAttribute("wml-confirm");
  if (confirm) {
    Question({
      Title: "Confirm",
      Message: confirm,
      Detached: false,
      Buttons: [
        { Label: "Yes" },
        { Label: "No", IsDefault: true }
      ]
    }).then(runEffect);
  } else {
    runEffect();
  }
}
const controller = Symbol();
class AbortControllerRegistry {
  constructor() {
    this[controller] = new AbortController();
  }
  /**
   * Returns an options object for addEventListener that ties the listener
   * to the AbortSignal from the current AbortController.
   *
   * @param {HTMLElement} element An HTML element
   * @param {string[]} triggers The list of active WML trigger events for the specified elements
   * @returns {AddEventListenerOptions}
   */
  set(element, triggers) {
    return { signal: this[controller].signal };
  }
  /**
   * Removes all registered event listeners.
   *
   * @returns {void}
   */
  reset() {
    this[controller].abort();
    this[controller] = new AbortController();
  }
}
const triggerMap = Symbol();
const elementCount = Symbol();
class WeakMapRegistry {
  constructor() {
    this[triggerMap] = /* @__PURE__ */ new WeakMap();
    this[elementCount] = 0;
  }
  /**
   * Sets the active triggers for the specified element.
   *
   * @param {HTMLElement} element An HTML element
   * @param {string[]} triggers The list of active WML trigger events for the specified element
   * @returns {AddEventListenerOptions}
   */
  set(element, triggers) {
    this[elementCount] += !this[triggerMap].has(element);
    this[triggerMap].set(element, triggers);
    return {};
  }
  /**
   * Removes all registered event listeners.
   *
   * @returns {void}
   */
  reset() {
    if (this[elementCount] <= 0)
      return;
    for (const element of document.body.querySelectorAll("*")) {
      if (this[elementCount] <= 0)
        break;
      const triggers = this[triggerMap].get(element);
      this[elementCount] -= typeof triggers !== "undefined";
      for (const trigger of triggers || [])
        element.removeEventListener(trigger, onWMLTriggered);
    }
    this[triggerMap] = /* @__PURE__ */ new WeakMap();
    this[elementCount] = 0;
  }
}
canAbortListeners() ? new AbortControllerRegistry() : new WeakMapRegistry();
window._wails = window._wails || {};
window._wails.invoke = invoke;
invoke("wails:runtime:ready");
function Greet(name) {
  let $resultPromise = (
    /** @type {any} */
    ByID(1411160069, name)
  );
  return $resultPromise;
}
const resultElement = document.getElementById("result");
const timeElement = document.getElementById("time");
window.doGreet = () => {
  let name = document.getElementById("name").value;
  if (!name) {
    name = "anonymous";
  }
  Greet(name).then((result) => {
    resultElement.innerText = result;
  }).catch((err) => {
    console.log(err);
  });
};
On("time", (time) => {
  timeElement.innerText = time.data;
});
