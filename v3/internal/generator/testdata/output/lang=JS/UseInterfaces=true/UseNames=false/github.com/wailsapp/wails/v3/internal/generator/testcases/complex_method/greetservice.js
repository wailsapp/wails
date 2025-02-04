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
import * as $models from "./models.js";

/**
 * Greet does XYZ
 * It has a multiline doc comment
 * The comment has even some * / traps!!
 * @param {string} str
 * @param {$models.Person[] | null} people
 * @param {{"AnotherCount": number, "AnotherOne": $models.Person | null}} $2
 * @param {{ [_: `${number}`]: boolean | null } | null} assoc
 * @param {(number | null)[] | null} $4
 * @param {string[]} other
 * @returns {Promise<[$models.Person, any, number[] | null]> & { cancel(): void }}
 */
export function Greet(str, people, $2, assoc, $4, ...other) {
    let $resultPromise = /** @type {any} */($Call.ByID(1411160069, str, people, $2, assoc, $4, other));
    return $resultPromise;
}
