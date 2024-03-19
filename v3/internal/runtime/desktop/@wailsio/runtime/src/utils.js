/*
 _     __     _ __
| |  / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

/**
 * Logs a message to the console with custom formatting.
 * @param {string} message - The message to be logged.
 * @return {void}
 */
export function debugLog(message) {
    // eslint-disable-next-line
    console.log(
        '%c wails3 %c ' + message + ' ',
        'background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem',
        'background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem'
    );
}

/**
 * Checks whether the browser supports removing listeners by triggering an AbortSignal
 * (see https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener#signal)
 *
 * @return {boolean}
 */
export function canAbortListeners() {
    if (!EventTarget || !AbortSignal || !AbortController)
        return false;

    let result = true;

    const target = new EventTarget();
    const controller = new AbortController();
    target.addEventListener('test', () => { result = false; }, { signal: controller.signal });
    controller.abort();
    target.dispatchEvent(new CustomEvent('test'));

    return result;
}

/***
 This technique for proper load detection is taken from HTMX:

 BSD 2-Clause License

 Copyright (c) 2020, Big Sky Software
 All rights reserved.

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice, this
 list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
 this list of conditions and the following disclaimer in the documentation
 and/or other materials provided with the distribution.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

 ***/

let isReady = false;
document.addEventListener('DOMContentLoaded', () => isReady = true);

export function whenReady(callback) {
    if (isReady || document.readyState === 'complete') {
        callback();
    } else {
        document.addEventListener('DOMContentLoaded', callback);
    }
}
