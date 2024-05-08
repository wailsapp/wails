// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

/**
 * GreetService is great
 * @module
 */

import {Call as $Call, Create as $Create} from "@wailsio/runtime";

import * as $models from "./models.js";

/**
 * Greet does XYZ
 * @param {string} name
 * @param {$models.Title} title
 * @returns {Promise<string> & { cancel(): void }}
 */
export function Greet(name, title) {
    let $resultPromise = $Call.ByID(1411160069, name, title);
    return /** @type {any} */($resultPromise);
}

/**
 * NewPerson creates a new person
 * @param {string} name
 * @returns {Promise<$models.Person | null> & { cancel(): void }}
 */
export function NewPerson(name) {
    let $resultPromise = $Call.ByID(1661412647, name);
    let $typingPromise = $resultPromise.then(($result) => {
        return $$createType1($result);
    });
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return /** @type {any} */($typingPromise);
}

// Private type creation functions
const $$createType0 = $models.Person.createFrom;
const $$createType1 = $Create.Nullable($$createType0);