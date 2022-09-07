/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */

export const callbacks = {};

/**
 * Returns a number from the native browser random function
 *
 * @returns number
 */
function cryptoRandom() {
	var array = new Uint32Array(1);
	return window.crypto.getRandomValues(array)[0];
}

/**
 * Returns a number using da old-skool Math.Random
 * I likes to call it LOLRandom
 *
 * @returns number
 */
function basicRandom() {
	return Math.random() * 9007199254740991;
}

// Pick a random number function based on browser capability
var randomFunc;
if (window.crypto) {
	randomFunc = cryptoRandom;
} else {
	randomFunc = basicRandom;
}


/**
 * Call sends a message to the backend to call the binding with the
 * given data. A promise is returned and will be completed when the
 * backend responds. This will be resolved when the call was successful
 * or rejected if an error is passed back.
 * There is a timeout mechanism. If the call doesn't respond in the given
 * time (in milliseconds) then the promise is rejected.
 *
 * @export
 * @param {string} name
 * @param {any=} args
 * @param {number=} timeout
 * @returns
 */
export function Call(name, args, timeout) {

	// Timeout infinite by default
	if (timeout == null) {
		timeout = 0;
	}

	// Create a promise
	return new Promise(function (resolve, reject) {

		// Create a unique callbackID
		var callbackID;
		do {
			callbackID = name + '-' + randomFunc();
		} while (callbacks[callbackID]);

		var timeoutHandle;
		// Set timeout
		if (timeout > 0) {
			timeoutHandle = setTimeout(function () {
				reject(Error('Call to ' + name + ' timed out. Request ID: ' + callbackID));
			}, timeout);
		}

		// Store callback
		callbacks[callbackID] = {
			timeoutHandle: timeoutHandle,
			reject: reject,
			resolve: resolve
		};

		try {
			const payload = {
				name,
				args,
				callbackID,
			};

            // Make the call
            window.WailsInvoke('C' + JSON.stringify(payload));
        } catch (e) {
            // eslint-disable-next-line
            console.error(e);
        }
    });
}

window.ObfuscatedCall = (id, args, timeout) => {

    // Timeout infinite by default
    if (timeout == null) {
        timeout = 0;
    }

    // Create a promise
    return new Promise(function (resolve, reject) {

        // Create a unique callbackID
        var callbackID;
        do {
            callbackID = id + '-' + randomFunc();
        } while (callbacks[callbackID]);

        var timeoutHandle;
        // Set timeout
        if (timeout > 0) {
            timeoutHandle = setTimeout(function () {
                reject(Error('Call to method ' + id + ' timed out. Request ID: ' + callbackID));
            }, timeout);
        }

        // Store callback
        callbacks[callbackID] = {
            timeoutHandle: timeoutHandle,
            reject: reject,
            resolve: resolve
        };

        try {
            const payload = {
				id,
				args,
				callbackID,
			};

            // Make the call
            window.WailsInvoke('c' + JSON.stringify(payload));
        } catch (e) {
            // eslint-disable-next-line
            console.error(e);
        }
    });
};


/**
 * Called by the backend to return data to a previously called
 * binding invocation
 *
 * @export
 * @param {string} incomingMessage
 */
export function Callback(incomingMessage) {
	// Parse the message
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
		console.error(error); // eslint-disable-line
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
