// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "/wails/runtime.js";

/**
 * Greet someone
 * @param {string} name
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Greet(name) {
    let $resultPromise = $Call.ByID(1411160069, name);
    return /** @type {any} */($resultPromise);
}
