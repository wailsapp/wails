(function () {
  "use strict";
  if (window === window.top) return;
  // WebKitGTK does not expose the sender frame to the native callback.
  // Block direct child-frame access to the privileged v2 IPC channel.
  // Same-origin children can still access the trusted parent bridge; they
  // share its origin authority, unlike cross-origin or opaque sandbox frames.
  var blocked = Object.freeze({ postMessage: function () {} });
  var handlers = Object.freeze({ external: blocked });
  var guarded = Object.freeze({ messageHandlers: handlers });
  try {
    Object.defineProperty(window, "webkit", {
      value: guarded, writable: false, configurable: false, enumerable: true
    });
  } catch (e) {
    Object.defineProperty(window.webkit.messageHandlers, "external", {
      value: blocked, writable: false, configurable: false
    });
  }
})();
