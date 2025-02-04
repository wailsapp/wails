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
import * as nobindingshere$0 from "../no_bindings_here/models.js";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * Get someone.
 * @param {$models.Alias} aliasValue
 * @returns {Promise<$models.Person> & { cancel(): void }}
 */
export function Get(aliasValue) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.Get", aliasValue));
    return $resultPromise;
}

/**
 * Apparently, aliases are all the rage right now.
 * @param {$models.AliasedPerson} p
 * @returns {Promise<$models.StrangelyAliasedPerson> & { cancel(): void }}
 */
export function GetButAliased(p) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.GetButAliased", p));
    return $resultPromise;
}

/**
 * Get someone quite different.
 * @returns {Promise<$models.GenericPerson<boolean>> & { cancel(): void }}
 */
export function GetButDifferent() {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.GetButDifferent"));
    return $resultPromise;
}

/**
 * @returns {Promise<nobindingshere$0.PrivatePerson> & { cancel(): void }}
 */
export function GetButForeignPrivateAlias() {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.GetButForeignPrivateAlias"));
    return $resultPromise;
}

/**
 * @returns {Promise<$models.AliasGroup> & { cancel(): void }}
 */
export function GetButGenericAliases() {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.GetButGenericAliases"));
    return $resultPromise;
}

/**
 * Greet a lot of unusual things.
 * @param {$models.EmptyAliasStruct} $0
 * @param {$models.EmptyStruct} $1
 * @returns {Promise<$models.AliasStruct> & { cancel(): void }}
 */
export function Greet($0, $1) {
    let $resultPromise = /** @type {any} */($Call.ByName("main.GreetService.Greet", $0, $1));
    return $resultPromise;
}
