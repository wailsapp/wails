/*
 _	   __	  _ __
| |	 / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/* jshint esversion: 9 */

import {newRuntimeCallerWithID, objectNames} from "./runtime";

let call = newRuntimeCallerWithID(objectNames.HTTP);

let HTTPFetch = 0;

/**
 * Perform an HTTP request
 * @param {Object} options - The request options
 * @param {string} options.url - The URL to request
 * @param {string} [options.method='GET'] - The HTTP method
 * @param {Object} [options.headers] - Request headers
 * @param {string} [options.body] - Request body
 * @param {number} [options.timeout] - Request timeout in seconds
 * @returns {Promise<Object>} The response object
 */
export function Fetch(options) {
    // Ensure we have required fields
    if (!options || !options.url) {
        return Promise.reject(new Error("URL is required"));
    }

    // Set defaults
    const request = {
        url: options.url,
        method: options.method || 'GET',
        headers: options.headers || {},
        body: options.body || '',
        timeout: options.timeout || 30
    };

    // For POST requests, we need to send the body differently
    if (request.body) {
        return call(HTTPFetch, null, JSON.stringify(request));
    }

    return call(HTTPFetch, null, JSON.stringify(request));
}

/**
 * Convenience method for GET requests
 * @param {string} url - The URL to request
 * @param {Object} [options] - Additional options
 * @returns {Promise<Object>} The response object
 */
export function Get(url, options = {}) {
    return Fetch({ ...options, url, method: 'GET' });
}

/**
 * Convenience method for POST requests
 * @param {string} url - The URL to request
 * @param {string|Object} body - The request body
 * @param {Object} [options] - Additional options
 * @returns {Promise<Object>} The response object
 */
export function Post(url, body, options = {}) {
    if (typeof body === 'object' && !(body instanceof String)) {
        body = JSON.stringify(body);
        options.headers = {
            'Content-Type': 'application/json',
            ...(options.headers || {})
        };
    }
    return Fetch({ ...options, url, method: 'POST', body });
}

/**
 * Convenience method for PUT requests
 * @param {string} url - The URL to request
 * @param {string|Object} body - The request body
 * @param {Object} [options] - Additional options
 * @returns {Promise<Object>} The response object
 */
export function Put(url, body, options = {}) {
    if (typeof body === 'object' && !(body instanceof String)) {
        body = JSON.stringify(body);
        options.headers = {
            'Content-Type': 'application/json',
            ...(options.headers || {})
        };
    }
    return Fetch({ ...options, url, method: 'PUT', body });
}

/**
 * Convenience method for DELETE requests
 * @param {string} url - The URL to request
 * @param {Object} [options] - Additional options
 * @returns {Promise<Object>} The response object
 */
export function Delete(url, options = {}) {
    return Fetch({ ...options, url, method: 'DELETE' });
}

/**
 * Convenience method for PATCH requests
 * @param {string} url - The URL to request
 * @param {string|Object} body - The request body
 * @param {Object} [options] - Additional options
 * @returns {Promise<Object>} The response object
 */
export function Patch(url, body, options = {}) {
    if (typeof body === 'object' && !(body instanceof String)) {
        body = JSON.stringify(body);
        options.headers = {
            'Content-Type': 'application/json',
            ...(options.headers || {})
        };
    }
    return Fetch({ ...options, url, method: 'PATCH', body });
}

/**
 * Convenience method for HEAD requests
 * @param {string} url - The URL to request
 * @param {Object} [options] - Additional options
 * @returns {Promise<Object>} The response object
 */
export function Head(url, options = {}) {
    return Fetch({ ...options, url, method: 'HEAD' });
}