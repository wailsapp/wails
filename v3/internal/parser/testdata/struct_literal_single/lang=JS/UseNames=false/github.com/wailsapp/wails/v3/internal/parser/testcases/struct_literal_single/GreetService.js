// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import {Person} from "./models.js";

/**
 * @param {number[]} $in
 * @returns {Promise<void>}
 */
export function ArrayInt($in) {
    let $resultPromise = $Call.ByID(3862002418, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {boolean} $in
 * @returns {Promise<boolean>}
 */
export function BoolInBoolOut($in) {
    let $resultPromise = $Call.ByID(2424639793, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function Float32InFloat32Out($in) {
    let $resultPromise = $Call.ByID(3132595881, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function Float64InFloat64Out($in) {
    let $resultPromise = $Call.ByID(2182412247, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * Greet someone
 * @param {string} name
 * @returns {Promise<string>}
 */
export function Greet(name) {
    let $resultPromise = $Call.ByID(1411160069, name);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function Int16InIntOut($in) {
    let $resultPromise = $Call.ByID(3306292566, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function Int16PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(1754277916, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function Int32InIntOut($in) {
    let $resultPromise = $Call.ByID(1909469092, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function Int32PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(4251088558, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function Int64InIntOut($in) {
    let $resultPromise = $Call.ByID(1343888303, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function Int64PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(2205561041, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function Int8InIntOut($in) {
    let $resultPromise = $Call.ByID(572240879, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function Int8PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(2189402897, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function IntInIntOut($in) {
    let $resultPromise = $Call.ByID(642881729, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function IntPointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(1066151743, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function IntPointerInputNamedOutputs($in) {
    let $resultPromise = $Call.ByID(2718999663, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {{ [_: `${number}`]: number }} $in
 * @returns {Promise<void>}
 */
export function MapIntInt($in) {
    let $resultPromise = $Call.ByID(2386486356, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {{ [_: string]: number }} $in
 * @returns {Promise<void>}
 */
export function MapIntPointerInt($in) {
    let $resultPromise = $Call.ByID(550413585, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {{ [_: `${number}`]: number[] }} $in
 * @returns {Promise<void>}
 */
export function MapIntSliceInt($in) {
    let $resultPromise = $Call.ByID(2900172572, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {{ [_: `${number}`]: number[] }} $in
 * @returns {Promise<{ [_: `${number}`]: number[] }>}
 */
export function MapIntSliceIntInMapIntSliceIntOut($in) {
    let $resultPromise = $Call.ByID(881980169, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * @returns {Promise<string>}
 */
export function NoInputsStringOut() {
    let $resultPromise = $Call.ByID(1075577233);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {boolean | null} $in
 * @returns {Promise<boolean | null>}
 */
export function PointerBoolInBoolOut($in) {
    let $resultPromise = $Call.ByID(3589606958, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function PointerFloat32InFloat32Out($in) {
    let $resultPromise = $Call.ByID(224675106, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function PointerFloat64InFloat64Out($in) {
    let $resultPromise = $Call.ByID(2124953624, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {{ [_: `${number}`]: number } | null} $in
 * @returns {Promise<void>}
 */
export function PointerMapIntInt($in) {
    let $resultPromise = $Call.ByID(3516977899, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {string | null} $in
 * @returns {Promise<string | null>}
 */
export function PointerStringInStringOut($in) {
    let $resultPromise = $Call.ByID(229603958, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {string[]} $in
 * @returns {Promise<string[]>}
 */
export function StringArrayInputNamedOutput($in) {
    let $resultPromise = $Call.ByID(3678582682, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * @param {string[]} $in
 * @returns {Promise<string[]>}
 */
export function StringArrayInputNamedOutputs($in) {
    let $resultPromise = $Call.ByID(319259595, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * @param {string[]} $in
 * @returns {Promise<string[]>}
 */
export function StringArrayInputStringArrayOut($in) {
    let $resultPromise = $Call.ByID(383995060, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType2($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * @param {string[]} $in
 * @returns {Promise<string>}
 */
export function StringArrayInputStringOut($in) {
    let $resultPromise = $Call.ByID(1091960237, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {Person} $in
 * @returns {Promise<Person>}
 */
export function StructInputStructOutput($in) {
    let $resultPromise = $Call.ByID(3835643147, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return Person.createFrom($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * @param {Person | null} $in
 * @returns {Promise<void>}
 */
export function StructPointerInputErrorOutput($in) {
    let $resultPromise = $Call.ByID(2447692557, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {Person | null} $in
 * @returns {Promise<Person | null>}
 */
export function StructPointerInputStructPointerOutput($in) {
    let $resultPromise = $Call.ByID(2943477349, $in);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType3($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function UInt16InUIntOut($in) {
    let $resultPromise = $Call.ByID(3401034892, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function UInt16PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(1236957573, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function UInt32InUIntOut($in) {
    let $resultPromise = $Call.ByID(1160383782, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function UInt32PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(1739300671, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function UInt64InUIntOut($in) {
    let $resultPromise = $Call.ByID(793803239, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function UInt64PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(1403757716, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function UInt8InUIntOut($in) {
    let $resultPromise = $Call.ByID(2988345717, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function UInt8PointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(518250834, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number} $in
 * @returns {Promise<number>}
 */
export function UIntInUIntOut($in) {
    let $resultPromise = $Call.ByID(2836661285, $in);
    return /** @type {any} */($resultPromise);
}

/**
 * @param {number | null} $in
 * @returns {Promise<number | null>}
 */
export function UIntPointerInAndOutput($in) {
    let $resultPromise = $Call.ByID(1367187362, $in);
    return /** @type {any} */($resultPromise);
}

// Internal type creation functions
const $$createType0 = $Create.Array($Create.Any);
const $$createType1 = $Create.Map($Create.Any, $$createType0);
const $$createType2 = $Create.Array($Create.Any);
const $$createType3 = $Create.Nullable(Person.createFrom);
