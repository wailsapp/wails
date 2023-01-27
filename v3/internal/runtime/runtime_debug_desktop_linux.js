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
    GetText: () => GetText,
    SetText: () => SetText
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
  function GetText() {
    return call("GetText");
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9pcGMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2V2ZW50cy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3AvY2xpcGJvYXJkLmpzIiwgImRlc2t0b3AvcnVudGltZS5qcyIsICJkZXNrdG9wL3dpbmRvdy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5cblxubGV0IHBvc3RNZXNzYWdlID0gbnVsbDtcblxuKGZ1bmN0aW9uICgpIHtcblx0Ly8gQ3JlZGl0OiBodHRwczovL3N0YWNrb3ZlcmZsb3cuY29tL2EvMjYzMTUyMVxuXHRsZXQgX2RlZXB0ZXN0ID0gZnVuY3Rpb24gKHMpIHtcblx0XHRsZXQgb2JqID0gd2luZG93W3Muc2hpZnQoKV07XG5cdFx0d2hpbGUgKG9iaiAmJiBzLmxlbmd0aCkgb2JqID0gb2JqW3Muc2hpZnQoKV07XG5cdFx0cmV0dXJuIG9iajtcblx0fTtcblx0bGV0IHdpbmRvd3MgPSBfZGVlcHRlc3QoW1wiY2hyb21lXCIsIFwid2Vidmlld1wiLCBcInBvc3RNZXNzYWdlXCJdKTtcblx0bGV0IG1hY19saW51eCA9IF9kZWVwdGVzdChbXCJ3ZWJraXRcIiwgXCJtZXNzYWdlSGFuZGxlcnNcIiwgXCJleHRlcm5hbFwiLCBcInBvc3RNZXNzYWdlXCJdKTtcblxuXHRpZiAoIXdpbmRvd3MgJiYgIW1hY19saW51eCkge1xuXHRcdGNvbnNvbGUuZXJyb3IoXCJVbnN1cHBvcnRlZCBQbGF0Zm9ybVwiKTtcblx0XHRyZXR1cm47XG5cdH1cblxuXHRpZiAod2luZG93cykge1xuXHRcdHBvc3RNZXNzYWdlID0gKG1lc3NhZ2UpID0+IHdpbmRvdy5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZShtZXNzYWdlKTtcblx0fVxuXHRpZiAobWFjX2xpbnV4KSB7XG5cdFx0cG9zdE1lc3NhZ2UgPSAobWVzc2FnZSkgPT4gd2luZG93LndlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2UobWVzc2FnZSk7XG5cdH1cbn0pKCk7XG5cbmV4cG9ydCBmdW5jdGlvbiBpbnZva2UobWVzc2FnZSwgaWQpIHtcblx0aWYoIGlkICYmIGlkICE9PSAtMSkge1xuXHRcdHBvc3RNZXNzYWdlKFwiV0lORE9XSUQ6XCIrIGlkICsgXCI6XCIgKyBtZXNzYWdlKTtcblx0fSBlbHNlIHtcblx0XHRwb3N0TWVzc2FnZShtZXNzYWdlKTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5leHBvcnQgY29uc3QgY2FsbGJhY2tzID0ge307XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciBmcm9tIHRoZSBuYXRpdmUgYnJvd3NlciByYW5kb20gZnVuY3Rpb25cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gY3J5cHRvUmFuZG9tKCkge1xuXHRsZXQgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XG5cdHJldHVybiB3aW5kb3cuY3J5cHRvLmdldFJhbmRvbVZhbHVlcyhhcnJheSlbMF07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIG51bWJlciB1c2luZyBkYSBvbGQtc2tvb2wgTWF0aC5SYW5kb21cbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cbiAqXG4gKiBAcmV0dXJucyBudW1iZXJcbiAqL1xuZnVuY3Rpb24gYmFzaWNSYW5kb20oKSB7XG5cdHJldHVybiBNYXRoLnJhbmRvbSgpICogOTAwNzE5OTI1NDc0MDk5MTtcbn1cblxuLy8gUGljayBhIHJhbmRvbSBudW1iZXIgZnVuY3Rpb24gYmFzZWQgb24gYnJvd3NlciBjYXBhYmlsaXR5XG5sZXQgcmFuZG9tRnVuYztcbmlmICh3aW5kb3cuY3J5cHRvKSB7XG5cdHJhbmRvbUZ1bmMgPSBjcnlwdG9SYW5kb207XG59IGVsc2Uge1xuXHRyYW5kb21GdW5jID0gYmFzaWNSYW5kb207XG59XG5cblxuLyoqXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXG4gKiBnaXZlbiBkYXRhLiBBIHByb21pc2UgaXMgcmV0dXJuZWQgYW5kIHdpbGwgYmUgY29tcGxldGVkIHdoZW4gdGhlXG4gKiBiYWNrZW5kIHJlc3BvbmRzLiBUaGlzIHdpbGwgYmUgcmVzb2x2ZWQgd2hlbiB0aGUgY2FsbCB3YXMgc3VjY2Vzc2Z1bFxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXG4gKiBUaGVyZSBpcyBhIHRpbWVvdXQgbWVjaGFuaXNtLiBJZiB0aGUgY2FsbCBkb2Vzbid0IHJlc3BvbmQgaW4gdGhlIGdpdmVuXG4gKiB0aW1lIChpbiBtaWxsaXNlY29uZHMpIHRoZW4gdGhlIHByb21pc2UgaXMgcmVqZWN0ZWQuXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWVcbiAqIEBwYXJhbSB7YW55PX0gYXJnc1xuICogQHBhcmFtIHtudW1iZXI9fSB0aW1lb3V0XG4gKiBAcmV0dXJuc1xuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChuYW1lLCBhcmdzLCB0aW1lb3V0KSB7XG5cblx0Ly8gVGltZW91dCBpbmZpbml0ZSBieSBkZWZhdWx0XG5cdGlmICh0aW1lb3V0ID09IG51bGwpIHtcblx0XHR0aW1lb3V0ID0gMDtcblx0fVxuXG5cdGxldCB3aW5kb3dJRCA9IHdpbmRvdy53YWlscy53aW5kb3cuSUQoKTtcblxuXHQvLyBDcmVhdGUgYSBwcm9taXNlXG5cdHJldHVybiBuZXcgUHJvbWlzZShmdW5jdGlvbiAocmVzb2x2ZSwgcmVqZWN0KSB7XG5cblx0XHQvLyBDcmVhdGUgYSB1bmlxdWUgY2FsbGJhY2tJRFxuXHRcdGxldCBjYWxsYmFja0lEO1xuXHRcdGRvIHtcblx0XHRcdGNhbGxiYWNrSUQgPSBuYW1lICsgJy0nICsgcmFuZG9tRnVuYygpO1xuXHRcdH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cblx0XHRsZXQgdGltZW91dEhhbmRsZTtcblx0XHQvLyBTZXQgdGltZW91dFxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xuXHRcdFx0dGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRyZWplY3QoRXJyb3IoJ0NhbGwgdG8gJyArIG5hbWUgKyAnIHRpbWVkIG91dC4gUmVxdWVzdCBJRDogJyArIGNhbGxiYWNrSUQpKTtcblx0XHRcdH0sIHRpbWVvdXQpO1xuXHRcdH1cblxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXG5cdFx0Y2FsbGJhY2tzW2NhbGxiYWNrSURdID0ge1xuXHRcdFx0dGltZW91dEhhbmRsZTogdGltZW91dEhhbmRsZSxcblx0XHRcdHJlamVjdDogcmVqZWN0LFxuXHRcdFx0cmVzb2x2ZTogcmVzb2x2ZVxuXHRcdH07XG5cblx0XHR0cnkge1xuXHRcdFx0Y29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0bmFtZSxcblx0XHRcdFx0YXJncyxcblx0XHRcdFx0Y2FsbGJhY2tJRCxcblx0XHRcdFx0d2luZG93SUQsXG5cdFx0XHR9O1xuXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XG4gICAgICAgIH1cbiAgICB9KTtcbn1cblxud2luZG93Lk9iZnVzY2F0ZWRDYWxsID0gKGlkLCBhcmdzLCB0aW1lb3V0KSA9PiB7XG5cbiAgICAvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcbiAgICBpZiAodGltZW91dCA9PSBudWxsKSB7XG4gICAgICAgIHRpbWVvdXQgPSAwO1xuICAgIH1cblxuICAgIC8vIENyZWF0ZSBhIHByb21pc2VcbiAgICByZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xuXG4gICAgICAgIC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXG4gICAgICAgIGxldCBjYWxsYmFja0lEO1xuICAgICAgICBkbyB7XG4gICAgICAgICAgICBjYWxsYmFja0lEID0gaWQgKyAnLScgKyByYW5kb21GdW5jKCk7XG4gICAgICAgIH0gd2hpbGUgKGNhbGxiYWNrc1tjYWxsYmFja0lEXSk7XG5cbiAgICAgICAgbGV0IHRpbWVvdXRIYW5kbGU7XG4gICAgICAgIC8vIFNldCB0aW1lb3V0XG4gICAgICAgIGlmICh0aW1lb3V0ID4gMCkge1xuICAgICAgICAgICAgdGltZW91dEhhbmRsZSA9IHNldFRpbWVvdXQoZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcignQ2FsbCB0byBtZXRob2QgJyArIGlkICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XG4gICAgICAgICAgICB9LCB0aW1lb3V0KTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFN0b3JlIGNhbGxiYWNrXG4gICAgICAgIGNhbGxiYWNrc1tjYWxsYmFja0lEXSA9IHtcbiAgICAgICAgICAgIHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXG4gICAgICAgICAgICByZWplY3Q6IHJlamVjdCxcbiAgICAgICAgICAgIHJlc29sdmU6IHJlc29sdmVcbiAgICAgICAgfTtcblxuICAgICAgICB0cnkge1xuICAgICAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcblx0XHRcdFx0aWQsXG5cdFx0XHRcdGFyZ3MsXG5cdFx0XHRcdGNhbGxiYWNrSUQsXG5cdFx0XHRcdHdpbmRvd0lEOiB3aW5kb3cud2FpbHMud2luZG93LklEKCksXG5cdFx0XHR9O1xuXG4gICAgICAgICAgICAvLyBNYWtlIHRoZSBjYWxsXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ2MnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZSk7XG4gICAgICAgIH1cbiAgICB9KTtcbn07XG5cblxuLyoqXG4gKiBDYWxsZWQgYnkgdGhlIGJhY2tlbmQgdG8gcmV0dXJuIGRhdGEgdG8gYSBwcmV2aW91c2x5IGNhbGxlZFxuICogYmluZGluZyBpbnZvY2F0aW9uXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGluY29taW5nTWVzc2FnZVxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbGJhY2soaW5jb21pbmdNZXNzYWdlKSB7XG5cdC8vIFBhcnNlIHRoZSBtZXNzYWdlXG5cdGxldCBtZXNzYWdlO1xuXHR0cnkge1xuXHRcdG1lc3NhZ2UgPSBKU09OLnBhcnNlKGluY29taW5nTWVzc2FnZSk7XG5cdH0gY2F0Y2ggKGUpIHtcblx0XHRjb25zdCBlcnJvciA9IGBJbnZhbGlkIEpTT04gcGFzc2VkIHRvIGNhbGxiYWNrOiAke2UubWVzc2FnZX0uIE1lc3NhZ2U6ICR7aW5jb21pbmdNZXNzYWdlfWA7XG5cdFx0cnVudGltZS5Mb2dEZWJ1ZyhlcnJvcik7XG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcblx0fVxuXHRsZXQgY2FsbGJhY2tJRCA9IG1lc3NhZ2UuY2FsbGJhY2tpZDtcblx0bGV0IGNhbGxiYWNrRGF0YSA9IGNhbGxiYWNrc1tjYWxsYmFja0lEXTtcblx0aWYgKCFjYWxsYmFja0RhdGEpIHtcblx0XHRjb25zdCBlcnJvciA9IGBDYWxsYmFjayAnJHtjYWxsYmFja0lEfScgbm90IHJlZ2lzdGVyZWQhISFgO1xuXHRcdGNvbnNvbGUuZXJyb3IoZXJyb3IpOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lXG5cdFx0dGhyb3cgbmV3IEVycm9yKGVycm9yKTtcblx0fVxuXHRjbGVhclRpbWVvdXQoY2FsbGJhY2tEYXRhLnRpbWVvdXRIYW5kbGUpO1xuXG5cdGRlbGV0ZSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XG5cblx0aWYgKG1lc3NhZ2UuZXJyb3IpIHtcblx0XHRjYWxsYmFja0RhdGEucmVqZWN0KG1lc3NhZ2UuZXJyb3IpO1xuXHR9IGVsc2Uge1xuXHRcdGNhbGxiYWNrRGF0YS5yZXNvbHZlKG1lc3NhZ2UucmVzdWx0KTtcblx0fVxufVxuIiwgIi8qXG4gXyAgICAgICBfXyAgICAgIF8gX19cbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG4vLyBEZWZpbmVzIGEgc2luZ2xlIGxpc3RlbmVyIHdpdGggYSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0byBjYWxsYmFja1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gRXZlbnRzT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykge1xuICAgIHJldHVybiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XG5cbiAgICAvLyBHZXQgdGhlIGV2ZW50IG5hbWVcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XG5cbiAgICAvLyBDaGVjayBpZiB3ZSBoYXZlIGFueSBsaXN0ZW5lcnMgZm9yIHRoaXMgZXZlbnRcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSkge1xuXG4gICAgICAgIC8vIEtlZXAgYSBsaXN0IG9mIGxpc3RlbmVyIGluZGV4ZXMgdG8gZGVzdHJveVxuICAgICAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uc2xpY2UoKTtcblxuICAgICAgICAvLyBJdGVyYXRlIGxpc3RlbmVyc1xuICAgICAgICBmb3IgKGxldCBjb3VudCA9IDA7IGNvdW50IDwgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGg7IGNvdW50ICs9IDEpIHtcblxuICAgICAgICAgICAgLy8gR2V0IG5leHQgbGlzdGVuZXJcbiAgICAgICAgICAgIGNvbnN0IGxpc3RlbmVyID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXVtjb3VudF07XG5cbiAgICAgICAgICAgIGxldCBkYXRhID0gZXZlbnREYXRhLmRhdGE7XG5cbiAgICAgICAgICAgIC8vIERvIHRoZSBjYWxsYmFja1xuICAgICAgICAgICAgY29uc3QgZGVzdHJveSA9IGxpc3RlbmVyLkNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgaWYgKGRlc3Ryb3kpIHtcbiAgICAgICAgICAgICAgICAvLyBpZiB0aGUgbGlzdGVuZXIgaW5kaWNhdGVkIHRvIGRlc3Ryb3kgaXRzZWxmLCBhZGQgaXQgdG8gdGhlIGRlc3Ryb3kgbGlzdFxuICAgICAgICAgICAgICAgIG5ld0V2ZW50TGlzdGVuZXJMaXN0LnNwbGljZShjb3VudCwgMSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cblxuICAgICAgICAvLyBVcGRhdGUgY2FsbGJhY2tzIHdpdGggbmV3IGxpc3Qgb2YgbGlzdGVuZXJzXG4gICAgICAgIGlmIChuZXdFdmVudExpc3RlbmVyTGlzdC5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gbmV3RXZlbnRMaXN0ZW5lckxpc3Q7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogTm90aWZ5IGluZm9ybXMgZnJvbnRlbmQgbGlzdGVuZXJzIHRoYXQgYW4gZXZlbnQgd2FzIGVtaXR0ZWQgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBub3RpZnlNZXNzYWdlIC0gZW5jb2RlZCBub3RpZmljYXRpb24gbWVzc2FnZVxuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNOb3RpZnkobm90aWZ5TWVzc2FnZSkge1xuICAgIC8vIFBhcnNlIHRoZSBtZXNzYWdlXG4gICAgbGV0IG1lc3NhZ2U7XG4gICAgdHJ5IHtcbiAgICAgICAgbWVzc2FnZSA9IEpTT04ucGFyc2Uobm90aWZ5TWVzc2FnZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zdCBlcnJvciA9ICdJbnZhbGlkIEpTT04gcGFzc2VkIHRvIE5vdGlmeTogJyArIG5vdGlmeU1lc3NhZ2U7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihlcnJvcik7XG4gICAgfVxuICAgIG5vdGlmeUxpc3RlbmVycyhtZXNzYWdlKTtcbn1cblxuLyoqXG4gKiBFbWl0IGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIGRhdGFcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNFbWl0KGV2ZW50TmFtZSkge1xuXG4gICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgbmFtZTogZXZlbnROYW1lLFxuICAgICAgICBkYXRhOiBbXS5zbGljZS5hcHBseShhcmd1bWVudHMpLnNsaWNlKDEpLFxuICAgIH07XG5cbiAgICAvLyBOb3RpZnkgSlMgbGlzdGVuZXJzXG4gICAgbm90aWZ5TGlzdGVuZXJzKHBheWxvYWQpO1xuXG4gICAgLy8gTm90aWZ5IEdvIGxpc3RlbmVyc1xuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRUUnICsgSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xufVxuXG5mdW5jdGlvbiByZW1vdmVMaXN0ZW5lcihldmVudE5hbWUpIHtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJzXG4gICAgZGVsZXRlIGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV07XG5cbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXG4gICAgd2luZG93LldhaWxzSW52b2tlKCdFWCcgKyBldmVudE5hbWUpO1xufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcmVzIGNhbiBiZSB1bnJlZ2lzdGVyZWQgdmlhIGBhZGRpdGlvbmFsRXZlbnROYW1lc2BcbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcblxuICAgIGlmIChhZGRpdGlvbmFsRXZlbnROYW1lcy5sZW5ndGggPiAwKSB7XG4gICAgICAgIGFkZGl0aW9uYWxFdmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IHtcbiAgICAgICAgICAgIHJlbW92ZUxpc3RlbmVyKGV2ZW50TmFtZSlcbiAgICAgICAgfSlcbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGFsbCBldmVudCBsaXN0ZW5lcnMgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cbiAqL1xuIGV4cG9ydCBmdW5jdGlvbiBFdmVudHNPZmZBbGwoKSB7XG4gICAgY29uc3QgZXZlbnROYW1lcyA9IE9iamVjdC5rZXlzKGV2ZW50TGlzdGVuZXJzKTtcbiAgICBmb3IgKGxldCBpID0gMDsgaSAhPT0gZXZlbnROYW1lcy5sZW5ndGg7IGkrKykge1xuICAgICAgICByZW1vdmVMaXN0ZW5lcihldmVudE5hbWVzW2ldKTtcbiAgICB9XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBFdmVudHNPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbiBmdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcbiAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG5cbiAgICAvLyBDbGVhbiB1cCBpZiB0aGVyZSBhcmUgbm8gZXZlbnQgbGlzdGVuZXJzIGxlZnRcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGggPT09IDApIHtcbiAgICAgICAgcmVtb3ZlTGlzdGVuZXIoZXZlbnROYW1lKTtcbiAgICB9XG59XG4iLCAiLypcbiBfICAgICAgIF9fICAgICAgXyBfXyAgICBcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKSBcbnxfXy98X18vXFxfXyxfL18vXy9fX19fLyAgXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xuXG5pbXBvcnQge0NhbGx9IGZyb20gJy4vY2FsbHMnO1xuXG4vLyBUaGlzIGlzIHdoZXJlIHdlIGJpbmQgZ28gbWV0aG9kIHdyYXBwZXJzXG53aW5kb3cuZ28gPSB7fTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldEJpbmRpbmdzKGJpbmRpbmdzTWFwKSB7XG5cdHRyeSB7XG5cdFx0YmluZGluZ3NNYXAgPSBKU09OLnBhcnNlKGJpbmRpbmdzTWFwKTtcblx0fSBjYXRjaCAoZSkge1xuXHRcdGNvbnNvbGUuZXJyb3IoZSk7XG5cdH1cblxuXHQvLyBJbml0aWFsaXNlIHRoZSBiaW5kaW5ncyBtYXBcblx0d2luZG93LmdvID0gd2luZG93LmdvIHx8IHt9O1xuXG5cdC8vIEl0ZXJhdGUgcGFja2FnZSBuYW1lc1xuXHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcCkuZm9yRWFjaCgocGFja2FnZU5hbWUpID0+IHtcblxuXHRcdC8vIENyZWF0ZSBpbm5lciBtYXAgaWYgaXQgZG9lc24ndCBleGlzdFxuXHRcdHdpbmRvdy5nb1twYWNrYWdlTmFtZV0gPSB3aW5kb3cuZ29bcGFja2FnZU5hbWVdIHx8IHt9O1xuXG5cdFx0Ly8gSXRlcmF0ZSBzdHJ1Y3QgbmFtZXNcblx0XHRPYmplY3Qua2V5cyhiaW5kaW5nc01hcFtwYWNrYWdlTmFtZV0pLmZvckVhY2goKHN0cnVjdE5hbWUpID0+IHtcblxuXHRcdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XG5cdFx0XHR3aW5kb3cuZ29bcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdID0gd2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSB8fCB7fTtcblxuXHRcdFx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXBbcGFja2FnZU5hbWVdW3N0cnVjdE5hbWVdKS5mb3JFYWNoKChtZXRob2ROYW1lKSA9PiB7XG5cblx0XHRcdFx0d2luZG93LmdvW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXVttZXRob2ROYW1lXSA9IGZ1bmN0aW9uICgpIHtcblxuXHRcdFx0XHRcdC8vIE5vIHRpbWVvdXQgYnkgZGVmYXVsdFxuXHRcdFx0XHRcdGxldCB0aW1lb3V0ID0gMDtcblxuXHRcdFx0XHRcdC8vIEFjdHVhbCBmdW5jdGlvblxuXHRcdFx0XHRcdGZ1bmN0aW9uIGR5bmFtaWMoKSB7XG5cdFx0XHRcdFx0XHRjb25zdCBhcmdzID0gW10uc2xpY2UuY2FsbChhcmd1bWVudHMpO1xuXHRcdFx0XHRcdFx0cmV0dXJuIENhbGwoW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXS5qb2luKCcuJyksIGFyZ3MsIHRpbWVvdXQpO1xuXHRcdFx0XHRcdH1cblxuXHRcdFx0XHRcdC8vIEFsbG93IHNldHRpbmcgdGltZW91dCB0byBmdW5jdGlvblxuXHRcdFx0XHRcdGR5bmFtaWMuc2V0VGltZW91dCA9IGZ1bmN0aW9uIChuZXdUaW1lb3V0KSB7XG5cdFx0XHRcdFx0XHR0aW1lb3V0ID0gbmV3VGltZW91dDtcblx0XHRcdFx0XHR9O1xuXG5cdFx0XHRcdFx0Ly8gQWxsb3cgZ2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXG5cdFx0XHRcdFx0ZHluYW1pYy5nZXRUaW1lb3V0ID0gZnVuY3Rpb24gKCkge1xuXHRcdFx0XHRcdFx0cmV0dXJuIHRpbWVvdXQ7XG5cdFx0XHRcdFx0fTtcblxuXHRcdFx0XHRcdHJldHVybiBkeW5hbWljO1xuXHRcdFx0XHR9KCk7XG5cdFx0XHR9KTtcblx0XHR9KTtcblx0fSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gR2V0VGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChcIkdldFRleHRcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIGZvciAobGV0IGtleSBpbiBhcmdzKSB7XG4gICAgICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChrZXksIGFyZ3Nba2V5XSk7XG4gICAgICAgIH1cbiAgICB9XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZmV0Y2godXJsKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiY29udGVudC10eXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiY29udGVudC10eXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgaWQpIHtcbiAgICBpZiAoIWlkIHx8IGlkID09PSAtMSkge1xuICAgICAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncykge1xuICAgICAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCBhcmdzKTtcbiAgICAgICAgfTtcbiAgICB9XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3MpIHtcbiAgICAgICAgYXJncyA9IGFyZ3MgfHwge307XG4gICAgICAgIGFyZ3NbXCJ3aW5kb3dJRFwiXSA9IGlkO1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGwob2JqZWN0ICsgXCIuXCIgKyBtZXRob2QsIGFyZ3MpO1xuICAgIH1cbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdXaW5kb3coaWQpIHtcbiAgICBsZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJ3aW5kb3dcIiwgaWQpO1xuICAgIHJldHVybiB7XG4gICAgICAgIC8vIFJlbG9hZDogKCkgPT4gY2FsbCgnV1InKSxcbiAgICAgICAgLy8gUmVsb2FkQXBwOiAoKSA9PiBjYWxsKCdXUicpLFxuICAgICAgICAvLyBTZXRTeXN0ZW1EZWZhdWx0VGhlbWU6ICgpID0+IGNhbGwoJ1dBU0RUJyksXG4gICAgICAgIC8vIFNldExpZ2h0VGhlbWU6ICgpID0+IGNhbGwoJ1dBTFQnKSxcbiAgICAgICAgLy8gU2V0RGFya1RoZW1lOiAoKSA9PiBjYWxsKCdXQURUJyksXG4gICAgICAgIENlbnRlcjogKCkgPT4gY2FsbCgnQ2VudGVyJyksXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IGNhbGwoJ1NldFRpdGxlJywge3RpdGxlfSksXG4gICAgICAgIEZ1bGxzY3JlZW46ICgpID0+IGNhbGwoJ0Z1bGxzY3JlZW4nKSxcbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiBjYWxsKCdVbkZ1bGxzY3JlZW4nKSxcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoJ1NldFNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG4gICAgICAgIFNpemU6ICgpID0+IHsgcmV0dXJuIGNhbGwoJ1NpemUnKSB9LFxuICAgICAgICBTZXRNYXhTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbCgnU2V0TWF4U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcbiAgICAgICAgU2V0TWluU2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoJ1NldE1pblNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAoYikgPT4gY2FsbCgnU2V0QWx3YXlzT25Ub3AnLCB7YWx3YXlzT25Ub3A6Yn0pLFxuICAgICAgICBTZXRQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoJ1NldFBvc2l0aW9uJywge3gseX0pLFxuICAgICAgICBQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUG9zaXRpb24nKSB9LFxuICAgICAgICBIaWRlOiAoKSA9PiBjYWxsKCdIaWRlJyksXG4gICAgICAgIE1heGltaXNlOiAoKSA9PiBjYWxsKCdNYXhpbWlzZScpLFxuICAgICAgICBTaG93OiAoKSA9PiBjYWxsKCdTaG93JyksXG4gICAgICAgIFRvZ2dsZU1heGltaXNlOiAoKSA9PiBjYWxsKCdUb2dnbGVNYXhpbWlzZScpLFxuICAgICAgICBVbk1heGltaXNlOiAoKSA9PiBjYWxsKCdVbk1heGltaXNlJyksXG4gICAgICAgIE1pbmltaXNlOiAoKSA9PiBjYWxsKCdNaW5pbWlzZScpLFxuICAgICAgICBVbk1pbmltaXNlOiAoKSA9PiBjYWxsKCdVbk1pbmltaXNlJyksXG4gICAgICAgIFNldEJhY2tncm91bmRDb2xvdXI6IChyLCBnLCBiLCBhKSA9PiBjYWxsKCdTZXRCYWNrZ3JvdW5kQ29sb3VyJywge1IsIEcsIEIsIEF9KSxcbiAgICB9XG59XG5cbi8vIGV4cG9ydCBmdW5jdGlvbiBJc0Z1bGxzY3JlZW46ICgpPT4gLy8gICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzRnVsbHNjcmVlblwiKSxcbi8vXG5cbi8vIGV4cG9ydCBmdW5jdGlvbiBJc01heGltaXNlZDogKCk9PiAvLyAgICAgcmV0dXJuIENhbGwoXCI6d2FpbHM6V2luZG93SXNNYXhpbWlzZWRcIiksXG4vL1xuXG4vLyBleHBvcnQgZnVuY3Rpb24gSXNNaW5pbWlzZWQ6ICgpPT4gLy8gICAgIHJldHVybiBDYWxsKFwiOndhaWxzOldpbmRvd0lzTWluaW1pc2VkXCIpLFxuLy9cblxuLy8gZXhwb3J0IGZ1bmN0aW9uIElzTm9ybWFsOiAoKT0+IC8vICAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpXaW5kb3dJc05vcm1hbFwiKSxcbi8vXG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtpbnZva2V9IGZyb20gXCIuL2lwYy5qc1wiO1xuaW1wb3J0IHtDYWxsYmFjaywgY2FsbGJhY2tzfSBmcm9tICcuL2NhbGxzJztcbmltcG9ydCB7RXZlbnRzTm90aWZ5LCBldmVudExpc3RlbmVyc30gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQge1NldEJpbmRpbmdzfSBmcm9tIFwiLi9iaW5kaW5nc1wiO1xuXG5cbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tICcuL2NsaXBib2FyZCc7XG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XG5cbi8vIGV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpIHtcbi8vICAgICByZXR1cm4gQ2FsbChcIjp3YWlsczpFbnZpcm9ubWVudFwiKTtcbi8vIH1cblxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXG53aW5kb3cud2FpbHMgPSB7XG4gICAgQ2FsbGJhY2ssXG4gICAgY2FsbGJhY2tzLFxuICAgIEV2ZW50c05vdGlmeSxcbiAgICBldmVudExpc3RlbmVycyxcbiAgICBTZXRCaW5kaW5ncyxcbn07XG5cblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWUoaWQpIHtcbiAgICByZXR1cm4ge1xuICAgICAgICAvLyBMb2c6IG5ld0xvZyhpZCksXG4gICAgICAgIC8vIEJyb3dzZXI6IG5ld0Jyb3dzZXIoaWQpLFxuICAgICAgICAvLyBTY3JlZW46IG5ld1NjcmVlbihpZCksXG4gICAgICAgIC8vIEV2ZW50czogbmV3RXZlbnRzKGlkKSxcbiAgICAgICAgQ2xpcGJvYXJkOiB7XG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcbiAgICAgICAgfSxcbiAgICAgICAgV2luZG93OiBuZXdXaW5kb3coaWQpLFxuICAgICAgICBTaG93OiAoKSA9PiBpbnZva2UoXCJTXCIpLFxuICAgICAgICBIaWRlOiAoKSA9PiBpbnZva2UoXCJIXCIpLFxuICAgICAgICBRdWl0OiAoKSA9PiBpbnZva2UoXCJRXCIpLFxuICAgICAgICAvLyBHZXRXaW5kb3c6IGZ1bmN0aW9uICh3aW5kb3dJRCkge1xuICAgICAgICAvLyAgICAgaWYgKCF3aW5kb3dJRCkge1xuICAgICAgICAvLyAgICAgICAgIHJldHVybiB0aGlzLldpbmRvdztcbiAgICAgICAgLy8gICAgIH1cbiAgICAgICAgLy8gICAgIHJldHVybiBuZXdXaW5kb3cod2luZG93SUQpO1xuICAgICAgICAvLyB9XG4gICAgfVxufVxuXG53aW5kb3cucnVudGltZSA9IG5ld1J1bnRpbWUoLTEpO1xuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBYUEsTUFBSSxjQUFjO0FBRWxCLEdBQUMsV0FBWTtBQUVaLFFBQUksWUFBWSxTQUFVLEdBQUc7QUFDNUIsVUFBSSxNQUFNLE9BQU8sRUFBRSxNQUFNLENBQUM7QUFDMUIsYUFBTyxPQUFPLEVBQUU7QUFBUSxjQUFNLElBQUksRUFBRSxNQUFNLENBQUM7QUFDM0MsYUFBTztBQUFBLElBQ1I7QUFDQSxRQUFJLFVBQVUsVUFBVSxDQUFDLFVBQVUsV0FBVyxhQUFhLENBQUM7QUFDNUQsUUFBSSxZQUFZLFVBQVUsQ0FBQyxVQUFVLG1CQUFtQixZQUFZLGFBQWEsQ0FBQztBQUVsRixRQUFJLENBQUMsV0FBVyxDQUFDLFdBQVc7QUFDM0IsY0FBUSxNQUFNLHNCQUFzQjtBQUNwQztBQUFBLElBQ0Q7QUFFQSxRQUFJLFNBQVM7QUFDWixvQkFBYyxDQUFDLFlBQVksT0FBTyxPQUFPLFFBQVEsWUFBWSxPQUFPO0FBQUEsSUFDckU7QUFDQSxRQUFJLFdBQVc7QUFDZCxvQkFBYyxDQUFDLFlBQVks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
