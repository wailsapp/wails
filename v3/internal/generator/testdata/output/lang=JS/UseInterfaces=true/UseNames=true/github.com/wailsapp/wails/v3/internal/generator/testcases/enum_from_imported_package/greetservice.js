// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "/wails/runtime.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as services$0 from "./services/models.js";

/**
 * Greet does XYZ
 * @param {string} name
 * @param {services$0.Title} title
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Greet(name, title) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.Greet", name, title));
    return $resultPromise;
}
