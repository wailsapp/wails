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
 * Get someone.
 * @param {$models.Alias} aliasValue
 * @returns {Promise<$models.Person> & { cancel(): void }}
 */
export function Get(aliasValue) {
    let $resultPromise = /** @type {any} */($Call.ByID(1928502664, aliasValue));
    return $resultPromise;
}

/**
 * Apparently, aliases are all the rage right now.
 * @param {$models.AliasedPerson} p
 * @returns {Promise<$models.StrangelyAliasedPerson> & { cancel(): void }}
 */
export function GetButAliased(p) {
    let $resultPromise = /** @type {any} */($Call.ByID(1896499664, p));
    return $resultPromise;
}

/**
 * Get someone quite different.
 * @returns {Promise<$models.GenericPerson<boolean>> & { cancel(): void }}
 */
export function GetButDifferent() {
    let $resultPromise = /** @type {any} */($Call.ByID(2240931744));
    return $resultPromise;
}

/**
 * Greet a lot of unusual things.
 * @param {$models.EmptyAliasStruct} $0
 * @param {$models.EmptyStruct} $1
 * @returns {Promise<$models.AliasStruct> & { cancel(): void }}
 */
export function Greet($0, $1) {
    let $resultPromise = /** @type {any} */($Call.ByID(1411160069, $0, $1));
    return $resultPromise;
}
