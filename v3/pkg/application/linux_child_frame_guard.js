(function () {
  "use strict";
  // WebKitGTK gives the native "external" handler no frame information, so
  // child frames are filtered here, before any page script can run.
  try {
    if (window === window.top) return;
    var wk = window.webkit;
    var handler = wk && wk.messageHandlers && wk.messageHandlers.external;
    if (!handler || typeof handler.postMessage !== "function") return;
    var original = handler.postMessage;
    var post = function postMessage(message) {
      // Coerce once so the string checked is the string sent.
      var text = String(message);
      if (text.indexOf("wails:") === 0) return; // reserved: main frame only
      return original.call(handler, text);
    };
    var guarded = Object.freeze({ postMessage: post });
    var handlers = Object.freeze({ external: guarded });
    var fake = Object.freeze({ messageHandlers: handlers });
    try {
      Object.defineProperty(window, "webkit", { value: fake, writable: false, configurable: false, enumerable: true });
    } catch (e) {
      Object.defineProperty(wk.messageHandlers, "external", { value: guarded, writable: false, configurable: false });
    }
  } catch (e) {}
})();
