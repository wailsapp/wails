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
 * @typedef {import("./models.js").Title} Title
 */

/**
 * Greet does XYZ
 * @param {string} name
 * @param {Title} title
 * @returns {Promise<string>}
 */
export function Greet(name, title) {
    let $resultPromise = $Call.ByID(1411160069, name, title);
    return /** @type {any} */($resultPromise);
}

/**
 * NewPerson creates a new person
 * @param {string} name
 * @returns {Promise<Person | null>}
 */
export function NewPerson(name) {
    let $resultPromise = $Call.ByID(1661412647, name);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType0($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

// Internal type creation functions
const $$createType0 = $Create.Nullable(Person.createFrom);
