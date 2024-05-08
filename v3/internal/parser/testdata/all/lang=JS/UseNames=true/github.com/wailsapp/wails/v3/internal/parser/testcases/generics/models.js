// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

import {Create as $Create} from "@wailsio/runtime";

/**
 * @template T
 */
export class Person {
    /**
     * Creates a new Person instance.
     * @param {Partial<Person<T>>} [$$source = {}] - The source object to create the Person.
     */
    constructor($$source = {}) {
        if (!("Name" in $$source)) {
            /**
             * @member
             * @type {T | null}
             */
            this["Name"] = null;
        }
        if (!("AliasedField" in $$source)) {
            /**
             * @member
             * @type {number}
             */
            this["AliasedField"] = 0;
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Person instance from a string or object.
     * Generic types also need creation functions for each type parameter.
     * @template T
     * @param {(any) => T} $$createT
     * @param {any} [$$source = {}]
     * @returns {Person<T>}
     */
    static createFrom($$createT, $$source = {}) {
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("Name" in $$parsedSource) {
            $$parsedSource["Name"] = $$createT($$parsedSource["Name"]);
        }
        return new Person(/** @type {Partial<Person<T>>} */($$parsedSource));
    }
}
