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
import * as $models from "./internal.js";

/**
 * Greet does XYZ
 * @param {$models.Person} person
 * @param {$models.Embedded1} emb
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Greet(person, emb) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.Greet", person, emb));
    return $resultPromise;
}
