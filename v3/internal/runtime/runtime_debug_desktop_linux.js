(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

  // desktop/ipc.js
  var postMessage = null;
  (function() {
    let _deeptest = function(s) {
      let obj = window[s.shift()];
      while (obj && s.length)
        obj = obj[s.shift()];
      return obj;
    };
    let windows = _deeptest(["chrome", "webview", "postMessage"]);
    let mac_linux = _deeptest(["webkit", "messageHandlers", "external", "postMessage"]);
    if (!windows && !mac_linux) {
      console.error("Unsupported Platform");
      return;
    }
    if (windows) {
      postMessage = (message) => window.chrome.webview.postMessage(message);
    }
    if (mac_linux) {
      postMessage = (message) => window.webkit.messageHandlers.external.postMessage(message);
    }
  })();
  function invoke(message, id) {
    if (id && id !== -1) {
      postMessage("WINDOWID:" + id + ":" + message);
    } else {
      postMessage(message);
    }
  }

  // desktop/calls.js
  var callbacks = {};
  function cryptoRandom() {
    let array = new Uint32Array(1);
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
    let windowID = window.wails.window.ID();
    return new Promise(function(resolve, reject) {
      let callbackID;
      do {
        callbackID = name + "-" + randomFunc();
      } while (callbacks[callbackID]);
      let timeoutHandle;
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
          callbackID,
          windowID
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
      let callbackID;
      do {
        callbackID = id + "-" + randomFunc();
      } while (callbacks[callbackID]);
      let timeoutHandle;
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
          callbackID,
          windowID: window.wails.window.ID()
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

  // desktop/events.js
  var eventListeners = {};
  function notifyListeners(eventData) {
    let eventName = eventData.name;
    if (eventListeners[eventName]) {
      const newEventListenerList = eventListeners[eventName].slice();
      for (let count = 0; count < eventListeners[eventName].length; count += 1) {
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
  function removeListener(eventName) {
    delete eventListeners[eventName];
    window.WailsInvoke("EX" + eventName);
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

  // desktop/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    SetText: () => SetText,
    Text: () => Text
  });

  // desktop/runtime.js
  var runtimeURL = window.location.origin + "/wails/runtime";
  function runtimeCall(method, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("method", method);
    if (args) {
      for (let key in args) {
        url.searchParams.append(key, args[key]);
      }
    }
    return new Promise((resolve, reject) => {
      fetch(url).then((response) => {
        if (response.ok) {
          if (response.headers.get("content-type") && response.headers.get("content-type").indexOf("application/json") !== -1) {
            return response.json();
          } else {
            return response.text();
          }
        }
        reject(Error(response.statusText));
      }).then((data) => resolve(data)).catch((error) => reject(error));
    });
  }
  function newRuntimeCaller(object, id) {
    if (!id || id === -1) {
      return function(method, args) {
        return runtimeCall(object + "." + method, args);
      };
    }
    return function(method, args) {
      args = args || {};
      args["windowID"] = id;
      return runtimeCall(object + "." + method, args);
    };
  }

  // desktop/clipboard.js
  var call = newRuntimeCaller("clipboard");
  function SetText(text) {
    return call("SetText", { text });
  }
  function Text() {
    return call("Text");
  }

  // desktop/window.js
  function newWindow(id) {
    let call2 = newRuntimeCaller("window", id);
    return {
      // Reload: () => call('WR'),
      // ReloadApp: () => call('WR'),
      // SetSystemDefaultTheme: () => call('WASDT'),
      // SetLightTheme: () => call('WALT'),
      // SetDarkTheme: () => call('WADT'),
      Center: () => call2("Center"),
      SetTitle: (title) => call2("SetTitle", { title }),
      Fullscreen: () => call2("Fullscreen"),
      UnFullscreen: () => call2("UnFullscreen"),
      SetSize: (width, height) => call2("SetSize", { width, height }),
      Size: () => {
        return call2("Size");
      },
      SetMaxSize: (width, height) => call2("SetMaxSize", { width, height }),
      SetMinSize: (width, height) => call2("SetMinSize", { width, height }),
      SetAlwaysOnTop: (b) => call2("SetAlwaysOnTop", { alwaysOnTop: b }),
      SetPosition: (x, y) => call2("SetPosition", { x, y }),
      Position: () => {
        return call2("Position");
      },
      Screen: () => {
        return call2("Screen");
      },
      Hide: () => call2("Hide"),
      Maximise: () => call2("Maximise"),
      Show: () => call2("Show"),
      ToggleMaximise: () => call2("ToggleMaximise"),
      UnMaximise: () => call2("UnMaximise"),
      Minimise: () => call2("Minimise"),
      UnMinimise: () => call2("UnMinimise"),
      SetBackgroundColour: (r, g, b, a) => call2("SetBackgroundColour", { R, G, B, A })
    };
  }

  // desktop/main.js
  window.wails = {
    Callback,
    callbacks,
    EventsNotify,
    eventListeners,
    SetBindings
  };
  function newRuntime(id) {
    return {
      // Log: newLog(id),
      // Browser: newBrowser(id),
      // Screen: newScreen(id),
      // Events: newEvents(id),
      Clipboard: {
        ...clipboard_exports
      },
      Window: newWindow(id),
      Show: () => invoke("S"),
      Hide: () => invoke("H"),
      Quit: () => invoke("Q")
      // GetWindow: function (windowID) {
      //     if (!windowID) {
      //         return this.Window;
      //     }
      //     return newWindow(windowID);
      // }
    };
  }
  window.runtime = newRuntime(-1);
  if (true) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
  }
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9pcGMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2V2ZW50cy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3AvY2xpcGJvYXJkLmpzIiwgImRlc2t0b3AvcnVudGltZS5qcyIsICJkZXNrdG9wL3dpbmRvdy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cblxubGV0IHBvc3RNZXNzYWdlID0gbnVsbDtcblxuKGZ1bmN0aW9uICgpIHtcblx0Ly8gQ3JlZGl0OiBodHRwczovL3N0YWNrb3ZlcmZsb3cuY29tL2EvMjYzMTUyMVxuXHRsZXQgX2RlZXB0ZXN0ID0gZnVuY3Rpb24gKHMpIHtcblx0XHRsZXQgb2JqID0gd2luZG93W3Muc2hpZnQoKV07XG5cdFx0d2hpbGUgKG9iaiAmJiBzLmxlbmd0aCkgb2JqID0gb2JqW3Muc2hpZnQoKV07XG5cdFx0cmV0dXJuIG9iajtcblx0fTtcblx0bGV0IHdpbmRvd3MgPSBfZGVlcHRlc3QoW1wiY2hyb21lXCIsIFwid2Vidmlld1wiLCBcInBvc3RNZXNzYWdlXCJdKTtcblx0bGV0IG1hY19saW51eCA9IF9kZWVwdGVzdChbXCJ3ZWJraXRcIiwgXCJtZXNzYWdlSGFuZGxlcnNcIiwgXCJleHRlcm5hbFwiLCBcInBvc3RNZXNzYWdlXCJdKTtcblxuXHRpZiAoIXdpbmRvd3MgJiYgIW1hY19saW51eCkge1xuXHRcdGNvbnNvbGUuZXJyb3IoXCJVbnN1cHBvcnRlZCBQbGF0Zm9ybVwiKTtcblx0XHRyZXR1cm47XG5cdH1cblxuXHRpZiAod2luZG93cykge1xuXHRcdHBvc3RNZXNzYWdlID0gKG1lc3NhZ2UpID0+IHdpbmRvdy5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZShtZXNzYWdlKTtcblx0fVxuXHRpZiAobWFjX2xpbnV4KSB7XG5cdFx0cG9zdE1lc3NhZ2UgPSAobWVzc2FnZSkgPT4gd2luZG93LndlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2UobWVzc2FnZSk7XG5cdH1cbn0pKCk7XG5cbmV4cG9ydCBmdW5jdGlvbiBpbnZva2UobWVzc2FnZSwgaWQpIHtcblx0aWYoIGlkICYmIGlkICE9PSAtMSkge1xuXHRcdHBvc3RNZXNzYWdlKFwiV0lORE9XSUQ6XCIrIGlkICsgXCI6XCIgKyBtZXNzYWdlKTtcblx0fSBlbHNlIHtcblx0XHRwb3N0TWVzc2FnZShtZXNzYWdlKTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gY3J5cHRvUmFuZG9tKCkge1xuXHRsZXQgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciB1c2luZyBkYSBvbGQtc2tvb2wgTWF0aC5SYW5kb21cbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XG5cdHJldHVybiBNYXRoLnJhbmRvbSgpICogOTAwNzE5OTI1NDc0MDk5MTtcbn1cblxuLy8gUGljayBhIHJhbmRvbSBudW1iZXIgZnVuY3Rpb24gYmFzZWQgb24gYnJvd3NlciBjYXBhYmlsaXR5XG5sZXQgcmFuZG9tRnVuYztcbmlmICh3aW5kb3cuY3J5cHRvKSB7XG5cdHJhbmRvbUZ1bmMgPSBjcnlwdG9SYW5kb207XG59IGVsc2Uge1xuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XG59XG5cblxuLyoqXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXG4gKiBiYWNrZW5kIHJlc3BvbmRzLiBUaGlzIHdpbGwgYmUgcmVzb2x2ZWQgd2hlbiB0aGUgY2FsbCB3YXMgc3VjY2Vzc2Z1bFxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXG4gKiB0aW1lIChpbiBtaWxsaXNlY29uZHMpIHRoZW4gdGhlIHByb21pc2UgaXMgcmVqZWN0ZWQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWVcbiAqIEBwYXJhbSB7YW55PX0gYXJnc1xuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XG4gKiBAcmV0dXJuc1xuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XG5cblx0Ly8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcblx0XHR0aW1lb3V0ID0gMDtcblx0fVxuXG5cdGxldCB3aW5kb3dJRCA9IHdpbmRvdy53YWlscy53aW5kb3cuSUQoKTtcblxuXHQvLyBDcmVhdGUgYSBwcm9taXNlXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cblx0XHQvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuXHRcdGxldCBjYWxsYmFja0lEO1xuXHRcdGRvIHtcblx0XHRcdGNhbGxiYWNrSUQgPSBuYW1lICsgJy0nICsgcmFuZG9tRnVuYygpO1xuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cblx0XHRsZXQgdGltZW91dEhhbmRsZTtcblx0XHQvLyBTZXQgdGltZW91dFxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xuXHRcdFx0dGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcblx0XHRcdH0sIHRpbWVvdXQpO1xuXHRcdH1cblxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXG5cdFx0Y2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcblx0XHRcdHJlamVjdDogcmVqZWN0LFxuXHRcdFx0cmVzb2x2ZTogcmVzb2x2ZVxuXHRcdH07XG5cblx0XHR0cnkge1xuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0bmFtZSxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdFx0d2luZG93SUQsXG5cdFx0XHR9O1xuXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XG4gICAgICAgIH1cbiAgICB9KTtcbn1cblxud2luZG93Lk9iZnVzY2F0ZWRDYWxsID0gKGlkLCBhcmdzLCB0aW1lb3V0KSA9PiB7XG5cbiAgICAvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcbiAgICBpZiAodGltZW91dCA9PSBudWxsKSB7XG4gICAgICAgIHRpbWVvdXQgPSAwO1xuICAgIH1cblxuICAgIC8vIENyZWF0ZSBhIHByb21pc2VcbiAgICByZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xuXG4gICAgICAgIC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXG4gICAgICAgIGxldCBjYWxsYmFja0lEO1xuICAgICAgICBkbyB7XG4gICAgICAgICAgICBjYWxsYmFja0lEID0gaWQgKyAnLScgKyByYW5kb21GdW5jKCk7XG4gICAgICAgIH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cbiAgICAgICAgbGV0IHRpbWVvdXRIYW5kbGU7XG4gICAgICAgIC8vIFNldCB0aW1lb3V0XG4gICAgICAgIGlmICh0aW1lb3V0ID4gMCkge1xuICAgICAgICAgICAgdGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcignQ2FsbCB0byBtZXRob2QgJyArIGlkICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XG4gICAgICAgICAgICB9LCB0aW1lb3V0KTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFN0b3JlIGNhbGxiYWNrXG4gICAgICAgIGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXG4gICAgICAgICAgICByZWplY3Q6IHJlamVjdCxcbiAgICAgICAgICAgIHJlc29sdmU6IHJlc29sdmVcbiAgICAgICAgfTtcblxuICAgICAgICB0cnkge1xuICAgICAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0aWQsXG5cdFx0XHRcdGFyZ3MsXG5cdFx0XHRcdGNhbGxiYWNrSUQsXG5cdFx0XHRcdHdpbmRvd0lEOiB3aW5kb3cud2FpbHMud2luZG93LklEKCksXG5cdFx0XHR9O1xuXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ2MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XG4gICAgICAgIH1cbiAgICB9KTtcbn07XG5cblxuLyoqXG4gKiBDYWxsZWQgYnkgdGhlIGJhY2tlbmQgdG8gcmV0dXJuIGRhdGEgdG8gYSBwcmV2aW91c2x5IGNhbGxlZFxuICogYmluZGluZyBpbnZvY2F0aW9uXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGluY29taW5nTWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbGJhY2soaW5jb21pbmdNZXNzYWdlKSB7XG5cdC8vIFBhcnNlIHRoZSBtZXNzYWdlXG5cdGxldCBtZXNzYWdlO1xuXHR0cnkge1xuXHRcdG1lc3NhZ2UgPSBKU09OLnBhcnNlKGluY29taW5nTWVzc2FnZSk7XG5cdH0gY2F0Y2ggKGUpIHtcblx0XHRjb25zdCBlcnJvciA9IGBJbnZhbGlkIEpTT04gcGFzc2VkIHRvIGNhbGxiYWNrOiAke2UubWVzc2FnZX0uIE1lc3NhZ2U6ICR7aW5jb21pbmdNZXNzYWdlfWA7XG5cdFx0cnVudGltZS5Mb2dEZWJ1ZyhlcnJvcik7XG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcblx0fVxuXHRsZXQgY2FsbGJhY2tJRCA9IG1lc3NhZ2UuY2FsbGJhY2tpZDtcblx0bGV0IGNhbGxiYWNrRGF0YSA9IGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcblx0aWYgKCFjYWxsYmFja0RhdGEpIHtcblx0XHRjb25zdCBlcnJvciA9IGBDYWxsYmFjayAnJHtjYWxsYmFja0lEfScgbm90IHJlZ2lzdGVyZWQhISFgO1xuXHRcdGNvbnNvbGUuZXJyb3IoZXJyb3IpOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcblx0fVxuXHRjbGVhclRpbWVvdXQoY2FsbGJhY2tEYXRhLnRpbWVvdXRIYW5kbGUpO1xuXG5cdGRlbGV0ZSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XG5cblx0aWYgKG1lc3NhZ2UuZXJyb3IpIHtcblx0XHRjYWxsYmFja0RhdGEucmVqZWN0KG1lc3NhZ2UuZXJyb3IpO1xuXHR9IGVsc2Uge1xuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XG5cbiAgICAvLyBHZXQgdGhlIGV2ZW50IG5hbWVcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XG5cbiAgICAvLyBDaGVjayBpZiB3ZSBoYXZlIGFueSBsaXN0ZW5lcnMgZm9yIHRoaXMgZXZlbnRcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSkge1xuXG4gICAgICAgIC8vIEtlZXAgYSBsaXN0IG9mIGxpc3RlbmVyIGluZGV4ZXMgdG8gZGVzdHJveVxuICAgICAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uc2xpY2UoKTtcblxuICAgICAgICAvLyBJdGVyYXRlIGxpc3RlbmVyc1xuICAgICAgICBmb3IgKGxldCBjb3VudCA9IDA7IGNvdW50IDwgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGg7IGNvdW50ICs9IDEpIHtcblxuICAgICAgICAgICAgLy8gR2V0IG5leHQgbGlzdGVuZXJcbiAgICAgICAgICAgIGNvbnN0IGxpc3RlbmVyID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXVtjb3VudF07XG5cbiAgICAgICAgICAgIGxldCBkYXRhID0gZXZlbnREYXRhLmRhdGE7XG5cbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xuICAgICAgICAgICAgY29uc3QgZGVzdHJveSA9IGxpc3RlbmVyLkNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgaWYgKGRlc3Ryb3kpIHtcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxuICAgICAgICAgICAgICAgIG5ld0V2ZW50TGlzdGVuZXJMaXN0LnNwbGljZShjb3VudCwgMSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cblxuICAgICAgICAvLyBVcGRhdGUgY2FsbGJhY2tzIHdpdGggbmV3IGxpc3Qgb2YgbGlzdGVuZXJzXG4gICAgICAgIGlmIChuZXdFdmVudExpc3RlbmVyTGlzdC5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gbmV3RXZlbnRMaXN0ZW5lckxpc3Q7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogTm90aWZ5IGluZm9ybXMgZnJvbnRlbmQgbGlzdGVuZXJzIHRoYXQgYW4gZXZlbnQgd2FzIGVtaXR0ZWQgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBub3RpZnlNZXNzYWdlIC0gZW5jb2RlZCBub3RpZmljYXRpb24gbWVzc2FnZVxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNOb3RpZnkobm90aWZ5TWVzc2FnZSkge1xuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXG4gICAgbGV0IG1lc3NhZ2U7XG4gICAgdHJ5IHtcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zdCBlcnJvciA9ICdJbnZhbGlkIEpTT04gcGFzc2VkIHRvIE5vdGlmeTogJyArIG5vdGlmeU1lc3NhZ2U7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XG4gICAgfVxuICAgIG5vdGlmeUxpc3RlbmVycyhtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNFbWl0KGV2ZW50TmFtZSkge1xuXG4gICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgbmFtZTogZXZlbnROYW1lLFxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxuICAgIH07XG5cbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXG4gICAgbm90aWZ5TGlzdGVuZXJzKHBheWxvYWQpO1xuXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRUUnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xufVxuXG5mdW5jdGlvbiByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpIHtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJzXG4gICAgZGVsZXRlIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV07XG5cbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFWCcgKyBldmVudE5hbWUpO1xufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcmVzIGNhbiBiZSB1bnJlZ2lzdGVyZWQgdmlhIGBhZGRpdGlvbmFsRXZlbnROYW1lc2BcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcblxuICAgIGlmIChhZGRpdGlvbmFsRXZlbnROYW1lcy5sZW5ndGggPiAwKSB7XG4gICAgICAgIGFkZGl0aW9uYWxFdmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IHtcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcbiAgICAgICAgfSlcbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGFsbCBldmVudCBsaXN0ZW5lcnMgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cbiAqL1xuIGV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmZBbGwoKSB7XG4gICAgY29uc3QgZXZlbnROYW1lcyA9IE9iamVjdC5rZXlzKGV2ZW50TGlzdGVuZXJzKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSAhPT0gZXZlbnROYW1lcy5sZW5ndGg7IGkrKykge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWVzW2ldKTtcbiAgICB9XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBFdmVudHNPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbiBmdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcbiAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG5cbiAgICAvLyBDbGVhbiB1cCBpZiB0aGVyZSBhcmUgbm8gZXZlbnQgbGlzdGVuZXJzIGxlZnRcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGggPT09IDApIHtcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcbiAgICB9XG59XG4iLCAiLypcbiBfICAgICAgIF9fICAgICAgXyBfXyAgICBcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKSBcbnxfXy98X18vXFxfXyxfL18vXy9fX19fLyAgXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5pbXBvcnQge0NhbGx9IGZyb20gJy4vY2FsbHMnO1xuXG4vLyBUaGlzIGlzIHdoZXJlIHdlIGJpbmQgZ28gbWV0aG9kIHdyYXBwZXJzXG53aW5kb3cuZ28gPSB7fTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldEJpbmRpbmdzKGJpbmRpbmdzTWFwKSB7XG5cdHRyeSB7XG5cdFx0YmluZGluZ3NNYXAgPSBKU09OLnBhcnNlKGJpbmRpbmdzTWFwKTtcblx0fSBjYXRjaCAoZSkge1xuXHRcdGNvbnNvbGUuZXJyb3IoZSk7XG5cdH1cblxuXHQvLyBJbml0aWFsaXNlIHRoZSBiaW5kaW5ncyBtYXBcblx0d2luZG93LmdvID0gd2luZG93LmdvIHx8IHt9O1xuXG5cdC8vIEl0ZXJhdGUgcGFja2FnZSBuYW1lc1xuXHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcCkuZm9yRWFjaCgocGFja2FnZU5hbWUpID0+IHtcblxuXHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxuXHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdIHx8IHt9O1xuXG5cdFx0Ly8gSXRlcmF0ZSBzdHJ1Y3QgbmFtZXNcblx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV0pLmZvckVhY2goKHN0cnVjdE5hbWUpID0+IHtcblxuXHRcdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XG5cdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSB8fCB7fTtcblxuXHRcdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdKS5mb3JFYWNoKChtZXRob2ROYW1lKSA9PiB7XG5cblx0XHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXVttZXRob2ROYW1lXSA9IGZ1bmN0aW9uICgpIHtcblxuXHRcdFx0XHRcdC8vIE5vIHRpbWVvdXQgYnkgZGVmYXVsdFxuXHRcdFx0XHRcdGxldCB0aW1lb3V0ID0gMDtcblxuXHRcdFx0XHRcdC8vIEFjdHVhbCBmdW5jdGlvblxuXHRcdFx0XHRcdGZ1bmN0aW9uIGR5bmFtaWMoKSB7XG5cdFx0XHRcdFx0XHRjb25zdCBhcmdzID0gW10uc2xpY2UuY2FsbChhcmd1bWVudHMpO1xuXHRcdFx0XHRcdFx0cmV0dXJuIENhbGwoW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXS5qb2luKCcuJyksIGFyZ3MsIHRpbWVvdXQpO1xuXHRcdFx0XHRcdH1cblxuXHRcdFx0XHRcdC8vIEFsbG93IHNldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxuXHRcdFx0XHRcdGR5bmFtaWMuc2V0VGltZW91dCA9IGZ1bmN0aW9uIChuZXdUaW1lb3V0KSB7XG5cdFx0XHRcdFx0XHR0aW1lb3V0ID0gbmV3VGltZW91dDtcblx0XHRcdFx0XHR9O1xuXG5cdFx0XHRcdFx0Ly8gQWxsb3cgZ2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZHluYW1pYy5nZXRUaW1lb3V0ID0gZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRcdFx0cmV0dXJuIHRpbWVvdXQ7XG5cdFx0XHRcdFx0fTtcblxuXHRcdFx0XHRcdHJldHVybiBkeW5hbWljO1xuXHRcdFx0XHR9KCk7XG5cdFx0XHR9KTtcblx0XHR9KTtcblx0fSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChcIlRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIGZvciAobGV0IGtleSBpbiBhcmdzKSB7XG4gICAgICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChrZXksIGFyZ3Nba2V5XSk7XG4gICAgICAgIH1cbiAgICB9XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZmV0Y2godXJsKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiY29udGVudC10eXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiY29udGVudC10eXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgaWQpIHtcbiAgICBpZiAoIWlkIHx8IGlkID09PSAtMSkge1xuICAgICAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncykge1xuICAgICAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICAgICAgfTtcbiAgICB9XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgYXJncyA9IGFyZ3MgfHwge307XG4gICAgICAgIGFyZ3NbXCJ3aW5kb3dJRFwiXSA9IGlkO1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGwob2JqZWN0ICsgXCIuXCIgKyBtZXRob2QsIGFyZ3MpO1xuICAgIH1cbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdXaW5kb3coaWQpIHtcbiAgICBsZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJ3aW5kb3dcIiwgaWQpO1xuICAgIHJldHVybiB7XG4gICAgICAgIC8vIFJlbG9hZDogKCkgPT4gY2FsbCgnV1InKSxcbiAgICAgICAgLy8gUmVsb2FkQXBwOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBTZXRTeXN0ZW1EZWZhdWx0VGhlbWU6ICgpID0+IGNhbGwoJ1dBU0RUJyksXG4gICAgICAgIC8vIFNldExpZ2h0VGhlbWU6ICgpID0+IGNhbGwoJ1dBTFQnKSxcbiAgICAgICAgLy8gU2V0RGFya1RoZW1lOiAoKSA9PiBjYWxsKCdXQURUJyksXG4gICAgICAgIENlbnRlcjogKCkgPT4gY2FsbCgnQ2VudGVyJyksXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IGNhbGwoJ1NldFRpdGxlJywge3RpdGxlfSksXG4gICAgICAgIEZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ0Z1bGxzY3JlZW4nKSxcbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdVbkZ1bGxzY3JlZW4nKSxcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoJ1NldFNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG4gICAgICAgIFNpemU6ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1NpemUnKSB9LFxuICAgICAgICBTZXRNYXhTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWF4U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0TWluU2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoJ1NldE1pblNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAoYikgPT4gY2FsbCgnU2V0QWx3YXlzT25Ub3AnLCB7YWx3YXlzT25Ub3A6Yn0pLFxuICAgICAgICBTZXRQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoJ1NldFBvc2l0aW9uJywge3gseX0pLFxuICAgICAgICBQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUG9zaXRpb24nKSB9LFxuICAgICAgICBTY3JlZW46ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1NjcmVlbicpIH0sXG4gICAgICAgIEhpZGU6ICgpID0+IGNhbGwoJ0hpZGUnKSxcbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IGNhbGwoJ01heGltaXNlJyksXG4gICAgICAgIFNob3c6ICgpID0+IGNhbGwoJ1Nob3cnKSxcbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1RvZ2dsZU1heGltaXNlJyksXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IGNhbGwoJ1VuTWF4aW1pc2UnKSxcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IGNhbGwoJ01pbmltaXNlJyksXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwoJ1VuTWluaW1pc2UnKSxcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IGNhbGwoJ1NldEJhY2tncm91bmRDb2xvdXInLCB7UiwgRywgQiwgQX0pLFxuICAgIH1cbn1cblxuLy8gZXhwb3J0IGZ1bmN0aW9uIElzRnVsbHNjcmVlbjogKCk9PiAvLyAgICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNGdWxsc2NyZWVuXCIpLFxuLy9cblxuLy8gZXhwb3J0IGZ1bmN0aW9uIElzTWF4aW1pc2VkOiAoKT0+IC8vICAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc01heGltaXNlZFwiKSxcbi8vXG5cbi8vIGV4cG9ydCBmdW5jdGlvbiBJc01pbmltaXNlZDogKCk9PiAvLyAgICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNaW5pbWlzZWRcIiksXG4vL1xuXG4vLyBleHBvcnQgZnVuY3Rpb24gSXNOb3JtYWw6ICgpPT4gLy8gICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTm9ybWFsXCIpLFxuLy9cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge2ludm9rZX0gZnJvbSBcIi4vaXBjLmpzXCI7XG5pbXBvcnQge0NhbGxiYWNrLCBjYWxsYmFja3N9IGZyb20gJy4vY2FsbHMnO1xuaW1wb3J0IHtFdmVudHNOb3RpZnksIGV2ZW50TGlzdGVuZXJzfSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7U2V0QmluZGluZ3N9IGZyb20gXCIuL2JpbmRpbmdzXCI7XG5cblxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gJy4vY2xpcGJvYXJkJztcbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcblxuLy8gZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xuLy8gICAgIHJldHVybiBDYWxsKFwiOndhaWxzOkVudmlyb25tZW50XCIpO1xuLy8gfVxuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy53YWlscyA9IHtcbiAgICBDYWxsYmFjayxcbiAgICBjYWxsYmFja3MsXG4gICAgRXZlbnRzTm90aWZ5LFxuICAgIGV2ZW50TGlzdGVuZXJzLFxuICAgIFNldEJpbmRpbmdzLFxufTtcblxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZShpZCkge1xuICAgIHJldHVybiB7XG4gICAgICAgIC8vIExvZzogbmV3TG9nKGlkKSxcbiAgICAgICAgLy8gQnJvd3NlcjogbmV3QnJvd3NlcihpZCksXG4gICAgICAgIC8vIFNjcmVlbjogbmV3U2NyZWVuKGlkKSxcbiAgICAgICAgLy8gRXZlbnRzOiBuZXdFdmVudHMoaWQpLFxuICAgICAgICBDbGlwYm9hcmQ6IHtcbiAgICAgICAgICAgIC4uLkNsaXBib2FyZFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyhpZCksXG4gICAgICAgIFNob3c6ICgpID0+IGludm9rZShcIlNcIiksXG4gICAgICAgIEhpZGU6ICgpID0+IGludm9rZShcIkhcIiksXG4gICAgICAgIFF1aXQ6ICgpID0+IGludm9rZShcIlFcIiksXG4gICAgICAgIC8vIEdldFdpbmRvdzogZnVuY3Rpb24gKHdpbmRvd0lEKSB7XG4gICAgICAgIC8vICAgICBpZiAoIXdpbmRvd0lEKSB7XG4gICAgICAgIC8vICAgICAgICAgcmV0dXJuIHRoaXMuV2luZG93O1xuICAgICAgICAvLyAgICAgfVxuICAgICAgICAvLyAgICAgcmV0dXJuIG5ld1dpbmRvdyh3aW5kb3dJRCk7XG4gICAgICAgIC8vIH1cbiAgICB9XG59XG5cbndpbmRvdy5ydW50aW1lID0gbmV3UnVudGltZSgtMSk7XG5cbmlmIChERUJVRykge1xuICAgIGNvbnNvbGUubG9nKFwiV2FpbHMgdjMuMC4wIERlYnVnIE1vZGUgRW5hYmxlZFwiKTtcbn1cblxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7QUFhQSxNQUFJLGNBQWM7QUFFbEIsR0FBQyxXQUFZO0FBRVosUUFBSSxZQUFZLFNBQVUsR0FBRztBQUM1QixVQUFJLE1BQU0sT0FBTyxFQUFFLE1BQU0sQ0FBQztBQUMxQixhQUFPLE9BQU8sRUFBRTtBQUFRLGNBQU0sSUFBSSxFQUFFLE1BQU0sQ0FBQztBQUMzQyxhQUFPO0FBQUEsSUFDUjtBQUNBLFFBQUksVUFBVSxVQUFVLENBQUMsVUFBVSxXQUFXLGFBQWEsQ0FBQztBQUM1RCxRQUFJLFlBQVksVUFBVSxDQUFDLFVBQVUsbUJBQW1CLFlBQVksYUFBYSxDQUFDO0FBRWxGLFFBQUksQ0FBQyxXQUFXLENBQUMsV0FBVztBQUMzQixjQUFRLE1BQU0sc0JBQXNCO0FBQ3BDO0FBQUEsSUFDRDtBQUVBLFFBQUksU0FBUztBQUNaLG9CQUFjLENBQUMsWUFBWSxPQUFPLE9BQU8sUUFBUSxZQUFZLE9BQU87QUFB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
