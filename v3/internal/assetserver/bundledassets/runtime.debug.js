// \0rolldown/runtime.js
var __defProp = Object.defineProperty;
var __esmMin = (fn, res, err) => () => {
	if (err) throw err[0];
	try {
		return fn && (res = fn(fn = 0)), res;
	} catch (e) {
		throw err = [e], e;
	}
};
var __exportAll = (all, no_symbols) => {
	let target = {};
	for (var name in all) {
		__defProp(target, name, {
			get: all[name],
			enumerable: true
		});
	}
	if (!no_symbols) {
		__defProp(target, Symbol.toStringTag, { value: "Module" });
	}
	return target;
};

// src/environment.ts
var hasDOM;
function init_environment() {
	return (init_environment = __esmMin((() => {
		hasDOM = typeof window !== "undefined" && typeof document !== "undefined";
	})))();
}

// src/nanoid.ts
function nanoid(size = 21) {
	const source = typeof crypto !== "undefined" && typeof crypto.getRandomValues === "function" ? crypto : null;
	if (source) {
		const bytes = new Uint8Array(size);
		source.getRandomValues(bytes);
		let id = "";
		for (let i = 0; i < size; i++) id += urlAlphabet[bytes[i] & 63];
		return id;
	}
	let id = "";
	let i = size | 0;
	while (i--) id += urlAlphabet[Math.random() * 64 | 0];
	return id;
}
var urlAlphabet;
function init_nanoid() {
	return (init_nanoid = __esmMin((() => {
		urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
	})))();
}

// \0@oxc-project+runtime@0.155.0/helpers/esm/typeof.js
function _typeof(o) {
	"@babel/helpers - typeof";
	return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function(o) {
		return typeof o;
	} : function(o) {
		return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o;
	}, _typeof(o);
}
function init_typeof() {
	return (init_typeof = __esmMin((() => {})))();
}

// \0@oxc-project+runtime@0.155.0/helpers/esm/toPrimitive.js
function toPrimitive(t, r) {
	if ("object" != _typeof(t) || !t) return t;
	var e = t[Symbol.toPrimitive];
	if (void 0 !== e) {
		var i = e.call(t, r || "default");
		if ("object" != _typeof(i)) return i;
		throw new TypeError("@@toPrimitive must return a primitive value.");
	}
	return ("string" === r ? String : Number)(t);
}
function init_toPrimitive() {
	return (init_toPrimitive = __esmMin((() => {
		/* @__PURE__ */ init_typeof();
	})))();
}

// \0@oxc-project+runtime@0.155.0/helpers/esm/toPropertyKey.js
function toPropertyKey(t) {
	var i = toPrimitive(t, "string");
	return "symbol" == _typeof(i) ? i : i + "";
}
function init_toPropertyKey() {
	return (init_toPropertyKey = __esmMin((() => {
		/* @__PURE__ */ init_typeof();
		init_toPrimitive();
	})))();
}

// \0@oxc-project+runtime@0.155.0/helpers/esm/defineProperty.js
function _defineProperty(e, r, t) {
	return (r = toPropertyKey(r)) in e ? Object.defineProperty(e, r, {
		value: t,
		enumerable: !0,
		configurable: !0,
		writable: !0
	}) : e[r] = t, e;
}
function init_defineProperty() {
	return (init_defineProperty = __esmMin((() => {
		init_toPropertyKey();
	})))();
}

// \0@oxc-project+runtime@0.155.0/helpers/esm/objectSpread2.js
function ownKeys(e, r) {
	var t = Object.keys(e);
	if (Object.getOwnPropertySymbols) {
		var o = Object.getOwnPropertySymbols(e);
		r && (o = o.filter(function(r) {
			return Object.getOwnPropertyDescriptor(e, r).enumerable;
		})), t.push.apply(t, o);
	}
	return t;
}
function _objectSpread2(e) {
	for (var r = 1; r < arguments.length; r++) {
		var t = null != arguments[r] ? arguments[r] : {};
		r % 2 ? ownKeys(Object(t), !0).forEach(function(r) {
			_defineProperty(e, r, t[r]);
		}) : Object.getOwnPropertyDescriptors ? Object.defineProperties(e, Object.getOwnPropertyDescriptors(t)) : ownKeys(Object(t)).forEach(function(r) {
			Object.defineProperty(e, r, Object.getOwnPropertyDescriptor(t, r));
		});
	}
	return e;
}
function init_objectSpread2() {
	return (init_objectSpread2 = __esmMin((() => {
		init_defineProperty();
	})))();
}

// src/runtime.ts
function runtimeURL() {
	return window.location.origin + "/wails/runtime";
}
function setTransport(transport) {
	customTransport = transport;
}
function getTransport() {
	return customTransport;
}
function getOpenHarmonyTransport() {
	var _window$crypto$random, _window$crypto, _window$crypto$random2;
	if (detectedOpenHarmonyTransport !== void 0) return detectedOpenHarmonyTransport;
	detectedOpenHarmonyTransport = null;
	if (typeof window === "undefined") return null;
	const hostWindow = window;
	const host = hostWindow.wailsOpenHarmony;
	if (!host || typeof host.invoke !== "function" || typeof host.request !== "function") return null;
	const pending = /* @__PURE__ */ new Map();
	const session = (_window$crypto$random = (_window$crypto = window.crypto) === null || _window$crypto === void 0 || (_window$crypto$random2 = _window$crypto.randomUUID) === null || _window$crypto$random2 === void 0 ? void 0 : _window$crypto$random2.call(_window$crypto)) !== null && _window$crypto$random !== void 0 ? _window$crypto$random : nanoid();
	let nextID = 0;
	hostWindow.wailsOpenHarmonySession = session;
	hostWindow.wailsOpenHarmonyResolve = (id, payload) => {
		const call = pending.get(id);
		if (!call) return;
		pending.delete(id);
		window.clearTimeout(call.timer);
		try {
			var _reply$error, _reply$text;
			const parsed = JSON.parse(payload);
			if (typeof parsed !== "object" || parsed === null) throw new Error("Invalid runtime reply");
			const reply = parsed;
			if (!reply.ok) throw new Error((_reply$error = reply.error) !== null && _reply$error !== void 0 ? _reply$error : "OpenHarmony runtime call failed");
			call.resolve((_reply$text = reply.text) !== null && _reply$text !== void 0 ? _reply$text : reply.data);
		} catch (error) {
			call.reject(error);
		}
	};
	if (!hostWindow.wails) hostWindow.wails = {};
	if (!hostWindow.wails.invoke) hostWindow.wails.invoke = host.invoke.bind(host);
	detectedOpenHarmonyTransport = { call(objectID, method, windowName, args) {
		const id = String(++nextID);
		const payload = JSON.stringify({
			object: objectID,
			method,
			windowName,
			args,
			clientId
		});
		const { promise, resolve, reject } = Promise.withResolvers();
		const timer = window.setTimeout(() => {
			pending.delete(id);
			reject(/* @__PURE__ */ new Error("OpenHarmony runtime call timed out"));
		}, 15e3);
		pending.set(id, {
			resolve,
			reject,
			timer
		});
		if (!host.request(session, id, payload)) {
			pending.delete(id);
			window.clearTimeout(timer);
			reject(/* @__PURE__ */ new Error("OpenHarmony bridge unavailable"));
		}
		return promise;
	} };
	return detectedOpenHarmonyTransport;
}
function newRuntimeCaller(object, windowName = "") {
	return function(method, args = null) {
		return runtimeCallWithID(object, method, windowName, args);
	};
}
async function runtimeCallWithID(objectID, method, windowName, args) {
	var _response$headers$get, _response$headers$get2;
	if (customTransport) return customTransport.call(objectID, method, windowName, args);
	const openHarmony = getOpenHarmonyTransport();
	if (openHarmony) return openHarmony.call(objectID, method, windowName, args);
	let url = new URL(runtimeURL());
	let body = {
		object: objectID,
		method
	};
	if (args !== null && args !== void 0) body.args = args;
	let headers = {
		["x-wails-client-id"]: clientId,
		["Content-Type"]: "application/json"
	};
	if (windowName) headers["x-wails-window-name"] = windowName;
	const bodyStr = JSON.stringify(body);
	let response;
	if (bodyStr.length > CHUNK_THRESHOLD$1) response = await sendChunked(url, headers, bodyStr);
	else response = await fetch(url, {
		method: "POST",
		headers,
		body: bodyStr
	});
	if (!response.ok) {
		const ct = response.headers.get("Content-Type");
		if (ct === null || ct === void 0 ? void 0 : ct.includes("application/json")) {
			const json = await response.json();
			let err;
			switch (json.kind) {
				case "ReferenceError":
					err = new ReferenceError(json.message);
					break;
				case "TypeError":
					err = new TypeError(json.message);
					break;
				case "RuntimeError":
					err = new RuntimeError(json.message);
					break;
				default: err = new Error(json.message);
			}
			err.cause = json.cause;
			throw err;
		}
		throw new Error(await response.text());
	}
	if (((_response$headers$get = (_response$headers$get2 = response.headers.get("Content-Type")) === null || _response$headers$get2 === void 0 ? void 0 : _response$headers$get2.indexOf("application/json")) !== null && _response$headers$get !== void 0 ? _response$headers$get : -1) !== -1) return response.json();
	else return response.text();
}
async function sendChunked(url, headers, bodyStr) {
	const chunkId = nanoid();
	const bodyBytes = new TextEncoder().encode(bodyStr);
	const totalChunks = Math.ceil(bodyBytes.length / CHUNK_THRESHOLD$1);
	for (let i = 0; i < totalChunks - 1; i++) {
		const chunk = bodyBytes.subarray(i * CHUNK_THRESHOLD$1, (i + 1) * CHUNK_THRESHOLD$1);
		const resp = await fetch(url, {
			method: "POST",
			headers: _objectSpread2(_objectSpread2({}, headers), {}, {
				"x-wails-chunk-id": chunkId,
				"x-wails-chunk-index": String(i),
				"x-wails-chunk-total": String(totalChunks)
			}),
			body: chunk
		});
		if (!resp.ok) throw new Error(await resp.text());
	}
	return fetch(url, {
		method: "POST",
		headers: _objectSpread2(_objectSpread2({}, headers), {}, {
			"x-wails-chunk-id": chunkId,
			"x-wails-chunk-index": String(totalChunks - 1),
			"x-wails-chunk-total": String(totalChunks)
		}),
		body: bodyBytes.subarray((totalChunks - 1) * CHUNK_THRESHOLD$1)
	});
}
var _window$wails, CHUNK_THRESHOLD$1, RuntimeError, objectNames, clientId, customTransport, detectedOpenHarmonyTransport, androidBridge;
function init_runtime() {
	return (init_runtime = __esmMin((() => {
		init_nanoid();
		init_environment();
		init_objectSpread2();
		;
		CHUNK_THRESHOLD$1 = 524288;
		RuntimeError = class extends Error {
			constructor(message, options) {
				super(message, options);
				this.name = "RuntimeError";
			}
		};
		objectNames = Object.freeze({
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
			CancelCall: 10,
			IOS: 11,
			Android: 12
		});
		clientId = nanoid();
		customTransport = null;
		;
		androidBridge = hasDOM && typeof ((_window$wails = window.wails) === null || _window$wails === void 0 ? void 0 : _window$wails.invokeAsync) === "function" ? window.wails : null;
		if (androidBridge) {
			const pending = /* @__PURE__ */ new Map();
			window._wailsAndroidCallback = (id, response, error) => {
				const promise = pending.get(id);
				if (!promise) return;
				pending.delete(id);
				if (error) {
					promise.reject(new Error(error));
					return;
				}
				try {
					const envelope = JSON.parse(response !== null && response !== void 0 ? response : "{}");
					if (!envelope.ok) {
						var _envelope$error;
						promise.reject(new Error((_envelope$error = envelope.error) !== null && _envelope$error !== void 0 ? _envelope$error : "unknown runtime call error"));
						return;
					}
					promise.resolve("text" in envelope ? envelope.text : envelope.data);
				} catch (e) {
					promise.reject(e);
				}
			};
			customTransport = { call(objectID, method, windowName, args) {
				return new Promise((resolve, reject) => {
					const id = nanoid();
					pending.set(id, {
						resolve,
						reject
					});
					try {
						androidBridge.invokeAsync(id, JSON.stringify({
							object: objectID,
							method,
							windowName,
							args: args !== null && args !== void 0 ? args : null,
							clientId
						}));
					} catch (e) {
						pending.delete(id);
						reject(e);
					}
				});
			} };
		}
	})))();
}

// src/system.ts
var system_exports = /* @__PURE__ */ __exportAll({
	Capabilities: () => Capabilities,
	Environment: () => Environment,
	IsAMD64: () => IsAMD64,
	IsARM: () => IsARM,
	IsARM64: () => IsARM64,
	IsAndroid: () => IsAndroid,
	IsDarkMode: () => IsDarkMode,
	IsDebug: () => IsDebug,
	IsDesktop: () => IsDesktop,
	IsIOS: () => IsIOS,
	IsLinux: () => IsLinux,
	IsMac: () => IsMac,
	IsMobile: () => IsMobile,
	IsWindows: () => IsWindows,
	invoke: () => invoke
});
function invoke(msg) {
	_invoke === null || _invoke === void 0 || _invoke(msg);
}
function IsDarkMode() {
	return call$10(SystemIsDarkMode);
}
async function Capabilities() {
	return call$10(SystemCapabilities);
}
function Environment() {
	return call$10(SystemEnvironment);
}
function IsWindows() {
	var _window$_wails;
	return ((_window$_wails = window._wails) === null || _window$_wails === void 0 || (_window$_wails = _window$_wails.environment) === null || _window$_wails === void 0 ? void 0 : _window$_wails.OS) === "windows";
}
function IsLinux() {
	var _window$_wails2;
	return ((_window$_wails2 = window._wails) === null || _window$_wails2 === void 0 || (_window$_wails2 = _window$_wails2.environment) === null || _window$_wails2 === void 0 ? void 0 : _window$_wails2.OS) === "linux";
}
function IsMac() {
	var _window$_wails3;
	return ((_window$_wails3 = window._wails) === null || _window$_wails3 === void 0 || (_window$_wails3 = _window$_wails3.environment) === null || _window$_wails3 === void 0 ? void 0 : _window$_wails3.OS) === "darwin";
}
function IsIOS() {
	var _window$_wails4;
	return ((_window$_wails4 = window._wails) === null || _window$_wails4 === void 0 || (_window$_wails4 = _window$_wails4.environment) === null || _window$_wails4 === void 0 ? void 0 : _window$_wails4.OS) === "ios";
}
function IsAndroid() {
	var _window$_wails5;
	return ((_window$_wails5 = window._wails) === null || _window$_wails5 === void 0 || (_window$_wails5 = _window$_wails5.environment) === null || _window$_wails5 === void 0 ? void 0 : _window$_wails5.OS) === "android";
}
function IsMobile() {
	return IsIOS() || IsAndroid();
}
function IsDesktop() {
	return IsMac() || IsWindows() || IsLinux();
}
function IsAMD64() {
	var _window$_wails6;
	return ((_window$_wails6 = window._wails) === null || _window$_wails6 === void 0 || (_window$_wails6 = _window$_wails6.environment) === null || _window$_wails6 === void 0 ? void 0 : _window$_wails6.Arch) === "amd64";
}
function IsARM() {
	var _window$_wails7;
	return ((_window$_wails7 = window._wails) === null || _window$_wails7 === void 0 || (_window$_wails7 = _window$_wails7.environment) === null || _window$_wails7 === void 0 ? void 0 : _window$_wails7.Arch) === "arm";
}
function IsARM64() {
	var _window$_wails8;
	return ((_window$_wails8 = window._wails) === null || _window$_wails8 === void 0 || (_window$_wails8 = _window$_wails8.environment) === null || _window$_wails8 === void 0 ? void 0 : _window$_wails8.Arch) === "arm64";
}
function IsDebug() {
	var _window$_wails9;
	return Boolean((_window$_wails9 = window._wails) === null || _window$_wails9 === void 0 || (_window$_wails9 = _window$_wails9.environment) === null || _window$_wails9 === void 0 ? void 0 : _window$_wails9.Debug);
}
var call$10, SystemIsDarkMode, SystemEnvironment, SystemCapabilities, _invoke;
function init_system() {
	return (init_system = __esmMin((() => {
		init_runtime();
		call$10 = newRuntimeCaller(objectNames.System);
		SystemIsDarkMode = 0;
		SystemEnvironment = 1;
		SystemCapabilities = 2;
		_invoke = (function() {
			try {
				var _window$chrome, _window$webkit, _window$wails;
				if ((_window$chrome = window.chrome) === null || _window$chrome === void 0 || (_window$chrome = _window$chrome.webview) === null || _window$chrome === void 0 ? void 0 : _window$chrome.postMessage) return window.chrome.webview.postMessage.bind(window.chrome.webview);
				else if ((_window$webkit = window.webkit) === null || _window$webkit === void 0 || (_window$webkit = _window$webkit.messageHandlers) === null || _window$webkit === void 0 || (_window$webkit = _window$webkit["external"]) === null || _window$webkit === void 0 ? void 0 : _window$webkit.postMessage) return window.webkit.messageHandlers["external"].postMessage.bind(window.webkit.messageHandlers["external"]);
				else if ((_window$wails = window.wails) === null || _window$wails === void 0 ? void 0 : _window$wails.invoke) {
					const bridge = window.wails;
					return (msg) => bridge.invoke(typeof msg === "string" ? msg : JSON.stringify(msg));
				}
			} catch (e) {}
			console.warn("\n%c⚠️ Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n", "background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;", "background: transparent;", "color: #ffffff; font-style: italic; font-weight: bold;");
			return null;
		})();
	})))();
}

// src/utils.ts
function debugLog(message) {
	console.log("%c wails3 %c " + message + " ", "background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem", "background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem");
}
function canTrackButtons() {
	return new MouseEvent("mousedown").buttons === 0;
}
function canAbortListeners() {
	if (!EventTarget || !AbortSignal || !AbortController) return false;
	let result = true;
	const target = new EventTarget();
	const controller = new AbortController();
	target.addEventListener("test", () => {
		result = false;
	}, { signal: controller.signal });
	controller.abort();
	target.dispatchEvent(new CustomEvent("test"));
	return result;
}
function eventTarget(event) {
	if (event.target instanceof HTMLElement) return event.target;
	else if (!(event.target instanceof HTMLElement) && event.target instanceof Node) {
		var _event$target$parentE;
		return (_event$target$parentE = event.target.parentElement) !== null && _event$target$parentE !== void 0 ? _event$target$parentE : document.body;
	} else return document.body;
}
function whenReady(callback) {
	if (isReady || document.readyState === "complete") callback();
	else document.addEventListener("DOMContentLoaded", callback);
}
var isReady;
function init_utils() {
	return (init_utils = __esmMin((() => {
		init_environment();
		isReady = false;
		if (hasDOM) document.addEventListener("DOMContentLoaded", () => {
			isReady = true;
		});
	})))();
}

// src/contextmenu.ts
function openContextMenu(id, x, y, data) {
	call$9(ContextMenuOpen, {
		id,
		x,
		y,
		data
	});
}
function contextMenuHandler(event) {
	const target = eventTarget(event);
	const customContextMenu = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu").trim();
	if (customContextMenu) {
		event.preventDefault();
		const data = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu-data");
		openContextMenu(customContextMenu, event.clientX, event.clientY, data);
	} else processDefaultContextMenu(event, target);
}
function processDefaultContextMenu(event, target) {
	if (IsDebug()) return;
	switch (window.getComputedStyle(target).getPropertyValue("--default-contextmenu").trim()) {
		case "show": return;
		case "hide":
			event.preventDefault();
			return;
	}
	if (target.isContentEditable) return;
	const selection = window.getSelection();
	const hasSelection = selection && selection.toString().length > 0;
	if (hasSelection) for (let i = 0; i < selection.rangeCount; i++) {
		const rects = selection.getRangeAt(i).getClientRects();
		for (let j = 0; j < rects.length; j++) {
			const rect = rects[j];
			if (document.elementFromPoint(rect.left, rect.top) === target) return;
		}
	}
	if (target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement) {
		if (hasSelection || !target.readOnly && !target.disabled) return;
	}
	event.preventDefault();
}
var call$9, ContextMenuOpen;
function init_contextmenu() {
	return (init_contextmenu = __esmMin((() => {
		init_runtime();
		init_system();
		init_utils();
		init_environment();
		if (hasDOM) window.addEventListener("contextmenu", contextMenuHandler);
		call$9 = newRuntimeCaller(objectNames.ContextMenu);
		ContextMenuOpen = 0;
	})))();
}

// src/flags.ts
var flags_exports = /* @__PURE__ */ __exportAll({ GetFlag: () => GetFlag });
function GetFlag(key) {
	try {
		return window._wails.flags[key];
	} catch (e) {
		throw new Error("Unable to retrieve flag '" + key + "': " + e, { cause: e });
	}
}
function init_flags() {
	return (init_flags = __esmMin((() => {})))();
}

// src/drag.ts
function isMobile() {
	var _window$_wails;
	const os = (_window$_wails = window._wails) === null || _window$_wails === void 0 || (_window$_wails = _window$_wails.environment) === null || _window$_wails === void 0 ? void 0 : _window$_wails.OS;
	if (os === "ios" || os === "android") return true;
	const ua = navigator.userAgent || navigator.vendor || window.opera || "";
	return /android|iphone|ipad|ipod|iemobile|wpdesktop/i.test(ua);
}
function tryInitDragHandlers() {
	if (dragInitDone) return;
	if (isMobile()) return;
	window.addEventListener("mousedown", update, { capture: true });
	window.addEventListener("mousemove", update, { capture: true });
	window.addEventListener("mouseup", update, { capture: true });
	for (const ev of [
		"click",
		"contextmenu",
		"dblclick"
	]) window.addEventListener(ev, suppressEvent, { capture: true });
	dragInitDone = true;
}
function suppressEvent(event) {
	if (event.type === "dblclick" && IsMac() && window === window.top && isDraggableEvent(event)) {
		event.stopImmediatePropagation();
		event.stopPropagation();
		event.preventDefault();
		invoke("wails:drag:doubleclick");
		return;
	}
	if (dragging || resizing) {
		event.stopImmediatePropagation();
		event.stopPropagation();
		event.preventDefault();
	}
}
function isDraggableEvent(event) {
	const target = eventTarget(event);
	return window.getComputedStyle(target).getPropertyValue("--wails-draggable").trim() === "drag" && event.offsetX >= 0 && event.offsetX < target.clientWidth && event.offsetY >= 0 && event.offsetY < target.clientHeight;
}
function update(event) {
	let eventType, eventButtons = event.buttons;
	switch (event.type) {
		case "mousedown":
			eventType = MouseDown;
			if (!buttonsTracked) eventButtons = buttons | 1 << event.button;
			break;
		case "mouseup":
			eventType = MouseUp;
			if (!buttonsTracked) eventButtons = buttons & ~(1 << event.button);
			break;
		default:
			eventType = MouseMove;
			if (!buttonsTracked) eventButtons = buttons;
	}
	let released = buttons & ~eventButtons;
	let pressed = eventButtons & ~buttons;
	buttons = eventButtons;
	if (eventType === MouseDown && !(pressed & event.button)) {
		released |= 1 << event.button;
		pressed |= 1 << event.button;
	}
	if (eventType !== MouseMove && resizing || dragging && (eventType === MouseDown || event.button !== 0)) {
		event.stopImmediatePropagation();
		event.stopPropagation();
		event.preventDefault();
	}
	if (released & 1) primaryUp(event);
	if (pressed & 1) primaryDown(event);
	if (eventType === MouseMove) onMouseMove(event);
}
function primaryDown(event) {
	canDrag = false;
	canResize = false;
	if (!IsWindows()) {
		if (event.type === "mousedown" && event.button === 0 && event.detail !== 1) return;
	}
	if (resizeEdge) {
		if (event.type !== "mousedown") return;
		canResize = true;
		return;
	}
	canDrag = isDraggableEvent(event);
}
function primaryUp(event) {
	canDrag = false;
	dragging = false;
	canResize = false;
	resizing = false;
}
function setResize(edge) {
	if (edge) {
		if (!resizeEdge) defaultCursor = document.body.style.cursor;
		document.body.style.cursor = cursorForEdge[edge];
	} else if (!edge && resizeEdge) document.body.style.cursor = defaultCursor;
	resizeEdge = edge || "";
}
function onMouseMove(event) {
	if (canResize && resizeEdge) {
		resizing = true;
		invoke("wails:resize:" + resizeEdge);
	} else if (canDrag) {
		dragging = true;
		invoke("wails:drag");
	}
	if (dragging || resizing) {
		canDrag = canResize = false;
		return;
	}
	if (!resizable || !IsWindows() && !(IsLinux() && GetFlag("frameless"))) {
		if (resizeEdge) setResize();
		return;
	}
	const resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
	const resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
	const cornerExtra = GetFlag("resizeCornerExtra") || 10;
	const scrollbarWidth = Math.max(0, window.innerWidth - document.documentElement.clientWidth);
	const scrollbarHeight = Math.max(0, window.innerHeight - document.documentElement.clientHeight);
	const rightContentEdge = window.innerWidth - scrollbarWidth;
	const bottomContentEdge = window.innerHeight - scrollbarHeight;
	const rightBorder = event.clientX < rightContentEdge && rightContentEdge - event.clientX < resizeHandleWidth;
	const leftBorder = event.clientX < resizeHandleWidth;
	const topBorder = event.clientY < resizeHandleHeight;
	const bottomBorder = event.clientY < bottomContentEdge && bottomContentEdge - event.clientY < resizeHandleHeight;
	const rightCorner = event.clientX < rightContentEdge && rightContentEdge - event.clientX < resizeHandleWidth + cornerExtra;
	const leftCorner = event.clientX < resizeHandleWidth + cornerExtra;
	const topCorner = event.clientY < resizeHandleHeight + cornerExtra;
	const bottomCorner = event.clientY < bottomContentEdge && bottomContentEdge - event.clientY < resizeHandleHeight + cornerExtra;
	if (!leftCorner && !topCorner && !bottomCorner && !rightCorner) setResize();
	else if (rightCorner && bottomCorner) setResize("se-resize");
	else if (leftCorner && bottomCorner) setResize("sw-resize");
	else if (leftCorner && topCorner) setResize("nw-resize");
	else if (topCorner && rightCorner) setResize("ne-resize");
	else if (leftBorder) setResize("w-resize");
	else if (topBorder) setResize("n-resize");
	else if (bottomBorder) setResize("s-resize");
	else if (rightBorder) setResize("e-resize");
	else setResize();
}
var canDrag, dragging, resizable, canResize, resizing, resizeEdge, defaultCursor, buttons, buttonsTracked, dragInitDone, MouseDown, MouseUp, MouseMove, cursorForEdge;
function init_drag() {
	return (init_drag = __esmMin((() => {
		init_system();
		/* @__PURE__ */ init_flags();
		init_utils();
		init_environment();
		canDrag = false;
		dragging = false;
		resizable = false;
		canResize = false;
		resizing = false;
		resizeEdge = "";
		defaultCursor = "auto";
		buttons = 0;
		buttonsTracked = false;
		if (hasDOM) {
			buttonsTracked = canTrackButtons();
			window._wails = window._wails || {};
			window._wails.setResizable = (value) => {
				resizable = value;
				if (!resizable) {
					canResize = resizing = false;
					setResize();
				}
			};
		}
		dragInitDone = false;
		if (hasDOM) {
			tryInitDragHandlers();
			document.addEventListener("DOMContentLoaded", tryInitDragHandlers, { once: true });
			let dragEnvPolls = 0;
			const dragEnvPoll = window.setInterval(() => {
				if (dragInitDone) {
					window.clearInterval(dragEnvPoll);
					return;
				}
				tryInitDragHandlers();
				if (++dragEnvPolls > 100) window.clearInterval(dragEnvPoll);
			}, 50);
		}
		MouseDown = 0;
		MouseUp = 1;
		MouseMove = 2;
		cursorForEdge = Object.freeze({
			"se-resize": "nwse-resize",
			"sw-resize": "nesw-resize",
			"nw-resize": "nwse-resize",
			"ne-resize": "nesw-resize",
			"w-resize": "ew-resize",
			"n-resize": "ns-resize",
			"s-resize": "ns-resize",
			"e-resize": "ew-resize"
		});
	})))();
}

// src/appregion.ts
function normaliseRegionKind(value) {
	const region = value.trim().toLowerCase();
	if (validRegions.has(region)) return region;
}
function nonClientRegionForElement(element) {
	if (!(element instanceof HTMLElement)) return;
	const region = normaliseRegionKind(window.getComputedStyle(element).getPropertyValue(regionProperty));
	if (!region) return;
	const parent = element.parentElement;
	if (parent) {
		if (normaliseRegionKind(window.getComputedStyle(parent).getPropertyValue(regionProperty)) === region) return;
	}
	return region;
}
function isVisible(element) {
	const style = window.getComputedStyle(element);
	return style.display !== "none" && style.visibility !== "hidden" && style.contentVisibility !== "hidden";
}
function elementRegion(element) {
	if (!(element instanceof HTMLElement)) return;
	const kind = nonClientRegionForElement(element);
	if (!kind || !isVisible(element)) return;
	const rect = element.getBoundingClientRect();
	if (rect.width <= 0 || rect.height <= 0) return;
	const scale = window.devicePixelRatio || 1;
	const left = Math.floor(rect.left * scale);
	const top = Math.floor(rect.top * scale);
	const right = Math.ceil(rect.right * scale);
	const bottom = Math.ceil(rect.bottom * scale);
	if (right <= left || bottom <= top) return;
	return {
		kind,
		left,
		top,
		right,
		bottom
	};
}
function regionElements() {
	const elements = [];
	if (document.documentElement) elements.push(document.documentElement);
	if (document.body) {
		elements.push(document.body);
		for (const element of document.body.querySelectorAll("*")) elements.push(element);
	}
	return elements;
}
function observeRegionElements(elements) {
	var _resizeObserver;
	if (typeof ResizeObserver === "undefined") return;
	(_resizeObserver = resizeObserver) !== null && _resizeObserver !== void 0 || (resizeObserver = new ResizeObserver(scheduleUpdate));
	const nextElements = new Set(elements);
	for (const element of observedElements) if (!nextElements.has(element)) resizeObserver.unobserve(element);
	for (const element of nextElements) if (!observedElements.has(element)) resizeObserver.observe(element);
	observedElements = nextElements;
}
function updateNonClientRegions() {
	updatePending = false;
	const elements = regionElements();
	const regions = [];
	const activeElements = [];
	for (const element of elements) {
		const region = elementRegion(element);
		if (region) {
			regions.push(region);
			activeElements.push(element);
		}
	}
	observeRegionElements(activeElements);
	const payload = JSON.stringify({
		version: 1,
		regions
	});
	if (payload === lastPayload) return;
	lastPayload = payload;
	invoke("wails:non-client-region:" + payload);
}
function scheduleUpdate() {
	if (updatePending) return;
	updatePending = true;
	window.requestAnimationFrame(updateNonClientRegions);
}
function startNonClientRegionTracking() {
	var _window$visualViewpor, _window$visualViewpor2;
	if (trackingStarted) return;
	trackingStarted = true;
	scheduleUpdate();
	new MutationObserver(scheduleUpdate).observe(document.documentElement, {
		attributes: true,
		childList: true,
		subtree: true
	});
	window.addEventListener("resize", scheduleUpdate);
	window.addEventListener("scroll", scheduleUpdate, true);
	(_window$visualViewpor = window.visualViewport) === null || _window$visualViewpor === void 0 || _window$visualViewpor.addEventListener("resize", scheduleUpdate);
	(_window$visualViewpor2 = window.visualViewport) === null || _window$visualViewpor2 === void 0 || _window$visualViewpor2.addEventListener("scroll", scheduleUpdate);
}
function tryStartNonClientRegionTracking() {
	var _window$_wails$enviro, _window$_wails$flags;
	const os = (_window$_wails$enviro = window._wails.environment) === null || _window$_wails$enviro === void 0 ? void 0 : _window$_wails$enviro.OS;
	if (os === void 0) return false;
	const enabled = (_window$_wails$flags = window._wails.flags) === null || _window$_wails$flags === void 0 ? void 0 : _window$_wails$flags.nonClientRegionTracking;
	if (os === "windows") {
		if (enabled === true) whenReady(startNonClientRegionTracking);
		return true;
	}
	return true;
}
var regionProperty, runtimeConfigReadyEvent, validRegions, updatePending, lastPayload, observedElements, resizeObserver, trackingStarted;
function init_appregion() {
	return (init_appregion = __esmMin((() => {
		init_system();
		init_utils();
		init_environment();
		regionProperty = "--wails-non-client-region";
		runtimeConfigReadyEvent = "wails:runtime-config-ready";
		validRegions = /* @__PURE__ */ new Set([
			"caption",
			"minimize",
			"maximize",
			"close"
		]);
		if (hasDOM) window._wails = window._wails || {};
		updatePending = false;
		lastPayload = "";
		observedElements = /* @__PURE__ */ new Set();
		;
		trackingStarted = false;
		if (hasDOM && !tryStartNonClientRegionTracking()) window.addEventListener(runtimeConfigReadyEvent, tryStartNonClientRegionTracking, { once: true });
	})))();
}

// src/application.ts
var application_exports = /* @__PURE__ */ __exportAll({
	Hide: () => Hide,
	Quit: () => Quit,
	Show: () => Show
});
function Hide() {
	return call$8(HideMethod$1);
}
function Show() {
	return call$8(ShowMethod$1);
}
function Quit() {
	return call$8(QuitMethod);
}
var call$8, HideMethod$1, ShowMethod$1, QuitMethod;
function init_application() {
	return (init_application = __esmMin((() => {
		init_runtime();
		call$8 = newRuntimeCaller(objectNames.Application);
		HideMethod$1 = 0;
		ShowMethod$1 = 1;
		QuitMethod = 2;
	})))();
}

// src/browser.ts
var browser_exports = /* @__PURE__ */ __exportAll({ OpenURL: () => OpenURL });
function OpenURL(url) {
	return call$7(BrowserOpenURL, { url: url.toString() });
}
var call$7, BrowserOpenURL;
function init_browser() {
	return (init_browser = __esmMin((() => {
		init_runtime();
		call$7 = newRuntimeCaller(objectNames.Browser);
		BrowserOpenURL = 0;
	})))();
}

// src/callable.ts
function isCallableRefApply(value) {
	if (isDDA(value)) return true;
	if (!value) return false;
	if (typeof value !== "function" && typeof value !== "object") return false;
	try {
		reflectApply(value, null, badArrayLike);
	} catch (e) {
		if (e !== isCallableMarker) return false;
	}
	return !isES6ClassFn(value) && tryFunctionObject(value);
}
function isCallableNoRefApply(value) {
	if (isDDA(value)) return true;
	if (!value) return false;
	if (typeof value !== "function" && typeof value !== "object") return false;
	if (hasToStringTag) return tryFunctionObject(value);
	if (isES6ClassFn(value)) return false;
	var strClass = toStr.call(value);
	if (strClass !== fnClass && strClass !== genClass && !/^\[object HTML/.test(strClass)) return false;
	return tryFunctionObject(value);
}
var fnToStr, reflectApply, badArrayLike, isCallableMarker, constructorRegex, isES6ClassFn, tryFunctionObject, toStr, objectClass, fnClass, genClass, ddaClass, ddaClass2, ddaClass3, hasToStringTag, isIE68, isDDA, all, callable_default;
function init_callable() {
	return (init_callable = __esmMin((() => {
		fnToStr = Function.prototype.toString;
		reflectApply = typeof Reflect === "object" && Reflect !== null && Reflect.apply;
		;
		;
		if (typeof reflectApply === "function" && typeof Object.defineProperty === "function") try {
			badArrayLike = Object.defineProperty({}, "length", { get: function() {
				throw isCallableMarker;
			} });
			isCallableMarker = {};
			reflectApply(function() {
				throw 42;
			}, null, badArrayLike);
		} catch (_) {
			if (_ !== isCallableMarker) reflectApply = null;
		}
		else reflectApply = null;
		constructorRegex = /^\s*class\b/;
		isES6ClassFn = function isES6ClassFunction(value) {
			try {
				var fnStr = fnToStr.call(value);
				return constructorRegex.test(fnStr);
			} catch (e) {
				return false;
			}
		};
		tryFunctionObject = function tryFunctionToStr(value) {
			try {
				if (isES6ClassFn(value)) return false;
				fnToStr.call(value);
				return true;
			} catch (e) {
				return false;
			}
		};
		toStr = Object.prototype.toString;
		objectClass = "[object Object]";
		fnClass = "[object Function]";
		genClass = "[object GeneratorFunction]";
		ddaClass = "[object HTMLAllCollection]";
		ddaClass2 = "[object HTML document.all class]";
		ddaClass3 = "[object HTMLCollection]";
		hasToStringTag = typeof Symbol === "function" && !!Symbol.toStringTag;
		isIE68 = !(0 in [,]);
		isDDA = function isDocumentDotAll() {
			return false;
		};
		if (typeof document === "object") {
			all = document.all;
			if (toStr.call(all) === toStr.call(document.all)) isDDA = function isDocumentDotAll(value) {
				if ((isIE68 || !value) && (typeof value === "undefined" || typeof value === "object")) try {
					var str = toStr.call(value);
					return (str === ddaClass || str === ddaClass2 || str === ddaClass3 || str === objectClass) && value("") == null;
				} catch (e) {}
				return false;
			};
		}
		callable_default = reflectApply ? isCallableRefApply : isCallableNoRefApply;
	})))();
}

// src/cancellable.ts
function cancellerFor(promise, state) {
	let cancellationPromise = void 0;
	return (reason) => {
		if (!state.settled) {
			state.settled = true;
			state.reason = reason;
			promise.reject(reason);
			Promise.prototype.then.call(promise.promise, void 0, (err) => {
				if (err !== reason) throw err;
			});
		}
		if (!state.reason || !promise.oncancelled) return;
		cancellationPromise = new Promise((resolve) => {
			try {
				resolve(promise.oncancelled(state.reason.cause));
			} catch (err) {
				Promise.reject(new CancelledRejectionError(promise.promise, err, "Unhandled exception in oncancelled callback."));
			}
		}).catch((reason) => {
			Promise.reject(new CancelledRejectionError(promise.promise, reason, "Unhandled rejection in oncancelled callback."));
		});
		promise.oncancelled = null;
		return cancellationPromise;
	};
}
function resolverFor(promise, state) {
	return (value) => {
		if (state.resolving) return;
		state.resolving = true;
		if (value === promise.promise) {
			if (state.settled) return;
			state.settled = true;
			promise.reject(/* @__PURE__ */ new TypeError("A promise cannot be resolved with itself."));
			return;
		}
		if (value != null && (typeof value === "object" || typeof value === "function")) {
			let then;
			try {
				then = value.then;
			} catch (err) {
				state.settled = true;
				promise.reject(err);
				return;
			}
			if (callable_default(then)) {
				try {
					let cancel = value.cancel;
					if (callable_default(cancel)) {
						const oncancelled = (cause) => {
							Reflect.apply(cancel, value, [cause]);
						};
						if (state.reason) cancellerFor(_objectSpread2(_objectSpread2({}, promise), {}, { oncancelled }), state)(state.reason);
						else promise.oncancelled = oncancelled;
					}
				} catch (_unused) {}
				const newState = {
					root: state.root,
					resolving: false,
					get settled() {
						return this.root.settled;
					},
					set settled(value) {
						this.root.settled = value;
					},
					get reason() {
						return this.root.reason;
					}
				};
				const rejector = rejectorFor(promise, newState);
				try {
					Reflect.apply(then, value, [resolverFor(promise, newState), rejector]);
				} catch (err) {
					rejector(err);
				}
				return;
			}
		}
		if (state.settled) return;
		state.settled = true;
		promise.resolve(value);
	};
}
function rejectorFor(promise, state) {
	return (reason) => {
		if (state.resolving) return;
		state.resolving = true;
		if (state.settled) {
			try {
				if (reason instanceof CancelError && state.reason instanceof CancelError && Object.is(reason.cause, state.reason.cause)) return;
			} catch (_unused2) {}
			Promise.reject(new CancelledRejectionError(promise.promise, reason));
		} else {
			state.settled = true;
			promise.reject(reason);
		}
	};
}
function cancelAll(parent, values, cause) {
	const results = [];
	for (const value of values) {
		let cancel;
		try {
			if (!callable_default(value.then)) continue;
			cancel = value.cancel;
			if (!callable_default(cancel)) continue;
		} catch (_unused3) {
			continue;
		}
		let result;
		try {
			result = Reflect.apply(cancel, value, [cause]);
		} catch (err) {
			Promise.reject(new CancelledRejectionError(parent, err, "Unhandled exception in cancel method."));
			continue;
		}
		if (!result) continue;
		results.push((result instanceof Promise ? result : Promise.resolve(result)).catch((reason) => {
			Promise.reject(new CancelledRejectionError(parent, reason, "Unhandled rejection in cancel method."));
		}));
	}
	return Promise.all(results);
}
function identity(x) {
	return x;
}
function thrower(reason) {
	throw reason;
}
function errorMessage(err) {
	try {
		if (err instanceof Error || typeof err !== "object" || err.toString !== Object.prototype.toString) return "" + err;
	} catch (_unused4) {}
	try {
		return JSON.stringify(err);
	} catch (_unused5) {}
	try {
		return Object.prototype.toString.call(err);
	} catch (_unused6) {}
	return "<could not convert error to string>";
}
function currentBarrier(promise) {
	var _promise$barrierSym;
	let pwr = (_promise$barrierSym = promise[barrierSym]) !== null && _promise$barrierSym !== void 0 ? _promise$barrierSym : {};
	if (!("promise" in pwr)) Object.assign(pwr, promiseWithResolvers());
	if (promise[barrierSym] == null) {
		pwr.resolve();
		promise[barrierSym] = pwr;
	}
	return pwr.promise;
}
var _Symbol$species, CancelError, CancelledRejectionError, barrierSym, cancelImplSym, species, CancellablePromise, promiseWithResolvers;
function init_cancellable() {
	return (init_cancellable = __esmMin((() => {
		init_callable();
		init_objectSpread2();
		;
		CancelError = class extends Error {
			constructor(message, options) {
				super(message, options);
				this.name = "CancelError";
			}
		};
		CancelledRejectionError = class extends Error {
			constructor(promise, reason, info) {
				super((info !== null && info !== void 0 ? info : "Unhandled rejection in cancelled promise.") + " Reason: " + errorMessage(reason), { cause: reason });
				this.promise = promise;
				this.name = "CancelledRejectionError";
			}
		};
		barrierSym = Symbol("barrier");
		cancelImplSym = Symbol("cancelImpl");
		species = (_Symbol$species = Symbol.species) !== null && _Symbol$species !== void 0 ? _Symbol$species : Symbol("speciesPolyfill");
		CancellablePromise = class CancellablePromise extends Promise {
			constructor(executor, oncancelled) {
				let resolve;
				let reject;
				super((res, rej) => {
					resolve = res;
					reject = rej;
				});
				if (this.constructor[species] !== Promise) throw new TypeError("CancellablePromise does not support transparent subclassing. Please refrain from overriding the [Symbol.species] static property.");
				let promise = {
					promise: this,
					resolve,
					reject,
					get oncancelled() {
						var _oncancelled;
						return (_oncancelled = oncancelled) !== null && _oncancelled !== void 0 ? _oncancelled : null;
					},
					set oncancelled(cb) {
						oncancelled = cb !== null && cb !== void 0 ? cb : void 0;
					}
				};
				const state = {
					get root() {
						return state;
					},
					resolving: false,
					settled: false
				};
				Object.defineProperties(this, {
					[barrierSym]: {
						configurable: false,
						enumerable: false,
						writable: true,
						value: null
					},
					[cancelImplSym]: {
						configurable: false,
						enumerable: false,
						writable: false,
						value: cancellerFor(promise, state)
					}
				});
				const rejector = rejectorFor(promise, state);
				try {
					executor(resolverFor(promise, state), rejector);
				} catch (err) {
					if (state.resolving) console.log("Unhandled exception in CancellablePromise executor.", err);
					else rejector(err);
				}
			}
			cancel(cause) {
				return new CancellablePromise((resolve) => {
					Promise.all([this[cancelImplSym](new CancelError("Promise cancelled.", { cause })), currentBarrier(this)]).then(() => resolve(), () => resolve());
				});
			}
			cancelOn(signal) {
				if (signal.aborted) this.cancel(signal.reason);
				else signal.addEventListener("abort", () => void this.cancel(signal.reason), { capture: true });
				return this;
			}
			then(onfulfilled, onrejected, oncancelled) {
				if (!(this instanceof CancellablePromise)) throw new TypeError("CancellablePromise.prototype.then called on an invalid object.");
				if (!callable_default(onfulfilled)) onfulfilled = identity;
				if (!callable_default(onrejected)) onrejected = thrower;
				if (onfulfilled === identity && onrejected == thrower) return new CancellablePromise((resolve) => resolve(this));
				const barrier = {};
				this[barrierSym] = barrier;
				return new CancellablePromise((resolve, reject) => {
					super.then((value) => {
						var _barrier$resolve;
						if (this[barrierSym] === barrier) this[barrierSym] = null;
						(_barrier$resolve = barrier.resolve) === null || _barrier$resolve === void 0 || _barrier$resolve.call(barrier);
						try {
							resolve(onfulfilled(value));
						} catch (err) {
							reject(err);
						}
					}, (reason) => {
						var _barrier$resolve2;
						if (this[barrierSym] === barrier) this[barrierSym] = null;
						(_barrier$resolve2 = barrier.resolve) === null || _barrier$resolve2 === void 0 || _barrier$resolve2.call(barrier);
						try {
							resolve(onrejected(reason));
						} catch (err) {
							reject(err);
						}
					});
				}, async (cause) => {
					try {
						return oncancelled === null || oncancelled === void 0 ? void 0 : oncancelled(cause);
					} finally {
						await this.cancel(cause);
					}
				});
			}
			catch(onrejected, oncancelled) {
				return this.then(void 0, onrejected, oncancelled);
			}
			finally(onfinally, oncancelled) {
				if (!(this instanceof CancellablePromise)) throw new TypeError("CancellablePromise.prototype.finally called on an invalid object.");
				if (!callable_default(onfinally)) return this.then(onfinally, onfinally, oncancelled);
				return this.then((value) => CancellablePromise.resolve(onfinally()).then(() => value), (reason) => CancellablePromise.resolve(onfinally()).then(() => {
					throw reason;
				}), oncancelled);
			}
			static get [species]() {
				return Promise;
			}
			static all(values) {
				let collected = Array.from(values);
				const promise = collected.length === 0 ? CancellablePromise.resolve(collected) : new CancellablePromise((resolve, reject) => {
					Promise.all(collected).then(resolve, reject);
				}, (cause) => cancelAll(promise, collected, cause));
				return promise;
			}
			static allSettled(values) {
				let collected = Array.from(values);
				const promise = collected.length === 0 ? CancellablePromise.resolve(collected) : new CancellablePromise((resolve, reject) => {
					Promise.allSettled(collected).then(resolve, reject);
				}, (cause) => cancelAll(promise, collected, cause));
				return promise;
			}
			static any(values) {
				let collected = Array.from(values);
				const promise = collected.length === 0 ? CancellablePromise.resolve(collected) : new CancellablePromise((resolve, reject) => {
					Promise.any(collected).then(resolve, reject);
				}, (cause) => cancelAll(promise, collected, cause));
				return promise;
			}
			static race(values) {
				let collected = Array.from(values);
				const promise = new CancellablePromise((resolve, reject) => {
					Promise.race(collected).then(resolve, reject);
				}, (cause) => cancelAll(promise, collected, cause));
				return promise;
			}
			static cancel(cause) {
				const p = new CancellablePromise(() => {});
				p.cancel(cause);
				return p;
			}
			static timeout(milliseconds, cause) {
				const promise = new CancellablePromise(() => {});
				if (AbortSignal && typeof AbortSignal === "function" && AbortSignal.timeout && typeof AbortSignal.timeout === "function") AbortSignal.timeout(milliseconds).addEventListener("abort", () => void promise.cancel(cause));
				else setTimeout(() => void promise.cancel(cause), milliseconds);
				return promise;
			}
			static sleep(milliseconds, value) {
				return new CancellablePromise((resolve) => {
					setTimeout(() => resolve(value), milliseconds);
				});
			}
			static reject(reason) {
				return new CancellablePromise((_, reject) => reject(reason));
			}
			static resolve(value) {
				if (value instanceof CancellablePromise) return value;
				return new CancellablePromise((resolve) => resolve(value));
			}
			static withResolvers() {
				let result = { oncancelled: null };
				result.promise = new CancellablePromise((resolve, reject) => {
					result.resolve = resolve;
					result.reject = reject;
				}, (cause) => {
					var _result$oncancelled;
					(_result$oncancelled = result.oncancelled) === null || _result$oncancelled === void 0 || _result$oncancelled.call(result, cause);
				});
				return result;
			}
		};
		promiseWithResolvers = Promise.withResolvers;
		if (promiseWithResolvers && typeof promiseWithResolvers === "function") promiseWithResolvers = promiseWithResolvers.bind(Promise);
		else promiseWithResolvers = function() {
			let resolve;
			let reject;
			return {
				promise: new Promise((res, rej) => {
					resolve = res;
					reject = rej;
				}),
				resolve,
				reject
			};
		};
	})))();
}

// src/calls.ts
var calls_exports = /* @__PURE__ */ __exportAll({
	ByID: () => ByID,
	ByName: () => ByName,
	Call: () => Call,
	RuntimeError: () => RuntimeError
});
function generateID() {
	let result;
	do
		result = nanoid();
	while (callResponses.has(result));
	return result;
}
function Call(options) {
	const id = generateID();
	const result = CancellablePromise.withResolvers();
	callResponses.set(id, {
		resolve: result.resolve,
		reject: result.reject
	});
	const request = call$6(CallBinding, Object.assign({ "call-id": id }, options));
	let running = true;
	request.then((res) => {
		running = false;
		callResponses.delete(id);
		result.resolve(res);
	}, (err) => {
		running = false;
		callResponses.delete(id);
		result.reject(err);
	});
	const cancel = () => {
		callResponses.delete(id);
		return cancelCall(CancelMethod, { "call-id": id }).catch((err) => {
			console.error("Error while requesting binding call cancellation:", err);
		});
	};
	result.oncancelled = () => {
		if (running) return cancel();
		else return request.then(cancel);
	};
	return result.promise;
}
function ByName(methodName, ...args) {
	return Call({
		methodName,
		args
	});
}
function ByID(methodID, ...args) {
	return Call({
		methodID,
		args
	});
}
var call$6, cancelCall, callResponses, CallBinding, CancelMethod;
function init_calls() {
	return (init_calls = __esmMin((() => {
		init_cancellable();
		init_runtime();
		init_nanoid();
		init_environment();
		if (hasDOM) window._wails = window._wails || {};
		call$6 = newRuntimeCaller(objectNames.Call);
		cancelCall = newRuntimeCaller(objectNames.CancelCall);
		callResponses = /* @__PURE__ */ new Map();
		CallBinding = 0;
		CancelMethod = 0;
	})))();
}

// src/clipboard.ts
var clipboard_exports = /* @__PURE__ */ __exportAll({
	SetText: () => SetText,
	Text: () => Text
});
function SetText(text) {
	return call$5(ClipboardSetText, { text });
}
function Text() {
	return call$5(ClipboardText);
}
var call$5, ClipboardSetText, ClipboardText;
function init_clipboard() {
	return (init_clipboard = __esmMin((() => {
		init_runtime();
		call$5 = newRuntimeCaller(objectNames.Clipboard);
		ClipboardSetText = 0;
		ClipboardText = 1;
	})))();
}

// src/create.ts
var create_exports = /* @__PURE__ */ __exportAll({
	Any: () => Any,
	Array: () => Array$1,
	ByteSlice: () => ByteSlice,
	DateFromTime: () => DateFromTime,
	Events: () => Events$1,
	Map: () => Map$1,
	Nullable: () => Nullable,
	Struct: () => Struct
});
function Any(source) {
	return source;
}
function ByteSlice(source) {
	return source == null ? "" : source;
}
function Array$1(element) {
	if (element === Any) return (source) => source === null ? [] : source;
	return (source) => {
		if (source === null) return [];
		for (let i = 0; i < source.length; i++) source[i] = element(source[i]);
		return source;
	};
}
function Map$1(key, value) {
	if (value === Any) return (source) => source === null ? {} : source;
	return (source) => {
		if (source === null) return {};
		for (const key in source) source[key] = value(source[key]);
		return source;
	};
}
function Nullable(element) {
	if (element === Any) return Any;
	return (source) => source === null ? null : element(source);
}
function Struct(createField) {
	let allAny = true;
	for (const name in createField) if (createField[name] !== Any) {
		allAny = false;
		break;
	}
	if (allAny) return Any;
	return (source) => {
		for (const name in createField) if (name in source) source[name] = createField[name](source[name]);
		return source;
	};
}
function DateFromTime(source) {
	return new Date(source);
}
var Events$1;
function init_create() {
	return (init_create = __esmMin((() => {
		Events$1 = {};
	})))();
}

// src/dialogs.ts
var dialogs_exports = /* @__PURE__ */ __exportAll({
	Error: () => Error$1,
	Info: () => Info,
	OpenFile: () => OpenFile,
	Question: () => Question,
	SaveFile: () => SaveFile,
	Warning: () => Warning
});
function dialog(type, options = {}) {
	return call$4(type, options);
}
function Info(options) {
	return dialog(DialogInfo, options);
}
function Warning(options) {
	return dialog(DialogWarning, options);
}
function Error$1(options) {
	return dialog(DialogError, options);
}
function Question(options) {
	return dialog(DialogQuestion, options);
}
function OpenFile(options) {
	var _dialog;
	return (_dialog = dialog(DialogOpenFile, options)) !== null && _dialog !== void 0 ? _dialog : [];
}
function SaveFile(options) {
	return dialog(DialogSaveFile, options);
}
var call$4, DialogInfo, DialogWarning, DialogError, DialogQuestion, DialogOpenFile, DialogSaveFile;
function init_dialogs() {
	return (init_dialogs = __esmMin((() => {
		init_runtime();
		init_environment();
		if (hasDOM) window._wails = window._wails || {};
		call$4 = newRuntimeCaller(objectNames.Dialog);
		DialogInfo = 0;
		DialogWarning = 1;
		DialogError = 2;
		DialogQuestion = 3;
		DialogOpenFile = 4;
		DialogSaveFile = 5;
	})))();
}

// src/listener.ts
function listenerOff(listener) {
	let listeners = eventListeners.get(listener.eventName);
	if (!listeners) return;
	listeners = listeners.filter((l) => l !== listener);
	if (listeners.length === 0) eventListeners.delete(listener.eventName);
	else eventListeners.set(listener.eventName, listeners);
}
var eventListeners, Listener;
function init_listener() {
	return (init_listener = __esmMin((() => {
		eventListeners = /* @__PURE__ */ new Map();
		Listener = class {
			constructor(eventName, callback, maxCallbacks) {
				this.eventName = eventName;
				this.callback = callback;
				this.maxCallbacks = maxCallbacks || -1;
			}
			dispatch(data) {
				try {
					this.callback(data);
				} catch (err) {
					console.error(err);
				}
				if (this.maxCallbacks === -1) return false;
				this.maxCallbacks -= 1;
				return this.maxCallbacks === 0;
			}
		};
	})))();
}

// src/event_types.ts
var Types;
function init_event_types() {
	return (init_event_types = __esmMin((() => {
		Types = Object.freeze({
			Windows: Object.freeze({
				APMPowerSettingChange: "windows:APMPowerSettingChange",
				APMPowerStatusChange: "windows:APMPowerStatusChange",
				APMResumeAutomatic: "windows:APMResumeAutomatic",
				APMResumeSuspend: "windows:APMResumeSuspend",
				APMSuspend: "windows:APMSuspend",
				ApplicationStarted: "windows:ApplicationStarted",
				SystemThemeChanged: "windows:SystemThemeChanged",
				WebViewNavigationCompleted: "windows:WebViewNavigationCompleted",
				WindowActive: "windows:WindowActive",
				WindowBackgroundErase: "windows:WindowBackgroundErase",
				WindowClickActive: "windows:WindowClickActive",
				WindowClosing: "windows:WindowClosing",
				WindowDidMove: "windows:WindowDidMove",
				WindowDidResize: "windows:WindowDidResize",
				WindowDPIChanged: "windows:WindowDPIChanged",
				WindowDragDrop: "windows:WindowDragDrop",
				WindowDragEnter: "windows:WindowDragEnter",
				WindowDragLeave: "windows:WindowDragLeave",
				WindowDragOver: "windows:WindowDragOver",
				WindowEndMove: "windows:WindowEndMove",
				WindowEndResize: "windows:WindowEndResize",
				WindowFullscreen: "windows:WindowFullscreen",
				WindowHide: "windows:WindowHide",
				WindowInactive: "windows:WindowInactive",
				WindowKeyDown: "windows:WindowKeyDown",
				WindowKeyUp: "windows:WindowKeyUp",
				WindowKillFocus: "windows:WindowKillFocus",
				WindowNonClientHit: "windows:WindowNonClientHit",
				WindowNonClientMouseDown: "windows:WindowNonClientMouseDown",
				WindowNonClientMouseLeave: "windows:WindowNonClientMouseLeave",
				WindowNonClientMouseMove: "windows:WindowNonClientMouseMove",
				WindowNonClientMouseUp: "windows:WindowNonClientMouseUp",
				WindowPaint: "windows:WindowPaint",
				WindowRestore: "windows:WindowRestore",
				WindowSetFocus: "windows:WindowSetFocus",
				WindowShow: "windows:WindowShow",
				WindowStartMove: "windows:WindowStartMove",
				WindowStartResize: "windows:WindowStartResize",
				WindowUnFullscreen: "windows:WindowUnFullscreen",
				WindowZOrderChanged: "windows:WindowZOrderChanged",
				WindowMinimise: "windows:WindowMinimise",
				WindowUnMinimise: "windows:WindowUnMinimise",
				WindowMaximise: "windows:WindowMaximise",
				WindowUnMaximise: "windows:WindowUnMaximise"
			}),
			Mac: Object.freeze({
				ApplicationDidBecomeActive: "mac:ApplicationDidBecomeActive",
				ApplicationDidChangeBackingProperties: "mac:ApplicationDidChangeBackingProperties",
				ApplicationDidChangeEffectiveAppearance: "mac:ApplicationDidChangeEffectiveAppearance",
				ApplicationDidChangeIcon: "mac:ApplicationDidChangeIcon",
				ApplicationDidChangeOcclusionState: "mac:ApplicationDidChangeOcclusionState",
				ApplicationDidChangeScreenParameters: "mac:ApplicationDidChangeScreenParameters",
				ApplicationDidChangeStatusBarFrame: "mac:ApplicationDidChangeStatusBarFrame",
				ApplicationDidChangeStatusBarOrientation: "mac:ApplicationDidChangeStatusBarOrientation",
				ApplicationDidChangeTheme: "mac:ApplicationDidChangeTheme",
				ApplicationDidFinishLaunching: "mac:ApplicationDidFinishLaunching",
				ApplicationDidHide: "mac:ApplicationDidHide",
				ApplicationDidResignActive: "mac:ApplicationDidResignActive",
				ApplicationDidUnhide: "mac:ApplicationDidUnhide",
				ApplicationDidUpdate: "mac:ApplicationDidUpdate",
				ApplicationDidWake: "mac:ApplicationDidWake",
				ApplicationScreensDidSleep: "mac:ApplicationScreensDidSleep",
				ApplicationScreensDidWake: "mac:ApplicationScreensDidWake",
				ApplicationShouldHandleReopen: "mac:ApplicationShouldHandleReopen",
				ApplicationWillBecomeActive: "mac:ApplicationWillBecomeActive",
				ApplicationWillFinishLaunching: "mac:ApplicationWillFinishLaunching",
				ApplicationWillHide: "mac:ApplicationWillHide",
				ApplicationWillResignActive: "mac:ApplicationWillResignActive",
				ApplicationWillSleep: "mac:ApplicationWillSleep",
				ApplicationWillTerminate: "mac:ApplicationWillTerminate",
				ApplicationWillUnhide: "mac:ApplicationWillUnhide",
				ApplicationWillUpdate: "mac:ApplicationWillUpdate",
				MenuDidAddItem: "mac:MenuDidAddItem",
				MenuDidBeginTracking: "mac:MenuDidBeginTracking",
				MenuDidClose: "mac:MenuDidClose",
				MenuDidDisplayItem: "mac:MenuDidDisplayItem",
				MenuDidEndTracking: "mac:MenuDidEndTracking",
				MenuDidHighlightItem: "mac:MenuDidHighlightItem",
				MenuDidOpen: "mac:MenuDidOpen",
				MenuDidPopUp: "mac:MenuDidPopUp",
				MenuDidRemoveItem: "mac:MenuDidRemoveItem",
				MenuDidSendAction: "mac:MenuDidSendAction",
				MenuDidSendActionToItem: "mac:MenuDidSendActionToItem",
				MenuDidUpdate: "mac:MenuDidUpdate",
				MenuWillAddItem: "mac:MenuWillAddItem",
				MenuWillBeginTracking: "mac:MenuWillBeginTracking",
				MenuWillDisplayItem: "mac:MenuWillDisplayItem",
				MenuWillEndTracking: "mac:MenuWillEndTracking",
				MenuWillHighlightItem: "mac:MenuWillHighlightItem",
				MenuWillOpen: "mac:MenuWillOpen",
				MenuWillPopUp: "mac:MenuWillPopUp",
				MenuWillRemoveItem: "mac:MenuWillRemoveItem",
				MenuWillSendAction: "mac:MenuWillSendAction",
				MenuWillSendActionToItem: "mac:MenuWillSendActionToItem",
				MenuWillUpdate: "mac:MenuWillUpdate",
				WebViewDidCommitNavigation: "mac:WebViewDidCommitNavigation",
				WebViewDidFinishNavigation: "mac:WebViewDidFinishNavigation",
				WebViewDidReceiveServerRedirectForProvisionalNavigation: "mac:WebViewDidReceiveServerRedirectForProvisionalNavigation",
				WebViewDidStartProvisionalNavigation: "mac:WebViewDidStartProvisionalNavigation",
				WindowDidBecomeKey: "mac:WindowDidBecomeKey",
				WindowDidBecomeMain: "mac:WindowDidBecomeMain",
				WindowDidBeginSheet: "mac:WindowDidBeginSheet",
				WindowDidChangeAlpha: "mac:WindowDidChangeAlpha",
				WindowDidChangeBackingLocation: "mac:WindowDidChangeBackingLocation",
				WindowDidChangeBackingProperties: "mac:WindowDidChangeBackingProperties",
				WindowDidChangeCollectionBehavior: "mac:WindowDidChangeCollectionBehavior",
				WindowDidChangeEffectiveAppearance: "mac:WindowDidChangeEffectiveAppearance",
				WindowDidChangeOcclusionState: "mac:WindowDidChangeOcclusionState",
				WindowDidChangeOrderingMode: "mac:WindowDidChangeOrderingMode",
				WindowDidChangeScreen: "mac:WindowDidChangeScreen",
				WindowDidChangeScreenParameters: "mac:WindowDidChangeScreenParameters",
				WindowDidChangeScreenProfile: "mac:WindowDidChangeScreenProfile",
				WindowDidChangeScreenSpace: "mac:WindowDidChangeScreenSpace",
				WindowDidChangeScreenSpaceProperties: "mac:WindowDidChangeScreenSpaceProperties",
				WindowDidChangeSharingType: "mac:WindowDidChangeSharingType",
				WindowDidChangeSpace: "mac:WindowDidChangeSpace",
				WindowDidChangeSpaceOrderingMode: "mac:WindowDidChangeSpaceOrderingMode",
				WindowDidChangeTitle: "mac:WindowDidChangeTitle",
				WindowDidChangeToolbar: "mac:WindowDidChangeToolbar",
				WindowDidDeminiaturize: "mac:WindowDidDeminiaturize",
				WindowDidEndSheet: "mac:WindowDidEndSheet",
				WindowDidEnterFullScreen: "mac:WindowDidEnterFullScreen",
				WindowDidEnterVersionBrowser: "mac:WindowDidEnterVersionBrowser",
				WindowDidExitFullScreen: "mac:WindowDidExitFullScreen",
				WindowDidExitVersionBrowser: "mac:WindowDidExitVersionBrowser",
				WindowDidExpose: "mac:WindowDidExpose",
				WindowDidFocus: "mac:WindowDidFocus",
				WindowDidMiniaturize: "mac:WindowDidMiniaturize",
				WindowDidMove: "mac:WindowDidMove",
				WindowDidOrderOffScreen: "mac:WindowDidOrderOffScreen",
				WindowDidOrderOnScreen: "mac:WindowDidOrderOnScreen",
				WindowDidResignKey: "mac:WindowDidResignKey",
				WindowDidResignMain: "mac:WindowDidResignMain",
				WindowDidResize: "mac:WindowDidResize",
				WindowDidUpdate: "mac:WindowDidUpdate",
				WindowDidUpdateAlpha: "mac:WindowDidUpdateAlpha",
				WindowDidUpdateCollectionBehavior: "mac:WindowDidUpdateCollectionBehavior",
				WindowDidUpdateCollectionProperties: "mac:WindowDidUpdateCollectionProperties",
				WindowDidUpdateShadow: "mac:WindowDidUpdateShadow",
				WindowDidUpdateTitle: "mac:WindowDidUpdateTitle",
				WindowDidUpdateToolbar: "mac:WindowDidUpdateToolbar",
				WindowDidZoom: "mac:WindowDidZoom",
				WindowFileDraggingEntered: "mac:WindowFileDraggingEntered",
				WindowFileDraggingExited: "mac:WindowFileDraggingExited",
				WindowFileDraggingPerformed: "mac:WindowFileDraggingPerformed",
				WindowHide: "mac:WindowHide",
				WindowMaximise: "mac:WindowMaximise",
				WindowUnMaximise: "mac:WindowUnMaximise",
				WindowMinimise: "mac:WindowMinimise",
				WindowUnMinimise: "mac:WindowUnMinimise",
				WindowShouldClose: "mac:WindowShouldClose",
				WindowShow: "mac:WindowShow",
				WindowWillBecomeKey: "mac:WindowWillBecomeKey",
				WindowWillBecomeMain: "mac:WindowWillBecomeMain",
				WindowWillBeginSheet: "mac:WindowWillBeginSheet",
				WindowWillChangeOrderingMode: "mac:WindowWillChangeOrderingMode",
				WindowWillClose: "mac:WindowWillClose",
				WindowWillDeminiaturize: "mac:WindowWillDeminiaturize",
				WindowWillEnterFullScreen: "mac:WindowWillEnterFullScreen",
				WindowWillEnterVersionBrowser: "mac:WindowWillEnterVersionBrowser",
				WindowWillExitFullScreen: "mac:WindowWillExitFullScreen",
				WindowWillExitVersionBrowser: "mac:WindowWillExitVersionBrowser",
				WindowWillFocus: "mac:WindowWillFocus",
				WindowWillMiniaturize: "mac:WindowWillMiniaturize",
				WindowWillMove: "mac:WindowWillMove",
				WindowWillOrderOffScreen: "mac:WindowWillOrderOffScreen",
				WindowWillOrderOnScreen: "mac:WindowWillOrderOnScreen",
				WindowWillResignMain: "mac:WindowWillResignMain",
				WindowWillResize: "mac:WindowWillResize",
				WindowWillUnfocus: "mac:WindowWillUnfocus",
				WindowWillUpdate: "mac:WindowWillUpdate",
				WindowWillUpdateAlpha: "mac:WindowWillUpdateAlpha",
				WindowWillUpdateCollectionBehavior: "mac:WindowWillUpdateCollectionBehavior",
				WindowWillUpdateCollectionProperties: "mac:WindowWillUpdateCollectionProperties",
				WindowWillUpdateShadow: "mac:WindowWillUpdateShadow",
				WindowWillUpdateTitle: "mac:WindowWillUpdateTitle",
				WindowWillUpdateToolbar: "mac:WindowWillUpdateToolbar",
				WindowWillUpdateVisibility: "mac:WindowWillUpdateVisibility",
				WindowWillUseStandardFrame: "mac:WindowWillUseStandardFrame",
				WindowZoomIn: "mac:WindowZoomIn",
				WindowZoomOut: "mac:WindowZoomOut",
				WindowZoomReset: "mac:WindowZoomReset",
				WebViewWebContentProcessDidTerminate: "mac:WebViewWebContentProcessDidTerminate"
			}),
			Linux: Object.freeze({
				ApplicationStartup: "linux:ApplicationStartup",
				SystemDidWake: "linux:SystemDidWake",
				SystemThemeChanged: "linux:SystemThemeChanged",
				SystemWillSleep: "linux:SystemWillSleep",
				WindowDeleteEvent: "linux:WindowDeleteEvent",
				WindowDidMove: "linux:WindowDidMove",
				WindowDidResize: "linux:WindowDidResize",
				WindowFocusIn: "linux:WindowFocusIn",
				WindowFocusOut: "linux:WindowFocusOut",
				WindowLoadStarted: "linux:WindowLoadStarted",
				WindowLoadRedirected: "linux:WindowLoadRedirected",
				WindowLoadCommitted: "linux:WindowLoadCommitted",
				WindowLoadFinished: "linux:WindowLoadFinished"
			}),
			iOS: Object.freeze({
				ApplicationDidBecomeActive: "ios:ApplicationDidBecomeActive",
				ApplicationDidEnterBackground: "ios:ApplicationDidEnterBackground",
				ApplicationDidFinishLaunching: "ios:ApplicationDidFinishLaunching",
				ApplicationDidReceiveMemoryWarning: "ios:ApplicationDidReceiveMemoryWarning",
				ApplicationWillEnterForeground: "ios:ApplicationWillEnterForeground",
				ApplicationWillResignActive: "ios:ApplicationWillResignActive",
				ApplicationWillTerminate: "ios:ApplicationWillTerminate",
				WindowDidLoad: "ios:WindowDidLoad",
				WindowWillAppear: "ios:WindowWillAppear",
				WindowDidAppear: "ios:WindowDidAppear",
				WindowWillDisappear: "ios:WindowWillDisappear",
				WindowDidDisappear: "ios:WindowDidDisappear",
				WindowSafeAreaInsetsChanged: "ios:WindowSafeAreaInsetsChanged",
				WindowOrientationChanged: "ios:WindowOrientationChanged",
				WindowTouchBegan: "ios:WindowTouchBegan",
				WindowTouchMoved: "ios:WindowTouchMoved",
				WindowTouchEnded: "ios:WindowTouchEnded",
				WindowTouchCancelled: "ios:WindowTouchCancelled",
				WebViewDidStartNavigation: "ios:WebViewDidStartNavigation",
				WebViewDidFinishNavigation: "ios:WebViewDidFinishNavigation",
				WebViewDidFailNavigation: "ios:WebViewDidFailNavigation",
				WebViewDecidePolicyForNavigationAction: "ios:WebViewDecidePolicyForNavigationAction",
				BatteryChanged: "ios:BatteryChanged",
				NetworkChanged: "ios:NetworkChanged",
				ThemeChanged: "ios:ThemeChanged",
				ScreenLocked: "ios:ScreenLocked",
				ScreenUnlocked: "ios:ScreenUnlocked"
			}),
			Android: Object.freeze({
				ActivityCreated: "android:ActivityCreated",
				ActivityStarted: "android:ActivityStarted",
				ActivityResumed: "android:ActivityResumed",
				ActivityPaused: "android:ActivityPaused",
				ActivityStopped: "android:ActivityStopped",
				ActivityDestroyed: "android:ActivityDestroyed",
				ApplicationLowMemory: "android:ApplicationLowMemory",
				WebViewPageStarted: "android:WebViewPageStarted",
				WebViewPageFinished: "android:WebViewPageFinished",
				BatteryChanged: "android:BatteryChanged",
				NetworkChanged: "android:NetworkChanged",
				ThemeChanged: "android:ThemeChanged",
				ScreenLocked: "android:ScreenLocked",
				ScreenUnlocked: "android:ScreenUnlocked"
			}),
			Openharmony: Object.freeze({
				AbilityCreated: "openharmony:AbilityCreated",
				AbilityStarted: "openharmony:AbilityStarted",
				AbilityForeground: "openharmony:AbilityForeground",
				AbilityBackground: "openharmony:AbilityBackground",
				AbilityDestroyed: "openharmony:AbilityDestroyed",
				ApplicationLowMemory: "openharmony:ApplicationLowMemory",
				WebViewPageStarted: "openharmony:WebViewPageStarted",
				WebViewPageFinished: "openharmony:WebViewPageFinished",
				BatteryChanged: "openharmony:BatteryChanged",
				NetworkChanged: "openharmony:NetworkChanged",
				ThemeChanged: "openharmony:ThemeChanged",
				ScreenLocked: "openharmony:ScreenLocked",
				ScreenUnlocked: "openharmony:ScreenUnlocked"
			}),
			Common: Object.freeze({
				ApplicationOpenedWithFile: "common:ApplicationOpenedWithFile",
				ApplicationStarted: "common:ApplicationStarted",
				ApplicationLaunchedWithUrl: "common:ApplicationLaunchedWithUrl",
				ThemeChanged: "common:ThemeChanged",
				SystemDidWake: "common:SystemDidWake",
				SystemWillSleep: "common:SystemWillSleep",
				WindowClosing: "common:WindowClosing",
				WindowDidMove: "common:WindowDidMove",
				WindowDidResize: "common:WindowDidResize",
				WindowDPIChanged: "common:WindowDPIChanged",
				WindowFilesDropped: "common:WindowFilesDropped",
				WindowFocus: "common:WindowFocus",
				WindowFullscreen: "common:WindowFullscreen",
				WindowHide: "common:WindowHide",
				WindowLostFocus: "common:WindowLostFocus",
				WindowMaximise: "common:WindowMaximise",
				WindowMinimise: "common:WindowMinimise",
				WindowToggleFrameless: "common:WindowToggleFrameless",
				WindowRestore: "common:WindowRestore",
				WindowRuntimeReady: "common:WindowRuntimeReady",
				WindowShow: "common:WindowShow",
				WindowUnFullscreen: "common:WindowUnFullscreen",
				WindowUnMaximise: "common:WindowUnMaximise",
				WindowUnMinimise: "common:WindowUnMinimise",
				WindowZoom: "common:WindowZoom",
				WindowZoomIn: "common:WindowZoomIn",
				WindowZoomOut: "common:WindowZoomOut",
				WindowZoomReset: "common:WindowZoomReset",
				BatteryChanged: "common:BatteryChanged",
				NetworkChanged: "common:NetworkChanged",
				ScreenLocked: "common:ScreenLocked",
				ScreenUnlocked: "common:ScreenUnlocked",
				LowMemory: "common:LowMemory"
			})
		});
	})))();
}

// src/events.ts
var events_exports = /* @__PURE__ */ __exportAll({
	Emit: () => Emit,
	Off: () => Off,
	OffAll: () => OffAll,
	On: () => On,
	OnMultiple: () => OnMultiple,
	Once: () => Once,
	Types: () => Types,
	WailsEvent: () => WailsEvent
});
function dispatchWailsEvent(event) {
	let listeners = eventListeners.get(event.name);
	if (!listeners) return;
	let wailsEvent = new WailsEvent(event.name, event.name in Events$1 ? Events$1[event.name](event.data) : event.data);
	if ("sender" in event) wailsEvent.sender = event.sender;
	const expired = /* @__PURE__ */ new Set();
	for (const listener of listeners.slice()) if (listener.dispatch(wailsEvent)) expired.add(listener);
	if (expired.size > 0) {
		const live = eventListeners.get(event.name);
		if (live) {
			const remaining = live.filter((l) => !expired.has(l));
			if (remaining.length === 0) eventListeners.delete(event.name);
			else eventListeners.set(event.name, remaining);
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
function Once(eventName, callback) {
	return OnMultiple(eventName, callback, 1);
}
function Off(...eventNames) {
	eventNames.forEach((eventName) => eventListeners.delete(eventName));
}
function OffAll() {
	eventListeners.clear();
}
function Emit(name, data) {
	return call$3(EmitMethod, new WailsEvent(name, data));
}
var call$3, EmitMethod, WailsEvent;
function init_events() {
	return (init_events = __esmMin((() => {
		init_runtime();
		init_listener();
		init_create();
		init_event_types();
		init_environment();
		if (hasDOM) {
			window._wails = window._wails || {};
			window._wails.dispatchWailsEvent = dispatchWailsEvent;
		}
		call$3 = newRuntimeCaller(objectNames.Events);
		EmitMethod = 0;
		WailsEvent = class {
			constructor(name, data) {
				this.name = name;
				this.data = data !== null && data !== void 0 ? data : null;
			}
		};
	})))();
}

// src/screens.ts
var screens_exports = /* @__PURE__ */ __exportAll({
	GetAll: () => GetAll,
	GetByID: () => GetByID,
	GetByIndex: () => GetByIndex,
	GetCurrent: () => GetCurrent,
	GetPrimary: () => GetPrimary
});
function GetAll() {
	return call$2(getAll);
}
function GetPrimary() {
	return call$2(getPrimary);
}
function GetCurrent() {
	return call$2(getCurrent);
}
function GetByID(id) {
	return call$2(getByID, { id });
}
function GetByIndex(index) {
	return call$2(getByIndex, { index });
}
var call$2, getAll, getPrimary, getCurrent, getByID, getByIndex;
function init_screens() {
	return (init_screens = __esmMin((() => {
		init_runtime();
		call$2 = newRuntimeCaller(objectNames.Screens);
		getAll = 0;
		getPrimary = 1;
		getCurrent = 2;
		getByID = 3;
		getByIndex = 4;
	})))();
}

// src/ios.ts
var ios_exports = /* @__PURE__ */ __exportAll({
	Device: () => Device$1,
	Haptics: () => Haptics$1
});
var call$1, HapticsImpact, DeviceInfo$1, Haptics$1, Device$1;
function init_ios() {
	return (init_ios = __esmMin((() => {
		init_runtime();
		call$1 = newRuntimeCaller(objectNames.IOS);
		HapticsImpact = 0;
		DeviceInfo$1 = 1;
		;
		(function(_Haptics) {
			function Impact(style = "medium") {
				return call$1(HapticsImpact, { style });
			}
			_Haptics.Impact = Impact;
		})(Haptics$1 || (Haptics$1 = {}));
		;
		(function(_Device) {
			function Info() {
				return call$1(DeviceInfo$1);
			}
			_Device.Info = Info;
		})(Device$1 || (Device$1 = {}));
	})))();
}

// src/android.ts
var android_exports = /* @__PURE__ */ __exportAll({
	Device: () => Device,
	Haptics: () => Haptics,
	Toast: () => Toast
});
var call, HapticsVibrate, DeviceInfo, ToastShow, Haptics, Device, Toast;
function init_android() {
	return (init_android = __esmMin((() => {
		init_runtime();
		call = newRuntimeCaller(objectNames.Android);
		HapticsVibrate = 0;
		DeviceInfo = 1;
		ToastShow = 2;
		;
		(function(_Haptics) {
			function Vibrate(durationMs = 100) {
				return call(HapticsVibrate, { duration: durationMs });
			}
			_Haptics.Vibrate = Vibrate;
		})(Haptics || (Haptics = {}));
		;
		(function(_Device) {
			function Info() {
				return call(DeviceInfo);
			}
			_Device.Info = Info;
		})(Device || (Device = {}));
		;
		(function(_Toast) {
			function Show(message) {
				return call(ToastShow, { message });
			}
			_Toast.Show = Show;
		})(Toast || (Toast = {}));
	})))();
}

// src/updater.ts
var updater_exports = /* @__PURE__ */ __exportAll({ Events: () => Events });
var Events;
function init_updater() {
	return (init_updater = __esmMin((() => {
		Events = Object.freeze({
			CheckStarted: "wails:updater:check-started",
			UpdateAvailable: "wails:updater:update-available",
			NoUpdate: "wails:updater:no-update",
			DownloadStarted: "wails:updater:download-started",
			DownloadProgress: "wails:updater:download-progress",
			DownloadComplete: "wails:updater:download-complete",
			Verifying: "wails:updater:verifying",
			Installing: "wails:updater:installing",
			UpdateReady: "wails:updater:update-ready",
			Error: "wails:updater:error",
			Meta: "wails:updater:meta",
			User: Object.freeze({
				Install: "wails:updater:user:install",
				Restart: "wails:updater:user:restart",
				Skip: "wails:updater:user:skip",
				Remind: "wails:updater:user:remind",
				Cancel: "wails:updater:user:cancel"
			}),
			Window: Object.freeze({ Ready: "wails:updater:window:ready" })
		});
	})))();
}

// src/window.ts
function getDropTargetElement(element) {
	if (!element) return null;
	return element.closest(`[${DROP_TARGET_ATTRIBUTE}]`);
}
function canResolveFilePaths() {
	var _window$chrome, _window$_wails;
	if (((_window$chrome = window.chrome) === null || _window$chrome === void 0 || (_window$chrome = _window$chrome.webview) === null || _window$chrome === void 0 ? void 0 : _window$chrome.postMessageWithAdditionalObjects) == null) return false;
	return ((_window$_wails = window._wails) === null || _window$_wails === void 0 || (_window$_wails = _window$_wails.flags) === null || _window$_wails === void 0 ? void 0 : _window$_wails.enableFileDrop) === true;
}
function resolveFilePaths(x, y, files) {
	var _window$chrome2;
	if ((_window$chrome2 = window.chrome) === null || _window$chrome2 === void 0 || (_window$chrome2 = _window$chrome2.webview) === null || _window$chrome2 === void 0 ? void 0 : _window$chrome2.postMessageWithAdditionalObjects) window.chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
}
function cleanupNativeDrag() {
	nativeDragActive = false;
	if (currentDropTarget) {
		currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
		currentDropTarget = null;
	}
}
function handleDragEnter() {
	var _window$_wails2;
	if (((_window$_wails2 = window._wails) === null || _window$_wails2 === void 0 || (_window$_wails2 = _window$_wails2.flags) === null || _window$_wails2 === void 0 ? void 0 : _window$_wails2.enableFileDrop) === false) return;
	nativeDragActive = true;
}
function handleDragLeave() {
	cleanupNativeDrag();
}
function handleDragOver(x, y) {
	var _window$_wails3;
	if (!nativeDragActive) return;
	if (((_window$_wails3 = window._wails) === null || _window$_wails3 === void 0 || (_window$_wails3 = _window$_wails3.flags) === null || _window$_wails3 === void 0 ? void 0 : _window$_wails3.enableFileDrop) === false) return;
	const dropTarget = getDropTargetElement(document.elementFromPoint(x, y));
	if (currentDropTarget && currentDropTarget !== dropTarget) currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
	if (dropTarget) {
		dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
		currentDropTarget = dropTarget;
	} else currentDropTarget = null;
}
function setupDropTargetListeners() {
	const docElement = document.documentElement;
	let dragEnterCounter = 0;
	docElement.addEventListener("dragenter", (event) => {
		var _event$dataTransfer, _window$_wails5;
		if (!((_event$dataTransfer = event.dataTransfer) === null || _event$dataTransfer === void 0 ? void 0 : _event$dataTransfer.types.includes("Files"))) return;
		event.preventDefault();
		if (((_window$_wails5 = window._wails) === null || _window$_wails5 === void 0 || (_window$_wails5 = _window$_wails5.flags) === null || _window$_wails5 === void 0 ? void 0 : _window$_wails5.enableFileDrop) === false) {
			event.dataTransfer.dropEffect = "none";
			return;
		}
		dragEnterCounter++;
		const dropTarget = getDropTargetElement(document.elementFromPoint(event.clientX, event.clientY));
		if (currentDropTarget && currentDropTarget !== dropTarget) currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
		if (dropTarget) {
			dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
			event.dataTransfer.dropEffect = "copy";
			currentDropTarget = dropTarget;
		} else {
			event.dataTransfer.dropEffect = "none";
			currentDropTarget = null;
		}
	}, false);
	docElement.addEventListener("dragover", (event) => {
		var _event$dataTransfer2, _window$_wails6;
		if (!((_event$dataTransfer2 = event.dataTransfer) === null || _event$dataTransfer2 === void 0 ? void 0 : _event$dataTransfer2.types.includes("Files"))) return;
		event.preventDefault();
		if (((_window$_wails6 = window._wails) === null || _window$_wails6 === void 0 || (_window$_wails6 = _window$_wails6.flags) === null || _window$_wails6 === void 0 ? void 0 : _window$_wails6.enableFileDrop) === false) {
			event.dataTransfer.dropEffect = "none";
			return;
		}
		const dropTarget = getDropTargetElement(document.elementFromPoint(event.clientX, event.clientY));
		if (currentDropTarget && currentDropTarget !== dropTarget) currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
		if (dropTarget) {
			if (!dropTarget.classList.contains(DROP_TARGET_ACTIVE_CLASS)) dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
			event.dataTransfer.dropEffect = "copy";
			currentDropTarget = dropTarget;
		} else {
			event.dataTransfer.dropEffect = "none";
			currentDropTarget = null;
		}
	}, false);
	docElement.addEventListener("dragleave", (event) => {
		var _event$dataTransfer3, _window$_wails7;
		if (!((_event$dataTransfer3 = event.dataTransfer) === null || _event$dataTransfer3 === void 0 ? void 0 : _event$dataTransfer3.types.includes("Files"))) return;
		event.preventDefault();
		if (((_window$_wails7 = window._wails) === null || _window$_wails7 === void 0 || (_window$_wails7 = _window$_wails7.flags) === null || _window$_wails7 === void 0 ? void 0 : _window$_wails7.enableFileDrop) === false) return;
		if (event.relatedTarget === null) return;
		dragEnterCounter--;
		if (dragEnterCounter === 0 || currentDropTarget && !currentDropTarget.contains(event.relatedTarget)) {
			if (currentDropTarget) {
				currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
				currentDropTarget = null;
			}
			dragEnterCounter = 0;
		}
	}, false);
	docElement.addEventListener("drop", (event) => {
		var _event$dataTransfer4, _window$_wails8;
		if (!((_event$dataTransfer4 = event.dataTransfer) === null || _event$dataTransfer4 === void 0 ? void 0 : _event$dataTransfer4.types.includes("Files"))) return;
		event.preventDefault();
		if (((_window$_wails8 = window._wails) === null || _window$_wails8 === void 0 || (_window$_wails8 = _window$_wails8.flags) === null || _window$_wails8 === void 0 ? void 0 : _window$_wails8.enableFileDrop) === false) return;
		dragEnterCounter = 0;
		if (currentDropTarget) {
			currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
			currentDropTarget = null;
		}
		if (canResolveFilePaths()) {
			const files = [];
			if (event.dataTransfer.items) {
				for (const item of event.dataTransfer.items) if (item.kind === "file") {
					const file = item.getAsFile();
					if (file) files.push(file);
				}
			} else if (event.dataTransfer.files) for (const file of event.dataTransfer.files) files.push(file);
			if (files.length > 0) resolveFilePaths(event.clientX, event.clientY, files);
		}
	}, false);
}
var DROP_TARGET_ATTRIBUTE, DROP_TARGET_ACTIVE_CLASS, currentDropTarget, PositionMethod, CenterMethod, CloseMethod, DisableSizeConstraintsMethod, EnableSizeConstraintsMethod, FocusMethod, ForceReloadMethod, FullscreenMethod, GetScreenMethod, GetZoomMethod, HeightMethod, HideMethod, IsFocusedMethod, IsFullscreenMethod, IsMaximisedMethod, IsMinimisedMethod, MaximiseMethod, MinimiseMethod, NameMethod, OpenDevToolsMethod, RelativePositionMethod, ReloadMethod, ResizableMethod, RestoreMethod, SetPositionMethod, SetAlwaysOnTopMethod, SetBackgroundColourMethod, SetFramelessMethod, SetFullscreenButtonEnabledMethod, SetMaxSizeMethod, SetMinSizeMethod, SetRelativePositionMethod, SetResizableMethod, SetSizeMethod, SetTitleMethod, SetZoomMethod, ShowMethod, SizeMethod, ToggleFullscreenMethod, ToggleMaximiseMethod, ToggleFramelessMethod, UnFullscreenMethod, UnMaximiseMethod, UnMinimiseMethod, WidthMethod, ZoomMethod, ZoomInMethod, ZoomOutMethod, ZoomResetMethod, SnapAssistMethod, FilesDropped, PrintMethod, SetScreenMethod, nativeDragActive, callerSym, Window, thisWindow;
function init_window() {
	return (init_window = __esmMin((() => {
		init_runtime();
		DROP_TARGET_ATTRIBUTE = "data-file-drop-target";
		DROP_TARGET_ACTIVE_CLASS = "file-drop-target-active";
		currentDropTarget = null;
		PositionMethod = 0;
		CenterMethod = 1;
		CloseMethod = 2;
		DisableSizeConstraintsMethod = 3;
		EnableSizeConstraintsMethod = 4;
		FocusMethod = 5;
		ForceReloadMethod = 6;
		FullscreenMethod = 7;
		GetScreenMethod = 8;
		GetZoomMethod = 9;
		HeightMethod = 10;
		HideMethod = 11;
		IsFocusedMethod = 12;
		IsFullscreenMethod = 13;
		IsMaximisedMethod = 14;
		IsMinimisedMethod = 15;
		MaximiseMethod = 16;
		MinimiseMethod = 17;
		NameMethod = 18;
		OpenDevToolsMethod = 19;
		RelativePositionMethod = 20;
		ReloadMethod = 21;
		ResizableMethod = 22;
		RestoreMethod = 23;
		SetPositionMethod = 24;
		SetAlwaysOnTopMethod = 25;
		SetBackgroundColourMethod = 26;
		SetFramelessMethod = 27;
		SetFullscreenButtonEnabledMethod = 28;
		SetMaxSizeMethod = 29;
		SetMinSizeMethod = 30;
		SetRelativePositionMethod = 31;
		SetResizableMethod = 32;
		SetSizeMethod = 33;
		SetTitleMethod = 34;
		SetZoomMethod = 35;
		ShowMethod = 36;
		SizeMethod = 37;
		ToggleFullscreenMethod = 38;
		ToggleMaximiseMethod = 39;
		ToggleFramelessMethod = 40;
		UnFullscreenMethod = 41;
		UnMaximiseMethod = 42;
		UnMinimiseMethod = 43;
		WidthMethod = 44;
		ZoomMethod = 45;
		ZoomInMethod = 46;
		ZoomOutMethod = 47;
		ZoomResetMethod = 48;
		SnapAssistMethod = 49;
		FilesDropped = 50;
		PrintMethod = 51;
		SetScreenMethod = 52;
		nativeDragActive = false;
		callerSym = Symbol("caller");
		Window = class Window {
			constructor(name = "") {
				this[callerSym] = newRuntimeCaller(objectNames.Window, name);
				for (const method of Object.getOwnPropertyNames(Window.prototype)) if (method !== "constructor" && typeof this[method] === "function") this[method] = this[method].bind(this);
			}
			Get(name) {
				return new Window(name);
			}
			Position() {
				return this[callerSym](PositionMethod);
			}
			Center() {
				return this[callerSym](CenterMethod);
			}
			Close() {
				return this[callerSym](CloseMethod);
			}
			DisableSizeConstraints() {
				return this[callerSym](DisableSizeConstraintsMethod);
			}
			EnableSizeConstraints() {
				return this[callerSym](EnableSizeConstraintsMethod);
			}
			Focus() {
				return this[callerSym](FocusMethod);
			}
			ForceReload() {
				return this[callerSym](ForceReloadMethod);
			}
			Fullscreen() {
				return this[callerSym](FullscreenMethod);
			}
			GetScreen() {
				return this[callerSym](GetScreenMethod);
			}
			GetZoom() {
				return this[callerSym](GetZoomMethod);
			}
			Height() {
				return this[callerSym](HeightMethod);
			}
			Hide() {
				return this[callerSym](HideMethod);
			}
			IsFocused() {
				return this[callerSym](IsFocusedMethod);
			}
			IsFullscreen() {
				return this[callerSym](IsFullscreenMethod);
			}
			IsMaximised() {
				return this[callerSym](IsMaximisedMethod);
			}
			IsMinimised() {
				return this[callerSym](IsMinimisedMethod);
			}
			Maximise() {
				return this[callerSym](MaximiseMethod);
			}
			Minimise() {
				return this[callerSym](MinimiseMethod);
			}
			Name() {
				return this[callerSym](NameMethod);
			}
			OpenDevTools() {
				return this[callerSym](OpenDevToolsMethod);
			}
			RelativePosition() {
				return this[callerSym](RelativePositionMethod);
			}
			Reload() {
				return this[callerSym](ReloadMethod);
			}
			Resizable() {
				return this[callerSym](ResizableMethod);
			}
			Restore() {
				return this[callerSym](RestoreMethod);
			}
			SetPosition(x, y) {
				return this[callerSym](SetPositionMethod, {
					x,
					y
				});
			}
			SetAlwaysOnTop(alwaysOnTop) {
				return this[callerSym](SetAlwaysOnTopMethod, { alwaysOnTop });
			}
			SetBackgroundColour(r, g, b, a) {
				return this[callerSym](SetBackgroundColourMethod, {
					r,
					g,
					b,
					a
				});
			}
			SetFrameless(frameless) {
				return this[callerSym](SetFramelessMethod, { frameless });
			}
			SetFullscreenButtonEnabled(enabled) {
				return this[callerSym](SetFullscreenButtonEnabledMethod, { enabled });
			}
			SetMaxSize(width, height) {
				return this[callerSym](SetMaxSizeMethod, {
					width,
					height
				});
			}
			SetMinSize(width, height) {
				return this[callerSym](SetMinSizeMethod, {
					width,
					height
				});
			}
			SetRelativePosition(x, y) {
				return this[callerSym](SetRelativePositionMethod, {
					x,
					y
				});
			}
			SetResizable(resizable) {
				return this[callerSym](SetResizableMethod, { resizable });
			}
			SetSize(width, height) {
				return this[callerSym](SetSizeMethod, {
					width,
					height
				});
			}
			SetTitle(title) {
				return this[callerSym](SetTitleMethod, { title });
			}
			SetZoom(zoom) {
				return this[callerSym](SetZoomMethod, { zoom });
			}
			Show() {
				return this[callerSym](ShowMethod);
			}
			Size() {
				return this[callerSym](SizeMethod);
			}
			ToggleFullscreen() {
				return this[callerSym](ToggleFullscreenMethod);
			}
			ToggleMaximise() {
				return this[callerSym](ToggleMaximiseMethod);
			}
			ToggleFrameless() {
				return this[callerSym](ToggleFramelessMethod);
			}
			UnFullscreen() {
				return this[callerSym](UnFullscreenMethod);
			}
			UnMaximise() {
				return this[callerSym](UnMaximiseMethod);
			}
			UnMinimise() {
				return this[callerSym](UnMinimiseMethod);
			}
			Width() {
				return this[callerSym](WidthMethod);
			}
			Zoom() {
				return this[callerSym](ZoomMethod);
			}
			ZoomIn() {
				return this[callerSym](ZoomInMethod);
			}
			ZoomOut() {
				return this[callerSym](ZoomOutMethod);
			}
			ZoomReset() {
				return this[callerSym](ZoomResetMethod);
			}
			HandlePlatformFileDrop(filenames, x, y) {
				var _window$_wails4;
				if (((_window$_wails4 = window._wails) === null || _window$_wails4 === void 0 || (_window$_wails4 = _window$_wails4.flags) === null || _window$_wails4 === void 0 ? void 0 : _window$_wails4.enableFileDrop) === false) return;
				const dropTarget = getDropTargetElement(document.elementFromPoint(x, y));
				if (!dropTarget) return;
				const elementDetails = {
					id: dropTarget.id,
					classList: Array.from(dropTarget.classList),
					attributes: {}
				};
				for (let i = 0; i < dropTarget.attributes.length; i++) {
					const attr = dropTarget.attributes[i];
					elementDetails.attributes[attr.name] = attr.value;
				}
				const payload = {
					filenames,
					x,
					y,
					elementDetails
				};
				this[callerSym](FilesDropped, payload);
				cleanupNativeDrag();
			}
			SetScreen(screenID) {
				return this[callerSym](SetScreenMethod, { screenID });
			}
			SnapAssist() {
				return this[callerSym](SnapAssistMethod);
			}
			Print() {
				return this[callerSym](PrintMethod);
			}
		};
		thisWindow = new Window("");
		if (typeof window !== "undefined" && typeof document !== "undefined") setupDropTargetListeners();
	})))();
}

// src/wml.ts
var wml_exports = /* @__PURE__ */ __exportAll({
	Enable: () => Enable,
	Reload: () => Reload
});
function sendEvent(eventName, data = null) {
	Emit(eventName, data);
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
		if (choice !== "Yes") return;
		const eventType = element.getAttribute("wml-event") || element.getAttribute("data-wml-event");
		const targetWindow = element.getAttribute("wml-target-window") || element.getAttribute("data-wml-target-window") || "";
		const windowMethod = element.getAttribute("wml-window") || element.getAttribute("data-wml-window");
		const url = element.getAttribute("wml-openurl") || element.getAttribute("data-wml-openurl");
		if (eventType !== null) sendEvent(eventType);
		if (windowMethod !== null) callWindowMethod(targetWindow, windowMethod);
		if (url !== null) OpenURL(url);
	}
	const confirm = element.getAttribute("wml-confirm") || element.getAttribute("data-wml-confirm");
	if (confirm) Question({
		Title: "Confirm",
		Message: confirm,
		Detached: false,
		Buttons: [{ Label: "Yes" }, {
			Label: "No",
			IsDefault: true
		}]
	}).then(runEffect);
	else runEffect();
}
function addWMLListeners(element) {
	const triggerRegExp = /\S+/g;
	const triggerAttr = element.getAttribute("wml-trigger") || element.getAttribute("data-wml-trigger") || "click";
	const triggers = [];
	let match;
	while ((match = triggerRegExp.exec(triggerAttr)) !== null) triggers.push(match[0]);
	const options = triggerRegistry.set(element, triggers);
	for (const trigger of triggers) element.addEventListener(trigger, onWMLTriggered, options);
}
function Enable() {
	whenReady(Reload);
}
function Reload() {
	triggerRegistry.reset();
	document.body.querySelectorAll("[wml-event], [wml-window], [wml-openurl], [data-wml-event], [data-wml-window], [data-wml-openurl]").forEach(addWMLListeners);
}
var controllerSym, triggerMapSym, elementCountSym, AbortControllerRegistry, WeakMapRegistry, triggerRegistry;
function init_wml() {
	return (init_wml = __esmMin((() => {
		init_browser();
		init_dialogs();
		init_events();
		init_utils();
		init_window();
		controllerSym = Symbol("controller");
		triggerMapSym = Symbol("triggerMap");
		elementCountSym = Symbol("elementCount");
		AbortControllerRegistry = class {
			constructor() {
				this[controllerSym] = new AbortController();
			}
			set(element, triggers) {
				return { signal: this[controllerSym].signal };
			}
			reset() {
				this[controllerSym].abort();
				this[controllerSym] = new AbortController();
			}
		};
		WeakMapRegistry = class {
			constructor() {
				this[triggerMapSym] = /* @__PURE__ */ new WeakMap();
				this[elementCountSym] = 0;
			}
			set(element, triggers) {
				if (!this[triggerMapSym].has(element)) this[elementCountSym]++;
				this[triggerMapSym].set(element, triggers);
				return {};
			}
			reset() {
				if (this[elementCountSym] <= 0) return;
				for (const element of document.body.querySelectorAll("*")) {
					if (this[elementCountSym] <= 0) break;
					const triggers = this[triggerMapSym].get(element);
					if (triggers != null) this[elementCountSym]--;
					for (const trigger of triggers || []) element.removeEventListener(trigger, onWMLTriggered);
				}
				this[triggerMapSym] = /* @__PURE__ */ new WeakMap();
				this[elementCountSym] = 0;
			}
		};
		triggerRegistry = canAbortListeners() ? new AbortControllerRegistry() : new WeakMapRegistry();
	})))();
}

// src/stream.ts
function namedPageGeneration() {
	try {
		const marker = window.name.lastIndexOf(GENERATION_NAME_MARKER);
		if (marker < 0) return 0;
		const value = Number(window.name.slice(marker + 29));
		return Number.isSafeInteger(value) && value > 0 ? value : 0;
	} catch (_unused) {
		return 0;
	}
}
function rememberNamedPageGeneration(generation) {
	try {
		const marker = window.name.lastIndexOf(GENERATION_NAME_MARKER);
		const applicationName = marker < 0 ? window.name : window.name.slice(0, marker);
		window.name = applicationName + GENERATION_NAME_MARKER + generation;
	} catch (_unused2) {}
}
function nextPageGeneration() {
	if (!hasDOM) return 1;
	const origin = typeof performance !== "undefined" && Number.isFinite(performance.timeOrigin) ? performance.timeOrigin : Date.now();
	const pageTime = Math.max(1, Math.floor(origin * 1e3));
	let current = namedPageGeneration();
	try {
		const raw = window.sessionStorage.getItem(GENERATION_KEY);
		const stored = raw === null ? 0 : Number(raw);
		if (Number.isSafeInteger(stored) && stored > current) current = stored;
	} catch (_unused3) {}
	const next = Number.isSafeInteger(current) && current >= 0 && current < Number.MAX_SAFE_INTEGER ? Math.max(current + 1, pageTime) : pageTime;
	try {
		window.sessionStorage.setItem(GENERATION_KEY, String(next));
	} catch (_unused4) {}
	rememberNamedPageGeneration(next);
	return next;
}
function streamGlobals() {
	const fresh = () => ({
		session: nanoid(),
		generation: nextPageGeneration(),
		nextConnID: 1,
		connections: /* @__PURE__ */ new Map(),
		polling: false
	});
	if (!hasDOM) return fresh();
	const w = window._wails = window._wails || {};
	return w.__stream = w.__stream || fresh();
}
function streamURL(endpoint) {
	return window.location.origin + "/wails/stream/" + endpoint;
}
function sleep(ms, signal) {
	if (!signal) return new Promise((resolve) => setTimeout(resolve, ms));
	if (signal.aborted) return Promise.reject(new StreamRequestCancelled());
	return new Promise((resolve, reject) => {
		const onAbort = () => {
			clearTimeout(timer);
			reject(new StreamRequestCancelled());
		};
		const timer = setTimeout(() => {
			signal.removeEventListener("abort", onAbort);
			resolve();
		}, ms);
		signal.addEventListener("abort", onAbort, { once: true });
	});
}
function Stream(name) {
	if (hasDOM) {
		var _window$_wails;
		const factory = (_window$_wails = window._wails) === null || _window$_wails === void 0 ? void 0 : _window$_wails.streamFactory;
		if (typeof factory === "function") return factory(name);
	}
	return new WailsSocket(name);
}
function defineHandlerProperty(target, type) {
	let current = null;
	Object.defineProperty(target, "on" + type, {
		get: () => current,
		set(fn) {
			if (current) target.removeEventListener(type, current);
			current = typeof fn === "function" ? fn : null;
			if (current) target.addEventListener(type, current);
		},
		configurable: true,
		enumerable: true
	});
}
function JSONStream(name) {
	const socket = Stream(name);
	const decoder = new TextDecoder();
	const parse = (payload) => JSON.parse(typeof payload === "string" ? payload : decoder.decode(payload));
	if (socket instanceof WailsSocket) socket._decode = (payload) => parse(payload);
	else {
		const native = socket;
		native.binaryType = "arraybuffer";
		const wrapped = /* @__PURE__ */ new WeakMap();
		const decodedEvents = /* @__PURE__ */ new WeakMap();
		const add = native.addEventListener.bind(native);
		const remove = native.removeEventListener.bind(native);
		const decodeEvent = (ev) => {
			var _decodedEvents$get;
			if (decodedEvents.has(ev)) return (_decodedEvents$get = decodedEvents.get(ev)) !== null && _decodedEvents$get !== void 0 ? _decodedEvents$get : null;
			try {
				const value = parse(ev.data);
				const decoded = new MessageEvent("message", { data: value });
				decodedEvents.set(ev, decoded);
				return decoded;
			} catch (_unused6) {
				decodedEvents.set(ev, null);
				native.dispatchEvent(new Event("error"));
				return null;
			}
		};
		const decode = (listener) => {
			const existing = wrapped.get(listener);
			if (existing) return existing;
			const fn = (ev) => {
				const decoded = decodeEvent(ev);
				if (!decoded) return;
				if (typeof listener === "function") listener.call(native, decoded);
				else listener.handleEvent(decoded);
			};
			wrapped.set(listener, fn);
			return fn;
		};
		native.addEventListener = ((type, listener, opts) => {
			add(type, type === "message" && listener ? decode(listener) : listener, opts);
		});
		native.removeEventListener = ((type, listener, opts) => {
			const fn = type === "message" && listener ? wrapped.get(listener) : void 0;
			remove(type, fn !== null && fn !== void 0 ? fn : listener, opts);
		});
		let handler = null;
		Object.defineProperty(native, "onmessage", {
			get: () => handler,
			set(fn) {
				if (handler) native.removeEventListener("message", handler);
				handler = typeof fn === "function" ? fn : null;
				if (handler) native.addEventListener("message", handler);
			},
			configurable: true,
			enumerable: true
		});
	}
	const send = socket.send.bind(socket);
	socket.send = (value) => send(JSON.stringify(value));
	return socket;
}
function toBytesSync(data) {
	if (typeof data === "string") return textEncoder.encode(data);
	if (typeof Blob !== "undefined" && data instanceof Blob) return null;
	if (ArrayBuffer.isView(data)) return new Uint8Array(data.buffer, data.byteOffset, data.byteLength).slice();
	return new Uint8Array(data).slice();
}
function abortable(promise, signal) {
	if (!signal) return promise;
	if (signal.aborted) return Promise.reject(new StreamRequestCancelled());
	return new Promise((resolve, reject) => {
		const onAbort = () => reject(new StreamRequestCancelled());
		signal.addEventListener("abort", onAbort, { once: true });
		promise.then((value) => {
			signal.removeEventListener("abort", onAbort);
			resolve(value);
		}, (err) => {
			signal.removeEventListener("abort", onAbort);
			reject(err);
		});
	});
}
async function toBytes(data, signal) {
	if (typeof data === "string") return textEncoder.encode(data);
	if (typeof Blob !== "undefined" && data instanceof Blob) return new Uint8Array(await abortable(data.arrayBuffer(), signal));
	if (ArrayBuffer.isView(data)) return new Uint8Array(data.buffer, data.byteOffset, data.byteLength).slice();
	return new Uint8Array(data).slice();
}
async function postFrame(connID, kind, body, name, signal) {
	const headers = {
		[HDR_SESSION]: G.session,
		[HDR_GENERATION]: String(G.generation),
		[HDR_CONN]: String(connID),
		[HDR_KIND]: String(kind),
		"Content-Type": "application/octet-stream"
	};
	if (name !== void 0) headers[HDR_NAME] = name;
	if (body.byteLength <= CHUNK_THRESHOLD) {
		await postWithRetry(headers, body, signal);
		return;
	}
	const chunkID = nanoid();
	const total = Math.ceil(body.byteLength / CHUNK_THRESHOLD);
	for (let i = 0; i < total; i++) {
		const slice = body.subarray(i * CHUNK_THRESHOLD, (i + 1) * CHUNK_THRESHOLD);
		await postWithRetry(_objectSpread2(_objectSpread2({}, headers), {}, {
			[HDR_CHUNK]: chunkID,
			[HDR_CHUNK_INDEX]: String(i),
			[HDR_CHUNK_TOTAL]: String(total)
		}), slice, signal);
	}
}
async function postWithRetry(headers, body, signal) {
	let wait = 1;
	for (;;) {
		if (signal === null || signal === void 0 ? void 0 : signal.aborted) throw new StreamRequestCancelled();
		const resp = await fetch(streamURL("send"), {
			method: "POST",
			headers,
			body,
			signal
		});
		if (resp.ok) return;
		if (resp.status !== 429) throw new Error(await resp.text());
		await sleep(wait, signal);
		wait = Math.min(wait * 2, 50);
	}
}
async function postBatch(connID, frames, signal) {
	let start = 0;
	while (start < frames.length) {
		if (frames[start].byteLength > CHUNK_THRESHOLD) {
			await postFrame(connID, KIND_DATA, frames[start], void 0, signal);
			start++;
			continue;
		}
		let end = start;
		let size = 4;
		while (end < frames.length && end - start < MAX_BATCH_FRAMES) {
			const frame = frames[end];
			if (frame.byteLength > CHUNK_THRESHOLD || size + 4 + frame.byteLength > CHUNK_THRESHOLD) break;
			size += 4 + frame.byteLength;
			end++;
		}
		if (end === start) {
			await postFrame(connID, KIND_DATA, frames[start], void 0, signal);
			start++;
			continue;
		}
		const group = frames.slice(start, end);
		if (group.length === 1) await postFrame(connID, KIND_DATA, group[0], void 0, signal);
		else await postBatchGroup(connID, group, signal);
		start = end;
	}
}
async function postBatchGroup(connID, frames, signal) {
	if (frames.length === 0 || frames.length > MAX_BATCH_FRAMES) throw new Error("invalid stream batch size");
	let body = buildBatch(frames);
	let sent = 0;
	let wait = 1;
	for (;;) {
		var _resp$headers$get;
		if (signal === null || signal === void 0 ? void 0 : signal.aborted) throw new StreamRequestCancelled();
		const resp = await fetch(streamURL("send"), {
			method: "POST",
			headers: {
				[HDR_SESSION]: G.session,
				[HDR_GENERATION]: String(G.generation),
				[HDR_CONN]: String(connID),
				[HDR_KIND]: String(KIND_DATA),
				[HDR_BATCH]: String(frames.length - sent),
				"Content-Type": "application/octet-stream"
			},
			body,
			signal
		});
		if (resp.ok) return;
		if (resp.status !== 429) throw new Error(await resp.text());
		const accepted = Number((_resp$headers$get = resp.headers.get(HDR_BATCH)) !== null && _resp$headers$get !== void 0 ? _resp$headers$get : 0);
		const remaining = frames.length - sent;
		if (!Number.isInteger(accepted) || accepted < 0 || accepted >= remaining) {
			if (accepted !== 0) throw new Error("invalid stream batch acknowledgement");
		}
		sent += accepted;
		body = buildBatch(frames.slice(sent));
		await sleep(wait, signal);
		wait = Math.min(wait * 2, 50);
	}
}
function buildBatch(frames) {
	let size = 4;
	for (const f of frames) size += 4 + f.byteLength;
	const out = new Uint8Array(size);
	const view = new DataView(out.buffer);
	view.setUint32(0, frames.length);
	let off = 4;
	for (const f of frames) {
		view.setUint32(off, f.byteLength);
		off += 4;
		out.set(f, off);
		off += f.byteLength;
	}
	return out;
}
function startPolling() {
	if (G.polling || !hasDOM) return;
	G.polling = true;
	pollLoop();
}
async function pollLoop() {
	let backoff = 0;
	const abort = new AbortController();
	G.pollAbort = abort;
	while (G.connections.size > 0) try {
		const resp = await fetch(streamURL("poll"), {
			method: "GET",
			headers: {
				[HDR_SESSION]: G.session,
				[HDR_GENERATION]: String(G.generation)
			},
			cache: "no-store",
			signal: abort.signal
		});
		if (resp.status === 410) {
			closeAll(1001, "session closed");
			break;
		}
		if (!resp.ok && resp.status !== 204) {
			if (!isRetryablePollStatus(resp.status)) {
				failAll("stream poll failed: " + resp.status);
				break;
			}
			throw new Error("retryable stream poll failure: " + resp.status);
		}
		backoff = 0;
		if (resp.status !== 204) {
			let buf;
			try {
				buf = await resp.arrayBuffer();
			} catch (_unused7) {
				failAll("stream poll response body failed");
				break;
			}
			const frames = decodeFrames(buf);
			if (frames === null) {
				closeAll(1002, "unrecognised stream framing");
				break;
			}
			deliver(frames);
		}
	} catch (err) {
		if (abort.signal.aborted || G.connections.size === 0) break;
		backoff = backoff ? Math.min(backoff * 2, BACKOFF_MAX) : BACKOFF_MIN;
		try {
			await sleep(backoff, abort.signal);
		} catch (_unused8) {
			break;
		}
	}
	if (G.pollAbort === abort) G.pollAbort = void 0;
	G.polling = false;
	if (G.connections.size > 0) startPolling();
}
function isRetryablePollStatus(status) {
	return status === 408 || status === 425 || status === 429 || status >= 500 && status <= 599;
}
function decodeFrames(buf) {
	if (buf.byteLength < 9) return null;
	const dv = new DataView(buf);
	if (dv.getUint8(0) !== 87 || dv.getUint8(1) !== 83 || dv.getUint8(2) !== 49 || dv.getUint8(3) !== 0 || (dv.getUint8(4) & -2) !== 0) return null;
	const count = dv.getUint32(5);
	const frames = [];
	let off = 9;
	for (let i = 0; i < count; i++) {
		if (off + 9 > buf.byteLength) return null;
		const connID = dv.getUint32(off);
		off += 4;
		const kind = dv.getUint8(off);
		off += 1;
		const len = dv.getUint32(off);
		off += 4;
		if (kind > KIND_ERROR || len > buf.byteLength - off) return null;
		const payload = buf.slice(off, off + len);
		off += len;
		frames.push({
			connID,
			kind,
			payload
		});
	}
	return off === buf.byteLength ? frames : null;
}
function deliver(frames) {
	for (const frame of frames) {
		const connID = frame.connID;
		const kind = frame.kind;
		const payload = frame.payload;
		const conn = G.connections.get(connID);
		if (!conn) continue;
		switch (kind) {
			case KIND_DATA:
				conn._message(payload);
				break;
			case KIND_OPEN:
				conn._opened();
				break;
			case KIND_CLOSE:
				conn._closed(1e3, "", true);
				break;
			case KIND_ERROR: conn._fail(new Error(new TextDecoder().decode(payload)));
		}
	}
}
function closeAll(code, reason) {
	for (const conn of [...G.connections.values()]) conn._closed(code, reason, false);
}
function failAll(reason) {
	for (const conn of [...G.connections.values()]) conn._fail(new Error(reason));
}
var _WailsSocket, GENERATION_KEY, GENERATION_NAME_MARKER, G, HDR_SESSION, HDR_GENERATION, HDR_CONN, HDR_KIND, HDR_NAME, HDR_CHUNK, HDR_CHUNK_INDEX, HDR_CHUNK_TOTAL, HDR_BATCH, KIND_DATA, KIND_OPEN, KIND_CLOSE, KIND_ERROR, CHUNK_THRESHOLD, MAX_BATCH_FRAMES, BACKOFF_MIN, BACKOFF_MAX, textEncoder, StreamRequestCancelled, WailsSocket;
function init_stream() {
	return (init_stream = __esmMin((() => {
		init_nanoid();
		init_environment();
		init_objectSpread2();
		;
		GENERATION_KEY = "__wails_stream_generation";
		GENERATION_NAME_MARKER = "__wails_stream_generation__:";
		G = streamGlobals();
		HDR_SESSION = "x-wails-stream-session";
		HDR_GENERATION = "x-wails-stream-generation";
		HDR_CONN = "x-wails-stream-conn";
		HDR_KIND = "x-wails-stream-kind";
		HDR_NAME = "x-wails-stream-name";
		HDR_CHUNK = "x-wails-stream-chunk";
		HDR_CHUNK_INDEX = "x-wails-stream-chunk-index";
		HDR_CHUNK_TOTAL = "x-wails-stream-chunk-total";
		HDR_BATCH = "x-wails-stream-batch";
		KIND_DATA = 0;
		KIND_OPEN = 1;
		KIND_CLOSE = 2;
		KIND_ERROR = 3;
		CHUNK_THRESHOLD = 524288;
		MAX_BATCH_FRAMES = 4096;
		BACKOFF_MIN = 250;
		BACKOFF_MAX = 5e3;
		textEncoder = new TextEncoder();
		StreamRequestCancelled = class extends Error {
			constructor() {
				super("stream request cancelled");
				this.name = "AbortError";
			}
		};
		WailsSocket = class WailsSocket extends EventTarget {
			get bufferedAmount() {
				return this._buffered;
			}
			constructor(name) {
				super();
				this.CONNECTING = 0;
				this.OPEN = 1;
				this.CLOSING = 2;
				this.CLOSED = 3;
				this.protocol = "";
				this.extensions = "";
				this.binaryType = "arraybuffer";
				this.readyState = WailsSocket.CONNECTING;
				this.onopen = null;
				this.onmessage = null;
				this.onclose = null;
				this.onerror = null;
				this._decode = (payload) => payload;
				this._buffered = 0;
				this._chain = Promise.resolve();
				this._pending = [];
				this._flushing = false;
				this._openAbort = new AbortController();
				this._sendAbort = new AbortController();
				this.name = name;
				this._id = G.nextConnID++;
				this.url = hasDOM ? streamURL("poll") + "#" + encodeURIComponent(name) : "";
				for (const type of [
					"open",
					"message",
					"close",
					"error"
				]) defineHandlerProperty(this, type);
				G.connections.set(this._id, this);
				this._chain = this._chain.then(() => postFrame(this._id, KIND_OPEN, /* @__PURE__ */ new Uint8Array(0), name, this._openAbort.signal)).catch((err) => {
					this._fail(err);
				});
				startPolling();
			}
			send(data) {
				if (this.readyState === WailsSocket.CONNECTING) throw new DOMException("Still in CONNECTING state.", "InvalidStateError");
				if (this.readyState !== WailsSocket.OPEN) return;
				const immediate = toBytesSync(data);
				const snapshot = immediate ? Promise.resolve(immediate) : toBytes(data, this._sendAbort.signal);
				snapshot.catch(() => {});
				const blobSize = !immediate && typeof Blob !== "undefined" && data instanceof Blob ? data.size : 0;
				if (blobSize) this._buffered += blobSize;
				if (immediate) this._buffered += immediate.byteLength;
				this._pending.push(snapshot);
				if (this._flushing) return;
				this._flushing = true;
				this._chain = this._chain.then(async () => {
					try {
						while (this._pending.length > 0) {
							const queued = this._pending.splice(0, MAX_BATCH_FRAMES);
							const batch = await Promise.all(queued);
							if (this.readyState === WailsSocket.CLOSED) break;
							const total = batch.reduce((n, b) => n + b.byteLength, 0);
							try {
								await postBatch(this._id, batch, this._sendAbort.signal);
							} finally {
								this._buffered = Math.max(0, this._buffered - total);
							}
						}
					} finally {
						this._flushing = false;
					}
				}).catch((err) => {
					this._flushing = false;
					this._fail(err);
				});
			}
			close(code = 1e3, reason = "") {
				if (this.readyState === WailsSocket.CLOSING || this.readyState === WailsSocket.CLOSED) return;
				this.readyState = WailsSocket.CLOSING;
				this._pending.length = 0;
				this._buffered = 0;
				this._openAbort.abort();
				this._sendAbort.abort();
				this._chain = this._chain.then(() => postFrame(this._id, KIND_CLOSE, /* @__PURE__ */ new Uint8Array(0))).catch(() => {}).then(() => {
					this._closed(code, reason, true);
				});
			}
			_opened() {
				if (this.readyState !== WailsSocket.CONNECTING) return;
				this.readyState = WailsSocket.OPEN;
				this.dispatchEvent(new Event("open"));
			}
			_message(payload) {
				if (this.readyState !== WailsSocket.OPEN) return;
				let data;
				try {
					data = this._decode(payload);
				} catch (_unused5) {
					this.dispatchEvent(new Event("error"));
					return;
				}
				if (data === payload && this.binaryType === "blob") data = new Blob([payload]);
				this.dispatchEvent(new MessageEvent("message", { data }));
			}
			_fail(err) {
				if (this.readyState === WailsSocket.CLOSED || this.readyState === WailsSocket.CLOSING) return;
				this.dispatchEvent(new Event("error"));
				this._closed(1006, err instanceof Error ? err.message : String(err), false);
			}
			_closed(code, reason, wasClean) {
				var _G$pollAbort;
				if (this.readyState === WailsSocket.CLOSED) return;
				this.readyState = WailsSocket.CLOSED;
				this._openAbort.abort();
				this._sendAbort.abort();
				G.connections.delete(this._id);
				if (G.connections.size === 0) (_G$pollAbort = G.pollAbort) === null || _G$pollAbort === void 0 || _G$pollAbort.abort();
				this._pending.length = 0;
				this._buffered = 0;
				const ev = typeof CloseEvent === "function" ? new CloseEvent("close", {
					code,
					reason,
					wasClean
				}) : Object.assign(new Event("close"), {
					code,
					reason,
					wasClean
				});
				this.dispatchEvent(ev);
			}
		};
		_WailsSocket = WailsSocket;
		_WailsSocket.CONNECTING = 0;
		_WailsSocket.OPEN = 1;
		_WailsSocket.CLOSING = 2;
		_WailsSocket.CLOSED = 3;
	})))();
}

// ../../compiled/main.js
function init_main() {
	return (init_main = __esmMin((() => {
		init_src();
		init_wml();
		init_utils();
		window.wails = src_exports;
		Enable();
		debugLog("Wails Runtime Loaded");
	})))();
}

// src/index.ts
var src_exports = /* @__PURE__ */ __exportAll({
	Android: () => android_exports,
	Application: () => application_exports,
	Browser: () => browser_exports,
	Call: () => calls_exports,
	CancelError: () => CancelError,
	CancellablePromise: () => CancellablePromise,
	CancelledRejectionError: () => CancelledRejectionError,
	Clipboard: () => clipboard_exports,
	Create: () => create_exports,
	Dialogs: () => dialogs_exports,
	Events: () => events_exports,
	Flags: () => flags_exports,
	IOS: () => ios_exports,
	JSONStream: () => JSONStream,
	Screens: () => screens_exports,
	Stream: () => Stream,
	System: () => system_exports,
	Updater: () => updater_exports,
	WML: () => wml_exports,
	WailsSocket: () => WailsSocket,
	Window: () => thisWindow,
	clientId: () => clientId,
	getTransport: () => getTransport,
	loadOptionalScript: () => loadOptionalScript,
	objectNames: () => objectNames,
	setTransport: () => setTransport
});
function loadOptionalScript(url) {
	return fetch(url, { method: "HEAD" }).then((response) => {
		if (response.ok) {
			if ((response.headers.get("content-type") || "").toLowerCase().includes("javascript")) {
				const script = document.createElement("script");
				script.src = url;
				document.head.appendChild(script);
			}
		}
	}).catch(() => {});
}
function init_src() {
	return (init_src = __esmMin((() => {
		init_environment();
		init_contextmenu();
		init_drag();
		init_appregion();
		init_application();
		init_browser();
		init_calls();
		init_clipboard();
		init_create();
		init_dialogs();
		init_events();
		/* @__PURE__ */ init_flags();
		init_screens();
		init_system();
		init_ios();
		init_android();
		init_updater();
		init_window();
		init_wml();
		init_stream();
		init_cancellable();
		init_runtime();
		init_main();
		if (hasDOM) window._wails = window._wails || {};
		if (hasDOM) {
			window._wails.invoke = invoke;
			window._wails.clientId = clientId;
		}
		if (hasDOM) window._wails.handlePlatformFileDrop = thisWindow.HandlePlatformFileDrop.bind(thisWindow);
		if (hasDOM) {
			window._wails.handleDragEnter = handleDragEnter;
			window._wails.handleDragLeave = handleDragLeave;
			window._wails.handleDragOver = handleDragOver;
		}
		if (hasDOM) invoke("wails:runtime:ready");
		if (hasDOM) loadOptionalScript("/wails/custom.js");
	})))();
}

init_src();
export { android_exports as Android, application_exports as Application, browser_exports as Browser, calls_exports as Call, CancelError, CancellablePromise, CancelledRejectionError, clipboard_exports as Clipboard, create_exports as Create, dialogs_exports as Dialogs, events_exports as Events, flags_exports as Flags, ios_exports as IOS, JSONStream, screens_exports as Screens, Stream, system_exports as System, updater_exports as Updater, wml_exports as WML, WailsSocket, thisWindow as Window, clientId, getTransport, loadOptionalScript, objectNames, setTransport };
//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicnVudGltZS5kZWJ1Zy5qcyIsIm5hbWVzIjpbIkNIVU5LX1RIUkVTSE9MRCIsImNhbGwiLCJjYWxsIiwiY2FsbCIsIkhpZGVNZXRob2QiLCJTaG93TWV0aG9kIiwiY2FsbCIsImlzQ2FsbGFibGUiLCJjYWxsIiwiY2FsbCIsIkFycmF5IiwiTWFwIiwiRXZlbnRzIiwiY2FsbCIsIkVycm9yIiwiQ3JlYXRlIiwiY2FsbCIsImNhbGwiLCJjYWxsIiwiRGV2aWNlSW5mbyIsIkhhcHRpY3MiLCJEZXZpY2UiLCJXaW5kb3ciLCJSdW50aW1lIiwiU3lzdGVtLmludm9rZSIsIldpbmRvdyJdLCJzb3VyY2VzIjpbIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9lbnZpcm9ubWVudC50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9uYW5vaWQudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9zeXN0ZW0udHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZHJhZy50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHByZWdpb24udHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb24udHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYnJvd3Nlci50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxsYWJsZS50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYW5jZWxsYWJsZS50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY3JlYXRlLnRzIiwiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RpYWxvZ3MudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvbGlzdGVuZXIudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvaW9zLnRzIiwiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FuZHJvaWQudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXBkYXRlci50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvd21sLnRzIiwiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N0cmVhbS50cyIsIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9jb21waWxlZC9tYWluLmpzIiwiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIl0sInNvdXJjZXNDb250ZW50IjpbIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBUcnVlIHdoZW4gcnVubmluZyBpbnNpZGUgYSBicm93c2VyL3dlYnZpZXcgd2l0aCBhIERPTSBhdmFpbGFibGUuXHJcbiAqIEZhbHNlIHVuZGVyIHNlcnZlci1zaWRlIHJlbmRlcmluZyAoZS5nLiBgbmV4dCBidWlsZGAgcHJlcmVuZGVyaW5nKSxcclxuICogd2hlcmUgYXBwbGljYXRpb24gY29kZSBtYXkgaW1wb3J0IHRoZSBydW50aW1lIG1vZHVsZSBldmVuIHRob3VnaCBub1xyXG4gKiBXYWlscyBBUElzIGNhbiBhY3R1YWxseSBiZSB1c2VkICgjNDY3OSkuIE1vZHVsZXMgbXVzdCBub3QgdG91Y2hcclxuICogYHdpbmRvd2AvYGRvY3VtZW50YCBhdCBpbXBvcnQgdGltZSBleGNlcHQgYmVoaW5kIHRoaXMgZ3VhcmQuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgaGFzRE9NID0gdHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgZG9jdW1lbnQgIT09IFwidW5kZWZpbmVkXCI7XHJcbiIsIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxyXG5cclxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXHJcbi8vXHJcbi8vIENvcHlyaWdodCAyMDE3IEFuZHJleSBTaXRuaWsgPGFuZHJleUBzaXRuaWsucnU+XHJcbi8vXHJcbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2ZcclxuLy8gdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbCBpblxyXG4vLyB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvXHJcbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXHJcbi8vIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbyxcclxuLy8gICAgIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxyXG4vL1xyXG4vLyAgICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXHJcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXHJcbi8vXHJcbi8vICAgICBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXHJcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTXHJcbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxyXG4vLyBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVJcclxuLy8gSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU5cclxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cclxuXHJcbi8vIFRoaXMgYWxwaGFiZXQgdXNlcyBgQS1aYS16MC05Xy1gIHN5bWJvbHMuXHJcbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cclxuLy8gUmVmZXJlbmNlcyB0byB0aGUgc2FtZSBmaWxlICh3b3JrcyBib3RoIGZvciBnemlwIGFuZCBicm90bGkpOlxyXG4vLyBgJ3VzZWAsIGBhbmRvbWAsIGFuZCBgcmljdCdgXHJcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XHJcbi8vIGAtMjZUYCwgYDE5ODNgLCBgNDBweGAsIGA3NXB4YCwgYGJ1c2hgLCBgamFja2AsIGBtaW5kYCwgYHZlcnlgLCBhbmQgYHdvbGZgXHJcbmNvbnN0IHVybEFscGhhYmV0ID1cclxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xyXG5cclxuLy8gTU9ESUZJRUQgRk9SIFdBSUxTOiB0aGUgdXBzdHJlYW0gaW1wbGVtZW50YXRpb24gZHJhd3MgZnJvbSBNYXRoLnJhbmRvbSgpLFxyXG4vLyB3aGljaCBpcyBub3QgYSBjcnlwdG9ncmFwaGljIHNvdXJjZS4gVGhlc2UgaWRzIGFyZSB1c2VkIGZvciB0aGUgcnVudGltZVxyXG4vLyBjbGllbnQgaWQsIGJpbmRpbmcgY2FsbCBpZHMsIGNodW5rIGlkcywgYW5kIGRlc2t0b3Agc3RyZWFtIHNlc3Npb24gaWRzLiBUaGV5XHJcbi8vIGFyZSBjb3JyZWxhdGlvbiBpZGVudGlmaWVycyByYXRoZXIgdGhhbiBhdXRoZW50aWNhdGlvbiBjYXBhYmlsaXRpZXMsIGJ1dCBhXHJcbi8vIHN0cm9uZ2VyIHNvdXJjZSBtYWtlcyBhY2NpZGVudGFsIGNvbGxpc2lvbnMgbmVnbGlnaWJsZSBldmVuIGFjcm9zcyBtYW55XHJcbi8vIGNvbmN1cnJlbnRseSBsb2FkZWQgcnVudGltZXMuIFRoZSBhbHBoYWJldCBpcyA2NCBjaGFyYWN0ZXJzLCBzbyBtYXNraW5nIGFcclxuLy8gcmFuZG9tIGJ5dGUgd2l0aCA2MyBpcyB1bmlmb3JtIOKAlCBubyByZWplY3Rpb24gc2FtcGxpbmcgbmVlZGVkLlxyXG5leHBvcnQgZnVuY3Rpb24gbmFub2lkKHNpemU6IG51bWJlciA9IDIxKTogc3RyaW5nIHtcclxuICAgIGNvbnN0IHNvdXJjZSA9IHR5cGVvZiBjcnlwdG8gIT09IFwidW5kZWZpbmVkXCIgJiYgdHlwZW9mIGNyeXB0by5nZXRSYW5kb21WYWx1ZXMgPT09IFwiZnVuY3Rpb25cIlxyXG4gICAgICAgID8gY3J5cHRvXHJcbiAgICAgICAgOiBudWxsO1xyXG5cclxuICAgIGlmIChzb3VyY2UpIHtcclxuICAgICAgICBjb25zdCBieXRlcyA9IG5ldyBVaW50OEFycmF5KHNpemUpO1xyXG4gICAgICAgIHNvdXJjZS5nZXRSYW5kb21WYWx1ZXMoYnl0ZXMpO1xyXG4gICAgICAgIGxldCBpZCA9ICcnO1xyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2l6ZTsgaSsrKSB7XHJcbiAgICAgICAgICAgIGlkICs9IHVybEFscGhhYmV0W2J5dGVzW2ldICYgNjNdO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gaWQ7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gRmFsbGJhY2sgZm9yIGVudmlyb25tZW50cyB3aXRob3V0IFdlYiBDcnlwdG8gKHZlcnkgb2xkIGVuZ2luZXMsIHNvbWVcclxuICAgIC8vIHNlcnZlci1zaWRlIHJlbmRlcmluZyBjb250ZXh0cykuIElkcyByZW1haW4gdW5pcXVlIGJ1dCBub3QgdW5ndWVzc2FibGUuXHJcbiAgICBsZXQgaWQgPSAnJztcclxuICAgIGxldCBpID0gc2l6ZSB8IDA7XHJcbiAgICB3aGlsZSAoaS0tKSB7XHJcbiAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXTtcclxuICAgIH1cclxuICAgIHJldHVybiBpZDtcclxufVxyXG4iLCIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcclxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcclxuXHJcbi8vIFJlc29sdmVkIGxhemlseTogd2luZG93IGRvZXMgbm90IGV4aXN0IHdoZW4gdGhlIG1vZHVsZSBpcyBpbXBvcnRlZCBkdXJpbmdcclxuLy8gc2VydmVyLXNpZGUgcmVuZGVyaW5nICgjNDY3OSksIGFuZCBub3RoaW5nIGNhbiBjYWxsIHRoZSBydW50aW1lIHRoZXJlLlxyXG5mdW5jdGlvbiBydW50aW1lVVJMKCk6IHN0cmluZyB7XHJcbiAgICByZXR1cm4gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcclxufVxyXG5cclxuLy8gU3RheSB1bmRlciBXZWJWaWV3MidzIH4yTUIgcmVxdWVzdCBib2R5IGJ1ZmZlcmluZyBsaW1pdCBpbiBXZWJSZXNvdXJjZVJlcXVlc3RlZC5cclxuY29uc3QgQ0hVTktfVEhSRVNIT0xEID0gNTEyICogMTAyNDtcclxuXHJcbi8vIFJlLWV4cG9ydCBuYW5vaWQgZm9yIGN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zXHJcbmV4cG9ydCB7IG5hbm9pZCB9O1xyXG5cclxudHlwZSBDYWxsRXJyb3JUeXBlID0ge1xyXG4gICAgbWVzc2FnZTogc3RyaW5nLFxyXG4gICAgY2F1c2U/OiB1bmtub3duLFxyXG4gICAga2luZDogXCJSZWZlcmVuY2VFcnJvclwiIHwgXCJUeXBlRXJyb3JcIiB8IFwiUnVudGltZUVycm9yXCJcclxufVxyXG5cclxuLyoqXHJcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdGhyb3duIGluIGNhc2UgdGhlIGJvdW5kIG1ldGhvZCByZXR1cm5zIGFuIGVycm9yLlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBSdW50aW1lRXJyb3IjbmFtZX0gcHJvcGVydHkgaXMgXCJSdW50aW1lRXJyb3JcIi5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBSdW50aW1lRXJyb3IgZXh0ZW5kcyBFcnJvciB7XHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgUnVudGltZUVycm9yIGluc3RhbmNlLlxyXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XHJcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJSdW50aW1lRXJyb3JcIjtcclxuICAgIH1cclxufVxyXG5cclxuLy8gT2JqZWN0IE5hbWVzXHJcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xyXG4gICAgQ2FsbDogMCxcclxuICAgIENsaXBib2FyZDogMSxcclxuICAgIEFwcGxpY2F0aW9uOiAyLFxyXG4gICAgRXZlbnRzOiAzLFxyXG4gICAgQ29udGV4dE1lbnU6IDQsXHJcbiAgICBEaWFsb2c6IDUsXHJcbiAgICBXaW5kb3c6IDYsXHJcbiAgICBTY3JlZW5zOiA3LFxyXG4gICAgU3lzdGVtOiA4LFxyXG4gICAgQnJvd3NlcjogOSxcclxuICAgIENhbmNlbENhbGw6IDEwLFxyXG4gICAgSU9TOiAxMSxcclxuICAgIEFuZHJvaWQ6IDEyLFxyXG59KTtcclxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xyXG5cclxuLyoqXHJcbiAqIFJ1bnRpbWVUcmFuc3BvcnQgZGVmaW5lcyB0aGUgaW50ZXJmYWNlIGZvciBjdXN0b20gSVBDIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnMuXHJcbiAqIEltcGxlbWVudCB0aGlzIGludGVyZmFjZSB0byB1c2UgV2ViU29ja2V0cywgY3VzdG9tIHByb3RvY29scywgb3IgYW55IG90aGVyXHJcbiAqIHRyYW5zcG9ydCBtZWNoYW5pc20gaW5zdGVhZCBvZiB0aGUgZGVmYXVsdCBIVFRQIGZldGNoLlxyXG4gKi9cclxuZXhwb3J0IGludGVyZmFjZSBSdW50aW1lVHJhbnNwb3J0IHtcclxuICAgIC8qKlxyXG4gICAgICogU2VuZCBhIHJ1bnRpbWUgY2FsbCBhbmQgcmV0dXJuIHRoZSByZXNwb25zZS5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gb2JqZWN0SUQgLSBUaGUgV2FpbHMgb2JqZWN0IElEICgwPUNhbGwsIDE9Q2xpcGJvYXJkLCBldGMuKVxyXG4gICAgICogQHBhcmFtIG1ldGhvZCAtIFRoZSBtZXRob2QgSUQgdG8gY2FsbFxyXG4gICAgICogQHBhcmFtIHdpbmRvd05hbWUgLSBPcHRpb25hbCB3aW5kb3cgbmFtZVxyXG4gICAgICogQHBhcmFtIGFyZ3MgLSBBcmd1bWVudHMgdG8gcGFzcyAod2lsbCBiZSBKU09OIHN0cmluZ2lmaWVkIGlmIHByZXNlbnQpXHJcbiAgICAgKiBAcmV0dXJucyBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzcG9uc2UgZGF0YVxyXG4gICAgICovXHJcbiAgICBjYWxsKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55PjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb24gKGNhbiBiZSBzZXQgYnkgdXNlcilcclxuICovXHJcbmxldCBjdXN0b21UcmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsID0gbnVsbDtcclxuXHJcbi8qKlxyXG4gKiBTZXQgYSBjdXN0b20gdHJhbnNwb3J0IGZvciBhbGwgV2FpbHMgcnVudGltZSBjYWxscy5cclxuICogVGhpcyBhbGxvd3MgeW91IHRvIHJlcGxhY2UgdGhlIGRlZmF1bHQgSFRUUCBmZXRjaCB0cmFuc3BvcnQgd2l0aFxyXG4gKiBXZWJTb2NrZXRzLCBjdXN0b20gcHJvdG9jb2xzLCBvciBhbnkgb3RoZXIgbWVjaGFuaXNtLlxyXG4gKlxyXG4gKiBAcGFyYW0gdHJhbnNwb3J0IC0gWW91ciBjdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uXHJcbiAqXHJcbiAqIEBleGFtcGxlXHJcbiAqIGBgYHR5cGVzY3JpcHRcclxuICogaW1wb3J0IHsgc2V0VHJhbnNwb3J0IH0gZnJvbSAnL3dhaWxzL3J1bnRpbWUuanMnO1xyXG4gKlxyXG4gKiBjb25zdCB3c1RyYW5zcG9ydCA9IHtcclxuICogICBjYWxsOiBhc3luYyAob2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncykgPT4ge1xyXG4gKiAgICAgLy8gWW91ciBXZWJTb2NrZXQgaW1wbGVtZW50YXRpb25cclxuICogICB9XHJcbiAqIH07XHJcbiAqXHJcbiAqIHNldFRyYW5zcG9ydCh3c1RyYW5zcG9ydCk7XHJcbiAqIGBgYFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIHNldFRyYW5zcG9ydCh0cmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsKTogdm9pZCB7XHJcbiAgICBjdXN0b21UcmFuc3BvcnQgPSB0cmFuc3BvcnQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIGN1cnJlbnQgdHJhbnNwb3J0ICh1c2VmdWwgZm9yIGV4dGVuZGluZy93cmFwcGluZylcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBnZXRUcmFuc3BvcnQoKTogUnVudGltZVRyYW5zcG9ydCB8IG51bGwge1xyXG4gICAgcmV0dXJuIGN1c3RvbVRyYW5zcG9ydDtcclxufVxyXG5cclxuaW50ZXJmYWNlIE9wZW5IYXJtb255SG9zdCB7XHJcbiAgICBpbnZva2UobWVzc2FnZTogc3RyaW5nKTogdW5rbm93bjtcclxuICAgIHJlcXVlc3Qoc2Vzc2lvbjogc3RyaW5nLCBpZDogc3RyaW5nLCBwYXlsb2FkOiBzdHJpbmcpOiBib29sZWFuO1xyXG59XHJcblxyXG4vLyBsaWJlbnRyeS5zbyBpbmplY3RzIHRoZXNlIGdsb2JhbHMgaW50byB0aGUgQXJrV2ViIGRvY3VtZW50LlxyXG5pbnRlcmZhY2UgT3Blbkhhcm1vbnlXaW5kb3cgZXh0ZW5kcyBXaW5kb3cge1xyXG4gICAgd2FpbHNPcGVuSGFybW9ueT86IE9wZW5IYXJtb255SG9zdDtcclxuICAgIHdhaWxzT3Blbkhhcm1vbnlTZXNzaW9uPzogc3RyaW5nO1xyXG4gICAgd2FpbHNPcGVuSGFybW9ueVJlc29sdmU/OiAoaWQ6IHN0cmluZywgcGF5bG9hZDogc3RyaW5nKSA9PiB2b2lkO1xyXG4gICAgd2FpbHM/OiB7IGludm9rZT86IChtZXNzYWdlOiBzdHJpbmcpID0+IHVua25vd247IFtrZXk6IHN0cmluZ106IHVua25vd24gfTtcclxufVxyXG5cclxuaW50ZXJmYWNlIFBlbmRpbmdPcGVuSGFybW9ueUNhbGwge1xyXG4gICAgcmVzb2x2ZTogKHZhbHVlOiB1bmtub3duKSA9PiB2b2lkO1xyXG4gICAgcmVqZWN0OiAocmVhc29uOiB1bmtub3duKSA9PiB2b2lkO1xyXG4gICAgdGltZXI6IG51bWJlcjtcclxufVxyXG5cclxuLy8gVGhlIEdvIHNpZGUgZW1pdHMgdGhpcyBlbnZlbG9wZSB3aXRoIG9rL3RleHQvZGF0YS9lcnJvciBmaWVsZHMuXHJcbmludGVyZmFjZSBSdW50aW1lUmVwbHkge1xyXG4gICAgb2s/OiBib29sZWFuO1xyXG4gICAgdGV4dD86IHVua25vd247XHJcbiAgICBkYXRhPzogdW5rbm93bjtcclxuICAgIGVycm9yPzogc3RyaW5nO1xyXG59XHJcblxyXG5sZXQgZGV0ZWN0ZWRPcGVuSGFybW9ueVRyYW5zcG9ydDogUnVudGltZVRyYW5zcG9ydCB8IG51bGwgfCB1bmRlZmluZWQ7XHJcblxyXG4vKipcclxuICogQnVpbGQgYSB0cmFuc3BvcnQgYmFja2VkIGJ5IHRoZSBBcmtUUyBicmlkZ2UgaW5qZWN0ZWQgYnkgdGhlIE9wZW5IYXJtb255XHJcbiAqIHNoZWxsLiBUaGUgc2hlbGwgY2Fubm90IHJlYWQgUE9TVCBib2RpZXMgZnJvbSBXZWJSZXNvdXJjZVJlcXVlc3Qgb24gQVBJIDI2LFxyXG4gKiBzbyBydW50aW1lIGNhbGxzIGdvIHRocm91Z2ggaXRzIEpTIHByb3h5IGluc3RlYWQgb2YgdGhlIGRlZmF1bHQgZmV0Y2guXHJcbiAqIFJldHVybnMgbnVsbCBvbiBldmVyeSBvdGhlciBob3N0LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2V0T3Blbkhhcm1vbnlUcmFuc3BvcnQoKTogUnVudGltZVRyYW5zcG9ydCB8IG51bGwge1xyXG4gICAgaWYgKGRldGVjdGVkT3Blbkhhcm1vbnlUcmFuc3BvcnQgIT09IHVuZGVmaW5lZCkgcmV0dXJuIGRldGVjdGVkT3Blbkhhcm1vbnlUcmFuc3BvcnQ7XHJcbiAgICBkZXRlY3RlZE9wZW5IYXJtb255VHJhbnNwb3J0ID0gbnVsbDtcclxuICAgIGlmICh0eXBlb2Ygd2luZG93ID09PSAndW5kZWZpbmVkJykgcmV0dXJuIG51bGw7XHJcblxyXG4gICAgY29uc3QgaG9zdFdpbmRvdyA9IHdpbmRvdyBhcyBPcGVuSGFybW9ueVdpbmRvdztcclxuICAgIGNvbnN0IGhvc3QgPSBob3N0V2luZG93LndhaWxzT3Blbkhhcm1vbnk7XHJcbiAgICBpZiAoIWhvc3QgfHwgdHlwZW9mIGhvc3QuaW52b2tlICE9PSAnZnVuY3Rpb24nIHx8IHR5cGVvZiBob3N0LnJlcXVlc3QgIT09ICdmdW5jdGlvbicpIHJldHVybiBudWxsO1xyXG5cclxuICAgIGNvbnN0IHBlbmRpbmcgPSBuZXcgTWFwPHN0cmluZywgUGVuZGluZ09wZW5IYXJtb255Q2FsbD4oKTtcclxuICAgIGNvbnN0IHNlc3Npb24gPSB3aW5kb3cuY3J5cHRvPy5yYW5kb21VVUlEPy4oKSA/PyBuYW5vaWQoKTtcclxuICAgIGxldCBuZXh0SUQgPSAwO1xyXG5cclxuICAgIC8vIFJlcGxpZXMgYXJlIHJlc29sdmVkIGludG8gdGhlIGRvY3VtZW50IHRoYXQgc3RhcnRlZCB0aGUgY2FsbCBieVxyXG4gICAgLy8gY2hlY2tpbmcgdGhpcyBzZXNzaW9uIGFuZCBpbnZva2luZyB0aGUgbWF0Y2hpbmcgcmVzb2x2ZXIuXHJcbiAgICBob3N0V2luZG93LndhaWxzT3Blbkhhcm1vbnlTZXNzaW9uID0gc2Vzc2lvbjtcclxuICAgIGhvc3RXaW5kb3cud2FpbHNPcGVuSGFybW9ueVJlc29sdmUgPSAoaWQ6IHN0cmluZywgcGF5bG9hZDogc3RyaW5nKTogdm9pZCA9PiB7XHJcbiAgICAgICAgY29uc3QgY2FsbCA9IHBlbmRpbmcuZ2V0KGlkKTtcclxuICAgICAgICBpZiAoIWNhbGwpIHJldHVybjtcclxuICAgICAgICBwZW5kaW5nLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgd2luZG93LmNsZWFyVGltZW91dChjYWxsLnRpbWVyKTtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBjb25zdCBwYXJzZWQ6IHVua25vd24gPSBKU09OLnBhcnNlKHBheWxvYWQpO1xyXG4gICAgICAgICAgICBpZiAodHlwZW9mIHBhcnNlZCAhPT0gJ29iamVjdCcgfHwgcGFyc2VkID09PSBudWxsKSB0aHJvdyBuZXcgRXJyb3IoJ0ludmFsaWQgcnVudGltZSByZXBseScpO1xyXG4gICAgICAgICAgICBjb25zdCByZXBseSA9IHBhcnNlZCBhcyBSdW50aW1lUmVwbHk7XHJcbiAgICAgICAgICAgIGlmICghcmVwbHkub2spIHRocm93IG5ldyBFcnJvcihyZXBseS5lcnJvciA/PyAnT3Blbkhhcm1vbnkgcnVudGltZSBjYWxsIGZhaWxlZCcpO1xyXG4gICAgICAgICAgICBjYWxsLnJlc29sdmUocmVwbHkudGV4dCA/PyByZXBseS5kYXRhKTtcclxuICAgICAgICB9IGNhdGNoIChlcnJvcikge1xyXG4gICAgICAgICAgICBjYWxsLnJlamVjdChlcnJvcik7XHJcbiAgICAgICAgfVxyXG4gICAgfTtcclxuXHJcbiAgICAvLyBGcm9udGVuZHMgdGhhdCBjYWxsIHdpbmRvdy53YWlscy5pbnZva2UgZGlyZWN0bHkgbmVlZCB0aGUgc2FtZSBicmlkZ2UuXHJcbiAgICBpZiAoIWhvc3RXaW5kb3cud2FpbHMpIGhvc3RXaW5kb3cud2FpbHMgPSB7fTtcclxuICAgIGlmICghaG9zdFdpbmRvdy53YWlscy5pbnZva2UpIGhvc3RXaW5kb3cud2FpbHMuaW52b2tlID0gaG9zdC5pbnZva2UuYmluZChob3N0KTtcclxuXHJcbiAgICBkZXRlY3RlZE9wZW5IYXJtb255VHJhbnNwb3J0ID0ge1xyXG4gICAgICAgIGNhbGwob2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncykge1xyXG4gICAgICAgICAgICBjb25zdCBpZCA9IFN0cmluZygrK25leHRJRCk7XHJcbiAgICAgICAgICAgIGNvbnN0IHBheWxvYWQgPSBKU09OLnN0cmluZ2lmeSh7IG9iamVjdDogb2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncywgY2xpZW50SWQgfSk7XHJcbiAgICAgICAgICAgIGNvbnN0IHsgcHJvbWlzZSwgcmVzb2x2ZSwgcmVqZWN0IH0gPSBQcm9taXNlLndpdGhSZXNvbHZlcnM8dW5rbm93bj4oKTtcclxuICAgICAgICAgICAgY29uc3QgdGltZXIgPSB3aW5kb3cuc2V0VGltZW91dCgoKSA9PiB7XHJcbiAgICAgICAgICAgICAgICBwZW5kaW5nLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgICAgICAgICByZWplY3QobmV3IEVycm9yKCdPcGVuSGFybW9ueSBydW50aW1lIGNhbGwgdGltZWQgb3V0JykpO1xyXG4gICAgICAgICAgICB9LCAxNTAwMCk7XHJcbiAgICAgICAgICAgIHBlbmRpbmcuc2V0KGlkLCB7IHJlc29sdmUsIHJlamVjdCwgdGltZXIgfSk7XHJcbiAgICAgICAgICAgIGlmICghaG9zdC5yZXF1ZXN0KHNlc3Npb24sIGlkLCBwYXlsb2FkKSkge1xyXG4gICAgICAgICAgICAgICAgcGVuZGluZy5kZWxldGUoaWQpO1xyXG4gICAgICAgICAgICAgICAgd2luZG93LmNsZWFyVGltZW91dCh0aW1lcik7XHJcbiAgICAgICAgICAgICAgICByZWplY3QobmV3IEVycm9yKCdPcGVuSGFybW9ueSBicmlkZ2UgdW5hdmFpbGFibGUnKSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICAgICAgfSxcclxuICAgIH07XHJcbiAgICByZXR1cm4gZGV0ZWN0ZWRPcGVuSGFybW9ueVRyYW5zcG9ydDtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxyXG4gKlxyXG4gKiBAcGFyYW0gb2JqZWN0IC0gVGhlIG9iamVjdCB0byBpbnZva2UgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcmV0dXJuIFRoZSBuZXcgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcclxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kOiBudW1iZXIsIGFyZ3M6IGFueSA9IG51bGwpIHtcclxuICAgICAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0LCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xyXG4gICAgfTtcclxufVxyXG5cclxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcclxuICAgIC8vIEFuIGV4cGxpY2l0IGN1c3RvbSB0cmFuc3BvcnQgYWx3YXlzIHdpbnMuIE9uIE9wZW5IYXJtb255IHRoZSBpbmplY3RlZFxyXG4gICAgLy8gQXJrVFMgYnJpZGdlIGlzIHRoZSBkZWZhdWx0IGJlY2F1c2UgUE9TVCByZXF1ZXN0IGJvZGllcyBhcmUgbm90IGV4cG9zZWRcclxuICAgIC8vIGJ5IG9uSW50ZXJjZXB0UmVxdWVzdC5cclxuICAgIGlmIChjdXN0b21UcmFuc3BvcnQpIHtcclxuICAgICAgICByZXR1cm4gY3VzdG9tVHJhbnNwb3J0LmNhbGwob2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9XHJcbiAgICBjb25zdCBvcGVuSGFybW9ueSA9IGdldE9wZW5IYXJtb255VHJhbnNwb3J0KCk7XHJcbiAgICBpZiAob3Blbkhhcm1vbnkpIHtcclxuICAgICAgICByZXR1cm4gb3Blbkhhcm1vbnkuY2FsbChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBEZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKCkpO1xyXG5cclxuICAgIGxldCBib2R5OiB7IG9iamVjdDogbnVtYmVyOyBtZXRob2Q6IG51bWJlciwgYXJncz86IGFueSB9ID0ge1xyXG4gICAgICBvYmplY3Q6IG9iamVjdElELFxyXG4gICAgICBtZXRob2RcclxuICAgIH1cclxuICAgIGlmIChhcmdzICE9PSBudWxsICYmIGFyZ3MgIT09IHVuZGVmaW5lZCkge1xyXG4gICAgICBib2R5LmFyZ3MgPSBhcmdzO1xyXG4gICAgfVxyXG5cclxuICAgIGxldCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+ID0ge1xyXG4gICAgICAgIFtcIngtd2FpbHMtY2xpZW50LWlkXCJdOiBjbGllbnRJZCxcclxuICAgICAgICBbXCJDb250ZW50LVR5cGVcIl06IFwiYXBwbGljYXRpb24vanNvblwiXHJcbiAgICB9XHJcbiAgICBpZiAod2luZG93TmFtZSkge1xyXG4gICAgICAgIGhlYWRlcnNbXCJ4LXdhaWxzLXdpbmRvdy1uYW1lXCJdID0gd2luZG93TmFtZTtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCBib2R5U3RyID0gSlNPTi5zdHJpbmdpZnkoYm9keSk7XHJcbiAgICBsZXQgcmVzcG9uc2U6IFJlc3BvbnNlO1xyXG4gICAgaWYgKGJvZHlTdHIubGVuZ3RoID4gQ0hVTktfVEhSRVNIT0xEKSB7XHJcbiAgICAgICAgcmVzcG9uc2UgPSBhd2FpdCBzZW5kQ2h1bmtlZCh1cmwsIGhlYWRlcnMsIGJvZHlTdHIpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICByZXNwb25zZSA9IGF3YWl0IGZldGNoKHVybCwgeyBtZXRob2Q6ICdQT1NUJywgaGVhZGVycywgYm9keTogYm9keVN0ciB9KTtcclxuICAgIH1cclxuICAgIGlmICghcmVzcG9uc2Uub2spIHtcclxuICAgICAgY29uc3QgY3QgPSByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKTtcclxuICAgICAgaWYgKGN0Py5pbmNsdWRlcyhcImFwcGxpY2F0aW9uL2pzb25cIikpIHtcclxuICAgICAgICAgIGNvbnN0IGpzb246IENhbGxFcnJvclR5cGUgPSBhd2FpdCByZXNwb25zZS5qc29uKCk7XHJcbiAgICAgICAgICBsZXQgZXJyO1xyXG4gICAgICAgICAgc3dpdGNoIChqc29uLmtpbmQpIHtcclxuICAgICAgICAgICAgICBjYXNlIFwiUmVmZXJlbmNlRXJyb3JcIjogZXJyID0gbmV3IFJlZmVyZW5jZUVycm9yKGpzb24ubWVzc2FnZSk7IGJyZWFrO1xyXG4gICAgICAgICAgICAgIGNhc2UgXCJUeXBlRXJyb3JcIjogICAgICBlcnIgPSBuZXcgVHlwZUVycm9yKGpzb24ubWVzc2FnZSk7IGJyZWFrO1xyXG4gICAgICAgICAgICAgIGNhc2UgXCJSdW50aW1lRXJyb3JcIjogICBlcnIgPSBuZXcgUnVudGltZUVycm9yKGpzb24ubWVzc2FnZSk7IGJyZWFrO1xyXG4gICAgICAgICAgICAgIGRlZmF1bHQ6ICAgICAgICAgICAgICAgZXJyID0gbmV3IEVycm9yKGpzb24ubWVzc2FnZSk7XHJcbiAgICAgICAgICB9XHJcbiAgICAgICAgICBlcnIuY2F1c2UgPSBqc29uLmNhdXNlO1xyXG4gICAgICAgICAgdGhyb3cgZXJyXHJcbiAgICAgIH1cclxuICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIHNlbmRDaHVua2VkIHNwbGl0cyBhIGxhcmdlIHNlcmlhbGlzZWQgcmVxdWVzdCBib2R5IGludG8gQ0hVTktfVEhSRVNIT0xELXNpemVkXHJcbi8vIGJ5dGUgY2h1bmtzIGFuZCBzZW5kcyB0aGVtIHNlcmlhbGx5LiAgRW5jb2RpbmcgdG8gVVRGLTggYnl0ZXMgYmVmb3JlIHNsaWNpbmdcclxuLy8gcHJldmVudHMgY29ycnVwdGlvbiBvZiBub24tQk1QIGNoYXJhY3RlcnMgKHN1cnJvZ2F0ZSBwYWlycykgdGhhdCB3b3VsZCBvY2N1clxyXG4vLyB3aGVuIHNwbGl0dGluZyBhdCBKYXZhU2NyaXB0IHN0cmluZyBpbmRpY2VzLiAgVGhlIEdvIHRyYW5zcG9ydCBhc3NlbWJsZXMgdGhlXHJcbi8vIHJhdyBieXRlcyBiZWZvcmUgcHJvY2Vzc2luZy4gIE9ubHkgdGhlIGZpbmFsIGNodW5rJ3MgcmVzcG9uc2UgY2FycmllcyB0aGUgUlBDIHJlc3VsdC5cclxuYXN5bmMgZnVuY3Rpb24gc2VuZENodW5rZWQodXJsOiBVUkwsIGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4sIGJvZHlTdHI6IHN0cmluZyk6IFByb21pc2U8UmVzcG9uc2U+IHtcclxuICAgIGNvbnN0IGNodW5rSWQgPSBuYW5vaWQoKTtcclxuICAgIGNvbnN0IGJvZHlCeXRlcyA9IG5ldyBUZXh0RW5jb2RlcigpLmVuY29kZShib2R5U3RyKTtcclxuICAgIGNvbnN0IHRvdGFsQ2h1bmtzID0gTWF0aC5jZWlsKGJvZHlCeXRlcy5sZW5ndGggLyBDSFVOS19USFJFU0hPTEQpO1xyXG5cclxuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgdG90YWxDaHVua3MgLSAxOyBpKyspIHtcclxuICAgICAgICBjb25zdCBjaHVuayA9IGJvZHlCeXRlcy5zdWJhcnJheShpICogQ0hVTktfVEhSRVNIT0xELCAoaSArIDEpICogQ0hVTktfVEhSRVNIT0xEKTtcclxuICAgICAgICBjb25zdCByZXNwID0gYXdhaXQgZmV0Y2godXJsLCB7XHJcbiAgICAgICAgICAgIG1ldGhvZDogJ1BPU1QnLFxyXG4gICAgICAgICAgICBoZWFkZXJzOiB7XHJcbiAgICAgICAgICAgICAgICAuLi5oZWFkZXJzLFxyXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaWQnOiBjaHVua0lkLFxyXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaW5kZXgnOiBTdHJpbmcoaSksXHJcbiAgICAgICAgICAgICAgICAneC13YWlscy1jaHVuay10b3RhbCc6IFN0cmluZyh0b3RhbENodW5rcyksXHJcbiAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgIGJvZHk6IGNodW5rLFxyXG4gICAgICAgIH0pO1xyXG4gICAgICAgIGlmICghcmVzcC5vaykge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gZmV0Y2godXJsLCB7XHJcbiAgICAgICAgbWV0aG9kOiAnUE9TVCcsXHJcbiAgICAgICAgaGVhZGVyczoge1xyXG4gICAgICAgICAgICAuLi5oZWFkZXJzLFxyXG4gICAgICAgICAgICAneC13YWlscy1jaHVuay1pZCc6IGNodW5rSWQsXHJcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWluZGV4JzogU3RyaW5nKHRvdGFsQ2h1bmtzIC0gMSksXHJcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLXRvdGFsJzogU3RyaW5nKHRvdGFsQ2h1bmtzKSxcclxuICAgICAgICB9LFxyXG4gICAgICAgIGJvZHk6IGJvZHlCeXRlcy5zdWJhcnJheSgodG90YWxDaHVua3MgLSAxKSAqIENIVU5LX1RIUkVTSE9MRCksXHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEFuZHJvaWQgV2ViVmlldyBjYW5ub3QgZGVsaXZlciBmZXRjaCgpIFBPU1QgYm9kaWVzIHRvXHJcbiAqIHNob3VsZEludGVyY2VwdFJlcXVlc3QsIHNvIHRoZSBkZWZhdWx0IEhUVFAgdHJhbnNwb3J0IGNhbm5vdCByZWFjaCBHby5cclxuICogV2hlbiB0aGUgQW5kcm9pZCBKYXZhc2NyaXB0SW50ZXJmYWNlIGJyaWRnZSAod2luZG93LndhaWxzKSBpcyBwcmVzZW50LFxyXG4gKiByb3V0ZSBydW50aW1lIGNhbGxzIHRocm91Z2ggaXQgaW5zdGVhZC4gUmVzcG9uc2VzIGFycml2ZSB2aWFcclxuICogd2luZG93Ll93YWlsc0FuZHJvaWRDYWxsYmFjaywgaW52b2tlZCBieSB0aGUgSmF2YSBzaWRlLlxyXG4gKi9cclxuaW50ZXJmYWNlIEFuZHJvaWRKU0JyaWRnZSB7XHJcbiAgICBpbnZva2VBc3luYyhjYWxsYmFja0lEOiBzdHJpbmcsIHBheWxvYWQ6IHN0cmluZyk6IHZvaWQ7XHJcbn1cclxuXHJcbmNvbnN0IGFuZHJvaWRCcmlkZ2U6IEFuZHJvaWRKU0JyaWRnZSB8IG51bGwgPSBoYXNET00gJiZcclxuICAgIHR5cGVvZiAod2luZG93IGFzIGFueSkud2FpbHM/Lmludm9rZUFzeW5jID09PSBcImZ1bmN0aW9uXCIgPyAod2luZG93IGFzIGFueSkud2FpbHMgOiBudWxsO1xyXG5cclxuaWYgKGFuZHJvaWRCcmlkZ2UpIHtcclxuICAgIGNvbnN0IHBlbmRpbmcgPSBuZXcgTWFwPHN0cmluZywgeyByZXNvbHZlOiAodmFsdWU6IGFueSkgPT4gdm9pZDsgcmVqZWN0OiAocmVhc29uOiBhbnkpID0+IHZvaWQgfT4oKTtcclxuXHJcbiAgICAod2luZG93IGFzIGFueSkuX3dhaWxzQW5kcm9pZENhbGxiYWNrID0gKGlkOiBzdHJpbmcsIHJlc3BvbnNlOiBzdHJpbmcgfCBudWxsLCBlcnJvcjogc3RyaW5nIHwgbnVsbCkgPT4ge1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBwZW5kaW5nLmdldChpZCk7XHJcbiAgICAgICAgaWYgKCFwcm9taXNlKSByZXR1cm47XHJcbiAgICAgICAgcGVuZGluZy5kZWxldGUoaWQpO1xyXG4gICAgICAgIGlmIChlcnJvcikge1xyXG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgRXJyb3IoZXJyb3IpKTtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBjb25zdCBlbnZlbG9wZSA9IEpTT04ucGFyc2UocmVzcG9uc2UgPz8gXCJ7fVwiKTtcclxuICAgICAgICAgICAgaWYgKCFlbnZlbG9wZS5vaykge1xyXG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QobmV3IEVycm9yKGVudmVsb3BlLmVycm9yID8/IFwidW5rbm93biBydW50aW1lIGNhbGwgZXJyb3JcIikpO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIHByb21pc2UucmVzb2x2ZShcInRleHRcIiBpbiBlbnZlbG9wZSA/IGVudmVsb3BlLnRleHQgOiBlbnZlbG9wZS5kYXRhKTtcclxuICAgICAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KGUpO1xyXG4gICAgICAgIH1cclxuICAgIH07XHJcblxyXG4gICAgY3VzdG9tVHJhbnNwb3J0ID0ge1xyXG4gICAgICAgIGNhbGwob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcclxuICAgICAgICAgICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgICAgIGNvbnN0IGlkID0gbmFub2lkKCk7XHJcbiAgICAgICAgICAgICAgICBwZW5kaW5nLnNldChpZCwgeyByZXNvbHZlLCByZWplY3QgfSk7XHJcbiAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgIGFuZHJvaWRCcmlkZ2UuaW52b2tlQXN5bmMoaWQsIEpTT04uc3RyaW5naWZ5KHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgb2JqZWN0OiBvYmplY3RJRCxcclxuICAgICAgICAgICAgICAgICAgICAgICAgbWV0aG9kOiBtZXRob2QsXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHdpbmRvd05hbWU6IHdpbmRvd05hbWUsXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGFyZ3M6IGFyZ3MgPz8gbnVsbCxcclxuICAgICAgICAgICAgICAgICAgICAgICAgY2xpZW50SWQ6IGNsaWVudElkLFxyXG4gICAgICAgICAgICAgICAgICAgIH0pKTtcclxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICAgICAgICAgICAgICAvLyBEb24ndCBsZWFrIHRoZSBwZW5kaW5nIGVudHJ5IGlmIGRpc3BhdGNoIHRocm93cyBzeW5jaHJvbm91c2x5XHJcbiAgICAgICAgICAgICAgICAgICAgcGVuZGluZy5kZWxldGUoaWQpO1xyXG4gICAgICAgICAgICAgICAgICAgIHJlamVjdChlKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSk7XHJcbiAgICAgICAgfSxcclxuICAgIH07XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TeXN0ZW0pO1xyXG5cclxuY29uc3QgU3lzdGVtSXNEYXJrTW9kZSA9IDA7XHJcbmNvbnN0IFN5c3RlbUVudmlyb25tZW50ID0gMTtcclxuY29uc3QgU3lzdGVtQ2FwYWJpbGl0aWVzID0gMjtcclxuXHJcbmNvbnN0IF9pbnZva2UgPSAoZnVuY3Rpb24gKCkge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICAvLyBXaW5kb3dzIFdlYlZpZXcyXHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlKSB7XHJcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcpO1xyXG4gICAgICAgIH1cclxuICAgICAgICAvLyBtYWNPUy9pT1MgV0tXZWJWaWV3XHJcbiAgICAgICAgZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndlYmtpdD8ubWVzc2FnZUhhbmRsZXJzPy5bJ2V4dGVybmFsJ10/LnBvc3RNZXNzYWdlKSB7XHJcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXS5wb3N0TWVzc2FnZS5iaW5kKCh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgLy8gQW5kcm9pZCBXZWJWaWV3IC0gdXNlcyBhZGRKYXZhc2NyaXB0SW50ZXJmYWNlIHdoaWNoIGV4cG9zZXMgd2luZG93LndhaWxzLmludm9rZVxyXG4gICAgICAgIGVsc2UgaWYgKCh3aW5kb3cgYXMgYW55KS53YWlscz8uaW52b2tlKSB7XHJcbiAgICAgICAgICAgIGNvbnN0IGJyaWRnZSA9ICh3aW5kb3cgYXMgYW55KS53YWlscztcclxuICAgICAgICAgICAgcmV0dXJuIChtc2c6IGFueSkgPT4gYnJpZGdlLmludm9rZSh0eXBlb2YgbXNnID09PSAnc3RyaW5nJyA/IG1zZyA6IEpTT04uc3RyaW5naWZ5KG1zZykpO1xyXG4gICAgICAgIH1cclxuICAgIH0gY2F0Y2goZSkge31cclxuXHJcbiAgICBjb25zb2xlLndhcm4oJ1xcbiVj4pqg77iPIEJyb3dzZXIgRW52aXJvbm1lbnQgRGV0ZWN0ZWQgJWNcXG5cXG4lY09ubHkgVUkgcHJldmlld3MgYXJlIGF2YWlsYWJsZSBpbiB0aGUgYnJvd3Nlci4gRm9yIGZ1bGwgZnVuY3Rpb25hbGl0eSwgcGxlYXNlIHJ1biB0aGUgYXBwbGljYXRpb24gaW4gZGVza3RvcCBtb2RlLlxcbk1vcmUgaW5mb3JtYXRpb24gYXQ6IGh0dHBzOi8vdjMud2FpbHMuaW8vbGVhcm4vYnVpbGQvI3VzaW5nLWEtYnJvd3Nlci1mb3ItZGV2ZWxvcG1lbnRcXG4nLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjZmZmZmZmOyBjb2xvcjogIzAwMDAwMDsgZm9udC13ZWlnaHQ6IGJvbGQ7IHBhZGRpbmc6IDRweCA4cHg7IGJvcmRlci1yYWRpdXM6IDRweDsgYm9yZGVyOiAycHggc29saWQgIzAwMDAwMDsnLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiB0cmFuc3BhcmVudDsnLFxyXG4gICAgICAgICdjb2xvcjogI2ZmZmZmZjsgZm9udC1zdHlsZTogaXRhbGljOyBmb250LXdlaWdodDogYm9sZDsnKTtcclxuICAgIHJldHVybiBudWxsO1xyXG59KSgpO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2c6IGFueSk6IHZvaWQge1xyXG4gICAgX2ludm9rZT8uKG1zZyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtSXNEYXJrTW9kZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBGZXRjaGVzIHRoZSBjYXBhYmlsaXRpZXMgb2YgdGhlIGFwcGxpY2F0aW9uIGZyb20gdGhlIHNlcnZlci5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cclxuICovXHJcbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKTogUHJvbWlzZTxSZWNvcmQ8c3RyaW5nLCBhbnk+PiB7XHJcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1DYXBhYmlsaXRpZXMpO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIE9TSW5mbyB7XHJcbiAgICAvKiogVGhlIGJyYW5kaW5nIG9mIHRoZSBPUy4gKi9cclxuICAgIEJyYW5kaW5nOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIElEIG9mIHRoZSBPUy4gKi9cclxuICAgIElEOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIG5hbWUgb2YgdGhlIE9TLiAqL1xyXG4gICAgTmFtZTogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSB2ZXJzaW9uIG9mIHRoZSBPUy4gKi9cclxuICAgIFZlcnNpb246IHN0cmluZztcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBFbnZpcm9ubWVudEluZm8ge1xyXG4gICAgLyoqIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS4gKi9cclxuICAgIEFyY2g6IHN0cmluZztcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS4gKi9cclxuICAgIERlYnVnOiBib29sZWFuO1xyXG4gICAgLyoqIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS4gKi9cclxuICAgIE9TOiBzdHJpbmc7XHJcbiAgICAvKiogRGV0YWlscyBvZiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cclxuICAgIE9TSW5mbzogT1NJbmZvO1xyXG4gICAgLyoqIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uICovXHJcbiAgICBQbGF0Zm9ybUluZm86IFJlY29yZDxzdHJpbmcsIGFueT47XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgZW52aXJvbm1lbnQgZGV0YWlscy5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKTogUHJvbWlzZTxFbnZpcm9ubWVudEluZm8+IHtcclxuICAgIHJldHVybiBjYWxsKFN5c3RlbUVudmlyb25tZW50KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MuXHJcbiAqXHJcbiAqIEByZXR1cm4gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJ3aW5kb3dzXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cclxuICpcclxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJsaW51eFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwiZGFyd2luXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBpT1MuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgaU9TLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNJT1MoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TID09PSBcImlvc1wiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgQW5kcm9pZC5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBBbmRyb2lkLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBbmRyb2lkKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJhbmRyb2lkXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGFwcCBpcyBydW5uaW5nIG9uIGEgbW9iaWxlIE9TIChpT1Mgb3IgQW5kcm9pZCkuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgb24gaU9TIG9yIEFuZHJvaWQsIG90aGVyd2lzZSBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc01vYmlsZSgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiBJc0lPUygpIHx8IElzQW5kcm9pZCgpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBhcHAgaXMgcnVubmluZyBvbiBhIGRlc2t0b3AgT1MgKG1hY09TLCBXaW5kb3dzIG9yIExpbnV4KS5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBvbiBtYWNPUywgV2luZG93cyBvciBMaW51eCwgb3RoZXJ3aXNlIGZhbHNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGVza3RvcCgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiBJc01hYygpIHx8IElzV2luZG93cygpIHx8IElzTGludXgoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhbWQ2NFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFybVwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cclxuICpcclxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUsIG90aGVyd2lzZSByZXR1cm5zIGZhbHNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYXJtNjRcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlcG9ydHMgd2hldGhlciB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiBCb29sZWFuKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5EZWJ1Zyk7XHJcbn1cclxuXHJcbiIsIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxyXG4gKlxyXG4gKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBtZXNzYWdlIHRvIGJlIGxvZ2dlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlOiBhbnkpIHtcclxuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG4gICAgY29uc29sZS5sb2coXHJcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXHJcbiAgICApO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIHdoZXRoZXIgdGhlIHdlYnZpZXcgc3VwcG9ydHMgdGhlIHtAbGluayBNb3VzZUV2ZW50I2J1dHRvbnN9IHByb3BlcnR5LlxyXG4gKiBMb29raW5nIGF0IHlvdSBtYWNPUyBIaWdoIFNpZXJyYSFcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBjYW5UcmFja0J1dHRvbnMoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKG5ldyBNb3VzZUV2ZW50KCdtb3VzZWRvd24nKSkuYnV0dG9ucyA9PT0gMDtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXHJcbiAqIChzZWUgaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmcvZW4tVVMvZG9jcy9XZWIvQVBJL0V2ZW50VGFyZ2V0L2FkZEV2ZW50TGlzdGVuZXIjc2lnbmFsKS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBjYW5BYm9ydExpc3RlbmVycygpIHtcclxuICAgIGlmICghRXZlbnRUYXJnZXQgfHwgIUFib3J0U2lnbmFsIHx8ICFBYm9ydENvbnRyb2xsZXIpXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG5cclxuICAgIGxldCByZXN1bHQgPSB0cnVlO1xyXG5cclxuICAgIGNvbnN0IHRhcmdldCA9IG5ldyBFdmVudFRhcmdldCgpO1xyXG4gICAgY29uc3QgY29udHJvbGxlciA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xyXG4gICAgY29udHJvbGxlci5hYm9ydCgpO1xyXG4gICAgdGFyZ2V0LmRpc3BhdGNoRXZlbnQobmV3IEN1c3RvbUV2ZW50KCd0ZXN0JykpO1xyXG5cclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXNvbHZlcyB0aGUgY2xvc2VzdCBIVE1MRWxlbWVudCBhbmNlc3RvciBvZiBhbiBldmVudCdzIHRhcmdldC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBldmVudFRhcmdldChldmVudDogRXZlbnQpOiBIVE1MRWxlbWVudCB7XHJcbiAgICBpZiAoZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpIHtcclxuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0O1xyXG4gICAgfSBlbHNlIGlmICghKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSAmJiBldmVudC50YXJnZXQgaW5zdGFuY2VvZiBOb2RlKSB7XHJcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldC5wYXJlbnRFbGVtZW50ID8/IGRvY3VtZW50LmJvZHk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiBkb2N1bWVudC5ib2R5O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKioqXHJcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcclxuXHJcbiBCU0QgMi1DbGF1c2UgTGljZW5zZVxyXG5cclxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxyXG4gQWxsIHJpZ2h0cyByZXNlcnZlZC5cclxuXHJcbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcclxuIG1vZGlmaWNhdGlvbiwgYXJlIHBlcm1pdHRlZCBwcm92aWRlZCB0aGF0IHRoZSBmb2xsb3dpbmcgY29uZGl0aW9ucyBhcmUgbWV0OlxyXG5cclxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xyXG4gbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIuXHJcblxyXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxyXG4gdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvblxyXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXHJcblxyXG4gVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyBcIkFTIElTXCJcclxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcclxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxyXG4gRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1JTIEJFIExJQUJMRVxyXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcclxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXHJcbiBTRVJWSUNFUzsgTE9TUyBPRiBVU0UsIERBVEEsIE9SIFBST0ZJVFM7IE9SIEJVU0lORVNTIElOVEVSUlVQVElPTikgSE9XRVZFUlxyXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcclxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXHJcbiBPRiBUSElTIFNPRlRXQVJFLCBFVkVOIElGIEFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLlxyXG5cclxuICoqKi9cclxuXHJcbmxldCBpc1JlYWR5ID0gZmFsc2U7XHJcbmltcG9ydCB7IGhhc0RPTSB9IGZyb20gXCIuL2Vudmlyb25tZW50LmpzXCI7XHJcblxyXG5pZiAoaGFzRE9NKSB7XHJcbiAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgKCkgPT4geyBpc1JlYWR5ID0gdHJ1ZSB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjazogKCkgPT4gdm9pZCkge1xyXG4gICAgaWYgKGlzUmVhZHkgfHwgZG9jdW1lbnQucmVhZHlTdGF0ZSA9PT0gJ2NvbXBsZXRlJykge1xyXG4gICAgICAgIGNhbGxiYWNrKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCBjYWxsYmFjayk7XHJcbiAgICB9XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB7IElzRGVidWcgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcclxuaW1wb3J0IHsgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xyXG5cclxuLy8gc2V0dXBcclxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcclxuXHJcbmlmIChoYXNET00pIHtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XHJcbn1cclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcclxuXHJcbmNvbnN0IENvbnRleHRNZW51T3BlbiA9IDA7XHJcblxyXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQ6IHN0cmluZywgeDogbnVtYmVyLCB5OiBudW1iZXIsIGRhdGE6IGFueSk6IHZvaWQge1xyXG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudDogTW91c2VFdmVudCkge1xyXG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xyXG5cclxuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XHJcbiAgICBjb25zdCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpLnRyaW0oKTtcclxuXHJcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIGNvbnN0IGRhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xyXG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZGF0YSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQsIHRhcmdldCk7XHJcbiAgICB9XHJcbn1cclxuXHJcblxyXG4vKlxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuXHJcblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXHJcbiovXHJcbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQ6IE1vdXNlRXZlbnQsIHRhcmdldDogSFRNTEVsZW1lbnQpIHtcclxuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxyXG4gICAgaWYgKElzRGVidWcoKSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICBzd2l0Y2ggKHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCkpIHtcclxuICAgICAgICBjYXNlICdzaG93JzpcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGNhc2UgJ2hpZGUnOlxyXG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcclxuICAgIGlmICh0YXJnZXQuaXNDb250ZW50RWRpdGFibGUpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxyXG4gICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xyXG4gICAgY29uc3QgaGFzU2VsZWN0aW9uID0gc2VsZWN0aW9uICYmIHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDA7XHJcbiAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XHJcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XHJcbiAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XHJcbiAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcclxuICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xyXG4gICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xyXG4gICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IHRhcmdldCkge1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiB0YWcgaXMgaW5wdXQgb3IgdGV4dGFyZWEuXHJcbiAgICBpZiAodGFyZ2V0IGluc3RhbmNlb2YgSFRNTElucHV0RWxlbWVudCB8fCB0YXJnZXQgaW5zdGFuY2VvZiBIVE1MVGV4dEFyZWFFbGVtZW50KSB7XHJcbiAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIXRhcmdldC5yZWFkT25seSAmJiAhdGFyZ2V0LmRpc2FibGVkKSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyB0aGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5IGZyb20gdGhlIGZsYWcgbWFwLlxyXG4gKlxyXG4gKiBAcGFyYW0ga2V5IC0gVGhlIGtleSB0byByZXRyaWV2ZSB0aGUgdmFsdWUgZm9yLlxyXG4gKiBAcmV0dXJuIFRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXk6IHN0cmluZyk6IGFueSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleV07XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiVW5hYmxlIHRvIHJldHJpZXZlIGZsYWcgJ1wiICsga2V5ICsgXCInOiBcIiArIGUsIHsgY2F1c2U6IGUgfSk7XHJcbiAgICB9XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgaW52b2tlLCBJc1dpbmRvd3MsIElzTGludXgsIElzTWFjIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XHJcbmltcG9ydCB7IEdldEZsYWcgfSBmcm9tIFwiLi9mbGFncy5qc1wiO1xyXG5pbXBvcnQgeyBjYW5UcmFja0J1dHRvbnMsIGV2ZW50VGFyZ2V0IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcclxuXHJcbi8vIFNldHVwXHJcbmxldCBjYW5EcmFnID0gZmFsc2U7XHJcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xyXG5cclxubGV0IHJlc2l6YWJsZSA9IGZhbHNlO1xyXG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XHJcbmxldCByZXNpemluZyA9IGZhbHNlO1xyXG5sZXQgcmVzaXplRWRnZTogc3RyaW5nID0gXCJcIjtcclxubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcclxuXHJcbmxldCBidXR0b25zID0gMDtcclxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcclxuXHJcbmxldCBidXR0b25zVHJhY2tlZCA9IGZhbHNlO1xyXG5cclxuaWYgKGhhc0RPTSkge1xyXG4gICAgYnV0dG9uc1RyYWNrZWQgPSBjYW5UcmFja0J1dHRvbnMoKTtcclxuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG4gICAgd2luZG93Ll93YWlscy5zZXRSZXNpemFibGUgPSAodmFsdWU6IGJvb2xlYW4pOiB2b2lkID0+IHtcclxuICAgICAgICByZXNpemFibGUgPSB2YWx1ZTtcclxuICAgICAgICBpZiAoIXJlc2l6YWJsZSkge1xyXG4gICAgICAgICAgICAvLyBTdG9wIHJlc2l6aW5nIGlmIGluIHByb2dyZXNzLlxyXG4gICAgICAgICAgICBjYW5SZXNpemUgPSByZXNpemluZyA9IGZhbHNlO1xyXG4gICAgICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgICAgICB9XHJcbiAgICB9O1xyXG59XHJcblxyXG4vLyBEZWZlciBhdHRhY2hpbmcgbW91c2UgbGlzdGVuZXJzIHVudGlsIHdlIGtub3cgd2UncmUgbm90IG9uIG1vYmlsZS5cclxubGV0IGRyYWdJbml0RG9uZSA9IGZhbHNlO1xyXG5mdW5jdGlvbiBpc01vYmlsZSgpOiBib29sZWFuIHtcclxuICAgIGNvbnN0IG9zID0gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TO1xyXG4gICAgaWYgKG9zID09PSBcImlvc1wiIHx8IG9zID09PSBcImFuZHJvaWRcIikgcmV0dXJuIHRydWU7XHJcbiAgICAvLyBGYWxsYmFjayBoZXVyaXN0aWMgaWYgZW52aXJvbm1lbnQgbm90IHlldCBzZXRcclxuICAgIGNvbnN0IHVhID0gbmF2aWdhdG9yLnVzZXJBZ2VudCB8fCBuYXZpZ2F0b3IudmVuZG9yIHx8ICh3aW5kb3cgYXMgYW55KS5vcGVyYSB8fCBcIlwiO1xyXG4gICAgcmV0dXJuIC9hbmRyb2lkfGlwaG9uZXxpcGFkfGlwb2R8aWVtb2JpbGV8d3BkZXNrdG9wL2kudGVzdCh1YSk7XHJcbn1cclxuZnVuY3Rpb24gdHJ5SW5pdERyYWdIYW5kbGVycygpOiB2b2lkIHtcclxuICAgIGlmIChkcmFnSW5pdERvbmUpIHJldHVybjtcclxuICAgIGlmIChpc01vYmlsZSgpKSByZXR1cm47XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xyXG4gICAgZm9yIChjb25zdCBldiBvZiBbJ2NsaWNrJywgJ2NvbnRleHRtZW51JywgJ2RibGNsaWNrJ10pIHtcclxuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihldiwgc3VwcHJlc3NFdmVudCwgeyBjYXB0dXJlOiB0cnVlIH0pO1xyXG4gICAgfVxyXG4gICAgZHJhZ0luaXREb25lID0gdHJ1ZTtcclxufVxyXG5pZiAoaGFzRE9NKSB7XHJcbiAgICAvLyBBdHRlbXB0IGltbWVkaWF0ZSBpbml0IChpbiBjYXNlIGVudmlyb25tZW50IGFscmVhZHkgcHJlc2VudClcclxuICAgIHRyeUluaXREcmFnSGFuZGxlcnMoKTtcclxuICAgIC8vIEFsc28gYXR0ZW1wdCBvbiBET00gcmVhZHlcclxuICAgIGRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCB0cnlJbml0RHJhZ0hhbmRsZXJzLCB7IG9uY2U6IHRydWUgfSk7XHJcbiAgICAvLyBBcyBhIGxhc3QgcmVzb3J0LCBwb2xsIGZvciBlbnZpcm9ubWVudCBmb3IgYSBzaG9ydCBwZXJpb2RcclxuICAgIGxldCBkcmFnRW52UG9sbHMgPSAwO1xyXG4gICAgY29uc3QgZHJhZ0VudlBvbGwgPSB3aW5kb3cuc2V0SW50ZXJ2YWwoKCkgPT4ge1xyXG4gICAgICAgIGlmIChkcmFnSW5pdERvbmUpIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyByZXR1cm47IH1cclxuICAgICAgICB0cnlJbml0RHJhZ0hhbmRsZXJzKCk7XHJcbiAgICAgICAgaWYgKCsrZHJhZ0VudlBvbGxzID4gMTAwKSB7IHdpbmRvdy5jbGVhckludGVydmFsKGRyYWdFbnZQb2xsKTsgfVxyXG4gICAgfSwgNTApO1xyXG59XHJcblxyXG5mdW5jdGlvbiBzdXBwcmVzc0V2ZW50KGV2ZW50OiBFdmVudCkge1xyXG4gICAgaWYgKFxyXG4gICAgICAgIGV2ZW50LnR5cGUgPT09IFwiZGJsY2xpY2tcIlxyXG4gICAgICAgICYmIElzTWFjKClcclxuICAgICAgICAmJiB3aW5kb3cgPT09IHdpbmRvdy50b3BcclxuICAgICAgICAmJiBpc0RyYWdnYWJsZUV2ZW50KGV2ZW50IGFzIE1vdXNlRXZlbnQpXHJcbiAgICApIHtcclxuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIGludm9rZShcIndhaWxzOmRyYWc6ZG91YmxlY2xpY2tcIik7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFN1cHByZXNzIGNsaWNrIGV2ZW50cyB3aGlsZSByZXNpemluZyBvciBkcmFnZ2luZy5cclxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xyXG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xyXG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGlzRHJhZ2dhYmxlRXZlbnQoZXZlbnQ6IE1vdXNlRXZlbnQpOiBib29sZWFuIHtcclxuICAgIGNvbnN0IHRhcmdldCA9IGV2ZW50VGFyZ2V0KGV2ZW50KTtcclxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KTtcclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgc3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpLnRyaW0oKSA9PT0gXCJkcmFnXCJcclxuICAgICAgICAmJiBldmVudC5vZmZzZXRYID49IDBcclxuICAgICAgICAmJiBldmVudC5vZmZzZXRYIDwgdGFyZ2V0LmNsaWVudFdpZHRoXHJcbiAgICAgICAgJiYgZXZlbnQub2Zmc2V0WSA+PSAwXHJcbiAgICAgICAgJiYgZXZlbnQub2Zmc2V0WSA8IHRhcmdldC5jbGllbnRIZWlnaHRcclxuICAgICk7XHJcbn1cclxuXHJcbi8vIFVzZSBjb25zdGFudHMgdG8gYXZvaWQgY29tcGFyaW5nIHN0cmluZ3MgbXVsdGlwbGUgdGltZXMuXHJcbmNvbnN0IE1vdXNlRG93biA9IDA7XHJcbmNvbnN0IE1vdXNlVXAgICA9IDE7XHJcbmNvbnN0IE1vdXNlTW92ZSA9IDI7XHJcblxyXG5mdW5jdGlvbiB1cGRhdGUoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcclxuICAgIC8vIFdpbmRvd3Mgc3VwcHJlc3NlcyBtb3VzZSBldmVudHMgYXQgdGhlIGVuZCBvZiBkcmFnZ2luZyBvciByZXNpemluZyxcclxuICAgIC8vIHNvIHdlIG5lZWQgdG8gYmUgc21hcnQgYW5kIHN5bnRoZXNpemUgYnV0dG9uIGV2ZW50cy5cclxuXHJcbiAgICBsZXQgZXZlbnRUeXBlOiBudW1iZXIsIGV2ZW50QnV0dG9ucyA9IGV2ZW50LmJ1dHRvbnM7XHJcbiAgICBzd2l0Y2ggKGV2ZW50LnR5cGUpIHtcclxuICAgICAgICBjYXNlICdtb3VzZWRvd24nOlxyXG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZURvd247XHJcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyB8ICgxIDw8IGV2ZW50LmJ1dHRvbik7IH1cclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgY2FzZSAnbW91c2V1cCc6XHJcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlVXA7XHJcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyAmIH4oMSA8PCBldmVudC5idXR0b24pOyB9XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIGRlZmF1bHQ6XHJcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlTW92ZTtcclxuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zOyB9XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgfVxyXG5cclxuICAgIGxldCByZWxlYXNlZCA9IGJ1dHRvbnMgJiB+ZXZlbnRCdXR0b25zO1xyXG4gICAgbGV0IHByZXNzZWQgPSBldmVudEJ1dHRvbnMgJiB+YnV0dG9ucztcclxuXHJcbiAgICBidXR0b25zID0gZXZlbnRCdXR0b25zO1xyXG5cclxuICAgIC8vIFN5bnRoZXNpemUgYSByZWxlYXNlLXByZXNzIHNlcXVlbmNlIGlmIHdlIGRldGVjdCBhIHByZXNzIG9mIGFuIGFscmVhZHkgcHJlc3NlZCBidXR0b24uXHJcbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZURvd24gJiYgIShwcmVzc2VkICYgZXZlbnQuYnV0dG9uKSkge1xyXG4gICAgICAgIHJlbGVhc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XHJcbiAgICAgICAgcHJlc3NlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFN1cHByZXNzIGFsbCBidXR0b24gZXZlbnRzIGR1cmluZyBkcmFnZ2luZyBhbmQgcmVzaXppbmcsXHJcbiAgICAvLyB1bmxlc3MgdGhpcyBpcyBhIG1vdXNldXAgZXZlbnQgdGhhdCBpcyBlbmRpbmcgYSBkcmFnIGFjdGlvbi5cclxuICAgIGlmIChcclxuICAgICAgICBldmVudFR5cGUgIT09IE1vdXNlTW92ZSAvLyBGYXN0IHBhdGggZm9yIG1vdXNlbW92ZVxyXG4gICAgICAgICYmIHJlc2l6aW5nXHJcbiAgICAgICAgfHwgKFxyXG4gICAgICAgICAgICBkcmFnZ2luZ1xyXG4gICAgICAgICAgICAmJiAoXHJcbiAgICAgICAgICAgICAgICBldmVudFR5cGUgPT09IE1vdXNlRG93blxyXG4gICAgICAgICAgICAgICAgfHwgZXZlbnQuYnV0dG9uICE9PSAwXHJcbiAgICAgICAgICAgIClcclxuICAgICAgICApXHJcbiAgICApIHtcclxuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIEhhbmRsZSByZWxlYXNlc1xyXG4gICAgaWYgKHJlbGVhc2VkICYgMSkgeyBwcmltYXJ5VXAoZXZlbnQpOyB9XHJcbiAgICAvLyBIYW5kbGUgcHJlc3Nlc1xyXG4gICAgaWYgKHByZXNzZWQgJiAxKSB7IHByaW1hcnlEb3duKGV2ZW50KTsgfVxyXG5cclxuICAgIC8vIEhhbmRsZSBtb3VzZW1vdmVcclxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlTW92ZSkgeyBvbk1vdXNlTW92ZShldmVudCk7IH07XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHByaW1hcnlEb3duKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XHJcbiAgICAvLyBSZXNldCByZWFkaW5lc3Mgc3RhdGUuXHJcbiAgICBjYW5EcmFnID0gZmFsc2U7XHJcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcclxuXHJcbiAgICAvLyBJZ25vcmUgcmVwZWF0ZWQgY2xpY2tzIG9uIG1hY09TIGFuZCBMaW51eC5cclxuICAgIGlmICghSXNXaW5kb3dzKCkpIHtcclxuICAgICAgICBpZiAoZXZlbnQudHlwZSA9PT0gJ21vdXNlZG93bicgJiYgZXZlbnQuYnV0dG9uID09PSAwICYmIGV2ZW50LmRldGFpbCAhPT0gMSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIGlmIChyZXNpemVFZGdlKSB7XHJcbiAgICAgICAgLy8gRG8gbm90IGFybSBlZGdlIHJlc2l6ZSBmcm9tIHN5bnRoZXNpemVkIHByZXNzZXMgb2JzZXJ2ZWQgb24gbW92ZS91cDpcclxuICAgICAgICAvLyBlbnRlcmluZyB0aGUgd2luZG93IHdpdGggdGhlIHByaW1hcnkgYnV0dG9uIGFscmVhZHkgaGVsZCBzaG91bGQgbm90XHJcbiAgICAgICAgLy8gc3RlYWwgYW5vdGhlciBnZXN0dXJlIGludG8gYSByZXNpemUuXHJcbiAgICAgICAgaWYgKGV2ZW50LnR5cGUgIT09ICdtb3VzZWRvd24nKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIFJlYWR5IHRvIHJlc2l6ZSBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lLlxyXG4gICAgICAgIGNhblJlc2l6ZSA9IHRydWU7XHJcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWcgb3BlcmF0aW9ucyB3aGVuIG9uIHJlc2l6ZSBlZGdlcy5cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gUmV0cmlldmUgdGFyZ2V0IGVsZW1lbnRcclxuICAgIC8vIFJlYWR5IHRvIGRyYWcgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZSBvbiBhIGRyYWdnYWJsZSBlbGVtZW50LlxyXG4gICAgLy8gSWdub3JlIGNsaWNrcyBvbiB0aGUgc2Nyb2xsYmFyLlxyXG4gICAgY2FuRHJhZyA9IGlzRHJhZ2dhYmxlRXZlbnQoZXZlbnQpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBwcmltYXJ5VXAoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcclxuICAgIC8vIFN0b3AgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLlxyXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xyXG4gICAgZHJhZ2dpbmcgPSBmYWxzZTtcclxuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xyXG4gICAgcmVzaXppbmcgPSBmYWxzZTtcclxufVxyXG5cclxuY29uc3QgY3Vyc29yRm9yRWRnZSA9IE9iamVjdC5mcmVlemUoe1xyXG4gICAgXCJzZS1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxyXG4gICAgXCJzdy1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxyXG4gICAgXCJudy1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxyXG4gICAgXCJuZS1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxyXG4gICAgXCJ3LXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxyXG4gICAgXCJuLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxyXG4gICAgXCJzLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxyXG4gICAgXCJlLXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxyXG59KVxyXG5cclxuZnVuY3Rpb24gc2V0UmVzaXplKGVkZ2U/OiBrZXlvZiB0eXBlb2YgY3Vyc29yRm9yRWRnZSk6IHZvaWQge1xyXG4gICAgaWYgKGVkZ2UpIHtcclxuICAgICAgICBpZiAoIXJlc2l6ZUVkZ2UpIHsgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yOyB9XHJcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBjdXJzb3JGb3JFZGdlW2VkZ2VdO1xyXG4gICAgfSBlbHNlIGlmICghZWRnZSAmJiByZXNpemVFZGdlKSB7XHJcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBkZWZhdWx0Q3Vyc29yO1xyXG4gICAgfVxyXG5cclxuICAgIHJlc2l6ZUVkZ2UgPSBlZGdlIHx8IFwiXCI7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XHJcbiAgICBpZiAoY2FuUmVzaXplICYmIHJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICAvLyBTdGFydCByZXNpemluZy5cclxuICAgICAgICByZXNpemluZyA9IHRydWU7XHJcbiAgICAgICAgaW52b2tlKFwid2FpbHM6cmVzaXplOlwiICsgcmVzaXplRWRnZSk7XHJcbiAgICB9IGVsc2UgaWYgKGNhbkRyYWcpIHtcclxuICAgICAgICAvLyBTdGFydCBkcmFnZ2luZy5cclxuICAgICAgICBkcmFnZ2luZyA9IHRydWU7XHJcbiAgICAgICAgaW52b2tlKFwid2FpbHM6ZHJhZ1wiKTtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcclxuICAgICAgICAvLyBFaXRoZXIgZHJhZyBvciByZXNpemUgaXMgb25nb2luZyxcclxuICAgICAgICAvLyByZXNldCByZWFkaW5lc3MgYW5kIHN0b3AgcHJvY2Vzc2luZy5cclxuICAgICAgICBjYW5EcmFnID0gY2FuUmVzaXplID0gZmFsc2U7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICghcmVzaXphYmxlIHx8ICghSXNXaW5kb3dzKCkgJiYgIShJc0xpbnV4KCkgJiYgR2V0RmxhZyhcImZyYW1lbGVzc1wiKSkpKSB7XHJcbiAgICAgICAgaWYgKHJlc2l6ZUVkZ2UpIHsgc2V0UmVzaXplKCk7IH1cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcclxuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xyXG5cclxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhcy5cclxuICAgIGNvbnN0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xyXG5cclxuICAgIC8vIFdoZW4gYSBzY3JvbGxiYXIgaXMgcHJlc2VudCBhdCB0aGUgd2luZG93IGVkZ2UgaXQgY29uc3VtZXMgbW91c2UgZXZlbnRzIGluIHRoYXQgc3RyaXAuXHJcbiAgICAvLyBTaGlmdCB0aGUgZWZmZWN0aXZlIGNvbnRlbnQgZWRnZSBpbndhcmQgc28gdGhlIHJlc2l6ZSB6b25lIHNpdHMganVzdCBiZWZvcmUgdGhlIHNjcm9sbGJhci5cclxuICAgIGNvbnN0IHNjcm9sbGJhcldpZHRoID0gTWF0aC5tYXgoMCwgd2luZG93LmlubmVyV2lkdGggLSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuY2xpZW50V2lkdGgpO1xyXG4gICAgY29uc3Qgc2Nyb2xsYmFySGVpZ2h0ID0gTWF0aC5tYXgoMCwgd2luZG93LmlubmVySGVpZ2h0IC0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LmNsaWVudEhlaWdodCk7XHJcbiAgICBjb25zdCByaWdodENvbnRlbnRFZGdlID0gd2luZG93LmlubmVyV2lkdGggLSBzY3JvbGxiYXJXaWR0aDtcclxuICAgIGNvbnN0IGJvdHRvbUNvbnRlbnRFZGdlID0gd2luZG93LmlubmVySGVpZ2h0IC0gc2Nyb2xsYmFySGVpZ2h0O1xyXG5cclxuICAgIGNvbnN0IHJpZ2h0Qm9yZGVyID0gZXZlbnQuY2xpZW50WCA8IHJpZ2h0Q29udGVudEVkZ2UgJiYgKHJpZ2h0Q29udGVudEVkZ2UgLSBldmVudC5jbGllbnRYKSA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgY29uc3QgbGVmdEJvcmRlciA9IGV2ZW50LmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGNvbnN0IHRvcEJvcmRlciA9IGV2ZW50LmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcbiAgICBjb25zdCBib3R0b21Cb3JkZXIgPSBldmVudC5jbGllbnRZIDwgYm90dG9tQ29udGVudEVkZ2UgJiYgKGJvdHRvbUNvbnRlbnRFZGdlIC0gZXZlbnQuY2xpZW50WSkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcblxyXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXIgYXJlYXMuXHJcbiAgICBjb25zdCByaWdodENvcm5lciA9IGV2ZW50LmNsaWVudFggPCByaWdodENvbnRlbnRFZGdlICYmIChyaWdodENvbnRlbnRFZGdlIC0gZXZlbnQuY2xpZW50WCkgPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBjb25zdCBsZWZ0Q29ybmVyID0gZXZlbnQuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGNvbnN0IHRvcENvcm5lciA9IGV2ZW50LmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG4gICAgY29uc3QgYm90dG9tQ29ybmVyID0gZXZlbnQuY2xpZW50WSA8IGJvdHRvbUNvbnRlbnRFZGdlICYmIChib3R0b21Db250ZW50RWRnZSAtIGV2ZW50LmNsaWVudFkpIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcclxuXHJcbiAgICBpZiAoIWxlZnRDb3JuZXIgJiYgIXRvcENvcm5lciAmJiAhYm90dG9tQ29ybmVyICYmICFyaWdodENvcm5lcikge1xyXG4gICAgICAgIC8vIE9wdGltaXNhdGlvbjogb3V0IG9mIGFsbCBjb3JuZXIgYXJlYXMgaW1wbGllcyBvdXQgb2YgYm9yZGVycy5cclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH1cclxuICAgIC8vIERldGVjdCBjb3JuZXJzLlxyXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgLy8gRGV0ZWN0IGJvcmRlcnMuXHJcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xyXG4gICAgLy8gT3V0IG9mIGJvcmRlciBhcmVhLlxyXG4gICAgZWxzZSBzZXRSZXNpemUoKTtcclxufVxyXG4iLCIvKlxyXG4gXyAgICAgX18gICAgXyBfX1xyXG58IHwgICAvIC9fX18oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgaW52b2tlIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XHJcbmltcG9ydCB7IHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XHJcbmltcG9ydCB7IGhhc0RPTSB9IGZyb20gXCIuL2Vudmlyb25tZW50LmpzXCI7XHJcblxyXG50eXBlIE5vbkNsaWVudFJlZ2lvbktpbmQgPSBcImNhcHRpb25cIiB8IFwibWluaW1pemVcIiB8IFwibWF4aW1pemVcIiB8IFwiY2xvc2VcIjtcclxuXHJcbmludGVyZmFjZSBOb25DbGllbnRSZWdpb24ge1xyXG4gICAga2luZDogTm9uQ2xpZW50UmVnaW9uS2luZDtcclxuICAgIGxlZnQ6IG51bWJlcjtcclxuICAgIHRvcDogbnVtYmVyO1xyXG4gICAgcmlnaHQ6IG51bWJlcjtcclxuICAgIGJvdHRvbTogbnVtYmVyO1xyXG59XHJcblxyXG4vKlxyXG4tLXdhaWxzLW5vbi1jbGllbnQtcmVnaW9uOiBjYXB0aW9uOyAgbWFya3MgYW4gYXJlYSB0aGF0IGNhbiBkcmFnIHRoZSB3aW5kb3dcclxuLS13YWlscy1ub24tY2xpZW50LXJlZ2lvbjogbWluaW1pemU7IG1hcmtzIGEgY3VzdG9tIG1pbmltaXplIGJ1dHRvblxyXG4tLXdhaWxzLW5vbi1jbGllbnQtcmVnaW9uOiBtYXhpbWl6ZTsgbWFya3MgYSBjdXN0b20gbWF4aW1pemUgYnV0dG9uXHJcbi0td2FpbHMtbm9uLWNsaWVudC1yZWdpb246IGNsb3NlOyAgICBtYXJrcyBhIGN1c3RvbSBjbG9zZSBidXR0b25cclxuKi9cclxuY29uc3QgcmVnaW9uUHJvcGVydHkgPSBcIi0td2FpbHMtbm9uLWNsaWVudC1yZWdpb25cIjtcclxuY29uc3QgcnVudGltZUNvbmZpZ1JlYWR5RXZlbnQgPSBcIndhaWxzOnJ1bnRpbWUtY29uZmlnLXJlYWR5XCI7XHJcbmNvbnN0IHZhbGlkUmVnaW9ucyA9IG5ldyBTZXQ8Tm9uQ2xpZW50UmVnaW9uS2luZD4oW1wiY2FwdGlvblwiLCBcIm1pbmltaXplXCIsIFwibWF4aW1pemVcIiwgXCJjbG9zZVwiXSk7XHJcblxyXG4vLyBTZXR1cFxyXG5pZiAoaGFzRE9NKSB7XHJcbiAgICB3aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxufVxyXG5cclxubGV0IHVwZGF0ZVBlbmRpbmcgPSBmYWxzZTtcclxubGV0IGxhc3RQYXlsb2FkID0gXCJcIjtcclxubGV0IG9ic2VydmVkRWxlbWVudHMgPSBuZXcgU2V0PEVsZW1lbnQ+KCk7XHJcbmxldCByZXNpemVPYnNlcnZlcjogUmVzaXplT2JzZXJ2ZXIgfCB1bmRlZmluZWQ7XHJcbmxldCB0cmFja2luZ1N0YXJ0ZWQgPSBmYWxzZTtcclxuXHJcbmZ1bmN0aW9uIG5vcm1hbGlzZVJlZ2lvbktpbmQodmFsdWU6IHN0cmluZyk6IE5vbkNsaWVudFJlZ2lvbktpbmQgfCB1bmRlZmluZWQge1xyXG4gICAgY29uc3QgcmVnaW9uID0gdmFsdWUudHJpbSgpLnRvTG93ZXJDYXNlKCk7XHJcbiAgICBpZiAodmFsaWRSZWdpb25zLmhhcyhyZWdpb24gYXMgTm9uQ2xpZW50UmVnaW9uS2luZCkpIHtcclxuICAgICAgICByZXR1cm4gcmVnaW9uIGFzIE5vbkNsaWVudFJlZ2lvbktpbmQ7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gdW5kZWZpbmVkO1xyXG59XHJcblxyXG5mdW5jdGlvbiBub25DbGllbnRSZWdpb25Gb3JFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQpOiBOb25DbGllbnRSZWdpb25LaW5kIHwgdW5kZWZpbmVkIHtcclxuICAgIGlmICghKGVsZW1lbnQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkpIHtcclxuICAgICAgICByZXR1cm4gdW5kZWZpbmVkO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XHJcbiAgICBjb25zdCByZWdpb24gPSBub3JtYWxpc2VSZWdpb25LaW5kKHN0eWxlLmdldFByb3BlcnR5VmFsdWUocmVnaW9uUHJvcGVydHkpKTtcclxuICAgIGlmICghcmVnaW9uKSB7XHJcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCBwYXJlbnQgPSBlbGVtZW50LnBhcmVudEVsZW1lbnQ7XHJcbiAgICBpZiAocGFyZW50KSB7XHJcbiAgICAgICAgY29uc3QgcGFyZW50U3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShwYXJlbnQpO1xyXG4gICAgICAgIC8vIFRoZSBDU1MgcHJvcGVydHkgaXMgaW5oZXJpdGVkLiBPbmx5IHJlcG9ydCB0aGUgb3V0ZXJtb3N0IGVsZW1lbnQgZm9yXHJcbiAgICAgICAgLy8gZWFjaCBjb250aWd1b3VzIHJlZ2lvbiBzbyBuYXRpdmUgaGl0IHRlc3Rpbmcgc2VlcyBzdGFibGUgcmVjdGFuZ2xlcy5cclxuICAgICAgICBpZiAobm9ybWFsaXNlUmVnaW9uS2luZChwYXJlbnRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKHJlZ2lvblByb3BlcnR5KSkgPT09IHJlZ2lvbikge1xyXG4gICAgICAgICAgICByZXR1cm4gdW5kZWZpbmVkO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gcmVnaW9uO1xyXG59XHJcblxyXG5mdW5jdGlvbiBpc1Zpc2libGUoZWxlbWVudDogSFRNTEVsZW1lbnQpOiBib29sZWFuIHtcclxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XHJcbiAgICByZXR1cm4gc3R5bGUuZGlzcGxheSAhPT0gXCJub25lXCIgJiZcclxuICAgICAgICBzdHlsZS52aXNpYmlsaXR5ICE9PSBcImhpZGRlblwiICYmXHJcbiAgICAgICAgc3R5bGUuY29udGVudFZpc2liaWxpdHkgIT09IFwiaGlkZGVuXCI7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGVsZW1lbnRSZWdpb24oZWxlbWVudDogRWxlbWVudCk6IE5vbkNsaWVudFJlZ2lvbiB8IHVuZGVmaW5lZCB7XHJcbiAgICBpZiAoIShlbGVtZW50IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpKSB7XHJcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCBraW5kID0gbm9uQ2xpZW50UmVnaW9uRm9yRWxlbWVudChlbGVtZW50KTtcclxuICAgIGlmICgha2luZCB8fCAhaXNWaXNpYmxlKGVsZW1lbnQpKSB7XHJcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCByZWN0ID0gZWxlbWVudC5nZXRCb3VuZGluZ0NsaWVudFJlY3QoKTtcclxuICAgIGlmIChyZWN0LndpZHRoIDw9IDAgfHwgcmVjdC5oZWlnaHQgPD0gMCkge1xyXG4gICAgICAgIHJldHVybiB1bmRlZmluZWQ7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gTmF0aXZlIGhpdCB0ZXN0aW5nIHJ1bnMgaW4gcGh5c2ljYWwgcGl4ZWxzLCB3aGlsZSBET00gZ2VvbWV0cnkgaXMgaW4gQ1NTIHBpeGVscy5cclxuICAgIGNvbnN0IHNjYWxlID0gd2luZG93LmRldmljZVBpeGVsUmF0aW8gfHwgMTtcclxuICAgIGNvbnN0IGxlZnQgPSBNYXRoLmZsb29yKHJlY3QubGVmdCAqIHNjYWxlKTtcclxuICAgIGNvbnN0IHRvcCA9IE1hdGguZmxvb3IocmVjdC50b3AgKiBzY2FsZSk7XHJcbiAgICBjb25zdCByaWdodCA9IE1hdGguY2VpbChyZWN0LnJpZ2h0ICogc2NhbGUpO1xyXG4gICAgY29uc3QgYm90dG9tID0gTWF0aC5jZWlsKHJlY3QuYm90dG9tICogc2NhbGUpO1xyXG5cclxuICAgIGlmIChyaWdodCA8PSBsZWZ0IHx8IGJvdHRvbSA8PSB0b3ApIHtcclxuICAgICAgICByZXR1cm4gdW5kZWZpbmVkO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiB7IGtpbmQsIGxlZnQsIHRvcCwgcmlnaHQsIGJvdHRvbSB9O1xyXG59XHJcblxyXG5mdW5jdGlvbiByZWdpb25FbGVtZW50cygpOiBFbGVtZW50W10ge1xyXG4gICAgY29uc3QgZWxlbWVudHM6IEVsZW1lbnRbXSA9IFtdO1xyXG5cclxuICAgIGlmIChkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQpIHtcclxuICAgICAgICBlbGVtZW50cy5wdXNoKGRvY3VtZW50LmRvY3VtZW50RWxlbWVudCk7XHJcbiAgICB9XHJcbiAgICBpZiAoZG9jdW1lbnQuYm9keSkge1xyXG4gICAgICAgIGVsZW1lbnRzLnB1c2goZG9jdW1lbnQuYm9keSk7XHJcbiAgICAgICAgLy8gQXBwZW5kIHZpYSBhIGxvb3A6IHNwcmVhZGluZyBhIGh1Z2UgTm9kZUxpc3QgaW50byBwdXNoKCkgb3ZlcmZsb3dzXHJcbiAgICAgICAgLy8gdGhlIGVuZ2luZSdzIGFyZ3VtZW50IGxpbWl0IG9uIHZlcnkgbGFyZ2UgZG9jdW1lbnRzLlxyXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoXCIqXCIpKSB7XHJcbiAgICAgICAgICAgIGVsZW1lbnRzLnB1c2goZWxlbWVudCk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBlbGVtZW50cztcclxufVxyXG5cclxuZnVuY3Rpb24gb2JzZXJ2ZVJlZ2lvbkVsZW1lbnRzKGVsZW1lbnRzOiBFbGVtZW50W10pOiB2b2lkIHtcclxuICAgIGlmICh0eXBlb2YgUmVzaXplT2JzZXJ2ZXIgPT09IFwidW5kZWZpbmVkXCIpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gVHJhY2sgc2l6ZSBjaGFuZ2VzIG9ubHkgZm9yIGFjdGl2ZSByZWdpb24gZWxlbWVudHMuIERPTSBzdHJ1Y3R1cmUgYW5kIHN0eWxlXHJcbiAgICAvLyBjaGFuZ2VzIGFyZSBjb3ZlcmVkIGJ5IE11dGF0aW9uT2JzZXJ2ZXIgaW4gc3RhcnROb25DbGllbnRSZWdpb25UcmFja2luZygpLlxyXG4gICAgcmVzaXplT2JzZXJ2ZXIgPz89IG5ldyBSZXNpemVPYnNlcnZlcihzY2hlZHVsZVVwZGF0ZSk7XHJcbiAgICBjb25zdCBuZXh0RWxlbWVudHMgPSBuZXcgU2V0KGVsZW1lbnRzKTtcclxuXHJcbiAgICBmb3IgKGNvbnN0IGVsZW1lbnQgb2Ygb2JzZXJ2ZWRFbGVtZW50cykge1xyXG4gICAgICAgIGlmICghbmV4dEVsZW1lbnRzLmhhcyhlbGVtZW50KSkge1xyXG4gICAgICAgICAgICByZXNpemVPYnNlcnZlci51bm9ic2VydmUoZWxlbWVudCk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBuZXh0RWxlbWVudHMpIHtcclxuICAgICAgICBpZiAoIW9ic2VydmVkRWxlbWVudHMuaGFzKGVsZW1lbnQpKSB7XHJcbiAgICAgICAgICAgIHJlc2l6ZU9ic2VydmVyLm9ic2VydmUoZWxlbWVudCk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIG9ic2VydmVkRWxlbWVudHMgPSBuZXh0RWxlbWVudHM7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHVwZGF0ZU5vbkNsaWVudFJlZ2lvbnMoKTogdm9pZCB7XHJcbiAgICB1cGRhdGVQZW5kaW5nID0gZmFsc2U7XHJcblxyXG4gICAgY29uc3QgZWxlbWVudHMgPSByZWdpb25FbGVtZW50cygpO1xyXG4gICAgY29uc3QgcmVnaW9uczogTm9uQ2xpZW50UmVnaW9uW10gPSBbXTtcclxuICAgIGNvbnN0IGFjdGl2ZUVsZW1lbnRzOiBFbGVtZW50W10gPSBbXTtcclxuXHJcbiAgICBmb3IgKGNvbnN0IGVsZW1lbnQgb2YgZWxlbWVudHMpIHtcclxuICAgICAgICBjb25zdCByZWdpb24gPSBlbGVtZW50UmVnaW9uKGVsZW1lbnQpO1xyXG4gICAgICAgIGlmIChyZWdpb24pIHtcclxuICAgICAgICAgICAgcmVnaW9ucy5wdXNoKHJlZ2lvbik7XHJcbiAgICAgICAgICAgIGFjdGl2ZUVsZW1lbnRzLnB1c2goZWxlbWVudCk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIG9ic2VydmVSZWdpb25FbGVtZW50cyhhY3RpdmVFbGVtZW50cyk7XHJcblxyXG4gICAgY29uc3QgcGF5bG9hZCA9IEpTT04uc3RyaW5naWZ5KHsgdmVyc2lvbjogMSwgcmVnaW9ucyB9KTtcclxuICAgIGlmIChwYXlsb2FkID09PSBsYXN0UGF5bG9hZCkge1xyXG4gICAgICAgIC8vIEF2b2lkIHNlbmRpbmcgZHVwbGljYXRlIG5hdGl2ZSBtZXNzYWdlcyBkdXJpbmcgcmVzaXplIG9yIHN0eWxlIGNodXJuLlxyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBsYXN0UGF5bG9hZCA9IHBheWxvYWQ7XHJcbiAgICBpbnZva2UoXCJ3YWlsczpub24tY2xpZW50LXJlZ2lvbjpcIiArIHBheWxvYWQpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBzY2hlZHVsZVVwZGF0ZSgpOiB2b2lkIHtcclxuICAgIGlmICh1cGRhdGVQZW5kaW5nKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIEJhdGNoIHJlZ2lvbiB1cGRhdGVzIHRvIGFuaW1hdGlvbiBmcmFtZXMgc28gbGF5b3V0IGlzIG1lYXN1cmVkIG9uY2UgcGVyIGZyYW1lLlxyXG4gICAgdXBkYXRlUGVuZGluZyA9IHRydWU7XHJcbiAgICB3aW5kb3cucmVxdWVzdEFuaW1hdGlvbkZyYW1lKHVwZGF0ZU5vbkNsaWVudFJlZ2lvbnMpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBzdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKCk6IHZvaWQge1xyXG4gICAgaWYgKHRyYWNraW5nU3RhcnRlZCkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICB0cmFja2luZ1N0YXJ0ZWQgPSB0cnVlO1xyXG4gICAgLy8gU2VuZCBhbiBpbml0aWFsIGVtcHR5IG9yIHBvcHVsYXRlZCByZWdpb24gbGlzdCBvbmNlIHRoZSBET00gaXMgcmVhZHkuXHJcbiAgICBzY2hlZHVsZVVwZGF0ZSgpO1xyXG5cclxuICAgIGNvbnN0IG11dGF0aW9uT2JzZXJ2ZXIgPSBuZXcgTXV0YXRpb25PYnNlcnZlcihzY2hlZHVsZVVwZGF0ZSk7XHJcbiAgICBtdXRhdGlvbk9ic2VydmVyLm9ic2VydmUoZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LCB7XHJcbiAgICAgICAgYXR0cmlidXRlczogdHJ1ZSxcclxuICAgICAgICBjaGlsZExpc3Q6IHRydWUsXHJcbiAgICAgICAgc3VidHJlZTogdHJ1ZSxcclxuICAgIH0pO1xyXG5cclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKFwicmVzaXplXCIsIHNjaGVkdWxlVXBkYXRlKTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKFwic2Nyb2xsXCIsIHNjaGVkdWxlVXBkYXRlLCB0cnVlKTtcclxuICAgIHdpbmRvdy52aXN1YWxWaWV3cG9ydD8uYWRkRXZlbnRMaXN0ZW5lcihcInJlc2l6ZVwiLCBzY2hlZHVsZVVwZGF0ZSk7XHJcbiAgICB3aW5kb3cudmlzdWFsVmlld3BvcnQ/LmFkZEV2ZW50TGlzdGVuZXIoXCJzY3JvbGxcIiwgc2NoZWR1bGVVcGRhdGUpO1xyXG59XHJcblxyXG5mdW5jdGlvbiB0cnlTdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKCk6IGJvb2xlYW4ge1xyXG4gICAgY29uc3Qgb3MgPSB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Py5PUztcclxuICAgIGlmIChvcyA9PT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IGVuYWJsZWQgPSB3aW5kb3cuX3dhaWxzLmZsYWdzPy5ub25DbGllbnRSZWdpb25UcmFja2luZztcclxuICAgIGlmIChvcyA9PT0gXCJ3aW5kb3dzXCIpIHtcclxuICAgICAgICBpZiAoZW5hYmxlZCA9PT0gdHJ1ZSkge1xyXG4gICAgICAgICAgICB3aGVuUmVhZHkoc3RhcnROb25DbGllbnRSZWdpb25UcmFja2luZyk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiB0cnVlO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiB0cnVlO1xyXG59XHJcblxyXG5pZiAoaGFzRE9NICYmICF0cnlTdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKCkpIHtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKHJ1bnRpbWVDb25maWdSZWFkeUV2ZW50LCB0cnlTdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nLCB7IG9uY2U6IHRydWUgfSk7XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcclxuXHJcbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xyXG5jb25zdCBTaG93TWV0aG9kID0gMTtcclxuY29uc3QgUXVpdE1ldGhvZCA9IDI7XHJcblxyXG4vKipcclxuICogSGlkZXMgYSBjZXJ0YWluIG1ldGhvZCBieSBjYWxsaW5nIHRoZSBIaWRlTWV0aG9kIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChIaWRlTWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcclxufVxyXG4iLCIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xyXG5cclxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xyXG5cclxuLyoqXHJcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMLlxyXG4gKlxyXG4gKiBAcGFyYW0gdXJsIC0gVGhlIFVSTCB0byBvcGVuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmw6IHN0cmluZyB8IFVSTCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmw6IHVybC50b1N0cmluZygpfSk7XHJcbn1cclxuIiwiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vaW5zcGVjdC1qcy9pcy1jYWxsYWJsZVxyXG5cclxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXHJcbi8vXHJcbi8vIENvcHlyaWdodCAoYykgMjAxNSBKb3JkYW4gSGFyYmFuZFxyXG4vL1xyXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5XHJcbi8vIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWxcclxuLy8gaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0c1xyXG4vLyB0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsXHJcbi8vIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpc1xyXG4vLyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxyXG4vL1xyXG4vLyBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcclxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cclxuLy9cclxuLy8gVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxyXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSxcclxuLy8gRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFXHJcbi8vIEFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVJcclxuLy8gTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSxcclxuLy8gT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEVcclxuLy8gU09GVFdBUkUuXHJcblxyXG52YXIgZm5Ub1N0ciA9IEZ1bmN0aW9uLnByb3RvdHlwZS50b1N0cmluZztcclxudmFyIHJlZmxlY3RBcHBseTogdHlwZW9mIFJlZmxlY3QuYXBwbHkgfCBmYWxzZSB8IG51bGwgPSB0eXBlb2YgUmVmbGVjdCA9PT0gJ29iamVjdCcgJiYgUmVmbGVjdCAhPT0gbnVsbCAmJiBSZWZsZWN0LmFwcGx5O1xyXG52YXIgYmFkQXJyYXlMaWtlOiBhbnk7XHJcbnZhciBpc0NhbGxhYmxlTWFya2VyOiBhbnk7XHJcbmlmICh0eXBlb2YgcmVmbGVjdEFwcGx5ID09PSAnZnVuY3Rpb24nICYmIHR5cGVvZiBPYmplY3QuZGVmaW5lUHJvcGVydHkgPT09ICdmdW5jdGlvbicpIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgYmFkQXJyYXlMaWtlID0gT2JqZWN0LmRlZmluZVByb3BlcnR5KHt9LCAnbGVuZ3RoJywge1xyXG4gICAgICAgICAgICBnZXQ6IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgICAgIHRocm93IGlzQ2FsbGFibGVNYXJrZXI7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgICAgICBpc0NhbGxhYmxlTWFya2VyID0ge307XHJcbiAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXRocm93LWxpdGVyYWxcclxuICAgICAgICByZWZsZWN0QXBwbHkoZnVuY3Rpb24gKCkgeyB0aHJvdyA0MjsgfSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcclxuICAgIH0gY2F0Y2ggKF8pIHtcclxuICAgICAgICBpZiAoXyAhPT0gaXNDYWxsYWJsZU1hcmtlcikge1xyXG4gICAgICAgICAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufSBlbHNlIHtcclxuICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XHJcbn1cclxuXHJcbnZhciBjb25zdHJ1Y3RvclJlZ2V4ID0gL15cXHMqY2xhc3NcXGIvO1xyXG52YXIgaXNFUzZDbGFzc0ZuID0gZnVuY3Rpb24gaXNFUzZDbGFzc0Z1bmN0aW9uKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgdmFyIGZuU3RyID0gZm5Ub1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICByZXR1cm4gY29uc3RydWN0b3JSZWdleC50ZXN0KGZuU3RyKTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7IC8vIG5vdCBhIGZ1bmN0aW9uXHJcbiAgICB9XHJcbn07XHJcblxyXG52YXIgdHJ5RnVuY3Rpb25PYmplY3QgPSBmdW5jdGlvbiB0cnlGdW5jdGlvblRvU3RyKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XHJcbiAgICAgICAgZm5Ub1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICByZXR1cm4gdHJ1ZTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcbn07XHJcbnZhciB0b1N0ciA9IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmc7XHJcbnZhciBvYmplY3RDbGFzcyA9ICdbb2JqZWN0IE9iamVjdF0nO1xyXG52YXIgZm5DbGFzcyA9ICdbb2JqZWN0IEZ1bmN0aW9uXSc7XHJcbnZhciBnZW5DbGFzcyA9ICdbb2JqZWN0IEdlbmVyYXRvckZ1bmN0aW9uXSc7XHJcbnZhciBkZGFDbGFzcyA9ICdbb2JqZWN0IEhUTUxBbGxDb2xsZWN0aW9uXSc7IC8vIElFIDExXHJcbnZhciBkZGFDbGFzczIgPSAnW29iamVjdCBIVE1MIGRvY3VtZW50LmFsbCBjbGFzc10nO1xyXG52YXIgZGRhQ2xhc3MzID0gJ1tvYmplY3QgSFRNTENvbGxlY3Rpb25dJzsgLy8gSUUgOS0xMFxyXG52YXIgaGFzVG9TdHJpbmdUYWcgPSB0eXBlb2YgU3ltYm9sID09PSAnZnVuY3Rpb24nICYmICEhU3ltYm9sLnRvU3RyaW5nVGFnOyAvLyBiZXR0ZXI6IHVzZSBgaGFzLXRvc3RyaW5ndGFnYFxyXG5cclxudmFyIGlzSUU2OCA9ICEoMCBpbiBbLF0pOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIG5vLXNwYXJzZS1hcnJheXMsIGNvbW1hLXNwYWNpbmdcclxuXHJcbnZhciBpc0REQTogKHZhbHVlOiBhbnkpID0+IGJvb2xlYW4gPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKCkgeyByZXR1cm4gZmFsc2U7IH07XHJcbmlmICh0eXBlb2YgZG9jdW1lbnQgPT09ICdvYmplY3QnKSB7XHJcbiAgICAvLyBGaXJlZm94IDMgY2Fub25pY2FsaXplcyBEREEgdG8gdW5kZWZpbmVkIHdoZW4gaXQncyBub3QgYWNjZXNzZWQgZGlyZWN0bHlcclxuICAgIHZhciBhbGwgPSBkb2N1bWVudC5hbGw7XHJcbiAgICBpZiAodG9TdHIuY2FsbChhbGwpID09PSB0b1N0ci5jYWxsKGRvY3VtZW50LmFsbCkpIHtcclxuICAgICAgICBpc0REQSA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwodmFsdWUpIHtcclxuICAgICAgICAgICAgLyogZ2xvYmFscyBkb2N1bWVudDogZmFsc2UgKi9cclxuICAgICAgICAgICAgLy8gaW4gSUUgNi04LCB0eXBlb2YgZG9jdW1lbnQuYWxsIGlzIFwib2JqZWN0XCIgYW5kIGl0J3MgdHJ1dGh5XHJcbiAgICAgICAgICAgIGlmICgoaXNJRTY4IHx8ICF2YWx1ZSkgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3VuZGVmaW5lZCcgfHwgdHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JykpIHtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgdmFyIHN0ciA9IHRvU3RyLmNhbGwodmFsdWUpO1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybiAoXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHN0ciA9PT0gZGRhQ2xhc3NcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczJcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczMgLy8gb3BlcmEgMTIuMTZcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBvYmplY3RDbGFzcyAvLyBJRSA2LThcclxuICAgICAgICAgICAgICAgICAgICApICYmIHZhbHVlKCcnKSA9PSBudWxsOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIGVxZXFlcVxyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZSkgeyAvKiovIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gaXNDYWxsYWJsZVJlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55ICB7XHJcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XHJcbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgdHJ5IHtcclxuICAgICAgICAocmVmbGVjdEFwcGx5IGFzIGFueSkodmFsdWUsIG51bGwsIGJhZEFycmF5TGlrZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgaWYgKGUgIT09IGlzQ2FsbGFibGVNYXJrZXIpIHsgcmV0dXJuIGZhbHNlOyB9XHJcbiAgICB9XHJcbiAgICByZXR1cm4gIWlzRVM2Q2xhc3NGbih2YWx1ZSkgJiYgdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBpc0NhbGxhYmxlTm9SZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSB7XHJcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XHJcbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKGhhc1RvU3RyaW5nVGFnKSB7IHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7IH1cclxuICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgdmFyIHN0ckNsYXNzID0gdG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICBpZiAoc3RyQ2xhc3MgIT09IGZuQ2xhc3MgJiYgc3RyQ2xhc3MgIT09IGdlbkNsYXNzICYmICEoL15cXFtvYmplY3QgSFRNTC8pLnRlc3Qoc3RyQ2xhc3MpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcclxufTtcclxuXHJcbmV4cG9ydCBkZWZhdWx0IHJlZmxlY3RBcHBseSA/IGlzQ2FsbGFibGVSZWZBcHBseSA6IGlzQ2FsbGFibGVOb1JlZkFwcGx5O1xyXG4iLCIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgaXNDYWxsYWJsZSBmcm9tIFwiLi9jYWxsYWJsZS5qc1wiO1xyXG5cclxuLyoqXHJcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdXNlZCBhcyByZWplY3Rpb24gcmVhc29uXHJcbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSBpcyBjYW5jZWxsZWQgc3VjY2Vzc2Z1bGx5LlxyXG4gKlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbEVycm9yXCJgLlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIGNhdXNlIHBhc3NlZCB0byB0aGUgY2FuY2VsIG1ldGhvZCwgaWYgYW55LlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIENhbmNlbEVycm9yIGV4dGVuZHMgRXJyb3Ige1xyXG4gICAgLyoqXHJcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxFcnJvcmAgaW5zdGFuY2UuXHJcbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxyXG4gICAgICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIHRvIGJlIGZvcndhcmRlZCB0byB0aGUgRXJyb3IgY29uc3RydWN0b3IuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcclxuICAgICAgICBzdXBlcihtZXNzYWdlLCBvcHRpb25zKTtcclxuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbEVycm9yXCI7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb25cclxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IHJlamVjdHMgYWZ0ZXIgYmVpbmcgY2FuY2VsbGVkLFxyXG4gKiBvciB3aGVuIHRoZSBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHRocm93cyBvciByZWplY3RzLlxyXG4gKlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCJgLlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIHJlYXNvbiB0aGUgcHJvbWlzZSByZWplY3RlZCB3aXRoLlxyXG4gKlxyXG4gKiBCZWNhdXNlIHRoZSBvcmlnaW5hbCBwcm9taXNlIHdhcyBjYW5jZWxsZWQsXHJcbiAqIGEgd3JhcHBlciBwcm9taXNlIHdpbGwgYmUgcGFzc2VkIHRvIHRoZSB1bmhhbmRsZWQgcmVqZWN0aW9uIGxpc3RlbmVyIGluc3RlYWQuXHJcbiAqIFRoZSB7QGxpbmsgcHJvbWlzZX0gcHJvcGVydHkgaG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIG9yaWdpbmFsIHByb21pc2UuXHJcbiAqL1xyXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IgZXh0ZW5kcyBFcnJvciB7XHJcbiAgICAvKipcclxuICAgICAqIEhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBwcm9taXNlIHRoYXQgd2FzIGNhbmNlbGxlZCBhbmQgdGhlbiByZWplY3RlZC5cclxuICAgICAqL1xyXG4gICAgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JgIGluc3RhbmNlLlxyXG4gICAgICogQHBhcmFtIHByb21pc2UgLSBUaGUgcHJvbWlzZSB0aGF0IGNhdXNlZCB0aGUgZXJyb3Igb3JpZ2luYWxseS5cclxuICAgICAqIEBwYXJhbSByZWFzb24gLSBUaGUgcmVqZWN0aW9uIHJlYXNvbi5cclxuICAgICAqIEBwYXJhbSBpbmZvIC0gQW4gb3B0aW9uYWwgaW5mb3JtYXRpdmUgbWVzc2FnZSBzcGVjaWZ5aW5nIHRoZSBjaXJjdW1zdGFuY2VzIGluIHdoaWNoIHRoZSBlcnJvciB3YXMgdGhyb3duLlxyXG4gICAgICogICAgICAgICAgICAgICBEZWZhdWx0cyB0byB0aGUgc3RyaW5nIGBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCJgLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3Rvcihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHJlYXNvbj86IGFueSwgaW5mbz86IHN0cmluZykge1xyXG4gICAgICAgIHN1cGVyKChpbmZvID8/IFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cIikgKyBcIiBSZWFzb246IFwiICsgZXJyb3JNZXNzYWdlKHJlYXNvbiksIHsgY2F1c2U6IHJlYXNvbiB9KTtcclxuICAgICAgICB0aGlzLnByb21pc2UgPSBwcm9taXNlO1xyXG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcIjtcclxuICAgIH1cclxufVxyXG5cclxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiA9ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcclxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciA9IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XHJcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyID0gKGNhdXNlPzogYW55KSA9PiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XHJcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4gPSAocmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4sIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IpID0+IHZvaWQ7XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xyXG4gICAgdGhlbjxUUmVzdWx0MSA9IFQsIFRSZXN1bHQyID0gbmV2ZXI+KG9uZnVsZmlsbGVkPzogKCh2YWx1ZTogVCkgPT4gVFJlc3VsdDEgfCBQcm9taXNlTGlrZTxUUmVzdWx0MT4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gVFJlc3VsdDIgfCBQcm9taXNlTGlrZTxUUmVzdWx0Mj4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQyPikgfCB1bmRlZmluZWQgfCBudWxsKTogQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MSB8IFRSZXN1bHQyPjtcclxuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxufVxyXG5cclxuLyoqXHJcbiAqIFdyYXBzIGEgY2FuY2VsbGFibGUgcHJvbWlzZSBhbG9uZyB3aXRoIGl0cyByZXNvbHV0aW9uIG1ldGhvZHMuXHJcbiAqIFRoZSBgb25jYW5jZWxsZWRgIGZpZWxkIHdpbGwgYmUgbnVsbCBpbml0aWFsbHkgYnV0IG1heSBiZSBzZXQgdG8gcHJvdmlkZSBhIGN1c3RvbSBjYW5jZWxsYXRpb24gZnVuY3Rpb24uXHJcbiAqL1xyXG5leHBvcnQgaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xyXG4gICAgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xyXG4gICAgcmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD47XHJcbiAgICByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yO1xyXG4gICAgb25jYW5jZWxsZWQ6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciB8IG51bGw7XHJcbn1cclxuXHJcbmludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSB7XHJcbiAgICByZWFkb25seSByb290OiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZTtcclxuICAgIHJlc29sdmluZzogYm9vbGVhbjtcclxuICAgIHNldHRsZWQ6IGJvb2xlYW47XHJcbiAgICByZWFzb24/OiBDYW5jZWxFcnJvcjtcclxufVxyXG5cclxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cclxuY29uc3QgYmFycmllclN5bSA9IFN5bWJvbChcImJhcnJpZXJcIik7XHJcbmNvbnN0IGNhbmNlbEltcGxTeW0gPSBTeW1ib2woXCJjYW5jZWxJbXBsXCIpO1xyXG5jb25zdCBzcGVjaWVzOiB0eXBlb2YgU3ltYm9sLnNwZWNpZXMgPSBTeW1ib2wuc3BlY2llcyA/PyBTeW1ib2woXCJzcGVjaWVzUG9seWZpbGxcIik7XHJcblxyXG4vKipcclxuICogQSBwcm9taXNlIHdpdGggYW4gYXR0YWNoZWQgbWV0aG9kIGZvciBjYW5jZWxsaW5nIGxvbmctcnVubmluZyBvcGVyYXRpb25zIChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWx9KS5cclxuICogQ2FuY2VsbGF0aW9uIGNhbiBvcHRpb25hbGx5IGJlIGJvdW5kIHRvIGFuIHtAbGluayBBYm9ydFNpZ25hbH1cclxuICogZm9yIGJldHRlciBjb21wb3NhYmlsaXR5IChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWxPbn0pLlxyXG4gKlxyXG4gKiBDYW5jZWxsaW5nIGEgcGVuZGluZyBwcm9taXNlIHdpbGwgcmVzdWx0IGluIGFuIGltbWVkaWF0ZSByZWplY3Rpb25cclxuICogd2l0aCBhbiBpbnN0YW5jZSBvZiB7QGxpbmsgQ2FuY2VsRXJyb3J9IGFzIHJlYXNvbixcclxuICogYnV0IHdob2V2ZXIgc3RhcnRlZCB0aGUgcHJvbWlzZSB3aWxsIGJlIHJlc3BvbnNpYmxlXHJcbiAqIGZvciBhY3R1YWxseSBhYm9ydGluZyB0aGUgdW5kZXJseWluZyBvcGVyYXRpb24uXHJcbiAqIFRvIHRoaXMgcHVycG9zZSwgdGhlIGNvbnN0cnVjdG9yIGFuZCBhbGwgY2hhaW5pbmcgbWV0aG9kc1xyXG4gKiBhY2NlcHQgb3B0aW9uYWwgY2FuY2VsbGF0aW9uIGNhbGxiYWNrcy5cclxuICpcclxuICogSWYgYSBgQ2FuY2VsbGFibGVQcm9taXNlYCBzdGlsbCByZXNvbHZlcyBhZnRlciBoYXZpbmcgYmVlbiBjYW5jZWxsZWQsXHJcbiAqIHRoZSByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWQuIElmIGl0IHJlamVjdHMsIHRoZSByZWFzb25cclxuICogd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLFxyXG4gKiB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBpbnN0YW5jZS5cclxuICogVG8gZmFjaWxpdGF0ZSB0aGUgaGFuZGxpbmcgb2YgY2FuY2VsbGF0aW9uIHJlcXVlc3RzLFxyXG4gKiBjYW5jZWxsZWQgYENhbmNlbGxhYmxlUHJvbWlzZWBzIHdpbGwgX25vdF8gcmVwb3J0IHVuaGFuZGxlZCBgQ2FuY2VsRXJyb3Jgc1xyXG4gKiB3aG9zZSBgY2F1c2VgIGZpZWxkIGlzIHRoZSBzYW1lIGFzIHRoZSBvbmUgd2l0aCB3aGljaCB0aGUgY3VycmVudCBwcm9taXNlIHdhcyBjYW5jZWxsZWQuXHJcbiAqXHJcbiAqIEFsbCB1c3VhbCBwcm9taXNlIG1ldGhvZHMgYXJlIGRlZmluZWQgYW5kIHJldHVybiBhIGBDYW5jZWxsYWJsZVByb21pc2VgXHJcbiAqIHdob3NlIGNhbmNlbCBtZXRob2Qgd2lsbCBjYW5jZWwgdGhlIHBhcmVudCBvcGVyYXRpb24gYXMgd2VsbCwgcHJvcGFnYXRpbmcgdGhlIGNhbmNlbGxhdGlvbiByZWFzb25cclxuICogdXB3YXJkcyB0aHJvdWdoIHByb21pc2UgY2hhaW5zLlxyXG4gKiBDb252ZXJzZWx5LCBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aWxsIG5vdCBhdXRvbWF0aWNhbGx5IGNhbmNlbCBkZXBlbmRlbnQgcHJvbWlzZXMgZG93bnN0cmVhbTpcclxuICogYGBgdHNcclxuICogbGV0IHJvb3QgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xyXG4gKiBsZXQgY2hpbGQxID0gcm9vdC50aGVuKCgpID0+IHsgLi4uIH0pO1xyXG4gKiBsZXQgY2hpbGQyID0gY2hpbGQxLnRoZW4oKCkgPT4geyAuLi4gfSk7XHJcbiAqIGxldCBjaGlsZDMgPSByb290LmNhdGNoKCgpID0+IHsgLi4uIH0pO1xyXG4gKiBjaGlsZDEuY2FuY2VsKCk7IC8vIENhbmNlbHMgY2hpbGQxIGFuZCByb290LCBidXQgbm90IGNoaWxkMiBvciBjaGlsZDNcclxuICogYGBgXHJcbiAqIENhbmNlbGxpbmcgYSBwcm9taXNlIHRoYXQgaGFzIGFscmVhZHkgc2V0dGxlZCBpcyBzYWZlIGFuZCBoYXMgbm8gY29uc2VxdWVuY2UuXHJcbiAqXHJcbiAqIFRoZSBgY2FuY2VsYCBtZXRob2QgcmV0dXJucyBhIHByb21pc2UgdGhhdCBfYWx3YXlzIGZ1bGZpbGxzX1xyXG4gKiBhZnRlciB0aGUgd2hvbGUgY2hhaW4gaGFzIHByb2Nlc3NlZCB0aGUgY2FuY2VsIHJlcXVlc3RcclxuICogYW5kIGFsbCBhdHRhY2hlZCBjYWxsYmFja3MgdXAgdG8gdGhhdCBtb21lbnQgaGF2ZSBydW4uXHJcbiAqXHJcbiAqIEFsbCBFUzIwMjQgcHJvbWlzZSBtZXRob2RzIChzdGF0aWMgYW5kIGluc3RhbmNlKSBhcmUgZGVmaW5lZCBvbiBDYW5jZWxsYWJsZVByb21pc2UsXHJcbiAqIGJ1dCBhY3R1YWwgYXZhaWxhYmlsaXR5IG1heSB2YXJ5IHdpdGggT1Mvd2VidmlldyB2ZXJzaW9uLlxyXG4gKlxyXG4gKiBJbiBsaW5lIHdpdGggdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmcsXHJcbiAqIGBDYW5jZWxsYWJsZVByb21pc2VgIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuXHJcbiAqIEV4dGVuZGVycyBzaG91bGQgdGFrZSBjYXJlIHRvIHByb3ZpZGUgdGhlaXIgb3duIG1ldGhvZCBpbXBsZW1lbnRhdGlvbnMuXHJcbiAqIFRoaXMgbWlnaHQgYmUgcmVjb25zaWRlcmVkIGluIGNhc2UgdGhlIHByb3Bvc2FsIGlzIHJldGlyZWQuXHJcbiAqXHJcbiAqIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBhIHdyYXBwZXIgYXJvdW5kIHRoZSBET00gUHJvbWlzZSBvYmplY3RcclxuICogYW5kIGlzIGNvbXBsaWFudCB3aXRoIHRoZSBbUHJvbWlzZXMvQSsgc3BlY2lmaWNhdGlvbl0oaHR0cHM6Ly9wcm9taXNlc2FwbHVzLmNvbS8pXHJcbiAqIChpdCBwYXNzZXMgdGhlIFtjb21wbGlhbmNlIHN1aXRlXShodHRwczovL2dpdGh1Yi5jb20vcHJvbWlzZXMtYXBsdXMvcHJvbWlzZXMtdGVzdHMpKVxyXG4gKiBpZiBzbyBpcyB0aGUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbi5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBDYW5jZWxsYWJsZVByb21pc2U8VD4gZXh0ZW5kcyBQcm9taXNlPFQ+IGltcGxlbWVudHMgUHJvbWlzZUxpa2U8VD4sIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xyXG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXHJcbiAgICAvKiogQGludGVybmFsICovXHJcbiAgICBwcml2YXRlIFtiYXJyaWVyU3ltXSE6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+IHwgbnVsbDtcclxuICAgIC8qKiBAaW50ZXJuYWwgKi9cclxuICAgIHByaXZhdGUgcmVhZG9ubHkgW2NhbmNlbEltcGxTeW1dITogKHJlYXNvbjogQ2FuY2VsRXJyb3IpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgYENhbmNlbGxhYmxlUHJvbWlzZWAuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGV4ZWN1dG9yIC0gQSBjYWxsYmFjayB1c2VkIHRvIGluaXRpYWxpemUgdGhlIHByb21pc2UuIFRoaXMgY2FsbGJhY2sgaXMgcGFzc2VkIHR3byBhcmd1bWVudHM6XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhIGByZXNvbHZlYCBjYWxsYmFjayB1c2VkIHRvIHJlc29sdmUgdGhlIHByb21pc2Ugd2l0aCBhIHZhbHVlXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBvciB0aGUgcmVzdWx0IG9mIGFub3RoZXIgcHJvbWlzZSAocG9zc2libHkgY2FuY2VsbGFibGUpLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGEgYHJlamVjdGAgY2FsbGJhY2sgdXNlZCB0byByZWplY3QgdGhlIHByb21pc2Ugd2l0aCBhIHByb3ZpZGVkIHJlYXNvbiBvciBlcnJvci5cclxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIHRoZSB2YWx1ZSBwcm92aWRlZCB0byB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAoaXQgaGFzIGEgYHRoZW5gIF9hbmRfIGEgYGNhbmNlbGAgbWV0aG9kKSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGZvcndhcmRlZCB0byB0aGF0IG9iamVjdCBhbmQgdGhlIG9uY2FuY2VsbGVkIHdpbGwgbm90IGJlIGludm9rZWQgYW55bW9yZS5cclxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIGFueSBvbmUgb2YgdGhlIHR3byBjYWxsYmFja3MgaXMgY2FsbGVkIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhlIHByb3ZpZGVkIHZhbHVlcyB3aWxsIGJlIGNhbmNlbGxlZCBhbmQgcmVzb2x2ZWQgYXMgdXN1YWwsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBidXQgdGhlaXIgcmVzdWx0cyB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqICAgICAgICAgICAgICAgICAgIEhvd2V2ZXIsIGlmIHRoZSByZXNvbHV0aW9uIHByb2Nlc3MgdWx0aW1hdGVseSBlbmRzIHVwIGluIGEgcmVqZWN0aW9uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGF0IGlzIG5vdCBkdWUgdG8gY2FuY2VsbGF0aW9uLCB0aGUgcmVqZWN0aW9uIHJlYXNvblxyXG4gICAgICogICAgICAgICAgICAgICAgICAgd2lsbCBiZSB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfVxyXG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqIEBwYXJhbSBvbmNhbmNlbGxlZCAtIEl0IGlzIHRoZSBjYWxsZXIncyByZXNwb25zaWJpbGl0eSB0byBlbnN1cmUgdGhhdCBhbnkgb3BlcmF0aW9uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBzdGFydGVkIGJ5IHRoZSBleGVjdXRvciBpcyBwcm9wZXJseSBoYWx0ZWQgdXBvbiBjYW5jZWxsYXRpb24uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBUaGlzIG9wdGlvbmFsIGNhbGxiYWNrIGNhbiBiZSB1c2VkIHRvIHRoYXQgcHVycG9zZS5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIEl0IHdpbGwgYmUgY2FsbGVkIF9zeW5jaHJvbm91c2x5XyB3aXRoIGEgY2FuY2VsbGF0aW9uIGNhdXNlXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aGVuIGNhbmNlbGxhdGlvbiBpcyByZXF1ZXN0ZWQsIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgcmVqZWN0ZWRcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9LCBidXQgX2JlZm9yZV9cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFueSB7QGxpbmsgdGhlbn0ve0BsaW5rIGNhdGNofS97QGxpbmsgZmluYWxseX0gY2FsbGJhY2sgcnVucy5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH1cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpbGwgb25seSBmdWxmaWxsIGFmdGVyIHRoZSBmb3JtZXIgaGFzIHNldHRsZWQuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBVbmhhbmRsZWQgZXhjZXB0aW9ucyBvciByZWplY3Rpb25zIGZyb20gdGhlIGNhbGxiYWNrIHdpbGwgYmUgd3JhcHBlZFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIHVuaGFuZGxlZCByZWplY3Rpb25zLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBjYWxsZWQgYmVmb3JlIGNhbmNlbGxhdGlvbiB3aXRoIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGlzIHByb21pc2Ugd2lsbCBiZSBkaXZlcnRlZCB0byB0aGF0IHByb21pc2UsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbmQgdGhlIG9yaWdpbmFsIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKGV4ZWN1dG9yOiBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpIHtcclxuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xyXG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xyXG4gICAgICAgIHN1cGVyKChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xyXG5cclxuICAgICAgICBpZiAoKHRoaXMuY29uc3RydWN0b3IgYXMgYW55KVtzcGVjaWVzXSAhPT0gUHJvbWlzZSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuIFBsZWFzZSByZWZyYWluIGZyb20gb3ZlcnJpZGluZyB0aGUgW1N5bWJvbC5zcGVjaWVzXSBzdGF0aWMgcHJvcGVydHkuXCIpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgbGV0IHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7XHJcbiAgICAgICAgICAgIHByb21pc2U6IHRoaXMsXHJcbiAgICAgICAgICAgIHJlc29sdmUsXHJcbiAgICAgICAgICAgIHJlamVjdCxcclxuICAgICAgICAgICAgZ2V0IG9uY2FuY2VsbGVkKCkgeyByZXR1cm4gb25jYW5jZWxsZWQgPz8gbnVsbDsgfSxcclxuICAgICAgICAgICAgc2V0IG9uY2FuY2VsbGVkKGNiKSB7IG9uY2FuY2VsbGVkID0gY2IgPz8gdW5kZWZpbmVkOyB9XHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgY29uc3Qgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xyXG4gICAgICAgICAgICBnZXQgcm9vdCgpIHsgcmV0dXJuIHN0YXRlOyB9LFxyXG4gICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxyXG4gICAgICAgICAgICBzZXR0bGVkOiBmYWxzZVxyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFNldHVwIGNhbmNlbGxhdGlvbiBzeXN0ZW0uXHJcbiAgICAgICAgdm9pZCBPYmplY3QuZGVmaW5lUHJvcGVydGllcyh0aGlzLCB7XHJcbiAgICAgICAgICAgIFtiYXJyaWVyU3ltXToge1xyXG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IHRydWUsXHJcbiAgICAgICAgICAgICAgICB2YWx1ZTogbnVsbFxyXG4gICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICBbY2FuY2VsSW1wbFN5bV06IHtcclxuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIHZhbHVlOiBjYW5jZWxsZXJGb3IocHJvbWlzZSwgc3RhdGUpXHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuXHJcbiAgICAgICAgLy8gUnVuIHRoZSBhY3R1YWwgZXhlY3V0b3IuXHJcbiAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBzdGF0ZSk7XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgZXhlY3V0b3IocmVzb2x2ZXJGb3IocHJvbWlzZSwgc3RhdGUpLCByZWplY3Rvcik7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHtcclxuICAgICAgICAgICAgICAgIGNvbnNvbGUubG9nKFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBDYW5jZWxsYWJsZVByb21pc2UgZXhlY3V0b3IuXCIsIGVycik7XHJcbiAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ2FuY2VscyBpbW1lZGlhdGVseSB0aGUgZXhlY3V0aW9uIG9mIHRoZSBvcGVyYXRpb24gYXNzb2NpYXRlZCB3aXRoIHRoaXMgcHJvbWlzZS5cclxuICAgICAqIFRoZSBwcm9taXNlIHJlamVjdHMgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0gaW5zdGFuY2UgYXMgcmVhc29uLFxyXG4gICAgICogd2l0aCB0aGUge0BsaW5rIENhbmNlbEVycm9yI2NhdXNlfSBwcm9wZXJ0eSBzZXQgdG8gdGhlIGdpdmVuIGFyZ3VtZW50LCBpZiBhbnkuXHJcbiAgICAgKlxyXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgYWZ0ZXIgdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZDtcclxuICAgICAqIHJlcGVhdGVkIGNhbGxzIGluIHBhcnRpY3VsYXIgYXJlIHNhZmUsIGJ1dCBvbmx5IHRoZSBmaXJzdCBvbmVcclxuICAgICAqIHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIGBDYW5jZWxFcnJvcmAgZXhjZXB0aW9uIF9uZWVkIG5vdF8gYmUgaGFuZGxlZCBleHBsaWNpdGx5IF9vbiB0aGUgcHJvbWlzZXMgdGhhdCBhcmUgYmVpbmcgY2FuY2VsbGVkOl9cclxuICAgICAqIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpdGggbm8gYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXIgZG9lcyBub3QgdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LlxyXG4gICAgICogVGhlcmVmb3JlLCB0aGUgZm9sbG93aW5nIGlkaW9tcyBhcmUgYWxsIGVxdWFsbHkgY29ycmVjdDpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLmNhbmNlbCgpO1xyXG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2FuY2VsKCk7XHJcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYXRjaCguLi4pLmNhbmNlbCgpO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBXaGVuZXZlciBzb21lIGNhbmNlbGxlZCBwcm9taXNlIGluIGEgY2hhaW4gcmVqZWN0cyB3aXRoIGEgYENhbmNlbEVycm9yYFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBjYW5jZWxsYXRpb24gY2F1c2UgYXMgaXRzZWxmLCB0aGUgZXJyb3Igd2lsbCBiZSBkaXNjYXJkZWQgc2lsZW50bHkuXHJcbiAgICAgKiBIb3dldmVyLCB0aGUgYENhbmNlbEVycm9yYCBfd2lsbCBzdGlsbCBiZSBkZWxpdmVyZWRfIHRvIGFsbCBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlcnNcclxuICAgICAqIGFkZGVkIGJ5IHtAbGluayB0aGVufSBhbmQgcmVsYXRlZCBtZXRob2RzOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XHJcbiAgICAgKiBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLmNhdGNoKGNvbnNvbGUubG9nKTtcclxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBBIENhbmNlbEVycm9yIGlzIHByaW50ZWQgdG8gdGhlIGNvbnNvbGUuXHJcbiAgICAgKiBgYGBcclxuICAgICAqIElmIHRoZSBgQ2FuY2VsRXJyb3JgIGlzIG5vdCBoYW5kbGVkIGRvd25zdHJlYW0gYnkgdGhlIHRpbWUgaXQgcmVhY2hlc1xyXG4gICAgICogYSBfbm9uLWNhbmNlbGxlZF8gcHJvbWlzZSwgaXQgX3dpbGxfIHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCxcclxuICAgICAqIGp1c3QgbGlrZSBub3JtYWwgcmVqZWN0aW9ucyB3b3VsZDpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xyXG4gICAgICogbGV0IGNoYWluZWQgPSBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLnRoZW4oKCkgPT4geyAuLi4gfSk7IC8vIE5vIGNhdGNoLi4uXHJcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gVW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCBvbiBjaGFpbmVkIVxyXG4gICAgICogYGBgXHJcbiAgICAgKiBUaGVyZWZvcmUsIGl0IGlzIGltcG9ydGFudCB0byBlaXRoZXIgY2FuY2VsIHdob2xlIHByb21pc2UgY2hhaW5zIGZyb20gdGhlaXIgdGFpbCxcclxuICAgICAqIGFzIHNob3duIGluIHRoZSBjb3JyZWN0IGlkaW9tcyBhYm92ZSwgb3IgdGFrZSBjYXJlIG9mIGhhbmRsaW5nIGVycm9ycyBldmVyeXdoZXJlLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIEEgY2FuY2VsbGFibGUgcHJvbWlzZSB0aGF0IF9mdWxmaWxsc18gYWZ0ZXIgdGhlIGNhbmNlbCBjYWxsYmFjayAoaWYgYW55KVxyXG4gICAgICogYW5kIGFsbCBoYW5kbGVycyBhdHRhY2hlZCB1cCB0byB0aGUgY2FsbCB0byBjYW5jZWwgaGF2ZSBydW4uXHJcbiAgICAgKiBJZiB0aGUgY2FuY2VsIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgYnkgYGNhbmNlbGBcclxuICAgICAqIHdpbGwgYWxzbyB3YWl0IGZvciB0aGF0IHRoZW5hYmxlIHRvIHNldHRsZS5cclxuICAgICAqIFRoaXMgZW5hYmxlcyBjYWxsZXJzIHRvIHdhaXQgZm9yIHRoZSBjYW5jZWxsZWQgb3BlcmF0aW9uIHRvIHRlcm1pbmF0ZVxyXG4gICAgICogd2l0aG91dCBiZWluZyBmb3JjZWQgdG8gaGFuZGxlIHBvdGVudGlhbCBlcnJvcnMgYXQgdGhlIGNhbGwgc2l0ZS5cclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKS50aGVuKCgpID0+IHtcclxuICAgICAqICAgICAvLyBDbGVhbnVwIGZpbmlzaGVkLCBpdCdzIHNhZmUgdG8gZG8gc29tZXRoaW5nIGVsc2UuXHJcbiAgICAgKiB9LCAoZXJyKSA9PiB7XHJcbiAgICAgKiAgICAgLy8gVW5yZWFjaGFibGU6IHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gY2FuY2VsIHdpbGwgbmV2ZXIgcmVqZWN0LlxyXG4gICAgICogfSk7XHJcbiAgICAgKiBgYGBcclxuICAgICAqIE5vdGUgdGhhdCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIF9ub3RfIGhhbmRsZSBpbXBsaWNpdGx5IGFueSByZWplY3Rpb25cclxuICAgICAqIHRoYXQgbWlnaHQgaGF2ZSBvY2N1cnJlZCBhbHJlYWR5IGluIHRoZSBjYW5jZWxsZWQgY2hhaW4uXHJcbiAgICAgKiBJdCB3aWxsIGp1c3QgdHJhY2sgd2hldGhlciByZWdpc3RlcmVkIGhhbmRsZXJzIGhhdmUgYmVlbiBleGVjdXRlZCBvciBub3QuXHJcbiAgICAgKiBUaGVyZWZvcmUsIHVuaGFuZGxlZCByZWplY3Rpb25zIHdpbGwgbmV2ZXIgYmUgc2lsZW50bHkgaGFuZGxlZCBieSBjYWxsaW5nIGNhbmNlbC5cclxuICAgICAqL1xyXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xyXG4gICAgICAgICAgICAvLyBJTlZBUklBTlQ6IHRoZSByZXN1bHQgb2YgdGhpc1tjYW5jZWxJbXBsU3ltXSBhbmQgdGhlIGJhcnJpZXIgZG8gbm90IGV2ZXIgcmVqZWN0LlxyXG4gICAgICAgICAgICAvLyBVbmZvcnR1bmF0ZWx5IG1hY09TIEhpZ2ggU2llcnJhIGRvZXMgbm90IHN1cHBvcnQgUHJvbWlzZS5hbGxTZXR0bGVkLlxyXG4gICAgICAgICAgICBQcm9taXNlLmFsbChbXHJcbiAgICAgICAgICAgICAgICB0aGlzW2NhbmNlbEltcGxTeW1dKG5ldyBDYW5jZWxFcnJvcihcIlByb21pc2UgY2FuY2VsbGVkLlwiLCB7IGNhdXNlIH0pKSxcclxuICAgICAgICAgICAgICAgIGN1cnJlbnRCYXJyaWVyKHRoaXMpXHJcbiAgICAgICAgICAgIF0pLnRoZW4oKCkgPT4gcmVzb2x2ZSgpLCAoKSA9PiByZXNvbHZlKCkpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQmluZHMgcHJvbWlzZSBjYW5jZWxsYXRpb24gdG8gdGhlIGFib3J0IGV2ZW50IG9mIHRoZSBnaXZlbiB7QGxpbmsgQWJvcnRTaWduYWx9LlxyXG4gICAgICogSWYgdGhlIHNpZ25hbCBoYXMgYWxyZWFkeSBhYm9ydGVkLCB0aGUgcHJvbWlzZSB3aWxsIGJlIGNhbmNlbGxlZCBpbW1lZGlhdGVseS5cclxuICAgICAqIFdoZW4gZWl0aGVyIGNvbmRpdGlvbiBpcyB2ZXJpZmllZCwgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZSB3aWxsIGJlIHNldFxyXG4gICAgICogdG8gdGhlIHNpZ25hbCdzIGFib3J0IHJlYXNvbiAoc2VlIHtAbGluayBBYm9ydFNpZ25hbCNyZWFzb259KS5cclxuICAgICAqXHJcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCAob3IgaWYgdGhlIHNpZ25hbCBhYm9ydHMpIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZC5cclxuICAgICAqIE9ubHkgdGhlIGZpcnN0IHNpZ25hbCB0byBhYm9ydCB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEZvciBtb3JlIGRldGFpbHMgYWJvdXQgdGhlIGNhbmNlbGxhdGlvbiBwcm9jZXNzLFxyXG4gICAgICogc2VlIHtAbGluayBjYW5jZWx9IGFuZCB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAgY29uc3RydWN0b3IuXHJcbiAgICAgKlxyXG4gICAgICogVGhpcyBtZXRob2QgZW5hYmxlcyBgYXdhaXRgaW5nIGNhbmNlbGxhYmxlIHByb21pc2VzIHdpdGhvdXQgaGF2aW5nXHJcbiAgICAgKiB0byBzdG9yZSB0aGVtIGZvciBmdXR1cmUgY2FuY2VsbGF0aW9uLCBlLmcuOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGF3YWl0IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCkuY2FuY2VsT24oc2lnbmFsKTtcclxuICAgICAqIGBgYFxyXG4gICAgICogaW5zdGVhZCBvZjpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBsZXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQgPSBsb25nUnVubmluZ09wZXJhdGlvbigpO1xyXG4gICAgICogYXdhaXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQ7XHJcbiAgICAgKiBgYGBcclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGlzIHByb21pc2UsIGZvciBtZXRob2QgY2hhaW5pbmcuXHJcbiAgICAgKi9cclxuICAgIGNhbmNlbE9uKHNpZ25hbDogQWJvcnRTaWduYWwpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIGlmIChzaWduYWwuYWJvcnRlZCkge1xyXG4gICAgICAgICAgICB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pXHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgc2lnbmFsLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKSwge2NhcHR1cmU6IHRydWV9KTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIHJldHVybiB0aGlzO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQXR0YWNoZXMgY2FsbGJhY2tzIGZvciB0aGUgcmVzb2x1dGlvbiBhbmQvb3IgcmVqZWN0aW9uIG9mIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYC5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXHJcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXHJcbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcclxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXHJcbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gb25mdWxmaWxsZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZXNvbHZlZC5cclxuICAgICAqIEBwYXJhbSBvbnJlamVjdGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVqZWN0ZWQuXHJcbiAgICAgKiBAcmV0dXJucyBBIGBDYW5jZWxsYWJsZVByb21pc2VgIGZvciB0aGUgY29tcGxldGlvbiBvZiB3aGljaGV2ZXIgY2FsbGJhY2sgaXMgZXhlY3V0ZWQuXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBob29rZWQgdXAgdG8gcHJvcGFnYXRlIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB1cCB0aGUgY2hhaW4sIGJ1dCBub3QgZG93bjpcclxuICAgICAqXHJcbiAgICAgKiAgIC0gaWYgdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgYmUgaW52b2tlZCB3aXRoIGEgYENhbmNlbEVycm9yYFxyXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSBfd2lsbCByZXNvbHZlIHJlZ3VsYXJseV8gd2l0aCBpdHMgcmVzdWx0O1xyXG4gICAgICogICAtIGNvbnZlcnNlbHksIGlmIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgX3RoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQgdG9vO19cclxuICAgICAqICAgICB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBzdGlsbCBiZSBpbnZva2VkIHdpdGggdGhlIHBhcmVudCdzIGBDYW5jZWxFcnJvcmAsXHJcbiAgICAgKiAgICAgYnV0IGl0cyByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWRcclxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhIGBDYW5jZWxFcnJvcmAgYXMgd2VsbC5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xyXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXHJcbiAgICAgKlxyXG4gICAgICogSWYgZWl0aGVyIGNhbGxiYWNrIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXHJcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICovXHJcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+IHtcclxuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS50aGVuIGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBOT1RFOiBUeXBlU2NyaXB0J3MgYnVpbHQtaW4gdHlwZSBmb3IgdGhlbiBpcyBicm9rZW4sXHJcbiAgICAgICAgLy8gYXMgaXQgYWxsb3dzIHNwZWNpZnlpbmcgYW4gYXJiaXRyYXJ5IFRSZXN1bHQxICE9IFQgZXZlbiB3aGVuIG9uZnVsZmlsbGVkIGlzIG5vdCBhIGZ1bmN0aW9uLlxyXG4gICAgICAgIC8vIFdlIGNhbm5vdCBmaXggaXQgaWYgd2Ugd2FudCB0byBDYW5jZWxsYWJsZVByb21pc2UgdG8gaW1wbGVtZW50IFByb21pc2VMaWtlPFQ+LlxyXG5cclxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25mdWxmaWxsZWQpKSB7IG9uZnVsZmlsbGVkID0gaWRlbnRpdHkgYXMgYW55OyB9XHJcbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9ucmVqZWN0ZWQpKSB7IG9ucmVqZWN0ZWQgPSB0aHJvd2VyOyB9XHJcblxyXG4gICAgICAgIGlmIChvbmZ1bGZpbGxlZCA9PT0gaWRlbnRpdHkgJiYgb25yZWplY3RlZCA9PSB0aHJvd2VyKSB7XHJcbiAgICAgICAgICAgIC8vIFNob3J0Y3V0IGZvciB0cml2aWFsIGFyZ3VtZW50cy5cclxuICAgICAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUpID0+IHJlc29sdmUodGhpcyBhcyBhbnkpKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGNvbnN0IGJhcnJpZXI6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0ge307XHJcbiAgICAgICAgdGhpc1tiYXJyaWVyU3ltXSA9IGJhcnJpZXI7XHJcblxyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgdm9pZCBzdXBlci50aGVuKFxyXG4gICAgICAgICAgICAgICAgKHZhbHVlKSA9PiB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cclxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xyXG5cclxuICAgICAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9uZnVsZmlsbGVkISh2YWx1ZSkpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICAgICAgKHJlYXNvbj8pID0+IHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxyXG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XHJcblxyXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25yZWplY3RlZCEocmVhc29uKSk7XHJcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgKTtcclxuICAgICAgICB9LCBhc3luYyAoY2F1c2U/KSA9PiB7XHJcbiAgICAgICAgICAgIC8vY2FuY2VsbGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBvbmNhbmNlbGxlZD8uKGNhdXNlKTtcclxuICAgICAgICAgICAgfSBmaW5hbGx5IHtcclxuICAgICAgICAgICAgICAgIGF3YWl0IHRoaXMuY2FuY2VsKGNhdXNlKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayBmb3Igb25seSB0aGUgcmVqZWN0aW9uIG9mIHRoZSBQcm9taXNlLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcclxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cclxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxyXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cclxuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIEl0IGlzIGVxdWl2YWxlbnQgdG9cclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBjYW5jZWxsYWJsZVByb21pc2UudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcclxuICAgICAqIGBgYFxyXG4gICAgICogYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXHJcbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcclxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXHJcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcclxuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxyXG4gICAgICpcclxuICAgICAqIElmIGBvbnJlamVjdGVkYCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxyXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cclxuICAgICAqL1xyXG4gICAgY2F0Y2g8VFJlc3VsdCA9IG5ldmVyPihvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gKFByb21pc2VMaWtlPFRSZXN1bHQ+IHwgVFJlc3VsdCkpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VCB8IFRSZXN1bHQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpcy50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayB0aGF0IGlzIGludm9rZWQgd2hlbiB0aGUgQ2FuY2VsbGFibGVQcm9taXNlIGlzIHNldHRsZWQgKGZ1bGZpbGxlZCBvciByZWplY3RlZCkuIFRoZVxyXG4gICAgICogcmVzb2x2ZWQgdmFsdWUgY2Fubm90IGJlIGFjY2Vzc2VkIG9yIG1vZGlmaWVkIGZyb20gdGhlIGNhbGxiYWNrLlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBzZXR0bGUgaW4gdGhlIHNhbWUgc3RhdGUgYXMgdGhlIG9yaWdpbmFsIG9uZVxyXG4gICAgICogYWZ0ZXIgdGhlIHByb3ZpZGVkIGNhbGxiYWNrIGhhcyBjb21wbGV0ZWQgZXhlY3V0aW9uLFxyXG4gICAgICogdW5sZXNzIHRoZSBjYWxsYmFjayB0aHJvd3Mgb3IgcmV0dXJucyBhIHJlamVjdGluZyBwcm9taXNlLFxyXG4gICAgICogaW4gd2hpY2ggY2FzZSB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCBhcyB3ZWxsLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcclxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cclxuICAgICAqIE9uY2UgdGhlIHBhcmVudCBwcm9taXNlIHNldHRsZXMsIHRoZSBgb25maW5hbGx5YCBjYWxsYmFjayB3aWxsIHJ1bixcclxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXHJcbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBpbXBsZW1lbnRlZCBpbiB0ZXJtcyBvZiB7QGxpbmsgdGhlbn0gYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXHJcbiAgICAgKiBJdCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxyXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXHJcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxyXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXHJcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cclxuICAgICAqXHJcbiAgICAgKiBJZiBgb25maW5hbGx5YCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxyXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cclxuICAgICAqL1xyXG4gICAgZmluYWxseShvbmZpbmFsbHk/OiAoKCkgPT4gdm9pZCkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcclxuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUuZmluYWxseSBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZmluYWxseSkpIHtcclxuICAgICAgICAgICAgcmV0dXJuIHRoaXMudGhlbihvbmZpbmFsbHksIG9uZmluYWxseSwgb25jYW5jZWxsZWQpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbihcclxuICAgICAgICAgICAgKHZhbHVlKSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB2YWx1ZSksXHJcbiAgICAgICAgICAgIChyZWFzb24/KSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB7IHRocm93IHJlYXNvbjsgfSksXHJcbiAgICAgICAgICAgIG9uY2FuY2VsbGVkLFxyXG4gICAgICAgICk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBXZSB1c2UgdGhlIGBbU3ltYm9sLnNwZWNpZXNdYCBzdGF0aWMgcHJvcGVydHksIGlmIGF2YWlsYWJsZSxcclxuICAgICAqIHRvIGRpc2FibGUgdGhlIGJ1aWx0LWluIGF1dG9tYXRpYyBzdWJjbGFzc2luZyBmZWF0dXJlcyBmcm9tIHtAbGluayBQcm9taXNlfS5cclxuICAgICAqIEl0IGlzIGNyaXRpY2FsIGZvciBwZXJmb3JtYW5jZSByZWFzb25zIHRoYXQgZXh0ZW5kZXJzIGRvIG5vdCBvdmVycmlkZSB0aGlzLlxyXG4gICAgICogT25jZSB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZ1xyXG4gICAgICogaXMgZWl0aGVyIGFjY2VwdGVkIG9yIHJldGlyZWQsIHRoaXMgaW1wbGVtZW50YXRpb24gd2lsbCBoYXZlIHRvIGJlIHJldmlzZWQgYWNjb3JkaW5nbHkuXHJcbiAgICAgKlxyXG4gICAgICogQGlnbm9yZVxyXG4gICAgICogQGludGVybmFsXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBnZXQgW3NwZWNpZXNdKCkge1xyXG4gICAgICAgIHJldHVybiBQcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xyXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUsIG9yIHJlamVjdGVkIHdoZW4gYW55IFByb21pc2UgaXMgcmVqZWN0ZWQuXHJcbiAgICAgKlxyXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgYWxsPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+W10+O1xyXG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogQXdhaXRlZDxUW1BdPjsgfT47XHJcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXHJcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxyXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcclxuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlIG9yIHJlamVjdC5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFQ+PltdPjtcclxuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFRbUF0+PjsgfT47XHJcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxyXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcclxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGxTZXR0bGVkKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVGhlIGFueSBmdW5jdGlvbiByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGlzIGZ1bGZpbGxlZCBieSB0aGUgZmlyc3QgZ2l2ZW4gcHJvbWlzZSB0byBiZSBmdWxmaWxsZWQsXHJcbiAgICAgKiBvciByZWplY3RlZCB3aXRoIGFuIEFnZ3JlZ2F0ZUVycm9yIGNvbnRhaW5pbmcgYW4gYXJyYXkgb2YgcmVqZWN0aW9uIHJlYXNvbnNcclxuICAgICAqIGlmIGFsbCBvZiB0aGUgZ2l2ZW4gcHJvbWlzZXMgYXJlIHJlamVjdGVkLlxyXG4gICAgICogSXQgcmVzb2x2ZXMgYWxsIGVsZW1lbnRzIG9mIHRoZSBwYXNzZWQgaXRlcmFibGUgdG8gcHJvbWlzZXMgYXMgaXQgcnVucyB0aGlzIGFsZ29yaXRobS5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBhbnk8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xyXG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xyXG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxyXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcclxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbnkoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XHJcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIG9yIHJlamVjdGVkIHdoZW4gYW55IG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyBhcmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQuXHJcbiAgICAgKlxyXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmFjZTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xyXG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xyXG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmFjZShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBjYW5jZWxsZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBjYW5jZWw8VCA9IG5ldmVyPihjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgY29uc3QgcCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xyXG4gICAgICAgIHAuY2FuY2VsKGNhdXNlKTtcclxuICAgICAgICByZXR1cm4gcDtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgY2FuY2Vsc1xyXG4gICAgICogYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBJZiB0aGUge0BsaW5rIEFib3J0U2lnbmFsLnRpbWVvdXR9IGZhY3RvcnkgbWV0aG9kIGlzIGF2YWlsYWJsZSxcclxuICAgICAqIGl0IGlzIHVzZWQgdG8gYmFzZSB0aGUgdGltZW91dCBvbiBfYWN0aXZlXyB0aW1lIHJhdGhlciB0aGFuIF9lbGFwc2VkXyB0aW1lLlxyXG4gICAgICogT3RoZXJ3aXNlLCBgdGltZW91dGAgZmFsbHMgYmFjayB0byB7QGxpbmsgc2V0VGltZW91dH0uXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyB0aW1lb3V0PFQgPSBuZXZlcj4obWlsbGlzZWNvbmRzOiBudW1iZXIsIGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XHJcbiAgICAgICAgaWYgKEFib3J0U2lnbmFsICYmIHR5cGVvZiBBYm9ydFNpZ25hbCA9PT0gJ2Z1bmN0aW9uJyAmJiBBYm9ydFNpZ25hbC50aW1lb3V0ICYmIHR5cGVvZiBBYm9ydFNpZ25hbC50aW1lb3V0ID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICAgICAgICAgIEFib3J0U2lnbmFsLnRpbWVvdXQobWlsbGlzZWNvbmRzKS5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpLCBtaWxsaXNlY29uZHMpO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHNsZWVwKG1pbGxpc2Vjb25kczogbnVtYmVyKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyXHJcbiAgICAgKiB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIHZhbHVlLlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHNsZWVwPFQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcclxuICAgIHN0YXRpYyBzbGVlcDxUID0gdm9pZD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlPzogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUpID0+IHtcclxuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiByZXNvbHZlKHZhbHVlISksIG1pbGxpc2Vjb25kcyk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlamVjdGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHJlYXNvbi5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlamVjdDxUID0gbmV2ZXI+KHJlYXNvbj86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKF8sIHJlamVjdCkgPT4gcmVqZWN0KHJlYXNvbikpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByZXNvbHZlKCk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICBzdGF0aWMgcmVzb2x2ZTxUID0gdm9pZD4odmFsdWU/OiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj4ge1xyXG4gICAgICAgIGlmICh2YWx1ZSBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkge1xyXG4gICAgICAgICAgICAvLyBPcHRpbWlzZSBmb3IgY2FuY2VsbGFibGUgcHJvbWlzZXMuXHJcbiAgICAgICAgICAgIHJldHVybiB2YWx1ZTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8YW55PigocmVzb2x2ZSkgPT4gcmVzb2x2ZSh2YWx1ZSkpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgYW5kIHJldHVybnMgaXQgaW4gYW4gb2JqZWN0LCBhbG9uZyB3aXRoIGl0cyByZXNvbHZlIGFuZCByZWplY3QgZnVuY3Rpb25zXHJcbiAgICAgKiBhbmQgYSBnZXR0ZXIvc2V0dGVyIGZvciB0aGUgY2FuY2VsbGF0aW9uIGNhbGxiYWNrLlxyXG4gICAgICpcclxuICAgICAqIFRoaXMgbWV0aG9kIGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyB3aXRoUmVzb2x2ZXJzPFQ+KCk6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xyXG4gICAgICAgIGxldCByZXN1bHQ6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7IG9uY2FuY2VsbGVkOiBudWxsIH0gYXMgYW55O1xyXG4gICAgICAgIHJlc3VsdC5wcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgIHJlc3VsdC5yZXNvbHZlID0gcmVzb2x2ZTtcclxuICAgICAgICAgICAgcmVzdWx0LnJlamVjdCA9IHJlamVjdDtcclxuICAgICAgICB9LCAoY2F1c2U/OiBhbnkpID0+IHsgcmVzdWx0Lm9uY2FuY2VsbGVkPy4oY2F1c2UpOyB9KTtcclxuICAgICAgICByZXR1cm4gcmVzdWx0O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgY2FuY2VsbGF0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXHJcbiAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIHJlc3VsdGluZyBmdW5jdGlvbiBkb2VzIG5vdCByZWplY3QuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYW5jZWxsZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKSB7XHJcbiAgICBsZXQgY2FuY2VsbGF0aW9uUHJvbWlzZTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0gdW5kZWZpbmVkO1xyXG5cclxuICAgIHJldHVybiAocmVhc29uOiBDYW5jZWxFcnJvcik6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9PiB7XHJcbiAgICAgICAgaWYgKCFzdGF0ZS5zZXR0bGVkKSB7XHJcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICBzdGF0ZS5yZWFzb24gPSByZWFzb247XHJcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XHJcblxyXG4gICAgICAgICAgICAvLyBBdHRhY2ggYW4gZXJyb3IgaGFuZGxlciB0aGF0IGlnbm9yZXMgdGhpcyBzcGVjaWZpYyByZWplY3Rpb24gcmVhc29uIGFuZCBub3RoaW5nIGVsc2UuXHJcbiAgICAgICAgICAgIC8vIEluIHRoZW9yeSwgYSBzYW5lIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24gYXQgdGhpcyBwb2ludFxyXG4gICAgICAgICAgICAvLyBzaG91bGQgYWx3YXlzIHJlamVjdCB3aXRoIG91ciBjYW5jZWxsYXRpb24gcmVhc29uLFxyXG4gICAgICAgICAgICAvLyBoZW5jZSB0aGUgaGFuZGxlciB3aWxsIG5ldmVyIHRocm93LlxyXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucHJvdG90eXBlLnRoZW4uY2FsbChwcm9taXNlLnByb21pc2UsIHVuZGVmaW5lZCwgKGVycikgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKGVyciAhPT0gcmVhc29uKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgdGhyb3cgZXJyO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9KTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIElmIHJlYXNvbiBpcyBub3Qgc2V0LCB0aGUgcHJvbWlzZSByZXNvbHZlZCByZWd1bGFybHksIGhlbmNlIHdlIG11c3Qgbm90IGNhbGwgb25jYW5jZWxsZWQuXHJcbiAgICAgICAgLy8gSWYgb25jYW5jZWxsZWQgaXMgdW5zZXQsIG5vIG5lZWQgdG8gZ28gYW55IGZ1cnRoZXIuXHJcbiAgICAgICAgaWYgKCFzdGF0ZS5yZWFzb24gfHwgIXByb21pc2Uub25jYW5jZWxsZWQpIHsgcmV0dXJuOyB9XHJcblxyXG4gICAgICAgIGNhbmNlbGxhdGlvblByb21pc2UgPSBuZXcgUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgcmVzb2x2ZShwcm9taXNlLm9uY2FuY2VsbGVkIShzdGF0ZS5yZWFzb24hLmNhdXNlKSk7XHJcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pLmNhdGNoKChyZWFzb24/KSA9PiB7XHJcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XHJcbiAgICAgICAgfSk7XHJcblxyXG4gICAgICAgIC8vIFVuc2V0IG9uY2FuY2VsbGVkIHRvIHByZXZlbnQgcmVwZWF0ZWQgY2FsbHMuXHJcbiAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG51bGw7XHJcblxyXG4gICAgICAgIHJldHVybiBjYW5jZWxsYXRpb25Qcm9taXNlO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVzb2x1dGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxyXG4gKi9cclxuZnVuY3Rpb24gcmVzb2x2ZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4ge1xyXG4gICAgcmV0dXJuICh2YWx1ZSkgPT4ge1xyXG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcclxuXHJcbiAgICAgICAgaWYgKHZhbHVlID09PSBwcm9taXNlLnByb21pc2UpIHtcclxuICAgICAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgVHlwZUVycm9yKFwiQSBwcm9taXNlIGNhbm5vdCBiZSByZXNvbHZlZCB3aXRoIGl0c2VsZi5cIikpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAodmFsdWUgIT0gbnVsbCAmJiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyB8fCB0eXBlb2YgdmFsdWUgPT09ICdmdW5jdGlvbicpKSB7XHJcbiAgICAgICAgICAgIGxldCB0aGVuOiBhbnk7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICB0aGVuID0gKHZhbHVlIGFzIGFueSkudGhlbjtcclxuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgICAgIHByb21pc2UucmVqZWN0KGVycik7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKHRoZW4pKSB7XHJcbiAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgIGxldCBjYW5jZWwgPSAodmFsdWUgYXMgYW55KS5jYW5jZWw7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUoY2FuY2VsKSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCBvbmNhbmNlbGxlZCA9IChjYXVzZT86IGFueSkgPT4ge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKHN0YXRlLnJlYXNvbikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gSWYgYWxyZWFkeSBjYW5jZWxsZWQsIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24uXHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSBjYW5jZWxsZXIgYWxnb3JpdGhtIGRvZXMgbm90IHJlamVjdFxyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gc28gaXQgY2FuIGJlIGRpc2NhcmRlZCBzYWZlbHkuXHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICB2b2lkIGNhbmNlbGxlckZvcih7IC4uLnByb21pc2UsIG9uY2FuY2VsbGVkIH0sIHN0YXRlKShzdGF0ZS5yZWFzb24pO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG9uY2FuY2VsbGVkO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCB7fVxyXG5cclxuICAgICAgICAgICAgICAgIGNvbnN0IG5ld1N0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcclxuICAgICAgICAgICAgICAgICAgICByb290OiBzdGF0ZS5yb290LFxyXG4gICAgICAgICAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHNldHRsZWQoKSB7IHJldHVybiB0aGlzLnJvb3Quc2V0dGxlZCB9LFxyXG4gICAgICAgICAgICAgICAgICAgIHNldCBzZXR0bGVkKHZhbHVlKSB7IHRoaXMucm9vdC5zZXR0bGVkID0gdmFsdWU7IH0sXHJcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHJlYXNvbigpIHsgcmV0dXJuIHRoaXMucm9vdC5yZWFzb24gfVxyXG4gICAgICAgICAgICAgICAgfTtcclxuXHJcbiAgICAgICAgICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIG5ld1N0YXRlKTtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseSh0aGVuLCB2YWx1ZSwgW3Jlc29sdmVyRm9yKHByb21pc2UsIG5ld1N0YXRlKSwgcmVqZWN0b3JdKTtcclxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICByZXR1cm47IC8vIElNUE9SVEFOVCFcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgcHJvbWlzZS5yZXNvbHZlKHZhbHVlKTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZWplY3Rpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cclxuICovXHJcbmZ1bmN0aW9uIHJlamVjdG9yRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yIHtcclxuICAgIHJldHVybiAocmVhc29uPykgPT4ge1xyXG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcclxuXHJcbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIGlmIChyZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBzdGF0ZS5yZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBPYmplY3QuaXMocmVhc29uLmNhdXNlLCBzdGF0ZS5yZWFzb24uY2F1c2UpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gU3dhbGxvdyBsYXRlIHJlamVjdGlvbnMgdGhhdCBhcmUgQ2FuY2VsRXJyb3JzIHdob3NlIGNhbmNlbGxhdGlvbiBjYXVzZSBpcyB0aGUgc2FtZSBhcyBvdXJzLlxyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSBjYXRjaCB7fVxyXG5cclxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24pKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYW5jZWxzIGFsbCB2YWx1ZXMgaW4gYW4gYXJyYXkgdGhhdCBsb29rIGxpa2UgY2FuY2VsbGFibGUgdGhlbmFibGVzLlxyXG4gKiBSZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGZ1bGZpbGxzIG9uY2UgYWxsIGNhbmNlbGxhdGlvbiBwcm9jZWR1cmVzIGZvciB0aGUgZ2l2ZW4gdmFsdWVzIGhhdmUgc2V0dGxlZC5cclxuICovXHJcbmZ1bmN0aW9uIGNhbmNlbEFsbChwYXJlbnQ6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgdmFsdWVzOiBhbnlbXSwgY2F1c2U/OiBhbnkpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIGNvbnN0IHJlc3VsdHM6IFByb21pc2U8dm9pZD5bXSA9IFtdO1xyXG5cclxuICAgIGZvciAoY29uc3QgdmFsdWUgb2YgdmFsdWVzKSB7XHJcbiAgICAgICAgbGV0IGNhbmNlbDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyO1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZSh2YWx1ZS50aGVuKSkgeyBjb250aW51ZTsgfVxyXG4gICAgICAgICAgICBjYW5jZWwgPSB2YWx1ZS5jYW5jZWw7XHJcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZShjYW5jZWwpKSB7IGNvbnRpbnVlOyB9XHJcbiAgICAgICAgfSBjYXRjaCB7IGNvbnRpbnVlOyB9XHJcblxyXG4gICAgICAgIGxldCByZXN1bHQ6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICByZXN1bHQgPSBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xyXG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XHJcbiAgICAgICAgICAgIGNvbnRpbnVlO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKCFyZXN1bHQpIHsgY29udGludWU7IH1cclxuICAgICAgICByZXN1bHRzLnB1c2goXHJcbiAgICAgICAgICAgIChyZXN1bHQgaW5zdGFuY2VvZiBQcm9taXNlICA/IHJlc3VsdCA6IFByb21pc2UucmVzb2x2ZShyZXN1bHQpKS5jYXRjaCgocmVhc29uPykgPT4ge1xyXG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xyXG4gICAgICAgICAgICB9KVxyXG4gICAgICAgICk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIFByb21pc2UuYWxsKHJlc3VsdHMpIGFzIGFueTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgaXRzIGFyZ3VtZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gaWRlbnRpdHk8VD4oeDogVCk6IFQge1xyXG4gICAgcmV0dXJuIHg7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUaHJvd3MgaXRzIGFyZ3VtZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gdGhyb3dlcihyZWFzb24/OiBhbnkpOiBuZXZlciB7XHJcbiAgICB0aHJvdyByZWFzb247XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBdHRlbXB0cyB2YXJpb3VzIHN0cmF0ZWdpZXMgdG8gY29udmVydCBhbiBlcnJvciB0byBhIHN0cmluZy5cclxuICovXHJcbmZ1bmN0aW9uIGVycm9yTWVzc2FnZShlcnI6IGFueSk6IHN0cmluZyB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGlmIChlcnIgaW5zdGFuY2VvZiBFcnJvciB8fCB0eXBlb2YgZXJyICE9PSAnb2JqZWN0JyB8fCBlcnIudG9TdHJpbmcgIT09IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcpIHtcclxuICAgICAgICAgICAgcmV0dXJuIFwiXCIgKyBlcnI7XHJcbiAgICAgICAgfVxyXG4gICAgfSBjYXRjaCB7fVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgcmV0dXJuIEpTT04uc3RyaW5naWZ5KGVycik7XHJcbiAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZy5jYWxsKGVycik7XHJcbiAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgcmV0dXJuIFwiPGNvdWxkIG5vdCBjb252ZXJ0IGVycm9yIHRvIHN0cmluZz5cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIGN1cnJlbnQgYmFycmllciBwcm9taXNlIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS4gSWYgbmVjZXNzYXJ5LCBpbml0aWFsaXNlcyB0aGUgYmFycmllci5cclxuICovXHJcbmZ1bmN0aW9uIGN1cnJlbnRCYXJyaWVyPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgbGV0IHB3cjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSBwcm9taXNlW2JhcnJpZXJTeW1dID8/IHt9O1xyXG4gICAgaWYgKCEoJ3Byb21pc2UnIGluIHB3cikpIHtcclxuICAgICAgICBPYmplY3QuYXNzaWduKHB3ciwgcHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4oKSk7XHJcbiAgICB9XHJcbiAgICBpZiAocHJvbWlzZVtiYXJyaWVyU3ltXSA9PSBudWxsKSB7XHJcbiAgICAgICAgcHdyLnJlc29sdmUhKCk7XHJcbiAgICAgICAgcHJvbWlzZVtiYXJyaWVyU3ltXSA9IHB3cjtcclxuICAgIH1cclxuICAgIHJldHVybiBwd3IucHJvbWlzZSE7XHJcbn1cclxuXHJcbi8vIFBvbHlmaWxsIFByb21pc2Uud2l0aFJlc29sdmVycy5cclxubGV0IHByb21pc2VXaXRoUmVzb2x2ZXJzID0gUHJvbWlzZS53aXRoUmVzb2x2ZXJzO1xyXG5pZiAocHJvbWlzZVdpdGhSZXNvbHZlcnMgJiYgdHlwZW9mIHByb21pc2VXaXRoUmVzb2x2ZXJzID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IHByb21pc2VXaXRoUmVzb2x2ZXJzLmJpbmQoUHJvbWlzZSk7XHJcbn0gZWxzZSB7XHJcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IGZ1bmN0aW9uIDxUPigpOiBQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XHJcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcclxuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IFByb21pc2U8VD4oKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XHJcbiAgICAgICAgcmV0dXJuIHsgcHJvbWlzZSwgcmVzb2x2ZSwgcmVqZWN0IH07XHJcbiAgICB9XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgQ2FuY2VsbGFibGVQcm9taXNlLCB0eXBlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnMgfSBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XHJcblxyXG4vLyBTZXR1cFxyXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xyXG5cclxuaWYgKGhhc0RPTSkge1xyXG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbn1cclxuXHJcbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8Q2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIiB8IFwib25jYW5jZWxsZWRcIj5cclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbGwpO1xyXG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYW5jZWxDYWxsKTtcclxuY29uc3QgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xyXG5cclxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xyXG5jb25zdCBDYW5jZWxNZXRob2QgPSAwXHJcblxyXG4vKipcclxuICogSG9sZHMgYWxsIHJlcXVpcmVkIGluZm9ybWF0aW9uIGZvciBhIGJpbmRpbmcgY2FsbC5cclxuICogTWF5IHByb3ZpZGUgZWl0aGVyIGEgbWV0aG9kIElEIG9yIGEgbWV0aG9kIG5hbWUsIGJ1dCBub3QgYm90aC5cclxuICovXHJcbmV4cG9ydCB0eXBlIENhbGxPcHRpb25zID0ge1xyXG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZElEOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZE5hbWU/OiBuZXZlcjtcclxuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cclxuICAgIGFyZ3M6IGFueVtdO1xyXG59IHwge1xyXG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZElEPzogbmV2ZXI7XHJcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZE5hbWU6IHN0cmluZztcclxuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cclxuICAgIGFyZ3M6IGFueVtdO1xyXG59O1xyXG5cclxuLy8gcnVudGltZS5qcyBuZWVkcyB0byB1c2UgUnVudGltZUVycm9yIGludGVybmFsbHkgdG8gcHJvcGVybHkgcGFyc2UgYW5kIHJldHVyblxyXG4vLyBlcnJvcnMgZm9yIGJpbmRpbmcgY2FsbHMsIHNvIGl0IGhhZCB0byBtb3ZlIHRoZXJlLiBFeHBvcnRpbmcgaGVyZSBhZ2FpbiB0b1xyXG4vLyBrZWVwIGZyb20gYnJlYWtpbmcgdGhlIHB1YmxpYyBDYWxsIGludGVyZmFjZS5cclxuZXhwb3J0IHsgUnVudGltZUVycm9yIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuLyoqXHJcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBJRCB1c2luZyB0aGUgbmFub2lkIGxpYnJhcnkuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgdW5pcXVlIElEIHRoYXQgZG9lcyBub3QgZXhpc3QgaW4gdGhlIGNhbGxSZXNwb25zZXMgc2V0LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGwgYSBib3VuZCBtZXRob2QgYWNjb3JkaW5nIHRvIHRoZSBnaXZlbiBjYWxsIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEluIGNhc2Ugb2YgZmFpbHVyZSwgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhbiBleGNlcHRpb25cclxuICogYW1vbmcgUmVmZXJlbmNlRXJyb3IgKHVua25vd24gbWV0aG9kKSwgVHlwZUVycm9yICh3cm9uZyBhcmd1bWVudCBjb3VudCBvciB0eXBlKSxcclxuICoge0BsaW5rIFJ1bnRpbWVFcnJvcn0gKG1ldGhvZCByZXR1cm5lZCBhbiBlcnJvciksIG9yIG90aGVyIChuZXR3b3JrIG9yIGludGVybmFsIGVycm9ycykuXHJcbiAqIFRoZSBleGNlcHRpb24gbWlnaHQgaGF2ZSBhIFwiY2F1c2VcIiBmaWVsZCB3aXRoIHRoZSB2YWx1ZSByZXR1cm5lZFxyXG4gKiBieSB0aGUgYXBwbGljYXRpb24tIG9yIHNlcnZpY2UtbGV2ZWwgZXJyb3IgbWFyc2hhbGluZyBmdW5jdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gQSBtZXRob2QgY2FsbCBkZXNjcmlwdG9yLlxyXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9uczogQ2FsbE9wdGlvbnMpOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XHJcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuXHJcbiAgICBjb25zdCByZXN1bHQgPSBDYW5jZWxsYWJsZVByb21pc2Uud2l0aFJlc29sdmVyczxhbnk+KCk7XHJcbiAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlOiByZXN1bHQucmVzb2x2ZSwgcmVqZWN0OiByZXN1bHQucmVqZWN0IH0pO1xyXG5cclxuICAgIGNvbnN0IHJlcXVlc3QgPSBjYWxsKENhbGxCaW5kaW5nLCBPYmplY3QuYXNzaWduKHsgXCJjYWxsLWlkXCI6IGlkIH0sIG9wdGlvbnMpKTtcclxuICAgIGxldCBydW5uaW5nID0gdHJ1ZTtcclxuXHJcbiAgICByZXF1ZXN0LnRoZW4oKHJlcykgPT4ge1xyXG4gICAgICAgIHJ1bm5pbmcgPSBmYWxzZTtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgcmVzdWx0LnJlc29sdmUocmVzKTtcclxuICAgIH0sIChlcnIpID0+IHtcclxuICAgICAgICBydW5uaW5nID0gZmFsc2U7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIHJlc3VsdC5yZWplY3QoZXJyKTtcclxuICAgIH0pO1xyXG5cclxuICAgIGNvbnN0IGNhbmNlbCA9ICgpID0+IHtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgcmV0dXJuIGNhbmNlbENhbGwoQ2FuY2VsTWV0aG9kLCB7XCJjYWxsLWlkXCI6IGlkfSkuY2F0Y2goKGVycikgPT4ge1xyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igd2hpbGUgcmVxdWVzdGluZyBiaW5kaW5nIGNhbGwgY2FuY2VsbGF0aW9uOlwiLCBlcnIpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfTtcclxuXHJcbiAgICByZXN1bHQub25jYW5jZWxsZWQgPSAoKSA9PiB7XHJcbiAgICAgICAgaWYgKHJ1bm5pbmcpIHtcclxuICAgICAgICAgICAgcmV0dXJuIGNhbmNlbCgpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHJldHVybiByZXF1ZXN0LnRoZW4oY2FuY2VsKTtcclxuICAgICAgICB9XHJcbiAgICB9O1xyXG5cclxuICAgIHJldHVybiByZXN1bHQucHJvbWlzZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgYm91bmQgbWV0aG9kIGJ5IG5hbWUgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cclxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cclxuICpcclxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXHJcbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShtZXRob2ROYW1lOiBzdHJpbmcsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2ROYW1lLCBhcmdzIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIG51bWVyaWMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cclxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cclxuICpcclxuICogQHBhcmFtIG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cclxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cclxuICogQHJldHVybiBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeUlEKG1ldGhvZElEOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2RJRCwgYXJncyB9KTtcclxufVxyXG4iLCIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DbGlwYm9hcmQpO1xyXG5cclxuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XHJcbmNvbnN0IENsaXBib2FyZFRleHQgPSAxO1xyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHRleHQgdG8gdGhlIENsaXBib2FyZC5cclxuICpcclxuICogQHBhcmFtIHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cclxuICogQHJldHVybiBBIFByb21pc2UgdGhhdCByZXNvbHZlcyB3aGVuIHRoZSBvcGVyYXRpb24gaXMgc3VjY2Vzc2Z1bC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQ6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpOiBQcm9taXNlPHN0cmluZz4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkVGV4dCk7XHJcbn1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIEFueSBpcyBhIGR1bW15IGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBzaW1wbGUgb3IgdW5rbm93biB0eXBlcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBBbnk8VCA9IGFueT4oc291cmNlOiBhbnkpOiBUIHtcclxuICAgIHJldHVybiBzb3VyY2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXHJcbiAqIG51bGwgc3RyaW5ncyB3aXRoIGVtcHR5IHN0cmluZ3MuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnl0ZVNsaWNlKHNvdXJjZTogYW55KTogc3RyaW5nIHtcclxuICAgIHJldHVybiAoKHNvdXJjZSA9PSBudWxsKSA/IFwiXCIgOiBzb3VyY2UpO1xyXG59XHJcblxyXG4vKipcclxuICogQXJyYXkgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XHJcbiAqIHdob3NlIGVsZW1lbnRzIGFyZSBvZiB0aGF0IHR5cGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQXJyYXk8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiBUW10ge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm4gW107XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XHJcbiAgICAgICAgICAgIHNvdXJjZVtpXSA9IGVsZW1lbnQoc291cmNlW2ldKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcclxuICogd2hvc2Uga2V5cyBhbmQgdmFsdWVzIGFyZSBvZiB0aG9zZSB0eXBlcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBNYXA8SyBleHRlbmRzIFByb3BlcnR5S2V5ID0gYW55LCBWID0gYW55PihrZXk6IChzb3VyY2U6IGFueSkgPT4gSywgdmFsdWU6IChzb3VyY2U6IGFueSkgPT4gVik6IChzb3VyY2U6IGFueSkgPT4gUmVjb3JkPEssIFY+IHtcclxuICAgIGlmICh2YWx1ZSA9PT0gQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyB7fSA6IHNvdXJjZSk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcclxuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XHJcbiAgICAgICAgICAgIHJldHVybiB7fTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZm9yIChjb25zdCBrZXkgaW4gc291cmNlKSB7XHJcbiAgICAgICAgICAgIHNvdXJjZVtrZXldID0gdmFsdWUoc291cmNlW2tleV0pO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gc291cmNlO1xyXG4gICAgfTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE51bGxhYmxlIHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXHJcbiAqIGFuZCByZXR1cm5zIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgbnVsbGFibGUgdmFsdWUgb2YgdGhhdCB0eXBlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE51bGxhYmxlPFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gKFQgfCBudWxsKSB7XHJcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIEFueTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IG51bGwgOiBlbGVtZW50KHNvdXJjZSkpO1xyXG59XHJcblxyXG4vKipcclxuICogU3RydWN0IHRha2VzIGFuIG9iamVjdCBtYXBwaW5nIGZpZWxkIG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9uc1xyXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBzdHJ1Y3QuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU3RydWN0KGNyZWF0ZUZpZWxkOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4pOlxyXG4gICAgPFUgZXh0ZW5kcyBSZWNvcmQ8c3RyaW5nLCBhbnk+ID0gYW55Pihzb3VyY2U6IGFueSkgPT4gVVxyXG57XHJcbiAgICBsZXQgYWxsQW55ID0gdHJ1ZTtcclxuICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xyXG4gICAgICAgIGlmIChjcmVhdGVGaWVsZFtuYW1lXSAhPT0gQW55KSB7XHJcbiAgICAgICAgICAgIGFsbEFueSA9IGZhbHNlO1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcbiAgICB9XHJcbiAgICBpZiAoYWxsQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIEFueTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xyXG4gICAgICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xyXG4gICAgICAgICAgICBpZiAobmFtZSBpbiBzb3VyY2UpIHtcclxuICAgICAgICAgICAgICAgIHNvdXJjZVtuYW1lXSA9IGNyZWF0ZUZpZWxkW25hbWVdKHNvdXJjZVtuYW1lXSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBEYXRlRnJvbVRpbWUgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IGNvbnZlcnRzIFJGQzMzMzkgc3RyaW5nc1xyXG4gKiAoZnJvbSBHbydzIHRpbWUuVGltZSBKU09OIG1hcnNoYWxpbmcpIHRvIEphdmFTY3JpcHQgRGF0ZSBvYmplY3RzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIERhdGVGcm9tVGltZShzb3VyY2U6IGFueSk6IERhdGUge1xyXG4gICAgcmV0dXJuIG5ldyBEYXRlKHNvdXJjZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXBzIGtub3duIGV2ZW50IG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdGhlaXIgZGF0YSB0eXBlcy5cclxuICogV2lsbCBiZSBtb25rZXktcGF0Y2hlZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgRXZlbnRzOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4gPSB7fTtcclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuLy8gc2V0dXBcclxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcclxuXHJcbmlmIChoYXNET00pIHtcclxuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG59XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5EaWFsb2cpO1xyXG5cclxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcclxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XHJcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xyXG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XHJcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcclxuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xyXG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xyXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xyXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXHJcbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC4gKi9cclxuICAgIEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xyXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cclxuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXHJcbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cclxuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXHJcbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xyXG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXHJcbiAgICBUaXRsZT86IHN0cmluZztcclxuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XHJcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXHJcbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xyXG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cclxuICAgIERldGFjaGVkPzogYm9vbGVhbjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBTYXZlRmlsZURpYWxvZ09wdGlvbnMge1xyXG4gICAgLyoqIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBGaWxlbmFtZT86IHN0cmluZztcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cclxuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cclxuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXHJcbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXHJcbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cclxuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cclxuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXHJcbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xyXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXHJcbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XHJcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xyXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcclxuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgVGl0bGU/OiBzdHJpbmc7XHJcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBNZXNzYWdlPzogc3RyaW5nO1xyXG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xyXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcclxuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXHJcbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgTWVzc2FnZURpYWxvZ09wdGlvbnMge1xyXG4gICAgLyoqIFRoZSB0aXRsZSBvZiB0aGUgZGlhbG9nIHdpbmRvdy4gKi9cclxuICAgIFRpdGxlPzogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSBtYWluIG1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgTWVzc2FnZT86IHN0cmluZztcclxuICAgIC8qKiBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBCdXR0b25zPzogQnV0dG9uW107XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLiAqL1xyXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIEJ1dHRvbiB7XHJcbiAgICAvKiogVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uICovXHJcbiAgICBMYWJlbD86IHN0cmluZztcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLiAqL1xyXG4gICAgSXNDYW5jZWw/OiBib29sZWFuO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci4gKi9cclxuICAgIElzRGVmYXVsdD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgRmlsZUZpbHRlciB7XHJcbiAgICAvKiogRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuICovXHJcbiAgICBEaXNwbGF5TmFtZT86IHN0cmluZztcclxuICAgIC8qKiBQYXR0ZXJuIHRvIG1hdGNoIGZvciB0aGUgZmlsdGVyLCBlLmcuIFwiKi50eHQ7Ki5tZFwiIGZvciB0ZXh0IG1hcmtkb3duIGZpbGVzLiAqL1xyXG4gICAgUGF0dGVybj86IHN0cmluZztcclxufVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSB0eXBlIC0gRGlhbG9nIHR5cGUuXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyBmb3IgdGhlIGRpYWxvZy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nLlxyXG4gKi9cclxuZnVuY3Rpb24gZGlhbG9nKHR5cGU6IG51bWJlciwgb3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMgfCBPcGVuRmlsZURpYWxvZ09wdGlvbnMgfCBTYXZlRmlsZURpYWxvZ09wdGlvbnMgPSB7fSk6IFByb21pc2U8YW55PiB7XHJcbiAgICByZXR1cm4gY2FsbCh0eXBlLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGFuIGluZm8gZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSB3YXJuaW5nIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhbiBlcnJvciBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIHF1ZXN0aW9uIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgb25lIG9yIG1vcmUgZmlsZXMgdG8gb3Blbi5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSBvciBsaXN0IG9mIGZpbGVzLCBvciBhIGJsYW5rIHN0cmluZy9lbXB0eSBsaXN0IGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbjogdHJ1ZSB9KTogUHJvbWlzZTxzdHJpbmdbXT47XHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogZmFsc2UgfCB1bmRlZmluZWQgfSk6IFByb21pc2U8c3RyaW5nPjtcclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+O1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT4geyByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKSA/PyBbXTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgYSBmaWxlIHRvIHNhdmUuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUsIG9yIGEgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9uczogU2F2ZUZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7IH1cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLy8gVGhlIGZvbGxvd2luZyB1dGlsaXRpZXMgaGF2ZSBiZWVuIGZhY3RvcmVkIG91dCBvZiAuL2V2ZW50cy50c1xyXG4vLyBmb3IgdGVzdGluZyBwdXJwb3Nlcy5cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXA8c3RyaW5nLCBMaXN0ZW5lcltdPigpO1xyXG5cclxuZXhwb3J0IGNsYXNzIExpc3RlbmVyIHtcclxuICAgIGV2ZW50TmFtZTogc3RyaW5nO1xyXG4gICAgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQ7XHJcbiAgICBtYXhDYWxsYmFja3M6IG51bWJlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQsIG1heENhbGxiYWNrczogbnVtYmVyKSB7XHJcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XHJcbiAgICAgICAgdGhpcy5jYWxsYmFjayA9IGNhbGxiYWNrO1xyXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xyXG4gICAgfVxyXG5cclxuICAgIGRpc3BhdGNoKGRhdGE6IGFueSk6IGJvb2xlYW4ge1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIHRoaXMuY2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZXJyKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyOiBMaXN0ZW5lcik6IHZvaWQge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChsaXN0ZW5lci5ldmVudE5hbWUpO1xyXG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGxpc3RlbmVyLmV2ZW50TmFtZSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChsaXN0ZW5lci5ldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuIiwiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIEN5bmh5cmNod3lkIHkgZmZlaWwgaG9uIHluIGF3dG9tYXRpZy4gUEVJRElXQ0ggw4IgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRXYWtlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFdha2VcIixcblx0XHRBcHBsaWNhdGlvblNjcmVlbnNEaWRTbGVlcDogXCJtYWM6QXBwbGljYXRpb25TY3JlZW5zRGlkU2xlZXBcIixcblx0XHRBcHBsaWNhdGlvblNjcmVlbnNEaWRXYWtlOiBcIm1hYzpBcHBsaWNhdGlvblNjcmVlbnNEaWRXYWtlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFNsZWVwOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxTbGVlcFwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdFdlYlZpZXdXZWJDb250ZW50UHJvY2Vzc0RpZFRlcm1pbmF0ZTogXCJtYWM6V2ViVmlld1dlYkNvbnRlbnRQcm9jZXNzRGlkVGVybWluYXRlXCIsXG5cdH0pLFxuXHRMaW51eDogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxuXHRcdFN5c3RlbURpZFdha2U6IFwibGludXg6U3lzdGVtRGlkV2FrZVwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRTeXN0ZW1XaWxsU2xlZXA6IFwibGludXg6U3lzdGVtV2lsbFNsZWVwXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibGludXg6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcblx0XHRXaW5kb3dMb2FkU3RhcnRlZDogXCJsaW51eDpXaW5kb3dMb2FkU3RhcnRlZFwiLFxuXHRcdFdpbmRvd0xvYWRSZWRpcmVjdGVkOiBcImxpbnV4OldpbmRvd0xvYWRSZWRpcmVjdGVkXCIsXG5cdFx0V2luZG93TG9hZENvbW1pdHRlZDogXCJsaW51eDpXaW5kb3dMb2FkQ29tbWl0dGVkXCIsXG5cdFx0V2luZG93TG9hZEZpbmlzaGVkOiBcImxpbnV4OldpbmRvd0xvYWRGaW5pc2hlZFwiLFxuXHR9KSxcblx0aU9TOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJpb3M6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZEVudGVyQmFja2dyb3VuZDogXCJpb3M6QXBwbGljYXRpb25EaWRFbnRlckJhY2tncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZzogXCJpb3M6QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmdcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nOiBcImlvczpBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsRW50ZXJGb3JlZ3JvdW5kOiBcImlvczpBcHBsaWNhdGlvbldpbGxFbnRlckZvcmVncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJpb3M6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0V2luZG93RGlkTG9hZDogXCJpb3M6V2luZG93RGlkTG9hZFwiLFxuXHRcdFdpbmRvd1dpbGxBcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxBcHBlYXJcIixcblx0XHRXaW5kb3dEaWRBcHBlYXI6IFwiaW9zOldpbmRvd0RpZEFwcGVhclwiLFxuXHRcdFdpbmRvd1dpbGxEaXNhcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxEaXNhcHBlYXJcIixcblx0XHRXaW5kb3dEaWREaXNhcHBlYXI6IFwiaW9zOldpbmRvd0RpZERpc2FwcGVhclwiLFxuXHRcdFdpbmRvd1NhZmVBcmVhSW5zZXRzQ2hhbmdlZDogXCJpb3M6V2luZG93U2FmZUFyZWFJbnNldHNDaGFuZ2VkXCIsXG5cdFx0V2luZG93T3JpZW50YXRpb25DaGFuZ2VkOiBcImlvczpXaW5kb3dPcmllbnRhdGlvbkNoYW5nZWRcIixcblx0XHRXaW5kb3dUb3VjaEJlZ2FuOiBcImlvczpXaW5kb3dUb3VjaEJlZ2FuXCIsXG5cdFx0V2luZG93VG91Y2hNb3ZlZDogXCJpb3M6V2luZG93VG91Y2hNb3ZlZFwiLFxuXHRcdFdpbmRvd1RvdWNoRW5kZWQ6IFwiaW9zOldpbmRvd1RvdWNoRW5kZWRcIixcblx0XHRXaW5kb3dUb3VjaENhbmNlbGxlZDogXCJpb3M6V2luZG93VG91Y2hDYW5jZWxsZWRcIixcblx0XHRXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZhaWxOYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkRmFpbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGVjaWRlUG9saWN5Rm9yTmF2aWdhdGlvbkFjdGlvbjogXCJpb3M6V2ViVmlld0RlY2lkZVBvbGljeUZvck5hdmlnYXRpb25BY3Rpb25cIixcblx0XHRCYXR0ZXJ5Q2hhbmdlZDogXCJpb3M6QmF0dGVyeUNoYW5nZWRcIixcblx0XHROZXR3b3JrQ2hhbmdlZDogXCJpb3M6TmV0d29ya0NoYW5nZWRcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiaW9zOlRoZW1lQ2hhbmdlZFwiLFxuXHRcdFNjcmVlbkxvY2tlZDogXCJpb3M6U2NyZWVuTG9ja2VkXCIsXG5cdFx0U2NyZWVuVW5sb2NrZWQ6IFwiaW9zOlNjcmVlblVubG9ja2VkXCIsXG5cdH0pLFxuXHRBbmRyb2lkOiBPYmplY3QuZnJlZXplKHtcblx0XHRBY3Rpdml0eUNyZWF0ZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eUNyZWF0ZWRcIixcblx0XHRBY3Rpdml0eVN0YXJ0ZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eVN0YXJ0ZWRcIixcblx0XHRBY3Rpdml0eVJlc3VtZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eVJlc3VtZWRcIixcblx0XHRBY3Rpdml0eVBhdXNlZDogXCJhbmRyb2lkOkFjdGl2aXR5UGF1c2VkXCIsXG5cdFx0QWN0aXZpdHlTdG9wcGVkOiBcImFuZHJvaWQ6QWN0aXZpdHlTdG9wcGVkXCIsXG5cdFx0QWN0aXZpdHlEZXN0cm95ZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eURlc3Ryb3llZFwiLFxuXHRcdEFwcGxpY2F0aW9uTG93TWVtb3J5OiBcImFuZHJvaWQ6QXBwbGljYXRpb25Mb3dNZW1vcnlcIixcblx0XHRXZWJWaWV3UGFnZVN0YXJ0ZWQ6IFwiYW5kcm9pZDpXZWJWaWV3UGFnZVN0YXJ0ZWRcIixcblx0XHRXZWJWaWV3UGFnZUZpbmlzaGVkOiBcImFuZHJvaWQ6V2ViVmlld1BhZ2VGaW5pc2hlZFwiLFxuXHRcdEJhdHRlcnlDaGFuZ2VkOiBcImFuZHJvaWQ6QmF0dGVyeUNoYW5nZWRcIixcblx0XHROZXR3b3JrQ2hhbmdlZDogXCJhbmRyb2lkOk5ldHdvcmtDaGFuZ2VkXCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImFuZHJvaWQ6VGhlbWVDaGFuZ2VkXCIsXG5cdFx0U2NyZWVuTG9ja2VkOiBcImFuZHJvaWQ6U2NyZWVuTG9ja2VkXCIsXG5cdFx0U2NyZWVuVW5sb2NrZWQ6IFwiYW5kcm9pZDpTY3JlZW5VbmxvY2tlZFwiLFxuXHR9KSxcblx0T3Blbmhhcm1vbnk6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFiaWxpdHlDcmVhdGVkOiBcIm9wZW5oYXJtb255OkFiaWxpdHlDcmVhdGVkXCIsXG5cdFx0QWJpbGl0eVN0YXJ0ZWQ6IFwib3Blbmhhcm1vbnk6QWJpbGl0eVN0YXJ0ZWRcIixcblx0XHRBYmlsaXR5Rm9yZWdyb3VuZDogXCJvcGVuaGFybW9ueTpBYmlsaXR5Rm9yZWdyb3VuZFwiLFxuXHRcdEFiaWxpdHlCYWNrZ3JvdW5kOiBcIm9wZW5oYXJtb255OkFiaWxpdHlCYWNrZ3JvdW5kXCIsXG5cdFx0QWJpbGl0eURlc3Ryb3llZDogXCJvcGVuaGFybW9ueTpBYmlsaXR5RGVzdHJveWVkXCIsXG5cdFx0QXBwbGljYXRpb25Mb3dNZW1vcnk6IFwib3Blbmhhcm1vbnk6QXBwbGljYXRpb25Mb3dNZW1vcnlcIixcblx0XHRXZWJWaWV3UGFnZVN0YXJ0ZWQ6IFwib3Blbmhhcm1vbnk6V2ViVmlld1BhZ2VTdGFydGVkXCIsXG5cdFx0V2ViVmlld1BhZ2VGaW5pc2hlZDogXCJvcGVuaGFybW9ueTpXZWJWaWV3UGFnZUZpbmlzaGVkXCIsXG5cdFx0QmF0dGVyeUNoYW5nZWQ6IFwib3Blbmhhcm1vbnk6QmF0dGVyeUNoYW5nZWRcIixcblx0XHROZXR3b3JrQ2hhbmdlZDogXCJvcGVuaGFybW9ueTpOZXR3b3JrQ2hhbmdlZFwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJvcGVuaGFybW9ueTpUaGVtZUNoYW5nZWRcIixcblx0XHRTY3JlZW5Mb2NrZWQ6IFwib3Blbmhhcm1vbnk6U2NyZWVuTG9ja2VkXCIsXG5cdFx0U2NyZWVuVW5sb2NrZWQ6IFwib3Blbmhhcm1vbnk6U2NyZWVuVW5sb2NrZWRcIixcblx0fSksXG5cdENvbW1vbjogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZTogXCJjb21tb246QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZVwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmw6IFwiY29tbW9uOkFwcGxpY2F0aW9uTGF1bmNoZWRXaXRoVXJsXCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImNvbW1vbjpUaGVtZUNoYW5nZWRcIixcblx0XHRTeXN0ZW1EaWRXYWtlOiBcImNvbW1vbjpTeXN0ZW1EaWRXYWtlXCIsXG5cdFx0U3lzdGVtV2lsbFNsZWVwOiBcImNvbW1vbjpTeXN0ZW1XaWxsU2xlZXBcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJjb21tb246V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJjb21tb246V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RmlsZXNEcm9wcGVkXCIsXG5cdFx0V2luZG93Rm9jdXM6IFwiY29tbW9uOldpbmRvd0ZvY3VzXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd0hpZGU6IFwiY29tbW9uOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dMb3N0Rm9jdXM6IFwiY29tbW9uOldpbmRvd0xvc3RGb2N1c1wiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1RvZ2dsZUZyYW1lbGVzczogXCJjb21tb246V2luZG93VG9nZ2xlRnJhbWVsZXNzXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1J1bnRpbWVSZWFkeTogXCJjb21tb246V2luZG93UnVudGltZVJlYWR5XCIsXG5cdFx0V2luZG93U2hvdzogXCJjb21tb246V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcImNvbW1vbjpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcblx0XHRCYXR0ZXJ5Q2hhbmdlZDogXCJjb21tb246QmF0dGVyeUNoYW5nZWRcIixcblx0XHROZXR3b3JrQ2hhbmdlZDogXCJjb21tb246TmV0d29ya0NoYW5nZWRcIixcblx0XHRTY3JlZW5Mb2NrZWQ6IFwiY29tbW9uOlNjcmVlbkxvY2tlZFwiLFxuXHRcdFNjcmVlblVubG9ja2VkOiBcImNvbW1vbjpTY3JlZW5VbmxvY2tlZFwiLFxuXHRcdExvd01lbW9yeTogXCJjb21tb246TG93TWVtb3J5XCIsXG5cdH0pLFxufSk7XG4iLCIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgZXZlbnRMaXN0ZW5lcnMsIExpc3RlbmVyLCBsaXN0ZW5lck9mZiB9IGZyb20gXCIuL2xpc3RlbmVyLmpzXCI7XHJcbmltcG9ydCB7IEV2ZW50cyBhcyBDcmVhdGUgfSBmcm9tIFwiLi9jcmVhdGUuanNcIjtcclxuaW1wb3J0IHsgVHlwZXMgfSBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xyXG5cclxuLy8gU2V0dXBcclxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcclxuXHJcbmlmIChoYXNET00pIHtcclxuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG4gICAgd2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XHJcbn1cclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkV2ZW50cyk7XHJcbmNvbnN0IEVtaXRNZXRob2QgPSAwO1xyXG5cclxuZXhwb3J0ICogZnJvbSBcIi4vZXZlbnRfdHlwZXMuanNcIjtcclxuXHJcbi8qKlxyXG4gKiBBIHRhYmxlIG9mIGRhdGEgdHlwZXMgZm9yIGFsbCBrbm93biBldmVudHMuXHJcbiAqIFdpbGwgYmUgbW9ua2V5LXBhdGNoZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxyXG4gKi9cclxuZXhwb3J0IGludGVyZmFjZSBDdXN0b21FdmVudHMge31cclxuXHJcbi8qKlxyXG4gKiBFaXRoZXIgYSBrbm93biBldmVudCBuYW1lIG9yIGFuIGFyYml0cmFyeSBzdHJpbmcuXHJcbiAqL1xyXG5leHBvcnQgdHlwZSBXYWlsc0V2ZW50TmFtZTxFIGV4dGVuZHMga2V5b2YgQ3VzdG9tRXZlbnRzID0ga2V5b2YgQ3VzdG9tRXZlbnRzPiA9IEUgfCAoc3RyaW5nICYge30pO1xyXG5cclxuLyoqXHJcbiAqIFVuaW9uIG9mIGFsbCBrbm93biBzeXN0ZW0gZXZlbnQgbmFtZXMuXHJcbiAqL1xyXG50eXBlIFN5c3RlbUV2ZW50TmFtZSA9IHtcclxuICAgIFtLIGluIGtleW9mICh0eXBlb2YgVHlwZXMpXTogKHR5cGVvZiBUeXBlcylbS11ba2V5b2YgKCh0eXBlb2YgVHlwZXMpW0tdKV1cclxufSBleHRlbmRzIChpbmZlciBNKSA/IE1ba2V5b2YgTV0gOiBuZXZlcjtcclxuXHJcbi8qKlxyXG4gKiBUaGUgZGF0YSB0eXBlIGFzc29jaWF0ZWQgdG8gYSBnaXZlbiBldmVudC5cclxuICovXHJcbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnREYXRhPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiA9XHJcbiAgICBFIGV4dGVuZHMga2V5b2YgQ3VzdG9tRXZlbnRzID8gQ3VzdG9tRXZlbnRzW0VdIDogKEUgZXh0ZW5kcyBTeXN0ZW1FdmVudE5hbWUgPyB2b2lkIDogYW55KTtcclxuXHJcbi8qKlxyXG4gKiBUaGUgdHlwZSBvZiBoYW5kbGVycyBmb3IgYSBnaXZlbiBldmVudC5cclxuICovXHJcbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnRDYWxsYmFjazxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4gPSAoZXY6IFdhaWxzRXZlbnQ8RT4pID0+IHZvaWQ7XHJcblxyXG4vKipcclxuICogUmVwcmVzZW50cyBhIHN5c3RlbSBldmVudCBvciBhIGN1c3RvbSBldmVudCBlbWl0dGVkIHRocm91Z2ggd2FpbHMtcHJvdmlkZWQgZmFjaWxpdGllcy5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiB7XHJcbiAgICAvKipcclxuICAgICAqIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICAgICAqL1xyXG4gICAgbmFtZTogRTtcclxuXHJcbiAgICAvKipcclxuICAgICAqIE9wdGlvbmFsIGRhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBlbWl0dGVkIGV2ZW50LlxyXG4gICAgICovXHJcbiAgICBkYXRhOiBXYWlsc0V2ZW50RGF0YTxFPjtcclxuXHJcbiAgICAvKipcclxuICAgICAqIE5hbWUgb2YgdGhlIG9yaWdpbmF0aW5nIHdpbmRvdy4gT21pdHRlZCBmb3IgYXBwbGljYXRpb24gZXZlbnRzLlxyXG4gICAgICogV2lsbCBiZSBvdmVycmlkZGVuIGlmIHNldCBtYW51YWxseS5cclxuICAgICAqL1xyXG4gICAgc2VuZGVyPzogc3RyaW5nO1xyXG5cclxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IEUsIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+KTtcclxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IFdhaWxzRXZlbnREYXRhPEU+IGV4dGVuZHMgbnVsbCB8IHZvaWQgPyBFIDogbmV2ZXIpXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBFLCBkYXRhPzogYW55KSB7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcclxuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhID8/IG51bGw7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudDogYW55KSB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xyXG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgbGV0IHdhaWxzRXZlbnQgPSBuZXcgV2FpbHNFdmVudChcclxuICAgICAgICBldmVudC5uYW1lLFxyXG4gICAgICAgIChldmVudC5uYW1lIGluIENyZWF0ZSkgPyBDcmVhdGVbZXZlbnQubmFtZV0oZXZlbnQuZGF0YSkgOiBldmVudC5kYXRhXHJcbiAgICApO1xyXG4gICAgaWYgKCdzZW5kZXInIGluIGV2ZW50KSB7XHJcbiAgICAgICAgd2FpbHNFdmVudC5zZW5kZXIgPSBldmVudC5zZW5kZXI7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gRGlzcGF0Y2ggdG8gYSBzbmFwc2hvdCwgdGhlbiByZW1vdmUgYWxsIGV4cGlyZWQgbGlzdGVuZXJzIGluIGEgc2luZ2xlXHJcbiAgICAvLyBwb3N0LWRpc3BhdGNoIGZpbHRlciBvZiB0aGUgbGl2ZSBtYXAuXHJcbiAgICAvLyAtIFdyaXRpbmcgdGhlIHNuYXBzaG90IGJhY2sgd2hvbGVzYWxlIHdvdWxkIHVuZG8gc3Vic2NyaXB0aW9uIGNoYW5nZXNcclxuICAgIC8vICAgbWFkZSBpbnNpZGUgYSBoYW5kbGVyICgjNDM5MykuXHJcbiAgICAvLyAtIENhbGxpbmcgbGlzdGVuZXJPZmYoKSBwZXIgZXhwaXJlZCBsaXN0ZW5lciBpbnNpZGUgdGhlIGxvb3AgaXMgTyhuwrIpXHJcbiAgICAvLyAgIHdoZW4gbWFueSBsaXN0ZW5lcnMgZXhwaXJlIG9uIHRoZSBzYW1lIGV2ZW50LlxyXG4gICAgLy8gRmlsdGVyaW5nIHRoZSBsaXZlIGFycmF5IG9uY2UgYWZ0ZXIgZGlzcGF0Y2ggaXMgTyhuKSBhbmQgc3RpbGwgaG9ub3Vyc1xyXG4gICAgLy8gYW55IGxpc3RlbmVycyBhZGRlZCBvciByZW1vdmVkIGJ5IGhhbmRsZXJzIGR1cmluZyBkaXNwYXRjaC5cclxuICAgIGNvbnN0IGV4cGlyZWQgPSBuZXcgU2V0PExpc3RlbmVyPigpO1xyXG4gICAgZm9yIChjb25zdCBsaXN0ZW5lciBvZiBsaXN0ZW5lcnMuc2xpY2UoKSkge1xyXG4gICAgICAgIGlmIChsaXN0ZW5lci5kaXNwYXRjaCh3YWlsc0V2ZW50KSkge1xyXG4gICAgICAgICAgICBleHBpcmVkLmFkZChsaXN0ZW5lcik7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG4gICAgaWYgKGV4cGlyZWQuc2l6ZSA+IDApIHtcclxuICAgICAgICBjb25zdCBsaXZlID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xyXG4gICAgICAgIGlmIChsaXZlKSB7XHJcbiAgICAgICAgICAgIGNvbnN0IHJlbWFpbmluZyA9IGxpdmUuZmlsdGVyKGwgPT4gIWV4cGlyZWQuaGFzKGwpKTtcclxuICAgICAgICAgICAgaWYgKHJlbWFpbmluZy5sZW5ndGggPT09IDApIHtcclxuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcclxuICAgICAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCByZW1haW5pbmcpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXHJcbiAqIEBwYXJhbSBtYXhDYWxsYmFja3MgLSBUaGUgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdGhlIGNhbGxiYWNrIGNhbiBiZSBjYWxsZWQgZm9yIHRoZSBldmVudC4gT25jZSB0aGUgbWF4aW11bSBudW1iZXIgaXMgcmVhY2hlZCwgdGhlIGNhbGxiYWNrIHdpbGwgbm8gbG9uZ2VyIGJlIGNhbGxlZC5cclxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihldmVudE5hbWU6IEUsIGNhbGxiYWNrOiBXYWlsc0V2ZW50Q2FsbGJhY2s8RT4sIG1heENhbGxiYWNrczogbnVtYmVyKSB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XHJcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcclxuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XHJcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xyXG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCB3aGVuIHRoZSBzcGVjaWZpZWQgZXZlbnQgb2NjdXJzLlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXHJcbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxyXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihldmVudE5hbWU6IEUsIGNhbGxiYWNrOiBXYWlsc0V2ZW50Q2FsbGJhY2s8RT4pOiAoKSA9PiB2b2lkIHtcclxuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxyXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbmNlPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihldmVudE5hbWU6IEUsIGNhbGxiYWNrOiBXYWlsc0V2ZW50Q2FsbGJhY2s8RT4pOiAoKSA9PiB2b2lkIHtcclxuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXHJcbiAqXHJcbiAqIEBwYXJhbSBldmVudE5hbWVzIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50cyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmYoLi4uZXZlbnROYW1lczogW1dhaWxzRXZlbnROYW1lLCAuLi5XYWlsc0V2ZW50TmFtZVtdXSk6IHZvaWQge1xyXG4gICAgZXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCk6IHZvaWQge1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEVtaXRzIGFuIGV2ZW50LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCB3aWxsIGJlIGZ1bGZpbGxlZCBvbmNlIHRoZSBldmVudCBoYXMgYmVlbiBlbWl0dGVkLiAgUmVzb2x2ZXMgdG8gdHJ1ZSBpZiB0aGUgZXZlbnQgd2FzIGNhbmNlbGxlZC5cclxuICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdFxyXG4gKiBAcGFyYW0gZGF0YSAtIFRoZSBkYXRhIHRoYXQgd2lsbCBiZSBzZW50IHdpdGggdGhlIGV2ZW50XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogRSwgZGF0YTogV2FpbHNFdmVudERhdGE8RT4pOiBQcm9taXNlPGJvb2xlYW4+XHJcbmV4cG9ydCBmdW5jdGlvbiBFbWl0PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiBleHRlbmRzIG51bGwgfCB2b2lkID8gRSA6IG5ldmVyKTogUHJvbWlzZTxib29sZWFuPlxyXG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogV2FpbHNFdmVudERhdGE8RT4sIGRhdGE/OiBhbnkpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgIHJldHVybiBjYWxsKEVtaXRNZXRob2QsICBuZXcgV2FpbHNFdmVudChuYW1lLCBkYXRhKSlcclxufVxyXG5cclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuZXhwb3J0IGludGVyZmFjZSBTaXplIHtcclxuICAgIC8qKiBUaGUgd2lkdGggb2YgYSByZWN0YW5ndWxhciBhcmVhLiAqL1xyXG4gICAgV2lkdGg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cclxuICAgIEhlaWdodDogbnVtYmVyO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFJlY3Qge1xyXG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cclxuICAgIFg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uICovXHJcbiAgICBZOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuICovXHJcbiAgICBXaWR0aDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHJlY3RhbmdsZS4gKi9cclxuICAgIEhlaWdodDogbnVtYmVyO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFNjcmVlbiB7XHJcbiAgICAvKiogVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBzY3JlZW4uICovXHJcbiAgICBJRDogc3RyaW5nO1xyXG4gICAgLyoqIEh1bWFuLXJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIE5hbWU6IHN0cmluZztcclxuICAgIC8qKiBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLiAqL1xyXG4gICAgU2NhbGVGYWN0b3I6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgWCBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBYOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgWTogbnVtYmVyO1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBTaXplOiBTaXplO1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC4gKi9cclxuICAgIEJvdW5kczogUmVjdDtcclxuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQgKGJlZm9yZSBzY2FsaW5nKS4gKi9cclxuICAgIFBoeXNpY2FsQm91bmRzOiBSZWN0O1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLiAqL1xyXG4gICAgV29ya0FyZWE6IFJlY3Q7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIFdvcmtBcmVhIG9mIHRoZSBzY3JlZW4gKGJlZm9yZSBzY2FsaW5nKS4gKi9cclxuICAgIFBoeXNpY2FsV29ya0FyZWE6IFJlY3Q7XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXHJcbiAgICBJc1ByaW1hcnk6IGJvb2xlYW47XHJcbiAgICAvKiogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBSb3RhdGlvbjogbnVtYmVyO1xyXG59XHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuU2NyZWVucyk7XHJcblxyXG5jb25zdCBnZXRBbGwgPSAwO1xyXG5jb25zdCBnZXRQcmltYXJ5ID0gMTtcclxuY29uc3QgZ2V0Q3VycmVudCA9IDI7XHJcbmNvbnN0IGdldEJ5SUQgPSAzO1xyXG5jb25zdCBnZXRCeUluZGV4ID0gNDtcclxuXHJcbi8qKlxyXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKTogUHJvbWlzZTxTY3JlZW5bXT4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCk6IFByb21pc2U8U2NyZWVuPiB7XHJcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgYSBzY3JlZW4gYnkgaXRzIHVuaXF1ZSBkaXNwbGF5IElELlxyXG4gKlxyXG4gKiBAcGFyYW0gaWQgLSBUaGUgdW5pcXVlIGlkZW50aWZpZXIgb2YgdGhlIHNjcmVlbi5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIG1hdGNoaW5nIFNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRCeUlEKGlkOiBzdHJpbmcpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0QnlJRCwgeyBpZCB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgYSBzY3JlZW4gYnkgaXRzIGluZGV4IGluIHRoZSBzY3JlZW4gbGlzdC5cclxuICpcclxuICogQHBhcmFtIGluZGV4IC0gVGhlIHplcm8tYmFzZWQgaW5kZXggb2YgdGhlIHNjcmVlbi5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIG1hdGNoaW5nIFNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRCeUluZGV4KGluZGV4OiBudW1iZXIpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0QnlJbmRleCwgeyBpbmRleCB9KTtcclxufVxyXG4iLCIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLklPUyk7XHJcblxyXG4vLyBNZXRob2QgSURzXHJcbmNvbnN0IEhhcHRpY3NJbXBhY3QgPSAwO1xyXG5jb25zdCBEZXZpY2VJbmZvID0gMTtcclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgSGFwdGljcyB7XHJcbiAgICBleHBvcnQgdHlwZSBJbXBhY3RTdHlsZSA9IFwibGlnaHRcInxcIm1lZGl1bVwifFwiaGVhdnlcInxcInNvZnRcInxcInJpZ2lkXCI7XHJcbiAgICBleHBvcnQgZnVuY3Rpb24gSW1wYWN0KHN0eWxlOiBJbXBhY3RTdHlsZSA9IFwibWVkaXVtXCIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gY2FsbChIYXB0aWNzSW1wYWN0LCB7IHN0eWxlIH0pO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgbmFtZXNwYWNlIERldmljZSB7XHJcbiAgICBleHBvcnQgaW50ZXJmYWNlIEluZm8ge1xyXG4gICAgICAgIG1vZGVsOiBzdHJpbmc7XHJcbiAgICAgICAgc3lzdGVtTmFtZTogc3RyaW5nO1xyXG4gICAgICAgIHN5c3RlbVZlcnNpb246IHN0cmluZztcclxuICAgICAgICBpc1NpbXVsYXRvcjogYm9vbGVhbjtcclxuICAgIH1cclxuICAgIGV4cG9ydCBmdW5jdGlvbiBJbmZvKCk6IFByb21pc2U8SW5mbz4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKERldmljZUluZm8pO1xyXG4gICAgfVxyXG59XHJcbiIsIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQW5kcm9pZCk7XHJcblxyXG4vLyBNZXRob2QgSURzIChtdXN0IG1hdGNoIG1lc3NhZ2Vwcm9jZXNzb3JfYW5kcm9pZC5nbylcclxuY29uc3QgSGFwdGljc1ZpYnJhdGUgPSAwO1xyXG5jb25zdCBEZXZpY2VJbmZvID0gMTtcclxuY29uc3QgVG9hc3RTaG93ID0gMjtcclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgSGFwdGljcyB7XHJcbiAgICAvKiogVmlicmF0ZSB0aGUgZGV2aWNlIGZvciB0aGUgZ2l2ZW4gZHVyYXRpb24gaW4gbWlsbGlzZWNvbmRzLiAqL1xyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFZpYnJhdGUoZHVyYXRpb25NczogbnVtYmVyID0gMTAwKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIGNhbGwoSGFwdGljc1ZpYnJhdGUsIHsgZHVyYXRpb246IGR1cmF0aW9uTXMgfSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgRGV2aWNlIHtcclxuICAgIGV4cG9ydCBpbnRlcmZhY2UgSW5mbyB7XHJcbiAgICAgICAgcGxhdGZvcm06IHN0cmluZztcclxuICAgICAgICBtYW51ZmFjdHVyZXI6IHN0cmluZztcclxuICAgICAgICBicmFuZDogc3RyaW5nO1xyXG4gICAgICAgIG1vZGVsOiBzdHJpbmc7XHJcbiAgICAgICAgZGV2aWNlOiBzdHJpbmc7XHJcbiAgICAgICAgdmVyc2lvbjogc3RyaW5nO1xyXG4gICAgICAgIHNka0ludDogbnVtYmVyO1xyXG4gICAgfVxyXG4gICAgLyoqIFJldHVybiBpbmZvcm1hdGlvbiBhYm91dCB0aGUgQW5kcm9pZCBkZXZpY2UuICovXHJcbiAgICBleHBvcnQgZnVuY3Rpb24gSW5mbygpOiBQcm9taXNlPEluZm8+IHtcclxuICAgICAgICByZXR1cm4gY2FsbChEZXZpY2VJbmZvKTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IG5hbWVzcGFjZSBUb2FzdCB7XHJcbiAgICAvKiogU2hvdyBhIHNob3J0IEFuZHJvaWQgdG9hc3QgbWVzc2FnZS4gKi9cclxuICAgIGV4cG9ydCBmdW5jdGlvbiBTaG93KG1lc3NhZ2U6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKFRvYXN0U2hvdywgeyBtZXNzYWdlIH0pO1xyXG4gICAgfVxyXG59XHJcbiIsIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBVcGRhdGVyIGV2ZW50IG5hbWUgY29uc3RhbnRzLlxyXG4gKlxyXG4gKiBVc2UgdGhlc2UgaW5zdGVhZCBvZiBoYXJkLWNvZGluZyBzdHJpbmcgbGl0ZXJhbHMgd2hlbiBzdWJzY3JpYmluZyB0b1xyXG4gKiB1cGRhdGVyIGV2ZW50cyBmcm9tIEphdmFTY3JpcHQ6XHJcbiAqXHJcbiAqICAgICBpbXBvcnQgeyBFdmVudHMsIFVwZGF0ZXIgfSBmcm9tIFwiQHdhaWxzaW8vcnVudGltZVwiO1xyXG4gKlxyXG4gKiAgICAgRXZlbnRzLk9uKFVwZGF0ZXIuRXZlbnRzLlVwZGF0ZUF2YWlsYWJsZSwgKGUpID0+IHtcclxuICogICAgICAgICBjb25zb2xlLmxvZyhcInVwZGF0ZSBmb3VuZDpcIiwgZS5kYXRhLnZlcnNpb24pO1xyXG4gKiAgICAgfSk7XHJcbiAqXHJcbiAqICAgICBFdmVudHMuT24oVXBkYXRlci5FdmVudHMuRG93bmxvYWRQcm9ncmVzcywgKGUpID0+IHtcclxuICogICAgICAgICBjb25zdCBwID0gZS5kYXRhO1xyXG4gKiAgICAgICAgIGNvbnNvbGUubG9nKGAke3Aud3JpdHRlbn0gLyAke3AudG90YWx9IGJ5dGVzYCk7XHJcbiAqICAgICB9KTtcclxuICpcclxuICogTWlycm9ycyB0aGUgR28tc2lkZSBjb25zdGFudHMgaW4gYHBrZy91cGRhdGVyL2V2ZW50cy5nb2AgYW5kIHRoZVxyXG4gKiB1c2VyLWFjdGlvbiBjb25zdGFudHMgaW4gYHBrZy91cGRhdGVyL3dpbmRvd19saWZlY3ljbGUuZ29gLiBBbnlcclxuICogY2hhbmdlcyBoZXJlIG11c3Qgc3RheSBpbiBzeW5jIHdpdGggdGhvc2UgZmlsZXMg4oCUIHRoZXJlJ3MgYW5cclxuICogaW50ZWdyYXRpb24gdGVzdCB0aGF0IGFzc2VydHMgdGhlIHN0cmluZ3MgbWF0Y2guXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgRXZlbnRzID0gT2JqZWN0LmZyZWV6ZSh7XHJcbiAgICAvKiogQSBDaGVjayByb3VuZC10cmlwIGlzIHN0YXJ0aW5nLiBQYXlsb2FkOiBudWxsLiAqL1xyXG4gICAgQ2hlY2tTdGFydGVkOiBcIndhaWxzOnVwZGF0ZXI6Y2hlY2stc3RhcnRlZFwiLFxyXG4gICAgLyoqIENoZWNrIGZvdW5kIGEgbmV3ZXIgcmVsZWFzZS4gUGF5bG9hZDogUmVsZWFzZS4gKi9cclxuICAgIFVwZGF0ZUF2YWlsYWJsZTogXCJ3YWlsczp1cGRhdGVyOnVwZGF0ZS1hdmFpbGFibGVcIixcclxuICAgIC8qKiBDaGVjayBjb25maXJtZWQgdGhlIGNhbGxlciBpcyB1cCB0byBkYXRlLiBQYXlsb2FkOiBudWxsLiAqL1xyXG4gICAgTm9VcGRhdGU6IFwid2FpbHM6dXBkYXRlcjpuby11cGRhdGVcIixcclxuICAgIC8qKiBEb3dubG9hZCBpcyBzdGFydGluZy4gUGF5bG9hZDogUmVsZWFzZS4gKi9cclxuICAgIERvd25sb2FkU3RhcnRlZDogXCJ3YWlsczp1cGRhdGVyOmRvd25sb2FkLXN0YXJ0ZWRcIixcclxuICAgIC8qKiBQZXJpb2RpYyBwcm9ncmVzcyB0aWNrIGR1cmluZyBkb3dubG9hZCAofjEwIEh6KS4gUGF5bG9hZDogUHJvZ3Jlc3MuICovXHJcbiAgICBEb3dubG9hZFByb2dyZXNzOiBcIndhaWxzOnVwZGF0ZXI6ZG93bmxvYWQtcHJvZ3Jlc3NcIixcclxuICAgIC8qKiBBbGwgYnl0ZXMgYXJlIG9uIGRpc2ssIGJ1dCB2ZXJpZmljYXRpb24gaGFzIG5vdCB5ZXQgc3RhcnRlZC4gUGF5bG9hZDogUmVsZWFzZS4gKi9cclxuICAgIERvd25sb2FkQ29tcGxldGU6IFwid2FpbHM6dXBkYXRlcjpkb3dubG9hZC1jb21wbGV0ZVwiLFxyXG4gICAgLyoqIFNpZ25hdHVyZSAvIGRpZ2VzdCB2ZXJpZmljYXRpb24gaGFzIHN0YXJ0ZWQuIFBheWxvYWQ6IFJlbGVhc2UuICovXHJcbiAgICBWZXJpZnlpbmc6IFwid2FpbHM6dXBkYXRlcjp2ZXJpZnlpbmdcIixcclxuICAgIC8qKiBUaGUgVXBkYXRlciBpcyBzd2FwcGluZyB0aGUgYmluYXJ5IGludG8gcGxhY2UuIFBheWxvYWQ6IFJlbGVhc2UuICovXHJcbiAgICBJbnN0YWxsaW5nOiBcIndhaWxzOnVwZGF0ZXI6aW5zdGFsbGluZ1wiLFxyXG4gICAgLyoqIFVwZGF0ZSBpcyBzdGFnZWQgYW5kIGEgcmVzdGFydCBpcyBwZW5kaW5nLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xyXG4gICAgVXBkYXRlUmVhZHk6IFwid2FpbHM6dXBkYXRlcjp1cGRhdGUtcmVhZHlcIixcclxuICAgIC8qKiBTb21ldGhpbmcgZmFpbGVkLiBQYXlsb2FkOiBFcnJvckluZm8geyBzdGFnZSwgbWVzc2FnZSwgcHJvdmlkZXIgfS4gKi9cclxuICAgIEVycm9yOiBcIndhaWxzOnVwZGF0ZXI6ZXJyb3JcIixcclxuICAgIC8qKiBIb3N0LXNpZGUgY29udGV4dCBkZWxpdmVyZWQgb25jZSBwZXIgc2Vzc2lvbi4gUGF5bG9hZDogTWV0YSB7IGN1cnJlbnRWZXJzaW9uLCBza2lwcGVkVmVyc2lvbiB9LiAqL1xyXG4gICAgTWV0YTogXCJ3YWlsczp1cGRhdGVyOm1ldGFcIixcclxuXHJcbiAgICAvKiogU3ViLW5hbWVzcGFjZTogdXNlci1hY3Rpb24gZXZlbnRzIHRoYXQgdGhlIFVJIGVtaXRzIEJBQ0sgdG8gdGhlIGhvc3QuICovXHJcbiAgICBVc2VyOiBPYmplY3QuZnJlZXplKHtcclxuICAgICAgICAvKiogVXNlciBjbGlja2VkIEluc3RhbGwgb24gYW4gQXZhaWxhYmxlIHVwZGF0ZS4gKi9cclxuICAgICAgICBJbnN0YWxsOiBcIndhaWxzOnVwZGF0ZXI6dXNlcjppbnN0YWxsXCIsXHJcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBSZXN0YXJ0ICYgQXBwbHkgb24gYSBSZWFkeSB1cGRhdGUuICovXHJcbiAgICAgICAgUmVzdGFydDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6cmVzdGFydFwiLFxyXG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgU2tpcCBUaGlzIFZlcnNpb24uICovXHJcbiAgICAgICAgU2tpcDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6c2tpcFwiLFxyXG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgUmVtaW5kIE1lIExhdGVyLiAqL1xyXG4gICAgICAgIFJlbWluZDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6cmVtaW5kXCIsXHJcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBDbG9zZSAvIENhbmNlbC4gKi9cclxuICAgICAgICBDYW5jZWw6IFwid2FpbHM6dXBkYXRlcjp1c2VyOmNhbmNlbFwiLFxyXG4gICAgfSksXHJcblxyXG4gICAgLyoqIFN1Yi1uYW1lc3BhY2U6IGZyYW1ld29yay1pbnRlcm5hbCBldmVudHMgdGhlIFVJIGVtaXRzIHRvIGNvb3JkaW5hdGVcclxuICAgICAqICB3aXRoIHRoZSBob3N0LiBNb3N0IGFwcCBjb2RlIGNhbiBpZ25vcmUgdGhlc2UuICovXHJcbiAgICBXaW5kb3c6IE9iamVjdC5mcmVlemUoe1xyXG4gICAgICAgIC8qKiBUaGUgd2luZG93IGZpbmlzaGVkIGxvYWRpbmcgYW5kIGFza3MgdGhlIGhvc3QgdG8gcmVwbGF5IGN1cnJlbnQgc3RhdGUuICovXHJcbiAgICAgICAgUmVhZHk6IFwid2FpbHM6dXBkYXRlcjp3aW5kb3c6cmVhZHlcIixcclxuICAgIH0pLFxyXG59KTtcclxuIiwiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5pbXBvcnQgdHlwZSB7IFNjcmVlbiB9IGZyb20gXCIuL3NjcmVlbnMuanNcIjtcclxuXHJcbi8vIERyb3AgdGFyZ2V0IGNvbnN0YW50c1xyXG5jb25zdCBEUk9QX1RBUkdFVF9BVFRSSUJVVEUgPSAnZGF0YS1maWxlLWRyb3AtdGFyZ2V0JztcclxuY29uc3QgRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTID0gJ2ZpbGUtZHJvcC10YXJnZXQtYWN0aXZlJztcclxubGV0IGN1cnJlbnREcm9wVGFyZ2V0OiBFbGVtZW50IHwgbnVsbCA9IG51bGw7XHJcblxyXG5jb25zdCBQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAwO1xyXG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xyXG5jb25zdCBDbG9zZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSAyO1xyXG5jb25zdCBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgPSAzO1xyXG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xyXG5jb25zdCBGb2N1c01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1O1xyXG5jb25zdCBGb3JjZVJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgPSA2O1xyXG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xyXG5jb25zdCBHZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA4O1xyXG5jb25zdCBHZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA5O1xyXG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcclxuY29uc3QgSGlkZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTE7XHJcbmNvbnN0IElzRm9jdXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDEyO1xyXG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcclxuY29uc3QgSXNNYXhpbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTQ7XHJcbmNvbnN0IElzTWluaW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE1O1xyXG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcclxuY29uc3QgTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTc7XHJcbmNvbnN0IE5hbWVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDE4O1xyXG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcclxuY29uc3QgUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMjA7XHJcbmNvbnN0IFJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDIxO1xyXG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcclxuY29uc3QgUmVzdG9yZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMjM7XHJcbmNvbnN0IFNldFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICA9IDI0O1xyXG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcclxuY29uc3QgU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCAgICAgICAgID0gMjY7XHJcbmNvbnN0IFNldEZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICAgICA9IDI3O1xyXG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcclxuY29uc3QgU2V0TWF4U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMjk7XHJcbmNvbnN0IFNldE1pblNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDMwO1xyXG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcclxuY29uc3QgU2V0UmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgID0gMzI7XHJcbmNvbnN0IFNldFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDMzO1xyXG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcclxuY29uc3QgU2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzU7XHJcbmNvbnN0IFNob3dNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM2O1xyXG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcclxuY29uc3QgVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgID0gMzg7XHJcbmNvbnN0IFRvZ2dsZU1heGltaXNlTWV0aG9kICAgICAgICAgICAgICA9IDM5O1xyXG5jb25zdCBUb2dnbGVGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgPSA0MDsgXHJcbmNvbnN0IFVuRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDQxO1xyXG5jb25zdCBVbk1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MjtcclxuY29uc3QgVW5NaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDM7XHJcbmNvbnN0IFdpZHRoTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDQ0O1xyXG5jb25zdCBab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSA0NTtcclxuY29uc3QgWm9vbUluTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gNDY7XHJcbmNvbnN0IFpvb21PdXRNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDQ3O1xyXG5jb25zdCBab29tUmVzZXRNZXRob2QgICAgICAgICAgICAgICAgICAgPSA0ODtcclxuY29uc3QgU25hcEFzc2lzdE1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDk7XHJcbmNvbnN0IEZpbGVzRHJvcHBlZCAgICAgICAgICAgICAgICAgICAgICA9IDUwO1xyXG5jb25zdCBQcmludE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1MTtcclxuY29uc3QgU2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gNTI7XHJcblxyXG4vKipcclxuICogRmluZHMgdGhlIG5lYXJlc3QgZHJvcCB0YXJnZXQgZWxlbWVudCBieSB3YWxraW5nIHVwIHRoZSBET00gdHJlZS5cclxuICovXHJcbmZ1bmN0aW9uIGdldERyb3BUYXJnZXRFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQgfCBudWxsKTogRWxlbWVudCB8IG51bGwge1xyXG4gICAgaWYgKCFlbGVtZW50KSB7XHJcbiAgICAgICAgcmV0dXJuIG51bGw7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gZWxlbWVudC5jbG9zZXN0KGBbJHtEUk9QX1RBUkdFVF9BVFRSSUJVVEV9XWApO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2sgaWYgd2UgY2FuIHVzZSBXZWJWaWV3MidzIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIChXaW5kb3dzKVxyXG4gKiBBbHNvIGNoZWNrcyB0aGF0IEVuYWJsZUZpbGVEcm9wIGlzIHRydWUgZm9yIHRoaXMgd2luZG93LlxyXG4gKi9cclxuZnVuY3Rpb24gY2FuUmVzb2x2ZUZpbGVQYXRocygpOiBib29sZWFuIHtcclxuICAgIC8vIE11c3QgaGF2ZSBXZWJWaWV3MidzIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIEFQSSAoV2luZG93cyBvbmx5KVxyXG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzID09IG51bGwpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcbiAgICAvLyBNdXN0IGhhdmUgRW5hYmxlRmlsZURyb3Agc2V0IHRvIHRydWUgZm9yIHRoaXMgd2luZG93XHJcbiAgICAvLyBUaGlzIGZsYWcgaXMgc2V0IGJ5IHRoZSBHbyBiYWNrZW5kIGR1cmluZyBydW50aW1lIGluaXRpYWxpemF0aW9uXHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSB0cnVlO1xyXG59XHJcblxyXG4vKipcclxuICogU2VuZCBmaWxlIGRyb3AgdG8gYmFja2VuZCB2aWEgV2ViVmlldzIgKFdpbmRvd3Mgb25seSlcclxuICovXHJcbmZ1bmN0aW9uIHJlc29sdmVGaWxlUGF0aHMoeDogbnVtYmVyLCB5OiBudW1iZXIsIGZpbGVzOiBGaWxlW10pOiB2b2lkIHtcclxuICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cykge1xyXG4gICAgICAgICh3aW5kb3cgYXMgYW55KS5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyhgZmlsZTpkcm9wOiR7eH06JHt5fWAsIGZpbGVzKTtcclxuICAgIH1cclxufVxyXG5cclxuLy8gTmF0aXZlIGRyYWcgc3RhdGUgKExpbnV4L21hY09TIGludGVyY2VwdCBET00gZHJhZyBldmVudHMpXHJcbmxldCBuYXRpdmVEcmFnQWN0aXZlID0gZmFsc2U7XHJcblxyXG4vKipcclxuICogQ2xlYW5zIHVwIG5hdGl2ZSBkcmFnIHN0YXRlIGFuZCBob3ZlciBlZmZlY3RzLlxyXG4gKiBDYWxsZWQgb24gZHJvcCBvciB3aGVuIGRyYWcgbGVhdmVzIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5mdW5jdGlvbiBjbGVhbnVwTmF0aXZlRHJhZygpOiB2b2lkIHtcclxuICAgIG5hdGl2ZURyYWdBY3RpdmUgPSBmYWxzZTtcclxuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCkge1xyXG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsZWQgZnJvbSBHbyB3aGVuIGEgZmlsZSBkcmFnIGVudGVycyB0aGUgd2luZG93IG9uIExpbnV4L21hY09TLlxyXG4gKi9cclxuZnVuY3Rpb24gaGFuZGxlRHJhZ0VudGVyKCk6IHZvaWQge1xyXG4gICAgLy8gQ2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XHJcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBhY3RpdmF0ZSBkcmFnIHN0YXRlXHJcbiAgICB9XHJcbiAgICBuYXRpdmVEcmFnQWN0aXZlID0gdHJ1ZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxlZCBmcm9tIEdvIHdoZW4gYSBmaWxlIGRyYWcgbGVhdmVzIHRoZSB3aW5kb3cgb24gTGludXgvbWFjT1MuXHJcbiAqL1xyXG5mdW5jdGlvbiBoYW5kbGVEcmFnTGVhdmUoKTogdm9pZCB7XHJcbiAgICBjbGVhbnVwTmF0aXZlRHJhZygpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbGVkIGZyb20gR28gZHVyaW5nIGZpbGUgZHJhZyB0byB1cGRhdGUgaG92ZXIgc3RhdGUgb24gTGludXgvbWFjT1MuXHJcbiAqIEBwYXJhbSB4IC0gWCBjb29yZGluYXRlIGluIENTUyBwaXhlbHNcclxuICogQHBhcmFtIHkgLSBZIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xyXG4gKi9cclxuZnVuY3Rpb24gaGFuZGxlRHJhZ092ZXIoeDogbnVtYmVyLCB5OiBudW1iZXIpOiB2b2lkIHtcclxuICAgIGlmICghbmF0aXZlRHJhZ0FjdGl2ZSkgcmV0dXJuO1xyXG4gICAgXHJcbiAgICAvLyBDaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xyXG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xyXG4gICAgfVxyXG4gICAgXHJcbiAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KTtcclxuICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudCh0YXJnZXRFbGVtZW50KTtcclxuICAgIFxyXG4gICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0ICYmIGN1cnJlbnREcm9wVGFyZ2V0ICE9PSBkcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgfVxyXG4gICAgXHJcbiAgICBpZiAoZHJvcFRhcmdldCkge1xyXG4gICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xyXG4gICAgfVxyXG59XHJcblxyXG5cclxuXHJcbi8vIEV4cG9ydCB0aGUgaGFuZGxlcnMgZm9yIHVzZSBieSBHbyB2aWEgaW5kZXgudHNcclxuZXhwb3J0IHsgaGFuZGxlRHJhZ0VudGVyLCBoYW5kbGVEcmFnTGVhdmUsIGhhbmRsZURyYWdPdmVyIH07XHJcblxyXG4vKipcclxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgcG9zaXRpb24gb2YgYSB3aW5kb3cuXHJcbiAqL1xyXG5pbnRlcmZhY2UgUG9zaXRpb24ge1xyXG4gICAgLyoqIFRoZSBob3Jpem9udGFsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICB4OiBudW1iZXI7XHJcbiAgICAvKiogVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICB5OiBudW1iZXI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBzaXplIG9mIGEgd2luZG93LlxyXG4gKi9cclxuaW50ZXJmYWNlIFNpemUge1xyXG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgd2luZG93LiAqL1xyXG4gICAgd2lkdGg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICBoZWlnaHQ6IG51bWJlcjtcclxufVxyXG5cclxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cclxuY29uc3QgY2FsbGVyU3ltID0gU3ltYm9sKFwiY2FsbGVyXCIpO1xyXG5cclxuY2xhc3MgV2luZG93IHtcclxuICAgIC8vIFByaXZhdGUgZmllbGRzLlxyXG4gICAgcHJpdmF0ZSBbY2FsbGVyU3ltXTogKG1lc3NhZ2U6IG51bWJlciwgYXJncz86IGFueSkgPT4gUHJvbWlzZTxhbnk+O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogSW5pdGlhbGlzZXMgYSB3aW5kb3cgb2JqZWN0IHdpdGggdGhlIHNwZWNpZmllZCBuYW1lLlxyXG4gICAgICpcclxuICAgICAqIEBwcml2YXRlXHJcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB0YXJnZXQgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcgPSAnJykge1xyXG4gICAgICAgIHRoaXNbY2FsbGVyU3ltXSA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuV2luZG93LCBuYW1lKVxyXG5cclxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xyXG4gICAgICAgIGZvciAoY29uc3QgbWV0aG9kIG9mIE9iamVjdC5nZXRPd25Qcm9wZXJ0eU5hbWVzKFdpbmRvdy5wcm90b3R5cGUpKSB7XHJcbiAgICAgICAgICAgIGlmIChcclxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXHJcbiAgICAgICAgICAgICAgICAmJiB0eXBlb2YgKHRoaXMgYXMgYW55KVttZXRob2RdID09PSBcImZ1bmN0aW9uXCJcclxuICAgICAgICAgICAgKSB7XHJcbiAgICAgICAgICAgICAgICAodGhpcyBhcyBhbnkpW21ldGhvZF0gPSAodGhpcyBhcyBhbnkpW21ldGhvZF0uYmluZCh0aGlzKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cclxuICAgICAqIEByZXR1cm5zIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXHJcbiAgICAgKi9cclxuICAgIEdldChuYW1lOiBzdHJpbmcpOiBXaW5kb3cge1xyXG4gICAgICAgIHJldHVybiBuZXcgV2luZG93KG5hbWUpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQb3NpdGlvbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cclxuICAgICAqL1xyXG4gICAgQ2VudGVyKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2VudGVyTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENsb3NlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBDbG9zZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENsb3NlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cclxuICAgICAqL1xyXG4gICAgRGlzYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRW5hYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXHJcbiAgICAgKi9cclxuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBGb2N1c2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIEZvY3VzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRm9jdXNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cclxuICAgICAqL1xyXG4gICAgRm9yY2VSZWxvYWQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb3JjZVJlbG9hZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTd2l0Y2hlcyB0aGUgd2luZG93IHRvIGZ1bGxzY3JlZW4gbW9kZS5cclxuICAgICAqL1xyXG4gICAgRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgc2NyZWVuIHRoYXQgdGhlIHdpbmRvdyBpcyBvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgc2NyZWVuIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG9uLlxyXG4gICAgICovXHJcbiAgICBHZXRTY3JlZW4oKTogUHJvbWlzZTxTY3JlZW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBjdXJyZW50IHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB6b29tIGxldmVsLlxyXG4gICAgICovXHJcbiAgICBHZXRab29tKCk6IFByb21pc2U8bnVtYmVyPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRab29tTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBIZWlnaHQoKTogUHJvbWlzZTxudW1iZXI+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhlaWdodE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBIaWRlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSGlkZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZm9jdXNlZC5cclxuICAgICAqL1xyXG4gICAgSXNGb2N1c2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNGb2N1c2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmdWxsc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBJc0Z1bGxzY3JlZW4oKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0Z1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWF4aW1pc2VkLlxyXG4gICAgICovXHJcbiAgICBJc01heGltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWF4aW1pc2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1pbmltaXNlZC5cclxuICAgICAqL1xyXG4gICAgSXNNaW5pbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc01pbmltaXNlZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNYXhpbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBNaW5pbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNaW5pbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgTmFtZSgpOiBQcm9taXNlPHN0cmluZz4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTmFtZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBPcGVucyB0aGUgZGV2ZWxvcG1lbnQgdG9vbHMgcGFuZS5cclxuICAgICAqL1xyXG4gICAgT3BlbkRldlRvb2xzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oT3BlbkRldlRvb2xzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBSZWxhdGl2ZVBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbGF0aXZlUG9zaXRpb25NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXHJcbiAgICAgKi9cclxuICAgIFJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbG9hZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGUuXHJcbiAgICAgKi9cclxuICAgIFJlc2l6YWJsZSgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc2l6YWJsZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXN0b3JlcyB0aGUgd2luZG93IHRvIGl0cyBwcmV2aW91cyBzdGF0ZSBpZiBpdCB3YXMgcHJldmlvdXNseSBtaW5pbWlzZWQsIG1heGltaXNlZCBvciBmdWxsc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBSZXN0b3JlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzdG9yZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBhbHdheXNPblRvcCAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgc3RheSBvbiB0b3AuXHJcbiAgICAgKi9cclxuICAgIFNldEFsd2F5c09uVG9wKGFsd2F5c09uVG9wOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRBbHdheXNPblRvcE1ldGhvZCwgeyBhbHdheXNPblRvcCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHIgLSBUaGUgZGVzaXJlZCByZWQgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cclxuICAgICAqIEBwYXJhbSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cclxuICAgICAqIEBwYXJhbSBiIC0gVGhlIGRlc2lyZWQgYmx1ZSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICogQHBhcmFtIGEgLSBUaGUgZGVzaXJlZCBhbHBoYSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICovXHJcbiAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyKHI6IG51bWJlciwgZzogbnVtYmVyLCBiOiBudW1iZXIsIGE6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCwgeyByLCBnLCBiLCBhIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGZyYW1lbGVzcyAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgZnJhbWVsZXNzLlxyXG4gICAgICovXHJcbiAgICBTZXRGcmFtZWxlc3MoZnJhbWVsZXNzOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRGcmFtZWxlc3NNZXRob2QsIHsgZnJhbWVsZXNzIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZW5hYmxlZCAtIFdoZXRoZXIgdGhlIGZ1bGxzY3JlZW4gYnV0dG9uIHNob3VsZCBiZSBlbmFibGVkLlxyXG4gICAgICovXHJcbiAgICBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZChlbmFibGVkOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCwgeyBlbmFibGVkIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWF4aW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldE1heFNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtaW5pbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRNaW5TaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRNaW5TaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFJlbGF0aXZlUG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gcmVzaXphYmxlIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSByZXNpemFibGUuXHJcbiAgICAgKi9cclxuICAgIFNldFJlc2l6YWJsZShyZXNpemFibGU6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFJlc2l6YWJsZU1ldGhvZCwgeyByZXNpemFibGUgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgdGl0bGUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gdGl0bGUgLSBUaGUgZGVzaXJlZCB0aXRsZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRUaXRsZSh0aXRsZTogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWwuXHJcbiAgICAgKi9cclxuICAgIFNldFpvb20oem9vbTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTaG93cyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2hvd01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTaXplKCk6IFByb21pc2U8U2l6ZT4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2l6ZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmdWxsc2NyZWVuIGFuZCBub3JtYWwuXHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZUZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVGdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxyXG4gICAgICovXHJcbiAgICBUb2dnbGVNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZyYW1lbGVzcyBhbmQgbm9ybWFsLlxyXG4gICAgICovXHJcbiAgICBUb2dnbGVGcmFtZWxlc3MoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVGcmFtZWxlc3NNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tZnVsbHNjcmVlbnMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgVW5GdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5GdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFVuLW1heGltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBVbk1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NYXhpbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBVbi1taW5pbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgVW5NaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuTWluaW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBXaWR0aCgpOiBQcm9taXNlPG51bWJlcj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oV2lkdGhNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogWm9vbXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgWm9vbSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXHJcbiAgICAgKi9cclxuICAgIFpvb21JbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21Jbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEZWNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cclxuICAgICAqL1xyXG4gICAgWm9vbU91dCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21PdXRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVzZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXHJcbiAgICAgKi9cclxuICAgIFpvb21SZXNldCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21SZXNldE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUy9MaW51eCBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXHJcbiAgICAgKiBHYXRoZXJzIGluZm9ybWF0aW9uIGFib3V0IHRoZSBkcm9wIHRhcmdldCBlbGVtZW50IGFuZCBzZW5kcyBpdCBiYWNrIHRvIHRoZSBHbyBiYWNrZW5kLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cclxuICAgICAqIEBwYXJhbSB4IC0gVGhlIHgtY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudCwgaW4gbG9naWNhbCAoQ1NTKSBwaXhlbHMgcmVsYXRpdmUgdG8gdGhlIHdlYnZpZXcuXHJcbiAgICAgKiBAcGFyYW0geSAtIFRoZSB5LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQsIGluIGxvZ2ljYWwgKENTUykgcGl4ZWxzIHJlbGF0aXZlIHRvIHRoZSB3ZWJ2aWV3LlxyXG4gICAgICovXHJcbiAgICBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wKGZpbGVuYW1lczogc3RyaW5nW10sIHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XHJcbiAgICAgICAgLy8gQ2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGlnbm9yZSB0aGUgZHJvcFxyXG4gICAgICAgIH1cclxuICAgICAgICBcclxuICAgICAgICBjb25zdCBlbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KTtcclxuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQoZWxlbWVudCk7XHJcblxyXG4gICAgICAgIGlmICghZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICAvLyBEcm9wIHdhcyBub3Qgb24gYSBkZXNpZ25hdGVkIGRyb3AgdGFyZ2V0IC0gaWdub3JlXHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGNvbnN0IGVsZW1lbnREZXRhaWxzID0ge1xyXG4gICAgICAgICAgICBpZDogZHJvcFRhcmdldC5pZCxcclxuICAgICAgICAgICAgY2xhc3NMaXN0OiBBcnJheS5mcm9tKGRyb3BUYXJnZXQuY2xhc3NMaXN0KSxcclxuICAgICAgICAgICAgYXR0cmlidXRlczoge30gYXMgeyBba2V5OiBzdHJpbmddOiBzdHJpbmcgfSxcclxuICAgICAgICB9O1xyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZHJvcFRhcmdldC5hdHRyaWJ1dGVzLmxlbmd0aDsgaSsrKSB7XHJcbiAgICAgICAgICAgIGNvbnN0IGF0dHIgPSBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXNbaV07XHJcbiAgICAgICAgICAgIGVsZW1lbnREZXRhaWxzLmF0dHJpYnV0ZXNbYXR0ci5uYW1lXSA9IGF0dHIudmFsdWU7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xyXG4gICAgICAgICAgICBmaWxlbmFtZXMsXHJcbiAgICAgICAgICAgIHgsXHJcbiAgICAgICAgICAgIHksXHJcbiAgICAgICAgICAgIGVsZW1lbnREZXRhaWxzLFxyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIHRoaXNbY2FsbGVyU3ltXShGaWxlc0Ryb3BwZWQsIHBheWxvYWQpO1xyXG4gICAgICAgIFxyXG4gICAgICAgIC8vIENsZWFuIHVwIG5hdGl2ZSBkcmFnIHN0YXRlIGFmdGVyIGRyb3BcclxuICAgICAgICBjbGVhbnVwTmF0aXZlRHJhZygpO1xyXG4gICAgfVxyXG4gIFxyXG4gICAgLyoqXHJcbiAgICAgKiBNb3ZlcyB0aGUgd2luZG93IHRvIHRoZSBjZW50ZXIgb2YgdGhlIHNwZWNpZmllZCBzY3JlZW4ncyB3b3JrIGFyZWEuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHNjcmVlbklEIC0gVGhlIElEIG9mIHRoZSB0YXJnZXQgc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBTZXRTY3JlZW4oc2NyZWVuSUQ6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2NyZWVuTWV0aG9kLCB7IHNjcmVlbklEIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qIFRyaWdnZXJzIFdpbmRvd3MgMTEgU25hcCBBc3Npc3QgZmVhdHVyZSAoV2luZG93cyBvbmx5KS5cclxuICAgICAqIFRoaXMgaXMgZXF1aXZhbGVudCB0byBwcmVzc2luZyBXaW4rWiBhbmQgc2hvd3Mgc25hcCBsYXlvdXQgb3B0aW9ucy5cclxuICAgICAqL1xyXG4gICAgU25hcEFzc2lzdCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNuYXBBc3Npc3RNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogT3BlbnMgdGhlIHByaW50IGRpYWxvZyBmb3IgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgUHJpbnQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQcmludE1ldGhvZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXHJcbiAqL1xyXG5jb25zdCB0aGlzV2luZG93ID0gbmV3IFdpbmRvdygnJyk7XHJcblxyXG4vKipcclxuICogU2V0cyB1cCBnbG9iYWwgZHJhZyBhbmQgZHJvcCBldmVudCBsaXN0ZW5lcnMgZm9yIGZpbGUgZHJvcHMuXHJcbiAqIEhhbmRsZXMgdmlzdWFsIGZlZWRiYWNrIChob3ZlciBzdGF0ZSkgYW5kIGZpbGUgZHJvcCBwcm9jZXNzaW5nLlxyXG4gKi9cclxuZnVuY3Rpb24gc2V0dXBEcm9wVGFyZ2V0TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZG9jRWxlbWVudCA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudDtcclxuICAgIGxldCBkcmFnRW50ZXJDb3VudGVyID0gMDtcclxuXHJcbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdlbnRlcicsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xyXG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJhZ3MsIGxldCBvdGhlciBkcmFncyBwYXNzIHRocm91Z2hcclxuICAgICAgICB9XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxyXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7IC8vIFNob3cgXCJubyBkcm9wXCIgY3Vyc29yXHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgZG9uJ3Qgc2hvdyBob3ZlciBlZmZlY3RzXHJcbiAgICAgICAgfVxyXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIrKztcclxuICAgICAgICBcclxuICAgICAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZKTtcclxuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XHJcblxyXG4gICAgICAgIC8vIFVwZGF0ZSBob3ZlciBzdGF0ZVxyXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBkcm9wVGFyZ2V0LmNsYXNzTGlzdC5hZGQoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgIH0sIGZhbHNlKTtcclxuXHJcbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdvdmVyJywgKGV2ZW50KSA9PiB7XHJcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gT25seSBoYW5kbGUgZmlsZSBkcmFnc1xyXG4gICAgICAgIH1cclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJzsgLy8gU2hvdyBcIm5vIGRyb3BcIiBjdXJzb3JcclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcclxuICAgICAgICB9XHJcbiAgICAgICAgXHJcbiAgICAgICAgLy8gVXBkYXRlIGRyb3AgdGFyZ2V0IGFzIGN1cnNvciBtb3Zlc1xyXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xyXG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudCh0YXJnZXRFbGVtZW50KTtcclxuICAgICAgICBcclxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBcclxuICAgICAgICBpZiAoZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBpZiAoIWRyb3BUYXJnZXQuY2xhc3NMaXN0LmNvbnRhaW5zKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUykpIHtcclxuICAgICAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcclxuICAgICAgICB9XHJcbiAgICB9LCBmYWxzZSk7XHJcblxyXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCAoZXZlbnQpID0+IHtcclxuICAgICAgICBpZiAoIWV2ZW50LmRhdGFUcmFuc2Zlcj8udHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIFxyXG4gICAgICAgIC8vIE9uIExpbnV4L1dlYktpdEdUSyBhbmQgbWFjT1MsIGRyYWdsZWF2ZSBmaXJlcyBpbW1lZGlhdGVseSB3aXRoIHJlbGF0ZWRUYXJnZXQ9bnVsbCB3aGVuIG5hdGl2ZVxyXG4gICAgICAgIC8vIGRyYWcgaGFuZGxpbmcgaXMgaW52b2x2ZWQuIElnbm9yZSB0aGVzZSBzcHVyaW91cyBldmVudHMgLSB3ZSdsbCBjbGVhbiB1cCBvbiBkcm9wIGluc3RlYWQuXHJcbiAgICAgICAgaWYgKGV2ZW50LnJlbGF0ZWRUYXJnZXQgPT09IG51bGwpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBcclxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyLS07XHJcbiAgICAgICAgXHJcbiAgICAgICAgaWYgKGRyYWdFbnRlckNvdW50ZXIgPT09IDAgfHwgXHJcbiAgICAgICAgICAgIChjdXJyZW50RHJvcFRhcmdldCAmJiAhY3VycmVudERyb3BUYXJnZXQuY29udGFpbnMoZXZlbnQucmVsYXRlZFRhcmdldCBhcyBOb2RlKSkpIHtcclxuICAgICAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG5cclxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJvcCcsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xyXG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJvcHNcclxuICAgICAgICB9XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxyXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyID0gMDtcclxuICAgICAgICBcclxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBoYW5kbGUgZmlsZSBkcm9wcyB2aWEgSmF2YVNjcmlwdFxyXG4gICAgICAgIC8vIE9uIG1hY09TL0xpbnV4LCBuYXRpdmUgY29kZSB3aWxsIGNhbGwgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcFxyXG4gICAgICAgIGlmIChjYW5SZXNvbHZlRmlsZVBhdGhzKCkpIHtcclxuICAgICAgICAgICAgY29uc3QgZmlsZXM6IEZpbGVbXSA9IFtdO1xyXG4gICAgICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGNvbnN0IGl0ZW0gb2YgZXZlbnQuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKGl0ZW0ua2luZCA9PT0gJ2ZpbGUnKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IGZpbGUgPSBpdGVtLmdldEFzRmlsZSgpO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZmlsZSkgZmlsZXMucHVzaChmaWxlKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0gZWxzZSBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyLmZpbGVzKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGNvbnN0IGZpbGUgb2YgZXZlbnQuZGF0YVRyYW5zZmVyLmZpbGVzKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgZmlsZXMucHVzaChmaWxlKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBcclxuICAgICAgICAgICAgaWYgKGZpbGVzLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgICAgIHJlc29sdmVGaWxlUGF0aHMoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZmlsZXMpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG59XHJcblxyXG4vLyBJbml0aWFsaXplIGxpc3RlbmVycyB3aGVuIHRoZSBzY3JpcHQgbG9hZHNcclxuaWYgKHR5cGVvZiB3aW5kb3cgIT09IFwidW5kZWZpbmVkXCIgJiYgdHlwZW9mIGRvY3VtZW50ICE9PSBcInVuZGVmaW5lZFwiKSB7XHJcbiAgICBzZXR1cERyb3BUYXJnZXRMaXN0ZW5lcnMoKTtcclxufVxyXG5cclxuZXhwb3J0IGRlZmF1bHQgdGhpc1dpbmRvdztcclxuIiwiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcclxuaW1wb3J0IHsgUXVlc3Rpb24gfSBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XHJcbmltcG9ydCB7IEVtaXQgfSBmcm9tIFwiLi9ldmVudHMuanNcIjtcclxuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XHJcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93LmpzXCI7XHJcblxyXG4vKipcclxuICogU2VuZHMgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgb3B0aW9uYWwgZGF0YS5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXHJcbiAqIEBwYXJhbSBbZGF0YT1udWxsXSAtIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XHJcbiAgICBFbWl0KGV2ZW50TmFtZSwgZGF0YSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAqXHJcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lOiBzdHJpbmcsIG1ldGhvZE5hbWU6IHN0cmluZykge1xyXG4gICAgY29uc3QgdGFyZ2V0V2luZG93ID0gV2luZG93LkdldCh3aW5kb3dOYW1lKTtcclxuICAgIGNvbnN0IG1ldGhvZCA9ICh0YXJnZXRXaW5kb3cgYXMgYW55KVttZXRob2ROYW1lXTtcclxuXHJcbiAgICBpZiAodHlwZW9mIG1ldGhvZCAhPT0gXCJmdW5jdGlvblwiKSB7XHJcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWV0aG9kLmNhbGwodGFyZ2V0V2luZG93KTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXQuXHJcbiAqL1xyXG5mdW5jdGlvbiBvbldNTFRyaWdnZXJlZChldjogRXZlbnQpOiB2b2lkIHtcclxuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0IGFzIEVsZW1lbnQ7XHJcblxyXG4gICAgZnVuY3Rpb24gcnVuRWZmZWN0KGNob2ljZSA9IFwiWWVzXCIpIHtcclxuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcclxuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XHJcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XHJcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW51cmwnKTtcclxuXHJcbiAgICAgICAgaWYgKGV2ZW50VHlwZSAhPT0gbnVsbClcclxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XHJcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcclxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcclxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XHJcblxyXG4gICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICBRdWVzdGlvbih7XHJcbiAgICAgICAgICAgIFRpdGxlOiBcIkNvbmZpcm1cIixcclxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcclxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxyXG4gICAgICAgICAgICBCdXR0b25zOiBbXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XHJcbiAgICAgICAgICAgIF1cclxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJ1bkVmZmVjdCgpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBjb250cm9sbGVyU3ltID0gU3ltYm9sKFwiY29udHJvbGxlclwiKTtcclxuY29uc3QgdHJpZ2dlck1hcFN5bSA9IFN5bWJvbChcInRyaWdnZXJNYXBcIik7XHJcbmNvbnN0IGVsZW1lbnRDb3VudFN5bSA9IFN5bWJvbChcImVsZW1lbnRDb3VudFwiKTtcclxuXHJcbi8qKlxyXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXHJcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXHJcbiAqL1xyXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XHJcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cclxuICAgIFtjb250cm9sbGVyU3ltXTogQWJvcnRDb250cm9sbGVyO1xyXG5cclxuICAgIGNvbnN0cnVjdG9yKCkge1xyXG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIGFuIG9wdGlvbnMgb2JqZWN0IGZvciBhZGRFdmVudExpc3RlbmVyIHRoYXQgdGllcyB0aGUgbGlzdGVuZXJcclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxyXG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xyXG4gICAgICovXHJcbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xyXG4gICAgICAgIHJldHVybiB7IHNpZ25hbDogdGhpc1tjb250cm9sbGVyU3ltXS5zaWduYWwgfTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzIGFuZCByZXNldHMgdGhlIHJlZ2lzdHJ5LlxyXG4gICAgICovXHJcbiAgICByZXNldCgpOiB2b2lkIHtcclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dLmFib3J0KCk7XHJcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFdlYWtNYXBSZWdpc3RyeSBtYXBzIGFjdGl2ZSB0cmlnZ2VyIGV2ZW50cyB0byBlYWNoIERPTSBlbGVtZW50IHRocm91Z2ggYSBXZWFrTWFwLlxyXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcclxuICogY29sbGVjdGlvbiBvZiB0aGUgaW52b2x2ZWQgZWxlbWVudHMuXHJcbiAqL1xyXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xyXG4gICAgLyoqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy4gKi9cclxuICAgIFt0cmlnZ2VyTWFwU3ltXTogV2Vha01hcDxFbGVtZW50LCBzdHJpbmdbXT47XHJcbiAgICAvKiogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLiAqL1xyXG4gICAgW2VsZW1lbnRDb3VudFN5bV06IG51bWJlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcigpIHtcclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxyXG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XHJcbiAgICAgKi9cclxuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XHJcbiAgICAgICAgaWYgKCF0aGlzW3RyaWdnZXJNYXBTeW1dLmhhcyhlbGVtZW50KSkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0rKzsgfVxyXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcclxuICAgICAgICByZXR1cm4ge307XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cclxuICAgICAqL1xyXG4gICAgcmVzZXQoKTogdm9pZCB7XHJcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xyXG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuXHJcbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwU3ltXS5nZXQoZWxlbWVudCk7XHJcbiAgICAgICAgICAgIGlmICh0cmlnZ2VycyAhPSBudWxsKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXS0tOyB9XHJcblxyXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXHJcbiAgICAgICAgICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XHJcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcclxuICAgIH1cclxufVxyXG5cclxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xyXG5cclxuLyoqXHJcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cclxuICovXHJcbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50OiBFbGVtZW50KTogdm9pZCB7XHJcbiAgICBjb25zdCB0cmlnZ2VyUmVnRXhwID0gL1xcUysvZztcclxuICAgIGNvbnN0IHRyaWdnZXJBdHRyID0gKGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiKTtcclxuICAgIGNvbnN0IHRyaWdnZXJzOiBzdHJpbmdbXSA9IFtdO1xyXG5cclxuICAgIGxldCBtYXRjaDtcclxuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcclxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcclxuXHJcbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XHJcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVuYWJsZSgpOiB2b2lkIHtcclxuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpOiB2b2lkIHtcclxuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xyXG4gICAgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XSwgW3dtbC13aW5kb3ddLCBbd21sLW9wZW51cmxdLCBbZGF0YS13bWwtZXZlbnRdLCBbZGF0YS13bWwtd2luZG93XSwgW2RhdGEtd21sLW9wZW51cmxdJykuZm9yRWFjaChhZGRXTUxMaXN0ZW5lcnMpO1xyXG59XHJcbiIsIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBHb1N0cmVhbSDigJQgdGhlIFdlYlNvY2tldCBwcm9ncmFtbWluZyBtb2RlbCB3aXRob3V0IGEgbGlzdGVuaW5nIHNvY2tldC5cclxuICpcclxuICogQSBXZWJTb2NrZXQgY2Fubm90IGJlIHNwb2tlbiBvdmVyIGEgY3VzdG9tIFVSTCBzY2hlbWUsIHNvIGdldHRpbmcgb25lIGluc2lkZVxyXG4gKiBhIHdlYnZpZXcgbWVhbnMgYmluZGluZyBhIHJlYWwgVENQIHBvcnQuIFRoaXMgc3BlYWtzIHRoZSBzYW1lIHNoYXBlIG92ZXIgdGhlXHJcbiAqIGFzc2V0IHNlcnZlciB0aGUgYXBwIGFscmVhZHkgc2VydmVzOiBHb+KGkkpTIG92ZXIgb25lIGhlbGQgcG9sbCBwZXIgcGFnZSwgSlPihpJHb1xyXG4gKiBvdmVyIGEgbm9ybWFsIFBPU1QuXHJcbiAqXHJcbiAqIHtAbGluayBTdHJlYW19IHJldHVybnMgc3luY2hyb25vdXNseSB3aXRoIGByZWFkeVN0YXRlID09PSBDT05ORUNUSU5HYCwgdGhlIHdheVxyXG4gKiBgbmV3IFdlYlNvY2tldCh1cmwpYCBkb2VzLCBzbyBhIGNvbm5lY3Rpb24gY2FuIGJlIGNyZWF0ZWQgYXQgbW9kdWxlIHNjb3BlOlxyXG4gKlxyXG4gKiBgYGBqc1xyXG4gKiBleHBvcnQgY29uc3QgVGVsZW1ldHJ5ID0gU3RyZWFtKFwidGVsZW1ldHJ5XCIpOyAgIC8vIGNvbm5lY3RzIGluIHRoZSBiYWNrZ3JvdW5kXHJcbiAqIGBgYFxyXG4gKi9cclxuXHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gXCIuL25hbm9pZC5qc1wiO1xyXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xyXG5cclxuLy8gRXZlcnkgcGllY2Ugb2Ygc3RhdGUgdGhhdCB0d28gcnVudGltZSBpbnN0YW5jZXMgbXVzdCBhZ3JlZSBvbiBsaXZlcyBoZXJlLCBvblxyXG4vLyB3aW5kb3csIG5vdCBpbiBtb2R1bGUgc2NvcGUuIFRoZSBydW50aW1lIGNhbiBiZSBldmFsdWF0ZWQgbW9yZSB0aGFuIG9uY2UgaW4gYVxyXG4vLyBwYWdlIOKAlCB0aGUgcGxhdGZvcm0gaW5qZWN0cyBhIGNvcHkgYW5kIGFuIGFwcCBtYXkgaW1wb3J0IGl0IHRvbyDigJQgYW5kIGFcclxuLy8gcGVyLW1vZHVsZSByZWdpc3RyeSBtZWFudCBib3RoIGluc3RhbmNlcyBhbGxvY2F0ZWQgY29ubmVjdGlvbiBpZCAxLCByYW5cclxuLy8gY29tcGV0aW5nIHBvbGxzIGFnYWluc3Qgb25lIHNlc3Npb24sIGFuZCBkaXNjYXJkZWQgZnJhbWVzIGFkZHJlc3NlZCB0byBpZHNcclxuLy8gb25seSB0aGUgb3RoZXIgaW5zdGFuY2Uga25ldyBhYm91dC4gVGhlIHNlc3Npb24gaWQgYWxvbmUgd2FzIG5vdCBlbm91Z2guXHJcbmludGVyZmFjZSBTdHJlYW1HbG9iYWxzIHtcclxuICAgIHNlc3Npb246IHN0cmluZztcclxuICAgIGdlbmVyYXRpb246IG51bWJlcjtcclxuICAgIG5leHRDb25uSUQ6IG51bWJlcjtcclxuICAgIGNvbm5lY3Rpb25zOiBNYXA8bnVtYmVyLCBXYWlsc1NvY2tldD47XHJcbiAgICBwb2xsaW5nOiBib29sZWFuO1xyXG4gICAgcG9sbEFib3J0PzogQWJvcnRDb250cm9sbGVyO1xyXG59XHJcblxyXG5jb25zdCBHRU5FUkFUSU9OX0tFWSA9IFwiX193YWlsc19zdHJlYW1fZ2VuZXJhdGlvblwiO1xyXG5jb25zdCBHRU5FUkFUSU9OX05BTUVfTUFSS0VSID0gXCJcXHUwMDFmX193YWlsc19zdHJlYW1fZ2VuZXJhdGlvbl9fOlwiO1xyXG5cclxuLy8gd2luZG93Lm5hbWUgYmVsb25ncyB0byB0aGUgYnJvd3NpbmcgY29udGV4dCByYXRoZXIgdGhhbiB0aGUgRG9jdW1lbnQsIHNvIGl0XHJcbi8vIHN1cnZpdmVzIGEgcmVsb2FkIGV2ZW4gd2hlbiBXZWIgU3RvcmFnZSBpcyBkaXNhYmxlZC4gUHJlc2VydmUgYW55IG5hbWUgdGhlXHJcbi8vIGFwcGxpY2F0aW9uIGFzc2lnbmVkIGFuZCByZXNlcnZlIG9ubHkgYSBzdWZmaXggZm9yIHRoZSBTdHJlYW0gZ2VuZXJhdGlvbi5cclxuZnVuY3Rpb24gbmFtZWRQYWdlR2VuZXJhdGlvbigpOiBudW1iZXIge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICBjb25zdCBtYXJrZXIgPSB3aW5kb3cubmFtZS5sYXN0SW5kZXhPZihHRU5FUkFUSU9OX05BTUVfTUFSS0VSKTtcclxuICAgICAgICBpZiAobWFya2VyIDwgMCkgcmV0dXJuIDA7XHJcbiAgICAgICAgY29uc3QgdmFsdWUgPSBOdW1iZXIod2luZG93Lm5hbWUuc2xpY2UobWFya2VyICsgR0VORVJBVElPTl9OQU1FX01BUktFUi5sZW5ndGgpKTtcclxuICAgICAgICByZXR1cm4gTnVtYmVyLmlzU2FmZUludGVnZXIodmFsdWUpICYmIHZhbHVlID4gMCA/IHZhbHVlIDogMDtcclxuICAgIH0gY2F0Y2gge1xyXG4gICAgICAgIHJldHVybiAwO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiByZW1lbWJlck5hbWVkUGFnZUdlbmVyYXRpb24oZ2VuZXJhdGlvbjogbnVtYmVyKTogdm9pZCB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGNvbnN0IG1hcmtlciA9IHdpbmRvdy5uYW1lLmxhc3RJbmRleE9mKEdFTkVSQVRJT05fTkFNRV9NQVJLRVIpO1xyXG4gICAgICAgIGNvbnN0IGFwcGxpY2F0aW9uTmFtZSA9IG1hcmtlciA8IDAgPyB3aW5kb3cubmFtZSA6IHdpbmRvdy5uYW1lLnNsaWNlKDAsIG1hcmtlcik7XHJcbiAgICAgICAgd2luZG93Lm5hbWUgPSBhcHBsaWNhdGlvbk5hbWUgKyBHRU5FUkFUSU9OX05BTUVfTUFSS0VSICsgZ2VuZXJhdGlvbjtcclxuICAgIH0gY2F0Y2gge1xyXG4gICAgICAgIC8vIFNvbWUgZW1iZWRkZWQgYnJvd3NlciBwb2xpY2llcyBjYW4gZGVueSBhY2Nlc3MgdG8gdGhlIGJyb3dzaW5nXHJcbiAgICAgICAgLy8gY29udGV4dCBuYW1lLiBUaGUgdGltZSBvcmlnaW4gcmVtYWlucyB0aGUgbGFzdC1yZXNvcnQgZmFsbGJhY2suXHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIFJlcXVlc3QgYXJyaXZhbCBvcmRlciBpcyBub3QgcGFnZSBvcmRlcjogYSByZXF1ZXN0IHN0YXJ0ZWQganVzdCBiZWZvcmUgYVxyXG4vLyBuYXZpZ2F0aW9uIGNhbiByZWFjaCBHbyBhZnRlciB0aGUgcmVwbGFjZW1lbnQgcGFnZSdzIGZpcnN0IHJlcXVlc3QuIEtlZXAgYVxyXG4vLyBtb25vdG9uaWMgY291bnRlciBpbiB0aGlzIHRvcC1sZXZlbCBicm93c2luZyBjb250ZXh0IHNvIEdvIGNhbiBkaXN0aW5ndWlzaFxyXG4vLyB0aG9zZSBwYWdlIGluY2FybmF0aW9ucyB3aXRob3V0IHRydXN0aW5nIHNlcnZlciBzY2hlZHVsaW5nIG9yZGVyLlxyXG5mdW5jdGlvbiBuZXh0UGFnZUdlbmVyYXRpb24oKTogbnVtYmVyIHtcclxuICAgIGlmICghaGFzRE9NKSByZXR1cm4gMTtcclxuICAgIC8vIEFuY2hvciB0aGUgdmFsdWUgdG8gdGhlIHBhZ2UncyBjcmVhdGlvbiB0aW1lIGFzIHdlbGwgYXMgdGhlIHN0b3JlZFxyXG4gICAgLy8gY291bnRlci4gSWYgYXBwbGljYXRpb24gY29kZSBjbGVhcnMgc2Vzc2lvblN0b3JhZ2UsIHRoZSBuZXh0IHJlbG9hZCBtdXN0XHJcbiAgICAvLyBzdGlsbCBiZSBuZXdlciB0aGFuIHRoZSBnZW5lcmF0aW9uIEdvIGhhcyBhbHJlYWR5IHJldGlyZWQ7IHJlc3RhcnRpbmdcclxuICAgIC8vIGZyb20gMSB3b3VsZCBtYWtlIGV2ZXJ5IHN0cmVhbSByZXF1ZXN0IHJlY2VpdmUgNDEwIHVudGlsIHRoZSBob3N0IGV4aXRzLlxyXG4gICAgLy8gUGVyZm9ybWFuY2UudGltZU9yaWdpbiBpcyBhYnNlbnQgb24gb2xkZXIgV2ViS2l0IHJlbGVhc2VzIHRoYXQgV2FpbHNcclxuICAgIC8vIHN0aWxsIHRhcmdldHMuIERhdGUubm93IGtlZXBzIHRob3NlIGVuZ2luZXMgb24gdGhlIHByb3RvY29sOyB0aGUgc3RvcmVkXHJcbiAgICAvLyBjb3VudGVyIGRpc2FtYmlndWF0ZXMgcmVsb2FkcyB0aGF0IGhhcHBlbiB3aXRoaW4gdGhlIHNhbWUgbWlsbGlzZWNvbmQuXHJcbiAgICBjb25zdCBvcmlnaW4gPSB0eXBlb2YgcGVyZm9ybWFuY2UgIT09IFwidW5kZWZpbmVkXCIgJiYgTnVtYmVyLmlzRmluaXRlKHBlcmZvcm1hbmNlLnRpbWVPcmlnaW4pXHJcbiAgICAgICAgPyBwZXJmb3JtYW5jZS50aW1lT3JpZ2luXHJcbiAgICAgICAgOiBEYXRlLm5vdygpO1xyXG4gICAgY29uc3QgcGFnZVRpbWUgPSBNYXRoLm1heCgxLCBNYXRoLmZsb29yKG9yaWdpbiAqIDEwMDApKTtcclxuICAgIGxldCBjdXJyZW50ID0gbmFtZWRQYWdlR2VuZXJhdGlvbigpO1xyXG4gICAgdHJ5IHtcclxuICAgICAgICBjb25zdCByYXcgPSB3aW5kb3cuc2Vzc2lvblN0b3JhZ2UuZ2V0SXRlbShHRU5FUkFUSU9OX0tFWSk7XHJcbiAgICAgICAgY29uc3Qgc3RvcmVkID0gcmF3ID09PSBudWxsID8gMCA6IE51bWJlcihyYXcpO1xyXG4gICAgICAgIGlmIChOdW1iZXIuaXNTYWZlSW50ZWdlcihzdG9yZWQpICYmIHN0b3JlZCA+IGN1cnJlbnQpIGN1cnJlbnQgPSBzdG9yZWQ7XHJcbiAgICB9IGNhdGNoIHtcclxuICAgICAgICAvLyBTdG9yYWdlIGNhbiBiZSBkaXNhYmxlZCBieSBwb2xpY3kuIHdpbmRvdy5uYW1lIHN1cHBsaWVzIHRoZVxyXG4gICAgICAgIC8vIGJyb3dzaW5nLWNvbnRleHQgY291bnRlciBpbiB0aGF0IGNhc2UuXHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgbmV4dCA9IE51bWJlci5pc1NhZmVJbnRlZ2VyKGN1cnJlbnQpICYmIGN1cnJlbnQgPj0gMCAmJiBjdXJyZW50IDwgTnVtYmVyLk1BWF9TQUZFX0lOVEVHRVJcclxuICAgICAgICA/IE1hdGgubWF4KGN1cnJlbnQgKyAxLCBwYWdlVGltZSlcclxuICAgICAgICA6IHBhZ2VUaW1lO1xyXG4gICAgdHJ5IHtcclxuICAgICAgICB3aW5kb3cuc2Vzc2lvblN0b3JhZ2Uuc2V0SXRlbShHRU5FUkFUSU9OX0tFWSwgU3RyaW5nKG5leHQpKTtcclxuICAgIH0gY2F0Y2gge1xyXG4gICAgICAgIC8vIFRoZSBicm93c2luZy1jb250ZXh0IGZhbGxiYWNrIGJlbG93IGlzIHN1ZmZpY2llbnQgd2hlbiBzdG9yYWdlIGlzXHJcbiAgICAgICAgLy8gdW5hdmFpbGFibGUuXHJcbiAgICB9XHJcbiAgICByZW1lbWJlck5hbWVkUGFnZUdlbmVyYXRpb24obmV4dCk7XHJcbiAgICByZXR1cm4gbmV4dDtcclxufVxyXG5cclxuZnVuY3Rpb24gc3RyZWFtR2xvYmFscygpOiBTdHJlYW1HbG9iYWxzIHtcclxuICAgIGNvbnN0IGZyZXNoID0gKCk6IFN0cmVhbUdsb2JhbHMgPT4gKHtcclxuICAgICAgICBzZXNzaW9uOiBuYW5vaWQoKSxcclxuICAgICAgICBnZW5lcmF0aW9uOiBuZXh0UGFnZUdlbmVyYXRpb24oKSxcclxuICAgICAgICBuZXh0Q29ubklEOiAxLFxyXG4gICAgICAgIGNvbm5lY3Rpb25zOiBuZXcgTWFwPG51bWJlciwgV2FpbHNTb2NrZXQ+KCksXHJcbiAgICAgICAgcG9sbGluZzogZmFsc2UsXHJcbiAgICB9KTtcclxuICAgIGlmICghaGFzRE9NKSByZXR1cm4gZnJlc2goKTtcclxuICAgIGNvbnN0IHcgPSAod2luZG93IGFzIGFueSkuX3dhaWxzID0gKHdpbmRvdyBhcyBhbnkpLl93YWlscyB8fCB7fTtcclxuICAgIHJldHVybiAody5fX3N0cmVhbSA9IHcuX19zdHJlYW0gfHwgZnJlc2goKSk7XHJcbn1cclxuXHJcbmNvbnN0IEcgPSBzdHJlYW1HbG9iYWxzKCk7XHJcblxyXG5jb25zdCBIRFJfU0VTU0lPTiA9IFwieC13YWlscy1zdHJlYW0tc2Vzc2lvblwiO1xyXG5jb25zdCBIRFJfR0VORVJBVElPTiA9IFwieC13YWlscy1zdHJlYW0tZ2VuZXJhdGlvblwiO1xyXG5jb25zdCBIRFJfQ09OTiA9IFwieC13YWlscy1zdHJlYW0tY29ublwiO1xyXG5jb25zdCBIRFJfS0lORCA9IFwieC13YWlscy1zdHJlYW0ta2luZFwiO1xyXG5jb25zdCBIRFJfTkFNRSA9IFwieC13YWlscy1zdHJlYW0tbmFtZVwiO1xyXG5jb25zdCBIRFJfQ0hVTksgPSBcIngtd2FpbHMtc3RyZWFtLWNodW5rXCI7XHJcbmNvbnN0IEhEUl9DSFVOS19JTkRFWCA9IFwieC13YWlscy1zdHJlYW0tY2h1bmstaW5kZXhcIjtcclxuY29uc3QgSERSX0NIVU5LX1RPVEFMID0gXCJ4LXdhaWxzLXN0cmVhbS1jaHVuay10b3RhbFwiO1xyXG5jb25zdCBIRFJfQkFUQ0ggPSBcIngtd2FpbHMtc3RyZWFtLWJhdGNoXCI7XHJcblxyXG5jb25zdCBLSU5EX0RBVEEgPSAwO1xyXG5jb25zdCBLSU5EX09QRU4gPSAxO1xyXG5jb25zdCBLSU5EX0NMT1NFID0gMjtcclxuY29uc3QgS0lORF9FUlJPUiA9IDM7XHJcblxyXG4vLyBNYXRjaGVzIHRoZSBydW50aW1lJ3Mgb3duIGxpbWl0OiBzdGF5IHVuZGVyIFdlYlZpZXcyJ3MgfjJNQiByZXF1ZXN0IGJvZHlcclxuLy8gYnVmZmVyaW5nIGxpbWl0IGluIFdlYlJlc291cmNlUmVxdWVzdGVkLlxyXG5jb25zdCBDSFVOS19USFJFU0hPTEQgPSA1MTIgKiAxMDI0O1xyXG5jb25zdCBNQVhfQkFUQ0hfRlJBTUVTID0gNDA5NjtcclxuXHJcbi8vIEJhY2tvZmYgZmxvb3IgYW5kIGNlaWxpbmcgZm9yIGEgZmFpbGluZyBwb2xsLiBUaGlzIGlzIHRoZSBcInJlYXNvbmFibGUgbG93ZXJcclxuLy8gbGltaXRcIiBvbiBwb2xsaW5nIOKAlCBhbiBlcnJvciBiYWNrb2ZmLCBub3QgYSBzdGVhZHktc3RhdGUgaW50ZXJ2YWwuIEluIG5vcm1hbFxyXG4vLyBvcGVyYXRpb24gdGhlIHBvbGwgcmUtaXNzdWVzIGltbWVkaWF0ZWx5LCBiZWNhdXNlIHRoZSBzZXJ2ZXIgaG9sZCBoYXMgYWxyZWFkeVxyXG4vLyBwcm92aWRlZCB3aGF0ZXZlciBkZWxheSB3YXMgYXBwcm9wcmlhdGUuXHJcbmNvbnN0IEJBQ0tPRkZfTUlOID0gMjUwO1xyXG5jb25zdCBCQUNLT0ZGX01BWCA9IDUwMDA7XHJcbmNvbnN0IHRleHRFbmNvZGVyID0gbmV3IFRleHRFbmNvZGVyKCk7XHJcblxyXG5mdW5jdGlvbiBzdHJlYW1VUkwoZW5kcG9pbnQ6IHN0cmluZyk6IHN0cmluZyB7XHJcbiAgICByZXR1cm4gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3N0cmVhbS9cIiArIGVuZHBvaW50O1xyXG59XHJcblxyXG5cclxuY2xhc3MgU3RyZWFtUmVxdWVzdENhbmNlbGxlZCBleHRlbmRzIEVycm9yIHtcclxuICAgIGNvbnN0cnVjdG9yKCkge1xyXG4gICAgICAgIHN1cGVyKFwic3RyZWFtIHJlcXVlc3QgY2FuY2VsbGVkXCIpO1xyXG4gICAgICAgIHRoaXMubmFtZSA9IFwiQWJvcnRFcnJvclwiO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBzbGVlcChtczogbnVtYmVyLCBzaWduYWw/OiBBYm9ydFNpZ25hbCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgaWYgKCFzaWduYWwpIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSkgPT4gc2V0VGltZW91dChyZXNvbHZlLCBtcykpO1xyXG4gICAgaWYgKHNpZ25hbC5hYm9ydGVkKSByZXR1cm4gUHJvbWlzZS5yZWplY3QobmV3IFN0cmVhbVJlcXVlc3RDYW5jZWxsZWQoKSk7XHJcblxyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBjb25zdCBvbkFib3J0ID0gKCkgPT4ge1xyXG4gICAgICAgICAgICBjbGVhclRpbWVvdXQodGltZXIpO1xyXG4gICAgICAgICAgICByZWplY3QobmV3IFN0cmVhbVJlcXVlc3RDYW5jZWxsZWQoKSk7XHJcbiAgICAgICAgfTtcclxuICAgICAgICBjb25zdCB0aW1lciA9IHNldFRpbWVvdXQoKCkgPT4ge1xyXG4gICAgICAgICAgICBzaWduYWwucmVtb3ZlRXZlbnRMaXN0ZW5lcihcImFib3J0XCIsIG9uQWJvcnQpO1xyXG4gICAgICAgICAgICByZXNvbHZlKCk7XHJcbiAgICAgICAgfSwgbXMpO1xyXG4gICAgICAgIHNpZ25hbC5hZGRFdmVudExpc3RlbmVyKFwiYWJvcnRcIiwgb25BYm9ydCwgeyBvbmNlOiB0cnVlIH0pO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBIHN0cmVhbSBjb25uZWN0aW9uLiBJbXBsZW1lbnRzIHRoZSB1c2VmdWwgc3Vic2V0IG9mIHRoZSBgV2ViU29ja2V0YFxyXG4gKiBpbnRlcmZhY2UsIHNvIHRoZSBzYW1lIGFwcGxpY2F0aW9uIGNvZGUgd29ya3MgYWdhaW5zdCBhIHJlYWwgV2ViU29ja2V0IGluXHJcbiAqIHNlcnZlciBidWlsZHMuXHJcbiAqXHJcbiAqIE9uZSBkZWxpYmVyYXRlIGRpdmVyZ2VuY2UgZnJvbSB0aGUgc3RhbmRhcmQ6IGBiaW5hcnlUeXBlYCBkZWZhdWx0cyB0b1xyXG4gKiBgXCJhcnJheWJ1ZmZlclwiYCByYXRoZXIgdGhhbiBgXCJibG9iXCJgLCBiZWNhdXNlIHN0cmVhbSBmcmFtZXMgYXJlIGFsd2F5c1xyXG4gKiBiaW5hcnkgYW5kIGEgQmxvYiB3b3VsZCBmb3JjZSBhbiBleHRyYSBhc3luYyBob3AgdG8gcmVhZCBldmVyeSBtZXNzYWdlLlxyXG4gKiBTZXR0aW5nIGl0IHRvIGBcImJsb2JcImAgYmVoYXZlcyBhcyB0aGUgc3RhbmRhcmQgZGVzY3JpYmVzLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIFdhaWxzU29ja2V0IGV4dGVuZHMgRXZlbnRUYXJnZXQge1xyXG4gICAgc3RhdGljIHJlYWRvbmx5IENPTk5FQ1RJTkcgPSAwO1xyXG4gICAgc3RhdGljIHJlYWRvbmx5IE9QRU4gPSAxO1xyXG4gICAgc3RhdGljIHJlYWRvbmx5IENMT1NJTkcgPSAyO1xyXG4gICAgc3RhdGljIHJlYWRvbmx5IENMT1NFRCA9IDM7XHJcblxyXG4gICAgcmVhZG9ubHkgQ09OTkVDVElORyA9IDA7XHJcbiAgICByZWFkb25seSBPUEVOID0gMTtcclxuICAgIHJlYWRvbmx5IENMT1NJTkcgPSAyO1xyXG4gICAgcmVhZG9ubHkgQ0xPU0VEID0gMztcclxuXHJcbiAgICAvKiogVGhlIHN0cmVhbSBuYW1lIHRoaXMgY29ubmVjdGlvbiB3YXMgb3BlbmVkIGFnYWluc3QuICovXHJcbiAgICByZWFkb25seSBuYW1lOiBzdHJpbmc7XHJcbiAgICByZWFkb25seSB1cmw6IHN0cmluZztcclxuXHJcbiAgICAvKiogQWx3YXlzIGVtcHR5OiBzdWJwcm90b2NvbHMgYW5kIGV4dGVuc2lvbnMgYXJlIG5vdCBuZWdvdGlhdGVkLiAqL1xyXG4gICAgcmVhZG9ubHkgcHJvdG9jb2wgPSBcIlwiO1xyXG4gICAgcmVhZG9ubHkgZXh0ZW5zaW9ucyA9IFwiXCI7XHJcblxyXG4gICAgYmluYXJ5VHlwZTogXCJhcnJheWJ1ZmZlclwiIHwgXCJibG9iXCIgPSBcImFycmF5YnVmZmVyXCI7XHJcbiAgICByZWFkeVN0YXRlOiBudW1iZXIgPSBXYWlsc1NvY2tldC5DT05ORUNUSU5HO1xyXG5cclxuICAgIC8vIERlY2xhcmVkIGZvciB0aGUgdHlwZTsgdGhlIGNvbnN0cnVjdG9yIHJlcGxhY2VzIGVhY2ggd2l0aCBhbiBhY2Nlc3NvclxyXG4gICAgLy8gdGhhdCByZWdpc3RlcnMgYSByZWFsIGxpc3RlbmVyLiBDYWxsaW5nIHRoZW0gZGlyZWN0bHkgaW5zdGVhZCB3b3VsZCBtYWtlXHJcbiAgICAvLyBhbnkgd3JhcHBlciAoc2VlIEpTT05TdHJlYW0pIHJlY2VpdmUgYm90aCB0aGUgcmF3IGFuZCB0aGUgZGVjb2RlZCBldmVudC5cclxuICAgIG9ub3BlbjogKChldjogRXZlbnQpID0+IHZvaWQpIHwgbnVsbCA9IG51bGw7XHJcbiAgICBvbm1lc3NhZ2U6ICgoZXY6IE1lc3NhZ2VFdmVudCkgPT4gdm9pZCkgfCBudWxsID0gbnVsbDtcclxuICAgIG9uY2xvc2U6ICgoZXY6IENsb3NlRXZlbnQpID0+IHZvaWQpIHwgbnVsbCA9IG51bGw7XHJcbiAgICBvbmVycm9yOiAoKGV2OiBFdmVudCkgPT4gdm9pZCkgfCBudWxsID0gbnVsbDtcclxuXHJcbiAgICAvKipcclxuICAgICAqIEBpbnRlcm5hbCBBcHBsaWVkIHRvIGV2ZXJ5IGluYm91bmQgcGF5bG9hZCBiZWZvcmUgZGlzcGF0Y2guIEpTT05TdHJlYW1cclxuICAgICAqIHJlcGxhY2VzIGl0LCBzbyBgb25tZXNzYWdlYCBhbmQgYGFkZEV2ZW50TGlzdGVuZXIoXCJtZXNzYWdlXCIsIOKApilgIHNlZSB0aGVcclxuICAgICAqIHNhbWUgZGVjb2RlZCB2YWx1ZSDigJQgaW50ZXJjZXB0aW5nIG9ubHkgYG9ubWVzc2FnZWAgbWFkZSB0aGUgdHdvIGRpc2FncmVlLlxyXG4gICAgICovXHJcbiAgICBfZGVjb2RlOiAocGF5bG9hZDogQXJyYXlCdWZmZXIpID0+IHVua25vd24gPSAocGF5bG9hZCkgPT4gcGF5bG9hZDtcclxuXHJcbiAgICAvKiogQGludGVybmFsICovIHJlYWRvbmx5IF9pZDogbnVtYmVyO1xyXG4gICAgcHJpdmF0ZSBfYnVmZmVyZWQgPSAwO1xyXG4gICAgLy8gU2VuZHMgYXJlIHNlcmlhbGlzZWQgcGVyIGNvbm5lY3Rpb246IGNvbmN1cnJlbnQgZmV0Y2goKSBQT1NUcyBkbyBub3RcclxuICAgIC8vIHByZXNlcnZlIG9yZGVyLCBhbmQgR28gcmVsaWVzIG9uIHNlbmQgb3JkZXIgYmVpbmcgdGhlIG9yZGVyIGl0IG9ic2VydmVzLlxyXG4gICAgcHJpdmF0ZSBfY2hhaW46IFByb21pc2U8dm9pZD4gPSBQcm9taXNlLnJlc29sdmUoKTtcclxuICAgIHByaXZhdGUgX3BlbmRpbmc6IFByb21pc2U8VWludDhBcnJheT5bXSA9IFtdO1xyXG4gICAgcHJpdmF0ZSBfZmx1c2hpbmcgPSBmYWxzZTtcclxuICAgIHByaXZhdGUgX29wZW5BYm9ydCA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIHByaXZhdGUgX3NlbmRBYm9ydCA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuXHJcbiAgICAvKiogQnl0ZXMgcXVldWVkIGJ5IHNlbmQoKSB0aGF0IGhhdmUgbm90IHlldCByZWFjaGVkIEdvLiAqL1xyXG4gICAgZ2V0IGJ1ZmZlcmVkQW1vdW50KCk6IG51bWJlciB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXMuX2J1ZmZlcmVkO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IHN0cmluZykge1xyXG4gICAgICAgIHN1cGVyKCk7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcclxuICAgICAgICB0aGlzLl9pZCA9IEcubmV4dENvbm5JRCsrO1xyXG4gICAgICAgIHRoaXMudXJsID0gaGFzRE9NID8gc3RyZWFtVVJMKFwicG9sbFwiKSArIFwiI1wiICsgZW5jb2RlVVJJQ29tcG9uZW50KG5hbWUpIDogXCJcIjtcclxuXHJcbiAgICAgICAgZm9yIChjb25zdCB0eXBlIG9mIFtcIm9wZW5cIiwgXCJtZXNzYWdlXCIsIFwiY2xvc2VcIiwgXCJlcnJvclwiXSkge1xyXG4gICAgICAgICAgICBkZWZpbmVIYW5kbGVyUHJvcGVydHkodGhpcywgdHlwZSk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBHLmNvbm5lY3Rpb25zLnNldCh0aGlzLl9pZCwgdGhpcyk7XHJcblxyXG4gICAgICAgIC8vIEZpcmUgYW5kIGZvcmdldDogdGhlIGFjayBhcnJpdmVzIGFzIGFuIG9wZW4gZnJhbWUgb24gdGhlIHBvbGwsIHdoaWNoXHJcbiAgICAgICAgLy8gaXMgd2hhdCBtb3ZlcyByZWFkeVN0YXRlIHRvIE9QRU4uXHJcbiAgICAgICAgdGhpcy5fY2hhaW4gPSB0aGlzLl9jaGFpbi50aGVuKCgpID0+XHJcbiAgICAgICAgICAgIHBvc3RGcmFtZSh0aGlzLl9pZCwgS0lORF9PUEVOLCBuZXcgVWludDhBcnJheSgwKSwgbmFtZSwgdGhpcy5fb3BlbkFib3J0LnNpZ25hbClcclxuICAgICAgICApLmNhdGNoKChlcnIpID0+IHtcclxuICAgICAgICAgICAgdGhpcy5fZmFpbChlcnIpO1xyXG4gICAgICAgIH0pO1xyXG5cclxuICAgICAgICBzdGFydFBvbGxpbmcoKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNlbmQgYSBmcmFtZSB0byBHby4gQWNjZXB0cyBhbnl0aGluZyBhIFdlYlNvY2tldCBkb2VzOyBzdHJpbmdzIGFyZVxyXG4gICAgICogZW5jb2RlZCBhcyBVVEYtOCwgc2luY2UgR28gcmVjZWl2ZXMgZXZlcnkgZnJhbWUgYXMgYSBieXRlIHNsaWNlLlxyXG4gICAgICovXHJcbiAgICBzZW5kKGRhdGE6IHN0cmluZyB8IEFycmF5QnVmZmVyTGlrZSB8IEFycmF5QnVmZmVyVmlldyB8IEJsb2IpOiB2b2lkIHtcclxuICAgICAgICBpZiAodGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DT05ORUNUSU5HKSB7XHJcbiAgICAgICAgICAgIC8vIE1hdGNoZXMgdGhlIFdlYlNvY2tldCBzcGVjaWZpY2F0aW9uLlxyXG4gICAgICAgICAgICB0aHJvdyBuZXcgRE9NRXhjZXB0aW9uKFwiU3RpbGwgaW4gQ09OTkVDVElORyBzdGF0ZS5cIiwgXCJJbnZhbGlkU3RhdGVFcnJvclwiKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgaWYgKHRoaXMucmVhZHlTdGF0ZSAhPT0gV2FpbHNTb2NrZXQuT1BFTikge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBSZXNvbHZlIG5vdyByYXRoZXIgdGhhbiBpbnNpZGUgdGhlIGNoYWluLiBNdXRhYmxlIGJpbmFyeSBpbnB1dHMgYXJlXHJcbiAgICAgICAgLy8gY29waWVkIHN5bmNocm9ub3VzbHkgc28gc2VuZCgpIHNuYXBzaG90cyB0aGVpciBieXRlcyBsaWtlIGEgbmF0aXZlXHJcbiAgICAgICAgLy8gV2ViU29ja2V0LCBldmVuIHdoZW4gYSBiYXRjaCB3YWl0cyBiZWhpbmQgYW4gaW4tZmxpZ2h0IHJlcXVlc3QuIEFcclxuICAgICAgICAvLyBCbG9iIGlzIGltbXV0YWJsZSBhbmQgaXMgdGhlcmVmb3JlIHNhZmUgdG8gcmVhZCBhc3luY2hyb25vdXNseSBvbmNlLlxyXG4gICAgICAgIGNvbnN0IGltbWVkaWF0ZSA9IHRvQnl0ZXNTeW5jKGRhdGEpO1xyXG4gICAgICAgIGNvbnN0IHNuYXBzaG90ID0gaW1tZWRpYXRlXHJcbiAgICAgICAgICAgID8gUHJvbWlzZS5yZXNvbHZlKGltbWVkaWF0ZSlcclxuICAgICAgICAgICAgOiB0b0J5dGVzKGRhdGEsIHRoaXMuX3NlbmRBYm9ydC5zaWduYWwpO1xyXG4gICAgICAgIC8vIGNsb3NlKCkgbWF5IGRpc2NhcmQgdGhpcyBwcm9taXNlIGJlZm9yZSB0aGUgZmx1c2ggaGFzIGRlcXVldWVkIGl0LlxyXG4gICAgICAgIC8vIEF0dGFjaCBhIHJlamVjdGlvbiBvYnNlcnZlciBpbW1lZGlhdGVseSBzbyBjYW5jZWxsaW5nIGFuIGFzeW5jaHJvbm91c1xyXG4gICAgICAgIC8vIEJsb2IgcmVhZCBjYW5ub3QgYmVjb21lIGFuIHVuaGFuZGxlZCByZWplY3Rpb247IFByb21pc2UuYWxsIHN0aWxsXHJcbiAgICAgICAgLy8gb2JzZXJ2ZXMgdGhlIG9yaWdpbmFsIHJlamVjdGlvbiB3aGVuIHRoZSBmbHVzaCBhbHJlYWR5IG93bnMgaXQuXHJcbiAgICAgICAgdm9pZCBzbmFwc2hvdC5jYXRjaCgoKSA9PiB7fSk7XHJcblxyXG4gICAgICAgIC8vIEEgQmxvYiBjYW5ub3QgYmUgY29udmVydGVkIHN5bmNocm9ub3VzbHksIGJ1dCBCbG9iLnNpemUgaXMgYXZhaWxhYmxlXHJcbiAgICAgICAgLy8gbm93LCBzbyBpdHMgYnl0ZXMgYXJlIHN0aWxsIGNvdW50ZWQgdGhlIG1vbWVudCBzZW5kKCkgcmV0dXJucy4gV2FpdGluZ1xyXG4gICAgICAgIC8vIGZvciBhcnJheUJ1ZmZlcigpIHRvIHJlc29sdmUgbGV0IGEgbGFyZ2UgQmxvYiBzbGlwIHBhc3QgY29kZSB1c2luZ1xyXG4gICAgICAgIC8vIGJ1ZmZlcmVkQW1vdW50IGZvciBiYWNrcHJlc3N1cmUuXHJcbiAgICAgICAgY29uc3QgYmxvYlNpemUgPVxyXG4gICAgICAgICAgICAhaW1tZWRpYXRlICYmIHR5cGVvZiBCbG9iICE9PSBcInVuZGVmaW5lZFwiICYmIGRhdGEgaW5zdGFuY2VvZiBCbG9iID8gZGF0YS5zaXplIDogMDtcclxuICAgICAgICBpZiAoYmxvYlNpemUpIHtcclxuICAgICAgICAgICAgdGhpcy5fYnVmZmVyZWQgKz0gYmxvYlNpemU7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBBY2NvdW50IGZvciB0aGUgYnl0ZXMgbm93LCBub3Qgd2hlbiB0aGUgY2hhaW4gcmVhY2hlcyB0aGVtLiBSZXBvcnRpbmdcclxuICAgICAgICAvLyB6ZXJvIHdoaWxlIGZyYW1lcyB3YWl0IGJlaGluZCBhbiBpbi1mbGlnaHQgYmF0Y2ggd291bGQgdGVsbCBhblxyXG4gICAgICAgIC8vIGFwcGxpY2F0aW9uIHVzaW5nIGJ1ZmZlcmVkQW1vdW50IGZvciBiYWNrcHJlc3N1cmUgdGhhdCBpdCBtYXkga2VlcFxyXG4gICAgICAgIC8vIHNlbmRpbmcuXHJcbiAgICAgICAgaWYgKGltbWVkaWF0ZSkge1xyXG4gICAgICAgICAgICB0aGlzLl9idWZmZXJlZCArPSBpbW1lZGlhdGUuYnl0ZUxlbmd0aDtcclxuICAgICAgICB9XHJcbiAgICAgICAgdGhpcy5fcGVuZGluZy5wdXNoKHNuYXBzaG90KTtcclxuXHJcbiAgICAgICAgLy8gUXVldWUgcmF0aGVyIHRoYW4gcG9zdC4gV2hhdGV2ZXIgYWNjdW11bGF0ZXMgd2hpbGUgYSByZXF1ZXN0IGlzIGluXHJcbiAgICAgICAgLy8gZmxpZ2h0IGdvZXMgb3V0IHRvZ2V0aGVyIGluIHRoZSBuZXh0IG9uZSwgc28gdGhlIHBlci1yZXF1ZXN0IGNvc3Qg4oCUXHJcbiAgICAgICAgLy8gYSBzY2hlbWUtaGFuZGxlciByb3VuZCB0cmlwLCBhYm91dCBlbGV2ZW4gY2dvIGNhbGxzIG9uIG1hY09TIOKAlCBpc1xyXG4gICAgICAgIC8vIGRpdmlkZWQgYWNyb3NzIHRoZSBiYXRjaC4gVW5kZXIgbGlnaHQgbG9hZCBhIGJhdGNoIGlzIG9uZSBmcmFtZSBhbmRcclxuICAgICAgICAvLyB0aGlzIGJlaGF2ZXMgZXhhY3RseSBhcyBiZWZvcmU7IHRoZSBiYXRjaGluZyBvbmx5IGFwcGVhcnMgd2hlbiB0aGVcclxuICAgICAgICAvLyBzZW5kZXIgaXMgb3V0cnVubmluZyB0aGUgdHJhbnNwb3J0LCB3aGljaCBpcyB3aGVuIGl0IGlzIG5lZWRlZC5cclxuICAgICAgICBpZiAodGhpcy5fZmx1c2hpbmcpIHJldHVybjtcclxuICAgICAgICB0aGlzLl9mbHVzaGluZyA9IHRydWU7XHJcblxyXG4gICAgICAgIHRoaXMuX2NoYWluID0gdGhpcy5fY2hhaW4udGhlbihhc3luYyAoKSA9PiB7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICB3aGlsZSAodGhpcy5fcGVuZGluZy5sZW5ndGggPiAwKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gQm91bmQgdGhlIHdvcmsgYW5kIHRlbXBvcmFyeSBwcm9taXNlIGFycmF5IGZvciBvbmUgcGFzcy5cclxuICAgICAgICAgICAgICAgICAgICAvLyBwb3N0QmF0Y2ggYXBwbGllcyB0aGUgd2lyZS1zaXplIGJvdW5kIGFmdGVyIHJlc29sdXRpb24uXHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcXVldWVkID0gdGhpcy5fcGVuZGluZy5zcGxpY2UoMCwgTUFYX0JBVENIX0ZSQU1FUyk7XHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgYmF0Y2ggPSBhd2FpdCBQcm9taXNlLmFsbChxdWV1ZWQpO1xyXG5cclxuICAgICAgICAgICAgICAgICAgICAvLyBBIHJlbW90ZSBjbG9zZSBtYXkgYXJyaXZlIHdoaWxlIGEgQmxvYiBpcyBiZWluZyByZWFkLlxyXG4gICAgICAgICAgICAgICAgICAgIC8vIFRoZSBjb25uZWN0aW9uIGlzIGFscmVhZHkgZ29uZSBpbiB0aGF0IGNhc2UsIHNvIGRvIG5vdFxyXG4gICAgICAgICAgICAgICAgICAgIC8vIGlzc3VlIGEgbGF0ZSBQT1NUIGZvciBieXRlcyB0aGUgcGVlciBjYW4gbm8gbG9uZ2VyIHRha2UuXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXMucmVhZHlTdGF0ZSA9PT0gV2FpbHNTb2NrZXQuQ0xPU0VEKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgdG90YWwgPSBiYXRjaC5yZWR1Y2UoKG4sIGIpID0+IG4gKyBiLmJ5dGVMZW5ndGgsIDApO1xyXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGF3YWl0IHBvc3RCYXRjaCh0aGlzLl9pZCwgYmF0Y2gsIHRoaXMuX3NlbmRBYm9ydC5zaWduYWwpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZmluYWxseSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIC8vIF9jbG9zZWQgcmVzZXRzIHRoZSBjb3VudGVyIGltbWVkaWF0ZWx5LiBBIHJlcXVlc3RcclxuICAgICAgICAgICAgICAgICAgICAgICAgLy8gYWxyZWFkeSBpbiBmbGlnaHQgbWF5IGZpbmlzaCBhZnRlcndhcmRzLCBhbmQgbXVzdFxyXG4gICAgICAgICAgICAgICAgICAgICAgICAvLyBub3QgZHJpdmUgYnVmZmVyZWRBbW91bnQgYmVsb3cgemVyby5cclxuICAgICAgICAgICAgICAgICAgICAgICAgdGhpcy5fYnVmZmVyZWQgPSBNYXRoLm1heCgwLCB0aGlzLl9idWZmZXJlZCAtIHRvdGFsKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0gZmluYWxseSB7XHJcbiAgICAgICAgICAgICAgICB0aGlzLl9mbHVzaGluZyA9IGZhbHNlO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSkuY2F0Y2goKGVycikgPT4ge1xyXG4gICAgICAgICAgICB0aGlzLl9mbHVzaGluZyA9IGZhbHNlO1xyXG4gICAgICAgICAgICB0aGlzLl9mYWlsKGVycik7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqIENsb3NlIHRoZSBjb25uZWN0aW9uLiBHbydzIGhhbmRsZXIgc2VlcyBSZWNlaXZlIGZhaWwgYW5kIHJldHVybnMuICovXHJcbiAgICBjbG9zZShjb2RlOiBudW1iZXIgPSAxMDAwLCByZWFzb246IHN0cmluZyA9IFwiXCIpOiB2b2lkIHtcclxuICAgICAgICBpZiAodGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DTE9TSU5HIHx8IHRoaXMucmVhZHlTdGF0ZSA9PT0gV2FpbHNTb2NrZXQuQ0xPU0VEKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICAgICAgdGhpcy5yZWFkeVN0YXRlID0gV2FpbHNTb2NrZXQuQ0xPU0lORztcclxuICAgICAgICB0aGlzLl9wZW5kaW5nLmxlbmd0aCA9IDA7XHJcbiAgICAgICAgdGhpcy5fYnVmZmVyZWQgPSAwO1xyXG4gICAgICAgIHRoaXMuX29wZW5BYm9ydC5hYm9ydCgpO1xyXG5cclxuICAgICAgICAvLyBBIGRhdGEgZnJhbWUgd2FpdGluZyBmb3IgcmVjZWl2ZXIgY2FwYWNpdHkgbXVzdCBub3QgcHJldmVudCB0aGVcclxuICAgICAgICAvLyByZXNlcnZlZCBjbG9zZSBjb250cm9sIGZyb20gZXZlciByZWFjaGluZyBHby4gQ2FuY2VsIG9ubHkgdGhlIGRhdGFcclxuICAgICAgICAvLyByZXRyeSBwYXRoOyB0aGUgY2hhaW4gYmVsb3cgc3RpbGwgb3JkZXJzIHRoZSBjbG9zZSBhZnRlciB0aGF0IHBhdGhcclxuICAgICAgICAvLyBoYXMgdW53b3VuZCwgcHJlc2VydmluZyBub3JtYWwgc2VuZC1iZWZvcmUtY2xvc2Ugb3JkZXJpbmcuXHJcbiAgICAgICAgdGhpcy5fc2VuZEFib3J0LmFib3J0KCk7XHJcbiAgICAgICAgdGhpcy5fY2hhaW4gPSB0aGlzLl9jaGFpbi50aGVuKCgpID0+XHJcbiAgICAgICAgICAgIHBvc3RGcmFtZSh0aGlzLl9pZCwgS0lORF9DTE9TRSwgbmV3IFVpbnQ4QXJyYXkoMCkpXHJcbiAgICAgICAgKS5jYXRjaCgoKSA9PiB7XHJcbiAgICAgICAgICAgIC8vIFRoZSBwZWVyIGlzIGFscmVhZHkgZ29uZTsgdGhlIGxvY2FsIGNsb3NlIGJlbG93IGlzIHdoYXQgbWF0dGVycy5cclxuICAgICAgICB9KS50aGVuKCgpID0+IHtcclxuICAgICAgICAgICAgdGhpcy5fY2xvc2VkKGNvZGUsIHJlYXNvbiwgdHJ1ZSk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqIEBpbnRlcm5hbCBDYWxsZWQgd2hlbiB0aGUgb3BlbiBmcmFtZSBjb21lcyBiYWNrIGZyb20gR28uICovXHJcbiAgICBfb3BlbmVkKCk6IHZvaWQge1xyXG4gICAgICAgIGlmICh0aGlzLnJlYWR5U3RhdGUgIT09IFdhaWxzU29ja2V0LkNPTk5FQ1RJTkcpIHJldHVybjtcclxuICAgICAgICB0aGlzLnJlYWR5U3RhdGUgPSBXYWlsc1NvY2tldC5PUEVOO1xyXG4gICAgICAgIHRoaXMuZGlzcGF0Y2hFdmVudChuZXcgRXZlbnQoXCJvcGVuXCIpKTtcclxuICAgIH1cclxuXHJcbiAgICAvKiogQGludGVybmFsIENhbGxlZCBmb3IgZWFjaCBkYXRhIGZyYW1lIGFkZHJlc3NlZCB0byB0aGlzIGNvbm5lY3Rpb24uICovXHJcbiAgICBfbWVzc2FnZShwYXlsb2FkOiBBcnJheUJ1ZmZlcik6IHZvaWQge1xyXG4gICAgICAgIGlmICh0aGlzLnJlYWR5U3RhdGUgIT09IFdhaWxzU29ja2V0Lk9QRU4pIHJldHVybjtcclxuICAgICAgICBsZXQgZGF0YTogdW5rbm93bjtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBkYXRhID0gdGhpcy5fZGVjb2RlKHBheWxvYWQpO1xyXG4gICAgICAgIH0gY2F0Y2gge1xyXG4gICAgICAgICAgICB0aGlzLmRpc3BhdGNoRXZlbnQobmV3IEV2ZW50KFwiZXJyb3JcIikpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGlmIChkYXRhID09PSBwYXlsb2FkICYmIHRoaXMuYmluYXJ5VHlwZSA9PT0gXCJibG9iXCIpIHtcclxuICAgICAgICAgICAgZGF0YSA9IG5ldyBCbG9iKFtwYXlsb2FkXSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHRoaXMuZGlzcGF0Y2hFdmVudChuZXcgTWVzc2FnZUV2ZW50KFwibWVzc2FnZVwiLCB7IGRhdGEgfSkpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKiBAaW50ZXJuYWwgKi9cclxuICAgIF9mYWlsKGVycjogdW5rbm93bik6IHZvaWQge1xyXG4gICAgICAgIC8vIEFuIGluLWZsaWdodCBzZW5kIGNhbiByZWplY3QgYWZ0ZXIgYSBjbGVhbiBjbG9zZSBoYXMgYWxyZWFkeSBhcnJpdmVkLlxyXG4gICAgICAgIC8vIFRoZSBzb2NrZXQncyBsaWZlY3ljbGUgaXMgY29tcGxldGUgdGhlbjogZW1pdHRpbmcgZXJyb3IgYWZ0ZXIgY2xvc2VcclxuICAgICAgICAvLyByZXZlcnNlcyB0aGUgV2ViU29ja2V0IGV2ZW50IG9yZGVyIGFuZCByZXBvcnRzIGEgZmFsc2UgZmFpbHVyZS5cclxuICAgICAgICBpZiAodGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DTE9TRUQgfHwgdGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DTE9TSU5HKSByZXR1cm47XHJcbiAgICAgICAgdGhpcy5kaXNwYXRjaEV2ZW50KG5ldyBFdmVudChcImVycm9yXCIpKTtcclxuICAgICAgICAvLyAxMDA2OiBjbG9zZWQgYWJub3JtYWxseSwgbm8gY2xvc2UgZnJhbWUuIFNhbWUgY29kZSBhIHJlYWwgV2ViU29ja2V0XHJcbiAgICAgICAgLy8gcmVwb3J0cyB3aGVuIHRoZSBjb25uZWN0aW9uIGRyb3BzIHdpdGhvdXQgYSBoYW5kc2hha2UuXHJcbiAgICAgICAgdGhpcy5fY2xvc2VkKDEwMDYsIGVyciBpbnN0YW5jZW9mIEVycm9yID8gZXJyLm1lc3NhZ2UgOiBTdHJpbmcoZXJyKSwgZmFsc2UpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKiBAaW50ZXJuYWwgKi9cclxuICAgIF9jbG9zZWQoY29kZTogbnVtYmVyLCByZWFzb246IHN0cmluZywgd2FzQ2xlYW46IGJvb2xlYW4pOiB2b2lkIHtcclxuICAgICAgICBpZiAodGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DTE9TRUQpIHJldHVybjtcclxuICAgICAgICB0aGlzLnJlYWR5U3RhdGUgPSBXYWlsc1NvY2tldC5DTE9TRUQ7XHJcbiAgICAgICAgdGhpcy5fb3BlbkFib3J0LmFib3J0KCk7XHJcbiAgICAgICAgdGhpcy5fc2VuZEFib3J0LmFib3J0KCk7XHJcbiAgICAgICAgRy5jb25uZWN0aW9ucy5kZWxldGUodGhpcy5faWQpO1xyXG4gICAgICAgIGlmIChHLmNvbm5lY3Rpb25zLnNpemUgPT09IDApIEcucG9sbEFib3J0Py5hYm9ydCgpO1xyXG4gICAgICAgIHRoaXMuX3BlbmRpbmcubGVuZ3RoID0gMDtcclxuICAgICAgICB0aGlzLl9idWZmZXJlZCA9IDA7XHJcblxyXG4gICAgICAgIGNvbnN0IGV2ID0gdHlwZW9mIENsb3NlRXZlbnQgPT09IFwiZnVuY3Rpb25cIlxyXG4gICAgICAgICAgICA/IG5ldyBDbG9zZUV2ZW50KFwiY2xvc2VcIiwgeyBjb2RlLCByZWFzb24sIHdhc0NsZWFuIH0pXHJcbiAgICAgICAgICAgIDogT2JqZWN0LmFzc2lnbihuZXcgRXZlbnQoXCJjbG9zZVwiKSwgeyBjb2RlLCByZWFzb24sIHdhc0NsZWFuIH0pIGFzIENsb3NlRXZlbnQ7XHJcbiAgICAgICAgdGhpcy5kaXNwYXRjaEV2ZW50KGV2KTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENvbm5lY3QgdG8gYSBzdHJlYW0gZGVjbGFyZWQgaW4gR28gd2l0aCBgYXBwLkhhbmRsZVN0cmVhbShuYW1lLCBoYW5kbGVyKWAuXHJcbiAqXHJcbiAqIFJldHVybnMgaW1tZWRpYXRlbHkgd2l0aCBgcmVhZHlTdGF0ZSA9PT0gQ09OTkVDVElOR2A7IGxpc3RlbiBmb3IgYG9wZW5gLCBvclxyXG4gKiBqdXN0IGNhbGwge0BsaW5rIFdhaWxzU29ja2V0LnNlbmR9IG9uY2Ugb3Blbi4gQ3JlYXRpbmcgYSBjb25uZWN0aW9uIGF0IG1vZHVsZVxyXG4gKiBzY29wZSBpcyBzdXBwb3J0ZWQgYW5kIGlzIHRoZSBpbnRlbmRlZCBzaGFwZSBmb3IgZ2VuZXJhdGVkIGJpbmRpbmdzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFN0cmVhbShuYW1lOiBzdHJpbmcpOiBXYWlsc1NvY2tldCB8IFdlYlNvY2tldCB7XHJcbiAgICAvLyBTZXJ2ZXIgYnVpbGRzIGluc3RhbGwgYSBmYWN0b3J5IHJldHVybmluZyBhIHJlYWwgV2ViU29ja2V0LCBzaW5jZSB0aGVyZVxyXG4gICAgLy8gaXMgYSBsaXN0ZW5lciB0aGVyZSB0byB1cGdyYWRlIG9uIGFuZCBub3RoaW5nIHRvIGVtdWxhdGUuIEJvdGggb2JqZWN0c1xyXG4gICAgLy8gcHJlc2VudCB0aGUgc2FtZSBpbnRlcmZhY2UsIHNvIG5vdGhpbmcgYWJvdmUgdGhpcyBsaW5lIGNhcmVzIHdoaWNoIGl0IGlzLlxyXG4gICAgaWYgKGhhc0RPTSkge1xyXG4gICAgICAgIGNvbnN0IGZhY3RvcnkgPSAod2luZG93IGFzIGFueSkuX3dhaWxzPy5zdHJlYW1GYWN0b3J5O1xyXG4gICAgICAgIGlmICh0eXBlb2YgZmFjdG9yeSA9PT0gXCJmdW5jdGlvblwiKSB7XHJcbiAgICAgICAgICAgIHJldHVybiBmYWN0b3J5KG5hbWUpO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuICAgIHJldHVybiBuZXcgV2FpbHNTb2NrZXQobmFtZSk7XHJcbn1cclxuXHJcbi8vIERlZmluZXMgYW4gYG9uPHR5cGU+YCBwcm9wZXJ0eSB0aGF0IHJlZ2lzdGVycyBhbmQgdW5yZWdpc3RlcnMgYSByZWFsIGV2ZW50XHJcbi8vIGxpc3RlbmVyLCB0aGUgd2F5IHRoZSBET00gZGVmaW5lcyB0aGVtLiBBc3NpZ25pbmcgcmVwbGFjZXMgdGhlIHByZXZpb3VzXHJcbi8vIGhhbmRsZXI7IHJlYWRpbmcgcmV0dXJucyBpdC5cclxuZnVuY3Rpb24gZGVmaW5lSGFuZGxlclByb3BlcnR5KHRhcmdldDogRXZlbnRUYXJnZXQsIHR5cGU6IHN0cmluZyk6IHZvaWQge1xyXG4gICAgbGV0IGN1cnJlbnQ6ICgoZXY6IGFueSkgPT4gdm9pZCkgfCBudWxsID0gbnVsbDtcclxuICAgIE9iamVjdC5kZWZpbmVQcm9wZXJ0eSh0YXJnZXQsIFwib25cIiArIHR5cGUsIHtcclxuICAgICAgICBnZXQ6ICgpID0+IGN1cnJlbnQsXHJcbiAgICAgICAgc2V0KGZuOiAoKGV2OiBhbnkpID0+IHZvaWQpIHwgbnVsbCkge1xyXG4gICAgICAgICAgICBpZiAoY3VycmVudCkgdGFyZ2V0LnJlbW92ZUV2ZW50TGlzdGVuZXIodHlwZSwgY3VycmVudCk7XHJcbiAgICAgICAgICAgIGN1cnJlbnQgPSB0eXBlb2YgZm4gPT09IFwiZnVuY3Rpb25cIiA/IGZuIDogbnVsbDtcclxuICAgICAgICAgICAgaWYgKGN1cnJlbnQpIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKHR5cGUsIGN1cnJlbnQpO1xyXG4gICAgICAgIH0sXHJcbiAgICAgICAgY29uZmlndXJhYmxlOiB0cnVlLFxyXG4gICAgICAgIGVudW1lcmFibGU6IHRydWUsXHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFRoZSBvYmplY3Qtc2hhcGVkIHZpZXcgb2YgYSBzdHJlYW0gcmV0dXJuZWQgYnkge0BsaW5rIEpTT05TdHJlYW19OiBgc2VuZGBcclxuICogdGFrZXMgYSB2YWx1ZSB0byBzdHJpbmdpZnksIGFuZCBgZXYuZGF0YWAgb24gYSBtZXNzYWdlIGlzIHRoZSBwYXJzZWQgcmVzdWx0LlxyXG4gKi9cclxuZXhwb3J0IGludGVyZmFjZSBKU09OU29ja2V0IGV4dGVuZHMgT21pdDxXYWlsc1NvY2tldCwgXCJzZW5kXCIgfCBcIm9ubWVzc2FnZVwiIHwgXCJhZGRFdmVudExpc3RlbmVyXCIgfCBcInJlbW92ZUV2ZW50TGlzdGVuZXJcIj4ge1xyXG4gICAgc2VuZCh2YWx1ZTogdW5rbm93bik6IHZvaWQ7XHJcbiAgICBvbm1lc3NhZ2U6ICgoZXY6IE1lc3NhZ2VFdmVudDxhbnk+KSA9PiB2b2lkKSB8IG51bGw7XHJcbiAgICBhZGRFdmVudExpc3RlbmVyKFxyXG4gICAgICAgIHR5cGU6IFwibWVzc2FnZVwiLFxyXG4gICAgICAgIGxpc3RlbmVyOiAoKHRoaXM6IEpTT05Tb2NrZXQsIGV2OiBNZXNzYWdlRXZlbnQ8YW55PikgPT4gYW55KSB8IEV2ZW50TGlzdGVuZXJPYmplY3QgfCBudWxsLFxyXG4gICAgICAgIG9wdGlvbnM/OiBib29sZWFuIHwgQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMsXHJcbiAgICApOiB2b2lkO1xyXG4gICAgYWRkRXZlbnRMaXN0ZW5lcihcclxuICAgICAgICB0eXBlOiBzdHJpbmcsXHJcbiAgICAgICAgbGlzdGVuZXI6IEV2ZW50TGlzdGVuZXJPckV2ZW50TGlzdGVuZXJPYmplY3QgfCBudWxsLFxyXG4gICAgICAgIG9wdGlvbnM/OiBib29sZWFuIHwgQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMsXHJcbiAgICApOiB2b2lkO1xyXG4gICAgcmVtb3ZlRXZlbnRMaXN0ZW5lcihcclxuICAgICAgICB0eXBlOiBcIm1lc3NhZ2VcIixcclxuICAgICAgICBsaXN0ZW5lcjogKCh0aGlzOiBKU09OU29ja2V0LCBldjogTWVzc2FnZUV2ZW50PGFueT4pID0+IGFueSkgfCBFdmVudExpc3RlbmVyT2JqZWN0IHwgbnVsbCxcclxuICAgICAgICBvcHRpb25zPzogYm9vbGVhbiB8IEV2ZW50TGlzdGVuZXJPcHRpb25zLFxyXG4gICAgKTogdm9pZDtcclxuICAgIHJlbW92ZUV2ZW50TGlzdGVuZXIoXHJcbiAgICAgICAgdHlwZTogc3RyaW5nLFxyXG4gICAgICAgIGxpc3RlbmVyOiBFdmVudExpc3RlbmVyT3JFdmVudExpc3RlbmVyT2JqZWN0IHwgbnVsbCxcclxuICAgICAgICBvcHRpb25zPzogYm9vbGVhbiB8IEV2ZW50TGlzdGVuZXJPcHRpb25zLFxyXG4gICAgKTogdm9pZDtcclxufVxyXG5cclxuLyoqXHJcbiAqIEEge0BsaW5rIFN0cmVhbX0gdGhhdCBzcGVha3Mgb2JqZWN0cyBpbnN0ZWFkIG9mIGJ5dGVzLlxyXG4gKlxyXG4gKiBgc2VuZCh2YWx1ZSlgIG1hcnNoYWxzIHdpdGggYEpTT04uc3RyaW5naWZ5YCwgYW5kIGBldi5kYXRhYCBpbiBhbiBgb25tZXNzYWdlYFxyXG4gKiBoYW5kbGVyIGlzIHRoZSBwYXJzZWQgb2JqZWN0IHJhdGhlciB0aGFuIGFuIGBBcnJheUJ1ZmZlcmAuIFBhaXJzIHdpdGhcclxuICogYFN0cmVhbUNvbm4uU2VuZEpTT05gIC8gYFJlY2VpdmVKU09OYCBvbiB0aGUgR28gc2lkZS5cclxuICpcclxuICogYGBganNcclxuICogY29uc3QgcyA9IEpTT05TdHJlYW0oXCJ0ZWxlbWV0cnlcIik7XHJcbiAqIHMub25tZXNzYWdlID0gKGV2KSA9PiBjb25zb2xlLmxvZyhldi5kYXRhLnRlbXBlcmF0dXJlKTtcclxuICogcy5vbm9wZW4gICAgPSAoKSA9PiBzLnNlbmQoeyBzdWJzY3JpYmU6IFwic2Vuc29yc1wiIH0pO1xyXG4gKiBgYGBcclxuICpcclxuICogQSBmcmFtZSB0aGF0IGlzIG5vdCB2YWxpZCBKU09OIHJhaXNlcyBhbiBgZXJyb3JgIGV2ZW50IGFuZCBpcyBkcm9wcGVkLCByYXRoZXJcclxuICogdGhhbiB0aHJvd2luZyBmcm9tIHRoZSBwb2xsIGxvb3AgYW5kIHRha2luZyB0aGUgY29ubmVjdGlvbiBkb3duIHdpdGggaXQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSlNPTlN0cmVhbShuYW1lOiBzdHJpbmcpOiBKU09OU29ja2V0IHtcclxuICAgIGNvbnN0IHNvY2tldCA9IFN0cmVhbShuYW1lKTtcclxuICAgIGNvbnN0IGRlY29kZXIgPSBuZXcgVGV4dERlY29kZXIoKTtcclxuICAgIGNvbnN0IHBhcnNlID0gKHBheWxvYWQ6IEFycmF5QnVmZmVyIHwgc3RyaW5nKSA9PlxyXG4gICAgICAgIEpTT04ucGFyc2UodHlwZW9mIHBheWxvYWQgPT09IFwic3RyaW5nXCIgPyBwYXlsb2FkIDogZGVjb2Rlci5kZWNvZGUocGF5bG9hZCkpO1xyXG5cclxuICAgIGlmIChzb2NrZXQgaW5zdGFuY2VvZiBXYWlsc1NvY2tldCkge1xyXG4gICAgICAgIC8vIERlY29kaW5nIHRocm91Z2ggdGhlIGhvb2sgcmF0aGVyIHRoYW4gYnkgd3JhcHBpbmcgb25tZXNzYWdlIG1lYW5zXHJcbiAgICAgICAgLy8gYWRkRXZlbnRMaXN0ZW5lciBsaXN0ZW5lcnMgc2VlIHRoZSBwYXJzZWQgb2JqZWN0IHRvby5cclxuICAgICAgICBzb2NrZXQuX2RlY29kZSA9IChwYXlsb2FkKSA9PiBwYXJzZShwYXlsb2FkKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgLy8gU2VydmVyIGJ1aWxkcyBnZXQgYSBuYXRpdmUgV2ViU29ja2V0LCB3aG9zZSBkaXNwYXRjaCBuZXZlciBjb25zdWx0c1xyXG4gICAgICAgIC8vIF9kZWNvZGUuIFBhdGNoIGJvdGggbGlzdGVuZXIgc3R5bGVzIG9uIHRoaXMgaW5zdGFuY2Ugc28gSlNPTlN0cmVhbVxyXG4gICAgICAgIC8vIG1lYW5zIHRoZSBzYW1lIHRoaW5nIG9uIGVpdGhlciB0cmFuc3BvcnQsIHdoaWNoIGlzIHRoZSBwb2ludCBvZiB0aGVcclxuICAgICAgICAvLyBzaGFyZWQgQVBJLlxyXG4gICAgICAgIGNvbnN0IG5hdGl2ZSA9IHNvY2tldCBhcyBXZWJTb2NrZXQ7XHJcbiAgICAgICAgbmF0aXZlLmJpbmFyeVR5cGUgPSBcImFycmF5YnVmZmVyXCI7XHJcblxyXG4gICAgICAgIGNvbnN0IHdyYXBwZWQgPSBuZXcgV2Vha01hcDxvYmplY3QsIEV2ZW50TGlzdGVuZXI+KCk7XHJcbiAgICAgICAgY29uc3QgZGVjb2RlZEV2ZW50cyA9IG5ldyBXZWFrTWFwPEV2ZW50LCBNZXNzYWdlRXZlbnQgfCBudWxsPigpO1xyXG4gICAgICAgIGNvbnN0IGFkZCA9IG5hdGl2ZS5hZGRFdmVudExpc3RlbmVyLmJpbmQobmF0aXZlKTtcclxuICAgICAgICBjb25zdCByZW1vdmUgPSBuYXRpdmUucmVtb3ZlRXZlbnRMaXN0ZW5lci5iaW5kKG5hdGl2ZSk7XHJcblxyXG4gICAgICAgIGNvbnN0IGRlY29kZUV2ZW50ID0gKGV2OiBFdmVudCk6IE1lc3NhZ2VFdmVudCB8IG51bGwgPT4ge1xyXG4gICAgICAgICAgICBpZiAoZGVjb2RlZEV2ZW50cy5oYXMoZXYpKSByZXR1cm4gZGVjb2RlZEV2ZW50cy5nZXQoZXYpID8/IG51bGw7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICBjb25zdCB2YWx1ZSA9IHBhcnNlKChldiBhcyBNZXNzYWdlRXZlbnQpLmRhdGEpO1xyXG4gICAgICAgICAgICAgICAgY29uc3QgZGVjb2RlZCA9IG5ldyBNZXNzYWdlRXZlbnQoXCJtZXNzYWdlXCIsIHsgZGF0YTogdmFsdWUgfSk7XHJcbiAgICAgICAgICAgICAgICBkZWNvZGVkRXZlbnRzLnNldChldiwgZGVjb2RlZCk7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gZGVjb2RlZDtcclxuICAgICAgICAgICAgfSBjYXRjaCB7XHJcbiAgICAgICAgICAgICAgICBkZWNvZGVkRXZlbnRzLnNldChldiwgbnVsbCk7XHJcbiAgICAgICAgICAgICAgICBuYXRpdmUuZGlzcGF0Y2hFdmVudChuZXcgRXZlbnQoXCJlcnJvclwiKSk7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gbnVsbDtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIGNvbnN0IGRlY29kZSA9IChsaXN0ZW5lcjogYW55KTogRXZlbnRMaXN0ZW5lciA9PiB7XHJcbiAgICAgICAgICAgIGNvbnN0IGV4aXN0aW5nID0gd3JhcHBlZC5nZXQobGlzdGVuZXIgYXMgb2JqZWN0KTtcclxuICAgICAgICAgICAgaWYgKGV4aXN0aW5nKSByZXR1cm4gZXhpc3Rpbmc7XHJcblxyXG4gICAgICAgICAgICBjb25zdCBmbjogRXZlbnRMaXN0ZW5lciA9IChldikgPT4ge1xyXG4gICAgICAgICAgICAgICAgY29uc3QgZGVjb2RlZCA9IGRlY29kZUV2ZW50KGV2KTtcclxuICAgICAgICAgICAgICAgIGlmICghZGVjb2RlZCkgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgaWYgKHR5cGVvZiBsaXN0ZW5lciA9PT0gXCJmdW5jdGlvblwiKSBsaXN0ZW5lci5jYWxsKG5hdGl2ZSwgZGVjb2RlZCk7XHJcbiAgICAgICAgICAgICAgICBlbHNlIGxpc3RlbmVyLmhhbmRsZUV2ZW50KGRlY29kZWQpO1xyXG4gICAgICAgICAgICB9O1xyXG4gICAgICAgICAgICB3cmFwcGVkLnNldChsaXN0ZW5lciBhcyBvYmplY3QsIGZuKTtcclxuICAgICAgICAgICAgcmV0dXJuIGZuO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIG5hdGl2ZS5hZGRFdmVudExpc3RlbmVyID0gKCh0eXBlOiBzdHJpbmcsIGxpc3RlbmVyOiBhbnksIG9wdHM/OiBhbnkpID0+IHtcclxuICAgICAgICAgICAgYWRkKHR5cGUsIHR5cGUgPT09IFwibWVzc2FnZVwiICYmIGxpc3RlbmVyID8gZGVjb2RlKGxpc3RlbmVyKSA6IGxpc3RlbmVyLCBvcHRzKTtcclxuICAgICAgICB9KSBhcyB0eXBlb2YgbmF0aXZlLmFkZEV2ZW50TGlzdGVuZXI7XHJcblxyXG4gICAgICAgIG5hdGl2ZS5yZW1vdmVFdmVudExpc3RlbmVyID0gKCh0eXBlOiBzdHJpbmcsIGxpc3RlbmVyOiBhbnksIG9wdHM/OiBhbnkpID0+IHtcclxuICAgICAgICAgICAgY29uc3QgZm4gPSB0eXBlID09PSBcIm1lc3NhZ2VcIiAmJiBsaXN0ZW5lciA/IHdyYXBwZWQuZ2V0KGxpc3RlbmVyKSA6IHVuZGVmaW5lZDtcclxuICAgICAgICAgICAgcmVtb3ZlKHR5cGUsIGZuID8/IGxpc3RlbmVyLCBvcHRzKTtcclxuICAgICAgICB9KSBhcyB0eXBlb2YgbmF0aXZlLnJlbW92ZUV2ZW50TGlzdGVuZXI7XHJcblxyXG4gICAgICAgIGxldCBoYW5kbGVyOiAoKGV2OiBNZXNzYWdlRXZlbnQpID0+IHZvaWQpIHwgbnVsbCA9IG51bGw7XHJcbiAgICAgICAgT2JqZWN0LmRlZmluZVByb3BlcnR5KG5hdGl2ZSwgXCJvbm1lc3NhZ2VcIiwge1xyXG4gICAgICAgICAgICBnZXQ6ICgpID0+IGhhbmRsZXIsXHJcbiAgICAgICAgICAgIHNldChmbikge1xyXG4gICAgICAgICAgICAgICAgaWYgKGhhbmRsZXIpIG5hdGl2ZS5yZW1vdmVFdmVudExpc3RlbmVyKFwibWVzc2FnZVwiLCBoYW5kbGVyIGFzIEV2ZW50TGlzdGVuZXIpO1xyXG4gICAgICAgICAgICAgICAgaGFuZGxlciA9IHR5cGVvZiBmbiA9PT0gXCJmdW5jdGlvblwiID8gZm4gOiBudWxsO1xyXG4gICAgICAgICAgICAgICAgaWYgKGhhbmRsZXIpIG5hdGl2ZS5hZGRFdmVudExpc3RlbmVyKFwibWVzc2FnZVwiLCBoYW5kbGVyIGFzIEV2ZW50TGlzdGVuZXIpO1xyXG4gICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICBjb25maWd1cmFibGU6IHRydWUsXHJcbiAgICAgICAgICAgIGVudW1lcmFibGU6IHRydWUsXHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3Qgc2VuZCA9IHNvY2tldC5zZW5kLmJpbmQoc29ja2V0KTtcclxuICAgIChzb2NrZXQgYXMgdW5rbm93biBhcyBKU09OU29ja2V0KS5zZW5kID0gKHZhbHVlOiB1bmtub3duKSA9PiBzZW5kKEpTT04uc3RyaW5naWZ5KHZhbHVlKSk7XHJcblxyXG4gICAgcmV0dXJuIHNvY2tldCBhcyB1bmtub3duIGFzIEpTT05Tb2NrZXQ7XHJcbn1cclxuXHJcbi8vIFN5bmNocm9ub3VzIGNvbnZlcnNpb24gZm9yIGV2ZXJ5dGhpbmcgYnV0IGEgQmxvYi4gQXZvaWRzIGEgcHJvbWlzZSBwZXIgc2VuZFxyXG4vLyBhbmQsIG1vcmUgaW1wb3J0YW50bHksIGxldHMgYnVmZmVyZWRBbW91bnQgYmUgZXhhY3QgdGhlIG1vbWVudCBzZW5kKCkgcmV0dXJucy5cclxuZnVuY3Rpb24gdG9CeXRlc1N5bmMoZGF0YTogc3RyaW5nIHwgQXJyYXlCdWZmZXJMaWtlIHwgQXJyYXlCdWZmZXJWaWV3IHwgQmxvYik6IFVpbnQ4QXJyYXkgfCBudWxsIHtcclxuICAgIGlmICh0eXBlb2YgZGF0YSA9PT0gXCJzdHJpbmdcIikge1xyXG4gICAgICAgIHJldHVybiB0ZXh0RW5jb2Rlci5lbmNvZGUoZGF0YSk7XHJcbiAgICB9XHJcbiAgICBpZiAodHlwZW9mIEJsb2IgIT09IFwidW5kZWZpbmVkXCIgJiYgZGF0YSBpbnN0YW5jZW9mIEJsb2IpIHtcclxuICAgICAgICByZXR1cm4gbnVsbDtcclxuICAgIH1cclxuICAgIGlmIChBcnJheUJ1ZmZlci5pc1ZpZXcoZGF0YSkpIHtcclxuICAgICAgICByZXR1cm4gbmV3IFVpbnQ4QXJyYXkoZGF0YS5idWZmZXIsIGRhdGEuYnl0ZU9mZnNldCwgZGF0YS5ieXRlTGVuZ3RoKS5zbGljZSgpO1xyXG4gICAgfVxyXG4gICAgcmV0dXJuIG5ldyBVaW50OEFycmF5KGRhdGEgYXMgQXJyYXlCdWZmZXJMaWtlKS5zbGljZSgpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBhYm9ydGFibGU8VD4ocHJvbWlzZTogUHJvbWlzZTxUPiwgc2lnbmFsPzogQWJvcnRTaWduYWwpOiBQcm9taXNlPFQ+IHtcclxuICAgIGlmICghc2lnbmFsKSByZXR1cm4gcHJvbWlzZTtcclxuICAgIGlmIChzaWduYWwuYWJvcnRlZCkgcmV0dXJuIFByb21pc2UucmVqZWN0KG5ldyBTdHJlYW1SZXF1ZXN0Q2FuY2VsbGVkKCkpO1xyXG5cclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgY29uc3Qgb25BYm9ydCA9ICgpID0+IHJlamVjdChuZXcgU3RyZWFtUmVxdWVzdENhbmNlbGxlZCgpKTtcclxuICAgICAgICBzaWduYWwuYWRkRXZlbnRMaXN0ZW5lcihcImFib3J0XCIsIG9uQWJvcnQsIHsgb25jZTogdHJ1ZSB9KTtcclxuICAgICAgICBwcm9taXNlLnRoZW4oXHJcbiAgICAgICAgICAgICh2YWx1ZSkgPT4ge1xyXG4gICAgICAgICAgICAgICAgc2lnbmFsLnJlbW92ZUV2ZW50TGlzdGVuZXIoXCJhYm9ydFwiLCBvbkFib3J0KTtcclxuICAgICAgICAgICAgICAgIHJlc29sdmUodmFsdWUpO1xyXG4gICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICAoZXJyKSA9PiB7XHJcbiAgICAgICAgICAgICAgICBzaWduYWwucmVtb3ZlRXZlbnRMaXN0ZW5lcihcImFib3J0XCIsIG9uQWJvcnQpO1xyXG4gICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XHJcbiAgICAgICAgICAgIH0sXHJcbiAgICAgICAgKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5hc3luYyBmdW5jdGlvbiB0b0J5dGVzKFxyXG4gICAgZGF0YTogc3RyaW5nIHwgQXJyYXlCdWZmZXJMaWtlIHwgQXJyYXlCdWZmZXJWaWV3IHwgQmxvYixcclxuICAgIHNpZ25hbD86IEFib3J0U2lnbmFsLFxyXG4pOiBQcm9taXNlPFVpbnQ4QXJyYXk+IHtcclxuICAgIGlmICh0eXBlb2YgZGF0YSA9PT0gXCJzdHJpbmdcIikge1xyXG4gICAgICAgIHJldHVybiB0ZXh0RW5jb2Rlci5lbmNvZGUoZGF0YSk7XHJcbiAgICB9XHJcbiAgICBpZiAodHlwZW9mIEJsb2IgIT09IFwidW5kZWZpbmVkXCIgJiYgZGF0YSBpbnN0YW5jZW9mIEJsb2IpIHtcclxuICAgICAgICByZXR1cm4gbmV3IFVpbnQ4QXJyYXkoYXdhaXQgYWJvcnRhYmxlKGRhdGEuYXJyYXlCdWZmZXIoKSwgc2lnbmFsKSk7XHJcbiAgICB9XHJcbiAgICBpZiAoQXJyYXlCdWZmZXIuaXNWaWV3KGRhdGEpKSB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBVaW50OEFycmF5KGRhdGEuYnVmZmVyLCBkYXRhLmJ5dGVPZmZzZXQsIGRhdGEuYnl0ZUxlbmd0aCkuc2xpY2UoKTtcclxuICAgIH1cclxuICAgIHJldHVybiBuZXcgVWludDhBcnJheShkYXRhIGFzIEFycmF5QnVmZmVyTGlrZSkuc2xpY2UoKTtcclxufVxyXG5cclxuYXN5bmMgZnVuY3Rpb24gcG9zdEZyYW1lKGNvbm5JRDogbnVtYmVyLCBraW5kOiBudW1iZXIsIGJvZHk6IFVpbnQ4QXJyYXksIG5hbWU/OiBzdHJpbmcsIHNpZ25hbD86IEFib3J0U2lnbmFsKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICBjb25zdCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+ID0ge1xyXG4gICAgICAgIFtIRFJfU0VTU0lPTl06IEcuc2Vzc2lvbixcclxuICAgICAgICBbSERSX0dFTkVSQVRJT05dOiBTdHJpbmcoRy5nZW5lcmF0aW9uKSxcclxuICAgICAgICBbSERSX0NPTk5dOiBTdHJpbmcoY29ubklEKSxcclxuICAgICAgICBbSERSX0tJTkRdOiBTdHJpbmcoa2luZCksXHJcbiAgICAgICAgXCJDb250ZW50LVR5cGVcIjogXCJhcHBsaWNhdGlvbi9vY3RldC1zdHJlYW1cIixcclxuICAgIH07XHJcbiAgICBpZiAobmFtZSAhPT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgaGVhZGVyc1tIRFJfTkFNRV0gPSBuYW1lO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChib2R5LmJ5dGVMZW5ndGggPD0gQ0hVTktfVEhSRVNIT0xEKSB7XHJcbiAgICAgICAgYXdhaXQgcG9zdFdpdGhSZXRyeShoZWFkZXJzLCBib2R5LCBzaWduYWwpO1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBTcGxpdCBvdmVyc2l6ZWQgZnJhbWVzIHRoZSBzYW1lIHdheSB0aGUgcnVudGltZSBzcGxpdHMgb3ZlcnNpemVkIGNhbGxzLlxyXG4gICAgY29uc3QgY2h1bmtJRCA9IG5hbm9pZCgpO1xyXG4gICAgY29uc3QgdG90YWwgPSBNYXRoLmNlaWwoYm9keS5ieXRlTGVuZ3RoIC8gQ0hVTktfVEhSRVNIT0xEKTtcclxuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgdG90YWw7IGkrKykge1xyXG4gICAgICAgIGNvbnN0IHNsaWNlID0gYm9keS5zdWJhcnJheShpICogQ0hVTktfVEhSRVNIT0xELCAoaSArIDEpICogQ0hVTktfVEhSRVNIT0xEKTtcclxuICAgICAgICBhd2FpdCBwb3N0V2l0aFJldHJ5KHtcclxuICAgICAgICAgICAgLi4uaGVhZGVycyxcclxuICAgICAgICAgICAgW0hEUl9DSFVOS106IGNodW5rSUQsXHJcbiAgICAgICAgICAgIFtIRFJfQ0hVTktfSU5ERVhdOiBTdHJpbmcoaSksXHJcbiAgICAgICAgICAgIFtIRFJfQ0hVTktfVE9UQUxdOiBTdHJpbmcodG90YWwpLFxyXG4gICAgICAgIH0sIHNsaWNlLCBzaWduYWwpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBBIDQyOSBtZWFucyB0aGUgR28gaGFuZGxlciBoYXMgbm90IHRha2VuIGRlbGl2ZXJ5IG9mIHdoYXQgaXMgYWxyZWFkeSBxdWV1ZWQuXHJcbi8vIFJldHJ5aW5nIHRoZSBzYW1lIGZyYW1lIGtlZXBzIG9yZGVyaW5nICh0aGUgc2VuZCBjaGFpbiBoYXMgbm90IG1vdmVkIG9uKSBhbmRcclxuLy8gbGVhdmVzIHRoZSByZXF1ZXN0IHNsb3QgZnJlZSwgd2hpY2ggaXMgd2hhdCBzdG9wcyBhIGJhY2tlZC11cCBjb25uZWN0aW9uIGZyb21cclxuLy8gc3RhcnZpbmcgdGhlIHdpbmRvdydzIHBvbGwuXHJcbmFzeW5jIGZ1bmN0aW9uIHBvc3RXaXRoUmV0cnkoaGVhZGVyczogUmVjb3JkPHN0cmluZywgc3RyaW5nPiwgYm9keTogVWludDhBcnJheSwgc2lnbmFsPzogQWJvcnRTaWduYWwpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIC8vIE5ldmVyIHJldHJ5IHdpdGggemVybyBkZWxheS4gVGhlIHJlY2VpdmVyIGJlaW5nIGJlaGluZCBpcyBhIGNvbmRpdGlvblxyXG4gICAgLy8gdGhhdCB0YWtlcyB0aW1lIHRvIGNsZWFyLCBhbmQgYW4gaW1tZWRpYXRlIHJldHJ5IGlzIGEgYnVzeSBsb29wIG9mXHJcbiAgICAvLyBmZXRjaGVzIOKAlCBlYWNoIG9uZSBhIHNjaGVtZS1oYW5kbGVyIHJvdW5kIHRyaXAsIGFuZCBvbiB0aGUgaG9zdCBhIGNnb1xyXG4gICAgLy8gY2FsbC4gU3RhcnQgYXQgMSBtcyBhbmQgcmFtcC5cclxuICAgIGxldCB3YWl0ID0gMTtcclxuICAgIGZvciAoOzspIHtcclxuICAgICAgICBpZiAoc2lnbmFsPy5hYm9ydGVkKSB0aHJvdyBuZXcgU3RyZWFtUmVxdWVzdENhbmNlbGxlZCgpO1xyXG4gICAgICAgIGNvbnN0IHJlc3AgPSBhd2FpdCBmZXRjaChzdHJlYW1VUkwoXCJzZW5kXCIpLCB7XHJcbiAgICAgICAgICAgIG1ldGhvZDogXCJQT1NUXCIsXHJcbiAgICAgICAgICAgIGhlYWRlcnMsXHJcbiAgICAgICAgICAgIGJvZHk6IGJvZHkgYXMgQm9keUluaXQsXHJcbiAgICAgICAgICAgIHNpZ25hbCxcclxuICAgICAgICB9KTtcclxuICAgICAgICBpZiAocmVzcC5vaykgcmV0dXJuO1xyXG4gICAgICAgIGlmIChyZXNwLnN0YXR1cyAhPT0gNDI5KSB7XHJcbiAgICAgICAgICAgIHRocm93IG5ldyBFcnJvcihhd2FpdCByZXNwLnRleHQoKSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGF3YWl0IHNsZWVwKHdhaXQsIHNpZ25hbCk7XHJcbiAgICAgICAgd2FpdCA9IE1hdGgubWluKHdhaXQgKiAyLCA1MCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIHBvc3RCYXRjaCBzZW5kcyBzZXZlcmFsIGRhdGEgZnJhbWVzIGZvciBvbmUgY29ubmVjdGlvbiBpbiBib3VuZGVkIHJlcXVlc3RzLlxyXG4vLyBCb2R5OiBjb3VudCB1MzIsIHRoZW4gY291bnQgeCAoIGxlbiB1MzIgfCBwYXlsb2FkICkuIFRoZSBjb21iaW5lZCBib2R5IHN0YXlzXHJcbi8vIHdpdGhpbiBDSFVOS19USFJFU0hPTEQ6IHNldmVyYWwgaW5kaXZpZHVhbGx5IHNtYWxsIGZyYW1lcyBtdXN0IG5vdCByZWNyZWF0ZVxyXG4vLyB0aGUgV2ViVmlldzIgYm9keS1zaXplIHByb2JsZW0gdGhhdCBjaHVua2luZyBzaW5nbGUgbGFyZ2UgZnJhbWVzIGF2b2lkcy5cclxuYXN5bmMgZnVuY3Rpb24gcG9zdEJhdGNoKGNvbm5JRDogbnVtYmVyLCBmcmFtZXM6IFVpbnQ4QXJyYXlbXSwgc2lnbmFsPzogQWJvcnRTaWduYWwpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIGxldCBzdGFydCA9IDA7XHJcbiAgICB3aGlsZSAoc3RhcnQgPCBmcmFtZXMubGVuZ3RoKSB7XHJcbiAgICAgICAgaWYgKGZyYW1lc1tzdGFydF0uYnl0ZUxlbmd0aCA+IENIVU5LX1RIUkVTSE9MRCkge1xyXG4gICAgICAgICAgICBhd2FpdCBwb3N0RnJhbWUoY29ubklELCBLSU5EX0RBVEEsIGZyYW1lc1tzdGFydF0sIHVuZGVmaW5lZCwgc2lnbmFsKTtcclxuICAgICAgICAgICAgc3RhcnQrKztcclxuICAgICAgICAgICAgY29udGludWU7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBsZXQgZW5kID0gc3RhcnQ7XHJcbiAgICAgICAgbGV0IHNpemUgPSA0OyAvLyBmcmFtZSBjb3VudFxyXG4gICAgICAgIHdoaWxlIChlbmQgPCBmcmFtZXMubGVuZ3RoICYmIGVuZCAtIHN0YXJ0IDwgTUFYX0JBVENIX0ZSQU1FUykge1xyXG4gICAgICAgICAgICBjb25zdCBmcmFtZSA9IGZyYW1lc1tlbmRdO1xyXG4gICAgICAgICAgICBpZiAoZnJhbWUuYnl0ZUxlbmd0aCA+IENIVU5LX1RIUkVTSE9MRCB8fCBzaXplICsgNCArIGZyYW1lLmJ5dGVMZW5ndGggPiBDSFVOS19USFJFU0hPTEQpIHtcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIHNpemUgKz0gNCArIGZyYW1lLmJ5dGVMZW5ndGg7XHJcbiAgICAgICAgICAgIGVuZCsrO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gQSBuZWFyLXRocmVzaG9sZCBmcmFtZSBmaXRzIGFzIGEgbm9ybWFsIFBPU1QgYm9keSBidXQgbm90IGluc2lkZSBhXHJcbiAgICAgICAgLy8gYmF0Y2ggb25jZSBpdHMgY291bnQgYW5kIGxlbmd0aCBoZWFkZXJzIGFyZSBpbmNsdWRlZC5cclxuICAgICAgICBpZiAoZW5kID09PSBzdGFydCkge1xyXG4gICAgICAgICAgICBhd2FpdCBwb3N0RnJhbWUoY29ubklELCBLSU5EX0RBVEEsIGZyYW1lc1tzdGFydF0sIHVuZGVmaW5lZCwgc2lnbmFsKTtcclxuICAgICAgICAgICAgc3RhcnQrKztcclxuICAgICAgICAgICAgY29udGludWU7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBjb25zdCBncm91cCA9IGZyYW1lcy5zbGljZShzdGFydCwgZW5kKTtcclxuICAgICAgICBpZiAoZ3JvdXAubGVuZ3RoID09PSAxKSB7XHJcbiAgICAgICAgICAgIGF3YWl0IHBvc3RGcmFtZShjb25uSUQsIEtJTkRfREFUQSwgZ3JvdXBbMF0sIHVuZGVmaW5lZCwgc2lnbmFsKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBhd2FpdCBwb3N0QmF0Y2hHcm91cChjb25uSUQsIGdyb3VwLCBzaWduYWwpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBzdGFydCA9IGVuZDtcclxuICAgIH1cclxufVxyXG5cclxuYXN5bmMgZnVuY3Rpb24gcG9zdEJhdGNoR3JvdXAoY29ubklEOiBudW1iZXIsIGZyYW1lczogVWludDhBcnJheVtdLCBzaWduYWw/OiBBYm9ydFNpZ25hbCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgaWYgKGZyYW1lcy5sZW5ndGggPT09IDAgfHwgZnJhbWVzLmxlbmd0aCA+IE1BWF9CQVRDSF9GUkFNRVMpIHtcclxuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJpbnZhbGlkIHN0cmVhbSBiYXRjaCBzaXplXCIpO1xyXG4gICAgfVxyXG5cclxuICAgIGxldCBib2R5ID0gYnVpbGRCYXRjaChmcmFtZXMpO1xyXG4gICAgbGV0IHNlbnQgPSAwO1xyXG4gICAgbGV0IHdhaXQgPSAxO1xyXG4gICAgZm9yICg7Oykge1xyXG4gICAgICAgIGlmIChzaWduYWw/LmFib3J0ZWQpIHRocm93IG5ldyBTdHJlYW1SZXF1ZXN0Q2FuY2VsbGVkKCk7XHJcbiAgICAgICAgY29uc3QgcmVzcCA9IGF3YWl0IGZldGNoKHN0cmVhbVVSTChcInNlbmRcIiksIHtcclxuICAgICAgICAgICAgbWV0aG9kOiBcIlBPU1RcIixcclxuICAgICAgICAgICAgaGVhZGVyczoge1xyXG4gICAgICAgICAgICAgICAgW0hEUl9TRVNTSU9OXTogRy5zZXNzaW9uLFxyXG4gICAgICAgICAgICAgICAgW0hEUl9HRU5FUkFUSU9OXTogU3RyaW5nKEcuZ2VuZXJhdGlvbiksXHJcbiAgICAgICAgICAgICAgICBbSERSX0NPTk5dOiBTdHJpbmcoY29ubklEKSxcclxuICAgICAgICAgICAgICAgIFtIRFJfS0lORF06IFN0cmluZyhLSU5EX0RBVEEpLFxyXG4gICAgICAgICAgICAgICAgW0hEUl9CQVRDSF06IFN0cmluZyhmcmFtZXMubGVuZ3RoIC0gc2VudCksXHJcbiAgICAgICAgICAgICAgICBcIkNvbnRlbnQtVHlwZVwiOiBcImFwcGxpY2F0aW9uL29jdGV0LXN0cmVhbVwiLFxyXG4gICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICBib2R5OiBib2R5IGFzIEJvZHlJbml0LFxyXG4gICAgICAgICAgICBzaWduYWwsXHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgaWYgKHJlc3Aub2spIHJldHVybjtcclxuICAgICAgICBpZiAocmVzcC5zdGF0dXMgIT09IDQyOSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xyXG4gICAgICAgIH1cclxuICAgICAgICAvLyBUaGUgcmVjZWl2ZXIgdG9vayBhIHByZWZpeCBvZiB0aGUgYmF0Y2guIFJlc2VuZCBvbmx5IHdoYXQgaXMgbGVmdCxcclxuICAgICAgICAvLyB3aGljaCBrZWVwcyBvcmRlcmluZyB3aXRob3V0IHJlLWRlbGl2ZXJpbmcgYW55dGhpbmcuXHJcbiAgICAgICAgY29uc3QgYWNjZXB0ZWQgPSBOdW1iZXIocmVzcC5oZWFkZXJzLmdldChIRFJfQkFUQ0gpID8/IDApO1xyXG4gICAgICAgIGNvbnN0IHJlbWFpbmluZyA9IGZyYW1lcy5sZW5ndGggLSBzZW50O1xyXG4gICAgICAgIGlmICghTnVtYmVyLmlzSW50ZWdlcihhY2NlcHRlZCkgfHwgYWNjZXB0ZWQgPCAwIHx8IGFjY2VwdGVkID49IHJlbWFpbmluZykge1xyXG4gICAgICAgICAgICAvLyBaZXJvIGlzIHZhbGlkIChubyBwcm9ncmVzcyksIGJ1dCBhbiBvdXQtb2YtcmFuZ2UgYWNrbm93bGVkZ2VtZW50XHJcbiAgICAgICAgICAgIC8vIGlzIGEgcHJvdG9jb2wgZXJyb3IuIEtlZXAgemVybyBvbiB0aGUgcmV0cnkgcGF0aCBiZWxvdy5cclxuICAgICAgICAgICAgaWYgKGFjY2VwdGVkICE9PSAwKSB0aHJvdyBuZXcgRXJyb3IoXCJpbnZhbGlkIHN0cmVhbSBiYXRjaCBhY2tub3dsZWRnZW1lbnRcIik7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHNlbnQgKz0gYWNjZXB0ZWQ7XHJcbiAgICAgICAgYm9keSA9IGJ1aWxkQmF0Y2goZnJhbWVzLnNsaWNlKHNlbnQpKTtcclxuICAgICAgICBhd2FpdCBzbGVlcCh3YWl0LCBzaWduYWwpO1xyXG4gICAgICAgIHdhaXQgPSBNYXRoLm1pbih3YWl0ICogMiwgNTApO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBidWlsZEJhdGNoKGZyYW1lczogVWludDhBcnJheVtdKTogVWludDhBcnJheSB7XHJcbiAgICBsZXQgc2l6ZSA9IDQ7XHJcbiAgICBmb3IgKGNvbnN0IGYgb2YgZnJhbWVzKSBzaXplICs9IDQgKyBmLmJ5dGVMZW5ndGg7XHJcbiAgICBjb25zdCBvdXQgPSBuZXcgVWludDhBcnJheShzaXplKTtcclxuICAgIGNvbnN0IHZpZXcgPSBuZXcgRGF0YVZpZXcob3V0LmJ1ZmZlcik7XHJcbiAgICB2aWV3LnNldFVpbnQzMigwLCBmcmFtZXMubGVuZ3RoKTtcclxuICAgIGxldCBvZmYgPSA0O1xyXG4gICAgZm9yIChjb25zdCBmIG9mIGZyYW1lcykge1xyXG4gICAgICAgIHZpZXcuc2V0VWludDMyKG9mZiwgZi5ieXRlTGVuZ3RoKTtcclxuICAgICAgICBvZmYgKz0gNDtcclxuICAgICAgICBvdXQuc2V0KGYsIG9mZik7XHJcbiAgICAgICAgb2ZmICs9IGYuYnl0ZUxlbmd0aDtcclxuICAgIH1cclxuICAgIHJldHVybiBvdXQ7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHN0YXJ0UG9sbGluZygpOiB2b2lkIHtcclxuICAgIGlmIChHLnBvbGxpbmcgfHwgIWhhc0RPTSkgcmV0dXJuO1xyXG4gICAgRy5wb2xsaW5nID0gdHJ1ZTtcclxuICAgIHZvaWQgcG9sbExvb3AoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE9uZSBwb2xsIGZvciB0aGUgd2hvbGUgcGFnZSwgc2hhcmVkIGJ5IGV2ZXJ5IGNvbm5lY3Rpb24uXHJcbiAqXHJcbiAqIFRoZXJlIGlzIG5vIHBvbGxpbmcgaW50ZXJ2YWwgYW5kIG5vdGhpbmcgYWRhcHRpdmUgaGVyZSBvbiBwdXJwb3NlLiBUaGUgc2VydmVyXHJcbiAqIGhvbGRzIHRoZSByZXF1ZXN0IHVudGlsIGEgZnJhbWUgZXhpc3RzLCBzbyBkZWxpdmVyeSBsYXRlbmN5IGlzIGFscmVhZHkgfjAgYW5kXHJcbiAqIGEgY2xpZW50LXNpZGUgaW50ZXJ2YWwgY291bGQgb25seSBhZGQgdG8gaXQuIEZyYW1lcyB0aGF0IGFycml2ZSB3aGlsZSBhXHJcbiAqIHJlc3BvbnNlIGlzIGluIGZsaWdodCBzaW1wbHkgYWNjdW11bGF0ZSBhbmQgcmlkZSB0aGUgbmV4dCBvbmUsIHdoaWNoIG1ha2VzXHJcbiAqIHRoZSByb3VuZCB0cmlwIGl0c2VsZiB0aGUgYmF0Y2hpbmcgd2luZG93IOKAlCBpdCB3aWRlbnMgYXMgbG9hZCByaXNlcyB3aXRob3V0XHJcbiAqIGFueXRoaW5nIGhhdmluZyB0byBtZWFzdXJlIGl0LlxyXG4gKi9cclxuYXN5bmMgZnVuY3Rpb24gcG9sbExvb3AoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICBsZXQgYmFja29mZiA9IDA7XHJcbiAgICBjb25zdCBhYm9ydCA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIEcucG9sbEFib3J0ID0gYWJvcnQ7XHJcblxyXG4gICAgd2hpbGUgKEcuY29ubmVjdGlvbnMuc2l6ZSA+IDApIHtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBjb25zdCByZXNwID0gYXdhaXQgZmV0Y2goc3RyZWFtVVJMKFwicG9sbFwiKSwge1xyXG4gICAgICAgICAgICAgICAgbWV0aG9kOiBcIkdFVFwiLFxyXG4gICAgICAgICAgICAgICAgaGVhZGVyczoge1xyXG4gICAgICAgICAgICAgICAgICAgIFtIRFJfU0VTU0lPTl06IEcuc2Vzc2lvbixcclxuICAgICAgICAgICAgICAgICAgICBbSERSX0dFTkVSQVRJT05dOiBTdHJpbmcoRy5nZW5lcmF0aW9uKSxcclxuICAgICAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgICAgICBjYWNoZTogXCJuby1zdG9yZVwiLFxyXG4gICAgICAgICAgICAgICAgc2lnbmFsOiBhYm9ydC5zaWduYWwsXHJcbiAgICAgICAgICAgIH0pO1xyXG5cclxuICAgICAgICAgICAgaWYgKHJlc3Auc3RhdHVzID09PSA0MTApIHtcclxuICAgICAgICAgICAgICAgIC8vIFRoZSBzZXNzaW9uIGlzIGdvbmU6IGFwcCBzaHV0dGluZyBkb3duLCB3aW5kb3cgY2xvc2VkLCBvciB3ZVxyXG4gICAgICAgICAgICAgICAgLy8gd2VyZSByZWFwZWQgYXMgdW5yZWFjaGFibGUuIE5vdGhpbmcgdG8gcmV0cnkuXHJcbiAgICAgICAgICAgICAgICBjbG9zZUFsbCgxMDAxLCBcInNlc3Npb24gY2xvc2VkXCIpO1xyXG4gICAgICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgaWYgKCFyZXNwLm9rICYmIHJlc3Auc3RhdHVzICE9PSAyMDQpIHtcclxuICAgICAgICAgICAgICAgIGlmICghaXNSZXRyeWFibGVQb2xsU3RhdHVzKHJlc3Auc3RhdHVzKSkge1xyXG4gICAgICAgICAgICAgICAgICAgIGZhaWxBbGwoXCJzdHJlYW0gcG9sbCBmYWlsZWQ6IFwiICsgcmVzcC5zdGF0dXMpO1xyXG4gICAgICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwicmV0cnlhYmxlIHN0cmVhbSBwb2xsIGZhaWx1cmU6IFwiICsgcmVzcC5zdGF0dXMpO1xyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICBiYWNrb2ZmID0gMDtcclxuICAgICAgICAgICAgaWYgKHJlc3Auc3RhdHVzICE9PSAyMDQpIHtcclxuICAgICAgICAgICAgICAgIGxldCBidWY6IEFycmF5QnVmZmVyO1xyXG4gICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICBidWYgPSBhd2FpdCByZXNwLmFycmF5QnVmZmVyKCk7XHJcbiAgICAgICAgICAgICAgICB9IGNhdGNoIHtcclxuICAgICAgICAgICAgICAgICAgICAvLyBUaGUgR28gcXVldWUgd2FzIGNvbnN1bWVkIHdoZW4gdGhlIHN1Y2Nlc3NmdWwgcmVzcG9uc2VcclxuICAgICAgICAgICAgICAgICAgICAvLyB3YXMgcHJvZHVjZWQuIFJldHJ5aW5nIGNhbm5vdCByZWNvdmVyIGEgYm9keSB0aGF0IHdhc1xyXG4gICAgICAgICAgICAgICAgICAgIC8vIGludGVycnVwdGVkIGFmdGVyIHRoYXQgcG9pbnQsIHNvIHN1cmZhY2UgdGhlIGxvc3MgYW5kXHJcbiAgICAgICAgICAgICAgICAgICAgLy8gc3RvcCBpbnN0ZWFkIG9mIHNpbGVudGx5IGNvbnRpbnVpbmcgd2l0aCBtaXNzaW5nIGZyYW1lcy5cclxuICAgICAgICAgICAgICAgICAgICBmYWlsQWxsKFwic3RyZWFtIHBvbGwgcmVzcG9uc2UgYm9keSBmYWlsZWRcIik7XHJcbiAgICAgICAgICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICBjb25zdCBmcmFtZXMgPSBkZWNvZGVGcmFtZXMoYnVmKTtcclxuICAgICAgICAgICAgICAgIGlmIChmcmFtZXMgPT09IG51bGwpIHtcclxuICAgICAgICAgICAgICAgICAgICAvLyBBIGZyYW1pbmcgbWlzbWF0Y2ggbWVhbnMgYSBzdGFsZSBjYWNoZWQgcnVudGltZSBhZ2FpbnN0IGFcclxuICAgICAgICAgICAgICAgICAgICAvLyBuZXdlciBHbyBzaWRlLiBSZXRyeWluZyBjYW5ub3QgZml4IHRoYXQsIGFuZCBsZXR0aW5nIGl0XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gZmFsbCBpbnRvIHRoZSBiYWNrb2ZmIGJlbG93IHdvdWxkIHNwaW4gc2lsZW50bHkgZm9yZXZlcixcclxuICAgICAgICAgICAgICAgICAgICAvLyBzbyBmYWlsIGV2ZXJ5IGNvbm5lY3Rpb24gd2l0aCBhIHJlYXNvbiBpbnN0ZWFkLlxyXG4gICAgICAgICAgICAgICAgICAgIGNsb3NlQWxsKDEwMDIsIFwidW5yZWNvZ25pc2VkIHN0cmVhbSBmcmFtaW5nXCIpO1xyXG4gICAgICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgZGVsaXZlcihmcmFtZXMpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIC8vIFJlLXBvbGwgaW1tZWRpYXRlbHksIG9uIGRhdGEgYW5kIG9uIGFuIGV4cGlyZWQgaG9sZCBhbGlrZS5cclxuICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgaWYgKGFib3J0LnNpZ25hbC5hYm9ydGVkIHx8IEcuY29ubmVjdGlvbnMuc2l6ZSA9PT0gMCkgYnJlYWs7XHJcbiAgICAgICAgICAgIGJhY2tvZmYgPSBiYWNrb2ZmID8gTWF0aC5taW4oYmFja29mZiAqIDIsIEJBQ0tPRkZfTUFYKSA6IEJBQ0tPRkZfTUlOO1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgYXdhaXQgc2xlZXAoYmFja29mZiwgYWJvcnQuc2lnbmFsKTtcclxuICAgICAgICAgICAgfSBjYXRjaCB7XHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICBpZiAoRy5wb2xsQWJvcnQgPT09IGFib3J0KSBHLnBvbGxBYm9ydCA9IHVuZGVmaW5lZDtcclxuICAgIEcucG9sbGluZyA9IGZhbHNlO1xyXG5cclxuICAgIC8vIEEgY29ubmVjdGlvbiBvcGVuZWQgd2hpbGUgd2Ugd2VyZSB3aW5kaW5nIGRvd246IHBpY2sgdGhlIGxvb3AgYmFjayB1cC5cclxuICAgIGlmIChHLmNvbm5lY3Rpb25zLnNpemUgPiAwKSB7XHJcbiAgICAgICAgc3RhcnRQb2xsaW5nKCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGlzUmV0cnlhYmxlUG9sbFN0YXR1cyhzdGF0dXM6IG51bWJlcik6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIHN0YXR1cyA9PT0gNDA4IHx8IHN0YXR1cyA9PT0gNDI1IHx8IHN0YXR1cyA9PT0gNDI5IHx8IChzdGF0dXMgPj0gNTAwICYmIHN0YXR1cyA8PSA1OTkpO1xyXG59XHJcblxyXG5pbnRlcmZhY2UgSW5GcmFtZSB7XHJcbiAgICBjb25uSUQ6IG51bWJlcjtcclxuICAgIGtpbmQ6IG51bWJlcjtcclxuICAgIHBheWxvYWQ6IEFycmF5QnVmZmVyO1xyXG59XHJcblxyXG4vLyBWYWxpZGF0ZSB0aGUgY29tcGxldGUgcmVzcG9uc2UgYmVmb3JlIGRpc3BhdGNoaW5nIGFueSBwYXJ0IG9mIGl0LiBBIHNob3J0XHJcbi8vIHBsYXRmb3JtIHJlc3BvbnNlIG11c3Qgbm90IGRlbGl2ZXIgYSB2YWxpZCBwcmVmaXggYW5kIHRoZW4gdGhyb3cgUmFuZ2VFcnJvcjpcclxuLy8gdGhlIEdvIHF1ZXVlIGhhcyBhbHJlYWR5IGJlZW4gZHJhaW5lZCwgc28gcmV0cnlpbmcgY2Fubm90IHJlY292ZXIgdGhvc2VcclxuLy8gZnJhbWVzLiBUcmVhdCB0aGUgd2hvbGUgcmVzcG9uc2UgYXMgYSBwcm90b2NvbCBlcnJvciBpbnN0ZWFkLlxyXG5mdW5jdGlvbiBkZWNvZGVGcmFtZXMoYnVmOiBBcnJheUJ1ZmZlcik6IEluRnJhbWVbXSB8IG51bGwge1xyXG4gICAgaWYgKGJ1Zi5ieXRlTGVuZ3RoIDwgOSkgcmV0dXJuIG51bGw7XHJcbiAgICBjb25zdCBkdiA9IG5ldyBEYXRhVmlldyhidWYpO1xyXG4gICAgaWYgKGR2LmdldFVpbnQ4KDApICE9PSAweDU3IHx8IGR2LmdldFVpbnQ4KDEpICE9PSAweDUzIHx8XHJcbiAgICAgICAgZHYuZ2V0VWludDgoMikgIT09IDB4MzEgfHwgZHYuZ2V0VWludDgoMykgIT09IDB4MDAgfHxcclxuICAgICAgICAoZHYuZ2V0VWludDgoNCkgJiB+MSkgIT09IDApIHtcclxuICAgICAgICByZXR1cm4gbnVsbDtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCBjb3VudCA9IGR2LmdldFVpbnQzMig1KTtcclxuICAgIGNvbnN0IGZyYW1lczogSW5GcmFtZVtdID0gW107XHJcbiAgICBsZXQgb2ZmID0gOTtcclxuXHJcbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGNvdW50OyBpKyspIHtcclxuICAgICAgICBpZiAob2ZmICsgOSA+IGJ1Zi5ieXRlTGVuZ3RoKSByZXR1cm4gbnVsbDtcclxuICAgICAgICBjb25zdCBjb25uSUQgPSBkdi5nZXRVaW50MzIob2ZmKTsgb2ZmICs9IDQ7XHJcbiAgICAgICAgY29uc3Qga2luZCA9IGR2LmdldFVpbnQ4KG9mZik7IG9mZiArPSAxO1xyXG4gICAgICAgIGNvbnN0IGxlbiA9IGR2LmdldFVpbnQzMihvZmYpOyBvZmYgKz0gNDtcclxuICAgICAgICBpZiAoa2luZCA+IEtJTkRfRVJST1IgfHwgbGVuID4gYnVmLmJ5dGVMZW5ndGggLSBvZmYpIHJldHVybiBudWxsO1xyXG4gICAgICAgIGNvbnN0IHBheWxvYWQgPSBidWYuc2xpY2Uob2ZmLCBvZmYgKyBsZW4pOyBvZmYgKz0gbGVuO1xyXG4gICAgICAgIGZyYW1lcy5wdXNoKHsgY29ubklELCBraW5kLCBwYXlsb2FkIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBvZmYgPT09IGJ1Zi5ieXRlTGVuZ3RoID8gZnJhbWVzIDogbnVsbDtcclxufVxyXG5cclxuZnVuY3Rpb24gZGVsaXZlcihmcmFtZXM6IEluRnJhbWVbXSk6IHZvaWQge1xyXG4gICAgZm9yIChjb25zdCBmcmFtZSBvZiBmcmFtZXMpIHtcclxuICAgICAgICBjb25zdCBjb25uSUQgPSBmcmFtZS5jb25uSUQ7XHJcbiAgICAgICAgY29uc3Qga2luZCA9IGZyYW1lLmtpbmQ7XHJcbiAgICAgICAgY29uc3QgcGF5bG9hZCA9IGZyYW1lLnBheWxvYWQ7XHJcbiAgICAgICAgY29uc3QgY29ubiA9IEcuY29ubmVjdGlvbnMuZ2V0KGNvbm5JRCk7XHJcbiAgICAgICAgaWYgKCFjb25uKSBjb250aW51ZTtcclxuXHJcbiAgICAgICAgc3dpdGNoIChraW5kKSB7XHJcbiAgICAgICAgICAgIGNhc2UgS0lORF9EQVRBOlxyXG4gICAgICAgICAgICAgICAgY29ubi5fbWVzc2FnZShwYXlsb2FkKTtcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICBjYXNlIEtJTkRfT1BFTjpcclxuICAgICAgICAgICAgICAgIGNvbm4uX29wZW5lZCgpO1xyXG4gICAgICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgICAgIGNhc2UgS0lORF9DTE9TRTpcclxuICAgICAgICAgICAgICAgIGNvbm4uX2Nsb3NlZCgxMDAwLCBcIlwiLCB0cnVlKTtcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICBjYXNlIEtJTkRfRVJST1I6XHJcbiAgICAgICAgICAgICAgICBjb25uLl9mYWlsKG5ldyBFcnJvcihuZXcgVGV4dERlY29kZXIoKS5kZWNvZGUocGF5bG9hZCkpKTtcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gY2xvc2VBbGwoY29kZTogbnVtYmVyLCByZWFzb246IHN0cmluZyk6IHZvaWQge1xyXG4gICAgZm9yIChjb25zdCBjb25uIG9mIFsuLi5HLmNvbm5lY3Rpb25zLnZhbHVlcygpXSkge1xyXG4gICAgICAgIGNvbm4uX2Nsb3NlZChjb2RlLCByZWFzb24sIGZhbHNlKTtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gZmFpbEFsbChyZWFzb246IHN0cmluZyk6IHZvaWQge1xyXG4gICAgZm9yIChjb25zdCBjb25uIG9mIFsuLi5HLmNvbm5lY3Rpb25zLnZhbHVlcygpXSkge1xyXG4gICAgICAgIGNvbm4uX2ZhaWwobmV3IEVycm9yKHJlYXNvbikpO1xyXG4gICAgfVxyXG59XHJcbiIsIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XHJcblxyXG5pbXBvcnQgeyBFbmFibGUgYXMgRW5hYmxlV01MIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dtbFwiO1xyXG5pbXBvcnQgeyBkZWJ1Z0xvZyB9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy91dGlsc1wiO1xyXG5cclxud2luZG93LndhaWxzID0gUnVudGltZTtcclxuRW5hYmxlV01MKCk7XHJcblxyXG5pZiAoREVCVUcpIHtcclxuICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIilcclxufVxyXG4iLCIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vLyBTZXR1cFxyXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xyXG5cclxuaWYgKGhhc0RPTSkge1xyXG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbn1cclxuXHJcbmltcG9ydCBcIi4vY29udGV4dG1lbnUuanNcIjtcclxuaW1wb3J0IFwiLi9kcmFnLmpzXCI7XHJcbmltcG9ydCBcIi4vYXBwcmVnaW9uLmpzXCI7XHJcblxyXG4vLyBSZS1leHBvcnQgcHVibGljIEFQSVxyXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi9hcHBsaWNhdGlvbi5qc1wiO1xyXG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcclxuaW1wb3J0ICogYXMgQ2FsbCBmcm9tIFwiLi9jYWxscy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkLmpzXCI7XHJcbmltcG9ydCAqIGFzIENyZWF0ZSBmcm9tIFwiLi9jcmVhdGUuanNcIjtcclxuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XHJcbmltcG9ydCAqIGFzIEV2ZW50cyBmcm9tIFwiLi9ldmVudHMuanNcIjtcclxuaW1wb3J0ICogYXMgRmxhZ3MgZnJvbSBcIi4vZmxhZ3MuanNcIjtcclxuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XHJcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcclxuaW1wb3J0ICogYXMgSU9TIGZyb20gXCIuL2lvcy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBBbmRyb2lkIGZyb20gXCIuL2FuZHJvaWQuanNcIjtcclxuaW1wb3J0ICogYXMgVXBkYXRlciBmcm9tIFwiLi91cGRhdGVyLmpzXCI7XHJcbmltcG9ydCBXaW5kb3csIHsgaGFuZGxlRHJhZ0VudGVyLCBoYW5kbGVEcmFnTGVhdmUsIGhhbmRsZURyYWdPdmVyIH0gZnJvbSBcIi4vd2luZG93LmpzXCI7XHJcbmltcG9ydCAqIGFzIFdNTCBmcm9tIFwiLi93bWwuanNcIjtcclxuaW1wb3J0IHsgU3RyZWFtLCBKU09OU3RyZWFtLCBXYWlsc1NvY2tldCwgdHlwZSBKU09OU29ja2V0IH0gZnJvbSBcIi4vc3RyZWFtLmpzXCI7XHJcblxyXG5leHBvcnQge1xyXG4gICAgQXBwbGljYXRpb24sXHJcbiAgICBCcm93c2VyLFxyXG4gICAgQ2FsbCxcclxuICAgIENsaXBib2FyZCxcclxuICAgIERpYWxvZ3MsXHJcbiAgICBFdmVudHMsXHJcbiAgICBGbGFncyxcclxuICAgIFNjcmVlbnMsXHJcbiAgICBTeXN0ZW0sXHJcbiAgICBJT1MsXHJcbiAgICBBbmRyb2lkLFxyXG4gICAgVXBkYXRlcixcclxuICAgIFdpbmRvdyxcclxuICAgIFdNTCxcclxuICAgIFN0cmVhbSxcclxuICAgIEpTT05TdHJlYW0sXHJcbiAgICBXYWlsc1NvY2tldCxcclxuICAgIHR5cGUgSlNPTlNvY2tldFxyXG59O1xyXG5cclxuLyoqXHJcbiAqIEFuIGludGVybmFsIHV0aWxpdHkgY29uc3VtZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxyXG4gKlxyXG4gKiBAaWdub3JlXHJcbiAqL1xyXG5leHBvcnQgeyBDcmVhdGUgfTtcclxuXHJcbmV4cG9ydCAqIGZyb20gXCIuL2NhbmNlbGxhYmxlLmpzXCI7XHJcblxyXG4vLyBFeHBvcnQgdHJhbnNwb3J0IGludGVyZmFjZXMgYW5kIHV0aWxpdGllc1xyXG5leHBvcnQge1xyXG4gICAgc2V0VHJhbnNwb3J0LFxyXG4gICAgZ2V0VHJhbnNwb3J0LFxyXG4gICAgdHlwZSBSdW50aW1lVHJhbnNwb3J0LFxyXG4gICAgb2JqZWN0TmFtZXMsXHJcbiAgICBjbGllbnRJZCxcclxufSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5pbXBvcnQgeyBjbGllbnRJZCB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbi8vIE5vdGlmeSBiYWNrZW5kXHJcbmlmIChoYXNET00pIHtcclxuICAgIHdpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcclxuICAgIHdpbmRvdy5fd2FpbHMuY2xpZW50SWQgPSBjbGllbnRJZDtcclxufVxyXG5cclxuLy8gUmVnaXN0ZXIgcGxhdGZvcm0gaGFuZGxlcnMgKGludGVybmFsIEFQSSlcclxuLy8gTm90ZTogV2luZG93IGlzIHRoZSB0aGlzV2luZG93IGluc3RhbmNlIChkZWZhdWx0IGV4cG9ydCBmcm9tIHdpbmRvdy50cylcclxuLy8gQmluZGluZyBlbnN1cmVzICd0aGlzJyBjb3JyZWN0bHkgcmVmZXJzIHRvIHRoZSBjdXJyZW50IHdpbmRvdyBpbnN0YW5jZVxyXG5pZiAoaGFzRE9NKSB7XHJcbiAgICB3aW5kb3cuX3dhaWxzLmhhbmRsZVBsYXRmb3JtRmlsZURyb3AgPSBXaW5kb3cuSGFuZGxlUGxhdGZvcm1GaWxlRHJvcC5iaW5kKFdpbmRvdyk7XHJcbn1cclxuXHJcbi8vIExpbnV4LXNwZWNpZmljIGRyYWcgaGFuZGxlcnMgKEdUSyBpbnRlcmNlcHRzIERPTSBkcmFnIGV2ZW50cylcclxuaWYgKGhhc0RPTSkge1xyXG4gICAgd2luZG93Ll93YWlscy5oYW5kbGVEcmFnRW50ZXIgPSBoYW5kbGVEcmFnRW50ZXI7XHJcbiAgICB3aW5kb3cuX3dhaWxzLmhhbmRsZURyYWdMZWF2ZSA9IGhhbmRsZURyYWdMZWF2ZTtcclxuICAgIHdpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ092ZXIgPSBoYW5kbGVEcmFnT3ZlcjtcclxufVxyXG5cclxuaWYgKGhhc0RPTSkge1xyXG4gICAgU3lzdGVtLmludm9rZShcIndhaWxzOnJ1bnRpbWU6cmVhZHlcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2FkcyBhIHNjcmlwdCBmcm9tIHRoZSBnaXZlbiBVUkwgaWYgaXQgZXhpc3RzLlxyXG4gKiBVc2VzIEhFQUQgcmVxdWVzdCB0byBjaGVjayBleGlzdGVuY2UsIHRoZW4gaW5qZWN0cyBhIHNjcmlwdCB0YWcuXHJcbiAqIFNpbGVudGx5IGlnbm9yZXMgaWYgdGhlIHNjcmlwdCBkb2Vzbid0IGV4aXN0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIGxvYWRPcHRpb25hbFNjcmlwdCh1cmw6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGZldGNoKHVybCwgeyBtZXRob2Q6ICdIRUFEJyB9KVxyXG4gICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcclxuICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAvLyBWZXJpZnkgdGhlIHJlc3BvbnNlIGlzIGFjdHVhbGx5IEphdmFTY3JpcHQgYW5kIG5vdCBhbiBIVE1MIGZhbGxiYWNrXHJcbiAgICAgICAgICAgICAgICAvLyAoZS5nLiBWaXRlIGRldiBzZXJ2ZXIgcmV0dXJucyBpbmRleC5odG1sIGZvciB1bmtub3duIHJvdXRlcylcclxuICAgICAgICAgICAgICAgIGNvbnN0IGNvbnRlbnRUeXBlID0gKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KCdjb250ZW50LXR5cGUnKSB8fCAnJykudG9Mb3dlckNhc2UoKTtcclxuICAgICAgICAgICAgICAgIGlmIChjb250ZW50VHlwZS5pbmNsdWRlcygnamF2YXNjcmlwdCcpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3Qgc2NyaXB0ID0gZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgnc2NyaXB0Jyk7XHJcbiAgICAgICAgICAgICAgICAgICAgc2NyaXB0LnNyYyA9IHVybDtcclxuICAgICAgICAgICAgICAgICAgICBkb2N1bWVudC5oZWFkLmFwcGVuZENoaWxkKHNjcmlwdCk7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KVxyXG4gICAgICAgIC5jYXRjaCgoKSA9PiB7fSk7IC8vIFNpbGVudGx5IGlnbm9yZSAtIHNjcmlwdCBpcyBvcHRpb25hbFxyXG59XHJcblxyXG4vLyBMb2FkIGN1c3RvbS5qcyBpZiBhdmFpbGFibGUgKHVzZWQgYnkgc2VydmVyIG1vZGUgZm9yIFdlYlNvY2tldCBldmVudHMsIGV0Yy4pXHJcbmlmIChoYXNET00pIHtcclxuICAgIGxvYWRPcHRpb25hbFNjcmlwdCgnL3dhaWxzL2N1c3RvbS5qcycpO1xyXG59XHJcbiJdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztFQWlCYSxTQUFTLE9BQU8sV0FBVyxlQUFlLE9BQU8sYUFBYTs7Ozs7QUNzQjNFLFNBQWdCLE9BQU8sT0FBZSxJQUFZO0NBQzlDLE1BQU0sU0FBUyxPQUFPLFdBQVcsZUFBZSxPQUFPLE9BQU8sb0JBQW9CLGFBQzVFLFNBQ0E7Q0FFTixJQUFJLFFBQVE7RUFDUixNQUFNLFFBQVEsSUFBSSxXQUFXLElBQUk7RUFDakMsT0FBTyxnQkFBZ0IsS0FBSztFQUM1QixJQUFJLEtBQUs7RUFDVCxLQUFLLElBQUks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