// Wails runtime JS
(function () {
	window.wails = window.wails || {};
	window.backend = {};

	/****************** Utility Functions ************************/

	// -------------- Random --------------
	// AwesomeRandom
	function cryptoRandom() {
		var array = new Uint32Array(1);
		return window.crypto.getRandomValues(array)[0];
	}


	// LOLRandom
	function basicRandom() {
		return Math.random() * 9007199254740991;
	}

	// Pick one based on browser capability
	var randomFunc;
	if (window.crypto) {
		randomFunc = cryptoRandom;
	} else {
		randomFunc = basicRandom;
	}

	// -------------- Identifiers ---------------

	function isValidIdentifier(name) {
		// Don't xss yourself :-)
		try {
			new Function('var ' + name);
			return true;
		} catch (e) {
			return false;
		}
	}

	// -------------- JS ----------------
	function addScript(js, callbackID) {
		var script = document.createElement('script');
		script.text = js;
		document.body.appendChild(script);
		window.wails.Events.Emit(callbackID);
	}

	// -------------- CSS ---------------
	// Adapted from webview - thanks zserge!
	function injectCSS(css) {
		var elem = document.createElement('style');
		elem.setAttribute('type', 'text/css');
		if (elem.styleSheet) {
			elem.styleSheet.cssText = css;
		} else {
			elem.appendChild(document.createTextNode(css));
		}
		var head = document.head || document.getElementsByTagName('head')[0];
		head.appendChild(elem);
	}

	/************************* Bindings *************************/

	var bindingsBasePath = window.backend;

	// Creates the path given in the bindings path
	function addBindingPath(pathSections) {
		// Start at the base path
		var currentPath = bindingsBasePath;
		// for each section of the given path
		for (var sectionIndex in pathSections) {

			var section = pathSections[sectionIndex];

			// Is section a valid javascript identifier?
			if (!isValidIdentifier(section)) {
				var errMessage = section + ' is not a valid javascript identifier.';
				var err = new Error(errMessage);
				return [null, err];
			}

			// Add if doesn't exist
			if (!currentPath[section]) {
				currentPath[section] = {};
			}
			// update current path to new path
			currentPath = currentPath[section];
		}
		return [currentPath, null];
	}

	function newBinding(bindingName) {

		// Get all the sections of the binding
		var bindingSections = bindingName.split('.').splice(1);

		// Get the actual function/method call name
		var callName = bindingSections.pop();

		// Add path to binding
		var bs = addBindingPath(bindingSections);
		var pathToBinding = bs[0];
		var err = bs[1];

		if (err != null) {
			// We need to return an error
			return err;
		}

		// Add binding call
		pathToBinding[callName] = function () {

			// No timeout by default
			var timeout = 0;

			// Actual function
			function dynamic() {
				var args = [].slice.call(arguments);
				return call(bindingName, args, timeout);
			}

			// Allow setting timeout to function
			dynamic.setTimeout = function (newTimeout) {
				timeout = newTimeout;
			};

			// Allow getting timeout to function
			dynamic.getTimeout = function () {
				return timeout;
			};

			return dynamic;
		}();
	}

	/************************************************************/

	/*************************** Calls **************************/

	var callbacks = {};

	// Call sends a message to the backend to call the binding with the
	// given data. A promise is returned and will be completed when the
	// backend responds. This will be resolved when the call was successful
	// or rejected if an error is passed back.
	// There is a timeout mechanism. If the call doesn't respond in the given
	// time (in milliseconds) then the promise is rejected.

	function call(bindingName, data, timeout) {

		// Timeout infinite by default
		if (timeout == null || timeout == undefined) {
			timeout = 0;
		}

		// Create a promise
		return new Promise(function (resolve, reject) {

			// Create a unique callbackID
			var callbackID;
			do {
				callbackID = bindingName + '-' + randomFunc();
			} while (callbacks[callbackID]);

			// Set timeout
			if (timeout > 0) {
				var timeoutHandle = setTimeout(function () {
					reject(Error('Call to ' + bindingName + ' timed out. Request ID: ' + callbackID));
				}, timeout);
			}

			// Store callback
			callbacks[callbackID] = {
				timeoutHandle: timeoutHandle,
				reject: reject,
				resolve: resolve
			};
			try {
				var payloaddata = JSON.stringify(data);
				// Create the message
				var message = {
					type: 'call',
					callbackid: callbackID,
					payload: {
						bindingName: bindingName,
						data: payloaddata,
					}
				};

				// Make the call
				var payload = JSON.stringify(message);
				external.invoke(payload);
			} catch (e) {
				// eslint-disable-next-line
				console.error(e);
			}
		});
	}

	// systemCall is used to call wails methods from the frontend
	function systemCall(method, data) {
		return call('.wails.' + method, data);
	}

	// Called by the backend to return data to a previously called
	// binding invocation
	function callback(incomingMessage) {

		// Decode the message - Credit: https://stackoverflow.com/a/13865680
		incomingMessage = decodeURIComponent(incomingMessage.replace(/\s+/g, '').replace(/[0-9a-f]{2}/g, '%$&'));

		// Parse the message
		var message;
		try {
			message = JSON.parse(incomingMessage);
		} catch (e) {
			window.wails.Log.Debug('Invalid JSON passed to callback: ' + e.message);
			window.wails.Log.Debug('Message: ' + incomingMessage);
			return;
		}
		var callbackID = message.callbackid;
		var callbackData = callbacks[callbackID];
		if (!callbackData) {
			// eslint-disable-next-line
			console.error('Callback \'' + callbackID + '\' not registed!!!');
			return;
		}
		clearTimeout(callbackData.timeoutHandle);
		delete callbacks[callbackID];
		if (message.error) {
			return callbackData.reject(message.error);
		}
		return callbackData.resolve(message.data);
	}

	/************************************************************/


	/************************** Events **************************/

	var eventListeners = {};

	// Registers event listeners
	function on(eventName, callback) {
		eventListeners[eventName] = eventListeners[eventName] || [];
		eventListeners[eventName].push(callback);
	}


	// notify informs frontend listeners that an event was emitted with the given data
	function notify(eventName, data) {
		if (eventListeners[eventName]) {
			eventListeners[eventName].forEach(function (element) {
				var parsedData = [];
				// Parse data if we have it
				if (data) {
					try {
						parsedData = JSON.parse(data);
					} catch (e) {
						window.wails.Log.Error('Invalid JSON data sent to notify. Event name = ' + eventName);
					}
				}
				element.apply(null, parsedData);
			});
		}
	}

	// emit an event with the given name and data
	function emit(eventName) {

		// Calculate the data
		var data = JSON.stringify([].slice.apply(arguments).slice(1));

		// Notify backend
		var message = {
			type: 'event',
			payload: {
				name: eventName,
				data: data,
			}
		};
		external.invoke(JSON.stringify(message));
	}

	function deprecatedEventsFunction(fn, oldName) {
		var newName = oldName[0].toUpperCase() + oldName.substring(1);
		return function (eventName, eventData) {
			// eslint-disable-next-line
			console.warn('Method events.' + oldName + ' has been deprecated. Please use Events.' + newName);
			return fn(eventName, eventData);
		};
	}

	// Deprecated Events calls
	window.wails.events = {
		emit: deprecatedEventsFunction(emit, 'emit'),
		on: deprecatedEventsFunction(on, 'on'),
	};

	// Events calls
	window.wails.Events = {
		Emit: emit,
		On: on
	};



	/************************************************************/

	/************************* Browser **************************/


	function OpenURL(url) {
		return systemCall('Browser.OpenURL', url);
	}

	function OpenFile(filename) {
		return systemCall('Browser.OpenFile', filename);
	}

	window.wails.Browser = {
		OpenURL,
		OpenFile,
	};

	/************************* Logging **************************/

	// Sends a log message to the backend with the given
	// level + message
	function sendLogMessage(level, message) {

		// Log Message
		message = {
			type: 'log',
			payload: {
				level: level,
				message: message,
			}
		};
		external.invoke(JSON.stringify(message));
	}

	function deprecatedLogFunction(fn, oldName) {
		var newName = oldName[0].toUpperCase() + oldName.substring(1);
		return function (message) {
			// eslint-disable-next-line
			console.warn('Method Log.' + oldName + ' has been deprecated. Please use Log.' + newName);
			return fn(message);
		};
	}

	function logDebug(message) {
		sendLogMessage('debug', message);
	}
	function logInfo(message) {
		sendLogMessage('info', message);
	}
	function logWarning(message) {
		sendLogMessage('warning', message);
	}
	function logError(message) {
		sendLogMessage('error', message);
	}
	function logFatal(message) {
		sendLogMessage('fatal', message);
	}

	window.wails.log = {
		debug: deprecatedLogFunction(logDebug, 'debug'),
		info: deprecatedLogFunction(logInfo, 'info'),
		warning: deprecatedLogFunction(logWarning, 'warning'),
		error: deprecatedLogFunction(logError, 'error'),
		fatal: deprecatedLogFunction(logFatal, 'fatal'),
	};

	window.wails.Log = {
		Debug: logDebug,
		Info: logInfo,
		Warning: logWarning,
		Error: logError,
		Fatal: logFatal,
	};


	/************************** Exports *************************/

	window.wails._ = {
		newBinding: newBinding,
		callback: callback,
		notify: notify,
		sendLogMessage: sendLogMessage,
		callbacks: callbacks,
		injectCSS: injectCSS,
		addScript: addScript,
	};


	/************************************************************/

	// Notify backend that the runtime has finished loading
	window.wails.Events.Emit('wails:loaded');

})();