// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import {Embedded1, Person} from "./models.js";

/**
 * Greet does XYZ
 * @param {Person} person
 * @param {Embedded1} emb
 * @returns {Promise<string>}
 */
export function Greet(person, emb) {
    let $resultPromise = $Call.ByID(1411160069, person, emb);
    return /** @type {any} */($resultPromise);
}
