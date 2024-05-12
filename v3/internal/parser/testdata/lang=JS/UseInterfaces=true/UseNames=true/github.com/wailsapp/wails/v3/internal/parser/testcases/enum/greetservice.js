// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "/wails/runtime.js";

import * as $models from "./models.js";

/**
 * Greet does XYZ
 * @param {string} name
 * @param {$models.Title} title
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Greet(name, title) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.Greet", name, title));
    return $resultPromise;
}

/**
 * NewPerson creates a new person
 * @param {string} name
 * @returns {Promise<$models.Person | null> & { cancel(): void }}
 */
export function NewPerson(name) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.NewPerson", name));
    return $resultPromise;
}
