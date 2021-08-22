(() => {
    var __defProp = Object.defineProperty;
    var __markAsModule = (target) => __defProp(target, "__esModule", {value: true});
    var __export = (target, all) => {
        __markAsModule(target);
        for (var name in all)
            __defProp(target, name, {get: all[name], enumerable: true});
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
        constructor(callback, maxCallbacks) {
            maxCallbacks = maxCallbacks || -1;
            this.Callback = (data) => {
                callback.apply(null, data);
                if (maxCallbacks === -1) {
                    return false;
                }
                maxCallbacks -= 1;
                return maxCallbacks === 0;
            };
        }
    };
    var eventListeners = {};

    function EventsOnMultiple(eventName, callback, maxCallbacks) {
        eventListeners[eventName] = eventListeners[eventName] || [];
        const thisListener = new Listener(callback, maxCallbacks);
        eventListeners[eventName].push(thisListener);
    }

    function EventsOn(eventName, callback) {
        EventsOnMultiple(eventName, callback, -1);
    }

    function EventsOnce(eventName, callback) {
        EventsOnMultiple(eventName, callback, 1);
    }

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
            eventListeners[eventName] = newEventListenerList;
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

    function EventsOff(eventName) {
        eventListeners.delete(eventName);
        window.WailsInvoke("EX" + eventName);
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
        return new Promise(function (resolve, reject) {
            var callbackID;
            do {
                callbackID = name + "-" + randomFunc();
            } while (callbacks[callbackID]);
            var timeoutHandle;
            if (timeout > 0) {
                timeoutHandle = setTimeout(function () {
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

    function Callback(incomingMessage) {
        var message;
        try {
            message = JSON.parse(incomingMessage);
        } catch (e) {
            const error = `Invalid JSON passed to callback: ${e.message}. Message: ${incomingMessage}`;
            wails.LogDebug(error);
            throw new Error(error);
        }
        var callbackID = message.callbackid;
        var callbackData = callbacks[callbackID];
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
    window.backend = {};

    function SetBindings(bindingsMap) {
        try {
            bindingsMap = JSON.parse(bindingsMap);
        } catch (e) {
            console.error(e);
        }
        window.backend = window.backend || {};
        Object.keys(bindingsMap).forEach((packageName) => {
            window.backend[packageName] = window.backend[packageName] || {};
            Object.keys(bindingsMap[packageName]).forEach((structName) => {
                window.backend[packageName][structName] = window.backend[packageName][structName] || {};
                Object.keys(bindingsMap[packageName][structName]).forEach((methodName) => {
                    window.backend[packageName][structName][methodName] = function () {
                        let timeout = 0;

                        function dynamic() {
                            const args = [].slice.call(arguments);
                            return Call([packageName, structName, methodName].join("."), args, timeout);
                        }

                        dynamic.setTimeout = function (newTimeout) {
                            timeout = newTimeout;
                        };
                        dynamic.getTimeout = function () {
                            return timeout;
                        };
                        return dynamic;
                    }();
                });
            });
        });
    }

    // desktop/main.js
    window.backend = {};
    window.runtime = {
        ...log_exports,
        EventsOn,
        EventsOnce,
        EventsOnMultiple,
        EventsEmit,
        EventsOff
    };
    window.wails = {
        Callback,
        EventsNotify,
        SetBindings,
        eventListeners,
        callbacks
    };
    window.wails.SetBindings(window.wailsbindings);
    delete window.wails.SetBindings;
    delete window.wailsbindings;
    window.addEventListener("mousedown", (e) => {
        let currentElement = e.target;
        while (currentElement != null) {
            if (currentElement.hasAttribute("data-wails-no-drag")) {
                break;
            } else if (currentElement.hasAttribute("data-wails-drag")) {
                window.WailsInvoke("drag");
                break;
            }
            currentElement = currentElement.parentElement;
        }
    });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9ldmVudHMuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL2JpbmRpbmdzLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcclxuIF8gICAgICAgX18gICAgICBfIF9fXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBsaWdodHdlaWdodCBmcmFtZXdvcmsgZm9yIHdlYi1saWtlIGFwcHNcclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogNiAqL1xyXG5cclxuLyoqXHJcbiAqIFNlbmRzIGEgbG9nIG1lc3NhZ2UgdG8gdGhlIGJhY2tlbmQgd2l0aCB0aGUgZ2l2ZW4gbGV2ZWwgKyBtZXNzYWdlXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBsZXZlbFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZnVuY3Rpb24gc2VuZExvZ01lc3NhZ2UobGV2ZWwsIG1lc3NhZ2UpIHtcclxuXHJcblx0Ly8gTG9nIE1lc3NhZ2UgZm9ybWF0OlxyXG5cdC8vIGxbdHlwZV1bbWVzc2FnZV1cclxuXHR3aW5kb3cuV2FpbHNJbnZva2UoJ0wnICsgbGV2ZWwgKyBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gdHJhY2UgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dUcmFjZShtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1QnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dQcmludChtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ1AnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gZGVidWcgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dEZWJ1ZyhtZXNzYWdlKSB7XHJcblx0c2VuZExvZ01lc3NhZ2UoJ0QnLCBtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIExvZyB0aGUgZ2l2ZW4gaW5mbyBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ0luZm8obWVzc2FnZSkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdJJywgbWVzc2FnZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBMb2cgdGhlIGdpdmVuIHdhcm5pbmcgbWVzc2FnZSB3aXRoIHRoZSBiYWNrZW5kXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2VcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBMb2dXYXJuaW5nKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnVycsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBlcnJvciBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ0Vycm9yKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnRScsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogTG9nIHRoZSBnaXZlbiBmYXRhbCBtZXNzYWdlIHdpdGggdGhlIGJhY2tlbmRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIExvZ0ZhdGFsKG1lc3NhZ2UpIHtcclxuXHRzZW5kTG9nTWVzc2FnZSgnRicsIG1lc3NhZ2UpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgTG9nIGxldmVsIHRvIHRoZSBnaXZlbiBsb2cgbGV2ZWxcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge251bWJlcn0gbG9nbGV2ZWxcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRMb2dMZXZlbChsb2dsZXZlbCkge1xyXG5cdHNlbmRMb2dNZXNzYWdlKCdTJywgbG9nbGV2ZWwpO1xyXG59XHJcblxyXG4vLyBMb2cgbGV2ZWxzXHJcbmV4cG9ydCBjb25zdCBMb2dMZXZlbCA9IHtcclxuXHRUUkFDRTogMSxcclxuXHRERUJVRzogMixcclxuXHRJTkZPOiAzLFxyXG5cdFdBUk5JTkc6IDQsXHJcblx0RVJST1I6IDUsXHJcbn07XHJcbiIsICIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX19cclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGxpZ2h0d2VpZ2h0IGZyYW1ld29yayBmb3Igd2ViLWxpa2UgYXBwc1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbi8vIERlZmluZXMgYSBzaW5nbGUgbGlzdGVuZXIgd2l0aCBhIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRvIGNhbGxiYWNrXHJcblxyXG4vKipcclxuICogVGhlIExpc3RlbmVyIGNsYXNzIGRlZmluZXMgYSBsaXN0ZW5lciEgOi0pXHJcbiAqXHJcbiAqIEBjbGFzcyBMaXN0ZW5lclxyXG4gKi9cclxuY2xhc3MgTGlzdGVuZXIge1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIExpc3RlbmVyLlxyXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcclxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcclxuICAgICAqIEBtZW1iZXJvZiBMaXN0ZW5lclxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XHJcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxyXG4gICAgICAgIG1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXHJcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxyXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xyXG4gICAgICAgICAgICBjYWxsYmFjay5hcHBseShudWxsLCBkYXRhKTtcclxuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxyXG4gICAgICAgICAgICBpZiAobWF4Q2FsbGJhY2tzID09PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcclxuICAgICAgICAgICAgbWF4Q2FsbGJhY2tzIC09IDE7XHJcbiAgICAgICAgICAgIHJldHVybiBtYXhDYWxsYmFja3MgPT09IDA7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0ge307XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGBtYXhDYWxsYmFja3NgIHRpbWVzIGJlZm9yZSBiZWluZyBkZXN0cm95ZWRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0gfHwgW107XHJcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XHJcbiAgICBldmVudExpc3RlbmVyc1tldmVudE5hbWVdLnB1c2godGhpc0xpc3RlbmVyKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIEV2ZW50c09uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXZlbnRzT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICBFdmVudHNPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBub3RpZnlMaXN0ZW5lcnMoZXZlbnREYXRhKSB7XHJcblxyXG4gICAgLy8gR2V0IHRoZSBldmVudCBuYW1lXHJcbiAgICBsZXQgZXZlbnROYW1lID0gZXZlbnREYXRhLm5hbWU7XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgd2UgaGF2ZSBhbnkgbGlzdGVuZXJzIGZvciB0aGlzIGV2ZW50XHJcbiAgICBpZiAoZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSkge1xyXG5cclxuICAgICAgICAvLyBLZWVwIGEgbGlzdCBvZiBsaXN0ZW5lciBpbmRleGVzIHRvIGRlc3Ryb3lcclxuICAgICAgICBjb25zdCBuZXdFdmVudExpc3RlbmVyTGlzdCA9IGV2ZW50TGlzdGVuZXJzW2V2ZW50TmFtZV0uc2xpY2UoKTtcclxuXHJcbiAgICAgICAgLy8gSXRlcmF0ZSBsaXN0ZW5lcnNcclxuICAgICAgICBmb3IgKGxldCBjb3VudCA9IDA7IGNvdW50IDwgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXS5sZW5ndGg7IGNvdW50ICs9IDEpIHtcclxuXHJcbiAgICAgICAgICAgIC8vIEdldCBuZXh0IGxpc3RlbmVyXHJcbiAgICAgICAgICAgIGNvbnN0IGxpc3RlbmVyID0gZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXVtjb3VudF07XHJcblxyXG4gICAgICAgICAgICBsZXQgZGF0YSA9IGV2ZW50RGF0YS5kYXRhO1xyXG5cclxuICAgICAgICAgICAgLy8gRG8gdGhlIGNhbGxiYWNrXHJcbiAgICAgICAgICAgIGNvbnN0IGRlc3Ryb3kgPSBsaXN0ZW5lci5DYWxsYmFjayhkYXRhKTtcclxuICAgICAgICAgICAgaWYgKGRlc3Ryb3kpIHtcclxuICAgICAgICAgICAgICAgIC8vIGlmIHRoZSBsaXN0ZW5lciBpbmRpY2F0ZWQgdG8gZGVzdHJveSBpdHNlbGYsIGFkZCBpdCB0byB0aGUgZGVzdHJveSBsaXN0XHJcbiAgICAgICAgICAgICAgICBuZXdFdmVudExpc3RlbmVyTGlzdC5zcGxpY2UoY291bnQsIDEpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBVcGRhdGUgY2FsbGJhY2tzIHdpdGggbmV3IGxpc3Qgb2YgbGlzdGVuZXJzXHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnNbZXZlbnROYW1lXSA9IG5ld0V2ZW50TGlzdGVuZXJMaXN0O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogTm90aWZ5IGluZm9ybXMgZnJvbnRlbmQgbGlzdGVuZXJzIHRoYXQgYW4gZXZlbnQgd2FzIGVtaXR0ZWQgd2l0aCB0aGUgZ2l2ZW4gZGF0YVxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBub3RpZnlNZXNzYWdlIC0gZW5jb2RlZCBub3RpZmljYXRpb24gbWVzc2FnZVxyXG5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNOb3RpZnkobm90aWZ5TWVzc2FnZSkge1xyXG4gICAgLy8gUGFyc2UgdGhlIG1lc3NhZ2VcclxuICAgIGxldCBtZXNzYWdlO1xyXG4gICAgdHJ5IHtcclxuICAgICAgICBtZXNzYWdlID0gSlNPTi5wYXJzZShub3RpZnlNZXNzYWdlKTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICBjb25zdCBlcnJvciA9ICdJbnZhbGlkIEpTT04gcGFzc2VkIHRvIE5vdGlmeTogJyArIG5vdGlmeU1lc3NhZ2U7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGVycm9yKTtcclxuICAgIH1cclxuICAgIG5vdGlmeUxpc3RlbmVycyhtZXNzYWdlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEVtaXQgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgZGF0YVxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFdmVudHNFbWl0KGV2ZW50TmFtZSkge1xyXG5cclxuICAgIGNvbnN0IHBheWxvYWQgPSB7XHJcbiAgICAgICAgbmFtZTogZXZlbnROYW1lLFxyXG4gICAgICAgIGRhdGE6IFtdLnNsaWNlLmFwcGx5KGFyZ3VtZW50cykuc2xpY2UoMSksXHJcbiAgICB9O1xyXG5cclxuICAgIC8vIE5vdGlmeSBKUyBsaXN0ZW5lcnNcclxuICAgIG5vdGlmeUxpc3RlbmVycyhwYXlsb2FkKTtcclxuXHJcbiAgICAvLyBOb3RpZnkgR28gbGlzdGVuZXJzXHJcbiAgICB3aW5kb3cuV2FpbHNJbnZva2UoJ0VFJyArIEpTT04uc3RyaW5naWZ5KHBheWxvYWQpKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIEV2ZW50c09mZihldmVudE5hbWUpIHtcclxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lcnNcclxuICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xyXG5cclxuICAgIC8vIE5vdGlmeSBHbyBsaXN0ZW5lcnNcclxuICAgIHdpbmRvdy5XYWlsc0ludm9rZSgnRVgnICsgZXZlbnROYW1lKTtcclxufSIsICIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX19cclxufCB8ICAgICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGxpZ2h0d2VpZ2h0IGZyYW1ld29yayBmb3Igd2ViLWxpa2UgYXBwc1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDYgKi9cclxuXHJcbmV4cG9ydCBjb25zdCBjYWxsYmFja3MgPSB7fTtcclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgbnVtYmVyIGZyb20gdGhlIG5hdGl2ZSBicm93c2VyIHJhbmRvbSBmdW5jdGlvblxyXG4gKlxyXG4gKiBAcmV0dXJucyBudW1iZXJcclxuICovXHJcbmZ1bmN0aW9uIGNyeXB0b1JhbmRvbSgpIHtcclxuXHR2YXIgYXJyYXkgPSBuZXcgVWludDMyQXJyYXkoMSk7XHJcblx0cmV0dXJuIHdpbmRvdy5jcnlwdG8uZ2V0UmFuZG9tVmFsdWVzKGFycmF5KVswXTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBudW1iZXIgdXNpbmcgZGEgb2xkLXNrb29sIE1hdGguUmFuZG9tXHJcbiAqIEkgbGlrZXMgdG8gY2FsbCBpdCBMT0xSYW5kb21cclxuICpcclxuICogQHJldHVybnMgbnVtYmVyXHJcbiAqL1xyXG5mdW5jdGlvbiBiYXNpY1JhbmRvbSgpIHtcclxuXHRyZXR1cm4gTWF0aC5yYW5kb20oKSAqIDkwMDcxOTkyNTQ3NDA5OTE7XHJcbn1cclxuXHJcbi8vIFBpY2sgYSByYW5kb20gbnVtYmVyIGZ1bmN0aW9uIGJhc2VkIG9uIGJyb3dzZXIgY2FwYWJpbGl0eVxyXG52YXIgcmFuZG9tRnVuYztcclxuaWYgKHdpbmRvdy5jcnlwdG8pIHtcclxuXHRyYW5kb21GdW5jID0gY3J5cHRvUmFuZG9tO1xyXG59IGVsc2Uge1xyXG5cdHJhbmRvbUZ1bmMgPSBiYXNpY1JhbmRvbTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBDYWxsIHNlbmRzIGEgbWVzc2FnZSB0byB0aGUgYmFja2VuZCB0byBjYWxsIHRoZSBiaW5kaW5nIHdpdGggdGhlXHJcbiAqIGdpdmVuIGRhdGEuIEEgcHJvbWlzZSBpcyByZXR1cm5lZCBhbmQgd2lsbCBiZSBjb21wbGV0ZWQgd2hlbiB0aGVcclxuICogYmFja2VuZCByZXNwb25kcy4gVGhpcyB3aWxsIGJlIHJlc29sdmVkIHdoZW4gdGhlIGNhbGwgd2FzIHN1Y2Nlc3NmdWxcclxuICogb3IgcmVqZWN0ZWQgaWYgYW4gZXJyb3IgaXMgcGFzc2VkIGJhY2suXHJcbiAqIFRoZXJlIGlzIGEgdGltZW91dCBtZWNoYW5pc20uIElmIHRoZSBjYWxsIGRvZXNuJ3QgcmVzcG9uZCBpbiB0aGUgZ2l2ZW5cclxuICogdGltZSAoaW4gbWlsbGlzZWNvbmRzKSB0aGVuIHRoZSBwcm9taXNlIGlzIHJlamVjdGVkLlxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBhcmdzXHJcbiAqIEBwYXJhbSB7bnVtYmVyPX0gdGltZW91dFxyXG4gKiBAcmV0dXJuc1xyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwobmFtZSwgYXJncywgdGltZW91dCkge1xyXG5cclxuXHQvLyBUaW1lb3V0IGluZmluaXRlIGJ5IGRlZmF1bHRcclxuXHRpZiAodGltZW91dCA9PSBudWxsKSB7XHJcblx0XHR0aW1lb3V0ID0gMDtcclxuXHR9XHJcblxyXG5cdC8vIENyZWF0ZSBhIHByb21pc2VcclxuXHRyZXR1cm4gbmV3IFByb21pc2UoZnVuY3Rpb24gKHJlc29sdmUsIHJlamVjdCkge1xyXG5cclxuXHRcdC8vIENyZWF0ZSBhIHVuaXF1ZSBjYWxsYmFja0lEXHJcblx0XHR2YXIgY2FsbGJhY2tJRDtcclxuXHRcdGRvIHtcclxuXHRcdFx0Y2FsbGJhY2tJRCA9IG5hbWUgKyAnLScgKyByYW5kb21GdW5jKCk7XHJcblx0XHR9IHdoaWxlIChjYWxsYmFja3NbY2FsbGJhY2tJRF0pO1xyXG5cclxuXHRcdHZhciB0aW1lb3V0SGFuZGxlO1xyXG5cdFx0Ly8gU2V0IHRpbWVvdXRcclxuXHRcdGlmICh0aW1lb3V0ID4gMCkge1xyXG5cdFx0XHR0aW1lb3V0SGFuZGxlID0gc2V0VGltZW91dChmdW5jdGlvbiAoKSB7XHJcblx0XHRcdFx0cmVqZWN0KEVycm9yKCdDYWxsIHRvICcgKyBuYW1lICsgJyB0aW1lZCBvdXQuIFJlcXVlc3QgSUQ6ICcgKyBjYWxsYmFja0lEKSk7XHJcblx0XHRcdH0sIHRpbWVvdXQpO1xyXG5cdFx0fVxyXG5cclxuXHRcdC8vIFN0b3JlIGNhbGxiYWNrXHJcblx0XHRjYWxsYmFja3NbY2FsbGJhY2tJRF0gPSB7XHJcblx0XHRcdHRpbWVvdXRIYW5kbGU6IHRpbWVvdXRIYW5kbGUsXHJcblx0XHRcdHJlamVjdDogcmVqZWN0LFxyXG5cdFx0XHRyZXNvbHZlOiByZXNvbHZlXHJcblx0XHR9O1xyXG5cclxuXHRcdHRyeSB7XHJcblx0XHRcdGNvbnN0IHBheWxvYWQgPSB7XHJcblx0XHRcdFx0bmFtZSxcclxuXHRcdFx0XHRhcmdzLFxyXG5cdFx0XHRcdGNhbGxiYWNrSUQsXHJcblx0XHRcdH07XHJcblxyXG5cdFx0XHQvLyBNYWtlIHRoZSBjYWxsXHJcblx0XHRcdHdpbmRvdy5XYWlsc0ludm9rZSgnQycgKyBKU09OLnN0cmluZ2lmeShwYXlsb2FkKSk7XHJcblx0XHR9IGNhdGNoIChlKSB7XHJcblx0XHRcdC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG5cdFx0XHRjb25zb2xlLmVycm9yKGUpO1xyXG5cdFx0fVxyXG5cdH0pO1xyXG59XHJcblxyXG5cclxuXHJcbi8qKlxyXG4gKiBDYWxsZWQgYnkgdGhlIGJhY2tlbmQgdG8gcmV0dXJuIGRhdGEgdG8gYSBwcmV2aW91c2x5IGNhbGxlZFxyXG4gKiBiaW5kaW5nIGludm9jYXRpb25cclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gaW5jb21pbmdNZXNzYWdlXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbGJhY2soaW5jb21pbmdNZXNzYWdlKSB7XHJcblx0Ly8gRGVjb2RlIHRoZSBtZXNzYWdlIC0gQ3JlZGl0OiBodHRwczovL3N0YWNrb3ZlcmZsb3cuY29tL2EvMTM4NjU2ODBcclxuXHQvL2luY29taW5nTWVzc2FnZSA9IGRlY29kZVVSSUNvbXBvbmVudChpbmNvbWluZ01lc3NhZ2UucmVwbGFjZSgvXFxzKy9nLCAnJykucmVwbGFjZSgvWzAtOWEtZl17Mn0vZywgJyUkJicpKTtcclxuXHJcblx0Ly8gUGFyc2UgdGhlIG1lc3NhZ2VcclxuXHR2YXIgbWVzc2FnZTtcclxuXHR0cnkge1xyXG5cdFx0bWVzc2FnZSA9IEpTT04ucGFyc2UoaW5jb21pbmdNZXNzYWdlKTtcclxuXHR9IGNhdGNoIChlKSB7XHJcblx0XHRjb25zdCBlcnJvciA9IGBJbnZhbGlkIEpTT04gcGFzc2VkIHRvIGNhbGxiYWNrOiAke2UubWVzc2FnZX0uIE1lc3NhZ2U6ICR7aW5jb21pbmdNZXNzYWdlfWA7XHJcblx0XHR3YWlscy5Mb2dEZWJ1ZyhlcnJvcik7XHJcblx0XHR0aHJvdyBuZXcgRXJyb3IoZXJyb3IpO1xyXG5cdH1cclxuXHR2YXIgY2FsbGJhY2tJRCA9IG1lc3NhZ2UuY2FsbGJhY2tpZDtcclxuXHR2YXIgY2FsbGJhY2tEYXRhID0gY2FsbGJhY2tzW2NhbGxiYWNrSURdO1xyXG5cdGlmICghY2FsbGJhY2tEYXRhKSB7XHJcblx0XHRjb25zdCBlcnJvciA9IGBDYWxsYmFjayAnJHtjYWxsYmFja0lEfScgbm90IHJlZ2lzdGVyZWQhISFgO1xyXG5cdFx0Y29uc29sZS5lcnJvcihlcnJvcik7IC8vIGVzbGludC1kaXNhYmxlLWxpbmVcclxuXHRcdHRocm93IG5ldyBFcnJvcihlcnJvcik7XHJcblx0fVxyXG5cdGNsZWFyVGltZW91dChjYWxsYmFja0RhdGEudGltZW91dEhhbmRsZSk7XHJcblxyXG5cdGRlbGV0ZSBjYWxsYmFja3NbY2FsbGJhY2tJRF07XHJcblxyXG5cdGlmIChtZXNzYWdlLmVycm9yKSB7XHJcblx0XHRjYWxsYmFja0RhdGEucmVqZWN0KG1lc3NhZ2UuZXJyb3IpO1xyXG5cdH0gZWxzZSB7XHJcblx0XHRjYWxsYmFja0RhdGEucmVzb2x2ZShtZXNzYWdlLnJlc3VsdCk7XHJcblx0fVxyXG59XHJcbiIsICIvKlxyXG4gXyAgICAgICBfXyAgICAgIF8gX18gICAgXHJcbnwgfCAgICAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gICkgXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fLyAgXHJcblRoZSBsaWdodHdlaWdodCBmcmFtZXdvcmsgZm9yIHdlYi1saWtlIGFwcHNcclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA2ICovXHJcblxyXG5pbXBvcnQgeyBDYWxsIH0gZnJvbSAnLi9jYWxscyc7XHJcblxyXG53aW5kb3cuYmFja2VuZCA9IHt9O1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIFNldEJpbmRpbmdzKGJpbmRpbmdzTWFwKSB7XHJcblx0dHJ5IHtcclxuXHRcdGJpbmRpbmdzTWFwID0gSlNPTi5wYXJzZShiaW5kaW5nc01hcCk7XHJcblx0fSBjYXRjaCAoZSkge1xyXG5cdFx0Y29uc29sZS5lcnJvcihlKTtcclxuXHR9XHJcblxyXG5cdC8vIEluaXRpYWxpc2UgdGhlIGJhY2tlbmQgbWFwXHJcblx0d2luZG93LmJhY2tlbmQgPSB3aW5kb3cuYmFja2VuZCB8fCB7fTtcclxuXHJcblx0Ly8gSXRlcmF0ZSBwYWNrYWdlIG5hbWVzXHJcblx0T2JqZWN0LmtleXMoYmluZGluZ3NNYXApLmZvckVhY2goKHBhY2thZ2VOYW1lKSA9PiB7XHJcblxyXG5cdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XHJcblx0XHR3aW5kb3cuYmFja2VuZFtwYWNrYWdlTmFtZV0gPSB3aW5kb3cuYmFja2VuZFtwYWNrYWdlTmFtZV0gfHwge307XHJcblxyXG5cdFx0Ly8gSXRlcmF0ZSBzdHJ1Y3QgbmFtZXNcclxuXHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXSkuZm9yRWFjaCgoc3RydWN0TmFtZSkgPT4ge1xyXG5cclxuXHRcdFx0Ly8gQ3JlYXRlIGlubmVyIG1hcCBpZiBpdCBkb2Vzbid0IGV4aXN0XHJcblx0XHRcdHdpbmRvdy5iYWNrZW5kW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSA9IHdpbmRvdy5iYWNrZW5kW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSB8fCB7fTtcclxuXHJcblx0XHRcdE9iamVjdC5rZXlzKGJpbmRpbmdzTWFwW3BhY2thZ2VOYW1lXVtzdHJ1Y3ROYW1lXSkuZm9yRWFjaCgobWV0aG9kTmFtZSkgPT4ge1xyXG5cclxuXHRcdFx0XHR3aW5kb3cuYmFja2VuZFtwYWNrYWdlTmFtZV1bc3RydWN0TmFtZV1bbWV0aG9kTmFtZV0gPSBmdW5jdGlvbiAoKSB7XHJcblxyXG5cdFx0XHRcdFx0Ly8gTm8gdGltZW91dCBieSBkZWZhdWx0XHJcblx0XHRcdFx0XHRsZXQgdGltZW91dCA9IDA7XHJcblxyXG5cdFx0XHRcdFx0Ly8gQWN0dWFsIGZ1bmN0aW9uXHJcblx0XHRcdFx0XHRmdW5jdGlvbiBkeW5hbWljKCkge1xyXG5cdFx0XHRcdFx0XHRjb25zdCBhcmdzID0gW10uc2xpY2UuY2FsbChhcmd1bWVudHMpO1xyXG5cdFx0XHRcdFx0XHRyZXR1cm4gQ2FsbChbcGFja2FnZU5hbWUsIHN0cnVjdE5hbWUsIG1ldGhvZE5hbWVdLmpvaW4oJy4nKSwgYXJncywgdGltZW91dCk7XHJcblx0XHRcdFx0XHR9XHJcblxyXG5cdFx0XHRcdFx0Ly8gQWxsb3cgc2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXHJcblx0XHRcdFx0XHRkeW5hbWljLnNldFRpbWVvdXQgPSBmdW5jdGlvbiAobmV3VGltZW91dCkge1xyXG5cdFx0XHRcdFx0XHR0aW1lb3V0ID0gbmV3VGltZW91dDtcclxuXHRcdFx0XHRcdH07XHJcblxyXG5cdFx0XHRcdFx0Ly8gQWxsb3cgZ2V0dGluZyB0aW1lb3V0IHRvIGZ1bmN0aW9uXHJcblx0XHRcdFx0XHRkeW5hbWljLmdldFRpbWVvdXQgPSBmdW5jdGlvbiAoKSB7XHJcblx0XHRcdFx0XHRcdHJldHVybiB0aW1lb3V0O1xyXG5cdFx0XHRcdFx0fTtcclxuXHJcblx0XHRcdFx0XHRyZXR1cm4gZHluYW1pYztcclxuXHRcdFx0XHR9KCk7XHJcblx0XHRcdH0pO1xyXG5cdFx0fSk7XHJcblx0fSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgbGlnaHR3ZWlnaHQgZnJhbWV3b3JrIGZvciB3ZWItbGlrZSBhcHBzXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xyXG5pbXBvcnQge1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMsXHJcbiAgICBFdmVudHNFbWl0LFxyXG4gICAgRXZlbnRzTm90aWZ5LFxyXG4gICAgRXZlbnRzT2ZmLFxyXG4gICAgRXZlbnRzT24sXHJcbiAgICBFdmVudHNPbmNlLFxyXG4gICAgRXZlbnRzT25NdWx0aXBsZVxyXG59IGZyb20gJy4vZXZlbnRzJztcclxuaW1wb3J0IHtDYWxsYmFjaywgY2FsbGJhY2tzfSBmcm9tICcuL2NhbGxzJztcclxuaW1wb3J0IHtTZXRCaW5kaW5nc30gZnJvbSBcIi4vYmluZGluZ3NcIjtcclxuXHJcbi8vIEJhY2tlbmQgaXMgd2hlcmUgdGhlIEdvIHN0cnVjdCB3cmFwcGVycyBnZXQgYm91bmQgdG9cclxud2luZG93LmJhY2tlbmQgPSB7fTtcclxuXHJcbndpbmRvdy5ydW50aW1lID0ge1xyXG4gICAgLi4uTG9nLFxyXG4gICAgRXZlbnRzT24sXHJcbiAgICBFdmVudHNPbmNlLFxyXG4gICAgRXZlbnRzT25NdWx0aXBsZSxcclxuICAgIEV2ZW50c0VtaXQsXHJcbiAgICBFdmVudHNPZmYsXHJcbn07XHJcblxyXG4vLyBJbml0aWFsaXNlIGdsb2JhbCBpZiBub3QgYWxyZWFkeVxyXG53aW5kb3cud2FpbHMgPSB7XHJcbiAgICBDYWxsYmFjayxcclxuICAgIEV2ZW50c05vdGlmeSxcclxuICAgIFNldEJpbmRpbmdzLFxyXG4gICAgZXZlbnRMaXN0ZW5lcnMsXHJcbiAgICBjYWxsYmFja3NcclxufTtcclxuXHJcbndpbmRvdy53YWlscy5TZXRCaW5kaW5ncyh3aW5kb3cud2FpbHNiaW5kaW5ncyk7XHJcbmRlbGV0ZSB3aW5kb3cud2FpbHMuU2V0QmluZGluZ3M7XHJcbmRlbGV0ZSB3aW5kb3cud2FpbHNiaW5kaW5ncztcclxuXHJcbi8vIFNldHVwIGRyYWcgaGFuZGxlclxyXG4vLyBCYXNlZCBvbiBjb2RlIGZyb206IGh0dHBzOi8vZ2l0aHViLmNvbS9wYXRyMG51cy9EZXNrR2FwXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCAoZSkgPT4ge1xyXG4gICAgbGV0IGN1cnJlbnRFbGVtZW50ID0gZS50YXJnZXQ7XHJcbiAgICB3aGlsZSAoY3VycmVudEVsZW1lbnQgIT0gbnVsbCkge1xyXG4gICAgICAgIGlmIChjdXJyZW50RWxlbWVudC5oYXNBdHRyaWJ1dGUoJ2RhdGEtd2FpbHMtbm8tZHJhZycpKSB7XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIH0gZWxzZSBpZiAoY3VycmVudEVsZW1lbnQuaGFzQXR0cmlidXRlKCdkYXRhLXdhaWxzLWRyYWcnKSkge1xyXG4gICAgICAgICAgICB3aW5kb3cuV2FpbHNJbnZva2UoXCJkcmFnXCIpO1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcbiAgICAgICAgY3VycmVudEVsZW1lbnQgPSBjdXJyZW50RWxlbWVudC5wYXJlbnRFbGVtZW50O1xyXG4gICAgfVxyXG59KTtcclxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWtCQSwwQkFBd0IsT0FBTyxTQUFTO0FBSXZDLFdBQU8sWUFBWSxNQUFNLFFBQVE7QUFBQTtBQVMzQixvQkFBa0IsU0FBUztBQUNqQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCxvQkFBa0IsU0FBUztBQUNqQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCxvQkFBa0IsU0FBUztBQUNqQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCxtQkFBaUIsU0FBUztBQUNoQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCxzQkFBb0IsU0FBUztBQUNuQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCxvQkFBa0IsU0FBUztBQUNqQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCxvQkFBa0IsU0FBUztBQUNqQyxtQkFBZSxLQUFLO0FBQUE7QUFTZCx1QkFBcUIsVUFBVTtBQUNyQyxtQkFBZSxLQUFLO0FBQUE7QUFJZCxNQUFNLFdBQVc7QUFBQSxJQUN2QixPQUFPO0FBQUEsSUFDUCxPQUFPO0FBQUEsSUFDUCxNQUFNO0FBQUEsSUFDTixTQUFTO0FBQUEsSUFDVCxPQUFPO0FBQUE7OztBQzdGUix1QkFBZTtBQUFBLElBT1gsWUFBWSxVQUFVLGNBQWM7QUFFaEMscUJBQWUsZ0JBQWdCO0FBRy9CLFdBQUssV0FBVyxDQUFDLFNBQVM7QUFDdEIsaUJBQVMsTUFBTSxNQUFNO0FBRXJCLFlBQUksaUJBQWlCLElBQUk7QUFDckIsaUJBQU87QUFBQTtBQUdYLHdCQUFnQjtBQUNoQixlQUFPLGlCQUFpQjtBQUFBO0FBQUE7QUFBQTtBQUs3QixNQUFNLGlCQUFpQjtBQVV2Qiw0QkFBMEIsV0FBVyxVQUFVLGNBQWM7QUFDaEUsbUJBQWUsYUFBYSxlQUFlLGNBQWM7QUFDekQsVUFBTSxlQUFl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bUJBQWEsUUFBUSxRQUFRO0FBQUE7QUFBQTs7O0FDOUgvQixTQUFPLFVBQVU7QUFFVix1QkFBcUIsYUFBYTtBQUN4QyxRQUFJO0FBQ0gsb0JBQWMsS0FBSyxNQUFNO0FBQUEsYUFDakIsR0FBUDtBQUNELGNBQVEsTUFBTTtBQUFBO0FBSWYsV0FBTyxVQUFVLE9BQU8sV0FBVztBQUduQyxXQUFPLEtBQUssYUFBYSxRQUFRLENBQUMsZ0JBQWdCO0FBR2pELGFBQU8sUUFBUSxlQUFlLE9BQU8sUUFBUSxnQkFBZ0I7QUFHN0QsYUFBTyxLQUFLLFlBQVksY0FBYyxRQUFRLENBQUMsZUFBZTtBQUc3RCxlQUFPLFFBQVEsYUFBYSxjQUFjLE9BQU8sUUFBUSxhQUFhLGVBQWU7QUFFckYsZUFBTyxLQUFLLFlBQVksYUFBYSxhQUFhLFFBQVEsQ0FBQyxlQUFlO0FBRXpFLGlCQUFPLFFBQVEsYUFBYSxZQUFZLGNBQWMsV0FBWTtBQUdqRSxnQkFBSSxVQUFVO0FBR2QsK0JBQW1CO0FBQ2xCLG9CQUFNLE9BQU8sR0FBRyxNQUFNLEtBQUs7QUFDM0IscUJBQU8sS0FBSyxDQUFDLGFBQWEsWUFBWSxZQUFZLEtBQUssTUFBTSxNQUFNO0FBQUE7QUFJcEUsb0JBQVEsYUFBYSxTQUFVLFlBQVk7QUFDMUMsd0JBQVU7QUFBQTtBQUlYLG9CQUFRLGFBQWEsV0FBWTtBQUNoQyxxQkFBTztBQUFBO0FBR1IsbUJBQU87QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNwQ1osU0FBTyxVQUFVO0FBRWpCLFNBQU8sVUFBVTtBQUFBLE9BQ1Y7QUFBQSxJQUNIO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBO0FBSUosU0FBTyxRQUFRO0FBQUEsSUFDWDtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQTtBQUdKLFNBQU8sTUFBTSxZQUFZLE9BQU87QUFDaEMsU0FBTyxPQUFPLE1BQU07QUFDcEIsU0FBTyxPQUFPO0FBSWQsU0FBTyxpQkFBaUIsYUFBYSxDQUFDLE1BQU07QUFDeEMsUUFBSSxpQkFBaUIsRUFBRTtBQUN2QixXQUFPLGtCQUFrQixNQUFNO0FBQzNCLFVBQUksZUFBZSxhQUFhLHVCQUF1QjtBQUNuRDtBQUFBLGlCQUNPLGVBQWUsYUFBYSxvQkFBb0I7QUFDdkQsZUFBTyxZQUFZO0FBQ25CO0FBQUE7QUFFSix1QkFBaUIsZUFBZTtBQUFBO0FBQUE7IiwKICAibmFtZXMiOiBbXQp9Cg==
